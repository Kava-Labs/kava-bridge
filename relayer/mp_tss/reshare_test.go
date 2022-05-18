package mp_tss_test

import (
	"encoding/json"
	"testing"

	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/test"
	"github.com/binance-chain/tss-lib/tss"
	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReshare(t *testing.T) {
	// newTotalPartyCount := 10
	// Number of participants in resharing -- t+1 + num new peers
	// newThreshold := 9

	newThreshold := test.TestThreshold

	err := logging.SetLogLevel("*", "debug")
	require.NoError(t, err)

	// 1. Get t+1 current keys
	oldKeys, oldPartyIDs, err := keygen.LoadKeygenTestFixturesRandomSet(keygen.TestThreshold+1, keygen.TestParticipants)
	require.NoError(t, err)
	require.Equal(t, keygen.TestThreshold+1, len(oldKeys))
	require.Equal(t, keygen.TestThreshold+1, len(oldPartyIDs))

	oldP2PCtx := tss.NewPeerContext(oldPartyIDs)

	// 2. Create new party IDs to add.. or replace ? confused
	newPartyIDs := tss.GenerateTestPartyIDs(test.TestParticipants)
	require.Len(t, newPartyIDs, test.TestParticipants)

	t.Logf("old partyIDs: %v", oldPartyIDs)
	t.Logf("new partyIDs: %v", newPartyIDs)

	// 3. Create and connect transport between peers
	oldTransports, newTransports := CreateAndConnectReSharingTransports(t, oldPartyIDs, newPartyIDs)

	newP2PCtx := tss.NewPeerContext(newPartyIDs)
	newPCount := len(newPartyIDs)
	//
	// oldCommittee := make([]*resharing.LocalParty, 0, len(oldPartyIDs))
	// newCommittee := make([]*resharing.LocalParty, 0, newPCount)
	// bothCommitteesPax := len(oldCommittee) + len(newCommittee)
	//
	// errCh := make(chan *tss.Error, bothCommitteesPax)
	// outCh := make(chan tss.Message, bothCommitteesPax)
	// endCh := make(chan keygen.LocalPartySaveData, bothCommitteesPax)

	// updater := test.SharedPartyUpdater

	// 4. Start resharing party for each peer
	outputAgg := make(chan keygen.LocalPartySaveData)
	errAgg := make(chan *tss.Error)

	// Start old parties
	for i, partyID := range oldPartyIDs {
		params := tss.NewReSharingParameters(tss.S256(), oldP2PCtx, newP2PCtx, partyID, test.TestParticipants, threshold, newPCount, newThreshold)

		// P := resharing.NewLocalParty(params, oldKeys[i], outCh, endCh).(*resharing.LocalParty) // discard old key data
		// oldCommittee = append(oldCommittee, P)

		outputCh, errCh := mp_tss.RunReshare(params, oldKeys[i], oldTransports[i])

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
		params := tss.NewReSharingParameters(tss.S256(), oldP2PCtx, newP2PCtx, partyID, test.TestParticipants, threshold, newPCount, newThreshold)
		t.Log(params.PartyID())

		save := keygen.NewLocalPartySaveData(newPCount)
		// Verify() fail is not from this
		save.LocalPreParams = ReadTestKey(i).LocalPreParams

		// require.True(t, save.Validate(), "new party save data should be valid")

		// P := resharing.NewLocalParty(params, save, outCh, endCh).(*resharing.LocalParty)
		// newCommittee = append(newCommittee, P)

		outputCh, errCh := mp_tss.RunReshare(params, save, newTransports[i])

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

	newKeys := make([]keygen.LocalPartySaveData, newPCount)

	for i := 0; i < len(oldPartyIDs)+len(newPartyIDs); i++ {
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

	require.Equal(t, len(newKeys), newPCount, "each party should get a key")

	// xj tests: BigXj == xj*G
	for j, key := range newKeys {
		// xj test: BigXj == xj*G
		xj := key.Xi
		gXj := crypto.ScalarBaseMult(tss.S256(), xj)
		BigXj := key.BigXj[j]
		assert.True(t, BigXj.Equals(gXj), "ensure BigX_j == g^x_j")
	}

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
