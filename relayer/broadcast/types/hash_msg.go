package types

import (
	"fmt"
	"strings"

	"github.com/multiformats/go-multibase"
)

var _ Message = (*HashMsg)(nil)

func (msg *HashMsg) Validate() error {
	if strings.TrimSpace(msg.MessageID) == "" {
		return ErrMsgIDEmpty
	}

	_, _, err := multibase.Decode(msg.MessageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	return nil
}

func (msg *HashMsg) GetHash() BroadcastMessageHash {
	var hash BroadcastMessageHash
	copy(hash[:], msg.Hash)

	return hash
}
