package session

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p-core/peer"
)

// GetLeader returns the leader of the given transaction hash and peer.IDSlice.
func GetLeader(txHash common.Hash, peerIDs peer.IDSlice) peer.ID {
	// Mutates peerIDs slice
	sort.Sort(peerIDs)

	leaderIndexBig := new(big.Int).Mod(txHash.Big(), big.NewInt(int64(len(peerIDs))))

	// Should not happen as there shouldn't be maxint number of peers.
	if !leaderIndexBig.IsUint64() {
		panic("leader index is not uint64")
	}

	return peerIDs[leaderIndexBig.Uint64()]
}
