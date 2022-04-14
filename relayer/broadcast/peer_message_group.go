package broadcast

import (
	"fmt"

	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// PeerMessageGroup is a group of the same message from different peers to be
// validated.
type PeerMessageGroup struct {
	Messages map[peer.ID]*MessageWithPeerMetadata
}

// NewPeerMessageGroup returns a new PeerMessageGroup
func NewPeerMessageGroup() *PeerMessageGroup {
	return &PeerMessageGroup{
		Messages: make(map[peer.ID]*MessageWithPeerMetadata),
	}
}

// Add adds a message to the group, returning true if it replaced a message with
// the same peer ID.
func (g *PeerMessageGroup) Add(msg *MessageWithPeerMetadata) bool {
	_, found := g.Messages[msg.PeerID]
	g.Messages[msg.PeerID] = msg

	return found
}

// GetMessageData returns the underlying MessageData for the group. This should
// be called *after* Validate() has been called and confirmed to have no errors.
func (g *PeerMessageGroup) GetMessageData() *types.MessageData {
	for _, msg := range g.Messages {
		return &msg.Message
	}

	return nil
}

// Len returns the number of messages in the group.
func (g *PeerMessageGroup) Len() int {
	return len(g.Messages)
}

// Validate returns true if all messages in the group are the same.
func (g *PeerMessageGroup) Validate() error {
	if len(g.Messages) == 0 {
		return nil
	}

	var messageID string
	var message *MessageWithPeerMetadata

	for _, msg := range g.Messages {
		// All messages checked against the first one in slice
		// TODO: Return the real message that is different from the others. e.g
		// If the first one is the wrong one, this reports the second one as
		// wrong.

		// Set messageID on first iteration
		if messageID == "" {
			messageID = msg.Message.ID
		}

		if msg.Message.ID != messageID {
			return fmt.Errorf(
				"message ID from peer %s mismatch: %q != %q",
				msg.PeerID, msg.Message.ID, messageID,
			)
		}

		if message == nil {
			message = msg
		}

		// TODO: Ensure all messages are the same and signed?
		if !msg.Message.Payload.Equal(message.Message.Payload) {
			return fmt.Errorf("message payloads do not match from peer %s", msg.PeerID)
		}
	}

	return nil
}
