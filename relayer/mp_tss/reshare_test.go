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

func TestReshare(t *testing.T) {
	newTotalPartyCount := 10
	// Number of participants in resharing -- t+1 + num new peers
	newThreshold := 9

	err := logging.SetLogLevel("*", "debug")
	require.NoError(t, err)

	// 1. Get t+1 current keys
	oldKeys := GetTestKeys(threshold + 1)
	require.Len(t, oldKeys, threshold+1)

	// Recreate party IDs from keys
	oldPartyIDs := GetTestPartyIDs(threshold + 1)
	require.Len(t, oldPartyIDs, threshold+1)

	// 2. Create new party IDs to add.. or replace ? confused
	newPartyIDs := tss.GenerateTestPartyIDs(newTotalPartyCount, len(oldPartyIDs))
	require.Len(t, newPartyIDs, newTotalPartyCount)

	t.Log(newPartyIDs)

	allPartyIDs := tss.UnSortedPartyIDs(append(oldPartyIDs, newPartyIDs...))
	require.Len(t, allPartyIDs, threshold+1+newTotalPartyCount)
	require.Len(t, allPartyIDs, 19) // 9 old to sign + 10 new

	// 3. Create and connect transport between peers
	transports := CreateAndConnectTransports(t, allPartyIDs)

	// Add committee lists to transports
	for _, transport := range transports {
		transport.AddOldCommittee(oldPartyIDs...)
		transport.AddNewCommittee(newPartyIDs...)
	}

	// 4. Start signing party for each peer
	outputAgg := make(chan keygen.LocalPartySaveData)
	errAgg := make(chan *tss.Error)

	// Start old parties
	for i, partyID := range oldPartyIDs {
		params := mp_tss.CreateReShareParams(
			oldPartyIDs.ToUnSorted(),
			newPartyIDs.ToUnSorted(),
			partyID,
			threshold,    // 8
			newThreshold, // 12
		)
		t.Log(params.PartyID())

		outputCh, errCh := mp_tss.RunReshare(params, oldKeys[i], transports[i])

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

	for i, partyID := range newPartyIDs {
		params := mp_tss.CreateReShareParams(
			oldPartyIDs.ToUnSorted(),
			newPartyIDs.ToUnSorted(),
			partyID,
			threshold,    // 8
			newThreshold, // 12
		)
		t.Log(params.PartyID())

		preParams := LoadTestPreParam(i)
		save := keygen.NewLocalPartySaveData(newTotalPartyCount)
		save.LocalPreParams = *preParams

		require.True(t, save.Validate(), "new party save data should be valid")

		outputCh, errCh := mp_tss.RunReshare(params, save, transports[len(oldPartyIDs)+i])

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

	t.Logf("started key reshare")

	var newKeys []keygen.LocalPartySaveData

	for range allPartyIDs {
		select {
		case output := <-outputAgg:
			bz, err := json.Marshal(&output)
			require.NoError(t, err)
			t.Log(string(bz))

			newKeys = append(newKeys, output)
		case err := <-errAgg:
			t.Fatal(err)
		}
	}

	require.Len(t, newKeys, newPartyIDs.Len()+oldPartyIDs.Len(), "each party should get a key")

	// New reshared pubkey should match old pubkey
	assert.Truef(t, oldKeys[0].ECDSAPub.Equals(newKeys[0].ECDSAPub), "reshared pubkey should match old pubkey")

	// make sure everyone has the same ECDSA public key
	for i, key := range newKeys {
		for j, key2 := range newKeys {
			// Skip self and previous keys
			if j <= i {
				continue
			}

			assert.Truef(t, key.ECDSAPub.Equals(key2.ECDSAPub), "key %v != %v", i, j)
		}
	}
}
