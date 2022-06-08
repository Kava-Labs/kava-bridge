package types

import (
	"context"
	"fmt"
	"strings"
	"time"

	prototypes "github.com/gogo/protobuf/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multibase"
	"github.com/pkg/errors"
)

const (
	MinimumTTLSeconds = 1
)

// NewBroadcastMessage creates a new BroadcastMessage with the payload marshaled as Any.
func NewBroadcastMessage(
	ctx context.Context,
	payload PeerMessage,
	hostID peer.ID,
	recipientsPeerIDs []peer.ID,
	TTLSeconds uint64,
) (BroadcastMessage, error) {
	messageID, err := NewBroadcastMessageID()
	if err != nil {
		return BroadcastMessage{}, err
	}

	anyPayload, err := prototypes.MarshalAny(payload)
	if err != nil {
		return BroadcastMessage{}, err
	}

	allPeerIDs := append(recipientsPeerIDs, hostID)
	allPeerIDs = dedupPeerIDs(allPeerIDs)

	// Add trace context to the message.
	traceCtx := NewTraceContext()
	traceCtx.Inject(ctx)

	return BroadcastMessage{
		ID:               messageID,
		From:             hostID,
		IsBroadcaster:    true,
		Payload:          *anyPayload,
		RecipientPeerIDs: allPeerIDs,
		Created:          time.Now().UTC(),
		TTLSeconds:       TTLSeconds,
		TraceContext:     traceCtx,
	}, nil
}

// Validate returns an error if the message is invalid.
func (msg *BroadcastMessage) Validate() error {
	if strings.TrimSpace(msg.ID) == "" {
		return ErrMsgIDEmpty
	}

	_, _, err := multibase.Decode(msg.ID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	if err := msg.From.Validate(); err != nil {
		return fmt.Errorf("invalid from peer.ID: %w", err)
	}

	if len(msg.RecipientPeerIDs) < 1 {
		return ErrMsgInsufficientRecipients
	}

	if duplicatePeerID, found := containsDuplicatePeerID(msg.RecipientPeerIDs); found {
		return fmt.Errorf("duplicate recipient peer ID in message: %s", duplicatePeerID)
	}

	for _, peerID := range msg.RecipientPeerIDs {
		if err := peerID.Validate(); err != nil {
			return fmt.Errorf("recipient peer ID is invalid: %w", err)
		}
	}

	if msg.TTLSeconds < MinimumTTLSeconds {
		return errors.Wrapf(ErrMsgTTLTooShort, "%d < %d seconds", msg.TTLSeconds, MinimumTTLSeconds)
	}

	if msg.Expired() {
		return errors.Wrapf(
			ErrMsgExpired,
			"%v + %v seconds < now (%v)",
			msg.Created, msg.TTLSeconds, time.Now().UTC(),
		)
	}

	return nil
}

// UnpackPayload unpacks the broadcast message payload into a PeerMessage.
func (msg *BroadcastMessage) UnpackPayload() (PeerMessage, error) {
	var payloadDyn prototypes.DynamicAny
	err := prototypes.UnmarshalAny(&msg.Payload, &payloadDyn)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal payload any: %w", err)
	}

	peerMsg, ok := payloadDyn.Message.(PeerMessage)
	if !ok {
		return nil, fmt.Errorf(
			"payload does not implement PeerMessage interface, got invalid payload type: %T",
			payloadDyn.Message,
		)
	}

	return peerMsg, nil
}

// Expired returns true if the TTL is exceeded since created time.
func (msg *BroadcastMessage) Expired() bool {
	// TTLSeconds converted to float, not converting duration to uint64 as it
	// will underflow.
	return time.Since(msg.Created).Seconds() > float64(msg.TTLSeconds)
}

func dedupPeerIDs(peerIDs []peer.ID) []peer.ID {
	seen := make(map[peer.ID]struct{})
	var deduped []peer.ID

	for _, peerID := range peerIDs {
		if _, ok := seen[peerID]; !ok {
			seen[peerID] = struct{}{}
			deduped = append(deduped, peerID)
		}
	}

	return deduped
}

func containsDuplicatePeerID(peerIDs []peer.ID) (peer.ID, bool) {
	seen := make(map[peer.ID]struct{})

	for _, peerID := range peerIDs {
		if _, found := seen[peerID]; found {
			return peerID, true
		}

		seen[peerID] = struct{}{}
	}

	return "", false
}
