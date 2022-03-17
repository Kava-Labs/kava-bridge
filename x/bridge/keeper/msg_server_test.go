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

func (suite *MsgServerSuite) TestMsg() {
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

				suite.TypedEventsContains(suite.GetEvents(), &types.EventBridgeEthereumToKava{
					Relayer:              msg.Relayer,
					EthereumErc20Address: msg.EthereumERC20Address,
					Receiver:             receiver.String(),
					Amount:               amount.String(),
					Sequence:             msg.Sequence.String(),
				})
			}
		})
	}
}
