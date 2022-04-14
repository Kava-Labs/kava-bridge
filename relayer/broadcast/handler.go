package broadcast

import "github.com/kava-labs/kava-bridge/relayer/types"

// BroadcastHandler defines the interface for handling broadcast messages.
type BroadcastHandler interface {
	// HandleRawMessage is called when a raw message from any peer is received.
	HandleRawMessage(msg *MessageWithPeerMetadata)
	// HandleValidatedMessage is called when a message is confirmed to be valid
	// from all peers.
	HandleValidatedMessage(msg *types.MessageData)
}

// NoOpBroadcastHandler is a BroadcastHandler that does nothing.
type NoOpBroadcastHandler struct{}

var _ BroadcastHandler = (*NoOpBroadcastHandler)(nil)

func (h *NoOpBroadcastHandler) HandleRawMessage(msg *MessageWithPeerMetadata) {}
func (h *NoOpBroadcastHandler) HandleValidatedMessage(msg *types.MessageData) {}
