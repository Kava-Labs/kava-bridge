package types_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		giveMsg types.JoinSessionMessage
		errArgs errArgs
	}{
		{
			name:    "invalid - empty",
			giveMsg: types.JoinSessionMessage{},
			errArgs: errArgs{
				expectPass: false,
				contains:   "invalid session type: <nil>",
			},
		},
		{
			name: "invalid - empty keygen session message",
			giveMsg: types.JoinSessionMessage{
				Session: &types.JoinSessionMessage_JoinKeygenSessionMessage{
					JoinKeygenSessionMessage: &types.JoinKeygenSessionMessage{},
				},
			},
			errArgs: errArgs{
				expectPass: false,
				contains:   "keygen session ID length incorrect: expected 32, got 0",
			},
		},
		{
			name: "invalid - empty signing session message",
			giveMsg: types.JoinSessionMessage{
				Session: &types.JoinSessionMessage_JoinSigningSessionMessage{
					JoinSigningSessionMessage: &types.JoinSigningSessionMessage{},
				},
			},
			errArgs: errArgs{
				expectPass: false,
				contains:   "invalid tx hash length: expected 32, got 0",
			},
		},
		{
			name: "invalid - empty peer session id part",
			giveMsg: types.JoinSessionMessage{
				Session: &types.JoinSessionMessage_JoinSigningSessionMessage{
					JoinSigningSessionMessage: &types.JoinSigningSessionMessage{
						TxHash: make([]byte, common.HashLength),
					},
				},
			},
			errArgs: errArgs{
				expectPass: false,
				contains:   "invalid peer session ID part length: expected 32, got 0",
			},
		},
		{
			name: "invalid - empty resharing session message",
			giveMsg: types.JoinSessionMessage{
				Session: &types.JoinSessionMessage_JoinResharingSessionMessage{
					JoinResharingSessionMessage: &types.JoinReSharingSessionMessage{},
				},
			},
			errArgs: errArgs{
				expectPass: false,
				contains:   "resharing session ID length incorrect: expected 32, got 0",
			},
		},
		// Valid
		{
			name:    "valid - keygen",
			giveMsg: types.NewJoinKeygenSessionMessage(types.KeygenSessionID{}),
			errArgs: errArgs{
				expectPass: true,
			},
		},
		{
			name:    "valid - signing",
			giveMsg: types.NewJoinSigningSessionMessage(common.Hash{}, types.PeerSessionIDPart{}),
			errArgs: errArgs{
				expectPass: true,
			},
		},
		{
			name:    "valid - resharing",
			giveMsg: types.NewJoinReSharingSessionMessage(types.ReSharingSessionID{}),
			errArgs: errArgs{
				expectPass: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.giveMsg.ValidateBasic()

			if tt.errArgs.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errArgs.contains)
			}
		})
	}
}
