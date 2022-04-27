package types_test

import (
	"testing"

	proto "github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/stretchr/testify/require"
)

func MustNewMessageData(id string, payload proto.Message) types.MessageData {
	msg, err := types.NewMessageData(id, payload)
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
		message types.MessageData
		errArgs errArgs
	}{
		{
			"valid",
			MustNewMessageData("hi", &prototypes.Empty{}),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - empty id",
			MustNewMessageData("", &prototypes.Empty{}),
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
			&types.Echo{
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
			msg, err := types.NewMessageData("an id", tc.payload)
			require.NoError(t, err)

			var unpacked prototypes.DynamicAny
			err = msg.UnpackPayload(&unpacked)
			require.NoError(t, err)

			require.Equal(t, tc.payload, unpacked.Message, "unpacked message should match original")
		})
	}
}
