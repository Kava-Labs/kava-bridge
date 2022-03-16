package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// RegisterBridgePair puts the erc20 bridge pair in the store with all of
// its corresponding external/internal to ID mappings.
func (k Keeper) RegisterBridgePair(
	ctx sdk.Context,
	pair types.ERC20BridgePair,
) {
	id := pair.GetID()

	k.setBridgePair(ctx, pair)
	k.setPairIDFromExternal(ctx, pair.GetExternalAddress(), id)
	k.setPairIDFromInternal(ctx, pair.GetInternalAddress(), id)
}

// setBridgePair puts the bridged address pair into the store.
func (k Keeper) setBridgePair(
	ctx sdk.Context,
	pair types.ERC20BridgePair,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20PairKeyPrefix)

	key := pair.GetID()
	bz := k.cdc.MustMarshal(&pair)

	store.Set(key, bz)
}

// GetBridgePair returns the ERC20 bridge pair with the provided pair ID
// from the store.
func (k Keeper) GetBridgePair(
	ctx sdk.Context,
	id []byte,
) (types.ERC20BridgePair, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20PairKeyPrefix)

	bz := store.Get(id)
	if bz == nil {
		return types.ERC20BridgePair{}, false
	}

	var pair types.ERC20BridgePair
	k.cdc.MustUnmarshal(bz, &pair)

	return pair, true
}

// IterateBridgePairs provides an iterator over all stored ERC20 bridge
// pairs. For each pair, cb will be called. If cb returns true, the iterator
// will close and stop.
func (k Keeper) IterateBridgePairs(
	ctx sdk.Context,
	cb func(pair types.ERC20BridgePair) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BridgedERC20PairKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pair types.ERC20BridgePair
		k.cdc.MustUnmarshal(iterator.Value(), &pair)

		if cb(pair) {
			break
		}
	}
}

// GetBridgePairFromExternal gets ERC20 bridge pair with the provided
// ExternalEVMAddress from the store.
func (k Keeper) GetBridgePairFromExternal(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
) (types.ERC20BridgePair, bool) {
	id, found := k.GetPairIDFromExternal(ctx, externalAddress)
	if !found {
		return types.ERC20BridgePair{}, false
	}

	return k.GetBridgePair(ctx, id)
}

// GetBridgePairFromInternal gets ERC20 bridge pair with the provided
// InternalEVMAddress from the store.
func (k Keeper) GetBridgePairFromInternal(
	ctx sdk.Context,
	internalAddress types.InternalEVMAddress,
) (types.ERC20BridgePair, bool) {
	id, found := k.GetPairIDFromInternal(ctx, internalAddress)
	if !found {
		return types.ERC20BridgePair{}, false
	}

	return k.GetBridgePair(ctx, id)
}

// -----------------------------------------------------------------------------
// External Address -> ERC20BridgePair ID

// GetPairIDFromExternal gets the erc20 bridge pair id from the given external address.
func (k Keeper) GetPairIDFromExternal(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
) ([]byte, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20PairByExternalKeyPrefix)
	id := store.Get(externalAddress.Bytes())
	return id, id != nil
}

// setPairIDFromExternal sets the erc20 bridge pair id for the given external address.
func (k Keeper) setPairIDFromExternal(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
	id []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20PairByExternalKeyPrefix)
	store.Set(externalAddress.Bytes(), id)
}

// -----------------------------------------------------------------------------
// Internal Address -> ERC20BridgePair ID

// GetPairIDFromInternal gets the erc20 bridge pair id from the given internal address.
func (k Keeper) GetPairIDFromInternal(
	ctx sdk.Context,
	internalAddress types.InternalEVMAddress,
) ([]byte, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20PairByInternalKeyPrefix)
	id := store.Get(internalAddress.Bytes())
	return id, id != nil
}

// setPairIDFromInternal sets the erc20 bridge pair id for the given internal address.
func (k Keeper) setPairIDFromInternal(
	ctx sdk.Context,
	internalAddress types.InternalEVMAddress,
	id []byte,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20PairByInternalKeyPrefix)
	store.Set(internalAddress.Bytes(), id)
}
