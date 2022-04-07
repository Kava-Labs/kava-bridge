package p2p

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/pnet"
)

type NodeOptions struct {
	Port              uint16
	NetworkPrivateKey pnet.PSK
	NodePrivateKey    crypto.PrivKey
}
