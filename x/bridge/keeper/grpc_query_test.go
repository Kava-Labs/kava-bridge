package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type GrpcQueryTestSuite struct {
	testutil.Suite

	msgServer types.MsgServer
}

func (suite *GrpcQueryTestSuite) SetupTest() {
	suite.Suite.SetupTest()
	suite.msgServer = keeper.NewMsgServerImpl(suite.App.BridgeKeeper)
}

func TestGrpcQueryTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcQueryTestSuite))
}

func (suite *GrpcQueryTestSuite) TestQueryERC20BridgePairs() {
	extContracts := []string{
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"0x000000000000000000000000000000000000000A",
		"A0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
	}

	var internalContracts []string

	for _, contractAddr := range extContracts {
		msg := types.NewMsgBridgeERC20FromEthereum(
			suite.RelayerAddress.String(),
			contractAddr,
			sdk.NewInt(10),
			"0x0000000000000000000000000000000000000001",
			sdk.NewInt(1),
		)

		receiver := types.InternalEVMAddress{}
		err := receiver.UnmarshalText([]byte(msg.Receiver))
		suite.Require().NoError(err)

		externalAddress := types.ExternalEVMAddress{}
		err = externalAddress.UnmarshalText([]byte(msg.EthereumERC20Address))
		suite.Require().NoError(err)

		_, err = suite.msgServer.BridgeERC20FromEthereum(sdk.WrapSDKContext(suite.Ctx), &msg)
		suite.Require().NoError(err)

		contractAddr, found := suite.App.BridgeKeeper.GetInternalERC20Address(suite.Ctx, externalAddress)
		suite.Require().True(found)

		internalContracts = append(internalContracts, contractAddr.String())
	}

	queriedBridgedERC20Pairs, err := suite.QueryClientBridge.ERC20BridgePairs(
		context.Background(),
		&types.QueryERC20BridgePairsRequest{},
	)
	suite.Require().NoError(err)

	suite.Require().Len(queriedBridgedERC20Pairs, len(extContracts), "queried erc20 pairs should match len of bridged contracts")

	suite.Require().Equal(extContracts, queriedBridgedERC20Pairs)
}
