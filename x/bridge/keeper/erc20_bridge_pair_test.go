package keeper_test

import (
	"testing"

	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/suite"
)

type BridgePairTestSuite struct {
	testutil.Suite
}

func TestBridgePairTestSuite(t *testing.T) {
	suite.Run(t, new(BridgePairTestSuite))
}

func (suite *BridgePairTestSuite) TestERC20PairIter() {
	pairs := types.NewERC20BridgePairs(
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000002"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000B"),
		),
		// Already registered in genesis state
		types.NewERC20BridgePair(
			testutil.MustNewExternalEVMAddressFromString("0x0000000000000000000000000000000000000001"),
			testutil.MustNewInternalEVMAddressFromString("0x000000000000000000000000000000000000000A"),
		),
	)

	for _, pair := range pairs {
		suite.App.BridgeKeeper.RegisterERC20BridgePair(suite.Ctx, pair)
	}

	var iterPairs types.ERC20BridgePairs
	suite.App.BridgeKeeper.IterateERC20BridgePairs(suite.Ctx, func(pair types.ERC20BridgePair) bool {
		iterPairs = append(iterPairs, pair)
		return false
	})

	suite.Require().Equal(pairs, iterPairs, "pairs from iterator should match pairs set in store")
}
