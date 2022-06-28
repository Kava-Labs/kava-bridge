package session

import (
	"sync"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
)

// SigningSessionState is the state of a signing session.
type SigningSessionState interface {
	SigningSessionState()
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

func (s *PickingLeaderState) SigningSessionState()              {}
func (s *LeaderWaitingForCandidatesState) SigningSessionState() {}
func (s *CandidateWaitingForLeaderState) SigningSessionState()  {}
func (s *SigningState) SigningSessionState()                    {}
func (s *DoneState) SigningSessionState()                       {}
func (s *DoneNonParticipantState) SigningSessionState()         {}
func (s *ErrorState) SigningSessionState()                      {}
