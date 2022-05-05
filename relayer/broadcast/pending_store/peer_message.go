package pending_store

import (
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// MessageWithPeerMetadata is a message with metadata about the peer that sent it.
type MessageWithPeerMetadata struct {
	Message types.Message

	// Not transmitted over wire, added when received.
	PeerID peer.ID
}
