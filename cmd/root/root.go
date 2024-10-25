package root

import (
	"fmt"

	"github.com/npitsillos/spinit/cmd/build"
	"github.com/npitsillos/spinit/cmd/copy"
	"github.com/npitsillos/spinit/cmd/deploy"
	"github.com/npitsillos/spinit/config"
	"github.com/spf13/cobra"
)

func firstRun() (bool, error) {
	return config.ConfigDirExists()
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "spinit",
		Short:   "spinit helps package and deploy applications to local Kubernetes clusters",
		Long:    "spinit helps package and deploy applications to local Kubernetes clusters",
		Version: "0.1.0",
		RunE:    run,
	}

	rootCmd.AddCommand(deploy.NewDeployCommand())
	rootCmd.AddCommand(build.NewBuildCommand())
	rootCmd.AddCommand(copy.NewCopyCommand())
	return rootCmd
}

func run(cmd *cobra.Command, args []string) error {

	firstRun, err := firstRun()
	if err != nil {
		return err
	}
	var cfg *config.Config
	if firstRun {
		cfg, err = config.InitSpinitConfig()
		if err != nil {
			return err
		}
	} else {
		cfg, err = config.LoadConfig()
		if err != nil {
			return err
		}
	}
	fmt.Println(cfg.SSH)
	fmt.Println(cfg.Nodes)
	return nil
}
