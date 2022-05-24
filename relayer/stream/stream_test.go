package stream_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/stream"
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
			&prototypes.Int64Value{
				Value: 123,
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
			&prototypes.StringValue{
				// Not quite max size since there's other data along with the Message
				// This may fail if more fields are added and should be decreased.
				Value: string(make([]byte, stream.MAX_MESSAGE_SIZE-300)),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"exceed size",
			&prototypes.StringValue{
				// This plus the other data will exceed the max size.
				Value: string(make([]byte, stream.MAX_MESSAGE_SIZE)),
			},
			errArgs{
				expectPass: false,
				contains:   io.ErrShortBuffer.Error(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msgAny, err := prototypes.MarshalAny(tc.payload)
			require.NoError(t, err)

			// Write/read from buffer
			var buf bytes.Buffer
			err = stream.NewProtoMessageWriter(&buf).WriteMsg(msgAny)
			require.NoError(t, err)

			var msgRead prototypes.Any
			err = stream.NewProtoMessageReader(&buf).ReadMsg(&msgRead)

			if tc.errArgs.expectPass {
				require.NoError(t, err)

				require.True(t, msgAny.Equal(msgRead), "unpacked message should match original")
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

	var msgRes prototypes.StringValue
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

	var msgRes prototypes.StringValue
	err := stream.NewProtoMessageReader(&buf).ReadMsg(&msgRes)

	require.Error(t, err)
	require.ErrorIs(t, io.EOF, err)
}
