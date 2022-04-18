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
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

type EVMHooksTestSuite struct {
	testutil.Suite

	msgServer types.MsgServer
	key1Addr  common.Address
	erc20Abi  abi.ABI
	pair      types.ERC20BridgePair
}

func TestEVMHooksTestSuite(t *testing.T) {
	suite.Run(t, new(EVMHooksTestSuite))
}

func (suite *EVMHooksTestSuite) SetupTest() {
	suite.Suite.SetupTest()

	suite.msgServer = keeper.NewMsgServerImpl(suite.App.BridgeKeeper)

	suite.erc20Abi = contract.ERC20MintableBurnableContract.ABI
	externalWethAddr := testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")

	// Bridge an asset to deploy the ERC20 asset and update store with pair
	suite.key1Addr = common.BytesToAddress(suite.Key1.PubKey().Address())
	suite.submitBridgeERC20Msg(
		externalWethAddr,
		testutil.MinWETHWithdrawAmount.Mul(sdk.NewInt(10)),
		suite.key1Addr,
	)

	var found bool
	suite.pair, found = suite.App.BridgeKeeper.GetBridgePairFromExternal(suite.Ctx, externalWethAddr)
	suite.Require().True(found, "bridge pair must exist after bridge")
}

func (suite *EVMHooksTestSuite) TestHooksSet() {
	suite.Require().PanicsWithValue("cannot set evm hooks twice", func() {
		suite.App.EvmKeeper.SetHooks(suite.App.BridgeKeeper.WithdrawHooks())
	})
}

func (suite *EVMHooksTestSuite) submitBridgeERC20Msg(
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

func (suite *EVMHooksTestSuite) Withdraw(
	contractAddr types.InternalEVMAddress,
	toAddr common.Address,
	amount *big.Int,
) *evmtypes.MsgEthereumTxResponse {
	// method is lowercase but event is upper
	data, err := suite.erc20Abi.Pack(
		"withdraw",
		toAddr,
		amount,
	)
	suite.Require().NoError(err)

	res, err := suite.SendTx(contractAddr, suite.key1Addr, suite.Key1, data)
	suite.Require().NoError(err)
	suite.Require().False(res.Failed(), "evm tx should not fail %v", res)

	return res
}

func (suite *EVMHooksTestSuite) TestERC20WithdrawUnpack() {
	withdrawAmount := testutil.MinWETHWithdrawAmount.BigInt()
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())

	// Send TX
	res := suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)

	containsWithdrawEvent := false

	for _, log := range res.Logs {
		eventID := log.Topics[0]

		event, err := suite.erc20Abi.EventByID(common.HexToHash(eventID))
		if err != nil {
			// invalid event for ERC20
			continue
		}

		if event.Name != types.ContractEventTypeWithdraw {
			continue
		}

		containsWithdrawEvent = true

		withdrawEvent, err := suite.erc20Abi.Unpack(types.ContractEventTypeWithdraw, log.Data)
		suite.Require().NoError(err)

		suite.Require().Len(withdrawEvent, 1, "withdraw event data should only have 1 item for amount")

		loggedAmount, ok := withdrawEvent[0].(*big.Int)
		suite.Require().True(ok, "withdraw event data should be *big.Int")
		suite.Require().Equal(withdrawAmount, loggedAmount)

		// 3 topics:
		// 0: Keccak-256 hash of Withdraw(address,address,uint256)
		// 1: address indexed sender
		// 2: address indexed toAddr
		suite.Require().Len(log.Topics, 3, "withdraw event should have 3 topics")

		// log.Topics is padded to 32 bytes, addresses are 20 bytes.
		// common.HexToAddress handles this, crops from left.
		senderAddr := common.HexToAddress(log.Topics[1])
		suite.Require().Equal(suite.key1Addr, senderAddr)

		logToAddr := common.HexToAddress(log.Topics[2])
		suite.Require().Equal(withdrawToAddr, logToAddr)
	}

	suite.Require().True(containsWithdrawEvent, "tx should contain Withdraw event")
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_BalanceChange() {
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := testutil.MinWETHWithdrawAmount.MulRaw(2).BigInt()

	balBefore := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.pair.GetInternalAddress(),
		types.NewInternalEVMAddress(suite.key1Addr),
	)

	// Send Withdraw TX
	_ = suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)

	balAfter := suite.GetERC20BalanceOf(
		contract.ERC20MintableBurnableContract.ABI,
		suite.pair.GetInternalAddress(),
		types.NewInternalEVMAddress(suite.key1Addr),
	)

	suite.Require().Equal(
		new(big.Int).Sub(balBefore, withdrawAmount),
		balAfter,
		"balance after withdraw should burn withdraw amount",
	)
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_MinimumWithdraw() {
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := testutil.MinWETHWithdrawAmount

	// Send Withdraw TX
	res := suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount.BigInt())

	suite.Require().False(res.Failed(), "minimum withdraw should not fail")
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_UnderMinimumWithdraw() {
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := testutil.MinWETHWithdrawAmount.Sub(sdk.OneInt())

	// Send Withdraw TX
	data, err := suite.erc20Abi.Pack(
		"withdraw",
		withdrawToAddr,
		withdrawAmount.BigInt(),
	)
	suite.Require().NoError(err)

	res, err := suite.SendTx(suite.pair.GetInternalAddress(), suite.key1Addr, suite.Key1, data)
	suite.Require().NoError(err)

	suite.Require().True(res.Failed(), "should fail for withdraws under minimum withdraw amount")
	suite.Require().Equal(evmtypes.ErrPostTxProcessing.Error(), res.VmError)
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_AfterERC20Disabled() {
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := testutil.MinWETHWithdrawAmount.BigInt()

	// Send Withdraw TX
	res := suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)
	suite.Require().False(res.Failed(), "minimum withdraw should not fail")

	// Remove all enabled ERC20 tokens
	params := suite.Keeper.GetParams(suite.Ctx)
	params.EnabledERC20Tokens = types.DefaultEnabledERC20Tokens
	suite.Keeper.SetParams(suite.Ctx, params)

	// Withdraw with ERC20 disabled should still fail
	data, err := suite.erc20Abi.Pack(
		"withdraw",
		withdrawToAddr,
		withdrawAmount,
	)
	suite.Require().NoError(err)

	res, err = suite.SendTx(suite.pair.GetInternalAddress(), suite.key1Addr, suite.Key1, data)
	suite.Require().NoError(err)

	suite.Require().True(res.Failed(), "withdraw with disabled ERC20 but with bridge pair state should fail")
	suite.Require().Equal(evmtypes.ErrPostTxProcessing.Error(), res.VmError)
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_SequenceIncrement() {
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := testutil.MinWETHWithdrawAmount.BigInt()

	beforeWithdrawSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)

	// Send Withdraw TX
	_ = suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)

	afterWithdrawSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(
		beforeWithdrawSeq.Add(sdk.OneInt()),
		afterWithdrawSeq,
		"next withdraw sequence should be incremented by 1",
	)
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_SequenceWrap() {
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := testutil.MinWETHWithdrawAmount.BigInt()

	// Set next sequence to max int
	suite.App.BridgeKeeper.SetNextWithdrawSequence(suite.Ctx, types.MaxWithdrawSequence)

	// Check it was set
	beforeWithdrawSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(
		types.MaxWithdrawSequence,
		beforeWithdrawSeq,
		"next withdraw sequence should be max int",
	)

	// Send Withdraw TX
	_ = suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)

	afterWithdrawSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(
		sdk.ZeroInt(),
		afterWithdrawSeq,
		"next withdraw sequence should wrap around to 0",
	)
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_EmitsEvent() {
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := testutil.MinWETHWithdrawAmount.BigInt()

	// Send Withdraw TX
	_ = suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)

	suite.EventsContains(suite.GetEvents(),
		sdk.NewEvent(
			types.EventTypeBridgeKavaToEthereum,
			sdk.NewAttribute(types.AttributeKeyEthereumERC20Address, suite.pair.GetExternalAddress().String()),
			sdk.NewAttribute(types.AttributeKeyKavaERC20Address, suite.pair.GetInternalAddress().String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, withdrawToAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, withdrawAmount.String()),
			sdk.NewAttribute(types.AttributeKeySequence, "1"),
		))

	// Second withdraw tx
	_ = suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)

	// Second one has incremented sequence
	suite.EventsContains(suite.GetEvents(),
		sdk.NewEvent(
			types.EventTypeBridgeKavaToEthereum,
			sdk.NewAttribute(types.AttributeKeyEthereumERC20Address, suite.pair.GetExternalAddress().String()),
			sdk.NewAttribute(types.AttributeKeyKavaERC20Address, suite.pair.GetInternalAddress().String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, withdrawToAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, withdrawAmount.String()),
			sdk.NewAttribute(types.AttributeKeySequence, "2"),
		))
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_IgnoreUnregisteredERC20() {
	token := types.NewEnabledERC20Token(
		testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		"Token Anyone Can Deploy",
		"TACD",
		18,
		sdk.NewInt(10_000_000_000_000_000),
	)

	// We are using keeper methods to deploy / mint ERC20 balance, but this can
	// be done by just regular EVM calls, ie. deploying the same
	// mintable/burnable ERC20 from some unknown account. We are only testing
	// emitted events are ignored by the hook if the contracts are not
	// registered in the state.
	unregisteredContractAddr, err := suite.App.BridgeKeeper.DeployMintableERC20Contract(suite.Ctx, token)
	suite.Require().NoError(err)
	suite.Require().Greater(len(unregisteredContractAddr.Address), 0)

	mintAmount := big.NewInt(10)
	err = suite.App.BridgeKeeper.MintERC20(suite.Ctx, unregisteredContractAddr, types.NewInternalEVMAddress(suite.key1Addr), mintAmount)
	suite.Require().NoError(err)

	// Withdraw / burn funds
	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := big.NewInt(10)

	// Send Withdraw TX to the erc20 contract that is not a registered pair
	_ = suite.Withdraw(unregisteredContractAddr, withdrawToAddr, withdrawAmount)

	suite.EventsDoNotContain(suite.GetEvents(), types.EventTypeBridgeKavaToEthereum)
}

func (suite *EVMHooksTestSuite) TestERC20Withdraw_BridgeDisabled() {
	// Disable bridge
	params := suite.Keeper.GetParams(suite.Ctx)
	params.BridgeEnabled = false
	suite.Keeper.SetParams(suite.Ctx, params)

	toKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	withdrawToAddr := common.BytesToAddress(toKey.PubKey().Address())
	withdrawAmount := big.NewInt(10)

	// Send Withdraw TX
	res := suite.Withdraw(suite.pair.GetInternalAddress(), withdrawToAddr, withdrawAmount)

	suite.Require().True(res.Failed())
}
