package mp_tss_test

import (
	"encoding/json"
	"fmt"
	"path"
	"testing"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/os"
)

func KeyPath(index int) string {
	return path.Join("test-fixtures", fmt.Sprintf("localparty-savedata%d.json", index))
}

func LoadKey(path string) keygen.LocalPartySaveData {
	bytes := os.MustReadFile(path)

	var key keygen.LocalPartySaveData
	if err := json.Unmarshal(bytes, &key); err != nil {
		panic(err)
	}

	return key
}

func LoadTestKey(index int) keygen.LocalPartySaveData {
	return LoadKey(KeyPath(index))
}

func WriteTestKey(index int, bz []byte) {
	os.MustWriteFile(KeyPath(index), bz, 0600)
}

func TestLoadKey(t *testing.T) {
	for i := 0; i < 2; i++ {
		key := LoadTestKey(i)
		require.True(t, key.Validate(), "test-fixture keys should be valid")
	}
}
