package types

import (
	"github.com/gogo/protobuf/proto"
)

// PeerMessage defines an interface that broadcast messages must implement.
type PeerMessage interface {
	proto.Message

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error
}
