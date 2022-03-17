package keeper_test

import (
	"context"
	"strings"
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
	// Fetch initial pairs since there's some already set in genesis
	initialBridgedERC20Pairs, err := suite.QueryClientBridge.ERC20BridgePairs(
		context.Background(),
		&types.QueryERC20BridgePairsRequest{},
	)

	suite.Require().NoError(err)
	extContracts := []string{
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		"0x000000000000000000000000000000000000000a",
		"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
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

		pair, found := suite.App.BridgeKeeper.GetBridgePairFromExternal(suite.Ctx, externalAddress)
		suite.Require().True(found)

		internalContracts = append(internalContracts, strings.ToLower(pair.GetInternalAddress().String()))
	}

	suite.Commit()

	queriedBridgedERC20Pairs, err := suite.QueryClientBridge.ERC20BridgePairs(
		context.Background(),
		&types.QueryERC20BridgePairsRequest{},
	)
	suite.Require().NoError(err)

	suite.Require().Lenf(
		queriedBridgedERC20Pairs.ERC20BridgePairs,
		len(extContracts)+len(initialBridgedERC20Pairs.ERC20BridgePairs),
		"queried erc20 pairs should match len of bridged contracts: %v",
		queriedBridgedERC20Pairs.ERC20BridgePairs,
	)

	var queriedExtAddrs []string
	var queriedIntAddrs []string

	for _, pair := range queriedBridgedERC20Pairs.ERC20BridgePairs {
		// ToLower since String() returns a checksum address which we don't care about
		queriedExtAddrs = append(queriedExtAddrs, strings.ToLower(pair.GetExternalAddress().String()))
		queriedIntAddrs = append(queriedIntAddrs, strings.ToLower(pair.GetInternalAddress().String()))
	}

	for _, addr := range extContracts {
		suite.Require().Containsf(queriedExtAddrs, addr, "queried pairs should contain new external addr %v", addr)
	}

	for _, addr := range internalContracts {
		suite.Require().Containsf(queriedIntAddrs, addr, "queried pairs should contain new internal addr %v", addr)
	}
}
