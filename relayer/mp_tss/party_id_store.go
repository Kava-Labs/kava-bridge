package mp_tss

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/binance-chain/tss-lib/tss"
	"github.com/libp2p/go-libp2p-core/peer"
)

// PartyIDStore keeps track of the party IDs of all libp2p peers.
type PartyIDStore struct {
	partyIDs map[peer.ID]*tss.PartyID // peer.ID -> PartyID
	peerIDs  map[string]peer.ID       // PartyID -> peer.ID
}

// NewPartyIDStore creates a new PartyIDStore.
func NewPartyIDStore() *PartyIDStore {
	return &PartyIDStore{
		partyIDs: make(map[peer.ID]*tss.PartyID),
		peerIDs:  make(map[string]peer.ID),
	}
}

// AddPeers adds a list of peer.IDs and monikers to the store.
func (s *PartyIDStore) AddPeers(
	peers []*PeerMetadata,
) error {
	var partyIDs tss.UnSortedPartyIDs
	// PartyID.Key -> PeerMetadata
	partyIDMap := make(map[string]*PeerMetadata)

	for _, peer := range peers {
		partyID, err := peer.GetUnsortedPartyID()
		if err != nil {
			return err
		}
		partyIDs = append(partyIDs, partyID)
		partyIDMap[hex.EncodeToString(partyID.Key)] = peer
	}

	sortedPartyIDs := tss.SortPartyIDs(partyIDs)

	for _, partyID := range sortedPartyIDs {
		peer := partyIDMap[hex.EncodeToString(partyID.Key)]
		if err := peer.SetSortedPartyID(partyID); err != nil {
			return err
		}

		// Set store with sorted partyID
		s.partyIDs[peer.PeerID] = partyID
		s.peerIDs[hex.EncodeToString(partyID.Key)] = peer.PeerID
	}

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

func (s *PartyIDStore) GetPeerID(partyID *tss.PartyID) (peer.ID, bool) {
	peerID, ok := s.peerIDs[hex.EncodeToString(partyID.Key)]
	return peerID, ok
}

func (s *PartyIDStore) String() string {
	str := "PartyIDStore{\npartyIDs:\n"

	for peerID, partyID := range s.partyIDs {
		str += fmt.Sprintf(
			" - %v -> {Index: %v, Id: %v, Moniker: %v, Key: %v}\n",
			peerID.ShortString(),
			partyID.Index,
			partyID.Id,
			partyID.Moniker,
			hex.EncodeToString(partyID.Key),
		)
	}

	str += "peerIDs:\n"

	for partyIDKey, peerID := range s.peerIDs {
		str += fmt.Sprintf(" - %v -> %v\n", partyIDKey, peerID.ShortString())
	}

	str += "}"

	return str
}

// 1. GetPartyID() from all parties
// 2. Sort party IDs
// 3. Update PartyID in each PeerMetadata to have the correct index
type PeerMetadata struct {
	PeerID  peer.ID
	Moniker string

	// partyIDCache is a cached partyID, set when GetPartyID is first called.
	partyIDCache *tss.PartyID
}

func NewPeerMetadata(peerID peer.ID, moniker string) *PeerMetadata {
	return &PeerMetadata{
		PeerID:  peerID,
		Moniker: moniker,
	}
}

// GetUnsortedPartyID get's the partyID of a with the partyID key set to the
// peerID's pubkey.
func (pm *PeerMetadata) GetUnsortedPartyID() (*tss.PartyID, error) {
	if pm.partyIDCache != nil {
		return pm.partyIDCache, nil
	}

	pubkey, err := pm.PeerID.ExtractPublicKey()
	if err != nil {
		return nil, fmt.Errorf("could not extract peer.ID pubkey: %w", err)
	}

	raw, err := pubkey.Raw()
	if err != nil {
		return nil, fmt.Errorf("could not get raw pubkey bytes: %w", err)
	}

	key := new(big.Int).SetBytes(raw)

	pm.partyIDCache = tss.NewPartyID(pm.Moniker, pm.Moniker, key)
	return pm.partyIDCache, nil
}

// SetSortedPartyID sets the partyID of a with the partyID key set to the
func (pm *PeerMetadata) SetSortedPartyID(partyID *tss.PartyID) error {
	if partyID.KeyInt().Cmp(pm.partyIDCache.KeyInt()) != 0 {
		return fmt.Errorf("partyID key does not match peerID pubkey")
	}

	if partyID.Index < 0 {
		return fmt.Errorf("partyID index is negative, must be sorted")
	}

	pm.partyIDCache = partyID

	return nil
}
