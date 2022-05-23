package session_test

import (
	"math/big"
	"math/rand"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/session"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func TestSelectLeader(t *testing.T) {
	txHash := common.BytesToHash([]byte("hello there"))

	expectedLeader := session.GetLeader(txHash, testutil.TestPeerIDs, 0)

	for i := 0; i < len(testutil.TestPeerIDs); i++ {
		// Make a copy of the slice so we can freely mutate it.
		// GetLeader() sorts in place.
		randomSortedPeerIDs := make([]peer.ID, len(testutil.TestPeerIDs))
		copy(randomSortedPeerIDs, testutil.TestPeerIDs)

		rand.Shuffle(len(testutil.TestPeerIDs), func(i, j int) {
			randomSortedPeerIDs[i], randomSortedPeerIDs[j] = randomSortedPeerIDs[j], randomSortedPeerIDs[i]
		})

		leader := session.GetLeader(txHash, randomSortedPeerIDs, 0)
		require.Equal(t, expectedLeader, leader, "leader should be the same for any order of input peers")
	}
}

func TestSelectLeader_CorrectIndex(t *testing.T) {
	sorted := make(peer.IDSlice, len(testutil.TestPeerIDs))
	copy(sorted, testutil.TestPeerIDs)
	sort.Sort(sorted)

	for i := 0; i < len(sorted)*2; i++ {
		// index as tx hash, hash % len(peerIDs) should be equal to i
		txHash := common.BigToHash(big.NewInt(int64(i)))

		// Sorts TestPeerIDs in place, TestPeerIDs is in order after this
		leader := session.GetLeader(txHash, sorted, 0)

		// use i % len(peerIDs) to get the index of the leader
		require.Equal(
			t,
			sorted[i%len(sorted)],
			leader,
			"leader should be the correct index based on hash",
		)
	}
}

func TestSelectLeader_Offset(t *testing.T) {
	sorted := make(peer.IDSlice, len(testutil.TestPeerIDs))
	copy(sorted, testutil.TestPeerIDs)
	sort.Sort(sorted)

	for i := 0; i < len(sorted)*2; i++ {
		// index as tx hash, hash % len(peerIDs) should be equal to i
		txHash := common.BigToHash(big.NewInt(0))

		// Sorts TestPeerIDs in place, TestPeerIDs is in order after this
		leader := session.GetLeader(txHash, sorted, int64(i))

		// use i % len(peerIDs) to get the index of the leader
		require.Equal(
			t,
			sorted[i%len(sorted)],
			leader,
			"leader should be the correct index based on hash offset",
		)
	}
}
