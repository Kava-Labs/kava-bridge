package mp_tss

import (
	"github.com/binance-chain/tss-lib/tss"
)

var Curve = tss.S256()

// CreateParams creates tss parameters for the given party IDs, local partyID,
// and threshold for tss.
func CreateParams(
	partyIDs tss.UnSortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
) *tss.Parameters {
	parties := tss.SortPartyIDs(partyIDs)

	ctx := tss.NewPeerContext(parties)
	return tss.NewParameters(
		Curve,                                    // secp256k1 curve
		ctx,                                      // PeerContext
		parties.FindByKey(localPartyID.KeyInt()), // Current party ID **after** sorted for correct Index
		len(parties),                             // Party count
		threshold,                                // Threshold
	)
}

func CreateReShareParams(
	oldPartyIDs tss.UnSortedPartyIDs,
	newPartyIDs tss.UnSortedPartyIDs,
	localPartyID *tss.PartyID,
	threshold int,
	newThreshold int,
) *tss.ReSharingParameters {
	oldCtx := tss.NewPeerContext(tss.SortPartyIDs(oldPartyIDs))
	newCtx := tss.NewPeerContext(tss.SortPartyIDs(newPartyIDs))

	return tss.NewReSharingParameters(
		Curve,            // secp256k1 curve
		oldCtx,           // Old PeerContext
		newCtx,           // New PeerContext with new peers
		localPartyID,     // Current party ID
		len(oldPartyIDs), // Current party count
		threshold,        // Current threshold
		len(newPartyIDs), // New party count
		newThreshold,     // New threshold
	)
}
