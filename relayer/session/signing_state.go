package session

import (
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
	joinMsgs  types.JoinSessionMessages

	// Not actually used in this state, but passed to SigningState
	transport mp_tss.Transporter
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
		localPart: localPart,
		joinMsgs:  nil,
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

func NewDoneState(signature tss_common.SignatureData) *DoneState {
	return &DoneState{
		signature: signature,
	}
}

func NewErrorState(err *tss.Error) *ErrorState {
	return &ErrorState{
		err: err,
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
