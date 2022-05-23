package types

import (
	fmt "fmt"

	"github.com/ethereum/go-ethereum/common"
)

const (
	PeerSessionIDPartLength  = 32
	KeygenSessionIDLength    = 32
	ReSharingSessionIDLength = 32
)

// PeerSessionIDPart is a peer part of signing session ID.
type PeerSessionIDPart [PeerSessionIDPartLength]byte

// KeygenSessionID is the ID for a keygen session.
type KeygenSessionID [KeygenSessionIDLength]byte

// ReSharingSessionID is the ID for a resharing session.
type ReSharingSessionID [ReSharingSessionIDLength]byte

var (
	_ TssMsg = (*JoinSessionMessage)(nil)
	_ TssMsg = (*JoinSigningSessionMessage)(nil)
	_ TssMsg = (*JoinKeygenSessionMessage)(nil)
	_ TssMsg = (*JoinReSharingSessionMessage)(nil)
)

// NewJoinSessionMessage creates a new JoinSessionMessage.
func NewJoinSessionMessage(session isJoinSessionMessage_Session) JoinSessionMessage {
	return JoinSessionMessage{Session: session}
}

// NewJoinSigningSessionMessage creates a new signing JoinSessionMessage.
func NewJoinSigningSessionMessage(
	tx_hash common.Hash,
	session_id_part PeerSessionIDPart,
) JoinSessionMessage {
	return NewJoinSessionMessage(&JoinSessionMessage_JoinSigningSessionMessage{
		JoinSigningSessionMessage: &JoinSigningSessionMessage{
			TxHash:            tx_hash.Bytes(),
			PeerSessionIDPart: session_id_part[:],
		},
	})
}

// NewJoinKeygenSessionMessage creates a new keygen JoinSessionMessage.
func NewJoinKeygenSessionMessage(keygen_session_id KeygenSessionID) JoinSessionMessage {
	return NewJoinSessionMessage(&JoinSessionMessage_JoinKeygenSessionMessage{
		JoinKeygenSessionMessage: &JoinKeygenSessionMessage{
			KeygenSessionID: keygen_session_id[:],
		},
	})
}

// NewJoinReSharingSessionMessage creates a new resharing JoinSessionMessage.
func NewJoinReSharingSessionMessage(resharing_session_id ReSharingSessionID) JoinSessionMessage {
	return NewJoinSessionMessage(&JoinSessionMessage_JoinResharingSessionMessage{
		JoinResharingSessionMessage: &JoinReSharingSessionMessage{
			ReSharingSessionID: resharing_session_id[:],
		},
	})
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

	if len(msg.PeerSessionIDPart) != PeerSessionIDPartLength {
		return fmt.Errorf(
			"invalid peer session ID part length: expected %d, got %d",
			PeerSessionIDPartLength,
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
