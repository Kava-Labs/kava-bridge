package signing

import (
	"fmt"
	"math/rand"

	"github.com/libp2p/go-libp2p-core/peer"
)

// PeerIDSelector defines an interface that selects t + 1 peers for signing.
type PeerIDSelector interface {
	// GetSigners returns the t + 1 slice of peers that are picked to sign a
	// message.
	GetSigners(threshold int, allPeerIDs peer.IDSlice) (peer.IDSlice, error)
}

// RandomPeerIDSelector is a PeerIDSelector that picks random t + 1 peers.
type RandomPeerIDSelector struct{}

var _ PeerIDSelector = (*RandomPeerIDSelector)(nil)

// GetSigners returns a random slice of t + 1 peers to sign a message.
func (s *RandomPeerIDSelector) GetSigners(
	threshold int,
	allPeerIDs peer.IDSlice,
) (peer.IDSlice, error) {
	if len(allPeerIDs) < threshold+1 {
		return nil, fmt.Errorf(
			"not enough peers to select signers, %d peers < %d t + 1",
			len(allPeerIDs), threshold+1,
		)
	}

	// Copy allPeerIDs to avoid mutating the original.
	allPeerIDsCopy := make(peer.IDSlice, len(allPeerIDs))
	copy(allPeerIDsCopy, allPeerIDs)

	rand.Shuffle(len(allPeerIDsCopy), func(i, j int) {
		allPeerIDsCopy[i], allPeerIDsCopy[j] = allPeerIDsCopy[j], allPeerIDsCopy[i]
	})

	return allPeerIDsCopy[:threshold+1], nil
}
