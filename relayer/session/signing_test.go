package session_test

import (
	"context"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/kava-labs/kava-bridge/relayer/session"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/require"
)

func TestSigningSession(t *testing.T) {
	store := session.NewSigningSessionStore()

	sess, _, err := session.NewSigningSession(
		context.Background(),
		store,
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
	require.NotNil(t, sess)

	// Non-leader state, does not accept JoinSessionMessage
	err = sess.Update(session.NewAddCandidateEvent(nil, types.JoinSessionMessage{}))
	require.Error(t, err)

	// Not signing yet
	err = sess.Update(session.NewAddSigningPartEvent(nil, nil, false))
	require.Error(t, err)

	// Only accepts StartSignerEvent
	err = sess.Update(session.NewStartSignerEvent(nil, nil, nil))
	require.Nil(t, err)

	// --------------------
	// Now in DoneNonParticipantState

	// Does not accept any more events
	err = sess.Update(session.NewAddCandidateEvent(nil, types.JoinSessionMessage{}))
	require.Error(t, err)

	err = sess.Update(session.NewStartSignerEvent(nil, nil, nil))
	require.Error(t, err)

	err = sess.Update(session.NewAddSigningPartEvent(nil, nil, false))
	require.Error(t, err)
}

func TestSigningSession_Leader(t *testing.T) {
	store := session.NewSigningSessionStore()

	sess, _, err := session.NewSigningSession(
		context.Background(),
		store,
		&broadcast.NoOpBroadcaster{},
		common.Hash{1},
		nil,                     // msg to sign
		2,                       // threshold
		testutil.TestPeerIDs[3], // current peer id -- must be 3 to be leader
		testutil.TestPeerIDs,
		nil, // current partyID
		nil, // partyID store
		keygen.LocalPartySaveData{},
	)

	require.NoError(t, err)
	require.NotNil(t, sess)

	// Leader accepts NewAddCandidateEvent
	err = sess.Update(session.NewAddCandidateEvent(
		nil,
		types.NewJoinSigningSessionMessage(
			testutil.TestPeerIDs[0],
			common.Hash{1},
			types.SigningSessionIDPart{1, 2, 3},
		),
	))
	require.NoError(t, err)

	// ----
	// Does not accept any other events

	err = sess.Update(session.NewStartSignerEvent(nil, nil, nil))
	require.Error(t, err)
	require.Equal(t, "not enough candidates to start signing", err.Error())

	err = sess.Update(session.NewAddSigningPartEvent(nil, nil, false))
	require.Error(t, err)
}
