package types

import (
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/peer"
)

// PeerMessage defines an interface that broadcast messages must implement.
type PeerMessage interface {
	proto.Message

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// GetSenderPeerID returns the peer ID of the peer that must send the message.
	GetSenderPeerID() peer.ID
}
