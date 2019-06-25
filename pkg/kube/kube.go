package kube

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient creates k8s client
func NewClient(kubeConfigFile string) (kubernetes.Interface, error) {
	var clientConfig *rest.Config
	var err error
	if len(kubeConfigFile) > 0 {
		clientConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			return nil, err
		}
	} else {
		clientConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}