package broadcast_test

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
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
		group   *broadcast.PeerMessageGroup
		errArgs errArgs
	}{
		{
			"valid - empty",
			broadcast.NewPeerMessageGroup(),
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - 1 item",
			&broadcast.PeerMessageGroup{
				Messages: map[peer.ID]*broadcast.MessageWithPeerMetadata{
					"peer1": {
						Message: types.MessageData{
							ID: "msg id 1",
							Payload: mustMarshalAny(&types.HelloRequest{
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
			&broadcast.PeerMessageGroup{
				Messages: map[peer.ID]*broadcast.MessageWithPeerMetadata{
					"peer1": {
						Message: types.MessageData{
							ID: "msg id 1",
							Payload: mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer1",
					},
					"peer2": {
						Message: types.MessageData{
							ID: "msg id 1",
							Payload: mustMarshalAny(&types.HelloRequest{
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
			&broadcast.PeerMessageGroup{
				Messages: map[peer.ID]*broadcast.MessageWithPeerMetadata{
					"peer1": {
						Message: types.MessageData{
							ID: "msg id 1",
							Payload: mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer1",
					},
					"peer2": {
						Message: types.MessageData{
							ID: "msg id 2",
							Payload: mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer2",
					},
				},
			},
			errArgs{
				expectPass: false,
				contains:   "mismatch: \"msg id 2\" != \"msg id 1\"",
			},
		},
		{
			"invalid - different payload",
			&broadcast.PeerMessageGroup{
				Messages: map[peer.ID]*broadcast.MessageWithPeerMetadata{
					"peer1": {
						Message: types.MessageData{
							ID: "msg id 1",
							Payload: mustMarshalAny(&types.HelloRequest{
								Message: "hello world",
							}),
						},
						PeerID: "peer1",
					},
					"peer2": {
						Message: types.MessageData{
							ID: "msg id 1",
							Payload: mustMarshalAny(&types.HelloRequest{
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
