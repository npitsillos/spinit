package deploy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDeployCommand() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy application",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := viper.Get("config")
			fmt.Println(cfg)
		},
	}

	return deployCmd
}
