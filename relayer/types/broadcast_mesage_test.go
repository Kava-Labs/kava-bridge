package types_test

import (
	"testing"
	"time"

	proto "github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func MustNewBroadcastMessage(
	payload proto.Message,
	hostPeerID peer.ID,
	recipients []peer.ID,
	TTLSeconds uint64,
) types.BroadcastMessage {
	msg, err := types.NewBroadcastMessage(payload, hostPeerID, recipients, TTLSeconds)
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
			MustNewBroadcastMessage(&prototypes.Empty{}, testutil.TestPeerIDs[0], []peer.ID{testutil.TestPeerIDs[1]}, 1),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - empty recipients",
			MustNewBroadcastMessage(&prototypes.Empty{}, testutil.TestPeerIDs[0], []peer.ID{}, 1),
			errArgs{
				expectPass: false,
				contains:   types.ErrMsgInsufficientRecipients.Error(),
			},
		},
		{
			"invalid - empty host ID",
			MustNewBroadcastMessage(&prototypes.Empty{}, "", []peer.ID{testutil.TestPeerIDs[0]}, 1),
			errArgs{
				expectPass: false,
				contains:   peer.ErrEmptyPeerID.Error(),
			},
		},
		{
			"invalid - 0 TTL",
			MustNewBroadcastMessage(&prototypes.Empty{}, testutil.TestPeerIDs[0], []peer.ID{testutil.TestPeerIDs[1]}, 0),
			errArgs{
				expectPass: false,
				contains:   types.ErrMsgTTLTooShort.Error(),
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

func TestMessageExpired(t *testing.T) {
	msg := MustNewBroadcastMessage(&prototypes.Empty{}, testutil.TestPeerIDs[0], testutil.TestPeerIDs[1:2], 1)
	require.False(t, msg.Expired())

	// 2 seconds to be > TTL and not >= TTL
	time.Sleep(2 * time.Second)
	require.True(t, msg.Expired())
	require.ErrorIs(t, msg.Validate(), types.ErrMsgExpired)
}

func TestMessageExpired_Future(t *testing.T) {
	// Message 5 seconds in future, ie. peers with out of sync times
	msg := types.BroadcastMessage{
		MustNewBroadcastMessageID(),
		testutil.TestPeerIDs[:2],
		prototypes.Any{},
		time.Now().Add(time.Second),
		1,
	}
	require.False(t, msg.Expired(), "duration since created should not underflow")

	// 2 seconds to be > TTL and not >= TTL
	time.Sleep(2 * time.Second)
	require.True(t, msg.Expired())
	require.ErrorIs(t, msg.Validate(), types.ErrMsgExpired)
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
			msg, err := types.NewBroadcastMessage(tc.payload, testutil.TestPeerIDs[0], nil, 1)
			require.NoError(t, err)

			var unpacked prototypes.DynamicAny
			err = msg.UnpackPayload(&unpacked)
			require.NoError(t, err)

			require.Equal(t, tc.payload, unpacked.Message, "unpacked message should match original")
		})
	}
}

func MustNewBroadcastMessageID() string {
	id, err := types.NewBroadcastMessageID()
	if err != nil {
		panic(err)
	}
	return id
}
