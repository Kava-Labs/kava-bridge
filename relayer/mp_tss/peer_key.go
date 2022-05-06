package mp_tss

import (
	"math/big"

	"github.com/libp2p/go-libp2p-core/peer"
)

func PeerIDToBigInt(peerID peer.ID) (*big.Int, error) {
	pubkey, err := peerID.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	pubkeyBytes, err := pubkey.Raw()
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(pubkeyBytes), nil
}
