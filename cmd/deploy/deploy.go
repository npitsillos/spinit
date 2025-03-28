package deploy

import (
	"path/filepath"

	"github.com/npitsillos/spinit/config"
	"github.com/npitsillos/spinit/helpers"
	"github.com/npitsillos/spinit/pkg/deploy"
	"github.com/spf13/cobra"
)

const (
	NAMESPACE_FLAG = "namespace"
	EXPOSE_FLAG    = "expose"
)

func NewDeployCommand() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy APP",
		Short: "deploy application",
		Long:  `Deploys application to the local cluster optionally exposes it on a domain.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			workDir := helpers.GetWorkingDir()
			appCfg, err := config.LoadAppConfig(filepath.Join(workDir, config.DefaultAppConfigFile))
			cobra.CheckErr(err)
			namespace, _ := cmd.Flags().GetString(NAMESPACE_FLAG)
			if namespace == "" {
				namespace = appCfg.AppName
			}
			deploy.CreateDeployment(appCfg, false, namespace)
		},
	}

	deployCmd.Flags().StringP(NAMESPACE_FLAG, "n", "", "Which namespace to deploy app in")
	deployCmd.Flags().BoolP(EXPOSE_FLAG, "e", false, "Whether to create a service & expose the app")

	return deployCmd
}
