package mp_tss_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"testing"

	logging "github.com/ipfs/go-log/v2"

	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	err := logging.SetLogLevel("*", "debug")
	require.NoError(t, err)

	// 1. Get party keys from file
	keys := GetTestKeys(threshold + 1)
	require.Len(t, keys, threshold+1)

	// Recreate party IDs from keys
	partyIDs := GetTestPartyIDs(threshold + 1)
	require.Len(t, partyIDs, threshold+1)

	// 2. Create and connect transport between peers
	transports := CreateAndConnectTransports(t, partyIDs)

	// 3. Create a hash to sign -- not necessarily the same as actual tx hash
	hash := sha256.Sum256([]byte("hello world"))
	hashBigInt := new(big.Int).SetBytes(hash[:])

	// 4. Start signing party for each peer
	outputAgg := make(chan common.SignatureData)
	errAgg := make(chan *tss.Error)

	for i, partyID := range partyIDs {
		params := mp_tss.CreateParams(partyIDs.ToUnSorted(), partyID, threshold)
		t.Log(params.PartyID())

		outputCh, errCh := mp_tss.RunSigner(hashBigInt, params, keys[i], transports[i])

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

	for range partyIDs {
		select {
		case output := <-outputAgg:
			bz, err := json.Marshal(&output)
			require.NoError(t, err)
			t.Log(string(bz))

			signatures = append(signatures, output)
		case err := <-errAgg:
			t.Fatal(err)
		}
	}

	// make sure everyone has the same signature
	require.True(t, bytes.Equal(signatures[0].Signature, signatures[1].Signature))
}
