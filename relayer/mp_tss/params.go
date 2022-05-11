package mp_tss

import (
	"github.com/binance-chain/tss-lib/tss"
)

// CreateParams creates tss parameters for the given party IDs, local partyID,
// and threshold for tss.
func CreateParams(
	partyIDs tss.UnSortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
) *tss.Parameters {
	parties := tss.SortPartyIDs(partyIDs)

	ctx := tss.NewPeerContext(parties)
	params := tss.NewParameters(tss.S256(), ctx, localPartyID, len(parties), threshold)

	return params
}
