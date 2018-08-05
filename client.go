package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var cs *kubernetes.Clientset

func newKubernetesClient() (*kubernetes.Clientset, error) {
	if cs != nil {
		return cs, nil
	}

	config, err := rest.InClusterConfig()
	//config, err := getLocalClientSet("dispatch")
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		cs = client
	}

	return client, err
}

func getLocalClientSet(context string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	).ClientConfig()
}
