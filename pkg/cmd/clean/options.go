// Copyright Contributors to the Open Cluster Management project

package init

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	genericclioptionsclusteradm "open-cluster-management.io/clusteradm/pkg/genericclioptions"
)

//Options: The structure holding all the command-line options
type Options struct {
	//ClusteradmFlags: The generic optiosn from the clusteradm cli-runtime.
	ClusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags
	Values          Values
	//The cluster manager resource name
	ClusterManageName string
	//The file to output the resources will be sent to the file.
	OutputFile string
	//If true the bootstrap token will be used instead of the service account token
	UseBootstrapToken bool
}

//Valus: The values used in the template
type Values struct {
	//The values related to the hub
	Hub Hub `json:"hub"`
}

//Hub: The hub values for the template
type Hub struct {
	//TokenID: A token id allowing the cluster to connect back to the hub
	TokenID string `json:"tokenID"`
	//TokenSecret: A token secret allowing the cluster to connect back to the hub
	TokenSecret string `json:"tokenSecret"`
}

func NewOptions(clusteradmFlags *genericclioptionsclusteradm.ClusteradmFlags, streams genericclioptions.IOStreams) *Options {
	return &Options{
		ClusteradmFlags: clusteradmFlags,
	}
}
