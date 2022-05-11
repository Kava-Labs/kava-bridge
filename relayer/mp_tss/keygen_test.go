package mp_tss_test

import (
	"testing"

	logging "github.com/ipfs/go-log/v2"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/stretchr/testify/require"
)

func TestKeygen(t *testing.T) {
	err := logging.SetLogLevel("*", "debug")
	require.NoError(t, err)

	count := 2
	threshold := 1 // 1 of 2

	// 1. Create party ID for each peer, share with other peers
	partyIDs := tss.GenerateTestPartyIDs(count)

	// 2. Generate pre-params and params for each peer
	var preParamsSlice []*keygen.LocalPreParams
	var paramsSlice []*tss.Parameters
	for i := range partyIDs {
		// Load from disk to avoid re-generating
		preParams := LoadTestPreParam(i)

		params, err := mp_tss.CreateParams(partyIDs.ToUnSorted(), partyIDs[i], threshold)
		require.NoError(t, err)

		preParamsSlice = append(preParamsSlice, preParams)
		paramsSlice = append(paramsSlice, params)
	}

	// 3. Create and connect transport between peers
	transports := CreateAndConnectTransports(t, partyIDs)

	// 4. Start keygen party for each peer
	errAgg := make(chan *tss.Error)
	outputAgg := make(chan keygen.LocalPartySaveData)

	for i := 0; i < count; i++ {
		outputCh, errCh := mp_tss.RunKeyGen(preParamsSlice[i], paramsSlice[i], transports[i])
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

	for range partyIDs {
		select {
		case output := <-outputAgg:
			keys = append(keys, output)
		case err := <-errAgg:
			t.Fatal(err)
		}
	}

	// make sure everyone has the same ECDSA public key
	require.True(t, keys[0].ECDSAPub.Equals(keys[1].ECDSAPub), "ECDSA public keys should be equal")

	// // Write to file for fixtures
	// for i, key := range keys {
	// 	bz, err := json.MarshalIndent(&key, "", "  ")
	// 	require.NoError(t, err)
	// 	t.Log(string(bz))
	//
	// 	WriteTestKey(i, bz)
	// }
}
