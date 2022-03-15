package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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
