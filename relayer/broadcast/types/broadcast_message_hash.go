package types

import "golang.org/x/crypto/sha3"

// Bytes returns the bytes of the entire broadcast message. Must be
// deterministic for the same message.
func (m *BroadcastMessage) Bytes() ([]byte, error) {
	return m.Marshal()
}

// Hash returns the hash of the message.
func (m *BroadcastMessage) Hash() ([32]byte, error) {
	bytes, err := m.Bytes()
	if err != nil {
		return [32]byte{}, err
	}

	return sha3.Sum256(bytes), nil
}
