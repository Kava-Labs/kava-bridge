package signer_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/kava-labs/kava-bridge/relayer/signer"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestSigner(t *testing.T) {
	numPeers := 5
	threshold := 3

	ctx := context.Background()
	done := make(chan bool)

	node_keys, err := testutil.GenerateNodeKeys(numPeers)
	require.NoError(t, err)

	peerIDs := testutil.PeerIDsFromKeys(node_keys)

	tss_keys, _ := testutil.GetTestTssKeys(numPeers)
	require.Len(t, tss_keys, numPeers)

	signers := make([]*signer.Signer, numPeers)
	for i := 0; i < numPeers; i++ {
		opts := p2p.NodeOptions{
			Port:              0,
			NetworkPrivateKey: []byte("network-private-key"),
			NodePrivateKey:    node_keys[i],
			PeerList:          peerIDs,
		}

		node, err := p2p.NewNode(ctx, opts, done)
		require.NoError(t, err)

		signers[i] = signer.NewSigner(
			node,
			fmt.Sprintf("node-%v", i),
			tss_keys[i],
			threshold,
		)
	}

	for _, signer := range signers {
		require.NoError(t, signer.Start())
	}

	txHash := common.BigToHash(big.NewInt(1))
	msgHash := big.NewInt(2)

	g := new(errgroup.Group)
	sigs := make([]tss_common.SignatureData, numPeers)

	for _, s := range signers {
		func(signer *signer.Signer) {
			g.Go(func() error {
				// The relayer will call this when there is a new signing output from
				// block syncing.
				sig, err := signer.SignMessage(txHash, msgHash)
				if err != nil {
					return err
				}

				sigs = append(sigs, sig)

				return nil
			})
		}(s)
	}

	err = g.Wait()
	require.NoError(t, err)

	for _, sig := range sigs {
		require.NotNil(t, sig)
	}

	// Verify signature
	pkX, pkY := tss_keys[0].ECDSAPub.X(), tss_keys[0].ECDSAPub.Y()
	pk := ecdsa.PublicKey{
		Curve: mp_tss.Curve,
		X:     pkX,
		Y:     pkY,
	}

	ok := ecdsa.Verify(
		&pk,                              // pubkey
		msgHash.Bytes(),                  // message hash
		new(big.Int).SetBytes(sigs[0].R), // R
		new(big.Int).SetBytes(sigs[0].S), // S
	)
	assert.True(t, ok, "ecdsa verify must pass")
}
