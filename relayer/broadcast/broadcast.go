package broadcast

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/sync/errgroup"

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

// Broadcaster is a reliable broadcaster to ensure that all connected peers
// receive the same message.
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

	// Raw incoming messages from other peers, NOT verified. Will contain
	// duplicate messages from different peers to be validated.
	incomingRawMessages chan *MessageWithPeerMetadata

	// Messages that all peers have confirmed to have received. Does not contain
	// any peer specific data as it does not originate from any specific peer.
	incomingValidatedMessages chan *types.MessageData

	// Messages that have been received but have not been validated to be same
	// from all peers.
	pendingMessagesLock sync.Mutex
	// TODO: Remove messages after expired.
	pendingMessages map[string]*PeerMessageGroup

	// Messages to send to all other peers
	outgoing chan *types.MessageData

	// External event handler
	handler BroadcastHandler

	ctx context.Context
}

// NewBroadcaster returns a new Broadcaster
func NewBroadcaster(
	ctx context.Context,
	host host.Host,
	options ...BroadcasterOption,
) (*Broadcaster, error) {
	b := &Broadcaster{
		host:                      host,
		inboundStreamsLock:        sync.Mutex{},
		inboundStreams:            make(map[peer.ID]network.Stream),
		outboundStreamsLock:       sync.Mutex{},
		outboundStreams:           make(map[peer.ID]network.Stream),
		peersLock:                 sync.RWMutex{},
		peers:                     make(map[peer.ID]struct{}),
		newPeers:                  make(chan peer.ID),
		incomingRawMessages:       make(chan *MessageWithPeerMetadata),
		incomingValidatedMessages: make(chan *types.MessageData, 1),
		pendingMessagesLock:       sync.Mutex{},
		pendingMessages:           make(map[string]*PeerMessageGroup),
		outgoing:                  make(chan *types.MessageData),
		handler:                   &NoOpBroadcastHandler{},
		ctx:                       ctx,
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

	go b.processLoop(ctx)

	return b, nil
}

// GetPeerCount returns the number of peers connected to the broadcaster. This
// does not include the current peer.
func (b *Broadcaster) GetPeerCount() int {
	b.peersLock.RLock()
	defer b.peersLock.RUnlock()

	return len(b.peers)
}

// processLoop handles incoming channel inputs.
func (b *Broadcaster) processLoop(ctx context.Context) {
	defer func() {
		b.inboundStreamsLock.Lock()
		b.outboundStreamsLock.Lock()

		// Close all streams, errors not important
		for _, ch := range b.inboundStreams {
			_ = ch.Reset()
		}

		for _, ch := range b.outboundStreams {
			_ = ch.Reset()
		}

		b.inboundStreamsLock.Unlock()
		b.outboundStreamsLock.Unlock()
	}()

	for {
		// TODO: Are these fine in the same loop? Not handled concurrently.

		select {
		case newPeerID := <-b.newPeers:
			b.handleNewPeer(ctx, newPeerID)
		case newRawMsg := <-b.incomingRawMessages:
			b.handleIncomingRawMsg(newRawMsg)
		case newValidatedMsg := <-b.incomingValidatedMessages:
			b.handleIncomingValidatedMsg(newValidatedMsg)
		case <-ctx.Done():
			log.Info("broadcast handler loop shutting down")
			return
		}
	}
}

// -----------------------------------------------------------------------------
// Peer messages handling

// BroadcastMessage marshals the proto.Message as Any, wraps it in MessageData,
// and it to all connected peers.
func (b *Broadcaster) BroadcastMessage(
	ctx context.Context,
	messageID string,
	pb proto.Message,
) error {
	// Wrap the proto message in the MessageData type.
	msg, err := types.NewMessageData(messageID, pb)
	if err != nil {
		return err
	}

	// Add the message to the pending messages map to keep track of responses
	// and to prevent re-broadcasting.
	b.pendingMessagesLock.Lock()
	// Could manually lock before broadcastRawMessage
	defer b.pendingMessagesLock.Unlock()

	_, found := b.pendingMessages[msg.ID]
	if found {
		return fmt.Errorf("cannot broadcast message that is already pending: %v", msg.ID)
	}
	b.pendingMessages[msg.ID] = NewPeerMessageGroup()

	if err := msg.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	return b.broadcastRawMessage(ctx, &msg)
}

// broadcastRawMessage sends a proto message to all connected peers without any
// marshalling or wrapping.
func (b *Broadcaster) broadcastRawMessage(
	ctx context.Context,
	pb proto.Message,
) error {
	log.Debugf("broadcast sending raw proto message: %s", pb.String())

	b.outboundStreamsLock.Lock()
	defer b.outboundStreamsLock.Unlock()

	// Run writes to peers in parallel
	g, _ := errgroup.WithContext(ctx)

	for _, ch := range b.outboundStreams {
		func(ch network.Stream) {
			g.Go(func() error {
				log.Debugf("sending message to peer: %s", ch.Conn().RemotePeer())

				// NewUint32DelimitedWriter has an internal buffer, bufio.NewWriter()
				// should not be necessary.
				return stream.NewProtoMessageWriter(ch).WriteMsg(pb)
			})
		}(ch)
	}

	return g.Wait()
}

// handleIncomingRawMsg handles all raw messages from other peers. This is
// before messages are verified to be received from all peers.
func (b *Broadcaster) handleIncomingRawMsg(msg *MessageWithPeerMetadata) {
	if err := msg.Message.Validate(); err != nil {
		log.Errorf("invalid message received from peer: %s", err)
		return
	}

	// This just dumps all incoming messages to the handler for logging or
	// testing purposes.
	go b.handler.HandleRawMessage(msg)

	// Check existing pending messages from other peers for the same message ID
	b.pendingMessagesLock.Lock()
	defer b.pendingMessagesLock.Unlock()

	peerMsgGroup, found := b.pendingMessages[msg.Message.ID]
	// First time we see this message
	if !found {
		log.Debugf("new message ID %s from peer %s, creating new pending group", msg.Message.ID, msg.PeerID)
		peerMsgGroup = NewPeerMessageGroup()

		// Rebroadcast to all other peers when first time seeing this message.
		go func() {
			// Send Payload, NOT the Message, as SendProtoMessage wraps it in a Message.
			if err := b.broadcastRawMessage(context.Background(), &msg.Message); err != nil {
				log.DPanic(
					"error rebroadcasting message %s: %s",
					msg.Message.ID,
					err,
				)
			}
		}()
	}

	if replaced := peerMsgGroup.Add(msg); replaced {
		// Panic in development, but error in production.
		log.DPanicw(
			"duplicate message ID received from same peer",
			"peerID", msg.PeerID,
			"message", msg,
		)
	}
	b.pendingMessages[msg.Message.ID] = peerMsgGroup

	log.Debugw(
		"added message to pending peer message group",
		"peerID", msg.PeerID,
		"messageID", msg.Message.ID,
		"newGroupLength", peerMsgGroup.Len(),
	)

	// Validate the message group -- could be done only when we have all messages
	// from all peers, or each time a message is added.
	if err := peerMsgGroup.Validate(); err != nil {
		log.Errorf("broadcast message validation failed: %s", err)

		// TODO: Reject additional messages for the same message ID? Or prune
		// them along with the expired messages
		delete(b.pendingMessages, msg.Message.ID)

		return
	}

	if peerMsgGroup.Len() == b.GetPeerCount() {
		log.Debugw(
			"pending peer message group complete",
			"messageID", msg.Message.ID,
			"groupLength", peerMsgGroup.Len(),
			"peerCount", b.GetPeerCount(),
		)

		// All peers have responded with the same message, send it to the valid
		// message channel to be handled.
		b.incomingValidatedMessages <- peerMsgGroup.GetMessageData()

		// Remove from pending messages
		delete(b.pendingMessages, msg.Message.ID)
	} else {
		log.Debugw(
			"peer message group still pending",
			"messageID", msg.Message.ID,
			"groupLength", peerMsgGroup.Len(),
			"peerCount", b.GetPeerCount(),
		)
	}
}

func (b *Broadcaster) handleIncomingValidatedMsg(msg *types.MessageData) {
	log.Infof("received validated message: %v", msg.String())

	go b.handler.HandleValidatedMessage(msg)
}

// -----------------------------------------------------------------------------
// Stream handling

// handleNewPeer opens a new stream with a newly connected peer.
func (b *Broadcaster) handleNewPeer(ctx context.Context, pid peer.ID) {
	s, err := b.host.NewStream(b.ctx, pid, ProtocolID)
	if err != nil {
		log.Errorf("failed to open new stream to peer: ", err, pid)

		return
	}

	log.Debugw("opened new stream to peer", "PeerID", pid)

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

	log.Debugf("starting stream reader for peer: %s", s.Conn().RemotePeer())

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
		peerMsg := MessageWithPeerMetadata{
			Message: msg,
			PeerID:  s.Conn().RemotePeer(),
		}

		log.Debugf("received message from peer: %s", peerMsg.PeerID)

		select {
		case b.incomingRawMessages <- &peerMsg:
		case <-b.ctx.Done():
			// Close is useless because the other side isn't reading.
			_ = s.Reset()
			return
		}
	}
}
