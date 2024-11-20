package root

import (
	"github.com/npitsillos/spinit/cmd/build"
	"github.com/npitsillos/spinit/cmd/config"
	"github.com/npitsillos/spinit/cmd/copy"
	"github.com/npitsillos/spinit/cmd/deploy"
	initConfig "github.com/npitsillos/spinit/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "spinit",
		Short:   "spinit helps package and deploy applications to local Kubernetes clusters",
		Long:    "spinit helps package and deploy applications to local Kubernetes clusters",
		Version: "0.1.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configExists, err := initConfig.ConfigDirExists()
			if err != nil {
				return err
			}
			var cfg *initConfig.Config
			if !configExists {
				cfg, err = initConfig.InitSpinitConfig()
				if err != nil {
					return err
				}
			} else {
				cfg, err = initConfig.LoadConfig()
				if err != nil {
					return err
				}
			}
			viper.Set("config", cfg)
			return nil
		},
	}

	rootCmd.AddCommand(deploy.NewDeployCommand())
	rootCmd.AddCommand(build.NewBuildCommand())
	rootCmd.AddCommand(copy.NewCopyCommand())
	rootCmd.AddCommand(config.NewInitCommand())
	return rootCmd
}
