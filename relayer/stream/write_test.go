package stream_test

import (
	"bytes"
	"testing"

	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/stream"
	"github.com/kava-labs/kava-bridge/relayer/stream/types"
	"github.com/stretchr/testify/require"
)

func TestReadWrite(t *testing.T) {
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
			&types.Echo{},
		},
		{
			"long",
			&types.Echo{
				Message: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			any, err := prototypes.MarshalAny(tc.payload)
			require.NoError(t, err)

			msg := &types.MessageData{
				Payload: any,
			}

			// Write/read from buffer
			var buf bytes.Buffer
			err = stream.WriteProtoMessage(&buf, msg)
			require.NoError(t, err)

			var msgRes types.MessageData
			err = stream.ReadProtoMessage(&buf, &msgRes)
			require.NoError(t, err)

			// Unpack response
			var unpacked prototypes.DynamicAny
			err = prototypes.UnmarshalAny(msgRes.Payload, &unpacked)
			require.NoError(t, err)

			require.Equal(t, tc.payload, unpacked.Message, "unpacked message should match original")
		})
	}
}
