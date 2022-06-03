package pending_store

import (
	"fmt"

	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// PeerMessageGroup is a group of the same message from different peers to be
// validated.
type PeerMessageGroup struct {
	IsSender bool
	Messages map[peer.ID]*MessageWithPeerMetadata
}

// NewPeerMessageGroup returns a new PeerMessageGroup
func NewPeerMessageGroup(isSender bool) *PeerMessageGroup {
	return &PeerMessageGroup{
		IsSender: isSender,
		Messages: make(map[peer.ID]*MessageWithPeerMetadata),
	}
}

// Add adds a message to the group, returning true if it replaced a message with
// the same peer ID.
func (g *PeerMessageGroup) Add(msg *MessageWithPeerMetadata) error {
	if _, found := g.Messages[msg.PeerID]; found {
		return fmt.Errorf("message from peer %s already exists", msg.PeerID)
	}

	g.Messages[msg.PeerID] = msg

	return nil
}

// Completed returns true if the number of received messages matches the number
// of recipients.
func (g *PeerMessageGroup) Completed(hostID peer.ID, recipients []peer.ID) bool {
	for _, recipient := range recipients {
		// Ignore current host peer as it may be in recipients list but won't be
		// contained in messages.
		if recipient == hostID {
			continue
		}

		// All other recipients need to be contained in messages
		if _, found := g.Messages[recipient]; !found {
			return false
		}
	}

	return true
}

// GetMessageData returns the underlying MessageData for the group. This should
// be called *after* Validate() has been called and confirmed to have no errors.
// This may return false if the group was created but did not receive any
// messages (ie. when broadcasting)
func (g *PeerMessageGroup) GetMessageData() (types.BroadcastMessage, bool) {
	for _, msg := range g.Messages {
		return msg.BroadcastMessage, true
	}

	return types.BroadcastMessage{}, false
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
			messageID = msg.BroadcastMessage.ID
		}

		if msg.BroadcastMessage.ID != messageID {
			return fmt.Errorf(
				"message ID from peer %s mismatch: %q != %q",
				msg.PeerID, msg.BroadcastMessage.ID, messageID,
			)
		}

		if message == nil {
			message = msg
		}

		// TODO: Ensure all messages are the same and signed?
		if !msg.BroadcastMessage.Payload.Equal(message.BroadcastMessage.Payload) {
			return fmt.Errorf("message payloads do not match from peer %s", msg.PeerID)
		}
	}

	return nil
}
