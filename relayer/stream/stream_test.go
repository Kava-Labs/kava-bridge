package stream_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/stream"
	"github.com/kava-labs/kava-bridge/relayer/stream/types"
	"github.com/stretchr/testify/require"
)

func TestReadWrite(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		payload proto.Message
		errArgs errArgs
	}{
		{
			"regular",
			&types.Echo{
				Message: "hello world",
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"empty",
			&types.Echo{},
			errArgs{
				expectPass: true,
			},
		},
		{
			"longish",
			&types.Echo{
				Message: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"large",
			&types.Echo{
				// Not quite max size since there's other data along with the Message
				Message: string(make([]byte, stream.MAX_MESSAGE_SIZE-100)),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"exceed size",
			&types.Echo{
				// This plus the other data will exceed the max size.
				Message: string(make([]byte, stream.MAX_MESSAGE_SIZE)),
			},
			errArgs{
				expectPass: false,
				contains:   io.ErrShortBuffer.Error(),
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

			if tc.errArgs.expectPass {
				require.NoError(t, err)

				// Unpack response
				var unpacked prototypes.DynamicAny
				err = prototypes.UnmarshalAny(msgRes.Payload, &unpacked)
				require.NoError(t, err)

				require.Equal(t, tc.payload, unpacked.Message, "unpacked message should match original")
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}
