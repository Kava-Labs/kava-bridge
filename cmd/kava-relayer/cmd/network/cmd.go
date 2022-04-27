package network

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage the private p2p network",
	}

	cmds := []*cobra.Command{
		newConnectCmd(),
		newShowNodeIdCmd(),
		newGenerateNetworkSecretCmd(),
		newGenerateNodeKeyCmd(),
	}

	// Bind pflags here so that they show up in init command. If done in RunE(),
	// they will not be added to the initialized config file.
	for _, cmd := range cmds {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Fatal("unable to bind flags", err)
		}
	}

	cmd.AddCommand(cmds...)

	return cmd
}
