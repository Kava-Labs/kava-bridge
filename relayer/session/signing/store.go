package signing

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
)

// tx_hash -> session ID
type TssSessions map[common.Hash]Session

// session -> tx_hash to get a session from session ID instead of tx_hash
type SessionIDToTxHash map[string]common.Hash

type SigningSessionStore struct {
	sessions          TssSessions
	sessionIDToTxHash SessionIDToTxHash
}

func NewSigningSessionStore() *SigningSessionStore {
	return &SigningSessionStore{
		sessions:          make(TssSessions),
		sessionIDToTxHash: make(SessionIDToTxHash),
	}
}

func (s *SigningSessionStore) NewSession(txHash common.Hash) {
	s.sessions[txHash] = LeaderSelectState{}
}

func (s *SigningSessionStore) GetSessionFromTxHash(txHash common.Hash) (Session, bool) {
	session, ok := s.sessions[txHash]
	return session, ok
}

func (s *SigningSessionStore) GetSessionFromID(sessID string) (Session, bool) {
	txHash, ok := s.sessionIDToTxHash[sessID]
	if !ok {
		return nil, false
	}

	return s.GetSessionFromTxHash(txHash)
}

type Session interface {
}

type LeaderSelectState struct {
	// Randomly locally generated ID part that is sent to leader
	LocalIDPart types.SigningSessionIDPart

	TxHash common.Hash
}
