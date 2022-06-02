package mp_tss

import "github.com/binance-chain/tss-lib/tss"

// Transporter is the interface that defines the Send and Receive methods to
// transfer lib-tss messages between parties.
type Transporter interface {
	Send([]byte, *tss.MessageRouting, bool) error
	// Receive returns a channel that will be read by the local tss party. This
	// consists of ReceivedPartyState messages received from other parties.
	Receive() chan ReceivedPartyState
}

// ReceivedPartyState is a message received from another party
type ReceivedPartyState struct {
	WireBytes   []byte
	From        *tss.PartyID
	IsBroadcast bool
}

// NewReceivedPartyState returns a new ReceivedPartyState.
func NewReceivedPartyState(
	wireBytes []byte,
	from *tss.PartyID,
	isBroadcast bool,
) ReceivedPartyState {
	return ReceivedPartyState{
		WireBytes:   wireBytes,
		From:        from,
		IsBroadcast: isBroadcast,
	}
}
