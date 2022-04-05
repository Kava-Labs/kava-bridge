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
		msg := types.NewMsgBridgeEthereumToKava(
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

		_, err = suite.msgServer.BridgeEthereumToKava(sdk.WrapSDKContext(suite.Ctx), &msg)
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

func (suite *GrpcQueryTestSuite) TestQueryERC20BridgePair() {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name        string
		giveRequest types.QueryERC20BridgePairRequest
		wantRes     types.QueryERC20BridgePairResponse
		errArgs     errArgs
	}{
		{
			"valid - external address",
			types.QueryERC20BridgePairRequest{Address: "0x0000000000000000000000000000000000000001"},
			types.QueryERC20BridgePairResponse{
				ERC20BridgePair: types.NewERC20BridgePair(
					testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
					testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
				),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - internal address",
			types.QueryERC20BridgePairRequest{Address: "0x000000000000000000000000000000000000000A"},
			types.QueryERC20BridgePairResponse{
				ERC20BridgePair: types.NewERC20BridgePair(
					testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
					testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
				),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"not found",
			types.QueryERC20BridgePairRequest{Address: "0x0000000000000000000000000000000000000009"},
			types.QueryERC20BridgePairResponse{},
			errArgs{
				expectPass: false,
				contains:   "could not find an ERC20 bridge pair with the provided address",
			},
		},
		{
			"invalid address",
			types.QueryERC20BridgePairRequest{Address: "hi this is invalid"},
			types.QueryERC20BridgePairResponse{},
			errArgs{
				expectPass: false,
				contains:   "invalid hex address",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			res, err := suite.QueryClientBridge.ERC20BridgePair(
				context.Background(),
				&tc.giveRequest,
			)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.wantRes, *res)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func (suite *GrpcQueryTestSuite) TestQueryConversionPairs() {
	pairs, err := suite.QueryClientBridge.ConversionPairs(
		context.Background(),
		&types.QueryConversionPairsRequest{},
	)
	suite.Require().NoError(err)

	params, err := suite.QueryClientBridge.Params(
		context.Background(),
		&types.QueryParamsRequest{},
	)
	suite.Require().NoError(err)

	suite.Require().Equal(params.Params.EnabledConversionPairs, pairs.ConversionPairs)
}

func (suite *GrpcQueryTestSuite) TestQueryConversionPair() {
	type errArgs struct {
		expectPass bool
		contains   string
	}

	tests := []struct {
		name        string
		giveRequest types.QueryConversionPairRequest
		wantRes     types.QueryConversionPairResponse
		errArgs     errArgs
	}{
		{
			"valid - address",
			types.QueryConversionPairRequest{AddressOrDenom: "0x404F9466d758eA33eA84CeBE9E444b06533b369e"},
			types.QueryConversionPairResponse{
				ConversionPair: types.NewConversionPair(
					testutil.MustNewInternalEVMAddressFromString("0x404F9466d758eA33eA84CeBE9E444b06533b369e"),
					"erc20/usdc",
				),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid - denom",
			types.QueryConversionPairRequest{AddressOrDenom: "erc20/usdc"},
			types.QueryConversionPairResponse{
				ConversionPair: types.NewConversionPair(
					testutil.MustNewInternalEVMAddressFromString("0x404F9466d758eA33eA84CeBE9E444b06533b369e"),
					"erc20/usdc",
				),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"valid but not found",
			types.QueryConversionPairRequest{AddressOrDenom: "erc20/sdf"},
			types.QueryConversionPairResponse{},
			errArgs{
				expectPass: false,
				contains:   "could not find bridge pair with provided address or denom",
			},
		},
		{
			"invalid address",
			types.QueryConversionPairRequest{AddressOrDenom: "0x4"},
			types.QueryConversionPairResponse{},
			errArgs{
				expectPass: false,
				contains:   "invalid hex address or denom",
			},
		},
		{
			"invalid denom",
			types.QueryConversionPairRequest{AddressOrDenom: "this is not a valid denom"},
			types.QueryConversionPairResponse{},
			errArgs{
				expectPass: false,
				contains:   "invalid hex address or denom",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			res, err := suite.QueryClientBridge.ConversionPair(
				context.Background(),
				&tc.giveRequest,
			)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.wantRes, *res)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}
