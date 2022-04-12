package types

import (
	proto "github.com/gogo/protobuf/proto"
	prototypes "github.com/gogo/protobuf/types"
)

// NewMessageData creates a new MessageData with the payload marshaled as Any.
func NewMessageData(payload proto.Message) (MessageData, error) {
	anyPayload, err := prototypes.MarshalAny(payload)
	if err != nil {
		return MessageData{}, err
	}

	return MessageData{
		Payload: anyPayload,
	}, nil
}

// UnpackPayload unmarshals the payload message into the given proto.Message.
func (msg *MessageData) UnpackPayload(pb proto.Message) error {
	return prototypes.UnmarshalAny(msg.Payload, pb)
}
