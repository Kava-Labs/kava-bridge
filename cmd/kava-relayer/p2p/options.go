package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/libp2p/go-tcp-transport"
)

type NodeOptions struct {
	Port              uint16
	NetworkPrivateKey pnet.PSK
	NodePrivateKey    crypto.PrivKey
}

func NewNode(options NodeOptions) (host.Host, error) {
	libp2pOpts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", options.Port)),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.PrivateNetwork(options.NetworkPrivateKey),
		libp2p.Identity(options.NodePrivateKey),
		libp2p.DisableRelay(),
	}

	return libp2p.New(libp2pOpts...)
}
