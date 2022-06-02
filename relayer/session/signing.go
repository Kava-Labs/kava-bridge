package session

import (
	"fmt"
	"math/big"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	eth_common "github.com/ethereum/go-ethereum/common"
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

	transport mp_tss.Transporter

	outputChan chan tss_common.SignatureData
	errChan    chan *tss.Error
}

func (s *SigningSession) Initialize() {
	// Generate a random session ID part

	// Pick leader

	// If not leader send the session ID part to leader

	// If leader, wait for all parties to send their session ID parts
}

// Start starts the session in the background once the leader broadcasts a
// SigningPartyStartMessage.
func (s *SigningSession) StartSigner(
	key keygen.LocalPartySaveData,
	transport mp_tss.Transporter,
) error {
	s.transport = transport

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

// WaitForSignature returns the signature data from the session when completed.
func (s *SigningSession) WaitForSignature() (tss_common.SignatureData, error) {
	select {
	case sigData := <-s.outputChan:
		return sigData, nil
	case err := <-s.errChan:
		return tss_common.SignatureData{}, err
	}
}

/// --- fsm

func (s *SigningSession) Update(event SigningSessionEvent) error {
	switch ev := event.(type) {
	case *AddCandidateEvent:
		ev.Candidate
	case *StartSignerEvent:
	case *AddSigningPartEvent:
	default:
		return fmt.Errorf("unexpected event type: %T", event)
	}

	return nil
}

// States

type SigningSessionStateType int

const (
	SigningSessionStateType_LeaderWaitingForCandidates SigningSessionStateType = iota + 1
	SigningSessionStateType_CandidateWaitingForLeader
)

type SigningSessionState interface {
	State() SigningSessionStateType
}

var _ SigningSessionState = (*LeaderWaitingForCandidatesState)(nil)
var _ SigningSessionState = (*LeaderWaitingForCandidatesState)(nil)
var _ SigningSessionState = (*LeaderWaitingForCandidatesState)(nil)

type LeaderWaitingForCandidatesState struct {
	threshold int
	parts     map[peer.ID]types.SigningSessionIDPart
}

func (s *LeaderWaitingForCandidatesState) State() SigningSessionStateType {
	return SigningSessionStateType_LeaderWaitingForCandidates
}

type CandidatesWaitingForLeaderState struct {
	part types.SigningSessionIDPart
}

func (s *CandidatesWaitingForLeaderState) State() SigningSessionStateType {
	return SigningSessionStateType_CandidateWaitingForLeader
}

// Events

type SigningSessionEventType int

const (
	SigningSessionEventType_AddCandidate SigningSessionEventType = iota + 1
	SigningSessionEventType_StartSigner
	SigningSessionEventType_AddSigningPart
)

type SigningSessionEvent interface {
	EventType() SigningSessionEventType
}

var _ SigningSessionEvent = (*AddCandidateEvent)(nil)
var _ SigningSessionEvent = (*StartSignerEvent)(nil)
var _ SigningSessionEvent = (*AddSigningPartEvent)(nil)

type AddCandidateEvent struct {
	Candidate *tss.PartyID
}

// EventType
func (e *AddCandidateEvent) EventType() SigningSessionEventType {
	return SigningSessionEventType_AddCandidate
}

type StartSignerEvent struct {
	key       keygen.LocalPartySaveData
	transport mp_tss.Transporter
}

func (e *StartSignerEvent) EventType() SigningSessionEventType {
	return SigningSessionEventType_StartSigner
}

type AddSigningPartEvent struct {
	From        peer.ID
	Data        []byte
	IsBroadcast bool
}

func (e *AddSigningPartEvent) EventType() SigningSessionEventType {
	return SigningSessionEventType_AddSigningPart
}
