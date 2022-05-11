package mp_tss_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"testing"

	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/crypto/sha3"

	"github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	err := logging.SetLogLevel("*", "debug")
	require.NoError(t, err)

	count := 2
	threshold := 1

	// 1. Create party ID for each peer, share with other peers
	partyIDs := tss.GenerateTestPartyIDs(count)

	transports := CreateAndConnectTransports(t, partyIDs)

	msg := []byte("hello world")
	hash := sha3.Sum256(msg)
	hashBigInt := new(big.Int).SetBytes(hash[:])

	// 4. Start keygen party for each peer
	outputAgg := make(chan common.SignatureData)
	errAgg := make(chan *tss.Error)

	for i, partyID := range partyIDs {
		params, err := mp_tss.CreateParams(partyIDs.ToUnSorted(), partyID, threshold)
		require.NoError(t, err)

		key := LoadTestKey(i)

		outputCh, errCh := mp_tss.RunSigner(hashBigInt, params, key, transports[i])

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
