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
			"invalid - unknown address",
			sdk.AccAddress(suite.Key1.PubKey().Address()),
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "signer not authorized for bridge message: unauthorized",
			},
		},
		{
			"invalid - empty signer",
			sdk.AccAddress{},
			suite.Key1,
			errArgs{
				expectPass: false,
				contains:   "signer not authorized for bridge message: unauthorized",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			err := suite.App.BridgeKeeper.IsSignerAuthorized(suite.Ctx, tc.signer)

			if tc.errArgs.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.errArgs.contains)
			}
		})
	}
}
