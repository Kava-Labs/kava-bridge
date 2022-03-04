package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"

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

// BridgedInternalEVMAddress puts the bridged address pair into the store.
func (k Keeper) SetBridgedEVMAddress(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
	internalAddress types.InternalEVMAddress,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)
	store.Set(types.GetBridgedERC20Key(externalAddress), internalAddress.Bytes())
}

// GetBridgedInternalEVMAddress gets the internal EVM address mapped to the
// provided ExternalEVMAddress from the store.
func (k Keeper) GetBridgedInternalEVMAddress(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
) (types.InternalEVMAddress, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)
	bz := store.Get(types.GetBridgedERC20Key(externalAddress))
	if bz == nil {
		return types.InternalEVMAddress{}, false
	}

	return types.InternalEVMAddress{
		Address: common.BytesToAddress(bz),
	}, true
}

// DeleteBridgedEVMAddress removes an bridged address pair from the store.
func (k Keeper) DeleteBridgedEVMAddress(ctx sdk.Context, externalAddress types.ExternalEVMAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)
	store.Delete(types.GetBridgedERC20Key(externalAddress))
}
