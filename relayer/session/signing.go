package session

import (
	"context"
	"fmt"
	"math/big"

	tss_common "github.com/binance-chain/tss-lib/common"
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
	broadcaster *broadcast.Broadcaster,
	txHash eth_common.Hash,
	msgToSign *big.Int,
	currentPeerID peer.ID,
	peerIDs peer.IDSlice,
) (*SigningSession, error) {
	session := &SigningSession{}

	session, err := NewSigningSession(
		broadcaster,
		txHash,
		msgToSign,
		currentPeerID,
		peerIDs,
	)
	if err != nil {
		return nil, err
	}

	s.sessions[txHash] = session

	return session, nil
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
	broadcaster *broadcast.Broadcaster
	TxHash      eth_common.Hash
	MsgToSign   *big.Int

	currentPeerID peer.ID
	peerIDs       peer.IDSlice

	// FSM
	state SigningSessionState
}

func NewSigningSession(
	broadcaster *broadcast.Broadcaster,
	txHash eth_common.Hash,
	msgToSign *big.Int,
	currentPeerID peer.ID,
	peerIDs peer.IDSlice,
) (*SigningSession, error) {
	state := NewPickingLeaderState()

	session := &SigningSession{
		broadcaster:   broadcaster,
		TxHash:        txHash,
		MsgToSign:     msgToSign,
		currentPeerID: currentPeerID,
		peerIDs:       peerIDs,
		state:         state,
	}

	// Generate a random session ID part

	// Pick leader
	leaderPeerID, err := GetLeader(session.TxHash, session.peerIDs, state.leaderOffset)
	if err != nil {
		return nil, err
	}

	// Leader is the current peer, transition to LeaderWaitingForCandidatesState
	// No broadcast necessary, other peers know the leader too.
	if leaderPeerID == session.currentPeerID {
		session.state, err = NewLeaderWaitingForCandidatesState()
		if err != nil {
			return nil, err
		}

		return session, nil
	}

	// Not leader - send the session ID part to leader and transition to
	// CandidateWaitingForLeaderState
	newState, err := NewCandidateWaitingForLeaderState()
	if err != nil {
		return nil, err
	}

	session.state = newState

	msg := types.NewJoinSigningSessionMessage(
		currentPeerID,
		session.TxHash,
		newState.localPart,
	)

	err = session.broadcaster.BroadcastMessage(
		context.Background(),
		&msg,                    // join signing session message
		[]peer.ID{leaderPeerID}, // send to leader
		30,                      // ttl seconds
	)

	if err != nil {
		return nil, fmt.Errorf("failed to broadcast JoinSigningSessionMessage: %w", err)
	}

	return session, nil
}

// -----------------------------------------------------------------------------
// FSM

func (s *SigningSession) Update(event SigningSessionEvent) error {
	switch ev := event.(type) {
	case *AddCandidateEvent:
		return s.UpdateAddCandidateEvent(ev)
	case *StartSignerEvent:
		return s.UpdateStartSignerEvent(ev)
	case *AddSigningPartEvent:
	default:
		return fmt.Errorf("unexpected event type: %T", event)
	}

	return nil
}

func (s *SigningSession) UpdateAddCandidateEvent(
	ev *AddCandidateEvent,
) error {
	state, ok := s.state.(*LeaderWaitingForCandidatesState)
	if !ok {
		return fmt.Errorf("unexpected state type: %T", s.state)
	}

	// Update state
	_, found := state.parts[ev.peerID]
	if found {
		return fmt.Errorf("already added candidate")
	}

	state.parts[ev.peerID] = ev.sessionIDPart

	if len(state.parts) <= state.threshold {
		// Do nothing more, wait for more candidates
		return nil
	}

	// Greater than threshold, pick peers to participate in the signing session

	// Broadcast StartSignerEvent with picked peers

	// Transition to signing state
	s.state = NewSigningState()

	return nil
}

func (s *SigningSession) UpdateStartSignerEvent(
	ev *StartSignerEvent,
) error {
	// Only candidates receive this event, leader should not receive it and
	// will transition by itself.
	_, ok := s.state.(*CandidateWaitingForLeaderState)
	if !ok {
		return fmt.Errorf("unexpected state type: %T", s.state)
	}

	newState := NewSigningState(ev.transport)

	newState.outputChan, newState.errChan = mp_tss.RunSign(
		s.MsgToSign,
		ev.tssParams,
		ev.key,
		ev.transport,
	)

	// Transition CandidatesWaitingForLeaderState -> signing state
	s.state = newState

	return nil
}

func (s *SigningSession) AddSigningPartEvent(
	ev *AddSigningPartEvent,
) error {
	state, ok := s.state.(*SigningState)
	if !ok {
		return fmt.Errorf("unexpected state type: got %T, expected %T", s.state, &SigningState{})
	}

	// Receive signing part to transport
	state.transport.Receive() <- mp_tss.NewReceivedPartyState(
		ev.Data,
		ev.From,
		ev.IsBroadcast,
	)

	return nil
}

func (s *SigningSession) TryGetSignature() (tss_common.SignatureData, error) {
	state, ok := s.state.(*DoneState)
	if !ok {
		return tss_common.SignatureData{}, fmt.Errorf("unexpected state type: %T", s.state)
	}

	return state.signature, nil
}

// -----------------------------------------------------------------------------
// States

type SigningSessionStateType int

const (
	SigningSessionStateType_PickingLeader SigningSessionStateType = iota + 1
	SigningSessionStateType_LeaderWaitingForCandidates
	SigningSessionStateType_CandidateWaitingForLeader
	SigningSessionStateType_Signing
	SigningSessionStateType_Done
	SigningSessionStateType_Error
)

type SigningSessionState interface {
	State() SigningSessionStateType
}

var _ SigningSessionState = (*PickingLeaderState)(nil)
var _ SigningSessionState = (*LeaderWaitingForCandidatesState)(nil)
var _ SigningSessionState = (*CandidateWaitingForLeaderState)(nil)
var _ SigningSessionState = (*SigningState)(nil)
var _ SigningSessionState = (*DoneState)(nil)
var _ SigningSessionState = (*ErrorState)(nil)

type PickingLeaderState struct {
	// If picked leader is offline or unresponsive, this is incremented
	leaderOffset int64
}

type LeaderWaitingForCandidatesState struct {
	threshold int
	localPart types.SigningSessionIDPart
	parts     map[peer.ID]types.SigningSessionIDPart
}

type CandidateWaitingForLeaderState struct {
	localPart types.SigningSessionIDPart
}

type SigningState struct {
	transport mp_tss.Transporter

	outputChan chan tss_common.SignatureData
	errChan    chan *tss.Error
}

type DoneState struct {
	signature tss_common.SignatureData
}

type ErrorState struct {
	err *tss.Error
}

func NewPickingLeaderState() *PickingLeaderState {
	return &PickingLeaderState{
		leaderOffset: 0,
	}
}

func NewLeaderWaitingForCandidatesState() (*LeaderWaitingForCandidatesState, error) {
	localPart, err := types.NewSigningSessionIDPart()
	if err != nil {
		return nil, err
	}

	return &LeaderWaitingForCandidatesState{
		threshold: 0,
		localPart: localPart,
		parts:     make(map[peer.ID]types.SigningSessionIDPart),
	}, nil
}

func NewCandidateWaitingForLeaderState() (*CandidateWaitingForLeaderState, error) {
	localPart, err := types.NewSigningSessionIDPart()
	if err != nil {
		return nil, err
	}

	return &CandidateWaitingForLeaderState{
		localPart: localPart,
	}, nil
}

func NewSigningState(transport mp_tss.Transporter) *SigningState {
	return &SigningState{
		transport:  transport,
		outputChan: make(chan tss_common.SignatureData),
		errChan:    make(chan *tss.Error),
	}
}

func (s *PickingLeaderState) State() SigningSessionStateType {
	return SigningSessionStateType_PickingLeader
}

func (s *LeaderWaitingForCandidatesState) State() SigningSessionStateType {
	return SigningSessionStateType_LeaderWaitingForCandidates
}

func (s *CandidateWaitingForLeaderState) State() SigningSessionStateType {
	return SigningSessionStateType_CandidateWaitingForLeader
}

func (s *SigningState) State() SigningSessionStateType {
	return SigningSessionStateType_Signing
}

func (s *DoneState) State() SigningSessionStateType {
	return SigningSessionStateType_Done
}

func (s *ErrorState) State() SigningSessionStateType {
	return SigningSessionStateType_Error
}

// -----------------------------------------------------------------------------
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
	partyID       *tss.PartyID
	peerID        peer.ID
	sessionIDPart types.SigningSessionIDPart
}

type StartSignerEvent struct {
	tssParams *tss.Parameters
	key       keygen.LocalPartySaveData
	transport mp_tss.Transporter
}

type AddSigningPartEvent struct {
	From        *tss.PartyID
	Data        []byte
	IsBroadcast bool
}

// EventType
func (e *AddCandidateEvent) EventType() SigningSessionEventType {
	return SigningSessionEventType_AddCandidate
}

func (e *StartSignerEvent) EventType() SigningSessionEventType {
	return SigningSessionEventType_StartSigner
}

func (e *AddSigningPartEvent) EventType() SigningSessionEventType {
	return SigningSessionEventType_AddSigningPart
}

func NewAddCandidateEvent(
	partyID *tss.PartyID,
	peerID peer.ID,
	sessionIDPart types.SigningSessionIDPart,
) *AddCandidateEvent {
	return &AddCandidateEvent{
		partyID:       partyID,
		peerID:        peerID,
		sessionIDPart: sessionIDPart,
	}
}

func NewStartSignerEvent(
	tssParams *tss.Parameters,
	key keygen.LocalPartySaveData,
	transport mp_tss.Transporter,
) *StartSignerEvent {
	return &StartSignerEvent{
		tssParams: tssParams,
		key:       key,
		transport: transport,
	}
}

func NewAddSigningPartEvent(from peer.ID,
	data []byte,
	isBroadcast bool,
) *AddSigningPartEvent {
	return &AddSigningPartEvent{
		From:        from,
		Data:        data,
		IsBroadcast: isBroadcast,
	}
}
