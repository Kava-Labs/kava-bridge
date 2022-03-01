package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type ERC20TestSuite struct {
	suite.Suite

	App          app.TestApp
	bridgeKeeper keeper.Keeper
	Ctx          sdk.Context
}

func (suite *ERC20TestSuite) SetupTest() {
	suite.App = app.NewTestApp()
	suite.Ctx = suite.App.NewContext(true, tmproto.Header{Height: 1, Time: tmtime.Now()})

	suite.bridgeKeeper = suite.App.GetBridgeKeeper()
	_, addrs := app.GeneratePrivKeyAddressPairs(1)

	bridgeGs := types.NewGenesisState(types.NewParams(
		types.EnabledERC20Tokens{
			types.NewEnabledERC20Token(
				"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				"Wrapped ETH",
				"WETH",
				18,
			),
		},
		addrs[0],
	))

	evmGs := evmtypes.NewGenesisState(
		evmtypes.NewParams("ukava", true, true, evmtypes.DefaultChainConfig()),
		nil,
	)

	evmGsBytes := app.GenesisState{evmtypes.ModuleName: suite.App.AppCodec().MustMarshalJSON(evmGs)}

	moduleGs := suite.App.AppCodec().MustMarshalJSON(&bridgeGs)
	gs := app.GenesisState{types.ModuleName: moduleGs}
	suite.App.InitializeFromGenesisStates(evmGsBytes, gs)
}

func TestERC20TestSuite(t *testing.T) {
	suite.Run(t, new(ERC20TestSuite))
}

func (suite *ERC20TestSuite) TestDeployERC20() {
	token := types.NewEnabledERC20Token(
		"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		"Wrapped ETH",
		"WETH",
		18,
	)

	addr, err := suite.bridgeKeeper.DeployMintableERC20Contract(suite.Ctx, token)
	suite.Require().NoError(err)
	suite.Require().Greater(0, len(addr))
}
