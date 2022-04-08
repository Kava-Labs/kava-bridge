package network

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multibase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logging.Logger("connect")

func newConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [remote peer multiaddr]",
		Short: "Connects the relayer to peers",
		Args:  cobra.MaximumNArgs(1),
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

			privSharedKeyData, err := os.ReadFile(privateSharedKeyPath)
			if err != nil {
				return fmt.Errorf("could not read pre-shared key: %w", err)
			}
			_, privSharedKey, err := multibase.Decode(string(privSharedKeyData))
			if err != nil {
				return err
			}

			options := p2p.NodeOptions{
				Port:              uint16(port),
				NodePrivateKey:    privKey,
				NetworkPrivateKey: privSharedKey,
			}

			node, err := p2p.NewNode(options)
			if err != nil {
				return err
			}

			multiAddr, err := node.GetMultiAddress()
			if err != nil {
				return err
			}
			log.Info("host multiaddress: ", multiAddr)

			if len(args) == 1 {
				log.Info("connecting to remote peer: ", args[0])
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				peerAddr, err := multiaddr.NewMultiaddr(args[0])
				if err != nil {
					return fmt.Errorf("could not parse peer multiaddr: %w", err)
				}

				if err := node.Connect(ctx, peerAddr); err != nil {
					return err
				}
			}

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			log.Info("Received signal, shutting down...")
			return node.Host.Close()
		},
	}

	cmd.Flags().Uint16(p2pFlagPort, 0, "Host port to listen on (required)")
	cmd.Flags().String(p2pFlagPrivateKeyPath, "", "Path to the peer private key (required)")
	cmd.Flags().String(p2pFlagSharedKeyPath, "", "Path to the shared private network key (required)")

	// Ignore errors, err only if flags do not exist
	_ = cmd.MarkFlagRequired(p2pFlagPort)
	_ = cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)
	_ = cmd.MarkFlagRequired(p2pFlagSharedKeyPath)

	return cmd
}
