package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/testutil"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SequenceTestSuite struct {
	testutil.Suite
}

func TestSequenceTestSuite(t *testing.T) {
	suite.Run(t, new(SequenceTestSuite))
}

func (suite *SequenceTestSuite) TestInitialNextWithdrawSequence() {
	readSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(types.DefaultNextWithdrawSequence, readSeq)
	suite.Require().NotZero(readSeq)
}

func (suite *SequenceTestSuite) TestIncrementNextWithdrawSequence() {
	var seq = sdk.NewInt(123456)
	suite.App.BridgeKeeper.SetNextWithdrawSequence(suite.Ctx, seq)

	suite.Require().NoError(suite.App.BridgeKeeper.IncrementNextWithdrawSequence(suite.Ctx))

	// check seq was incremented
	readSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(seq.AddRaw(1), readSeq)
}

func (suite *SequenceTestSuite) TestIncrementNextWithdrawSequence_Wrap() {
	seq := types.MaxWithdrawSequence
	suite.App.BridgeKeeper.SetNextWithdrawSequence(suite.Ctx, seq)

	suite.Require().NoError(suite.App.BridgeKeeper.IncrementNextWithdrawSequence(suite.Ctx))

	// check seq wrapped back to 0
	readSeq, err := suite.App.BridgeKeeper.GetNextWithdrawSequence(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.ZeroInt(), readSeq)
}

func TestWrappingAddInt(t *testing.T) {
	tests := []struct {
		name    string
		i1      sdk.Int
		i2      sdk.Int
		wantSum sdk.Int
	}{
		{
			"1+1=2",
			sdk.OneInt(), sdk.OneInt(), sdk.NewInt(2),
		},
		{
			"0+1=1",
			sdk.ZeroInt(), sdk.OneInt(), sdk.NewInt(1),
		},
		{
			"100+120=220",
			sdk.NewInt(100), sdk.NewInt(120), sdk.NewInt(100 + 120),
		},
		{
			"max+1=0",
			types.MaxWithdrawSequence, sdk.NewInt(1), sdk.NewInt(0),
		},
		{
			"max+2=1",
			types.MaxWithdrawSequence, sdk.NewInt(2), sdk.NewInt(1),
		},
		{
			"max+max=max-1",
			types.MaxWithdrawSequence, types.MaxWithdrawSequence, types.MaxWithdrawSequence.SubRaw(1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sum := keeper.WrappingAddInt(tc.i1, tc.i2)
			require.NotEqual(t, sdk.Int{}, sum, "sum should not be default int value")
			// Use string match here due to non-equal bigint abs type: (big.nat) nil vs abs: (big.nat) {}
			require.Equalf(
				t,
				tc.wantSum.String(),
				sum.String(),
				"sum should match expected value, want %v but got %v", tc.wantSum, sum,
			)
		})
	}
}
