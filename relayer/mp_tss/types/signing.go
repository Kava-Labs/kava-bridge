package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p-core/peer"
)

func NewSigningPartMessage(
	sessionID AggregateSigningSessionID,
	data []byte,
	isBroadcast bool,
) SigningPartMessage {
	return SigningPartMessage{
		SessionID:   sessionID,
		Data:        data,
		IsBroadcast: isBroadcast,
	}
}

func (m *SigningPartMessage) GetSessionID() AggregateSigningSessionID {
	return m.SessionID
}

func (m *SigningPartMessage) ValidateBasic() error {
	// TODO: Replace with AggregateSigningSessionID
	if m.SessionID == nil {
		return fmt.Errorf("session id is nil")
	}

	if err := m.GetSessionID().Validate(); err != nil {
		return fmt.Errorf("session id is invalid: %w", err)
	}

	if m.Data == nil {
		return fmt.Errorf("data is nil")
	}

	return nil
}

// NewSigningPartyStartMessage returns a new SigningPartyStartMessage.
func NewSigningPartyStartMessage(
	txHash common.Hash,
	sessionID AggregateSigningSessionID,
	participatingPeerIDs []peer.ID,
) *SigningPartyStartMessage {
	return &SigningPartyStartMessage{
		TxHash:               txHash.Bytes(),
		SessionId:            sessionID,
		ParticipatingPeerIDs: participatingPeerIDs,
	}
}

func (m *SigningPartyStartMessage) GetSessionID() AggregateSigningSessionID {
	return m.SessionId
}

func (m *SigningPartyStartMessage) ValidateBasic() error {
	if m.TxHash == nil {
		return fmt.Errorf("txHash is nil")
	}

	if err := m.GetSessionID().Validate(); err != nil {
		return fmt.Errorf("session id is invalid: %w", err)
	}

	if len(m.ParticipatingPeerIDs) == 0 {
		return fmt.Errorf("participants are empty")
	}

	return nil
}
