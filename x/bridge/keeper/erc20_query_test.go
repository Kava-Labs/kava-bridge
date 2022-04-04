package keeper_test

import (
	"math/big"
	"testing"

	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
)

type ERC20QueryTestSuite struct {
	testutil.Suite

	contractAddr types.InternalEVMAddress
}

func TestERC20QueryTestSuite(t *testing.T) {
	suite.Run(t, new(ERC20QueryTestSuite))
}

func (suite *ERC20QueryTestSuite) SetupTest() {
	suite.Suite.SetupTest()
	suite.contractAddr = suite.DeployERC20()
}

func (suite *ERC20QueryTestSuite) TestERC20QueryBalanceOf_Empty() {
	bal, err := suite.App.BridgeKeeper.QueryERC20BalanceOf(
		suite.Ctx,
		suite.contractAddr,
		suite.Key1Addr,
	)
	suite.Require().NoError(err)
	suite.Require().True(bal.Cmp(big.NewInt(0)) == 0, "balance should be 0")
}

func (suite *ERC20QueryTestSuite) TestERC20QueryBalanceOf_NonEmpty() {
	// Mint some tokens for the address
	err := suite.App.BridgeKeeper.MintERC20(
		suite.Ctx,
		suite.contractAddr,
		suite.Key1Addr,
		big.NewInt(10),
	)
	suite.Require().NoError(err)

	bal, err := suite.App.BridgeKeeper.QueryERC20BalanceOf(
		suite.Ctx,
		suite.contractAddr,
		suite.Key1Addr,
	)
	suite.Require().NoError(err)
	suite.Require().Equal(big.NewInt(10), bal, "balance should be 10")
}
