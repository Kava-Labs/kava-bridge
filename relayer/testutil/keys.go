package testutil

import (
	"crypto/rand"
	"fmt"
	"path"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/os"
)

func p2pNodeTestKeyPath(index int) string {
	return path.Join("..", "testutil", "test-fixtures", fmt.Sprintf("libp2p-%02d.key", index))
}

func ReadP2pNodeTestKey(index int) crypto.PrivKey {
	path := p2pNodeTestKeyPath(index)

	bytes := os.MustReadFile(path)
	key, err := crypto.UnmarshalPrivateKey(bytes)

	if err != nil {
		panic(err)
	}

	return key
}

func GetTestP2pNodeKeys(count int) []crypto.PrivKey {
	var keys []crypto.PrivKey
	for i := 0; i < count; i++ {
		key := ReadP2pNodeTestKey(i)
		keys = append(keys, key)
	}

	return keys
}

func WriteP2pNodeTestKey(index int, key crypto.PrivKey) {
	bz, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		panic(err)
	}

	os.MustWriteFile(p2pNodeTestKeyPath(index), bz, 0600)
}

// GenerateNodeKeys returns a slice of private keys.
func GenerateNodeKeys(n int) ([]crypto.PrivKey, error) {
	var keys []crypto.PrivKey

	for i := 0; i < n; i++ {
		priv, _, err := crypto.GenerateKeyPairWithReader(
			crypto.Secp256k1,
			0, // bitsize ignored for secp256k1
			rand.Reader,
		)

		if err != nil {
			return nil, err
		}

		keys = append(keys, priv)
	}

	return keys, nil
}

// PeerIDsFromKeys returns a slice of peer IDs from the given slice of private keys.
func PeerIDsFromKeys(keys []crypto.PrivKey) []peer.ID {
	var out []peer.ID

	for _, key := range keys {
		id, err := peer.IDFromPrivateKey(key)
		if err != nil {
			panic(err)
		}

		out = append(out, id)
	}

	return out
}

func GetTestKeys(t *testing.T, numPeers int) (
	[]crypto.PrivKey,
	peer.IDSlice,
	[]keygen.LocalPartySaveData,
	tss.UnSortedPartyIDs,
) {
	nodeKeys := GetTestP2pNodeKeys(numPeers)
	require.Len(t, nodeKeys, numPeers)

	// Peer ID derived from private libp2p key
	peerIDs := PeerIDsFromKeys(nodeKeys)
	require.Len(t, peerIDs, numPeers)

	// Party ID derived from peer ID public key
	partyIDs := PartyIDsFromPeerIDs(peerIDs)
	require.Len(t, partyIDs, numPeers)

	tss_keys := GetTestTssKeys(numPeers)
	require.Len(t, tss_keys, numPeers)

	for i := range tss_keys {
		t.Logf("tss_keys[%d] = %+v", i, tss_keys[i].ShareID)
		t.Logf("partyIDs[%d] = %+v", i, partyIDs[i].KeyInt())
	}

	return nodeKeys, peerIDs, tss_keys, partyIDs
}
