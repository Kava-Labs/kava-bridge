package broadcast

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/kava-labs/kava-bridge/relayer/broadcast/pending_store"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/kava-labs/kava-bridge/relayer/stream"
)

var log = logging.Logger("broadcast")
var tracer = otel.Tracer("Broadcaster")

const (
	ProtocolID  = "/kava-relayer/broadcast/1.0.0"
	ServiceName = "kava-relayer.broadcast"
)

// P2PBroadcaster is a reliable broadcaster to ensure that all connected peers
// receive the same message.
type P2PBroadcaster struct {
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
	incomingRawMessages chan *pending_store.MessageWithPeerMetadata

	// Messages that all peers have confirmed to have received. Does not contain
	// any peer specific data as it does not originate from any specific peer.
	incomingValidatedMessages chan types.BroadcastMessage

	// Messages that have been sent/received but not validated by other peers yet.
	pendingMessagesStore *pending_store.PendingMessagesStore

	// Messages to send to all other peers
	outgoing chan *types.BroadcastMessage

	// External event handler
	handler BroadcastHandler

	// Message hook
	broadcasterHook BroadcasterHook

	ctx context.Context

	log *zap.SugaredLogger
}

var _ Broadcaster = &P2PBroadcaster{}

// NewBroadcaster returns a new Broadcaster
func NewBroadcaster(
	ctx context.Context,
	host host.Host,
	options ...BroadcasterOption,
) (*P2PBroadcaster, error) {
	b := &P2PBroadcaster{
		host:                      host,
		inboundStreamsLock:        sync.Mutex{},
		inboundStreams:            make(map[peer.ID]network.Stream),
		outboundStreamsLock:       sync.Mutex{},
		outboundStreams:           make(map[peer.ID]network.Stream),
		peersLock:                 sync.RWMutex{},
		peers:                     make(map[peer.ID]struct{}),
		newPeers:                  make(chan peer.ID),
		incomingRawMessages:       make(chan *pending_store.MessageWithPeerMetadata, 1),
		incomingValidatedMessages: make(chan types.BroadcastMessage, 1),
		pendingMessagesStore:      pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL),
		outgoing:                  make(chan *types.BroadcastMessage),
		handler:                   &NoOpBroadcastHandler{},
		broadcasterHook:           &noOpBroadcasterHook{},
		ctx:                       ctx,
		log:                       log.Named(host.ID().ShortString()),
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
func (b *P2PBroadcaster) GetPeerCount() int {
	b.peersLock.RLock()
	defer b.peersLock.RUnlock()

	return len(b.peers)
}

// processLoop handles incoming channel inputs.
func (b *P2PBroadcaster) processLoop(ctx context.Context) {
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
		// Validated messages have priority to not deadlock when raw message
		// pushes more than 1 incoming validated message before handled.
		select {
		case newValidatedMsg := <-b.incomingValidatedMessages:
			b.handleIncomingValidatedMsg(newValidatedMsg)
		default:
		}

		select {
		case newPeerID := <-b.newPeers:
			b.handleNewPeer(ctx, newPeerID)
		case newRawMsg := <-b.incomingRawMessages:
			b.handleIncomingRawMsg(newRawMsg)
		case <-ctx.Done():
			b.log.Info("broadcast handler loop shutting down")
			return
		}
	}
}

// -----------------------------------------------------------------------------
// Peer messages handling

// BroadcastMessage marshals the proto.Message as Any, wraps it in MessageData,
// and it to all connected peers.
func (b *P2PBroadcaster) BroadcastMessage(
	ctx context.Context,
	pb types.PeerMessage,
	recipients []peer.ID,
	TTLSeconds uint64,
) error {
	// Wrap the proto message in the MessageData type.
	msg, err := types.NewBroadcastMessage(ctx, pb, b.host.ID(), recipients, TTLSeconds)
	if err != nil {
		return err
	}

	ctx, span := tracer.Start(
		ctx,
		"broadcast message", trace.WithAttributes(
			attribute.String("typeUrl", msg.Payload.TypeUrl),
			attribute.String("messageID", msg.ID),
			attribute.String("fromPeerID", b.host.ID().ShortString()),
		),
	)
	defer span.End()

	if err := msg.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// Add the message to the pending messages map to keep track of responses
	// and to prevent re-broadcasting.
	// Does not block receiving messages while broadcasting
	created := b.pendingMessagesStore.TryNewGroup(msg.ID)
	if !created {
		return fmt.Errorf("cannot broadcast message that is is already pending")
	}

	return b.broadcastRawMessage(ctx, &msg, recipients)
}

// broadcastRawMessage sends a proto message to all connected peers without any
// marshalling or wrapping.
func (b *P2PBroadcaster) broadcastRawMessage(
	ctx context.Context,
	pb proto.Message,
	recipients []peer.ID,
) error {
	b.log.Debugw("sending raw proto message", "recipients", recipients)

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
			b.log.Debugf("no outbound stream for peer %v, opening new one", peerID)

			// Try to open a new stream to this peer.
			// s, err := b.host.NewStream(b.ctx, peerID, ProtocolID)
			// if err != nil {
			// 	return fmt.Errorf("failed to open new stream to peer %s: %w", peerID, err)
			// }
			//
			// b.outboundStreams[peerID] = s
			// ch = s
			panic("no outbound stream")
		}

		// Avoid capturing loop variable
		func(peerID peer.ID, ch network.Stream) {
			group.Go(func() error {
				b.broadcasterHook.BeforeBroadcastRawMessage(b, peerID, &pb)

				// Check if still connected to peer
				// TODO: Reconnect if not connected to peer.
				if b.host.Network().Connectedness(peerID) != network.Connected {
					return fmt.Errorf("peer %v is not connected", peerID.ShortString())
				}

				b.log.Debugf("sending message to peer: %s", peerID.ShortString())

				// NewUint32DelimitedWriter has an internal buffer, bufio.NewWriter()
				// should not be necessary.
				err := stream.NewProtoMessageWriter(ch).WriteMsg(pb)
				if err != nil {
					return fmt.Errorf("failed to write proto message to peer: %w", err)
				}

				b.log.Debugf("wrote message to peer stream: %s", peerID.ShortString())

				return nil
			})
		}(peerID, ch)
	}

	return group.Wait()
}

// handleIncomingRawMsg handles all raw messages from other peers. This is
// before messages are verified to be received from all peers.
func (b *P2PBroadcaster) handleIncomingRawMsg(msg *pending_store.MessageWithPeerMetadata) {
	b.log.Debugf("received raw message from %v", msg.PeerID.ShortString())

	if err := msg.BroadcastMessage.Validate(); err != nil {
		b.log.Warnf("invalid message received from peer: %s", err)
		return
	}

	// This just dumps all incoming messages to the handler for logging or
	// testing purposes.
	go b.handler.RawMessage(*msg)

	// Create recipients list including the sender.
	recipientsAndSender := msg.BroadcastMessage.GetAllPeersIDs()

	// Create new group if it doesn't exist already.
	created := b.pendingMessagesStore.TryNewGroup(msg.BroadcastMessage.ID)
	// First time we see this message, re-broadcast
	if created {
		b.log.Debugf(
			"new message ID %s from peer %s, creating new pending group",
			msg.BroadcastMessage.ID,
			msg.PeerID.ShortString(),
		)

		// Update IsBroadcaster for re-broadcasts
		msg.BroadcastMessage.IsBroadcaster = false

		// Rebroadcast to all other peers when first time seeing this message.
		// Run in a goroutine to avoid blocking the incoming message handler.
		go func() {
			// Send Payload, NOT the BroadcastMessage, as SendProtoMessage wraps it in a Message.
			if err := b.broadcastRawMessage(
				msg.Context,
				&msg.BroadcastMessage,
				recipientsAndSender,
			); err != nil {
				b.log.DPanicf(
					"error rebroadcasting message %s: %s",
					msg.BroadcastMessage.ID,
					err,
				)
			}
		}()
	}

	if err := b.pendingMessagesStore.AddMessage(*msg); err != nil {
		b.log.Errorf("error adding message %s to pending messages store: %s", msg.BroadcastMessage.ID, err)

		// Remove from pending messages -- if there is 1 invalid message then
		// the message should be cancelled
		if err := b.pendingMessagesStore.DeleteGroup(msg.BroadcastMessage.ID); err != nil {
			b.log.Warnw(
				"failed to remove invalid message group from pending messages store",
				"msgId", msg.BroadcastMessage.ID,
				"err", err,
			)
		}

		return
	}

	if msgData, completed := b.pendingMessagesStore.GroupIsCompleted(
		msg.BroadcastMessage.ID,
		b.host.ID(),
		recipientsAndSender,
	); completed {
		// All peers have responded with the same message, send it to the valid
		// message channel to be handled.
		b.incomingValidatedMessages <- msgData

		// Remove from pending messages
		if err := b.pendingMessagesStore.DeleteGroup(msg.BroadcastMessage.ID); err != nil {
			b.log.Warnw(
				"failed to remove completed message group from pending messages store",
				"msgId", msg.BroadcastMessage.ID,
				"err", err,
			)
		}
	}
}

func (b *P2PBroadcaster) handleIncomingValidatedMsg(msg types.BroadcastMessage) {
	if msg.From == b.host.ID() {
		b.log.Debugf("ignoring message from self")
		return
	}

	b.log.Infow(
		"received validated message",
		"ID", msg.ID,
		"From", msg.From,
		"To", msg.RecipientPeerIDs,
		"typeUrl", msg.Payload.TypeUrl,
	)

	go b.handler.ValidatedMessage(msg)
}

// -----------------------------------------------------------------------------
// Stream handling

// handleNewPeer opens a new stream with a newly connected peer.
func (b *P2PBroadcaster) handleNewPeer(ctx context.Context, pid peer.ID) {
	if b.host.Network().Connectedness(pid) != network.Connected {
		return
	}

	b.outboundStreamsLock.Lock()
	if _, found := b.outboundStreams[pid]; found {
		b.outboundStreamsLock.Unlock()
		b.log.Debugf("peer %s is already connected, skip opening new stream", pid.ShortString())

		return
	}

	s, err := b.host.NewStream(b.ctx, pid, ProtocolID)
	if err != nil {
		b.log.Errorf("failed to open new stream to peer %s: %v", pid.ShortString(), err)

		b.outboundStreamsLock.Unlock()
		return
	}

	b.log.Debugw("opened new stream to peer", "PeerID", pid.ShortString())

	// Locked earlier
	b.outboundStreams[pid] = s
	b.outboundStreamsLock.Unlock()

	b.peersLock.Lock()
	b.peers[pid] = struct{}{}
	b.peersLock.Unlock()
}

// handleNewStream handles a new incoming stream, initiated when a peer is connected.
func (b *P2PBroadcaster) handleNewStream(s network.Stream) {
	peer := s.Conn().RemotePeer()

	b.log.Debugf(
		"incoming stream from peer: %s (%v)",
		peer.ShortString(),
		s.Protocol(),
	)

	// Ensure only 1 incoming stream from the peer
	b.inboundStreamsLock.Lock()
	other, dup := b.inboundStreams[peer]
	if dup {
		b.log.Debugf("duplicate inbound stream from %s; resetting other stream", peer.ShortString())

		if err := other.Reset(); err != nil {
			b.log.Warnf("error resetting other stream: %s", err)
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

	b.log.Debugf(
		"starting stream reader from peer %s for protocol %s",
		s.Conn().RemotePeer().ShortString(),
		s.Protocol(),
	)

	// Iterate over all messages, unmarshalling all as types.MessageData
	r := stream.NewProtoMessageReader(s)
	for {
		var msg types.BroadcastMessage
		if err := r.ReadMsg(&msg); err != nil {
			// Error when closing stream
			b.log.Warnf("error reading stream message from peer %s: %s", s.Conn().RemotePeer().ShortString(), err)
			_ = s.Reset()

			return
		}

		ctx := context.Background()

		carrier := msg.TraceContext.GetCarrier()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
		ctx, span := tracer.Start(ctx, "receive message")

		// Attach additional peer metadata to the message
		peerMsg := pending_store.MessageWithPeerMetadata{
			BroadcastMessage: msg,
			PeerID:           s.Conn().RemotePeer(),
			Context:          ctx,
		}

		if msg.IsBroadcaster && msg.From != s.Conn().RemotePeer() {
			b.log.Warnf(
				"broadcast message from does not match sender: got %s, expected %s",
				msg.From,
				s.Conn().RemotePeer(),
			)

			continue
		}

		// TODO: Redundant unpack, when payload is used it will be unpacked again
		broadcastMsg, err := peerMsg.BroadcastMessage.UnpackPayload()
		if err != nil {
			b.log.Warnf(
				"error unpacking payload for message %s from peer %s: %s",
				msg.ID,
				s.Conn().RemotePeer(),
				err,
			)

			return
		}

		if err := broadcastMsg.ValidateBasic(); err != nil {
			b.log.Warnf("invalid message from peer %s: %s", s.Conn().RemotePeer(), err)

			continue
		}

		span.AddEvent("Message received", trace.WithAttributes(
			attribute.String("Sender", peer.ShortString()),
			attribute.String("Receiver", b.host.ID().ShortString()),
			attribute.Bool("IsRebroadcast", msg.IsBroadcaster),
			attribute.String("Payload.TypeUrl", peerMsg.BroadcastMessage.Payload.TypeUrl),
		))

		// TODO: Memory leak if something fails prior to this
		span.End()

		select {
		case b.incomingRawMessages <- &peerMsg:
		case <-b.ctx.Done():
			// Close is useless because the other side isn't reading.
			_ = s.Reset()
			return
		}
	}
}
