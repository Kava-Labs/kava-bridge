package mp_tss_test

import (
	"encoding/json"
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/os"
)

func PreParamPath(index int) string {
	return path.Join("test-fixtures", fmt.Sprintf("pre-params%d.json", index))
}

func LoadTestPreParam(index int) *keygen.LocalPreParams {
	path := PreParamPath(index)
	bytes := os.MustReadFile(path)

	var preParams *keygen.LocalPreParams
	if err := json.Unmarshal(bytes, &preParams); err != nil {
		panic(err)
	}

	return preParams
}

func TestPreParam(t *testing.T) {
	for i := 0; i < 2; i++ {
		preParams := LoadTestPreParam(i)
		require.True(t, preParams.Validate(), "test-fixture preparams should be valid")
	}
}

func TestCreatePreParamFixtures(t *testing.T) {
	// Takes a while to run -- only run when generating new fixtures or adding
	// more parties
	t.Skip("skip create pre-param fixtures")

	for i := 0; i < partyCount; i++ {
		t.Logf("creating pre-params for party %d", i)

		preParams, err := keygen.GeneratePreParams(1 * time.Minute)
		require.NoError(t, err)

		b, err := json.MarshalIndent(preParams, "", "  ")
		require.NoError(t, err)

		os.MustWriteFile(PreParamPath(i), b, 0644)
	}
}
