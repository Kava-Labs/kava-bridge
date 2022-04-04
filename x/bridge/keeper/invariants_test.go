package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type InvariantTestSuite struct {
	testutil.Suite

	invariants   map[string]map[string]sdk.Invariant
	contractAddr types.InternalEVMAddress
}

func TestInvariantTestSuite(t *testing.T) {
	suite.Run(t, new(InvariantTestSuite))
}

func (suite *InvariantTestSuite) SetupTest() {
	suite.Suite.SetupTest()

	suite.contractAddr = suite.DeployERC20()
	suite.invariants = make(map[string]map[string]sdk.Invariant)
	keeper.RegisterInvariants(suite, suite.Keeper)
}

func (suite *InvariantTestSuite) SetupValidState() {
	err := suite.Keeper.MintERC20(suite.Ctx, suite.contractAddr, suite.Key1Addr, big.NewInt(1000000))
	suite.Require().NoError(err)

	// key1 ERC20 bal -10000, sdk.Coin +1000
	// Module account balance 0 -> 1000
	_, err = suite.Keeper.CallEVM(
		suite.Ctx,
		contract.ERC20MintableBurnableContract.ABI,
		suite.Key1Addr.Address,
		suite.contractAddr,
		"convertToCoin",
		// convertToCoin ERC20 args
		suite.Key1Addr.Address,
		big.NewInt(1000),
	)
	suite.Require().NoError(err)
}

// RegisterRoutes implements sdk.InvariantRegistry
func (suite *InvariantTestSuite) RegisterRoute(moduleName string, route string, invariant sdk.Invariant) {
	_, exists := suite.invariants[moduleName]

	if !exists {
		suite.invariants[moduleName] = make(map[string]sdk.Invariant)
	}

	suite.invariants[moduleName][route] = invariant
}

func (suite *InvariantTestSuite) runInvariant(route string, invariant func(k keeper.Keeper) sdk.Invariant) (string, bool) {
	ctx := suite.Ctx
	registeredInvariant := suite.invariants[types.ModuleName][route]
	suite.Require().NotNil(registeredInvariant)

	// direct call
	dMessage, dBroken := invariant(suite.Keeper)(ctx)
	// registered call
	rMessage, rBroken := registeredInvariant(ctx)
	// all call
	aMessage, aBroken := keeper.AllInvariants(suite.Keeper)(ctx)

	// require matching values for direct call and registered call
	suite.Require().Equal(dMessage, rMessage, "expected registered invariant message to match")
	suite.Require().Equal(dBroken, rBroken, "expected registered invariant broken to match")
	// require matching values for direct call and all invariants call if broken
	suite.Require().Equal(dBroken, aBroken, "expected all invariant broken to match")
	if dBroken {
		suite.Require().Equal(dMessage, aMessage, "expected all invariant message to match")
	}

	return dMessage, dBroken
}

func (suite *InvariantTestSuite) TestBackedCoinsInvariant() {
	// default state is valid
	_, broken := suite.runInvariant("backed-coins", keeper.BackedCoinsInvariant)
	suite.Equal(false, broken)

	suite.SetupValidState()
	_, broken = suite.runInvariant("backed-coins", keeper.BackedCoinsInvariant)
	suite.Equal(false, broken)

	// break invariant creating more sdk.Coin than module account ERC20 tokens
	err := suite.BankKeeper.MintCoins(
		suite.Ctx, types.ModuleName,
		sdk.NewCoins(sdk.NewCoin("usdc", sdk.NewInt(1001))),
	)
	suite.Require().NoError(err)

	message, broken := suite.runInvariant("backed-coins", keeper.BackedCoinsInvariant)
	suite.Equal("bridge: backed coins broken invariant\ncoin supply is greater than module account ERC20 tokens\n", message)
	suite.Equal(true, broken)
}

func (suite *InvariantTestSuite) TestBridgePairs() {
	// default state is valid
	_, broken := suite.runInvariant("bridge-pairs", keeper.BridgePairsInvariant)
	suite.Equal(false, broken)

	suite.SetupValidState()
	_, broken = suite.runInvariant("bridge-pairs", keeper.BridgePairsInvariant)
	suite.Equal(false, broken)

	message, broken := suite.runInvariant("bridge-pairs", keeper.BridgePairsInvariant)
	suite.Equal("bridge: bridge pairs broken invariant\nminor balances not all less than overflow\n", message)
	suite.Equal(true, broken)
}
