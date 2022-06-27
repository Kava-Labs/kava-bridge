package broadcast

import (
	"context"

	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

type NoOpBroadcaster struct{}

var _ Broadcaster = (*NoOpBroadcaster)(nil)

func (b *NoOpBroadcaster) BroadcastMessage(
	ctx context.Context,
	pb types.PeerMessage,
	recipients []peer.ID,
	TTLSeconds uint64,
) error {
	return nil
}
