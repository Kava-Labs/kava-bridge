package mp_tss_test

import (
	"context"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/assert"
)

func TestKeygen(t *testing.T) {
	// err := logging.SetLogLevel("*", "debug")
	// require.NoError(t, err)

	// 1. Create party ID for each peer, share with other peers
	partyIDs := testutil.GetTestPartyIDs(testutil.TestPartyCount)

	// 2. Create and connect transport between peers
	transports := CreateAndConnectTransports(t, partyIDs)

	// 3. Make params and start peers
	errAgg := make(chan *tss.Error)
	outputAgg := make(chan keygen.LocalPartySaveData)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := range partyIDs {
		// Load from disk to avoid re-generating
		preParams := LoadTestPreParam(i)
		params := mp_tss.CreateParams(partyIDs, partyIDs[i], testutil.TestThreshold)

		outputCh, errCh := mp_tss.RunKeyGen(ctx, preParams, params, transports[i])
		go func(outputCh chan keygen.LocalPartySaveData, errCh chan *tss.Error) {
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

	t.Logf("started keygen")

	var keys []keygen.LocalPartySaveData

	// 4. Wait for all parties to finish.
	for range partyIDs {
		select {
		case output := <-outputAgg:
			keys = append(keys, output)
		case err := <-errAgg:
			t.Fatal(err)
		}
	}

	// make sure everyone has the same ECDSA public key
	for i, key := range keys {
		assert.True(t, key.Validate(), "key should be valid")
		assert.True(t, key.ValidateWithProof(), "key should be valid with proof")

		for j, key2 := range keys {
			// Skip self and previous keys
			if j <= i {
				continue
			}

			assert.Truef(t, key.ECDSAPub.Equals(key2.ECDSAPub), "key %v != %v", i, j)
		}
	}

	// // Write keys to file for test fixtures for signing
	// // Must be in the same order as PartyIDs
	// for i, partyID := range partyIDs {
	// 	// Search key for this partyID
	// 	for _, key := range keys {
	// 		if key.ShareID.Cmp(partyID.KeyInt()) == 0 {
	// 			assert.Equal(t, partyIDs[i].KeyInt(), key.ShareID, "saved key part should match party id")
	//
	// 			fmt.Printf("partyID = %v \n", partyID.KeyInt())
	// 			fmt.Printf("keyID   = %v \n", key.ShareID)
	//
	// 			testutil.WriteTestKey(i, key)
	// 			break
	// 		}
	// 	}
	// }
}
