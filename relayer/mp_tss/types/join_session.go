package types

import (
	fmt "fmt"
	"math/rand"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

var (
	_ types.PeerMessage = (*JoinSessionMessage)(nil)
)

// NewJoinSessionMessage creates a new JoinSessionMessage.
func NewJoinSessionMessage(
	peerID peer.ID,
	session isJoinSessionMessage_Session,
) JoinSessionMessage {
	return JoinSessionMessage{
		PeerID:  peerID,
		Session: session,
	}
}

// NewJoinSigningSessionMessage creates a new signing JoinSessionMessage.
func NewJoinSigningSessionMessage(
	peerID peer.ID,
	tx_hash common.Hash,
	session_id_part SigningSessionIDPart,
) JoinSessionMessage {
	return NewJoinSessionMessage(
		peerID,
		&JoinSessionMessage_JoinSigningSessionMessage{
			JoinSigningSessionMessage: &JoinSigningSessionMessage{
				TxHash:            tx_hash.Bytes(),
				PeerSessionIDPart: session_id_part[:],
			},
		},
	)
}

// NewJoinKeygenSessionMessage creates a new keygen JoinSessionMessage.
func NewJoinKeygenSessionMessage(
	peerID peer.ID,
	keygen_session_id KeygenSessionID,
) JoinSessionMessage {
	return NewJoinSessionMessage(
		peerID,
		&JoinSessionMessage_JoinKeygenSessionMessage{
			JoinKeygenSessionMessage: &JoinKeygenSessionMessage{
				KeygenSessionID: keygen_session_id[:],
			},
		},
	)
}

// NewJoinReSharingSessionMessage creates a new resharing JoinSessionMessage.
func NewJoinReSharingSessionMessage(
	peerID peer.ID,
	resharing_session_id ReSharingSessionID,
) JoinSessionMessage {
	return NewJoinSessionMessage(
		peerID,
		&JoinSessionMessage_JoinResharingSessionMessage{
			JoinResharingSessionMessage: &JoinReSharingSessionMessage{
				ReSharingSessionID: resharing_session_id[:],
			},
		},
	)
}

// ValidateBasic does a simple validation check that doesn't require access to
// any other information.
func (msg *JoinSessionMessage) ValidateBasic() error {
	switch session := msg.Session.(type) {
	case *JoinSessionMessage_JoinSigningSessionMessage:
		if err := session.JoinSigningSessionMessage.ValidateBasic(); err != nil {
			return err
		}
	case *JoinSessionMessage_JoinKeygenSessionMessage:
		if err := session.JoinKeygenSessionMessage.ValidateBasic(); err != nil {
			return err
		}
	case *JoinSessionMessage_JoinResharingSessionMessage:
		if err := session.JoinResharingSessionMessage.ValidateBasic(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid session type: %T", session)
	}

	return nil
}

// GetSenderPeerID returns the peer ID of the sender.
func (msg *JoinSessionMessage) GetSenderPeerID() peer.ID {
	return msg.PeerID
}

// ValidateBasic does a simple validation check that doesn't require access to
// any other information.
func (msg *JoinSigningSessionMessage) ValidateBasic() error {
	if len(msg.TxHash) != common.HashLength {
		return fmt.Errorf(
			"invalid tx hash length: expected %d, got %d",
			common.HashLength,
			len(msg.TxHash),
		)
	}

	if len(msg.PeerSessionIDPart) != SigningSessionIDPartLength {
		return fmt.Errorf(
			"invalid peer session ID part length: expected %d, got %d",
			SigningSessionIDPartLength,
			len(msg.PeerSessionIDPart),
		)
	}

	return nil
}

// GetTxHash returns the transaction hash.
func (msg *JoinSigningSessionMessage) GetTxHash() common.Hash {
	return common.BytesToHash(msg.TxHash)
}

// GetPeerSessionIDPart returns the SigningSessionIDPart
func (msg *JoinSigningSessionMessage) GetPeerSessionIDPart() SigningSessionIDPart {
	sessionPart := SigningSessionIDPart{}
	copy(sessionPart[:], msg.PeerSessionIDPart)

	return sessionPart
}

// ValidateBasic does a simple validation check that doesn't require access to
// any other information.
func (msg *JoinKeygenSessionMessage) ValidateBasic() error {
	if len(msg.KeygenSessionID) != KeygenSessionIDLength {
		return fmt.Errorf(
			"keygen session ID length incorrect: expected %d, got %d",
			KeygenSessionIDLength,
			len(msg.KeygenSessionID),
		)
	}

	return nil
}

// ValidateBasic does a simple validation check that doesn't require access to
// any other information.
func (msg *JoinReSharingSessionMessage) ValidateBasic() error {
	if len(msg.ReSharingSessionID) != ReSharingSessionIDLength {
		return fmt.Errorf(
			"resharing session ID length incorrect: expected %d, got %d",
			ReSharingSessionIDLength,
			len(msg.ReSharingSessionID),
		)
	}

	return nil
}

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

	var txHash *common.Hash
	var sessionID []byte
	pickedPeerIDs := make(peer.IDSlice, 0, len(pickedMsgs))

	for _, msg := range pickedMsgs {
		// Only for signing messages
		signingMsg := msg.GetJoinSigningSessionMessage()
		if signingMsg == nil {
			return nil, nil, fmt.Errorf("invalid join message type: %T", msg)
		}

		// All signing messages must have the same tx hash
		if txHash == nil {
			msgTxHash := signingMsg.GetTxHash()
			txHash = &msgTxHash
		} else if *txHash != signingMsg.GetTxHash() {
			return nil, nil, fmt.Errorf("mismatch tx hash")
		}

		// Append the session ID
		sessionID = append(sessionID, signingMsg.PeerSessionIDPart...)
		// Append picked peer ID
		pickedPeerIDs = append(pickedPeerIDs, msg.PeerID)
	}

	return sessionID, pickedPeerIDs, nil
}
