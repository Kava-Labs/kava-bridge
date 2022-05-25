package signing

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	mp_tss_types "github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
)

type AggregateSigningSessionID []byte

// GetAggregateSigningSessionID returns the aggregate session ID for the given
// signing session.
func GetAggregateSigningSessionID(
	joinMsgs mp_tss_types.JoinSessionMessages,
) (AggregateSigningSessionID, error) {
	if len(joinMsgs) == 0 {
		return nil, fmt.Errorf("no join messages provided")
	}

	// Copy join messages to avoid mutating the original.
	joinMsgsCopy := make(mp_tss_types.JoinSessionMessages, len(joinMsgs))
	copy(joinMsgsCopy, joinMsgs)

	// Sort the join messages by their peer IDs.
	sort.Sort(joinMsgsCopy)

	var txHash *common.Hash
	var sessionID []byte

	for _, msg := range joinMsgsCopy {
		// Only for signing messages
		signingMsg := msg.GetJoinSigningSessionMessage()
		if signingMsg == nil {
			return nil, ErrInvalidSessionType
		}

		// All signing messages must have the same tx hash
		if txHash == nil {
			msgTxHash := signingMsg.GetTxHash()
			txHash = &msgTxHash
		} else if *txHash != signingMsg.GetTxHash() {
			return nil, ErrMismatchedTxHash
		}

		// Append the session ID
		sessionID = append(sessionID, signingMsg.PeerSessionIDPart...)
	}

	return sessionID, nil
}

// IsPeerParticipant returns true if the given peer is a signer for the given
// aggregate session ID.
func IsPeerParticipant(
	peer_session_id_part mp_tss_types.SigningSessionIDPart,
	sessionID AggregateSigningSessionID,
) bool {
	for i := 0; i < len(sessionID); i += mp_tss_types.SigningSessionIDPartLength {
		chunk := sessionID[i : i+mp_tss_types.SigningSessionIDPartLength]

		// If the current peer's session ID part is contained in the aggregate
		if bytes.Equal(chunk, peer_session_id_part[:]) {
			return true
		}
	}

	return false
}
