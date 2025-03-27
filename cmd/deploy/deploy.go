package deploy

import (
	"github.com/npitsillos/spinit/pkg/deploy"
	"github.com/spf13/cobra"
)

func NewDeployCommand() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy APP",
		Short: "deploy application",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			deploy.CreateDeployment(args[0], false, "test")
		},
	}

	deployCmd.Flags().StringP("namespace", "n", "default", "Which namespace to deploy app in")
	deployCmd.Flags().BoolP("expose", "e", false, "Whether to create a service & expose the app")

	return deployCmd
}
