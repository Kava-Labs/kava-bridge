package mp_tss

import (
	"github.com/binance-chain/tss-lib/tss"
)

var curve = tss.S256()

// CreateParams creates tss parameters for the given party IDs, local partyID,
// and threshold for tss.
func CreateParams(
	partyIDs tss.UnSortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
) *tss.Parameters {
	parties := tss.SortPartyIDs(partyIDs)

	ctx := tss.NewPeerContext(parties)
	return tss.NewParameters(curve, ctx, localPartyID, len(parties), threshold)
}

func CreateReShareParams(
	oldPartyIDs tss.SortedPartyIDs,
	newPartyIDs tss.SortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
	newThreshold int,
) *tss.ReSharingParameters {
	oldCtx := tss.NewPeerContext(oldPartyIDs)
	newCtx := tss.NewPeerContext(newPartyIDs)

	return tss.NewReSharingParameters(
		curve,             // secp256k1 curve
		oldCtx,            // Old PeerContext
		newCtx,            // New PeerContext with new peers
		localPartyID,      // Current party ID
		oldPartyIDs.Len(), // Current party count
		threshold,         // Current threshold
		newPartyIDs.Len(), // New party count
		newThreshold,      // New threshold
	)
}
