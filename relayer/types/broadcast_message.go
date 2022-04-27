package types

import (
	"errors"
	"strings"

	proto "github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
)

var (
	ErrMsgIDEmpty = errors.New("message ID is empty")
)

// NewBroadcastMessage creates a new BroadcastMessage with the payload marshaled as Any.
func NewBroadcastMessage(id string, payload proto.Message) (BroadcastMessage, error) {
	anyPayload, err := prototypes.MarshalAny(payload)
	if err != nil {
		return BroadcastMessage{}, err
	}

	return BroadcastMessage{
		ID:      id,
		Payload: anyPayload,
	}, nil
}

// Validate returns an error if the message is invalid.
func (msg *BroadcastMessage) Validate() error {
	if strings.TrimSpace(msg.ID) == "" {
		return ErrMsgIDEmpty
	}

	return nil
}

// UnpackPayload unmarshals the payload message into the given proto.Message.
func (msg *BroadcastMessage) UnpackPayload(pb proto.Message) error {
	return prototypes.UnmarshalAny(msg.Payload, pb)
}
