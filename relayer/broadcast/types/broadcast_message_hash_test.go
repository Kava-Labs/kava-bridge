package types_test

import (
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/require"
)

func TestHash_SamePayload(t *testing.T) {
	msg := MustNewBroadcastMessage(
		&types.HelloRequest{
			Message: "hi",
		},
		testutil.TestPeerIDs[0],
		testutil.TestPeerIDs[1:2],
		5,
	)

	hash, err := msg.Hash()
	require.NoError(t, err)
	require.Len(t, hash, 32)

	msg2 := MustNewBroadcastMessage(
		&types.HelloRequest{
			Message: "hi",
		},
		testutil.TestPeerIDs[0],
		testutil.TestPeerIDs[1:2],
		5,
	)
	// Copy created ID and timestamp to be the same
	msg2.ID = msg.ID
	msg2.Created = msg.Created

	hash2, err := msg2.Hash()
	require.NoError(t, err)
	require.Len(t, hash2, 32)

	require.Equal(t, hash, hash2, "hash of same message should be the same")
}

func TestHash_DifferentPayload(t *testing.T) {
	msg := MustNewBroadcastMessage(
		&types.HelloRequest{
			Message: "hi",
		},
		testutil.TestPeerIDs[0],
		testutil.TestPeerIDs[1:2],
		5,
	)

	hash, err := msg.Hash()
	require.NoError(t, err)
	require.Len(t, hash, 32)

	msg2 := MustNewBroadcastMessage(
		&types.HelloRequest{
			Message: "bye",
		},
		testutil.TestPeerIDs[0],
		testutil.TestPeerIDs[1:2],
		5,
	)

	hash2, err := msg2.Hash()
	require.NoError(t, err)
	require.Len(t, hash2, 32)

	require.NotEqual(t, hash, hash2, "hash of different message should be different")
}

func TestHash_DifferentRecipients(t *testing.T) {
	msg := MustNewBroadcastMessage(
		&types.HelloRequest{
			Message: "hi",
		},
		testutil.TestPeerIDs[0],
		testutil.TestPeerIDs[1:2],
		5,
	)

	hash, err := msg.Hash()
	require.NoError(t, err)
	require.Len(t, hash, 32)

	t.Log(msg.RecipientPeerIDs)

	msg2 := msg
	// Should include host if directly mutating the recipient list
	msg2.RecipientPeerIDs = testutil.TestPeerIDs[0:3]

	t.Log(msg2.RecipientPeerIDs)

	hash2, err := msg2.Hash()
	require.NoError(t, err)
	require.Len(t, hash2, 32)

	require.NotEqual(t, hash, hash2, "hash with different recipients should be different")
}
