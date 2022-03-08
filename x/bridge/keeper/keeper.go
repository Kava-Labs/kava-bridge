package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// GetOrDeployInternalERC20 returns the internal EVM address
// mapped to the provided ExternalEVMAddress. This will either return from the
// store if it is already deployed, or will first deploy the internal ERC20
// contract and return the new address.
func (k Keeper) GetOrDeployInternalERC20(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
) (types.InternalEVMAddress, error) {
	internalAddress, found := k.GetInternalERC20Address(ctx, externalAddress)
	if found {
		// If external ERC20 address is already mapped in store, there is
		// already a ERC20 deployed on Kava EVM.
		return internalAddress, nil
	}

	// The first time this external ERC20 is being bridged.
	// Check params for enabled ERC20. This both ensures the ERC20 is
	// whitelisted and fetches required ERC20 metadata: name, symbol,
	// decimals.
	enabledToken, err := k.GetEnabledERC20Token(ctx, externalAddress.String())
	if err != nil {
		return types.InternalEVMAddress{}, err
	}

	// Deploy the ERC20 contract on the Kava EVM
	internalAddress, err = k.DeployMintableERC20Contract(ctx, enabledToken)
	if err != nil {
		return types.InternalEVMAddress{}, err
	}

	addrPair := types.NewERC20BridgePair(externalAddress, internalAddress)
	if err := addrPair.Validate(); err != nil {
		return types.InternalEVMAddress{}, err
	}

	// Save the internal ERC20 address to state so that it is mapped to the
	// external ERC20 address.
	k.SetERC20AddressPair(ctx, addrPair)

	return internalAddress, nil
}

// IsSignerAuthorized returns an error if the provided message signer does not
// match the relayer set in params.
func (k Keeper) IsSignerAuthorized(ctx sdk.Context, msgSigners []sdk.AccAddress) error {
	relayer := k.GetRelayer(ctx)

	if len(msgSigners) != 1 {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid number of signers, expected 1 but got %d",
			len(msgSigners),
		)
	}

	if !relayer.Equals(msgSigners[0]) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"signer not authorized for bridge message",
		)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Store methods

// BridgedInternalEVMAddress puts the bridged address pair into the store.
func (k Keeper) SetERC20AddressPair(
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

	return types.InternalEVMAddress{
		Address: common.BytesToAddress(pair.InternalERC20Address),
	}, true
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
