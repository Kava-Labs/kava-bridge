package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

// GetNodeID returns the node ID with the given private key.
func GetNodeID(privKey crypto.PrivKey) (peer.ID, error) {
	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		return "", fmt.Errorf("could not get peer ID: %w", err)
	}

	return peerID, nil
}
