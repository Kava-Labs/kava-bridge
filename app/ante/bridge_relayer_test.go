package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/app/ante"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	bridgetypes "github.com/kava-labs/kava-bridge/x/bridge/types"
)

var (
	_ sdk.AnteHandler = (&MockAnteHandler{}).AnteHandle
)

type MockAnteHandler struct {
	WasCalled bool
}

func (mah *MockAnteHandler) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	mah.WasCalled = true
	return ctx, nil
}

type bridgeRelayerAnteTestSuite struct {
	testutil.Suite
}

func TestBridgeRelayerAnteTestSuite(t *testing.T) {
	suite.Run(t, new(bridgeRelayerAnteTestSuite))
}

func (suite *bridgeRelayerAnteTestSuite) TestBridgeAnte_OnlyParmsRelayer() {
	txConfig := app.MakeEncodingConfig().TxConfig

	decorator := ante.NewBridgeRelayerDecorator(suite.App.BridgeKeeper)

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
			bridgeMsg := bridgetypes.NewMsgBridgeERC20FromEthereum(
				tc.signer.String(),
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				sdk.NewInt(1234),
				"0x4A59E9DDB116A04C5D40082D67C738D5C56DF124",
				sdk.NewInt(1),
			)

			tx, err := helpers.GenTx(
				txConfig,
				[]sdk.Msg{&bridgeMsg},
				sdk.NewCoins(), // no fee
				helpers.DefaultGenTxGas,
				"testing-chain-id",
				[]uint64{0},
				[]uint64{0},
				tc.key,
			)
			suite.Require().NoError(err)
			mmd := MockAnteHandler{}
			ctx := suite.Ctx.WithIsCheckTx(false) // run as it would be during block update ('DeliverTx'), not just checking entry to mempool

			_, err = decorator.AnteHandle(ctx, tx, false, mmd.AnteHandle)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
				suite.Require().True(mmd.WasCalled, "should continue to next ante handler")
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
				suite.Require().False(mmd.WasCalled, "should not continue tx")
			}
		})
	}
}
