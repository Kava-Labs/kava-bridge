package signing

import (
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	mp_tss_types "github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
)

// GetAggregateSigningSessionID returns the aggregate session ID for the given
// signing session.
func GetAggregateSigningSessionID(
	joinMsgs mp_tss_types.JoinSessionMessages,
) ([]byte, error) {
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
