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
		Example: "roller explorer start --home /path/to/home -f \"NEXT_PUBLIC_NETWORK_NAME=Awesome rollapp\" " +
			"-f \"NEXT_PUBLIC_NETWORK_CURRENCY_SYMBOL=TKN\"",
		Run: func(cmd *cobra.Command, args []string) {
			home := cmd.Flag(utils.FlagNames.Home).Value.String()
			backendEnvs, err := cmd.Flags().GetStringArray("backend-envs")
			utils.PrettifyErrorIfExists(err)

			frontendEnvs, err := cmd.Flags().GetStringArray("frontend-envs")
			utils.PrettifyErrorIfExists(err)

			explorer := explorer.NewExplorer(config.Blockscout, home)
			err = explorer.Start(backendEnvs, frontendEnvs)
			utils.PrettifyErrorIfExists(err)
		},
	}
	cmd.PersistentFlags().StringArrayP("backend-envs", "b", []string{}, "The environment variables for the backend service.")
	cmd.PersistentFlags().StringArrayP("frontend-envs", "f", []string{}, "The environment variables for the frontend service.")
	return cmd
}
