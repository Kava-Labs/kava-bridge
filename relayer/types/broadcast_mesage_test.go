package types_test

import (
	"testing"

	proto "github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func MustNewBroadcastMessage(id string, payload proto.Message, recipients []peer.ID) types.BroadcastMessage {
	msg, err := types.NewBroadcastMessage(id, payload, recipients)
	if err != nil {
		panic(err)
	}
	return msg
}

func TestValidateMessage(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		message types.BroadcastMessage
		errArgs errArgs
	}{
		{
			"valid",
			MustNewBroadcastMessage("hi", &prototypes.Empty{}, []peer.ID{peer.ID("QmQQGdG9Ybz2qXNmzXo9pT9VZpvZ2Zcq2R6zQmXo9FtZz")}),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - empty id",
			MustNewBroadcastMessage("", &prototypes.Empty{}, []peer.ID{peer.ID("QmQQGdG9Ybz2qXNmzXo9pT9VZpvZ2Zcq2R6zQmXo9FtZz")}),
			errArgs{
				expectPass: false,
				contains:   "message ID is empty",
			},
		},
		{
			"invalid - empty recipients",
			MustNewBroadcastMessage("", &prototypes.Empty{}, []peer.ID{}),
			errArgs{
				expectPass: false,
				contains:   "message ID is empty",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.message.Validate()

			if tc.errArgs.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func TestMarshalUnmarshalPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload proto.Message
	}{
		{
			"regular",
			&types.HelloRequest{
				Message: "hello world",
			},
		},
		{
			"empty",
			&prototypes.Empty{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := types.NewBroadcastMessage("an id", tc.payload, nil)
			require.NoError(t, err)

			var unpacked prototypes.DynamicAny
			err = msg.UnpackPayload(&unpacked)
			require.NoError(t, err)

			require.Equal(t, tc.payload, unpacked.Message, "unpacked message should match original")
		})
	}
}
