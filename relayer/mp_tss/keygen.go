package mp_tss

import (
	logging "github.com/ipfs/go-log/v2"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
)

var log = logging.Logger("mp_tss")

// RunKeyGen starts the local keygen party and handles incoming and outgoing
// messages to other parties.
func RunKeyGen(
	preParams *keygen.LocalPreParams,
	params *tss.Parameters,
	transport Transporter,
) (chan keygen.LocalPartySaveData, chan *tss.Error) {
	// outgoing messages to other peers
	outCh := make(chan tss.Message, 10)
	// error if keygen fails, contains culprits to blame
	errCh := make(chan *tss.Error, 10)
	// output data when keygen finished
	endCh := make(chan keygen.LocalPartySaveData, 10)

	log.Debugw("creating new local party")
	party := keygen.NewLocalParty(params, outCh, endCh, *preParams)
	log.Debugw("local party created", "partyID", party.PartyID())

	RunParty(party, errCh, outCh, transport, false)

	return endCh, errCh
}
