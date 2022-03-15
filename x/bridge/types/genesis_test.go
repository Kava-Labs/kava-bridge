package types_test

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMaxInt_Overflow(t *testing.T) {
	require.PanicsWithValue(t, "Int overflow", func() {
		types.MaxWithdrawSequence.Add(sdk.OneInt())
	})
}

func TestMaxInt(t *testing.T) {
	bitLen := types.MaxWithdrawSequence.BigInt().BitLen()
	require.Equal(t, 256, bitLen, "maxint bitlen should be 256, same as sdk.Int maxBitLen")

	bytes := types.MaxWithdrawSequence.BigInt().Bytes()
	require.Len(t, bytes, 32, "maxint should be 32 bytes, 256 bits")

	for _, b := range bytes {
		require.Equal(t, byte(math.MaxUint8), b, "all bytes should be max value")
	}
}
