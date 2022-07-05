package testutil_test

import (
	"testing"

	"github.com/binance-chain/tss-lib/test"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/stretchr/testify/require"
)

func TestGenerateNodeKeys(t *testing.T) {
	t.Skip("TestKeygen must be run if this is run again to be re-generated correct shares")

	keys, err := testutil.GenerateNodeKeys(test.TestParticipants)
	require.NoError(t, err)

	require.Len(t, keys, test.TestParticipants)

	for i, key := range keys {
		testutil.WriteP2pNodeTestKey(i, key)
	}
}
