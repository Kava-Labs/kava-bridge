package mp_tss

import (
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
)

type SetupOptions struct {
	PartyIDs tss.UnSortedPartyIDs
	PartyID  *tss.PartyID

	Threshold int
}

type SetupOutput struct {
	preParams *keygen.LocalPreParams
	params    *tss.Parameters
}

func CreateKeygenParams(options SetupOptions) (SetupOutput, error) {
	// When using the keygen party it is recommended that you pre-compute the
	// "safe primes" and Paillier secret beforehand because this can take some time.
	// This code will generate those parameters using a concurrency limit equal
	// to the number of available CPU cores.
	preParams, err := keygen.GeneratePreParams(1 * time.Minute)
	if err != nil {
		return SetupOutput{}, err
	}

	// Create a `*PartyID` for each participating peer on the network
	// (you should call `tss.NewPartyID` for each one)
	parties := tss.SortPartyIDs(options.PartyIDs)

	ctx := tss.NewPeerContext(parties)
	params := tss.NewParameters(tss.S256(), ctx, options.PartyID, len(parties), options.Threshold)

	return SetupOutput{
		preParams: preParams,
		params:    params,
	}, err
}
