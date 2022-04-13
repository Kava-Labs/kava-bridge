package broadcast

import (
	"context"
	"sync"

	"github.com/gogo/protobuf/proto"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/kava-labs/kava-bridge/relayer/stream"
	"github.com/kava-labs/kava-bridge/relayer/types"
)

var log = logging.Logger("broadcast")

const (
	ProtocolID  = "/kava-relayer/broadcast/1.0.0"
	ServiceName = "kava-relayer.broadcast"
)

type Broadcast struct {
	host host.Host

	inboundStreamsLock sync.Mutex
	inboundStreams     map[peer.ID]network.Stream

	outboundStreamsLock sync.Mutex
	outboundStreams     map[peer.ID]network.Stream

	// Map of peers
	peersLock sync.Mutex
	peers     map[peer.ID]struct{}

	newPeers chan peer.ID

	// Incoming messages from other peers
	incoming chan *types.MessageData

	// Messages to send to all other peers
	outgoing chan *types.MessageData

	ctx context.Context
}

func NewBroadcast(ctx context.Context, host host.Host) *Broadcast {
	b := &Broadcast{
		host:                host,
		inboundStreamsLock:  sync.Mutex{},
		inboundStreams:      make(map[peer.ID]network.Stream),
		outboundStreamsLock: sync.Mutex{},
		outboundStreams:     make(map[peer.ID]network.Stream),
		peersLock:           sync.Mutex{},
		peers:               make(map[peer.ID]struct{}),
		incoming:            make(chan *types.MessageData),
		outgoing:            make(chan *types.MessageData),
		ctx:                 ctx,
	}

	// Register peer notifications
	b.host.Network().Notify((*BroadcastNotif)(b))

	// Handle new incoming streams
	host.SetStreamHandler(ProtocolID, b.handleNewStream)

	go b.handleNewPeers(ctx)

	return b
}

func (b *Broadcast) handleNewPeers(ctx context.Context) {
	defer func() {
		// Clean up go routines.
		for _, ch := range b.inboundStreams {
			ch.Reset()
		}

		for _, ch := range b.outboundStreams {
			ch.Reset()
		}
	}()

	for {
		select {
		case newPeerID := <-b.newPeers:
			b.handleNewPeer(ctx, newPeerID)
		case <-ctx.Done():
			log.Info("broadcast new peer handler loop shutting down")
			return
		}
	}
}

func (b *Broadcast) SendProtoMessage(
	s network.Stream,
	pb proto.Message,
) error {
	msg, err := types.NewMessageData(pb)
	if err != nil {
		return err
	}

	return stream.NewProtoMessageWriter(s).WriteMsg(&msg)
}

// handleNewPeer opens a new stream with a newly connected peer.
func (b *Broadcast) handleNewPeer(ctx context.Context, pid peer.ID) {
	s, err := b.host.NewStream(b.ctx, pid, ProtocolID)
	if err != nil {
		log.Errorf("failed to open new stream to peer: ", err, pid)

		return
	}

	b.outboundStreamsLock.Lock()
	b.outboundStreams[pid] = s
	b.outboundStreamsLock.Unlock()
}

func (b *Broadcast) handleNewStream(s network.Stream) {
	peer := s.Conn().RemotePeer()

	b.inboundStreamsLock.Lock()
	other, dup := b.inboundStreams[peer]
	if dup {
		log.Debugf("duplicate inbound stream from %s; resetting other stream", peer)

		if err := other.Reset(); err != nil {
			log.Warnf("error resetting other stream: %s", err)
		}
	}
	b.inboundStreams[peer] = s
	b.inboundStreamsLock.Unlock()

	// If there's an error in a stream, remove it from the map.
	defer func() {
		b.inboundStreamsLock.Lock()
		if b.inboundStreams[peer] == s {
			delete(b.inboundStreams, peer)
		}
		b.inboundStreamsLock.Unlock()
	}()

	// Iterate over all messages, unmarshalling all as types.MessageData
	r := stream.NewProtoMessageReader(s)
	for {
		var msg types.MessageData
		err := r.ReadMsg(&msg)

		// TODO: Peer information e.g. peer ID can be attached to message here
		// with a wrapper type. (Not additional fields in types.Message as those
		// are read from the other peer)

		if err != nil {
			log.Errorf("error reading stream message from peer %s: %s", s.Conn().RemotePeer(), err)
			_ = s.Reset()

			return
		}

		select {
		case b.incoming <- &msg:
		case <-b.ctx.Done():
			// Close is useless because the other side isn't reading.
			_ = s.Reset()
			return
		}
	}
}
