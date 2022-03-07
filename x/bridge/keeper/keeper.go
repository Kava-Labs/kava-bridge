package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// Keeper keeper for the bridge module
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramSubspace paramtypes.Subspace
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	evmKeeper     types.EvmKeeper
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramstore paramtypes.Subspace,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	ek types.EvmKeeper,
) Keeper {
	if !paramstore.HasKeyTable() {
		paramstore = paramstore.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		paramSubspace: paramstore,
		bankKeeper:    bk,
		accountKeeper: ak,
		evmKeeper:     ek,
	}
}
