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

type Broadcaster struct {
	host host.Host

	inboundStreamsLock sync.Mutex
	inboundStreams     map[peer.ID]network.Stream

	outboundStreamsLock sync.Mutex
	outboundStreams     map[peer.ID]network.Stream

	// Map of peers
	peersLock sync.RWMutex
	peers     map[peer.ID]struct{}

	// Notifications of new peers
	newPeers chan peer.ID

	// Raw incomingRaw messages from other peers, NOT verified. Will contain
	// duplicate messages from different peers.
	incomingRaw chan *RPC

	// Messages to send to all other peers
	outgoing chan *types.MessageData

	// External event handler
	handler BroadcastHandler

	ctx context.Context
}

// NewBroadcaster returns a new Broadcaster
func NewBroadcaster(ctx context.Context, host host.Host, options ...BroadcasterOption) (*Broadcaster, error) {
	b := &Broadcaster{
		host:                host,
		inboundStreamsLock:  sync.Mutex{},
		inboundStreams:      make(map[peer.ID]network.Stream),
		outboundStreamsLock: sync.Mutex{},
		outboundStreams:     make(map[peer.ID]network.Stream),
		peersLock:           sync.RWMutex{},
		peers:               make(map[peer.ID]struct{}),
		newPeers:            make(chan peer.ID),
		incomingRaw:         make(chan *RPC),
		outgoing:            make(chan *types.MessageData),
		handler:             nil,
		ctx:                 ctx,
	}

	for _, opt := range options {
		err := opt(b)
		if err != nil {
			return nil, err
		}
	}

	// Register peer notifications
	b.host.Network().Notify((*BroadcastNotif)(b))

	// Handle new incoming streams
	host.SetStreamHandler(ProtocolID, b.handleNewStream)

	go b.handleNewPeers(ctx)

	return b, nil
}

// GetPeerCount returns the number of peers connected to the broadcaster.
func (b *Broadcaster) GetPeerCount() int {
	b.peersLock.RLock()
	defer b.peersLock.RUnlock()

	return len(b.peers)
}

func (b *Broadcaster) handleNewPeers(ctx context.Context) {
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
		case newMsg := <-b.incomingRaw:
			b.handleIncomingMsg(newMsg)
		case <-ctx.Done():
			log.Info("broadcast new peer handler loop shutting down")
			return
		}
	}
}

// SendProtoMessage sends a proto message to all connected peers.
func (b *Broadcaster) SendProtoMessage(
	pb proto.Message,
) error {
	msg, err := types.NewMessageData(pb)
	if err != nil {
		return err
	}

	// TODO: Make concurrent
	for _, ch := range b.outboundStreams {
		// NewUint32DelimitedWriter has an internal buffer, bufio.NewWriter()
		// should not be necessary.
		if err := stream.NewProtoMessageWriter(ch).WriteMsg(&msg); err != nil {
			return err
		}
	}

	return nil
}

func (b *Broadcaster) handleIncomingMsg(r *RPC) {

}

// handleNewPeer opens a new stream with a newly connected peer.
func (b *Broadcaster) handleNewPeer(ctx context.Context, pid peer.ID) {
	s, err := b.host.NewStream(b.ctx, pid, ProtocolID)
	if err != nil {
		log.Errorf("failed to open new stream to peer: ", err, pid)

		return
	}

	log.Debugf("opened new stream to peer: %s", pid)

	b.outboundStreamsLock.Lock()
	b.outboundStreams[pid] = s
	b.outboundStreamsLock.Unlock()

	b.peersLock.Lock()
	b.peers[pid] = struct{}{}
	b.peersLock.Unlock()
}

// handleNewStream handles a new incoming stream, initiated when a peer is connected.
func (b *Broadcaster) handleNewStream(s network.Stream) {
	log.Debugf("incoming stream from peer: %s", s.Conn().RemotePeer())
	peer := s.Conn().RemotePeer()

	// Ensure only 1 incoming stream from the peer
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

	log.Debugf("starting stream processor for peer: %s", s.Conn().RemotePeer())
	// Iterate over all messages, unmarshalling all as types.MessageData
	r := stream.NewProtoMessageReader(s)
	for {
		var msg types.MessageData
		if err := r.ReadMsg(&msg); err != nil {
			log.Errorf("error reading stream message from peer %s: %s", s.Conn().RemotePeer(), err)
			_ = s.Reset()

			return
		}

		// Attach additional peer metadata to the message
		rpc := RPC{
			Message: msg,
			PeerID:  s.Conn().RemotePeer(),
		}

		select {
		case b.incomingRaw <- &rpc:
		case <-b.ctx.Done():
			// Close is useless because the other side isn't reading.
			_ = s.Reset()
			return
		}
	}
}
