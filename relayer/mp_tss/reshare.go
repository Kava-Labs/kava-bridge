package mp_tss

import (
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/ecdsa/resharing"
	"github.com/binance-chain/tss-lib/tss"
)

// RunReshare starts the local reshare party and handles incoming and outgoing
// messages to other parties.
func RunReshare(
	params *tss.ReSharingParameters,
	key keygen.LocalPartySaveData,
	transport Transporter,
) (chan keygen.LocalPartySaveData, chan *tss.Error) {
	// outgoing messages to other peers
	outCh := make(chan tss.Message, params.OldAndNewPartyCount())
	// output reshared key when finished
	endCh := make(chan keygen.LocalPartySaveData, 1)
	// error if reshare fails, contains culprits to blame
	errCh := make(chan *tss.Error, 1)

	log.Debugw("creating new local party")
	party := resharing.NewLocalParty(params, key, outCh, endCh)
	log.Debugw("local resharing party created", "partyID", party.PartyID())

	RunParty(party, errCh, outCh, transport, true)

	return endCh, errCh
}
