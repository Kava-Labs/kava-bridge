package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
)

type ParamsTestSuite struct {
	suite.Suite
}

func (suite *ParamsTestSuite) SetupTest() {
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
						"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
						"Wrapped Ether",
						"WETH",
						18,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - hex length mismatch",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756C",
						"Wrapped Ether",
						"WETH",
						18,
					),
				},
				relayer: sdk.AccAddress("1234"),
			},
			errArgs{
				expectPass: false,
				contains:   "address is not a valid hex address",
			},
		},
		{
			"invalid - empty name",
			args{
				enabledERC20Tokens: types.EnabledERC20Tokens{
					types.NewEnabledERC20Token(
						"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
						"",
						"WETH",
						18,
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
						"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
						"Wrapped Ether",
						"",
						18,
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
						"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
						"Wrapped Ether",
						"WETH",
						0,
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
						"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
						"Wrapped Ether",
						"WETH",
						256,
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
						"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
						"Wrapped Ether",
						"WETH",
						18,
					),
				},
				relayer: nil,
			},
			errArgs{
				expectPass: false,
				contains:   "relayer cannot be nil",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			params := types.NewParams(tc.args.enabledERC20Tokens, tc.args.relayer)

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

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}
