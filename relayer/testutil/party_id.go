package testutil

import (
	"fmt"
	"math/big"

	"github.com/binance-chain/tss-lib/tss"
	"github.com/libp2p/go-libp2p-core/peer"
)

func PartyIDsFromPeerIDs(peerIDs []peer.ID) tss.UnSortedPartyIDs {
	partyIDs := make(tss.UnSortedPartyIDs, len(peerIDs))
	for i, peerID := range peerIDs {
		pubkey, err := peerID.ExtractPublicKey()
		if err != nil {
			panic(fmt.Errorf("could not extract peer.ID pubkey: %w", err))
		}

		raw, err := pubkey.Raw()
		if err != nil {
			panic(fmt.Errorf("could not get raw pubkey bytes: %w", err))
		}

		key := new(big.Int).SetBytes(raw)

		pMoniker := fmt.Sprintf("%d", i+1)
		partyIDs[i] = tss.NewPartyID(pMoniker, pMoniker, key)
	}

	return partyIDs
}

// GetTestPartyIDs returns a list of party IDs derived from fixture libp2p
// publickeys for testing.
func GetTestPartyIDs(count int) tss.UnSortedPartyIDs {
	nodeKeys := GetTestP2pNodeKeys(count)
	peerIDs := PeerIDsFromKeys(nodeKeys)

	return PartyIDsFromPeerIDs(peerIDs)
}
