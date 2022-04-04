package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
)

type MsgServerSuite struct {
	testutil.Suite

	msgServer types.MsgServer
}

func (suite *MsgServerSuite) SetupTest() {
	suite.Suite.SetupTest()
	suite.msgServer = keeper.NewMsgServerImpl(suite.App.BridgeKeeper)
}

func TestMsgServerSuite(t *testing.T) {
	suite.Run(t, new(MsgServerSuite))
}

func (suite *MsgServerSuite) TestBridgeEthereumToKava() {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		msg     types.MsgBridgeEthereumToKava
		errArgs errArgs
	}{
		{
			"valid - signer matches relayer in params",
			types.NewMsgBridgeEthereumToKava(
				suite.RelayerAddress.String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - signer mismatch",
			types.NewMsgBridgeEthereumToKava(
				sdk.AccAddress(suite.Key1.PubKey().Address()).String(),
				"0x000000000000000000000000000000000000000A",
				sdk.NewInt(10),
				"0x0000000000000000000000000000000000000001",
				sdk.NewInt(0),
			),
			errArgs{
				expectPass: false,
				contains:   "signer not authorized for bridge message: unauthorized",
			},
		},
		{
			"invalid - token not enabled",
			types.NewMsgBridgeEthereumToKava(
				suite.RelayerAddress.String(),
				"0x000000000000000000000000000000000000000B",
				sdk.NewInt(10),
				"0x0000000000000000000000000000000000000001",
				sdk.NewInt(0),
			),
			errArgs{
				expectPass: false,
				contains:   types.ErrERC20NotEnabled.Error(),
			},
		},
		{
			"invalid - malformed external address",
			types.NewMsgBridgeEthereumToKava(
				suite.RelayerAddress.String(),
				"hi",
				sdk.NewInt(10),
				"0x0000000000000000000000000000000000000001",
				sdk.NewInt(0),
			),
			errArgs{
				expectPass: false,
				contains:   "invalid EthereumERC20Address: string is not a hex address",
			},
		},
		{
			"invalid - malformed internal receiver address",
			types.NewMsgBridgeEthereumToKava(
				suite.RelayerAddress.String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(10),
				"hi",
				sdk.NewInt(0),
			),
			errArgs{
				expectPass: false,
				contains:   "invalid Receiver address: string is not a hex address",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.BridgeEthereumToKava(sdk.WrapSDKContext(suite.Ctx), &tc.msg)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func (suite *MsgServerSuite) TestMint() {
	extContractAddr := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

	tests := []struct {
		name        string
		receiver    string
		mintAmounts []sdk.Int
	}{
		{
			"valid - mint once",
			"0x0000000000000000000000000000000000000001",
			[]sdk.Int{
				sdk.NewInt(10),
			},
		},
		{
			"valid - mint multiple times",
			"0x0000000000000000000000000000000000000002",
			[]sdk.Int{
				sdk.NewInt(10),
				sdk.NewInt(13),
				sdk.NewInt(15),
				sdk.NewInt(18),
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			total := big.NewInt(0)

			for i, amount := range tc.mintAmounts {
				total = total.Add(total, amount.BigInt())
				msg := types.NewMsgBridgeEthereumToKava(
					suite.RelayerAddress.String(),
					extContractAddr,
					amount,
					tc.receiver,
					// sequence doesn't actually matter here, but we use index
					// just as a way to check, later the same sequence is re-emitted
					sdk.NewInt(int64(i)),
				)

				receiver := types.InternalEVMAddress{}
				err := receiver.UnmarshalText([]byte(msg.Receiver))
				suite.Require().NoError(err)

				externalAddress := types.ExternalEVMAddress{}
				err = externalAddress.UnmarshalText([]byte(msg.EthereumERC20Address))
				suite.Require().NoError(err)

				_, err = suite.msgServer.BridgeEthereumToKava(sdk.WrapSDKContext(suite.Ctx), &msg)
				suite.Require().NoError(err)

				pair, found := suite.App.BridgeKeeper.GetBridgePairFromExternal(suite.Ctx, externalAddress)
				suite.Require().True(found)

				bal := suite.GetERC20BalanceOf(
					contract.ERC20MintableBurnableContract.ABI,
					pair.GetInternalAddress(),
					receiver,
				)

				suite.Require().Equal(total, bal, "balance should match amount minted so far")

				suite.EventsContains(suite.GetEvents(),
					sdk.NewEvent(
						types.EventTypeBridgeEthereumToKava,
						sdk.NewAttribute(types.AttributeKeyRelayer, msg.Relayer),
						sdk.NewAttribute(types.AttributeKeyEthereumERC20Address, msg.EthereumERC20Address),
						sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
						sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
						sdk.NewAttribute(types.AttributeKeySequence, msg.Sequence.String()),
					))
			}
		})
	}
}

func (suite *MsgServerSuite) TestConvertCoinToERC20() {
	invoker, err := sdk.AccAddressFromBech32("kava123fxg0l602etulhhcdm0vt7l57qya5wjcrwhzz")
	suite.Require().NoError(err)

	err = suite.App.FundAccount(suite.Ctx, invoker, sdk.NewCoins(sdk.NewCoin("erc20/usdc", sdk.NewInt(10000))))
	suite.Require().NoError(err)

	contractAddr := suite.DeployERC20()

	pair := types.NewConversionPair(
		contractAddr,
		"erc20/usdc",
	)

	// Module account should have starting balance
	pairStartingBal := big.NewInt(10000)
	err = suite.App.BridgeKeeper.MintERC20(
		suite.Ctx,
		pair.GetAddress(), // contractAddr
		types.NewInternalEVMAddress(types.ModuleEVMAddress), //receiver
		pairStartingBal,
	)
	suite.Require().NoError(err)

	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name    string
		msg     types.MsgConvertCoinToERC20
		errArgs errArgs
	}{
		{
			"valid",
			types.NewMsgConvertCoinToERC20(
				invoker.String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewCoin("erc20/usdc", sdk.NewInt(1234)),
			),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - odd length hex address",
			types.NewMsgConvertCoinToERC20(
				invoker.String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc",
				sdk.NewCoin("erc20/usdc", sdk.NewInt(1234)),
			),
			errArgs{
				expectPass: false,
				contains:   "invalid Receiver address: string is not a hex address",
			},
		},
		// Amount coin is not validated by msg_server, but on msg itself
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.ConvertCoinToERC20(sdk.WrapSDKContext(suite.Ctx), &tc.msg)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)

				bal := suite.GetERC20BalanceOf(
					contract.ERC20MintableBurnableContract.ABI,
					pair.GetAddress(),
					testutil.MustNewInternalEVMAddressFromString(tc.msg.Receiver),
				)

				suite.Require().Equal(tc.msg.Amount.Amount.BigInt(), bal, "balance should match converted amount")

				// msg server event
				suite.EventsContains(suite.GetEvents(),
					sdk.NewEvent(
						sdk.EventTypeMessage,
						sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
						sdk.NewAttribute(sdk.AttributeKeySender, tc.msg.Initiator),
					))

				// keeper event
				suite.EventsContains(suite.GetEvents(),
					sdk.NewEvent(
						types.EventTypeConvertCoinToERC20,
						sdk.NewAttribute(types.AttributeKeyInitiator, tc.msg.Initiator),
						sdk.NewAttribute(types.AttributeKeyReceiver, tc.msg.Receiver),
						sdk.NewAttribute(types.AttributeKeyERC20Address, pair.GetAddress().String()),
						sdk.NewAttribute(types.AttributeKeyAmount, tc.msg.Amount.String()),
					))
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}
