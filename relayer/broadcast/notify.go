package broadcast

import (
	"github.com/libp2p/go-libp2p-core/network"
	ma "github.com/multiformats/go-multiaddr"
)

var _ network.Notifiee = (*BroadcastNotif)(nil)

type BroadcastNotif P2PBroadcaster

func (p *BroadcastNotif) OpenedStream(n network.Network, s network.Stream) {
}

func (p *BroadcastNotif) ClosedStream(n network.Network, s network.Stream) {
}

func (p *BroadcastNotif) Connected(n network.Network, c network.Conn) {
	// Ignore transient connections
	if c.Stat().Transient {
		p.log.Debugf("Ignoring transient connection: %s", c.RemotePeer())
		return
	}

	go func() {
		p.log.Debugf("Notifying of new peer: %s", c.RemotePeer().ShortString())
		p.newPeers <- c.RemotePeer()
	}()
}

func (p *BroadcastNotif) Disconnected(n network.Network, c network.Conn) {
}

func (p *BroadcastNotif) Listen(n network.Network, _ ma.Multiaddr) {
}

func (p *BroadcastNotif) ListenClose(n network.Network, _ ma.Multiaddr) {
}
