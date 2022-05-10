package mp_tss

import (
	"github.com/binance-chain/tss-lib/tss"
)

func CreateKeygenParams(
	partyIDs tss.UnSortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
) (*tss.Parameters, error) {
	// Create a `*PartyID` for each participating peer on the network
	// (you should call `tss.NewPartyID` for each one)
	parties := tss.SortPartyIDs(partyIDs)

	ctx := tss.NewPeerContext(parties)
	params := tss.NewParameters(tss.S256(), ctx, localPartyID, len(parties), threshold)

	return params, nil
}
