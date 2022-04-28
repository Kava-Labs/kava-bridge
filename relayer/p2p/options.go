package p2p

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/pnet"
)

type NodeOptions struct {
	Port              uint16
	NetworkPrivateKey pnet.PSK       // Shared secret for private network
	NodePrivateKey    crypto.PrivKey // Private key for current node
	PeerList          []peer.ID      // Allowlist of peers to connect to
	EchoRequiredPeers int
}
