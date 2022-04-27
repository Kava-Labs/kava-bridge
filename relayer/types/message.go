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

// NewMessageData creates a new MessageData with the payload marshaled as Any.
func NewMessageData(id string, payload proto.Message) (MessageData, error) {
	anyPayload, err := prototypes.MarshalAny(payload)
	if err != nil {
		return MessageData{}, err
	}

	return MessageData{
		ID:      id,
		Payload: anyPayload,
	}, nil
}

// Validate returns an error if the message is invalid.
func (msg *MessageData) Validate() error {
	if strings.TrimSpace(msg.ID) == "" {
		return ErrMsgIDEmpty
	}

	return nil
}

// UnpackPayload unmarshals the payload message into the given proto.Message.
func (msg *MessageData) UnpackPayload(pb proto.Message) error {
	return prototypes.UnmarshalAny(msg.Payload, pb)
}
