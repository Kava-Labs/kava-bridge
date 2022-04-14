package broadcast

// BroadcastHandler defines the interface for handling broadcast messages.
type BroadcastHandler interface {
	HandleRawMessage(msg *MessageWithPeerMetadata)
}

// NoOpBroadcastHandler is a BroadcastHandler that does nothing.
type NoOpBroadcastHandler struct{}

var _ BroadcastHandler = (*NoOpBroadcastHandler)(nil)

func (h *NoOpBroadcastHandler) HandleRawMessage(msg *MessageWithPeerMetadata) {}
