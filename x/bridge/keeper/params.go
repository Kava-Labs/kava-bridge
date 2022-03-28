package keeper

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// -----------------------------------------------------------------------------
// Relayer
// -----------------------------------------------------------------------------

// GetRelayer returns the relayer in params
func (k Keeper) GetRelayer(ctx sdk.Context) sdk.AccAddress {
	params := k.GetParams(ctx)
	return params.Relayer
}

// SetRelayer sets the relayer in params
func (k Keeper) SetRelayer(ctx sdk.Context, relayer sdk.AccAddress) {
	params := k.GetParams(ctx)
	params.Relayer = relayer
	k.SetParams(ctx, params)
}

// -----------------------------------------------------------------------------
// EnabledERC20Tokens
// -----------------------------------------------------------------------------

// GetEnabledERC20Token returns an EnabledERC20Token from the contract address
func (k Keeper) GetEnabledERC20Token(ctx sdk.Context, address types.ExternalEVMAddress) (types.EnabledERC20Token, error) {
	params := k.GetParams(ctx)
	for _, token := range params.EnabledERC20Tokens {
		if bytes.Equal(token.Address, address.Bytes()) {
			return token, nil
		}
	}

	return types.EnabledERC20Token{}, sdkerrors.Wrap(types.ErrERC20NotEnabled, address.String())
}

// GetEnabledERC20Tokens returns the all EnabledERC20Tokens
func (k Keeper) GetEnabledERC20Tokens(ctx sdk.Context) types.EnabledERC20Tokens {
	params := k.GetParams(ctx)
	return params.EnabledERC20Tokens
}

// -----------------------------------------------------------------------------
// EnabledConversionPairs
// -----------------------------------------------------------------------------

// GetEnabledConversionPair returns an ConversionPair from the internal contract address.
func (k Keeper) GetEnabledConversionPair(
	ctx sdk.Context,
	address types.InternalEVMAddress,
) (types.ConversionPair, error) {
	params := k.GetParams(ctx)
	for _, pair := range params.EnabledConversionPairs {
		if bytes.Equal(pair.KavaERC20Address, address.Bytes()) {
			return pair, nil
		}
	}

	return types.ConversionPair{}, sdkerrors.Wrap(types.ErrConversionNotEnabled, address.String())
}

// GetEnabledConversionPairs returns the all ConversionPairs
func (k Keeper) GetEnabledConversionPairs(ctx sdk.Context) types.ConversionPairs {
	params := k.GetParams(ctx)
	return params.EnabledConversionPairs
}
