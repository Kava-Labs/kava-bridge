package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/libp2p/go-tcp-transport"
	"github.com/spf13/viper"
)

func ParseOptions() []libp2p.Option {
	// No uint16
	port := viper.GetUint("p2p.port")

	// TODO:
	psk := pnet.PSK{}

	return []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.PrivateNetwork(psk),
		// libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}
}

func NewNode(opts ...libp2p.Option) (host.Host, error) {
	return libp2p.New(opts...)
}
