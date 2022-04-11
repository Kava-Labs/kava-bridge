package network

import (
	"context"
	"fmt"
	"os"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/multiformats/go-multibase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logging.Logger("connect")

func newConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connects the relayer to peers",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

			port := viper.GetUint(p2pFlagPort)
			pkPath := viper.GetString(p2pFlagPrivateKeyPath)
			privateSharedKeyPath := viper.GetString(p2pFlagSharedKeyPath)
			peerAddrStrings := viper.GetStringSlice(p2pFlagPeerMultiAddrs)

			privKeyData, err := os.ReadFile(pkPath)
			if err != nil {
				return fmt.Errorf("could not read private key file: %w", err)
			}
			privKey, err := p2p.UnmarshalPrivateKey(privKeyData)
			if err != nil {
				return err
			}

			privSharedKeyData, err := os.ReadFile(privateSharedKeyPath)
			if err != nil {
				return fmt.Errorf("could not read pre-shared key: %w", err)
			}
			_, privSharedKey, err := multibase.Decode(string(privSharedKeyData))
			if err != nil {
				return err
			}

			peerAddrInfos, err := p2p.ParseMultiAddrSlice(peerAddrStrings)
			if err != nil {
				return err
			}

			options := p2p.NodeOptions{
				Port:              uint16(port),
				NodePrivateKey:    privKey,
				NetworkPrivateKey: privSharedKey,
				// Require response from all peers
				EchoRequiredPeers: len(peerAddrInfos),
			}

			// Need to be buffered by 1 to not block
			done := make(chan bool, 1)

			node, err := p2p.NewNode(options, done)
			if err != nil {
				return err
			}

			multiAddr, err := node.GetMultiAddress()
			if err != nil {
				return err
			}

			log.Info("host multiaddress: ", multiAddr)
			log.Info("connecting to remote peers: ", peerAddrInfos)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if err := node.ConnectToPeers(ctx, peerAddrInfos); err != nil {
				return err
			}

			if err := node.EchoPeers(ctx); err != nil {
				return err
			}

			log.Info("waiting for all echo requests")

			<-done
			log.Info("Done! exiting...")
			return node.Close()
		},
	}

	cmd.Flags().Uint16(p2pFlagPort, 0, "Host port to listen on (required)")
	cmd.Flags().String(p2pFlagPrivateKeyPath, "", "Path to the peer private key (required)")
	cmd.Flags().String(p2pFlagSharedKeyPath, "", "Path to the shared private network key (required)")
	cmd.Flags().StringSlice(p2pFlagPeerMultiAddrs, []string{}, "List of peer multiaddrs (required)")

	// Ignore errors, err only if flags do not exist
	_ = cmd.MarkFlagRequired(p2pFlagPort)
	_ = cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)
	_ = cmd.MarkFlagRequired(p2pFlagSharedKeyPath)
	_ = cmd.MarkFlagRequired(p2pFlagPeerMultiAddrs)

	return cmd
}
