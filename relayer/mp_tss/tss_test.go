package mp_tss_test

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"time"

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
			if transport.PartyID.Index == otherTransport.PartyID.Index {
				t.Logf("skipping self for transport: %v == %v", transport.PartyID.Index, otherTransport.PartyID.Index)
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

	errAgg := make(chan *tss.Error)
	outputAgg := make(chan keygen.LocalPartySaveData)

	for _, errCh := range errChs {
		go func(errCh chan *tss.Error) {
			for err := range errCh {
				errAgg <- err
			}
		}(errCh)
	}

	for _, outputCh := range outputChs {
		go func(outputCh chan keygen.LocalPartySaveData) {
			for output := range outputCh {
				outputAgg <- output
			}
		}(outputCh)
	}

	go func() {
		for {
			t.Logf("goroutines: %v", runtime.NumGoroutine())
			time.Sleep(time.Second * 5)
		}
	}()

	var keys []keygen.LocalPartySaveData
	select {
	case output := <-outputAgg:
		keys = append(keys, output)
	case err := <-errAgg:
		t.Fatal(err)
	}

	// 5. Output keys
	bz, err := json.Marshal(&keys)
	require.NoError(t, err)
	fmt.Println(string(bz))
}
