package session_test

import (
	"context"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/session"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/require"
)

func TestSigningSessionStore(t *testing.T) {
	store := session.NewSigningSessionStore()

	sess, found := store.GetSessionFromTxHash(common.Hash{1})
	require.Nil(t, sess)
	require.False(t, found, "should not find session that does not exist")

	sess, _, err := store.NewSession(
		context.Background(),
		&broadcast.NoOpBroadcaster{},
		common.Hash{1},
		nil,
		0,
		"",
		testutil.TestPeerIDs,
		nil,
		nil,
		keygen.LocalPartySaveData{},
	)

	require.NoError(t, err)

	sess2, found := store.GetSessionFromTxHash(common.Hash{1})
	require.True(t, found, "should find session that exists")

	require.Equal(t, sess, sess2, "should return same session")
}
