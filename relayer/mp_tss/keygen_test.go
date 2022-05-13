package mp_tss_test

import (
	"encoding/json"
	"testing"

	logging "github.com/ipfs/go-log/v2"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const partyCount = 10
const threshold = 8

func TestKeygen(t *testing.T) {
	err := logging.SetLogLevel("*", "debug")
	require.NoError(t, err)

	// 1. Create party ID for each peer, share with other peers
	partyIDs := tss.GenerateTestPartyIDs(partyCount)

	// 2. Create and connect transport between peers
	transports := CreateAndConnectTransports(t, partyIDs)

	// 3. Make params and start peers
	errAgg := make(chan *tss.Error)
	outputAgg := make(chan keygen.LocalPartySaveData)

	for i := range partyIDs {
		// Load from disk to avoid re-generating
		preParams := LoadTestPreParam(i)
		params := mp_tss.CreateParams(partyIDs.ToUnSorted(), partyIDs[i], threshold)

		outputCh, errCh := mp_tss.RunKeyGen(preParams, params, transports[i])
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
		for j, key2 := range keys {
			// Skip self and previous keys
			if j <= i {
				continue
			}

			assert.Truef(t, key.ECDSAPub.Equals(key2.ECDSAPub), "key %v != %v", i, j)
		}
	}

	// // Write keys to file for test fixtures for signing
	for i, key := range keys {
		bz, err := json.MarshalIndent(&key, "", "  ")
		require.NoError(t, err)
		t.Log(string(bz))

		WriteTestKey(i, bz)
	}

	for i, partyID := range partyIDs {
		bz, err := json.MarshalIndent(&partyID, "", "  ")
		require.NoError(t, err)
		t.Log(string(bz))

		WriteTestPartyID(i, bz)
	}
}
