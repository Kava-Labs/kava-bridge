package broadcast

import "github.com/kava-labs/kava-bridge/relayer/types"

// BroadcastHandler defines the interface for handling broadcast messages.
type BroadcastHandler interface {
	// RawMessage is called when a raw message from any peer is received.
	RawMessage(msg MessageWithPeerMetadata)
	// ValidatedMessage is called when a message is confirmed to be the same
	// from all peers.
	ValidatedMessage(msg types.MessageData)
	// MismatchMessage is called when a message with the same message ID is
	// different from other peer messages, i.g. a faulty or malicious node.
	MismatchMessage(msg MessageWithPeerMetadata)
}

// NoOpBroadcastHandler is a BroadcastHandler that does nothing.
type NoOpBroadcastHandler struct{}

var _ BroadcastHandler = (*NoOpBroadcastHandler)(nil)

func (h *NoOpBroadcastHandler) RawMessage(msg MessageWithPeerMetadata)      {}
func (h *NoOpBroadcastHandler) ValidatedMessage(msg types.MessageData)      {}
func (h *NoOpBroadcastHandler) MismatchMessage(msg MessageWithPeerMetadata) {}
