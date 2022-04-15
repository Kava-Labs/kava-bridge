package broadcast

import "github.com/kava-labs/kava-bridge/relayer/types"

// BroadcastHandler defines the interface for handling broadcast messages.
type BroadcastHandler interface {
	// RawMessage is called when a raw message from any peer is received.
	RawMessage(msg *MessageWithPeerMetadata)
	// ValidatedMessage is called when a message is confirmed to be valid
	// from all peers.
	ValidatedMessage(msg *types.MessageData)
}

// NoOpBroadcastHandler is a BroadcastHandler that does nothing.
type NoOpBroadcastHandler struct{}

var _ BroadcastHandler = (*NoOpBroadcastHandler)(nil)

func (h *NoOpBroadcastHandler) RawMessage(msg *MessageWithPeerMetadata) {}
func (h *NoOpBroadcastHandler) ValidatedMessage(msg *types.MessageData) {}
