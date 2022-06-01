package mp_tss

import (
	"fmt"
	"math/big"

	"github.com/binance-chain/tss-lib/tss"
	"github.com/libp2p/go-libp2p-core/peer"
)

// PartyIDStore keeps track of the party IDs of all libp2p peers.
type PartyIDStore struct {
	partyIDs map[peer.ID]*tss.PartyID
}

// NewPartyIDStore creates a new PartyIDStore.
func NewPartyIDStore() *PartyIDStore {
	return &PartyIDStore{
		partyIDs: make(map[peer.ID]*tss.PartyID),
	}
}

// AddPeer adds a peer.ID and moniker to the store.
func (s *PartyIDStore) AddPeer(peerID peer.ID, moniker string) error {
	pubkey, err := peerID.ExtractPublicKey()
	if err != nil {
		return fmt.Errorf("could not extract peer.ID pubkey: %w", err)
	}

	raw, err := pubkey.Raw()
	if err != nil {
		return fmt.Errorf("could not get raw pubkey bytes: %w", err)
	}

	if len(raw) != 32 {
		return fmt.Errorf("pubkey raw bytes are not 32 bytes: %d", len(raw))
	}

	key := new(big.Int).SetBytes(raw)
	s.partyIDs[peerID] = tss.NewPartyID(peerID.String(), moniker, key)

	return nil
}

func (s *PartyIDStore) GetPartyID(peerID peer.ID) (*tss.PartyID, error) {
	if err := peerID.Validate(); err != nil {
		return nil, fmt.Errorf("invalid PeerID: %w", err)
	}

	if _, ok := s.partyIDs[peerID]; !ok {
		return nil, fmt.Errorf("peerID %s not found", peerID)
	}

	return s.partyIDs[peerID], nil
}
