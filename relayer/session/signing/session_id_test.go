package signing_test

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/kava-labs/kava-bridge/relayer/session/signing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAggregateSigningSessionID_Invalid(t *testing.T) {
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
				contains:   "no join messages provided",
			},
		},
		{
			"valid - 1 item",
			types.JoinSessionMessages{
				types.NewJoinSigningSessionMessage(
					"peerid",
					common.BytesToHash([]byte{0}),
					types.SigningSessionIDPart{1},
				),
			},
			types.SigningSessionIDPart{1}.Bytes(),
			errArgs{
				expectPass: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sessionID, err := signing.NewAggregateSigningSessionID(tc.messages)

			if tc.errArgs.expectPass {
				require.NoError(t, err)
				require.NotNil(t, sessionID)
				require.Equal(t, tc.wantSessionID, sessionID.Bytes())
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func TestGetAggregateSigningSessionID_Order(t *testing.T) {
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

		// Make sure each shuffled order produces the same result
		sessionID, err := signing.NewAggregateSigningSessionID(msgs)
		require.NoError(t, err)

		require.True(t, sessionID.Validate())

		require.Equal(t, expectedSessionID, sessionID.Bytes())
	}
}

func TestIsPeerParticipant(t *testing.T) {
	sessionID := signing.AggregateSigningSessionID(
		AppendSlices(
			types.SigningSessionIDPart{1}.Bytes(),
			types.SigningSessionIDPart{2}.Bytes(),
			types.SigningSessionIDPart{3}.Bytes(),
			types.SigningSessionIDPart{4}.Bytes(),
			types.SigningSessionIDPart{5}.Bytes(),
		),
	)

	require.True(t, sessionID.Validate())

	for i := 1; i <= 10; i++ {
		isParticipant := sessionID.IsPeerParticipant(types.SigningSessionIDPart{byte(i)})

		if i <= 5 {
			assert.True(t, isParticipant)
		} else {
			assert.False(t, isParticipant)
		}
	}
}

// AppendSlices concatenates all the given slices.
func AppendSlices(slices ...[]byte) []byte {
	var result []byte
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
