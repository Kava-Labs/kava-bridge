package types_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestNewERC20BridgePair(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}
	tests := []struct {
		name            string
		externalAddress types.ExternalEVMAddress
		internalAddress types.InternalEVMAddress
		errArgs         errArgs
	}{
		{
			"valid",
			types.ExternalEVMAddress{
				Address: common.HexToAddress("0x01"),
			},
			types.InternalEVMAddress{
				Address: common.HexToAddress("0x02"),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - same address",
			types.ExternalEVMAddress{
				Address: common.HexToAddress("0x01"),
			},
			types.InternalEVMAddress{
				Address: common.HexToAddress("0x01"),
			},
			errArgs{
				expectPass: false,
				contains:   "external and internal bytes are same",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pair := types.NewERC20BridgePair(tc.externalAddress, tc.internalAddress)

			require.Equal(t, pair.ExternalERC20Address, tc.externalAddress.Bytes())
			require.Equal(t, pair.InternalERC20Address, tc.internalAddress.Bytes())

			err := pair.Validate()
			if tc.errArgs.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}

}
