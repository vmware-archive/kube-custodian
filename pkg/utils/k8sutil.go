// Package utils taken from github.com/bitnami-labs/kubewatch source
package utils

import (
	"os"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClient returns a k8s clientset to the request from inside of cluster
func GetClient() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		logrus.Fatalf("Can not get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Can not create kubernetes client: %v", err)
	}

	return clientset
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

// GetClientOutOfCluster returns a k8s clientset to the request from outside of cluster
func GetClientOutOfCluster() (*kubernetes.Clientset, error) {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		logrus.Fatalf("Can not get kubernetes config: %v", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	return clientset, nil
}

// LabelsSubSet checks that subSet label keys are all present in Labels
func LabelsSubSet(Labels map[string]string, subSet []string) bool {
	labelKeys := []string{}
	for k := range Labels {
		labelKeys = append(labelKeys, k)
	}
	labelSet := sets.NewString(labelKeys...)

	if labelSet.HasAll(subSet...) {
		return true
	}
	return false
}
