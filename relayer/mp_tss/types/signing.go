package types

import "fmt"

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
