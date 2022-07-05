package types_test

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/stretchr/testify/require"
)

func TestGetSessionID_Invalid(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name          string
		messages      types.JoinSessionMessages
		wantSessionID []byte
		errArgs       errArgs
	}{
		{
			"invalid - empty",
			types.JoinSessionMessages{},
			nil,
			errArgs{
				expectPass: false,
				contains:   "empty join session messages",
			},
		},
		{
			"invalid - 1 msg",
			types.JoinSessionMessages{
				types.NewJoinSigningSessionMessage(
					"peerid",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{1},
				),
			},
			nil,
			errArgs{
				expectPass: false,
				contains:   "not enough peers to select participants, 1 (peers) < 2 (t + 1)",
			},
		},
		{
			"invalid - session message type",
			types.JoinSessionMessages{
				types.NewJoinKeygenSessionMessage(
					"peerid",
					types.KeygenSessionID{0},
				),
			},
			nil,
			errArgs{
				expectPass: false,
				contains:   "invalid join session type",
			},
		},
		{
			"invalid - 2 msgs same peerid",
			types.JoinSessionMessages{
				types.NewJoinSigningSessionMessage(
					"peerid",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{1},
				),
				types.NewJoinSigningSessionMessage(
					"peerid",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{2},
				),
			},
			nil,
			errArgs{
				expectPass: false,
				contains:   "duplicate peer ID",
			},
		},
		{
			"invalid - different txhash",
			types.JoinSessionMessages{
				types.NewJoinSigningSessionMessage(
					"peerid1",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{1},
				),
				types.NewJoinSigningSessionMessage(
					"peerid2",
					common.BytesToHash([]byte{1}),
					types.SigningSessionIDPart{2},
				),
				types.NewJoinSigningSessionMessage(
					"peerid3",
					common.BytesToHash([]byte{1}),
					types.SigningSessionIDPart{3},
				),
			},
			nil,
			errArgs{
				expectPass: false,
				contains:   "different tx hashes",
			},
		},
		{
			"invalid - duplicate id part",
			types.JoinSessionMessages{
				types.NewJoinSigningSessionMessage(
					"peerid1",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{1},
				),
				types.NewJoinSigningSessionMessage(
					"peerid2",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{2},
				),
				types.NewJoinSigningSessionMessage(
					"peerid3",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{2},
				),
			},
			nil,
			errArgs{
				expectPass: false,
				contains:   "duplicate peer session ID part",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sessionID, _, err := tc.messages.GetSessionID(1)

			if tc.errArgs.expectPass {
				require.NoError(t, err)
				require.Equal(t, tc.wantSessionID, sessionID.Bytes())
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func TestGetSessionID_Order(t *testing.T) {
	msgs := types.JoinSessionMessages{
		types.NewJoinSigningSessionMessage(
			"peer1",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{1},
		),
		types.NewJoinSigningSessionMessage(
			"peer2",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{2},
		),
		types.NewJoinSigningSessionMessage(
			"peer3",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{3},
		),
		types.NewJoinSigningSessionMessage(
			"peer4",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{4},
		),
		types.NewJoinSigningSessionMessage(
			"peer5",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{5},
		),
	}

	expectedSessionID := AppendSlices(
		types.SigningSessionIDPart{1}.Bytes(),
		types.SigningSessionIDPart{2}.Bytes(),
		types.SigningSessionIDPart{3}.Bytes(),
		types.SigningSessionIDPart{4}.Bytes(),
		types.SigningSessionIDPart{5}.Bytes(),
	)

	for i := 0; i < len(msgs)*2; i++ {
		// Shuffle order
		rand.Shuffle(len(msgs), func(i, j int) {
			msgs[i], msgs[j] = msgs[j], msgs[i]
		})

		threshold := 4

		// Make sure each shuffled order produces the same result
		sessionID, participantPeerIDs, err := msgs.GetSessionID(threshold)
		require.NoError(t, err)

		require.Len(t, participantPeerIDs, threshold+1, "there should be t+1 participants")
		require.NoError(t, sessionID.Validate(), "session id should be valid")
		require.Equal(t, expectedSessionID, sessionID.Bytes(), "session id should match expected")
	}
}

func TestGetSessionID_InvalidThreshold(t *testing.T) {
	msgs := types.JoinSessionMessages{
		types.NewJoinSigningSessionMessage(
			"peer1",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{1},
		),
		types.NewJoinSigningSessionMessage(
			"peer2",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{2},
		),
		types.NewJoinSigningSessionMessage(
			"peer3",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{3},
		),
		types.NewJoinSigningSessionMessage(
			"peer4",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{4},
		),
		types.NewJoinSigningSessionMessage(
			"peer5",
			common.BytesToHash([]byte{0}),
			types.SigningSessionIDPart{5},
		),
	}

	_, _, err := msgs.GetSessionID(0)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid threshold")
}
