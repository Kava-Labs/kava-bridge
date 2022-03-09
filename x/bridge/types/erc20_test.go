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
			types.NewExternalEVMAddress(common.HexToAddress("0x01")),
			types.NewInternalEVMAddress(common.HexToAddress("0x02")),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - same address",
			types.NewExternalEVMAddress(common.HexToAddress("0x01")),
			types.NewInternalEVMAddress(common.HexToAddress("0x01")),
			errArgs{
				expectPass: false,
				contains:   "external and internal bytes are same",
			},
		},
		{
			"invalid - zero external",
			types.NewExternalEVMAddress(common.HexToAddress("0x00")),
			types.NewInternalEVMAddress(common.HexToAddress("0x01")),
			errArgs{
				expectPass: false,
				contains:   "external address cannot be zero value",
			},
		},
		{
			"invalid - zero internal",
			types.NewExternalEVMAddress(common.HexToAddress("0x01")),
			types.NewInternalEVMAddress(common.HexToAddress("0x00")),
			errArgs{
				expectPass: false,
				contains:   "external address cannot be zero value",
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

func TestNewERC20BridgePair_Direct(t *testing.T) {
	type errArgs struct {
		expectPass bool
		contains   string
	}
	tests := []struct {
		name    string
		pair    types.ERC20BridgePair
		errArgs errArgs
	}{
		{
			"valid",
			types.ERC20BridgePair{
				ExternalERC20Address: common.HexToAddress("0x01").Bytes(),
				InternalERC20Address: common.HexToAddress("0x02").Bytes(),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - invalid external length",
			types.ERC20BridgePair{
				ExternalERC20Address: []byte{1},
				InternalERC20Address: []byte{2},
			},
			errArgs{
				expectPass: false,
				contains:   "external address length is 1 but expected 20",
			},
		},
		{
			"invalid - invalid internal length",
			types.ERC20BridgePair{
				ExternalERC20Address: common.HexToAddress("0x01").Bytes(),
				InternalERC20Address: []byte{2},
			},
			errArgs{
				expectPass: false,
				contains:   "internal address length is 1 but expected 20",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.pair.Validate()
			if tc.errArgs.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errArgs.contains)
			}
		})
	}
}
