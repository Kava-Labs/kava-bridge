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
	incomingValidatedMessages chan *types.BroadcastMessage

	// Messages that have been received but have not been validated to be same
	// from all peers.
	pendingMessagesLock sync.Mutex
	// TODO: Remove messages after expired. Add TTL to BroadcastMessage.
	pendingMessages map[string]*PeerMessageGroup

	// Messages to send to all other peers
	outgoing chan *types.BroadcastMessage

	// External event handler
	handler BroadcastHandler

	// Message hook
	broadcasterHook broadcasterHook

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
		incomingValidatedMessages: make(chan *types.BroadcastMessage, 1),
		pendingMessagesLock:       sync.Mutex{},
		pendingMessages:           make(map[string]*PeerMessageGroup),
		outgoing:                  make(chan *types.BroadcastMessage),
		handler:                   &NoOpBroadcastHandler{},
		broadcasterHook:           &noOpBroadcasterHook{},
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
	pb proto.Message,
	recipients []peer.ID,
) error {
	// Wrap the proto message in the MessageData type.
	msg, err := types.NewBroadcastMessage(pb, b.host.ID(), recipients)
	if err != nil {
		return err
	}

	if err := msg.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
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

	return b.broadcastRawMessage(ctx, &msg, recipients)
}

// broadcastRawMessage sends a proto message to all connected peers without any
// marshalling or wrapping.
func (b *Broadcaster) broadcastRawMessage(
	ctx context.Context,
	pb proto.Message,
	recipients []peer.ID,
) error {
	log.Debugf("broadcast sending raw proto message: %s", pb.String())

	b.outboundStreamsLock.Lock()
	defer b.outboundStreamsLock.Unlock()

	// Run writes to peers in parallel
	group, _ := errgroup.WithContext(ctx)

	// Only send messages in the recipient list.
	for _, peerID := range recipients {
		// Ignore self, will be contained in recipient list if not original broadcaster
		if peerID == b.host.ID() {
			continue
		}

		ch, ok := b.outboundStreams[peerID]
		if !ok {
			// TODO: Try to open a new stream to this peer.
			return fmt.Errorf("no outbound stream for peer: %v", peerID)
		}

		// Avoid capturing loop variable
		func(peerID peer.ID, ch network.Stream) {
			b.broadcasterHook.BeforeBroadcastRawMessage(b, peerID, &pb)

			group.Go(func() error {
				// Check if still connected to peer
				// TODO: Reconnect if not connected to peer.
				if b.host.Network().Connectedness(peerID) != network.Connected {
					return fmt.Errorf("peer %v is not connected", peerID)
				}

				log.Debugf("sending message to peer: %s", peerID)

				// NewUint32DelimitedWriter has an internal buffer, bufio.NewWriter()
				// should not be necessary.
				err := stream.NewProtoMessageWriter(ch).WriteMsg(pb)
				if err != nil {
					return fmt.Errorf("failed to write proto message to peer: %w", err)
				}
				return nil
			})
		}(peerID, ch)
	}

	return group.Wait()
}

// handleIncomingRawMsg handles all raw messages from other peers. This is
// before messages are verified to be received from all peers.
func (b *Broadcaster) handleIncomingRawMsg(msg *MessageWithPeerMetadata) {
	if err := msg.BroadcastMessage.Validate(); err != nil {
		log.Errorf("invalid message received from peer: %s", err)
		return
	}

	// This just dumps all incoming messages to the handler for logging or
	// testing purposes.
	go b.handler.RawMessage(*msg)

	// Check existing pending messages from other peers for the same message ID
	b.pendingMessagesLock.Lock()
	defer b.pendingMessagesLock.Unlock()

	peerMsgGroup, found := b.pendingMessages[msg.BroadcastMessage.ID]
	// First time we see this message
	if !found {
		log.Debugf("new message ID %s from peer %s, creating new pending group", msg.BroadcastMessage.ID, msg.PeerID)
		peerMsgGroup = NewPeerMessageGroup()

		// Rebroadcast to all other peers when first time seeing this message.
		go func() {
			// Send Payload, NOT the BroadcastMessage, as SendProtoMessage wraps it in a Message.
			if err := b.broadcastRawMessage(
				context.Background(),
				&msg.BroadcastMessage,
				msg.BroadcastMessage.RecipientPeerIDs,
			); err != nil {
				log.DPanic(
					"error rebroadcasting message %s: %s",
					msg.BroadcastMessage.ID,
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
	b.pendingMessages[msg.BroadcastMessage.ID] = peerMsgGroup

	log.Debugw(
		"added message to pending peer message group",
		"peerID", msg.PeerID,
		"messageID", msg.BroadcastMessage.ID,
		"newGroupLength", peerMsgGroup.Len(),
	)

	// Validate the message group -- could be done only when we have all messages
	// from all peers, or each time a message is added.
	if err := peerMsgGroup.Validate(); err != nil {
		log.Errorf("broadcast message validation failed: %s", err)
		go b.handler.MismatchMessage(*msg)

		// TODO: Reject additional messages for the same message ID? Or prune
		// them along with the expired messages
		delete(b.pendingMessages, msg.BroadcastMessage.ID)

		return
	}

	if peerMsgGroup.Completed(b.host.ID(), msg.BroadcastMessage.RecipientPeerIDs) {
		log.Debugw(
			"pending peer message group complete",
			"messageID", msg.BroadcastMessage.ID,
			"groupLength", peerMsgGroup.Len(),
			"recipientsLength", len(msg.BroadcastMessage.RecipientPeerIDs),
		)

		// All peers have responded with the same message, send it to the valid
		// message channel to be handled.
		b.incomingValidatedMessages <- peerMsgGroup.GetMessageData()

		// Remove from pending messages
		delete(b.pendingMessages, msg.BroadcastMessage.ID)
	} else {
		log.Debugw(
			"peer message group still pending",
			"messageID", msg.BroadcastMessage.ID,
			"groupLength", peerMsgGroup.Len(),
			"recipientsLength", len(msg.BroadcastMessage.RecipientPeerIDs),
		)
	}
}

func (b *Broadcaster) handleIncomingValidatedMsg(msg *types.BroadcastMessage) {
	log.Infof("received validated message: %v", msg.String())

	go b.handler.ValidatedMessage(*msg)
}

// -----------------------------------------------------------------------------
// Stream handling

// handleNewPeer opens a new stream with a newly connected peer.
func (b *Broadcaster) handleNewPeer(ctx context.Context, pid peer.ID) {
	s, err := b.host.NewStream(b.ctx, pid, ProtocolID)
	if err != nil {
		log.Errorf("failed to open new stream to peer %s: %v", pid, err)

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
		var msg types.BroadcastMessage
		if err := r.ReadMsg(&msg); err != nil {
			// Error when closing stream
			log.Warnf("error reading stream message from peer %s: %s", s.Conn().RemotePeer(), err)
			_ = s.Reset()

			return
		}

		// Attach additional peer metadata to the message
		peerMsg := MessageWithPeerMetadata{
			BroadcastMessage: msg,
			PeerID:           s.Conn().RemotePeer(),
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
