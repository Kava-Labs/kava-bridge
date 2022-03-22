package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
)

type ConversionTestSuite struct {
	testutil.Suite
}

func TestConversionTestSuite(t *testing.T) {
	suite.Run(t, new(ConversionTestSuite))
}

func (suite *ConversionTestSuite) TestMint() {
	pair := types.NewConversionPair(
		testutil.MustNewInternalEVMAddressFromString("000000000000000000000000000000000000000A"),
		"erc20/usdc",
	)

	amount := sdk.NewInt(100)
	recipient := suite.Key1.PubKey().Address().Bytes()

	err := suite.App.BridgeKeeper.MintConversionPairCoin(suite.Ctx, pair, amount, recipient)
	suite.Require().NoError(err)

	bal := suite.App.BankKeeper.GetBalance(suite.Ctx, recipient, pair.Denom)
	suite.Require().Equal(amount, bal.Amount, "minted amount should increase balance")
}

func (suite *ConversionTestSuite) TestBurn_InsufficientBalance() {
	pair := types.NewConversionPair(
		testutil.MustNewInternalEVMAddressFromString("000000000000000000000000000000000000000A"),
		"erc20/usdc",
	)

	amount := sdk.NewInt(100)
	recipient := suite.Key1.PubKey().Address().Bytes()

	err := suite.App.BridgeKeeper.BurnConversionPairCoin(suite.Ctx, pair, amount, recipient)
	suite.Require().Error(err)
	suite.Require().Equal("0erc20/usdc is smaller than 100erc20/usdc: insufficient funds", err.Error())
}

func (suite *ConversionTestSuite) TestBurn() {
	pair := types.NewConversionPair(
		testutil.MustNewInternalEVMAddressFromString("000000000000000000000000000000000000000A"),
		"erc20/usdc",
	)

	amount := sdk.NewInt(100)
	recipient := suite.Key1.PubKey().Address().Bytes()

	err := suite.App.BridgeKeeper.MintConversionPairCoin(suite.Ctx, pair, amount, recipient)
	suite.Require().NoError(err)

	bal := suite.App.BankKeeper.GetBalance(suite.Ctx, recipient, pair.Denom)
	suite.Require().Equal(amount, bal.Amount, "minted amount should increase balance")

	err = suite.App.BridgeKeeper.BurnConversionPairCoin(suite.Ctx, pair, amount, recipient)
	suite.Require().NoError(err)

	bal = suite.App.BankKeeper.GetBalance(suite.Ctx, recipient, pair.Denom)
	suite.Require().Equal(sdk.ZeroInt(), bal.Amount, "balance should be zero after burn")
}
