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

// tx_hash -> session ID
type TssSessions map[eth_common.Hash]*SigningSession

// session -> tx_hash to get a session from session ID instead of tx_hash
type SessionIDToTxHash map[string]eth_common.Hash

// SigningSessionStore keeps track of signing sessions.
type SigningSessionStore struct {
	mu                *sync.Mutex
	sessions          TssSessions
	sessionIDToTxHash SessionIDToTxHash
}

func NewSigningSessionStore() *SigningSessionStore {
	return &SigningSessionStore{
		mu:                &sync.Mutex{},
		sessions:          make(TssSessions),
		sessionIDToTxHash: make(SessionIDToTxHash),
	}
}

func (s *SigningSessionStore) NewSession(
	ctx context.Context,
	broadcaster *broadcast.Broadcaster,
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

func (s *SigningSessionStore) GetSessionFromTxHash(txHash eth_common.Hash) (*SigningSession, bool) {
	s.mu.Lock()
	// Session is returned as pointer, still possible for concurrent access to
	// specific sessions.
	session, ok := s.sessions[txHash]
	s.mu.Unlock()

	return session, ok
}

func (s *SigningSessionStore) SetSessionID(txHash eth_common.Hash, sessID types.AggregateSigningSessionID) {
	s.mu.Lock()
	s.sessionIDToTxHash[sessID.String()] = txHash
	s.mu.Unlock()
}

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
