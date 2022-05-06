package pending_store_test

import (
	"testing"
	"time"

	"github.com/kava-labs/kava-bridge/relayer/broadcast/pending_store"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/require"
)

func TestContainsGroup(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID := "test-msg-id"

	found := store.ContainsGroup(msgID)
	require.False(t, found, "should not find a group that doesn't exist")

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	found = store.ContainsGroup(msgID)
	require.True(t, found, "should be able to find a group that exists")
}

func TestTryNewGroup(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID := "test-msg-id"

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	created = store.TryNewGroup(msgID)
	require.False(t, created, "should not be able to create a group that already exists")
}

func TestDeleteGroup(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID := "test-msg-id"

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	err := store.DeleteGroup(msgID)
	require.NoError(t, err)

	err = store.DeleteGroup(msgID)
	require.ErrorIs(t, err, pending_store.ErrGroupNotFound, "should not be able to delete a group that doesn't exist")
}

func TestAddMessage_GroupNotExists(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	err = store.AddMessage(types.BroadcastMessage{
		ID: msgID,
	})
	require.ErrorIs(t, err, pending_store.ErrGroupNotFound, "should not be able to add to a group that doesn't exist")
}

func TestAddMessage_GroupExists(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	err = store.AddMessage(types.BroadcastMessage{
		ID: msgID,
	})
	require.NoError(t, err, "should not error when adding to a group that exists")
}

func TestAddMessage_DuplicateMessage(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	err = store.AddMessage(types.BroadcastMessage{
		ID: msgID,
	})
	require.NoError(t, err)

	// Invalid message should error
	err = store.AddMessage(types.BroadcastMessage{
		ID: msgID,
	})
	require.Error(t, err)
	require.Equal(t, "message already received", err.Error())
}

func TestGroupIsCompleted_NotExist(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	_, complete := store.GroupIsCompleted(msgID)
	require.False(t, complete, "should not be complete if group doesn't exist")
}

func TestGroupIsCompleted_Incomplete(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(pending_store.DEFAULT_CLEAR_EXPIRED_INTERVAL)

	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	msg := types.BroadcastMessage{
		ID: msgID,
		// Only 1 recipient
		RecipientPeerIDs: testutil.TestPeerIDs[1:2],
	}

	err = store.AddMessage(msg)
	require.NoError(t, err)

	// Requires 2
	_, complete := store.GroupIsCompleted(msgID)
	require.False(t, complete, "should not be complete if group is incomplete")

	msgHash, err := msg.Hash()
	require.NoError(t, err)

	err = store.AddHash(testutil.TestPeerIDs[2], types.NewHashMsg(msgID, msgHash))
	require.NoError(t, err)

	_, complete = store.GroupIsCompleted(msgID)
	require.True(t, complete, "should be complete when recipient hashes added")
}

func TestRemovesExpiredGroups(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(1 * time.Second)

	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	err = store.AddMessage(types.BroadcastMessage{
		ID:         msgID,
		Created:    time.Now().Add(-time.Hour),
		TTLSeconds: 1,
	})
	require.NoError(t, err)

	contains := store.ContainsGroup(msgID)
	require.True(t, contains, "should contain group after creation")

	time.Sleep(2 * time.Second)

	contains = store.ContainsGroup(msgID)
	require.False(t, contains, "should delete expired groups")
}

func TestKeepsNonExpiredGroups(t *testing.T) {
	store := pending_store.NewPendingMessagesStore(1 * time.Second)

	msgID, err := types.NewBroadcastMessageID()
	require.NoError(t, err)

	created := store.TryNewGroup(msgID)
	require.True(t, created)

	err = store.AddMessage(types.BroadcastMessage{
		ID:         msgID,
		Created:    time.Now(),
		TTLSeconds: 4,
	})
	require.NoError(t, err)

	contains := store.ContainsGroup(msgID)
	require.True(t, contains, "should contain group after creation")

	time.Sleep(2 * time.Second)

	contains = store.ContainsGroup(msgID)
	require.True(t, contains, "should not delete groups not yet expired")
}
