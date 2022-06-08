package testutil

import (
	"crypto/rand"
	"fmt"
	"path"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
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
