package network

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newGenerateNodeKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-node-key",
		Short: "Generates a node secret key",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

			privKey, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
			if err != nil {
				return fmt.Errorf("could not read from rand and generate keypair: %w", err)
			}

			b, err := crypto.MarshalPrivateKey(privKey)
			if err != nil {
				return fmt.Errorf("could not marshal generated private key: %w", err)
			}

			os.Stdout.Write(b)

			return nil
		},
	}

	return cmd
}
