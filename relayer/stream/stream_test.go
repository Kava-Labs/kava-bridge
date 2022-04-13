package stream_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/stream"
	"github.com/kava-labs/kava-bridge/relayer/types"
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
			&types.EchoRequest{
				Message: "hello world",
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"empty",
			&prototypes.Empty{},
			errArgs{
				expectPass: true,
			},
		},
		{
			"large",
			&types.EchoRequest{
				// Not quite max size since there's other data along with the Message
				Message: string(make([]byte, stream.MAX_MESSAGE_SIZE-100)),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"exceed size",
			&types.EchoRequest{
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
			msg, err := types.NewMessageData(tc.payload)
			require.NoError(t, err)

			// Write/read from buffer
			var buf bytes.Buffer
			err = stream.WriteProtoMessage(&buf, &msg)
			require.NoError(t, err)

			var msgRes types.MessageData
			err = stream.ReadProtoMessage(&buf, &msgRes)

			if tc.errArgs.expectPass {
				require.NoError(t, err)

				// Unpack response
				var unpacked prototypes.DynamicAny
				err = msgRes.UnpackPayload(&unpacked)
				require.NoError(t, err)

				require.Equal(t, tc.payload, unpacked.Message, "unpacked message should match original")
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func TestRead_ExceedSize(t *testing.T) {
	var buf bytes.Buffer

	// Max u32
	buf.Write([]byte{255, 255, 255, 255})

	var msgRes types.MessageData
	err := stream.ReadProtoMessage(&buf, &msgRes)

	require.Error(t, err)
	require.ErrorIs(t, io.ErrShortBuffer, err)
}

func TestRead_MaxSizeEmptyData(t *testing.T) {
	var buf bytes.Buffer
	maxSize := stream.MAX_MESSAGE_SIZE

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(maxSize))

	buf.Write(b)

	var msgRes types.MessageData
	err := stream.ReadProtoMessage(&buf, &msgRes)

	require.Error(t, err)
	require.ErrorIs(t, io.EOF, err)
}
