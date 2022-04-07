package network

import (
	"fmt"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newShowNodeIdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-node-id",
		Short: "Shows the node id",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

			pkPath := viper.GetString(p2pFlagPrivateKeyPath)

			peerID, err := p2p.GetNodeID(pkPath)
			if err != nil {
				return err
			}

			fmt.Println(peerID)

			return nil
		},
	}

	cmd.Flags().String(p2pFlagPrivateKeyPath, "", "Path to the peer private key (required)")
	_ = cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)

	return cmd
}
