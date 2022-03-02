package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/kava-labs/kava-bridge/x/bridge/keeper"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	accountKeeper types.AccountKeeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)

	// Ensure bridge module account is set
	if moduleAcc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); moduleAcc == nil {
		panic("the bridge module account has not been set")
	}
}

// ExportGenesis exports genesis state of the bridge module
func ExportGenesis(ctx sdk.Context, k keeper.Keeper, ak types.AccountKeeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}
