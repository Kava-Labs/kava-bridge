package mp_tss

import (
	"math/big"

	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/ecdsa/signing"
	"github.com/binance-chain/tss-lib/tss"
)

// RunSigner starts the local signing party and handles incoming and outgoing
// messages to other parties.
func RunSigner(
	msg *big.Int,
	params *tss.Parameters,
	key keygen.LocalPartySaveData,
	transport Transporter,
) (chan common.SignatureData, chan *tss.Error) {
	// outgoing messages to other peers - not one to not deadlock when a party
	// round is waiting for outgoing messages channel to clear
	outCh := make(chan tss.Message, params.PartyCount())
	// output signature when finished
	endCh := make(chan common.SignatureData, 1)
	// error if signing fails, contains culprits to blame
	errCh := make(chan *tss.Error, 1)

	log.Debugw("creating new local party")
	party := signing.NewLocalParty(msg, params, key, outCh, endCh)
	log.Debugw("local signing party created", "partyID", party.PartyID())

	RunParty(party, errCh, outCh, transport, false)

	return endCh, errCh
}
