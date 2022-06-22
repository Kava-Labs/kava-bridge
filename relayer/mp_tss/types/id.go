package types

const (
	SigningSessionIDPartLength = 32
	KeygenSessionIDLength      = 32
	ReSharingSessionIDLength   = 32
)

// SigningSessionIDPart is a peer part of signing session ID.
type SigningSessionIDPart [SigningSessionIDPartLength]byte

// Bytes returns the bytes of the SigningSessionIDPart.
func (s SigningSessionIDPart) Bytes() []byte {
	return s[:]
}

// KeygenSessionID is the ID for a keygen session.
type KeygenSessionID [KeygenSessionIDLength]byte

// ReSharingSessionID is the ID for a resharing session.
type ReSharingSessionID [ReSharingSessionIDLength]byte
