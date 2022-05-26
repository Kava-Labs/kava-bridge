package types

import (
	fmt "fmt"
	"math/rand"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p-core/peer"
)

// JoinSessionMessages is a slice of JoinSessionMessages.
type JoinSessionMessages []JoinSessionMessage

func (a JoinSessionMessages) Len() int           { return len(a) }
func (a JoinSessionMessages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a JoinSessionMessages) Less(i, j int) bool { return a[i].PeerID < a[j].PeerID }

// ValidateBasic does a simple validation check that doesn't require access to
// any other information.
func (jMsgs JoinSessionMessages) ValidateBasic() error {
	if len(jMsgs) == 0 {
		return fmt.Errorf("empty join session messages")
	}

	uniquePeerIDs := make(map[peer.ID]struct{})
	uniqueTxHashes := make(map[common.Hash]struct{})
	uniqueIDParts := make(map[SigningSessionIDPart]struct{})

	for _, jMsg := range jMsgs {
		if err := jMsg.ValidateBasic(); err != nil {
			return err
		}

		if jMsg.GetJoinSigningSessionMessage() == nil {
			return fmt.Errorf("invalid join session type: %T", jMsg.Session)
		}

		// Each peer can only have one join message.
		if _, found := uniquePeerIDs[jMsg.PeerID]; found {
			return fmt.Errorf("duplicate peer ID: %s", jMsg.PeerID)
		}
		uniquePeerIDs[jMsg.PeerID] = struct{}{}

		// If uniqueTxHashes is not empty, and this is a new hash, then we have
		// 2 different hashes. All hashes must be the same.
		txHash := jMsg.GetJoinSigningSessionMessage().GetTxHash()
		if _, found := uniqueTxHashes[txHash]; !found && len(uniqueTxHashes) != 0 {
			return fmt.Errorf("different tx hashes: %s not in %v", txHash, uniqueTxHashes)
		}
		uniqueTxHashes[txHash] = struct{}{}

		// Each peer should generate random session part individually, if there
		// are duplicate parts then there must be something wrong.
		idPart := jMsg.GetJoinSigningSessionMessage().GetPeerSessionIDPart()
		if _, found := uniqueIDParts[idPart]; found {
			return fmt.Errorf("duplicate peer session ID part: %s", idPart)
		}
		uniqueIDParts[idPart] = struct{}{}
	}

	return nil
}

// GetSessionID returns the session ID of a signing session with a randomized
// subset of participants.
func (jMsgs JoinSessionMessages) GetSessionID(threshold int) (
	AggregateSigningSessionID,
	peer.IDSlice,
	error,
) {
	if err := jMsgs.ValidateBasic(); err != nil {
		return nil, nil, err
	}

	if len(jMsgs) < threshold+1 {
		return nil, nil, fmt.Errorf(
			"not enough peers to select participants, %d (peers) < %d (t + 1)",
			len(jMsgs), threshold+1,
		)
	}

	if threshold < 1 {
		return nil, nil, fmt.Errorf("invalid threshold: %d", threshold)
	}

	// Copy allPeerIDs to avoid mutating the original.
	jMsgsCopy := make(JoinSessionMessages, len(jMsgs))
	copy(jMsgsCopy, jMsgs)

	rand.Shuffle(len(jMsgsCopy), func(i, j int) {
		jMsgsCopy[i], jMsgsCopy[j] = jMsgsCopy[j], jMsgsCopy[i]
	})

	pickedMsgs := jMsgsCopy[:threshold+1]

	// Sort the join messages by their peer IDs.
	sort.Sort(pickedMsgs)

	var sessionID []byte
	pickedPeerIDs := make(peer.IDSlice, 0, len(pickedMsgs))

	for _, msg := range pickedMsgs {
		// Only for signing messages, same hash for all messages already
		// already validated in ValidateBasic()
		signingMsg := msg.GetJoinSigningSessionMessage()

		// Append the session ID
		sessionID = append(sessionID, signingMsg.PeerSessionIDPart...)
		// Append picked peer ID
		pickedPeerIDs = append(pickedPeerIDs, msg.PeerID)
	}

	return sessionID, pickedPeerIDs, nil
}
