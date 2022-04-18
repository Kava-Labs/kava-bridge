package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
)

type ParamsTestSuite struct {
	testutil.Suite
}

func TestParamsSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestGetSetRelayer() {
	relayer := suite.App.BridgeKeeper.GetRelayer(suite.Ctx)
	suite.Require().Equal(suite.RelayerAddress, relayer, "relayer should match address set in genesis")

	newRelayer := sdk.AccAddress("hi")
	suite.NotPanics(func() {
		suite.App.BridgeKeeper.SetRelayer(suite.Ctx, newRelayer)
	})
	relayer = suite.App.BridgeKeeper.GetRelayer(suite.Ctx)
	suite.Require().Equal(newRelayer, relayer)
}

func (suite *ParamsTestSuite) TestGetEnabledERC20Token() {
	token, err := suite.App.BridgeKeeper.GetEnabledERC20TokenFromExternal(
		suite.Ctx,
		testutil.MustNewExternalEVMAddressFromString("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	)
	suite.Require().NoError(err)

	expectedToken := types.NewEnabledERC20Token(
		testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		"Wrapped Ether",
		"WETH",
		18,
		testutil.MinWETHWithdrawAmount,
	)

	suite.Require().Equal(expectedToken, token)
}

func (suite *ParamsTestSuite) TestGetEnabledERC20Token_NotFound() {
	_, err := suite.App.BridgeKeeper.GetEnabledERC20TokenFromExternal(
		suite.Ctx,
		testutil.MustNewExternalEVMAddressFromString("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc4"),
	)
	suite.Require().Error(err)
}

func (suite *ParamsTestSuite) TestGetEnabledERC20Tokens() {
	token := suite.App.BridgeKeeper.GetEnabledERC20Tokens(suite.Ctx)
	suite.Require().Len(token, 3)
}
