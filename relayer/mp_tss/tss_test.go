package mp_tss_test

import (
	"encoding/json"
	"fmt"
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
	threshold := 1 // 4 of 5

	// 1. Create party ID for each peer, share with other peers
	partyIDs := tss.GenerateTestPartyIDs(count)

	var setupOptions []mp_tss.SetupOptions
	for i := 0; i < count; i++ {
		setupOptions = append(setupOptions, mp_tss.SetupOptions{
			PartyIDs:  partyIDs.ToUnSorted(),
			PartyID:   partyIDs[i],
			Threshold: threshold,
		})
	}

	t.Logf("setupOptions: %+v", setupOptions)

	// 2. Generate pre-params and params for each peer
	var preParams []mp_tss.SetupOutput
	for _, opts := range setupOptions {
		params, err := mp_tss.CreateKeygenParams(opts)
		require.NoError(t, err)

		preParams = append(preParams, params)
	}

	t.Logf("preParams: %+v", preParams)

	// 3. Create transport between peers
	var transports []*mp_tss.MemoryTransporter
	for _, opts := range setupOptions {
		transports = append(transports, mp_tss.NewMemoryTransporter(opts.PartyID))
	}

	t.Logf("transports: %+v", transports)

	// Add transport receivers to each other
	for _, transport := range transports {
		for _, otherTransport := range transports {
			// Skip self
			if transport.PartyID == otherTransport.PartyID {
				continue
			}

			transport.AddTarget(otherTransport.PartyID, otherTransport.GetReceiver())
		}
	}

	t.Logf("transports connected: %+v", transports)

	// 4. Start keygen party for each peer
	var outputChs []chan keygen.LocalPartySaveData
	var errChs []chan *tss.Error
	for i := 0; i < count; i++ {
		outputCh, errCh := mp_tss.RunKeyGen(preParams[i], transports[i])
		outputChs = append(outputChs, outputCh)
		errChs = append(errChs, errCh)
	}

	t.Logf("started: %+v", outputChs)

	var keys []keygen.LocalPartySaveData
	for i := 0; i < count; i++ {
		select {
		case output := <-outputChs[i]:
			keys = append(keys, output)
		case err := <-errChs[i]:
			t.Fatal(err)
		}
	}

	// 5. Output keys
	bz, err := json.Marshal(&keys)
	require.NoError(t, err)
	fmt.Println(string(bz))
}
