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

func LoadPreParams(path string) *keygen.LocalPreParams {
	bytes := os.MustReadFile(path)

	var preParams *keygen.LocalPreParams
	if err := json.Unmarshal(bytes, &preParams); err != nil {
		panic(err)
	}

	return preParams
}

func LoadTestPreParam(index int) *keygen.LocalPreParams {
	path := path.Join("test-fixtures", fmt.Sprintf("pre-params%d.json", index))
	return LoadPreParams(path)
}

func TestPreParam(t *testing.T) {
	for i := 0; i < 2; i++ {
		preParams := LoadTestPreParam(i)
		require.True(t, preParams.Validate(), "test-fixture preparams should be valid")
	}
}
