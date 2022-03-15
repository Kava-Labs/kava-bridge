package bridge_test

import (
	"testing"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/x/bridge"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type genesisTestSuite struct {
	testutil.Suite
}

func (suite *genesisTestSuite) Test_InitGenesis_Validation() {
	type errArgs struct {
		expectPass bool
		panicErr   string
	}
	testStates := []struct {
		name         string
		genesisState types.GenesisState
		errArgs      errArgs
	}{
		{
			"valid",
			types.NewGenesisState(
				types.Params{
					EnabledERC20Tokens: types.EnabledERC20Tokens{
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("C02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
							"Wrapped Ether",
							"WETH",
							18,
						),
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("000000000000000000000000000000000000000A"),
							"Wrapped Kava",
							"WKAVA",
							6,
						),
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("A0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
							"USD Coin",
							"USDC",
							6,
						),
					},
					Relayer: sdk.AccAddress("hi"),
				},
				types.NewERC20BridgePairs(
					types.NewERC20BridgePair(
						testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
						testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
					),
				),
				types.DefaultNextWithdrawSequence,
			),
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - nil relayer",
			types.NewGenesisState(
				types.Params{
					EnabledERC20Tokens: types.EnabledERC20Tokens{
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("C02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
							"Wrapped Ether",
							"WETH",
							18,
						),
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
							"Wrapped Kava",
							"WKAVA",
							6,
						),
					},
					Relayer: nil,
				},
				types.NewERC20BridgePairs(
					types.NewERC20BridgePair(
						testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
						testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
					),
				),
				types.DefaultNextWithdrawSequence,
			),
			errArgs{
				expectPass: false,
				panicErr:   "value from ParamSetPair is invalid: relayer address cannot be nil",
			},
		},
		{
			"invalid - duplicate token address",
			types.NewGenesisState(
				types.Params{
					EnabledERC20Tokens: types.EnabledERC20Tokens{
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("C02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
							"Wrapped Ether",
							"WETH",
							18,
						),
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
							"Wrapped Kava but actually WETH",
							"WKAVA",
							6,
						),
					},
					Relayer: sdk.AccAddress("hi"),
				},
				types.NewERC20BridgePairs(
					types.NewERC20BridgePair(
						testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
						testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
					),
				),
				types.DefaultNextWithdrawSequence,
			),
			errArgs{
				expectPass: false,
				panicErr:   "value from ParamSetPair is invalid: found duplicate enabled ERC20 token address c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
		},
		{
			"invalid - zero token address",
			types.NewGenesisState(
				types.Params{
					EnabledERC20Tokens: types.EnabledERC20Tokens{
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("0000000000000000000000000000000000000000"),
							"Wrapped Ether",
							"WETH",
							18,
						),
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("000000000000000000000000000000000000000A"),
							"Wrapped Kava",
							"WKAVA",
							6,
						),
					},
					Relayer: sdk.AccAddress("hi"),
				},
				types.NewERC20BridgePairs(),
				types.DefaultNextWithdrawSequence,
			),
			errArgs{
				expectPass: false,
				panicErr:   "value from ParamSetPair is invalid: address cannot be zero value 0000000000000000000000000000000000000000",
			},
		},
		{
			"invalid - empty token name",
			types.NewGenesisState(
				types.Params{
					EnabledERC20Tokens: types.EnabledERC20Tokens{
						types.NewEnabledERC20Token(
							testutil.MustNewExternalEVMAddressFromString("C02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
							"",
							"WETH",
							18,
						),
					},
					Relayer: sdk.AccAddress("hi"),
				},
				types.NewERC20BridgePairs(),
				types.DefaultNextWithdrawSequence,
			),
			errArgs{
				expectPass: false,
				panicErr:   "value from ParamSetPair is invalid: name cannot be empty",
			},
		},
	}

	for _, tc := range testStates {
		suite.Run(tc.name, func() {
			if tc.errArgs.expectPass {
				suite.NotPanics(func() {
					bridge.InitGenesis(suite.Ctx, suite.App.BridgeKeeper, suite.App.AccountKeeper, tc.genesisState)
				})
			} else {
				suite.PanicsWithValue(tc.errArgs.panicErr, func() {
					bridge.InitGenesis(suite.Ctx, suite.App.BridgeKeeper, suite.App.AccountKeeper, tc.genesisState)
				}, "expected init genesis to panic with invalid state")
			}
		})
	}
}

func (suite *genesisTestSuite) Test_InitAndExportGenesis() {
	state := types.NewGenesisState(
		types.Params{
			EnabledERC20Tokens: types.EnabledERC20Tokens{
				types.NewEnabledERC20Token(
					testutil.MustNewExternalEVMAddressFromString("C02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
					"Wrapped Ether",
					"WETH",
					18,
				),
				types.NewEnabledERC20Token(
					testutil.MustNewExternalEVMAddressFromString("000000000000000000000000000000000000000A"),
					"Wrapped Kava",
					"WKAVA",
					6,
				),
			},
			Relayer: sdk.AccAddress("hello"),
		},
		types.NewERC20BridgePairs(
			types.NewERC20BridgePair(
				testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
				testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
			),
			types.NewERC20BridgePair(
				testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
				testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
			),
		),
		types.DefaultNextWithdrawSequence,
	)

	bridge.InitGenesis(suite.Ctx, suite.App.BridgeKeeper, suite.App.AccountKeeper, state)
	suite.Equal(state.Params, suite.App.BridgeKeeper.GetParams(suite.Ctx))

	exportedState := bridge.ExportGenesis(suite.Ctx, suite.App.BridgeKeeper, suite.App.AccountKeeper)
	suite.Equal(state, *exportedState)
}

func (suite *genesisTestSuite) Test_Marshall() {
	state := types.NewGenesisState(
		types.Params{
			EnabledERC20Tokens: types.EnabledERC20Tokens{
				types.NewEnabledERC20Token(
					testutil.MustNewExternalEVMAddressFromString("C02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
					"Wrapped Ether",
					"WETH",
					18,
				),
				types.NewEnabledERC20Token(
					testutil.MustNewExternalEVMAddressFromString("000000000000000000000000000000000000000A"),
					"Wrapped Kava",
					"WKAVA",
					6,
				),
			},
			Relayer: sdk.AccAddress("hello"),
		},
		types.NewERC20BridgePairs(
			types.NewERC20BridgePair(
				testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
				testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
			),
			types.NewERC20BridgePair(
				testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
				testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
			),
		),
		types.DefaultNextWithdrawSequence,
	)

	encodingCfg := app.MakeEncodingConfig()
	cdc := encodingCfg.Marshaler

	bz, err := cdc.Marshal(&state)
	suite.Require().NoError(err, "expected genesis state to marshal without error")

	var decodedState types.GenesisState
	err = cdc.Unmarshal(bz, &decodedState)
	suite.Require().NoError(err, "expected genesis state to unmarshal without error")

	suite.Equal(state, decodedState)
}

func (suite *genesisTestSuite) Test_LegacyJSONConversion() {
	state := types.NewGenesisState(
		types.Params{
			EnabledERC20Tokens: types.EnabledERC20Tokens{
				types.NewEnabledERC20Token(
					testutil.MustNewExternalEVMAddressFromString("C02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
					"Wrapped Ether",
					"WETH",
					18,
				),
				types.NewEnabledERC20Token(
					testutil.MustNewExternalEVMAddressFromString("000000000000000000000000000000000000000A"),
					"Wrapped Kava",
					"WKAVA",
					6,
				),
			},
			Relayer: sdk.AccAddress("hello"),
		},
		types.NewERC20BridgePairs(
			types.NewERC20BridgePair(
				testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
				testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
			),
			types.NewERC20BridgePair(
				testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
				testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
			),
		),
		types.DefaultNextWithdrawSequence,
	)

	encodingCfg := app.MakeEncodingConfig()
	cdc := encodingCfg.Marshaler
	legacyCdc := encodingCfg.Amino

	protoJson, err := cdc.MarshalJSON(&state)
	suite.Require().NoError(err, "expected genesis state to marshal amino json without error")

	aminoJson, err := legacyCdc.MarshalJSON(&state)
	suite.Require().NoError(err, "expected genesis state to marshal amino json without error")

	suite.JSONEq(string(protoJson), string(aminoJson), "expected json outputs to be equal")

	var importedState types.GenesisState
	err = cdc.UnmarshalJSON(aminoJson, &importedState)
	suite.Require().NoError(err, "expected amino json to unmarshall to proto without error")

	suite.Equal(state, importedState, "expected genesis state to be equal")
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(genesisTestSuite))
}
