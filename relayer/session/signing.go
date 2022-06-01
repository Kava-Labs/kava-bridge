package session

import (
	"math/big"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
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

func (s *SigningSessionStore) NewSession(txHash eth_common.Hash) *SigningSession {
	session := &SigningSession{}

	s.sessions[txHash] = session

	return session
}

func (s *SigningSessionStore) GetSessionFromTxHash(txHash eth_common.Hash) (*SigningSession, bool) {
	session, ok := s.sessions[txHash]
	return session, ok
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

type SigningSession struct {
	// Randomly locally generated ID part that is sent to leader
	LocalIDPart types.SigningSessionIDPart

	tssParams *tss.Parameters

	TxHash    eth_common.Hash
	MsgToSign *big.Int

	outputChan chan tss_common.SignatureData
	errChan    chan *tss.Error
}

// Start starts the session in the background
func (s *SigningSession) StartSigner(
	key keygen.LocalPartySaveData,
	transport mp_tss.Transporter,
) error {
	outputChan, errChan := mp_tss.RunSign(
		s.MsgToSign,
		s.tssParams,
		key,
		transport,
	)

	s.outputChan = outputChan
	s.errChan = errChan

	return nil
}
