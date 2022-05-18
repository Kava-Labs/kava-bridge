package mp_tss_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/mp_tss"

	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/test"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	// err := logging.SetLogLevel("*", "debug")
	// require.NoError(t, err)

	// 1. Get party keys from file
	keys := GetTestKeys(threshold + 1)
	require.Len(t, keys, threshold+1)
	//
	// // Recreate party IDs from keys
	// signPIDs := GetTestPartyIDs(threshold + 1)
	// require.Len(t, signPIDs, threshold+1)

	signPIDsUnsorted := make(tss.UnSortedPartyIDs, len(keys))
	for i, key := range keys {
		pMoniker := fmt.Sprintf("%d", i+1)
		signPIDsUnsorted[i] = tss.NewPartyID(pMoniker, pMoniker, key.ShareID)
	}

	signPIDs := tss.SortPartyIDs(signPIDsUnsorted)
	// Sort keys so they match keys order
	sort.Slice(keys, func(i, j int) bool { return keys[i].ShareID.Cmp(keys[j].ShareID) == -1 })

	require.Equal(t, keygen.TestThreshold+1, len(keys))
	require.Equal(t, keygen.TestThreshold+1, len(signPIDs))

	// 2. Create and connect transport between peers
	transports := CreateAndConnectTransports(t, signPIDs)

	// 3. Start signing party for each peer
	outputAgg := make(chan common.SignatureData, keygen.TestThreshold)
	errAgg := make(chan *tss.Error, keygen.TestThreshold)

	for i := range signPIDs {
		params := mp_tss.CreateParams(signPIDs.ToUnSorted(), signPIDs[i], keygen.TestThreshold)
		t.Log(params.PartyID())

		// big.Int message, would be message hash converted to big int
		outputCh, errCh := mp_tss.RunSigner(big.NewInt(1234), params, keys[i], transports[i])

		go func(outputCh chan common.SignatureData, errCh chan *tss.Error) {
			for {
				select {
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
		case output := <-outputAgg:
			bz, err := json.Marshal(&output)
			require.NoError(t, err)
			t.Log(string(bz))

			signatures = append(signatures, output)
		case err := <-errAgg:
			t.Logf("err: %v", err)
		}
	}

	require.Len(t, signatures, test.TestThreshold+1, "each party should get a signature")

	for i, sig := range signatures {
		for j, sig2 := range signatures {
			// Skip self and previous keys
			if j <= i {
				continue
			}

			// make sure everyone has the same signature
			assert.True(t, bytes.Equal(sig.Signature, sig2.Signature))
		}
	}
}
