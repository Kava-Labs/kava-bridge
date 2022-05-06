package mp_tss

import "github.com/binance-chain/tss-lib/tss"

type Transporter interface {
	Send([]byte, *tss.MessageRouting) error
	Receive() <-chan ReceivedPartyState
}

// ReceivedPartyState is a message received from another party
type ReceivedPartyState struct {
	wireBytes   []byte
	from        *tss.PartyID
	isBroadcast bool
}
