package network

import (
	"crypto/rand"
	"fmt"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/multiformats/go-multibase"
	"github.com/spf13/cobra"
)

func newGenerateNetworkSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-network-secret",
		Short: "Generates a network secret key",
		RunE: func(cmd *cobra.Command, args []string) error {
			b := make([]byte, p2p.PreSharedNetworkKeyLengthBytes)

			_, err := rand.Read(b)
			if err != nil {
				return fmt.Errorf("could not read from rand: %w", err)
			}

			s, err := multibase.Encode(multibase.Base58BTC, b)
			if err != nil {
				return fmt.Errorf("could not encode random bytes: %w", err)
			}

			fmt.Print(s)

			return nil
		},
	}

	return cmd
}
