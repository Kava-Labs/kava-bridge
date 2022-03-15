package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"

	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
)

type RelayerTestSuite struct {
	testutil.Suite
}

func TestRelayerTestSuite(t *testing.T) {
	suite.Run(t, new(RelayerTestSuite))
}

func (suite *RelayerTestSuite) TestPermission() {
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
		signers []sdk.AccAddress
		key     *ethsecp256k1.PrivKey
		errArgs errArgs
	}{
		{
			"valid - signer matches relayer in params",
			[]sdk.AccAddress{relayerAddr},
			suite.RelayerKey,
			errArgs{
				expectPass: true,
			},
		},
		{
			"invalid - multiple signers even if permissioned",
			[]sdk.AccAddress{
				relayerAddr,
				sdk.AccAddress(suite.Key1.PubKey().Address()),
			},
			suite.RelayerKey,
			errArgs{
				expectPass: false,
				contains:   "invalid number of signers",
			},
		},
		{
			"invalid - single unknown address",
			[]sdk.AccAddress{sdk.AccAddress(suite.Key1.PubKey().Address())},
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "signer not authorized for bridge message",
			},
		},
		{
			"invalid - multiple unknown addresses",
			[]sdk.AccAddress{
				sdk.AccAddress(suite.Key1.PubKey().Address()),
				sdk.AccAddress(suite.Key2.PubKey().Address()),
			},
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "invalid number of signers",
			},
		},
		{
			"invalid - empty signers",
			[]sdk.AccAddress{},
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "invalid number of signers",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			err := suite.App.BridgeKeeper.IsSignerAuthorized(suite.Ctx, tc.signers)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}
