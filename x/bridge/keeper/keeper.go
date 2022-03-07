package keeper

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	// Save the internal ERC20 address to state so that it is mapped to the
	// external ERC20 address.
	k.SetERC20AddressPair(ctx, externalAddress, internalAddress)

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
	externalAddress types.ExternalEVMAddress,
	internalAddress types.InternalEVMAddress,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)
	store.Set(types.GetBridgedERC20Key(externalAddress), internalAddress.Bytes())
}

// GetInternalERC20Address gets the internal EVM address mapped to the
// provided ExternalEVMAddress from the store.
func (k Keeper) GetInternalERC20Address(
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

// DeleteInternalERC20Address removes an bridged address pair from the store.
func (k Keeper) DeleteInternalERC20Address(ctx sdk.Context, externalAddress types.ExternalEVMAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20KeyPrefix)
	store.Delete(types.GetBridgedERC20Key(externalAddress))
}
