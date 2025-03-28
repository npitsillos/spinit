package deploy

import (
	"github.com/npitsillos/spinit/config"
	"github.com/npitsillos/spinit/pkg/k8s"
	"github.com/spf13/viper"
)

func CreateDeployment(appCfg *config.AppConfig, expose bool, namespace string) error {
	k8sClient, err := k8s.GetK8SClient()

	if err != nil {
		return err
	}
	if err := k8sClient.CreateNamespaceIfNotExists(namespace); err != nil {
		return err
	}
	if err := k8sClient.CreateDeployment(appCfg.AppName, appCfg.Build.Image, namespace); err != nil {
		return err
	}

	if expose {
		return k8sClient.ExposeApp(appCfg.AppName, namespace, viper.Get("config").(*config.Config).Domain)
	}
	return nil
}
