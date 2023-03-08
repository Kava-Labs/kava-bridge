package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

type ConversionHooksTestSuite struct {
	testutil.Suite

	msgServer             types.MsgServer
	key1Addr              common.Address
	erc20Abi              abi.ABI
	conversionPair        types.ConversionPair
	invalidConversionPair types.ConversionPair
}

func TestConversionHooksTestSuite(t *testing.T) {
	suite.Run(t, new(ConversionHooksTestSuite))
}

func (suite *ConversionHooksTestSuite) SetupTest() {
	suite.Suite.SetupTest()

	suite.msgServer = keeper.NewMsgServerImpl(suite.App.BridgeKeeper)

	suite.erc20Abi = contract.ERC20MintableBurnableContract.ABI
	externalWethAddr := testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

	// Bridge an asset to deploy the ERC20 asset and update store with pair
	suite.key1Addr = common.BytesToAddress(suite.Key1.PubKey().Address())
	suite.submitBridgeERC20Msg(externalWethAddr, sdk.NewInt(100), suite.key1Addr)

	bridgePair, found := suite.App.BridgeKeeper.GetBridgePairFromExternal(suite.Ctx, externalWethAddr)
	suite.Require().True(found, "bridge pair must exist after bridge")

	// Cannot be set in genesis since we need to deploy the erc20 contract and get internal addr
	suite.conversionPair = types.NewConversionPair(bridgePair.GetInternalAddress(), "erc20/usdc")

	// Cannot be set in genesis since we need to deploy the erc20 contract and get internal addr
	suite.conversionPair = types.NewConversionPair(bridgePair.GetInternalAddress(), "erc20/wbtc")

	// Create a bridge pair that is not enabled for conversion, does not need
	// to be enabled as a bridge pair, just that it is deployed to EVM.
	bridgePair2Addr, err := suite.App.BridgeKeeper.DeployMintableERC20Contract(
		suite.Ctx,
		types.NewEnabledERC20Token(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000000"),
			"Invald Token",
			"IT",
			18,
			sdk.NewInt(10_000_000_000_000_000),
		),
	)
	suite.Require().NoError(err)

	err = suite.App.BridgeKeeper.MintERC20(suite.Ctx, bridgePair2Addr, types.NewInternalEVMAddress(suite.key1Addr), big.NewInt(100))
	suite.Require().NoError(err)
	suite.invalidConversionPair = types.NewConversionPair(bridgePair2Addr, "erc20/invalid")
}

func (suite *ConversionHooksTestSuite) TestHooksSet() {
	suite.Require().PanicsWithValue("cannot set evm hooks twice", func() {
		suite.App.EvmKeeper.SetHooks(suite.App.BridgeKeeper.ConversionHooks())
	})
}

func (suite *ConversionHooksTestSuite) submitBridgeERC20Msg(
	contractAddr types.ExternalEVMAddress,
	amount sdk.Int,
	receiver common.Address,
) {
	msg := types.NewMsgBridgeEthereumToKava(
		suite.RelayerAddress.String(),
		contractAddr.String(),
		amount,
		receiver.String(),
		sdk.NewInt(1),
	)

	_, err := suite.msgServer.BridgeEthereumToKava(sdk.WrapSDKContext(suite.Ctx), &msg)
	suite.Require().NoError(err)
}

func (suite *ConversionHooksTestSuite) ConvertToCoin(
	contractAddr types.InternalEVMAddress,
	toKavaAddr sdk.AccAddress,
	amount *big.Int,
) *evmtypes.MsgEthereumTxResponse {
	// Prevents out of gas error
	suite.Commit()

	// method is lowercase but event is upper
	data, err := suite.erc20Abi.Pack(
		"convertToCoin",
		common.BytesToAddress(toKavaAddr.Bytes()),
		amount,
	)
	suite.Require().NoError(err)

	res, err := suite.SendTx(contractAddr, suite.key1Addr, suite.Key1, data)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), "evm tx should not fail %v", res)

	return res
}

func (suite *ConversionHooksTestSuite) TestConvertToCoin() {
	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(100)

	_ = suite.ConvertToCoin(suite.conversionPair.GetAddress(), toKavaAddr, amount)
}

func (suite *ConversionHooksTestSuite) TestConvertToCoin_BridgeDisabled() {
	// Disable bridge
	params := suite.Keeper.GetParams(suite.Ctx)
	params.BridgeEnabled = false
	suite.Keeper.SetParams(suite.Ctx, params)

	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(100)

	// Manually create tx and send
	suite.Commit()

	// method is lowercase but event is upper
	data, err := suite.erc20Abi.Pack(
		"convertToCoin",
		common.BytesToAddress(toKavaAddr.Bytes()),
		amount,
	)
	suite.Require().NoError(err)

	res, err := suite.SendTx(suite.conversionPair.GetAddress(), suite.key1Addr, suite.Key1, data)
	suite.Require().NoError(err)

	suite.Require().True(res.Failed(), "evm tx should fail if bridge is disabled")
	// Does not contain the BridgeDisabled error
	suite.Require().Equal(evmtypes.ErrPostTxProcessing.Error(), res.VmError)
}

func (suite *ConversionHooksTestSuite) TestConvertToCoin_Events() {
	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(100)

	res := suite.ConvertToCoin(suite.conversionPair.GetAddress(), toKavaAddr, amount)
	suite.Require().False(res.Failed(), "tx should not fail")

	coinAmount := sdk.NewCoin(suite.conversionPair.Denom, sdk.NewIntFromBigInt(amount))
	suite.EventsContains(suite.GetEvents(), sdk.NewEvent(
		types.EventTypeConvertERC20ToCoin,
		sdk.NewAttribute(types.AttributeKeyERC20Address, suite.conversionPair.GetAddress().String()),
		sdk.NewAttribute(types.AttributeKeyInitiator, suite.key1Addr.String()),
		sdk.NewAttribute(types.AttributeKeyReceiver, toKavaAddr.String()),
		sdk.NewAttribute(types.AttributeKeyAmount, coinAmount.String()),
	))
}

func (suite *ConversionHooksTestSuite) TestConvert_BalanceChange() {
	suite.Commit()

	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(100)

	balBefore := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.conversionPair.GetAddress(),
		types.NewInternalEVMAddress(suite.key1Addr),
	)
	balModuleBefore := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.conversionPair.GetAddress(),
		types.NewInternalEVMAddress(types.ModuleEVMAddress),
	)
	recipientBalBefore := suite.App.BankKeeper.GetBalance(suite.Ctx, toKavaAddr, suite.conversionPair.Denom)

	// Sends from key1
	res := suite.ConvertToCoin(suite.conversionPair.GetAddress(), toKavaAddr, amount)
	suite.Require().False(res.Failed(), "tx should not fail")

	balAfter := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.conversionPair.GetAddress(),
		types.NewInternalEVMAddress(suite.key1Addr),
	)
	balModuleAfter := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.conversionPair.GetAddress(),
		types.NewInternalEVMAddress(types.ModuleEVMAddress),
	)
	recipientBalAfter := suite.App.BankKeeper.GetBalance(suite.Ctx, toKavaAddr, suite.conversionPair.Denom)

	suite.Require().Equal(
		new(big.Int).Sub(balBefore, amount),
		balAfter,
		"evm initiator balance after convert should decrease by amount",
	)
	suite.Require().Equal(
		new(big.Int).Add(balModuleBefore, amount),
		balModuleAfter,
		"module balance after convert should increase by amount",
	)
	suite.Require().Equal(
		recipientBalBefore.Amount.Add(sdk.NewIntFromBigInt(amount)),
		recipientBalAfter.Amount,
		"kava receiver balance after convert should increase by amount",
	)
}

func (suite *ConversionHooksTestSuite) TestConvert_InsufficientBalance() {
	suite.Commit()

	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	// Bal is 100
	amount := big.NewInt(1000)

	// Sends from key1
	// method is lowercase but event is upper
	data, err := suite.erc20Abi.Pack(
		"convertToCoin",
		common.BytesToAddress(toKavaAddr.Bytes()),
		amount,
	)
	suite.Require().NoError(err)

	_, err = suite.SendTx(suite.conversionPair.GetAddress(), suite.key1Addr, suite.Key1, data)
	suite.Require().Error(err)
	suite.Require().Equal("execution reverted: ERC20: transfer amount exceeds balance", err.Error())
}

func (suite *ConversionHooksTestSuite) TestConvertToCoin_NotEnabled() {
	suite.Commit()
	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(100)

	balBefore := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.invalidConversionPair.GetAddress(),
		types.NewInternalEVMAddress(suite.key1Addr),
	)
	balModuleBefore := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.invalidConversionPair.GetAddress(),
		types.NewInternalEVMAddress(types.ModuleEVMAddress),
	)
	recipientBalBefore := suite.App.BankKeeper.GetBalance(suite.Ctx, toKavaAddr, suite.invalidConversionPair.Denom)

	res := suite.ConvertToCoin(suite.invalidConversionPair.GetAddress(), toKavaAddr, amount)
	suite.Require().False(res.Failed(), "tx should not fail")

	balAfter := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.invalidConversionPair.GetAddress(),
		types.NewInternalEVMAddress(suite.key1Addr),
	)
	balModuleAfter := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.invalidConversionPair.GetAddress(),
		types.NewInternalEVMAddress(types.ModuleEVMAddress),
	)
	recipientBalAfter := suite.App.BankKeeper.GetBalance(suite.Ctx, toKavaAddr, suite.invalidConversionPair.Denom)

	suite.Require().Equal(
		new(big.Int).Sub(balBefore, amount),
		balAfter,
		"evm initiator balance after non-enabled convert should decrease by amount",
	)
	suite.Require().Equal(
		new(big.Int).Add(balModuleBefore, amount),
		balModuleAfter,
		"module balance after non-enabled convert should increase by amount",
	)
	suite.Require().Equal(
		recipientBalBefore.Amount,
		recipientBalAfter.Amount,
		"recipient balance should NOT change for non enabled conversions",
	)

	suite.EventsDoNotContain(suite.GetEvents(), types.EventTypeConvertERC20ToCoin)
}

func (suite *ConversionHooksTestSuite) TestConvertToCoin_NotEnabled_BridgeDisabled() {
	params := suite.Keeper.GetParams(suite.Ctx)
	params.BridgeEnabled = false
	suite.Keeper.SetParams(suite.Ctx, params)

	// Same behavior when bridge is disabled
	suite.TestConvertToCoin_NotEnabled()
}
