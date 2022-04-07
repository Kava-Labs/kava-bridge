package network

import (
	"fmt"
	"os"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/multiformats/go-multibase"
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

			port := viper.GetUint(p2pFlagPort)
			pkPath := viper.GetString(p2pFlagPrivateKeyPath)
			privateSharedKeyPath := viper.GetString(p2pFlagSharedKeyPath)

			privKeyData, err := os.ReadFile(pkPath)
			if err != nil {
				return fmt.Errorf("could not read private key file: %w", err)
			}
			privKey, err := p2p.UnmarshalPrivateKey(privKeyData)
			if err != nil {
				return err
			}

			pskData, err := os.ReadFile(privateSharedKeyPath)
			if err != nil {
				return fmt.Errorf("could not read pre-shared key: %w", err)
			}
			_, psk, err := multibase.Decode(string(pskData))
			if err != nil {
				return err
			}

			options := p2p.NodeOptions{
				Port:              uint16(port),
				NodePrivateKey:    privKey,
				NetworkPrivateKey: psk,
			}

			node, err := p2p.NewNode(options)
			if err != nil {
				return err
			}

			// TODO: Do something with the node
			return node.Close()
		},
	}

	cmd.Flags().Uint16(p2pFlagPort, 0, "Host port to listen on")
	cmd.Flags().String(p2pFlagPrivateKeyPath, "", "Path to the peer private key (required)")
	cmd.Flags().String(p2pFlagSharedKeyPath, "", "Path to the shared private network key (required)")

	// Ignore errors, if flags do not exist
	_ = cmd.MarkFlagRequired(p2pFlagPort)
	_ = cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)
	_ = cmd.MarkFlagRequired(p2pFlagSharedKeyPath)

	return cmd
}
