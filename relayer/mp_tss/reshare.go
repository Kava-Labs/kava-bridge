package mp_tss

import (
	"context"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/ecdsa/resharing"
	"github.com/binance-chain/tss-lib/tss"
)

// RunReshare starts the local reshare party and handles incoming and outgoing
// messages to other parties.
func RunReshare(
	ctx context.Context,
	params *tss.ReSharingParameters,
	key keygen.LocalPartySaveData,
	transport Transporter,
) (chan keygen.LocalPartySaveData, chan *tss.Error) {
	// outgoing messages to other peers
	outCh := make(chan tss.Message, 1)
	// output reshared key when finished
	endCh := make(chan keygen.LocalPartySaveData, 1)
	// error if reshare fails, contains culprits to blame
	errCh := make(chan *tss.Error, 1)

	log.Debugw("creating new local party")
	party := resharing.NewLocalParty(params, key, outCh, endCh)
	log.Debugw("local resharing party created", "partyID", party.PartyID())

	RunParty(ctx, party, errCh, outCh, transport, true)

	return endCh, errCh
}
