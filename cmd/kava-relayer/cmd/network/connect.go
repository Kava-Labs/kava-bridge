package network

import (
	"github.com/kava-labs/kava-bridge/cmd/kava-relayer/p2p"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connects the relayer to peers",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

			port := viper.GetUint("p2p.port")

			options := p2p.NodeOptions{
				Port: uint16(port),
			}

			node, err := p2p.NewNode(options)
			if err != nil {
				return err
			}

			node.Close()

			return nil
		},
	}

	cmd.Flags().Uint16(p2pFlagPort, 0, "Host port to listen on")
	cmd.Flags().String(p2pFlagPrivateKeyPath, "", "Path to the peer private key (required)")
	cmd.Flags().String(p2pFlagSharedKeyPath, "", "Path to the shared private network key (required)")

	cmd.MarkFlagRequired(p2pFlagPort)
	cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)
	cmd.MarkFlagRequired(p2pFlagSharedKeyPath)

	return cmd
}
