package broadcast

import (
	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

type MessageWithPeerMetadata struct {
	Message types.MessageData

	// Not transmitted over wire, added when received.
	PeerID peer.ID
}
