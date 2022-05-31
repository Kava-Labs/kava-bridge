package types

import (
	"bytes"
)

// AggregateSigningSessionID is a signing session ID, consisting of sorted and
// concatenated session ID parts from each participating peer.
type AggregateSigningSessionID []byte

// Validate returns an error if the session ID is an invalid length.
func (sid AggregateSigningSessionID) Validate() bool {
	return len(sid)%SigningSessionIDPartLength == 0
}

// Bytes returns the byte representation of the aggregate signing session ID.
func (sid AggregateSigningSessionID) Bytes() []byte {
	return sid[:]
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
