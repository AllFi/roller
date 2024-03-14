package start

import (
	"github.com/dymensionxyz/roller/cmd/utils"
	"github.com/dymensionxyz/roller/config"
	"github.com/dymensionxyz/roller/explorer"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Commands for running the block explorer.",
		Run: func(cmd *cobra.Command, args []string) {
			home := cmd.Flag(utils.FlagNames.Home).Value.String()
			explorer := explorer.NewExplorer(config.Blockscout, home)
			err := explorer.Start()
			utils.PrettifyErrorIfExists(err)
		},
	}
	return cmd
}
