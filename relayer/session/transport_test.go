package session_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/kava-labs/kava-bridge/relayer/session"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/require"
)

func TestSessionTransport_SendBroadcast(t *testing.T) {
	tp := session.NewSessionTransport(
		&broadcast.NoOpBroadcaster{},
		types.AggregateSigningSessionID{1, 2, 3},
		mp_tss.NewPartyIDStore(),
		testutil.TestPeerIDs,
	)

	err := tp.Send(
		context.Background(),
		[]byte{1, 2, 3},
		&tss.MessageRouting{
			From:                    nil,
			To:                      nil,
			IsBroadcast:             true,
			IsToOldCommittee:        false,
			IsToOldAndNewCommittees: false,
		},
		false, // not resharing
	)

	require.NoError(t, err)
}

func TestSessionTransport_SendDirect(t *testing.T) {
	partyIDStore := mp_tss.NewPartyIDStore()

	partyIDStore.AddPeers([]*mp_tss.PeerMetadata{
		{
			PeerID:  testutil.TestPeerIDs[0],
			Moniker: "one",
		},
		{
			PeerID:  testutil.TestPeerIDs[1],
			Moniker: "two",
		},
	})

	tp := session.NewSessionTransport(
		&broadcast.NoOpBroadcaster{},
		types.AggregateSigningSessionID{1, 2, 3},
		partyIDStore,
		testutil.TestPeerIDs,
	)

	toPartyID, err := partyIDStore.GetPartyID(testutil.TestPeerIDs[0])
	require.NoError(t, err)

	err = tp.Send(
		context.Background(),
		[]byte{1, 2, 3},
		&tss.MessageRouting{
			From: nil,
			To: []*tss.PartyID{
				toPartyID,
			},
			IsBroadcast:             false,
			IsToOldCommittee:        false,
			IsToOldAndNewCommittees: false,
		},
		false, // not resharing
	)

	require.NoError(t, err)
}

func TestSessionTransport_SendNotFound(t *testing.T) {
	// Empty partyIDStore
	partyIDStore := mp_tss.NewPartyIDStore()

	tp := session.NewSessionTransport(
		&broadcast.NoOpBroadcaster{},
		types.AggregateSigningSessionID{1, 2, 3},
		partyIDStore,
		testutil.TestPeerIDs,
	)

	err := tp.Send(
		context.Background(),
		[]byte{1, 2, 3},
		&tss.MessageRouting{
			From: nil,
			To: []*tss.PartyID{
				tss.NewPartyID("id that doesn't exist", "id", big.NewInt(123)),
			},
			IsBroadcast:             false,
			IsToOldCommittee:        false,
			IsToOldAndNewCommittees: false,
		},
		false, // not resharing
	)

	require.Error(t, err)
}
