// Copyright Contributors to the Open Cluster Management project
package enable

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
	"open-cluster-management.io/clusteradm/pkg/cmd/addon/enable/scenario"
	"open-cluster-management.io/clusteradm/pkg/helpers"
	"open-cluster-management.io/clusteradm/pkg/helpers/apply"
)

// const appMgrAddonName = "application-manager"

type ClusterAddonInfo struct {
	ClusterName string
	NameSpace   string
	AddonName   string
}

func NewClusterAddonInfo(cn string, ns string, an string) ClusterAddonInfo {
	return ClusterAddonInfo{
		ClusterName: cn,
		NameSpace:   ns,
		AddonName:   an,
	}
}

func (o *Options) complete(cmd *cobra.Command, args []string) (err error) {
	klog.V(1).InfoS("enable options:", "dry-run", o.ClusteradmFlags.DryRun, "names", o.Names, "clusters", o.Clusters, "output-file", o.OutputFile)

	return nil
}

func (o *Options) Validate() error {
	if len(o.Names) == 0 {
		return fmt.Errorf("names is missing")
	}

	if len(o.Clusters) == 0 {
		return fmt.Errorf("clusters is misisng")
	}

	return nil
}

func (o *Options) Run() error {
	addons := sets.NewString(o.Names...)
	clusters := sets.NewString(o.Clusters...)

	klog.V(3).InfoS("values:", "addon", addons, "clusters", clusters)

	restConfig, err := o.ClusteradmFlags.KubectlFactory.ToRESTConfig()
	if err != nil {
		return err
	}
	clusterClient, err := clusterclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	kubeClient, apiExtensionsClient, dynamicClient, err := helpers.GetClients(o.ClusteradmFlags.KubectlFactory)
	if err != nil {
		return err
	}

	return o.runWithClient(clusterClient, kubeClient, apiExtensionsClient, dynamicClient, o.ClusteradmFlags.DryRun, addons.List(), clusters.List())
}

func (o *Options) runWithClient(clusterClient clusterclientset.Interface,
	kubeClient kubernetes.Interface,
	apiExtensionsClient apiextensionsclient.Interface,
	dynamicClient dynamic.Interface,
	dryRun bool,
	addons []string,
	clusters []string) error {

	for _, clusterName := range clusters {
		_, err := clusterClient.ClusterV1().ManagedClusters().Get(context.TODO(),
			clusterName,
			metav1.GetOptions{})
		if err != nil {
			return err
		}
	}

	output := make([]string, 0)
	reader := scenario.GetScenarioResourcesReader()

	applierBuilder := apply.NewApplierBuilder()
	applier := applierBuilder.WithClient(kubeClient, apiExtensionsClient, dynamicClient).Build()

	for _, addon := range addons {
		for _, clusterName := range clusters {
			cai := NewClusterAddonInfo(clusterName, o.Namespace, addon)
			out, err := applier.ApplyCustomResources(reader, cai, dryRun, "", "addons/app/addon.yaml")
			if err != nil {
				return err
			}
			output = append(output, out...)

			fmt.Fprintf(o.Streams.Out, "Deploying %s add-on to namespaces %s of managed cluster: %s.\n", addon, o.Namespace, clusterName)
		}
	}

	return apply.WriteOutput(o.OutputFile, output)
}
