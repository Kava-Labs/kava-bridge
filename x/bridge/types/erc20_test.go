package types_test

import (
	"testing"

	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - same address",
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			errArgs{
				expectPass: false,
				contains:   "external and internal bytes are same",
			},
		},
		{
			"invalid - zero external",
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000000"),
			testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			errArgs{
				expectPass: false,
				contains:   "external address cannot be zero value",
			},
		},
		{
			"invalid - zero internal",
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000000"),
			errArgs{
				expectPass: false,
				contains:   "internal address cannot be zero value",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pair := types.NewERC20BridgePair(tc.externalAddress, tc.internalAddress)

			require.Equal(t, pair.GetExternalAddress(), tc.externalAddress)
			require.Equal(t, pair.GetInternalAddress(), tc.internalAddress)

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
				ExternalERC20Address: testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001").Bytes(),
				InternalERC20Address: testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000002").Bytes(),
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
				ExternalERC20Address: testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001").Bytes(),
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

func TestNewERC20BridgePairs_Valid(t *testing.T) {
	pairs := types.NewERC20BridgePairs(
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
		),
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
		),
	)

	err := pairs.Validate()
	require.NoError(t, err)
}

func TestNewERC20BridgePairs_BasicInvalid(t *testing.T) {
	pairs := types.NewERC20BridgePairs(
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
		),
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
			testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000000"),
		),
	)

	err := pairs.Validate()
	require.Error(t, err)
	require.Equal(t, "internal address cannot be zero value 0x0000000000000000000000000000000000000000", err.Error())
}

func TestNewERC20BridgePairs_DuplicateInternal(t *testing.T) {
	pairs := types.NewERC20BridgePairs(
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
		),
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
		),
	)

	err := pairs.Validate()
	require.Error(t, err)
	require.Equal(t, "found duplicate enabled bridge pair internal address 0x000000000000000000000000000000000000000A", err.Error())
}

func TestNewERC20BridgePairs_DuplicateExternal(t *testing.T) {
	pairs := types.NewERC20BridgePairs(
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
		),
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
		),
	)

	err := pairs.Validate()
	require.Error(t, err)
	require.Equal(t, "found duplicate enabled bridge pair external address 0x0000000000000000000000000000000000000001", err.Error())
}

func TestGetID(t *testing.T) {
	pair := types.NewERC20BridgePair(
		testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
		testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
	)

	extAddr := pair.GetExternalAddress().Bytes()
	intAddr := pair.GetInternalAddress().Bytes()

	// Make a copy instead of append
	s := make([]byte, len(extAddr)+len(intAddr))
	copy(s, extAddr)
	copy(s[len(extAddr):], intAddr)

	expID := tmhash.Sum(s)

	id := pair.GetID()
	require.Equal(t, expID, id)
}
