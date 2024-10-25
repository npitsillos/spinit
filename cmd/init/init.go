package init

import (
	"github.com/npitsillos/spinit/config"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "initialise spinit config",
		Long: `Initialises spinit config by obtaining user input that provides
		information on the nodes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.InitSpinitConfig()
			return err
		},
	}

	return initCmd
}
