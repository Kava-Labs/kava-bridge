package broadcast

import (
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
)

// BroadcasterHook defines the interface for broadcaster hooks. This is not
// exported as it should only be used for testing.
type BroadcasterHook interface {
	// Run before a raw message is broadcasted and can be used to modify the
	// message.
	BeforeBroadcastRawMessage(b *P2PBroadcaster, target peer.ID, pb *proto.Message)
}

// NoOpBroadcasterHook is a broadcasterHook that does nothing.
type noOpBroadcasterHook struct{}

var _ BroadcasterHook = (*noOpBroadcasterHook)(nil)

func (h *noOpBroadcasterHook) BeforeBroadcastRawMessage(b *P2PBroadcaster, target peer.ID, pb *proto.Message) {
}
