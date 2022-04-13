package broadcast

import (
	"github.com/libp2p/go-libp2p-core/host"
)

type Broadcast struct {
	host host.Host

	// incoming messages from other peers, unverified
	incoming chan []byte

	// outgoing messages to other peers
	outgoing chan []byte
}

func (b Broadcast) Send(msg []byte) error {
	// TODO: Create new proto Message

	// TODO: Sign the message

	b.outgoing <- msg
	return nil
}
