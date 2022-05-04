package pending_store

import (
	"fmt"
	"sync"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

var log = logging.Logger("broadcast/pending_store")

const CLEAR_EXPIRED_INTERVAL = time.Second * 30

type PendingMessagesStore struct {
	pendingMessagesLock sync.RWMutex
	pendingMessages     map[string]*PeerMessageGroup
}

// NewPendingMessagesStore returns a new PendingMessagesStore and starts a
// background goroutine to remove expired message groups.
func NewPendingMessagesStore() *PendingMessagesStore {
	store := &PendingMessagesStore{
		pendingMessagesLock: sync.RWMutex{},
		pendingMessages:     make(map[string]*PeerMessageGroup),
	}

	go func() {
		// Leaky Ticker, will not be collected by GC but should live for
		// entirety of process.
		for range time.Tick(CLEAR_EXPIRED_INTERVAL) {
			store.pendingMessagesLock.Lock()

			for key, msg := range store.pendingMessages {
				msgData, found := msg.GetMessageData()
				// If broadcasting and did not receive any yet
				if !found {
					continue
				}

				if msgData.Expired() {
					delete(store.pendingMessages, key)
				}
			}

			store.pendingMessagesLock.Unlock()
		}
	}()

	return store
}

// NewGroup creates a new PeerMessageGroup for the given message ID.
func (pm *PendingMessagesStore) ContainsGroup(msgID string) bool {
	pm.pendingMessagesLock.RLock()
	defer pm.pendingMessagesLock.RUnlock()

	_, found := pm.pendingMessages[msgID]
	return found
}

// TryNewGroup creates a new PeerMessageGroup for the given message ID and returns
// true if it was created. Returns false if the group already exists.
func (pm *PendingMessagesStore) TryNewGroup(msgID string) bool {
	pm.pendingMessagesLock.Lock()
	defer pm.pendingMessagesLock.Unlock()

	if _, found := pm.pendingMessages[msgID]; found {
		log.Debug("message group already exists", "msgID", msgID)
		return false
	}

	pm.pendingMessages[msgID] = NewPeerMessageGroup()

	return true
}

// DeleteGroup deletes the group with the given message ID.
func (pm *PendingMessagesStore) DeleteGroup(msgID string) error {
	pm.pendingMessagesLock.Lock()
	defer pm.pendingMessagesLock.Unlock()

	if _, found := pm.pendingMessages[msgID]; !found {
		return ErrGroupNotFound
	}

	delete(pm.pendingMessages, msgID)

	return nil
}

// AddMessage adds a message to it's corresponding pending message group,
// returning an error if it is invalid.
func (pm *PendingMessagesStore) AddMessage(msg MessageWithPeerMetadata) error {
	pm.pendingMessagesLock.Lock()
	defer pm.pendingMessagesLock.Unlock()

	peerMsgGroup, found := pm.pendingMessages[msg.BroadcastMessage.ID]
	if !found {
		return ErrGroupNotFound
	}

	if err := peerMsgGroup.Add(&msg); err != nil {
		return err
	}

	// TODO: Optimize validate
	if err := peerMsgGroup.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	log.Debugw("added message to pending message group", "msgID", msg.BroadcastMessage.ID)
	pm.pendingMessages[msg.BroadcastMessage.ID] = peerMsgGroup

	return nil
}

// GroupIsCompleted returns (the broadcast message and true) if the number of
// received messages matches the number of recipients.
func (pm *PendingMessagesStore) GroupIsCompleted(
	msgID string,
	hostID peer.ID,
	recipients []peer.ID,
) (types.BroadcastMessage, bool) {
	pm.pendingMessagesLock.RLock()
	defer pm.pendingMessagesLock.RUnlock()

	peerMsgGroup, found := pm.pendingMessages[msgID]
	if !found {
		return types.BroadcastMessage{}, false
	}

	if !peerMsgGroup.Completed(hostID, recipients) {
		return types.BroadcastMessage{}, false
	}

	msgData, found := peerMsgGroup.GetMessageData()
	if !found {
		log.DPanicf("message data not found for completed message ID %s", msgID)
		return types.BroadcastMessage{}, false
	}

	return msgData, true
}
