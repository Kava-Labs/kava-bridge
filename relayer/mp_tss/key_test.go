package mp_tss_test

import (
	"encoding/json"
	"fmt"
	"path"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/os"
)

// -----------------------------------------------------------------------------
// keygen

func KeyPath(index int) string {
	return path.Join("test-fixtures", fmt.Sprintf("localparty-savedata%d.json", index))
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

func GetTestKeys(count int) []keygen.LocalPartySaveData {
	var keys []keygen.LocalPartySaveData
	for i := 0; i < count; i++ {
		key := ReadTestKey(i)
		keys = append(keys, key)
	}

	return keys
}

func WriteTestKey(index int, bz []byte) {
	os.MustWriteFile(KeyPath(index), bz, 0600)
}

// -----------------------------------------------------------------------------
// partyIDs

func PartyIDPath(index int) string {
	return path.Join("test-fixtures", fmt.Sprintf("partyid%d.json", index))
}

func ReadPartyID(index int) *tss.PartyID {
	path := PartyIDPath(index)

	bytes := os.MustReadFile(path)

	var partyID tss.PartyID
	if err := json.Unmarshal(bytes, &partyID); err != nil {
		panic(err)
	}

	return &partyID
}

func GetTestPartyIDs(count int) tss.SortedPartyIDs {
	var partyIDs []*tss.PartyID
	for i := 0; i < count; i++ {
		partyID := ReadPartyID(i)
		partyIDs = append(partyIDs, partyID)
	}

	return tss.SortPartyIDs(partyIDs)
}

func WriteTestPartyID(index int, bz []byte) {
	os.MustWriteFile(PartyIDPath(index), bz, 0600)
}

func TestLoadKey(t *testing.T) {
	for i := 0; i < partyCount; i++ {
		key := ReadTestKey(i)
		require.True(t, key.Validate(), "test-fixture keys should be valid")
	}
}
