package clear

import (
	"github.com/dymensionxyz/roller/cmd/utils"
	"github.com/dymensionxyz/roller/config"
	"github.com/dymensionxyz/roller/explorer"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clears the block explorer data.",
		Run: func(cmd *cobra.Command, args []string) {
			home := cmd.Flag(utils.FlagNames.Home).Value.String()
			explorer := explorer.NewExplorer(config.Blockscout, home)
			err := explorer.Clear()
			utils.PrettifyErrorIfExists(err)
		},
	}
	return cmd
}
