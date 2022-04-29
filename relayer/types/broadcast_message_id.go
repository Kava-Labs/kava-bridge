package types

import (
	"crypto/rand"
	"fmt"

	"github.com/multiformats/go-multibase"
)

const (
	MessageIDLengthBytes = 32
)

// NewBroadcastMessageID returns a new broadcast message ID. This consists of
// random 32 bytes base58 encoded.
func NewBroadcastMessageID() (string, error) {
	b := make([]byte, MessageIDLengthBytes)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("could not read from rand: %w", err)
	}

	s, err := multibase.Encode(multibase.Base58BTC, b)
	if err != nil {
		return "", fmt.Errorf("could not encode random bytes: %w", err)
	}

	return s, nil
}
