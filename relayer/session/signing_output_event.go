package session

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

type SigningSessionOutputEventType int

const (
	SigningSessionOutputEventType_StartLeaderSigner SigningSessionOutputEventType = iota + 1
)

func (e SigningSessionOutputEventType) String() string {
	switch e {
	case SigningSessionOutputEventType_StartLeaderSigner:
		return "StartLeaderSigner"
	default:
		return "Unknown"
	}
}

// SigningSessionOutputEvent is output events of a signing session to update
// the parent session store or to trigger another event from parent.
type SigningSessionOutputEvent interface {
	OutputEventType() SigningSessionOutputEventType
}

var _ SigningSessionOutputEvent = (*LeaderDoneOutputEvent)(nil)

// LeaderDoneOutputEvent is an output event when the leader is finished picking
// the participants and generating the aggregate session ID.
type LeaderDoneOutputEvent struct {
	TxHash                    common.Hash
	AggregateSigningSessionID types.AggregateSigningSessionID
	Participants              []peer.ID
}

func NewLeaderDoneOutputEvent(
	txHash common.Hash,
	aggregateSigningSessionID types.AggregateSigningSessionID,
	participants []peer.ID,
) *LeaderDoneOutputEvent {
	return &LeaderDoneOutputEvent{
		TxHash:                    txHash,
		AggregateSigningSessionID: aggregateSigningSessionID,
		Participants:              participants,
	}
}

// OutputEventType
func (e *LeaderDoneOutputEvent) OutputEventType() SigningSessionOutputEventType {
	return SigningSessionOutputEventType_StartLeaderSigner
}
