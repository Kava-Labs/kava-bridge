package types

import (
	"bytes"
	"encoding/hex"
	fmt "fmt"
)

// AggregateSigningSessionID is a signing session ID, consisting of sorted and
// concatenated session ID parts from each participating peer.
type AggregateSigningSessionID []byte

// Validate returns an error if the session ID is an invalid length.
func (sid AggregateSigningSessionID) Validate() error {
	if len(sid)%SigningSessionIDPartLength != 0 {
		return fmt.Errorf("invalid session ID length: %d", len(sid))
	}

	return nil
}

// Bytes returns the byte representation of the aggregate signing session ID.
func (sid AggregateSigningSessionID) Bytes() []byte {
	return sid[:]
}

// String returns the hex representation of the aggregate signing session ID.
func (sid AggregateSigningSessionID) String() string {
	return hex.EncodeToString(sid)
}

// IsPeerParticipant returns true if the given peer is a signer for the given
// aggregate session ID.
func (asid AggregateSigningSessionID) IsPeerParticipant(
	peer_session_id_part SigningSessionIDPart,
) bool {
	for i := 0; i < len(asid); i += SigningSessionIDPartLength {
		chunk := asid[i : i+SigningSessionIDPartLength]

		// If the current peer's session ID part is contained in the aggregate
		if bytes.Equal(chunk, peer_session_id_part[:]) {
			return true
		}
	}

	return false
}
