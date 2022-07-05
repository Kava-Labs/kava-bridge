package network

import (
	"fmt"
	"net"
	"os"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newShowNodeMultiAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-node-multi-address",
		Short: "Shows the multi address of the node that peers should use for communication",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

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

			port := viper.GetInt(p2pFlagPort)

			ipv4s, err := GetHostIPv4s()
			if err != nil {
				return err
			}

			for _, ipv4 := range ipv4s {
				multiAddress, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", ipv4, port, peerID))
				if err != nil {
					return err
				}

				fmt.Println(multiAddress)
			}

			return nil
		},
	}

	cmd.Flags().String(p2pFlagPrivateKeyPath, "", "Path to the peer private key (required)")
	_ = cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)
	cmd.Flags().String(p2pFlagPort, "", "TCP port that the node will listen on for peer requests")
	_ = cmd.MarkFlagRequired(p2pFlagPrivateKeyPath)
	return cmd
}

// GetHostIPv4s returns a slice lf valid IPv4 address bound to a
// network interface for the current host and error (if any).
func GetHostIPv4s() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var publicIPv4s []net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				publicIPv4s = append(publicIPv4s, ipnet.IP.To4())
			}
		}
	}

	if len(publicIPv4s) == 0 {
		return nil, fmt.Errorf("no public IPv4 address found")
	}

	return publicIPv4s, nil
}
