package testutil

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/test"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/os"
)

func KeyPath(index int) string {
	return path.Join("..", "testutil", "test-fixtures", fmt.Sprintf("localparty-savedata%d.json", index))
}

func ReadTestKey(index int) keygen.LocalPartySaveData {
	path := KeyPath(index)

	bytes := os.MustReadFile(path)

	var key keygen.LocalPartySaveData
	if err := json.Unmarshal(bytes, &key); err != nil {
		panic(err)
	}

	return key
}

func GetTestTssKeys(count int) ([]keygen.LocalPartySaveData, tss.SortedPartyIDs) {
	var keys []keygen.LocalPartySaveData
	for i := 0; i < count; i++ {
		key := ReadTestKey(i)
		keys = append(keys, key)
	}

	signPIDsUnsorted := make(tss.UnSortedPartyIDs, len(keys))
	for i, key := range keys {
		pMoniker := fmt.Sprintf("%d", i+1)
		signPIDsUnsorted[i] = tss.NewPartyID(pMoniker, pMoniker, key.ShareID)
	}

	signPIDs := tss.SortPartyIDs(signPIDsUnsorted)
	// Sort keys so they match keys order, signing/resharing will fail if the keys
	// are mismatched to the wrong party ID
	sort.Slice(keys, func(i, j int) bool { return keys[i].ShareID.Cmp(keys[j].ShareID) == -1 })

	return keys, signPIDs
}

func WriteTestKey(index int, bz []byte) {
	os.MustWriteFile(KeyPath(index), bz, 0600)
}

func TestLoadKey(t *testing.T) {
	for i := 0; i < test.TestParticipants; i++ {
		key := ReadTestKey(i)
		require.True(t, key.Validate(), "test-fixture keys should be valid")
	}
}
