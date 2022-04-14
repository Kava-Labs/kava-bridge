package broadcast

import "github.com/kava-labs/kava-bridge/relayer/types"

// BroadcastHandler defines the interface for handling broadcast messages.
type BroadcastHandler interface {
	HandleRawMessage(msg *MessageWithPeerMetadata)
	HandleValidatedMessage(msg *types.MessageData)
}

// NoOpBroadcastHandler is a BroadcastHandler that does nothing.
type NoOpBroadcastHandler struct{}

var _ BroadcastHandler = (*NoOpBroadcastHandler)(nil)

func (h *NoOpBroadcastHandler) HandleRawMessage(msg *MessageWithPeerMetadata) {}
func (h *NoOpBroadcastHandler) HandleValidatedMessage(msg *types.MessageData) {}
