package contract_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedERC20(t *testing.T) {
	require.NotEmpty(t, contract.ERC20MintableBurnableJSON)
}

func TestERC20Pack(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}
	tests := []struct {
		name    string
		method  string
		args    []interface{}
		errArgs errArgs
	}{
		{
			"constructor",
			"",
			[]interface{}{
				"Wrapped Ethereum",
				"WETH",
				uint8(18),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"constructor - wrong types",
			"",
			[]interface{}{
				"Wrapped Ethereum",
				1234,
				uint8(18),
			},
			errArgs{
				expectPass: false,
				contains:   "cannot use int as type string as argument",
			},
		},
		{
			"transfer",
			"transfer",
			[]interface{}{
				common.Address{},
				big.NewInt(1234),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"transfer - wrong types",
			"transfer",
			[]interface{}{
				common.Address{},
				1234,
			},
			errArgs{
				expectPass: false,
				contains:   "cannot use int as type ptr as argument",
			},
		},
		{
			"mint - should exist",
			"mint",
			[]interface{}{
				common.Address{},
				big.NewInt(1234),
			},
			errArgs{
				expectPass: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := contract.ERC20MintableBurnableContract.ABI.Pack(
				tc.method,
				tc.args...,
			)

			if tc.errArgs.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})

	}

}
