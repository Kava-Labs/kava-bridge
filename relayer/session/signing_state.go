package session

import (
	"sync"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
)

type SigningSessionStateType int

const (
	SigningSessionStateType_PickingLeader SigningSessionStateType = iota + 1
	SigningSessionStateType_LeaderWaitingForCandidates
	SigningSessionStateType_CandidateWaitingForLeader
	SigningSessionStateType_Signing
	SigningSessionStateType_Done
	SigningSessionStateType_DoneNonParticipant
	SigningSessionStateType_Error
)

// String returns the string representation of the state.
func (t SigningSessionStateType) String() string {
	switch t {
	case SigningSessionStateType_PickingLeader:
		return "PickingLeader"
	case SigningSessionStateType_LeaderWaitingForCandidates:
		return "LeaderWaitingForCandidates"
	case SigningSessionStateType_CandidateWaitingForLeader:
		return "CandidateWaitingForLeader"
	case SigningSessionStateType_Signing:
		return "Signing"
	case SigningSessionStateType_Done:
		return "Done"
	case SigningSessionStateType_DoneNonParticipant:
		return "DoneNonParticipant"
	case SigningSessionStateType_Error:
		return "Error"
	default:
		return "Unknown"
	}
}

// SigningSessionState is the state of a signing session.
type SigningSessionState interface {
	State() SigningSessionStateType
}

var _ SigningSessionState = (*PickingLeaderState)(nil)
var _ SigningSessionState = (*LeaderWaitingForCandidatesState)(nil)
var _ SigningSessionState = (*CandidateWaitingForLeaderState)(nil)
var _ SigningSessionState = (*SigningState)(nil)
var _ SigningSessionState = (*DoneState)(nil)
var _ SigningSessionState = (*DoneNonParticipantState)(nil)
var _ SigningSessionState = (*ErrorState)(nil)

// PickingLeaderState is the state of a signing session where each peer is
// picking the leader.
type PickingLeaderState struct {
	// If picked leader is offline or unresponsive, this is incremented
	leaderOffset int64
}

// LeaderWaitingForCandidatesState is the state of a signing session where the
// leader is waiting for candidates to join.
type LeaderWaitingForCandidatesState struct {
	// Leader's local part of the signing session ID
	localPart types.SigningSessionIDPart
	// Join messages received from other parties

	joinMsgsLock *sync.Mutex
	joinMsgs     types.JoinSessionMessages
}

// CandidateWaitingForLeaderState is the state of a signing session where the
// non-leaders wait for the leader to start the signing party.
type CandidateWaitingForLeaderState struct {
	// Local part of the signing session ID
	localPart types.SigningSessionIDPart
}

// SigningState is the state of a signing session where all participants are signing.
type SigningState struct {
	transport mp_tss.Transporter

	outputChan chan tss_common.SignatureData
	errChan    chan *tss.Error
}

// DoneState is the state when the signing session is done and a signature is output.
type DoneState struct {
	signature tss_common.SignatureData
}

// DoneNonParticipantState is the state when the peer is not a participant and
// thus has no signature.
type DoneNonParticipantState struct {
}

// ErrorState is the state when an error during signing occurs.
type ErrorState struct {
	err *tss.Error
}

// -----------------------------------------------------------------------------
// New State methods

// NewPickingLeaderState returns a new PickingLeaderState.
func NewPickingLeaderState() *PickingLeaderState {
	return &PickingLeaderState{
		leaderOffset: 0,
	}
}

// NewLeaderWaitingForCandidatesState returns a new LeaderWaitingForCandidatesState.
func NewLeaderWaitingForCandidatesState() (*LeaderWaitingForCandidatesState, error) {
	localPart, err := types.NewSigningSessionIDPart()
	if err != nil {
		return nil, err
	}

	return &LeaderWaitingForCandidatesState{
		localPart:    localPart,
		joinMsgsLock: &sync.Mutex{},
		joinMsgs:     nil,
	}, nil
}

// NewCandidateWaitingForLeaderState returns a new CandidateWaitingForLeaderState.
func NewCandidateWaitingForLeaderState() (*CandidateWaitingForLeaderState, error) {
	localPart, err := types.NewSigningSessionIDPart()
	if err != nil {
		return nil, err
	}

	return &CandidateWaitingForLeaderState{
		localPart: localPart,
	}, nil
}

// NewSigningState returns a new SigningState.
func NewSigningState(transport mp_tss.Transporter) *SigningState {
	return &SigningState{
		transport:  transport,
		outputChan: make(chan tss_common.SignatureData),
		errChan:    make(chan *tss.Error),
	}
}

//nolint:govet
// NewDoneState returns a new DoneState.
func NewDoneState(signature tss_common.SignatureData) *DoneState {
	return &DoneState{
		//nolint:govet
		signature: signature,
	}
}

// NewDoneNonParticipantState returns a new DoneNonParticipantState.
func NewDoneNonParticipantState() *DoneNonParticipantState {
	return &DoneNonParticipantState{}
}

// NewErrorState returns a new ErrorState.
func NewErrorState(err *tss.Error) *ErrorState {
	return &ErrorState{
		err: err,
	}
}

// -----------------------------------------------------------------------------
// SigningSessionState interface implementations

// State returns the state type of the session state.
func (s *PickingLeaderState) State() SigningSessionStateType {
	return SigningSessionStateType_PickingLeader
}

// State returns the state type of the session state.
func (s *LeaderWaitingForCandidatesState) State() SigningSessionStateType {
	return SigningSessionStateType_LeaderWaitingForCandidates
}

// State returns the state type of the session state.
func (s *CandidateWaitingForLeaderState) State() SigningSessionStateType {
	return SigningSessionStateType_CandidateWaitingForLeader
}

// State returns the state type of the session state.
func (s *SigningState) State() SigningSessionStateType {
	return SigningSessionStateType_Signing
}

// State returns the state type of the session state.
func (s *DoneState) State() SigningSessionStateType {
	return SigningSessionStateType_Done
}

// State returns the state type of the session state.
func (s *DoneNonParticipantState) State() SigningSessionStateType {
	return SigningSessionStateType_DoneNonParticipant
}

// State returns the state type of the session state.
func (s *ErrorState) State() SigningSessionStateType {
	return SigningSessionStateType_Error
}
