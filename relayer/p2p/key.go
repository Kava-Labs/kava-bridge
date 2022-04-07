package p2p

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/libp2p/go-libp2p-core/crypto"
	crypto_pb "github.com/libp2p/go-libp2p-core/crypto/pb"
)

const (
	PreSharedNetworkKeyLengthBytes = 32
)

// UnmarshalPrivateKey unmarshals a private key from bytes
func UnmarshalPrivateKey(data []byte) (crypto.PrivKey, error) {
	privKey, err := crypto.UnmarshalPrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode private key: %w", err)
	}

	if privKey.Type() != crypto_pb.KeyType_Secp256k1 {
		return nil, fmt.Errorf("invalid key type %s", privKey.Type().String())
	}

	rawKey, err := privKey.Raw()
	if err != nil {
		return nil, fmt.Errorf("error decoding private key: %w", err)
	}

	if len(rawKey) != btcec.PrivKeyBytesLen {
		return nil, fmt.Errorf(
			"invalid private key length %d, expected %d",
			len(rawKey),
			btcec.PrivKeyBytesLen,
		)
	}

	return privKey, nil
}
