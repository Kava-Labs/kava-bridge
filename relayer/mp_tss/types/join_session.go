package types

import (
	fmt "fmt"

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
