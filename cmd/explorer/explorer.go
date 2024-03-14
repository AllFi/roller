package explorer

import (
	"github.com/dymensionxyz/roller/cmd/explorer/clear"
	"github.com/dymensionxyz/roller/cmd/explorer/start"
	"github.com/dymensionxyz/roller/cmd/explorer/stop"
	"github.com/spf13/cobra"
)

func ExplorerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explorer",
		Short: "Commands for managing the block explorer.",
	}
	cmd.AddCommand(start.Cmd())
	cmd.AddCommand(stop.Cmd())
	cmd.AddCommand(clear.Cmd())
	return cmd
}
