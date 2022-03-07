package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
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

func (suite *MsgServerSuite) TestPermissioned() {
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
		signer  sdk.AccAddress
		key     *ethsecp256k1.PrivKey
		errArgs errArgs
	}{
		{
			"valid - signer matches relayer in params",
			relayerAddr,
			suite.RelayerKey,
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid",
			sdk.AccAddress(suite.Key1.PubKey().Address()),
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "signer not authorized for bridge message",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			msg := types.NewMsgBridgeERC20FromEthereum(
				tc.signer.String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			)

			_, err := suite.msgServer.BridgeERC20FromEthereum(sdk.WrapSDKContext(suite.Ctx), &msg)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}
