package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	proto "github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

var (
	ErrMsgIDEmpty         = errors.New("message ID is empty")
	ErrMsgRecipientsEmpty = errors.New("no recipient peer IDs in message")
)

// NewBroadcastMessage creates a new BroadcastMessage with the payload marshaled as Any.
func NewBroadcastMessage(
	id string,
	payload proto.Message,
	recipientsPeerIDs []peer.ID,
) (BroadcastMessage, error) {
	anyPayload, err := prototypes.MarshalAny(payload)
	if err != nil {
		return BroadcastMessage{}, err
	}

	return BroadcastMessage{
		ID:               id,
		Payload:          *anyPayload,
		RecipientPeerIDs: recipientsPeerIDs,
		Created:          time.Now().UTC(),
	}, nil
}

// Validate returns an error if the message is invalid.
func (msg *BroadcastMessage) Validate() error {
	if strings.TrimSpace(msg.ID) == "" {
		return ErrMsgIDEmpty
	}

	if len(msg.RecipientPeerIDs) == 0 {
		return ErrMsgRecipientsEmpty
	}

	for _, peerID := range msg.RecipientPeerIDs {
		if err := peerID.Validate(); err != nil {
			return fmt.Errorf("recipient peer ID is invalid: %w", err)
		}
	}

	return nil
}

// UnpackPayload unmarshals the payload message into the given proto.Message.
func (msg *BroadcastMessage) UnpackPayload(pb proto.Message) error {
	return prototypes.UnmarshalAny(&msg.Payload, pb)
}
