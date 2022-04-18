package types_test

import (
	bytes "bytes"
	"encoding/json"
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/yaml"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/kava-labs/kava-bridge/app"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

type ParamsTestSuite struct {
	suite.Suite
}

func (suite *ParamsTestSuite) SetupTest() {
	app.SetSDKConfig()
}

func (suite *ParamsTestSuite) TestParamValidation() {
	type args struct {
		enabledERC20Tokens types.EnabledERC20Tokens
		relayer            sdk.AccAddress
	}
	type errArgs struct {
		expectPass bool
		contains   string
	}
	testCases := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{
			"valid",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						18,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - duplicate token address",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						18,
						testutil.MinWETHWithdrawAmount,
					),
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						18,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: false,
				contains:   "found duplicate enabled ERC20 token address",
			},
		},
		{
			"invalid - zero address",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000000"),
						"Wrapped Ether",
						"WETH",
						18,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: false,
				contains:   "address cannot be zero value",
			},
		},
		{
			"invalid - empty name",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"",
						"WETH",
						18,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: false,
				contains:   "name cannot be empty",
			},
		},
		{
			"invalid - empty symbol",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"",
						18,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: false,
				contains:   "symbol cannot be empty",
			},
		},
		{
			"invalid - zero decimals",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						0,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: false,
				contains:   "decimals cannot be 0",
			},
		},
		{
			"invalid - decimals more than 8 bits",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						256,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: false,
				contains:   "decimals is too large, max 255",
			},
		},
		{
			"invalid - nil address",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						18,
						testutil.MinWETHWithdrawAmount,
					),
				},
				relayer: nil,
			},
			errArgs{
				expectPass: false,
				contains:   "relayer cannot be nil",
			},
		},
		{
			"invalid - zero minimum_withdraw_amount",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						18,
						sdk.ZeroInt(),
					),
				},
				relayer: nil,
			},
			errArgs{
				expectPass: false,
				contains:   "minimum withdraw amount must be positive",
			},
		},
		{
			"invalid - negative minimum_withdraw_amount",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
						"Wrapped Ether",
						"WETH",
						18,
						sdk.NewInt(-1),
					),
				},
				relayer: nil,
			},
			errArgs{
				expectPass: false,
				contains:   "minimum withdraw amount must be positive",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			params := types.NewParams(
				true,
				tc.args.enabledERC20Tokens,
				tc.args.relayer,
				types.DefaultConversionPairs,
			)

			err := params.Validate()
			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}

func (suite *ParamsTestSuite) TestDefault() {
	defaultParams := types.DefaultParams()

	suite.Require().NoError(defaultParams.Validate())

	suite.Require().Empty(defaultParams.EnabledERC20Tokens)
	suite.Require().Equal(types.DefaultEnabledERC20Tokens, defaultParams.EnabledERC20Tokens)
	suite.Require().Equal(types.DefaultRelayer, defaultParams.Relayer)
}

func (suite *ParamsTestSuite) TestUnmarshalJSON() {
	enabledTokens := types.NewEnabledERC20Tokens(
		types.NewEnabledERC20Token(
			testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			"Wrapped Ether",
			"WETH",
			18,
			testutil.MinWETHWithdrawAmount,
		),
		types.NewEnabledERC20Token(
			testutil.MustNewExternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
			"Wrapped Kava",
			"WKAVA",
			6,
			testutil.MinWKavaWithdrawAmount,
		),
		types.NewEnabledERC20Token(
			testutil.MustNewExternalEVMAddressFromString("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			"USD Coin",
			"USDC",
			6,
			testutil.MinUSDCWithdrawAmount,
		),
	)
	enabledTokensJson, err := json.Marshal(enabledTokens)
	suite.Require().NoError(err)

	relayer, err := sdk.AccAddressFromBech32("kava1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea")
	suite.Require().NoError(err)
	relayerJson, err := json.Marshal(relayer)
	suite.Require().NoError(err)

	data := fmt.Sprintf(`{
		"enabled_erc20_tokens": %s,
		"relayer": %s
	}`, string(enabledTokensJson), string(relayerJson))

	var params types.Params
	fmt.Println(data)
	err = json.Unmarshal([]byte(data), &params)
	suite.Require().NoError(err)

	suite.Require().Equal(enabledTokens, params.EnabledERC20Tokens)
	suite.Require().Equal(relayer, params.Relayer)
}

func (suite *ParamsTestSuite) TestMarshalYAML() {
	enabledTokens := types.NewEnabledERC20Tokens(
		types.NewEnabledERC20Token(
			testutil.MustNewExternalEVMAddressFromString("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			"Wrapped Ether",
			"WETH",
			18,
			testutil.MinWETHWithdrawAmount,
		),
		types.NewEnabledERC20Token(
			testutil.MustNewExternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
			"Wrapped Kava",
			"WKAVA",
			6,
			testutil.MinWKavaWithdrawAmount,
		),
		types.NewEnabledERC20Token(
			testutil.MustNewExternalEVMAddressFromString("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			"USD Coin",
			"USDC",
			6,
			testutil.MinUSDCWithdrawAmount,
		),
	)

	relayer, err := sdk.AccAddressFromBech32("kava1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea")
	suite.Require().NoError(err)

	conversionPairs := types.NewConversionPairs(
		types.NewConversionPair(
			testutil.MustNewInternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			"usdc",
		),
	)

	p := types.NewParams(
		true,
		enabledTokens,
		relayer,
		conversionPairs,
	)

	data, err := yaml.Marshal(p)
	suite.Require().NoError(err)

	var params map[string]interface{}
	err = yaml.Unmarshal(data, &params)
	suite.Require().NoError(err)

	_, ok := params["enabled_erc20_tokens"]
	suite.Require().True(ok, "enabled_erc20_tokens should exist in yaml")
	_, ok = params["relayer"]
	suite.Require().True(ok, "relayer should exist in yaml")
	_, ok = params["enabled_conversion_pairs"]
	suite.Require().True(ok, "enabled_conversion_pairs should exist in yaml")
}

func (suite *ParamsTestSuite) TestParamSetPairs_EnabledERC20Tokens() {
	suite.Require().Equal([]byte("EnabledERC20Tokens"), types.KeyEnabledERC20Tokens)
	defaultParams := types.DefaultParams()

	var paramSetPair *paramstypes.ParamSetPair
	for _, pair := range defaultParams.ParamSetPairs() {
		if bytes.Equal(pair.Key, types.KeyEnabledERC20Tokens) {
			paramSetPair = &pair
			break
		}
	}
	suite.Require().NotNil(paramSetPair)

	pairs, ok := paramSetPair.Value.(*types.EnabledERC20Tokens)
	suite.Require().True(ok)
	suite.Require().Equal(pairs, &defaultParams.EnabledERC20Tokens)

	suite.Require().Nil(paramSetPair.ValidatorFn(*pairs))
	suite.Require().EqualError(paramSetPair.ValidatorFn(struct{}{}), "invalid parameter type: struct {}")
}

func (suite *ParamsTestSuite) TestParamSetPairs_Relayer() {
	suite.Require().Equal([]byte("Relayer"), types.KeyRelayer)
	defaultParams := types.DefaultParams()

	var paramSetPair *paramstypes.ParamSetPair
	for _, pair := range defaultParams.ParamSetPairs() {
		if bytes.Equal(pair.Key, types.KeyRelayer) {
			paramSetPair = &pair
			break
		}
	}
	suite.Require().NotNil(paramSetPair)

	pairs, ok := paramSetPair.Value.(*sdk.AccAddress)
	suite.Require().True(ok)
	suite.Require().Equal(pairs, &defaultParams.Relayer)

	suite.Require().Nil(paramSetPair.ValidatorFn(*pairs))
	suite.Require().EqualError(paramSetPair.ValidatorFn(struct{}{}), "invalid parameter type: struct {}")
}

func (suite *ParamsTestSuite) TestParamSetPairs_EnabledConversionPairs() {
	suite.Require().Equal([]byte("EnabledConversionPairs"), types.KeyEnabledConversionPairs)
	defaultParams := types.DefaultParams()

	var paramSetPair *paramstypes.ParamSetPair
	for _, pair := range defaultParams.ParamSetPairs() {
		if bytes.Equal(pair.Key, types.KeyEnabledConversionPairs) {
			paramSetPair = &pair
			break
		}
	}
	suite.Require().NotNil(paramSetPair)

	pairs, ok := paramSetPair.Value.(*types.ConversionPairs)
	suite.Require().True(ok)
	suite.Require().Equal(pairs, &defaultParams.EnabledConversionPairs)

	suite.Require().Nil(paramSetPair.ValidatorFn(*pairs))
	suite.Require().EqualError(paramSetPair.ValidatorFn(struct{}{}), "invalid parameter type: struct {}")
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}
