package k8s

import (
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	*kubernetes.Clientset
}

var k8sClientInstance *K8sClient
var once sync.Once

func getKubeConfig() (*rest.Config, error) {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func GetK8SClient() (*K8sClient, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	once.Do(func() {
		k8sClientInstance = &K8sClient{clientSet}
	})
	return k8sClientInstance, err
}
