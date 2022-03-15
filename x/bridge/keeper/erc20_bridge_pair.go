package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// SetERC20BridgePair puts the bridged address pair into the store.
func (k Keeper) SetERC20BridgePair(
	ctx sdk.Context,
	pair types.ERC20BridgePair,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)
	bz := k.cdc.MustMarshal(&pair)
	store.Set(types.GetBridgedERC20Key(pair.ExternalERC20Address), bz)
}

// GetInternalERC20Address gets the internal EVM address mapped to the
// provided ExternalEVMAddress from the store.
func (k Keeper) GetInternalERC20Address(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
) (types.InternalEVMAddress, bool) {
	pair, found := k.GetERC20BridgePair(ctx, externalAddress)
	if !found {
		return types.InternalEVMAddress{}, false
	}

	return types.NewInternalEVMAddress(
		common.BytesToAddress(pair.InternalERC20Address),
	), true
}

// GetERC20BridgePair returns the ERC20 bridge pair with the provided
// ExternalEVMAddress from the store.
func (k Keeper) GetERC20BridgePair(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
) (types.ERC20BridgePair, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)
	bz := store.Get(types.GetBridgedERC20Key(externalAddress.Bytes()))
	if bz == nil {
		return types.ERC20BridgePair{}, false
	}

	var pair types.ERC20BridgePair
	k.cdc.MustUnmarshal(bz, &pair)

	return pair, true
}

// IterateERC20BridgePairs provides an iterator over all stored ERC20 bridge
// pairs. For each pair, cb will be called. If cb returns true, the iterator
// will close and stop.
func (k Keeper) IterateERC20BridgePairs(
	ctx sdk.Context,
	cb func(pair types.ERC20BridgePair) (stop bool),
) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var pair types.ERC20BridgePair
		k.cdc.MustUnmarshal(iterator.Value(), &pair)

		if cb(pair) {
			break
		}
	}
}
