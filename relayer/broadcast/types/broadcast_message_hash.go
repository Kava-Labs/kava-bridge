package types

import (
	"bytes"
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

// BroadcastMessageHash is the SHA3 256 hash of a broadcast message.
type BroadcastMessageHash [32]byte

// Equal returns true if the provided hash is equal.
func (m *BroadcastMessageHash) Equal(other BroadcastMessageHash) bool {
	return bytes.Equal(m[:], other[:])
}

// String returns the hex encoded hash.
func (m *BroadcastMessageHash) String() string {
	return hex.EncodeToString(m[:])
}

// Bytes returns the bytes of the entire broadcast message. Must be
// deterministic for the same message.
func (m *BroadcastMessage) Bytes() ([]byte, error) {
	return m.Marshal()
}

// Hash returns the hash of the message.
func (m *BroadcastMessage) Hash() (BroadcastMessageHash, error) {
	bytes, err := m.Bytes()
	if err != nil {
		return BroadcastMessageHash{}, err
	}

	return sha3.Sum256(bytes), nil
}
