package pending_store

import (
	"fmt"

	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// PeerMessageGroup is a group of the same message from different peers to be
// validated.
type PeerMessageGroup struct {
	// BroadcastedMessage will only be received by broadcaster peer.
	BroadcastedMessage types.BroadcastMessage
	// BroadcastedMessageReceived will be false before the message is received. This may
	// occur if a hash is received from a re-broadcasting node before the
	// broadcasted message.
	BroadcastedMessageReceived bool
	// Hash of messages from other peers.
	PeerMessageHashes map[peer.ID]types.BroadcastMessageHash
}

// NewPeerMessageGroup returns a new PeerMessageGroup
func NewPeerMessageGroup() *PeerMessageGroup {
	return &PeerMessageGroup{
		BroadcastedMessage:         types.BroadcastMessage{},
		BroadcastedMessageReceived: false,
		PeerMessageHashes:          make(map[peer.ID]types.BroadcastMessageHash),
	}
}

// AddMessage adds a broadcasted message to the group, returning an error if the
// message already exists.
func (g *PeerMessageGroup) AddMessage(msg types.BroadcastMessage) error {
	if g.BroadcastedMessageReceived {
		return fmt.Errorf("message already received")
	}

	g.BroadcastedMessage = msg
	g.BroadcastedMessageReceived = true

	return nil
}

// AddHash adds a message hash to the group, returning an error if a hash
// already exists for a peer. This does **not** check if the message is
// validated as a hash may be added before the broadcasted message is received.
func (g *PeerMessageGroup) AddHash(peerID peer.ID, hash types.BroadcastMessageHash) error {
	if _, found := g.PeerMessageHashes[peerID]; found {
		return fmt.Errorf("peer hash %s already exists", peerID)
	}

	g.PeerMessageHashes[peerID] = hash

	return nil
}

// Completed returns true if the number of received messages matches the number
// of recipients.
func (g *PeerMessageGroup) Completed() bool {
	return g.BroadcastedMessageReceived &&
		len(g.PeerMessageHashes) == len(g.BroadcastedMessage.RecipientPeerIDs)
}

// GetMessageData returns the underlying MessageData for the group. This should
// be called *after* Validate() has been called and confirmed to have no errors.
// This may return false if the group was created but did not receive any
// messages (ie. when broadcasting)
func (g *PeerMessageGroup) GetMessageData() (types.BroadcastMessage, bool) {
	return g.BroadcastedMessage, g.BroadcastedMessageReceived
}

// Validate returns nil if all message hashes are valid OR if it is still
// waiting on additional information e.g. the broadcasted message. This only
// returns an error if the entire group should be invalidated and discarded.
func (g *PeerMessageGroup) Validate() error {
	if !g.BroadcastedMessageReceived {
		return nil
	}

	broadcastMessageHash, err := g.BroadcastedMessage.Hash()
	if err != nil {
		return err
	}

	for peerID, hash := range g.PeerMessageHashes {
		if !broadcastMessageHash.Equal(hash) {
			return fmt.Errorf(
				"group contains invalid hash for peer %s, got %v, expected %v",
				peerID, hash, broadcastMessageHash,
			)
		}
	}

	return nil
}
