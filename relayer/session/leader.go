package session

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p-core/peer"
)

// GetLeader returns the leader of the given transaction hash and peer.IDSlice.
// The offset is only used when there is an inactive leader. The initial offset
// should be 0, and incremented only while the leader does not respond.
func GetLeader(txHash common.Hash, peerIDs peer.IDSlice, offset int64) (peer.ID, error) {
	if len(peerIDs) == 0 {
		return "", fmt.Errorf("no peers provided")
	}

	if offset < 0 {
		return "", fmt.Errorf("offset must be >= 0")
	}

	// Make a copy to prevent mutation of the original slice.
	copiedPeerIDs := make(peer.IDSlice, len(peerIDs))
	_ = copy(copiedPeerIDs, peerIDs)

	// Sort copy
	sort.Sort(copiedPeerIDs)

	// (hash + offset) % len(peerIDs)
	sum := new(big.Int).Add(txHash.Big(), big.NewInt(offset))
	leaderIndexBig := sum.Mod(sum, big.NewInt(int64(len(copiedPeerIDs))))

	// Should not happen as there shouldn't be maxint number of peers.
	if !leaderIndexBig.IsUint64() {
		panic("leader index is not uint64")
	}

	return copiedPeerIDs[leaderIndexBig.Uint64()], nil
}
