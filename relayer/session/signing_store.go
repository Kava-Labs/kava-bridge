package session

import (
	"context"
	"math/big"
	"sync"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// TssSessions maps a transaction hash to signing session
type TssSessions map[eth_common.Hash]*SigningSession

// SessionIDToTxHash maps session ID -> tx_hash, to get a session from session
// ID instead of tx_hash
type SessionIDToTxHash map[string]eth_common.Hash

// SigningSessionStore keeps track of signing sessions.
type SigningSessionStore struct {
	mu                *sync.Mutex
	sessions          TssSessions
	sessionIDToTxHash SessionIDToTxHash
}

// NewSigningSessionStore returns a new signing session store.
func NewSigningSessionStore() *SigningSessionStore {
	return &SigningSessionStore{
		mu:                &sync.Mutex{},
		sessions:          make(TssSessions),
		sessionIDToTxHash: make(SessionIDToTxHash),
	}
}

// NewSession adds and returns a new signing session. This session does not have
// a session ID and must be set later with SetSessionID.
func (s *SigningSessionStore) NewSession(
	ctx context.Context,
	broadcaster broadcast.Broadcaster,
	txHash eth_common.Hash,
	msgToSign *big.Int,
	threshold int,
	currentPeerID peer.ID,
	peerIDs peer.IDSlice,
	currentPartyID *tss.PartyID,
	partyIDStore *mp_tss.PartyIDStore,
	key keygen.LocalPartySaveData,
) (*SigningSession, <-chan SigningSessionResult, error) {
	session, resultChan, err := NewSigningSession(
		ctx,
		s,
		broadcaster,
		txHash,
		msgToSign,
		threshold,
		currentPeerID,
		peerIDs,
		currentPartyID,
		partyIDStore,
		key,
	)
	if err != nil {
		return nil, nil, err
	}

	s.mu.Lock()
	s.sessions[txHash] = session
	s.mu.Unlock()

	return session, resultChan, nil
}

// GetSessionFromTxHash returns the signing session for the given transaction hash.
func (s *SigningSessionStore) GetSessionFromTxHash(txHash eth_common.Hash) (*SigningSession, bool) {
	s.mu.Lock()
	// Session is returned as pointer, still possible for concurrent access to
	// specific sessions.
	session, ok := s.sessions[txHash]
	s.mu.Unlock()

	return session, ok
}

// SetSessionID sets the session ID for the given signing session transaction hash.
func (s *SigningSessionStore) SetSessionID(txHash eth_common.Hash, sessID types.AggregateSigningSessionID) {
	s.mu.Lock()
	s.sessionIDToTxHash[sessID.String()] = txHash
	s.mu.Unlock()
}

// GetSessionFromID returns the signing session for the given session ID.
func (s *SigningSessionStore) GetSessionFromID(
	sessID types.AggregateSigningSessionID,
) (*SigningSession, bool) {
	s.mu.Lock()
	txHash, ok := s.sessionIDToTxHash[sessID.String()]
	s.mu.Unlock()

	if !ok {
		return nil, false
	}

	return s.GetSessionFromTxHash(txHash)
}
