package load

import (
	"github.com/npitsillos/spinit/pkg/load"
	"github.com/spf13/cobra"
)

var (
	NODE_FLAG = "node"
)

func NewLoadCommand() *cobra.Command {
	loadCmd := &cobra.Command{
		Use:   "copy",
		Short: "copies image to nodes in cluster",
		Long: `Copies an image from local to the cluster's nodes. Passing the --node flag
		will only copy the image on the nodes specified.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, _ := cmd.Flags().GetStringArray(NODE_FLAG)
			load.LoadImageToNodes(args[0], nodes)
			return nil
		},
	}

	loadCmd.Flags().StringArrayP(NODE_FLAG, "n", []string{}, "Node to copy image to")

	return loadCmd
}
