package keeper_test

import (
	"fmt"
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

	msgServer      types.MsgServer
	key1Addr       common.Address
	erc20Abi       abi.ABI
	conversionPair types.ConversionPair
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
	suite.conversionPair = types.NewConversionPair(bridgePair.GetInternalAddress(), "erc20/weth")
	params := suite.App.BridgeKeeper.GetParams(suite.Ctx)
	params.EnabledConversionPairs = append(params.EnabledConversionPairs, suite.conversionPair)

	suite.App.BridgeKeeper.SetParams(suite.Ctx, params)
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
	// Fixes out of gas error
	suite.Commit()

	// LEFT zero padded
	var toKavaAddrBytes32 [32]byte
	copy(toKavaAddrBytes32[32-20:], toKavaAddr)

	// method is lowercase but event is upper
	data, err := suite.erc20Abi.Pack(
		"convertToCoin",
		toKavaAddrBytes32,
		amount,
	)
	suite.Require().NoError(err)

	res := suite.SendTx(contractAddr, suite.key1Addr, suite.Key1, data)
	suite.Require().False(res.Failed(), "evm tx should not fail %v", res)

	return res
}

func (suite *ConversionHooksTestSuite) TestConvertToCoin() {
	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(100)

	_ = suite.ConvertToCoin(suite.conversionPair.GetAddress(), toKavaAddr, amount)
}

func (suite *ConversionHooksTestSuite) TestConvert_BalanceChange() {
	suite.Commit()

	toKavaAddr := sdk.AccAddress(suite.Key2.PubKey().Address())
	amount := big.NewInt(100)

	fmt.Printf("orig raw bytes: %b\nbech32: %v (bytes: %b)\n", toKavaAddr, toKavaAddr.String(), toKavaAddr.Bytes())

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
	_ = suite.ConvertToCoin(suite.conversionPair.GetAddress(), toKavaAddr, amount)

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
