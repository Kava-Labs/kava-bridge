package signer

import (
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/pending_store"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
)

type BroadcastHandler struct {
	signer *Signer
}

var _ broadcast.BroadcastHandler = (*BroadcastHandler)(nil)

func (h *BroadcastHandler) RawMessage(msg pending_store.MessageWithPeerMetadata) {

}

func (h *BroadcastHandler) ValidatedMessage(msg types.BroadcastMessage) {
	h.signer.HandleBroadcastMessage(msg)
}

func (h *BroadcastHandler) MismatchMessage(msg pending_store.MessageWithPeerMetadata) {}
