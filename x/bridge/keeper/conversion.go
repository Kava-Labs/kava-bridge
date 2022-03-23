package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// MintConversionPairCoin mints the given amount of a ConversionPair denom and
// sends it to the provided address.
func (k Keeper) MintConversionPairCoin(
	ctx sdk.Context,
	pair types.ConversionPair,
	amount sdk.Int,
	recipient sdk.AccAddress,
) error {
	coin := sdk.NewCoin(pair.Denom, amount)
	coins := sdk.NewCoins(coin)

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins); err != nil {
		return err
	}

	return nil
}

// ConvertCoinToERC20 converts an sdk.Coin from the originating account to an
// ERC20 to the receiver account.
func (k Keeper) ConvertCoinToERC20(
	ctx sdk.Context,
	pair types.ConversionPair,
	amount sdk.Int,
	originAccount sdk.AccAddress,
	receiverAccount types.InternalEVMAddress,
) error {
	if err := k.BurnConversionPairCoin(ctx, pair, amount, originAccount); err != nil {
		return err
	}

	if err := k.UnlockERC20Tokens(ctx, pair, amount.BigInt(), receiverAccount); err != nil {
		return err
	}

	return nil
}

// BurnConversionPairCoin transfers the provided amount to the module account
// then burns it.
func (k Keeper) BurnConversionPairCoin(
	ctx sdk.Context,
	pair types.ConversionPair,
	amount sdk.Int,
	account sdk.AccAddress,
) error {
	coin := sdk.NewCoin(pair.Denom, amount)
	coins := sdk.NewCoins(coin)

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, account, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	return nil
}

// UnlockERC20Tokens transfers the given amount of a conversion pair ERC20 token
// to the provided account.
func (k Keeper) UnlockERC20Tokens(
	ctx sdk.Context,
	pair types.ConversionPair,
	amount *big.Int,
	receiver types.InternalEVMAddress,
) error {
	_, err := k.CallEVM(
		ctx,
		contract.ERC20MintableBurnableContract.ABI, // abi
		types.ModuleEVMAddress,                     // from addr
		pair.GetAddress(),                          // contract addr
		"transfer",                                 // method
		// Transfer ERC20 args
		receiver.Address,
		amount,
	)

	return err
}
