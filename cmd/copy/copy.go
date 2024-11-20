package copy

import (
	"github.com/npitsillos/spinit/pkg/copy"
	"github.com/spf13/cobra"
)

var (
	NODE_FLAG = "node"
)

func NewCopyCommand() *cobra.Command {
	copyCmd := &cobra.Command{
		Use:   "copy",
		Short: "copies image to nodes in cluster",
		Long: `Copies an image from local to the cluster's nodes. Passing the --node flag
		will only copy the image on the nodes specified.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, _ := cmd.Flags().GetStringArray(NODE_FLAG)
			copy.CopyImageToNodes(args[0], nodes)
			return nil
		},
	}

	copyCmd.Flags().StringArrayP(NODE_FLAG, "n", []string{}, "Node to copy image to")

	return copyCmd
}
