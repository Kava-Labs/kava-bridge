package session

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p-core/peer"
)

// GetLeader returns the leader of the given transaction hash and peer.IDSlice.
func GetLeader(txHash common.Hash, peerIDs peer.IDSlice) peer.ID {
	// Make a copy to prevent mutation of the original slice.
	copiedPeerIDs := make(peer.IDSlice, len(peerIDs))
	_ = copy(copiedPeerIDs, peerIDs)

	// Sort copy
	sort.Sort(copiedPeerIDs)

	leaderIndexBig := new(big.Int).Mod(txHash.Big(), big.NewInt(int64(len(copiedPeerIDs))))

	// Should not happen as there shouldn't be maxint number of peers.
	if !leaderIndexBig.IsUint64() {
		panic("leader index is not uint64")
	}

	return copiedPeerIDs[leaderIndexBig.Uint64()]
}
