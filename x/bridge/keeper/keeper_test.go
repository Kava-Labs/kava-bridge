package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"

	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type KeeperTestSuite struct {
	testutil.Suite
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestERC20_NotEnabled() {
	// WETH but last char changed
	extAddr := testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc4")

	_, err := suite.App.BridgeKeeper.GetOrDeployInternalERC20(suite.Ctx, extAddr)
	suite.Require().Error(err)
	suite.Require().ErrorIs(err, types.ErrERC20NotEnabled)
}

func (suite *KeeperTestSuite) TestERC20SaveDeploy() {
	extAddr := testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

	_, found := suite.App.BridgeKeeper.GetInternalERC20Address(suite.Ctx, extAddr)
	suite.Require().False(found, "internal ERC20 address should not be set before first bridge")

	firstInternal, err := suite.App.BridgeKeeper.GetOrDeployInternalERC20(suite.Ctx, extAddr)
	suite.Require().NoError(err)

	// Fetch from store
	savedInternal, found := suite.App.BridgeKeeper.GetInternalERC20Address(suite.Ctx, extAddr)
	suite.Require().True(found, "internal ERC20 address should be saved after first bridge")
	suite.Require().Equal(firstInternal, savedInternal, "deployed address should match saved internal ERC20 address")

	// Fetch addr again to make sure we get the same one and another ERC20 isn't deployed
	secondInternal, err := suite.App.BridgeKeeper.GetOrDeployInternalERC20(suite.Ctx, extAddr)
	suite.Require().NoError(err)

	suite.Require().Equal(firstInternal, secondInternal, "second call should return the saved internal ERC20 address")
}

func (suite *KeeperTestSuite) TestPermission() {
	relayerAddr := suite.App.BridgeKeeper.GetRelayer(suite.Ctx)
	// Check relayer address before actually testing
	suite.Require().Equal(relayerAddr, suite.RelayerAddress, "test suite relayer should match relayer set in params")
	suite.Require().NotEmpty(relayerAddr, "relayer address should not be empty")

	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		signers []sdk.AccAddress
		key     *ethsecp256k1.PrivKey
		errArgs errArgs
	}{
		{
			"valid - signer matches relayer in params",
			[]sdk.AccAddress{relayerAddr},
			suite.RelayerKey,
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - multiple signers even if permissioned",
			[]sdk.AccAddress{
				relayerAddr,
				sdk.AccAddress(suite.Key1.PubKey().Address()),
			},
			suite.RelayerKey,
			errArgs{
				expectPass: false,
				contains:   "invalid number of signers",
			},
		},
		{
			"invalid - single unknown address",
			[]sdk.AccAddress{sdk.AccAddress(suite.Key1.PubKey().Address())},
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "signer not authorized for bridge message",
			},
		},
		{
			"invalid - multiple unknown addresses",
			[]sdk.AccAddress{
				sdk.AccAddress(suite.Key1.PubKey().Address()),
				sdk.AccAddress(suite.Key2.PubKey().Address()),
			},
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "invalid number of signers",
			},
		},
		{
			"invalid - empty signers",
			[]sdk.AccAddress{},
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "invalid number of signers",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			err := suite.App.BridgeKeeper.IsSignerAuthorized(suite.Ctx, tc.signers)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestERC20PairIter() {
	pairs := types.NewERC20BridgePairs(
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
		),
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
		),
	)

	for _, pair := range pairs {
		suite.App.BridgeKeeper.SetERC20AddressPair(suite.Ctx, pair)
	}

	var iterPairs types.ERC20BridgePairs
	suite.App.BridgeKeeper.IterateERC20BridgePairs(suite.Ctx, func(pair types.ERC20BridgePair) bool {
		iterPairs = append(iterPairs, pair)
		return false
	})

	suite.Require().Equal(pairs, iterPairs, "pairs from iterator should match pairs set in store")
}

func (suite *KeeperTestSuite) TestInitialNextWithdrawSequence() {
	readSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(types.DefaultNextWithdrawSequence, readSeq)
	suite.Require().NotZero(readSeq)
}

func (suite *KeeperTestSuite) TestIncrementNextWithdrawSequence() {
	var seq = sdk.NewInt(123456)
	suite.App.BridgeKeeper.SetNextWithdrawSequence(suite.Ctx, seq)

	suite.Require().NoError(suite.App.BridgeKeeper.IncrementNextWithdrawSequence(suite.Ctx))

	// check seq was incremented
	readSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(seq.AddRaw(1), readSeq)
}

func (suite *KeeperTestSuite) TestIncrementNextWithdrawSequence_Wrap() {
	seq := types.MaxWithdrawSequence
	suite.App.BridgeKeeper.SetNextWithdrawSequence(suite.Ctx, seq)

	suite.Require().NoError(suite.App.BridgeKeeper.IncrementNextWithdrawSequence(suite.Ctx))

	// check seq wrapped back to 0
	readSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.ZeroInt(), readSeq)
}

func TestWrappingAddInt(t *testing.T) {
	tests := []struct {
		name    string
		i1      sdk.Int
		i2      sdk.Int
		wantSum sdk.Int
	}{
		{
			"1+1=2",
			sdk.OneInt(), sdk.OneInt(), sdk.NewInt(2),
		},
		{
			"0+1=1",
			sdk.ZeroInt(), sdk.OneInt(), sdk.NewInt(1),
		},
		{
			"100+120=220",
			sdk.NewInt(100), sdk.NewInt(120), sdk.NewInt(100 + 120),
		},
		{
			"max+1=0",
			types.MaxWithdrawSequence, sdk.NewInt(1), sdk.NewInt(0),
		},
		{
			"max+2=1",
			types.MaxWithdrawSequence, sdk.NewInt(2), sdk.NewInt(1),
		},
		{
			"max+max=max-1",
			types.MaxWithdrawSequence, types.MaxWithdrawSequence, types.MaxWithdrawSequence.SubRaw(1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sum := keeper.WrappingAddInt(tc.i1, tc.i2)
			require.NotEqual(t, sdk.Int{}, sum, "sum should not be default int value")
			// Use string match here due to non-equal bigint abs type: (big.nat) nil vs abs: (big.nat) {}
			require.Equalf(
				t,
				tc.wantSum.String(),
				sum.String(),
				"sum should match expected value, want %v but got %v", tc.wantSum, sum,
			)
		})
	}
}
