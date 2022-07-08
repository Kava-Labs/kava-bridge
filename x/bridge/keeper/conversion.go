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
	amount *big.Int,
	recipient sdk.AccAddress,
) (sdk.Coin, error) {
	coin := sdk.NewCoin(pair.Denom, sdk.NewIntFromBigInt(amount))
	coins := sdk.NewCoins(coin)

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return sdk.Coin{}, err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins); err != nil {
		return sdk.Coin{}, err
	}

	return coin, nil
}

// ConvertCoinToERC20 converts an sdk.Coin from the originating account to an
// ERC20 to the receiver account.
func (k Keeper) ConvertCoinToERC20(
	ctx sdk.Context,
	initiatorAccount sdk.AccAddress,
	receiverAccount types.InternalEVMAddress,
	coin sdk.Coin,
) error {
	params := k.GetParams(ctx)
	if !params.BridgeEnabled {
		return types.ErrBridgeDisabled
	}

	pair, err := k.GetEnabledConversionPairFromDenom(ctx, coin.Denom)
	if err != nil {
		// Coin not in enabled conversion pair list
		return err
	}

	if err := k.BurnConversionPairCoin(ctx, pair, coin, initiatorAccount); err != nil {
		return err
	}

	if err := k.UnlockERC20Tokens(ctx, pair, coin.Amount.BigInt(), receiverAccount); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeConvertCoinToERC20,
		sdk.NewAttribute(types.AttributeKeyInitiator, initiatorAccount.String()),
		sdk.NewAttribute(types.AttributeKeyReceiver, receiverAccount.String()),
		sdk.NewAttribute(types.AttributeKeyERC20Address, pair.GetAddress().String()),
		sdk.NewAttribute(types.AttributeKeyAmount, coin.String()),
	))

	return nil
}

// ConvertERC20ToCoin converts an ERC20 coin from the originating account to an
// sdk.Coin to the receiver account.
func (k Keeper) ConvertERC20ToCoin(
	ctx sdk.Context,
	initiator types.InternalEVMAddress,
	receiver sdk.AccAddress,
	contractAddr types.InternalEVMAddress,
	amount sdk.Int,
) error {
	params := k.GetParams(ctx)
	if !params.BridgeEnabled {
		return types.ErrBridgeDisabled
	}

	// Check that the contract is enabled to convert to coin
	pair, err := k.GetEnabledConversionPairFromERC20Address(ctx, contractAddr)
	if err != nil {
		// contract not in enabled conversion pair list
		return err
	}

	// lock erc20 tokens
	if err := k.LockERC20Tokens(ctx, pair, amount.BigInt(), initiator); err != nil {
		return err
	}

	// mint conversion pair coin
	coin, err := k.MintConversionPairCoin(ctx, pair, amount.BigInt(), receiver)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeConvertERC20ToCoin,
		sdk.NewAttribute(types.AttributeKeyERC20Address, contractAddr.String()),
		sdk.NewAttribute(types.AttributeKeyInitiator, initiator.String()),
		sdk.NewAttribute(types.AttributeKeyReceiver, receiver.String()),
		sdk.NewAttribute(types.AttributeKeyAmount, coin.String()),
	))

	return nil
}

// BurnConversionPairCoin transfers the provided amount to the module account
// then burns it.
func (k Keeper) BurnConversionPairCoin(
	ctx sdk.Context,
	pair types.ConversionPair,
	coin sdk.Coin,
	account sdk.AccAddress,
) error {
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

// LockERC20Tokens transfers the given amount of a conversion pair ERC20 token
// from the initiator account to the module account.
func (k Keeper) LockERC20Tokens(
	ctx sdk.Context,
	pair types.ConversionPair,
	amount *big.Int,
	initiator types.InternalEVMAddress,
) error {
	_, err := k.CallEVM(
		ctx,
		contract.ERC20MintableBurnableContract.ABI, // abi
		initiator.Address,                          // from addr
		pair.GetAddress(),                          // contract addr
		"transfer",                                 // method
		// Transfer ERC20 args
		types.ModuleEVMAddress,
		amount,
	)

	return err
}
