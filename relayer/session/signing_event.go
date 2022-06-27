package session

import (
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

type SigningSessionEventType int

const (
	SigningSessionEventType_AddCandidate SigningSessionEventType = iota + 1
	SigningSessionEventType_StartSigner
	SigningSessionEventType_AddSigningPart
)

func (e SigningSessionEventType) String() string {
	switch e {
	case SigningSessionEventType_AddCandidate:
		return "AddCandidate"
	case SigningSessionEventType_StartSigner:
		return "StartSigner"
	case SigningSessionEventType_AddSigningPart:
		return "AddSigningPart"
	default:
		return "Unknown"
	}
}

type SigningSessionEvent interface {
	EventType() SigningSessionEventType
}

var _ SigningSessionEvent = (*AddCandidateEvent)(nil)
var _ SigningSessionEvent = (*StartSignerEvent)(nil)
var _ SigningSessionEvent = (*AddSigningPartEvent)(nil)

type AddCandidateEvent struct {
	partyID *tss.PartyID

	joinMsg types.JoinSessionMessage
}

type StartSignerEvent struct {
	tssParams    *tss.Parameters
	transport    mp_tss.Transporter
	participants []peer.ID
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
	joinMsg types.JoinSessionMessage,
) *AddCandidateEvent {
	return &AddCandidateEvent{
		partyID: partyID,
		joinMsg: joinMsg,
	}
}

func NewStartSignerEvent(
	tssParams *tss.Parameters,
	transport mp_tss.Transporter,
	participants []peer.ID,
) *StartSignerEvent {
	return &StartSignerEvent{
		tssParams:    tssParams,
		transport:    transport,
		participants: participants,
	}
}

func NewAddSigningPartEvent(
	from *tss.PartyID,
	data []byte,
	isBroadcast bool,
) *AddSigningPartEvent {
	return &AddSigningPartEvent{
		From:        from,
		Data:        data,
		IsBroadcast: isBroadcast,
	}
}
