package broadcast

import (
	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// MessageWithPeerMetadata is a message with metadata about the peer that sent it.
type MessageWithPeerMetadata struct {
	Message types.MessageData

	// Not transmitted over wire, added when received.
	PeerID peer.ID
}
