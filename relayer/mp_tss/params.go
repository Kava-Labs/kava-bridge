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
	return tss.NewParameters(tss.S256(), ctx, localPartyID, len(parties), threshold)
}

func CreateReShareParams(
	partyIDs tss.UnSortedPartyIDs,
	newPartyIDs tss.UnSortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
	newThreshold int,
) *tss.ReSharingParameters {
	oldParties := tss.SortPartyIDs(partyIDs)
	newParties := tss.SortPartyIDs(newPartyIDs)

	oldCtx := tss.NewPeerContext(oldParties)
	newCtx := tss.NewPeerContext(newParties)

	return tss.NewReSharingParameters(
		tss.S256(),      // curve
		oldCtx,          // Old PeerContext
		newCtx,          // New PeerContext with new peers
		localPartyID,    // Current party ID
		len(oldParties), // Current party count
		threshold,       // Current threshold
		len(newParties), // New party count
		newThreshold,    // New threshold
	)
}
