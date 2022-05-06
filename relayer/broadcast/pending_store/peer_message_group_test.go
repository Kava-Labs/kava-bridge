package pending_store_test

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/pending_store"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func TestUniquePeerMessages(t *testing.T) {
	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	msg := types.BroadcastMessage{
		ID: msgID,
		Payload: *mustMarshalAny(&types.HelloRequest{
			Message: "hello world",
		}),
	}

	msgHash, err := msg.Hash()
	require.NoError(t, err)

	invalidMsgHash := types.BroadcastMessageHash{}

	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		group   *pending_store.PeerMessageGroup
		errArgs errArgs
	}{
		{
			"valid - empty",
			pending_store.NewPeerMessageGroup(),
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - with broadcasted message + no hashes",
			&pending_store.PeerMessageGroup{
				BroadcastedMessage:         msg,
				BroadcastedMessageReceived: true,
				PeerMessageHashes:          map[peer.ID]types.BroadcastMessageHash{},
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - broadcasted message + 1 hash",
			&pending_store.PeerMessageGroup{
				BroadcastedMessage:         msg,
				BroadcastedMessageReceived: true,
				PeerMessageHashes: map[peer.ID]types.BroadcastMessageHash{
					testutil.TestPeerIDs[0]: msgHash,
				},
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - broadcasted message + 2 hashes",
			&pending_store.PeerMessageGroup{
				BroadcastedMessage:         msg,
				BroadcastedMessageReceived: true,
				PeerMessageHashes: map[peer.ID]types.BroadcastMessageHash{
					testutil.TestPeerIDs[0]: msgHash,
					testutil.TestPeerIDs[1]: msgHash,
				},
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - broadcasted message + mismatch hashes",
			&pending_store.PeerMessageGroup{
				BroadcastedMessage:         msg,
				BroadcastedMessageReceived: true,
				PeerMessageHashes: map[peer.ID]types.BroadcastMessageHash{
					testutil.TestPeerIDs[0]: msgHash,
					testutil.TestPeerIDs[1]: invalidMsgHash,
				},
			},
			errArgs{
				expectPass: false,
				contains:   "group contains invalid hash for peer",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.group.Validate()

			if tc.errArgs.expectPass {
				require.NoError(t, err)

			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func mustMarshalAny(pb proto.Message) *prototypes.Any {
	any, err := prototypes.MarshalAny(pb)
	if err != nil {
		panic(err)
	}

	return any
}
