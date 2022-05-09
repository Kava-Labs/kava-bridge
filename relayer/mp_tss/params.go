package mp_tss

import (
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
)

func CreateKeygenParams(
	partyIDs tss.UnSortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
) (*keygen.LocalPreParams, *tss.Parameters, error) {
	// When using the keygen party it is recommended that you pre-compute the
	// "safe primes" and Paillier secret beforehand because this can take some time.
	// This code will generate those parameters using a concurrency limit equal
	// to the number of available CPU cores.
	preParams, err := keygen.GeneratePreParams(1 * time.Minute)
	if err != nil {
		return nil, nil, err
	}

	// Create a `*PartyID` for each participating peer on the network
	// (you should call `tss.NewPartyID` for each one)
	parties := tss.SortPartyIDs(partyIDs)

	ctx := tss.NewPeerContext(parties)
	params := tss.NewParameters(tss.S256(), ctx, localPartyID, len(parties), threshold)

	return preParams, params, nil
}
