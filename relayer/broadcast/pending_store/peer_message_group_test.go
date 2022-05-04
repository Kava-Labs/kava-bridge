package pending_store_test

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/pending_store"
	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
)

func TestUniquePeerMessages(t *testing.T) {
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
			"valid - 1 item",
			&pending_store.PeerMessageGroup{
				Messages: map[peer.ID]*pending_store.MessageWithPeerMetadata{
					"peer1": {
						BroadcastMessage: types.BroadcastMessage{
							ID: "msg id 1",
							Payload: *mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer1",
					},
				},
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - 2 same items",
			&pending_store.PeerMessageGroup{
				Messages: map[peer.ID]*pending_store.MessageWithPeerMetadata{
					"peer1": {
						BroadcastMessage: types.BroadcastMessage{
							ID: "msg id 1",
							Payload: *mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer1",
					},
					"peer2": {
						BroadcastMessage: types.BroadcastMessage{
							ID: "msg id 1",
							Payload: *mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer2",
					},
				},
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - different message id",
			&pending_store.PeerMessageGroup{
				Messages: map[peer.ID]*pending_store.MessageWithPeerMetadata{
					"peer1": {
						BroadcastMessage: types.BroadcastMessage{
							ID: "msg id 1",
							Payload: *mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer1",
					},
					"peer2": {
						BroadcastMessage: types.BroadcastMessage{
							ID: "msg id 2",
							Payload: *mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer2",
					},
				},
			},
			errArgs{
				expectPass: false,
				// Does not contain full message, non-deterministic map iteration
				// may swap the order.
				contains: "mismatch: \"msg id ",
			},
		},
		{
			"invalid - different payload",
			&pending_store.PeerMessageGroup{
				Messages: map[peer.ID]*pending_store.MessageWithPeerMetadata{
					"peer1": {
						BroadcastMessage: types.BroadcastMessage{
							ID: "msg id 1",
							Payload: *mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer1",
					},
					"peer2": {
						BroadcastMessage: types.BroadcastMessage{
							ID: "msg id 1",
							Payload: *mustMarshalAny(&types.HelloRequest{
								Message: "goodbye world",
							}),
						},
						PeerID: "peer2",
					},
				},
			},
			errArgs{
				expectPass: false,
				contains:   "message payloads do not match from peer",
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
