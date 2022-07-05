package session

import (
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/libp2p/go-libp2p-core/peer"
)

// SigningSessionEvent is implemented by all events that are part of a signing
// session.
type SigningSessionEvent interface {
	SigningSessionEvent()
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
func (e *AddCandidateEvent) SigningSessionEvent()   {}
func (e *StartSignerEvent) SigningSessionEvent()    {}
func (e *AddSigningPartEvent) SigningSessionEvent() {}

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
