package broadcast

import (
	"fmt"
)

// PeerMessageGroup is a group of the same message from different peers to be
// validated.
type PeerMessageGroup struct {
	Messages []*MessageWithPeerMetadata
}

func (g *PeerMessageGroup) Add(msg *MessageWithPeerMetadata) {
	g.Messages = append(g.Messages, msg)
}

func (g *PeerMessageGroup) Validate() error {
	for _, msg := range g.Messages {
		// TODO: Ensure all messages are the same and signed?
		if !msg.Message.Payload.Equal(g.Messages[0].Message.Payload) {
			return fmt.Errorf("message payloads do not match from peer %s", msg.PeerID)
		}
	}

	return nil
}
