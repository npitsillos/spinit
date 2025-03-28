package deploy

import (
	"github.com/npitsillos/spinit/config"
	"github.com/npitsillos/spinit/pkg/k8s"
	"github.com/spf13/viper"
)

func CreateDeployment(name string, expose bool, namespace string) error {
	k8sClient, err := k8s.GetK8SClient()

	if err != nil {
		return err
	}
	if err := k8sClient.CreateNamespaceIfNotExists(namespace); err != nil {
		return err
	}
	if err := k8sClient.CreateDeployment(name, namespace); err != nil {
		return err
	}

	if expose {
		return k8sClient.ExposeApp(name, namespace, viper.Get("config").(*config.Config).Domain)
	}
	return nil
}
