package types

import (
	"encoding/hex"
	"math/rand"
)

const (
	SigningSessionIDPartLength = 32
	KeygenSessionIDLength      = 32
	ReSharingSessionIDLength   = 32
)

// SigningSessionIDPart is a peer part of signing session ID.
type SigningSessionIDPart [SigningSessionIDPartLength]byte

// NewSigningSessionIDPart returns a new random SigningSessionIDPart.
func NewSigningSessionIDPart() (SigningSessionIDPart, error) {
	bytes := make([]byte, SigningSessionIDPartLength)
	if _, err := rand.Read(bytes); err != nil {
		return SigningSessionIDPart{}, err
	}

	var part SigningSessionIDPart
	copy(part[:], bytes)

	return part, nil
}

// Bytes returns the bytes of the SigningSessionIDPart.
func (s SigningSessionIDPart) Bytes() []byte {
	return s[:]
}

// String returns the hex string of the SigningSessionIDPart.
func (s SigningSessionIDPart) String() string {
	return hex.EncodeToString(s[:])
}

// KeygenSessionID is the ID for a keygen session.
type KeygenSessionID [KeygenSessionIDLength]byte

// ReSharingSessionID is the ID for a resharing session.
type ReSharingSessionID [ReSharingSessionIDLength]byte
