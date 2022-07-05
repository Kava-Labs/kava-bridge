package mp_tss_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/testutil"

	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	// err := logging.SetLogLevel("*", "debug")
	// require.NoError(t, err)

	_, _, keys, signPIDs := testutil.GetTestKeys(t, testutil.TestThreshold+1)

	// 2. Create and connect transport between peers
	transports := CreateAndConnectTransports(t, signPIDs)
	require.Len(t, transports, testutil.TestThreshold+1)

	// 3. Start signing party for each peer
	outputAgg := make(chan common.SignatureData, testutil.TestThreshold)
	errAgg := make(chan *tss.Error, testutil.TestThreshold)

	msgHash := big.NewInt(1234)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := range signPIDs {
		params := mp_tss.CreateParams(signPIDs, signPIDs[i], testutil.TestThreshold)
		t.Log(params.PartyID())

		// big.Int message, would be message hash converted to big int
		outputCh, errCh := mp_tss.RunSign(ctx, msgHash, params, keys[i], transports[i])

		go func(outputCh chan common.SignatureData, errCh chan *tss.Error) {
			for {
				select {
				//nolint:govet // https://github.com/bnb-chain/tss-lib/pull/167
				case output := <-outputCh:
					outputAgg <- output
				case err := <-errCh:
					errAgg <- err
				}
			}
		}(outputCh, errCh)
	}

	t.Logf("started signing")

	var signatures []common.SignatureData

	for range signPIDs {
		select {
		//nolint:govet
		case output := <-outputAgg:
			bz, err := json.Marshal(&output)
			require.NoError(t, err)
			t.Log(string(bz))

			//nolint:govet
			signatures = append(signatures, output)
		case err := <-errAgg:
			t.Logf("err: %v", err)
		}
	}

	require.Len(t, signatures, testutil.TestThreshold+1, "each party should get a signature")

	//nolint:govet
	for i, sig := range signatures {
		//nolint:govet
		for j, sig2 := range signatures {
			// Skip self and previous keys
			if j <= i {
				continue
			}

			// make sure everyone has the same signature
			assert.True(t, bytes.Equal(sig.Signature, sig2.Signature))
		}
	}

	// Verify signature
	pkX, pkY := keys[0].ECDSAPub.X(), keys[0].ECDSAPub.Y()
	pk := ecdsa.PublicKey{
		Curve: mp_tss.Curve,
		X:     pkX,
		Y:     pkY,
	}

	ok := ecdsa.Verify(
		&pk,                                    // pubkey
		msgHash.Bytes(),                        // message hash
		new(big.Int).SetBytes(signatures[0].R), // R
		new(big.Int).SetBytes(signatures[0].S), // S
	)
	assert.True(t, ok, "ecdsa verify must pass")

	t.Log("signature verified")
}
