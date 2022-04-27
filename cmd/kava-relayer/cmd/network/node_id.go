package network

import (
	"fmt"
	"os"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newShowNodeIdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-node-id",
		Short: "Shows the node id",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkPath := viper.GetString(p2pFlagPrivateKeyPath)

			privKeyData, err := os.ReadFile(pkPath)
			if err != nil {
				return fmt.Errorf("could not read private key file: %w", err)
			}
			privKey, err := p2p.UnmarshalPrivateKey(privKeyData)
			if err != nil {
				return err
			}

			peerID, err := p2p.GetNodeID(privKey)
			if err != nil {
				return err
			}

			fmt.Print(peerID.String())

			return nil
		},
	}

	cmd.Flags().String(p2pFlagPrivateKeyPath, "", "Path to the peer private key (required)")
	_ = cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)

	return cmd
}
