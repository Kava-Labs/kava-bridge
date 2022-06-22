package testutil

import (
	"crypto/rand"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

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
		id, _ := peer.IDFromPrivateKey(key)

		out = append(out, id)
	}

	return out
}
