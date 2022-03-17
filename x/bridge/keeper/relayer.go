package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IsSignerAuthorized returns an error if the provided message signer does not
// match the relayer set in params.
func (k Keeper) IsSignerAuthorized(ctx sdk.Context, signer sdk.AccAddress) error {
	relayer := k.GetRelayer(ctx)

	if !relayer.Equals(signer) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"signer not authorized for bridge message",
		)
	}

	return nil
}
