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
	"github.com/libp2p/go-libp2p-core/peer"
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
			&types.HelloRequest{
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
			&types.HelloRequest{
				// Not quite max size since there's other data along with the Message
				// This may fail if more fields are added and should be decreased.
				Message: string(make([]byte, stream.MAX_MESSAGE_SIZE-200)),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"exceed size",
			&types.HelloRequest{
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
			// Error if decoding invalid peer ID:
			// "length greater than remaining number of bytes in buffer"
			hostPeerID, err := peer.Decode("16Uiu2HAmTdEddBdw1JVs5tHhqQGaFPkqq64TwppmL2G8fYbZeZei")
			require.NoError(t, err)

			recipients, err := peer.Decode("16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9")
			require.NoError(t, err)

			msg, err := types.NewBroadcastMessage(
				"id",
				tc.payload,
				hostPeerID,
				[]peer.ID{
					recipients,
				},
			)
			require.NoError(t, err)

			// Write/read from buffer
			var buf bytes.Buffer
			err = stream.NewProtoMessageWriter(&buf).WriteMsg(&msg)
			require.NoError(t, err)

			var msgRes types.BroadcastMessage
			err = stream.NewProtoMessageReader(&buf).ReadMsg(&msgRes)

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

	var msgRes types.BroadcastMessage
	err := stream.NewProtoMessageReader(&buf).ReadMsg(&msgRes)

	require.Error(t, err)
	require.ErrorIs(t, io.ErrShortBuffer, err)
}

func TestRead_MaxSizeEmptyData(t *testing.T) {
	var buf bytes.Buffer
	maxSize := stream.MAX_MESSAGE_SIZE

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(maxSize))

	buf.Write(b)

	var msgRes types.BroadcastMessage
	err := stream.NewProtoMessageReader(&buf).ReadMsg(&msgRes)

	require.Error(t, err)
	require.ErrorIs(t, io.EOF, err)
}
