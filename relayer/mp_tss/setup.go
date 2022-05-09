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
	// Set up elliptic curve
	// use ECDSA, which is used by default
	// TODO: Check if this is necessary since tss.NewParameters accepts a curve.
	// tss.SetCurve(s256k1.S256())

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

	// You should keep a local mapping of `id` strings to `*PartyID` instances
	// so that an incoming message can have its origin party's `*PartyID`
	// recovered for passing to `UpdateFromBytes` (see below)
	partyIDMap := make(map[string]*tss.PartyID)
	for _, id := range parties {
		partyIDMap[id.Id] = id
	}

	return SetupOutput{
		preParams: preParams,
		params:    params,
	}, err
}
