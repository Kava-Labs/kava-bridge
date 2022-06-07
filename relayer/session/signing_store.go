package session

import (
	"context"
	"math/big"

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
	sessions          TssSessions
	sessionIDToTxHash SessionIDToTxHash
}

func NewSigningSessionStore() *SigningSessionStore {
	return &SigningSessionStore{
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
	partyIDStore *mp_tss.PartyIDStore,
) (*SigningSession, <-chan SigningSessionResult, error) {
	session, resultChan, err := NewSigningSession(
		ctx,
		broadcaster,
		txHash,
		msgToSign,
		threshold,
		currentPeerID,
		peerIDs,
		partyIDStore,
	)
	if err != nil {
		return nil, nil, err
	}

	s.sessions[txHash] = session

	return session, resultChan, nil
}

func (s *SigningSessionStore) GetSessionFromTxHash(txHash eth_common.Hash) (*SigningSession, bool) {
	session, ok := s.sessions[txHash]
	return session, ok
}

func (s *SigningSessionStore) SetSessionID(txHash eth_common.Hash, sessID types.AggregateSigningSessionID) {
	s.sessionIDToTxHash[sessID.String()] = txHash
}

func (s *SigningSessionStore) GetSessionFromID(
	sessID types.AggregateSigningSessionID,
) (*SigningSession, bool) {
	txHash, ok := s.sessionIDToTxHash[sessID.String()]
	if !ok {
		return nil, false
	}

	return s.GetSessionFromTxHash(txHash)
}
