package network

import "github.com/spf13/cobra"

func GetNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage the private p2p network",
	}

	cmds := []*cobra.Command{
		newConnectCmd(),
		newShowNodeIdCmd(),
		newShowNodeMultiAddressCmd(),
		newGenerateNetworkSecretCmd(),
		newGenerateNodeKeyCmd(),
	}

	cmd.AddCommand(cmds...)

	return cmd
}
