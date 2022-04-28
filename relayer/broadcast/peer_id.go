package broadcast

import "github.com/libp2p/go-libp2p-core/peer"

// PeerIDMapToSlice converts a map of peer.ID to a slice of peer.ID.
func PeerIDMapToSlice(peerIDs map[peer.ID]struct{}) []peer.ID {
	peerIDSlice := make([]peer.ID, len(peerIDs))
	for peerID := range peerIDs {
		peerIDSlice = append(peerIDSlice, peerID)
	}

	return peerIDSlice
}
