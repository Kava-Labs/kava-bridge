package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/kava-labs/kava-bridge/x/bridge/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants registers the bridge module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "bridge-pairs", BridgePairsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "bridge-pairs-index", BridgePairsIndexInvariant(k))
	ir.RegisterRoute(types.ModuleName, "backed-coins", BackedCoinsInvariant(k))
}

// AllInvariants runs all invariants of the bridge module
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		if res, stop := BridgePairsInvariant(k)(ctx); stop {
			return res, stop
		}

		if res, stop := BridgePairsIndexInvariant(k)(ctx); stop {
			return res, stop
		}

		res, stop := BackedCoinsInvariant(k)(ctx)
		return res, stop
	}
}

// BridgePairsInvariant iterates all bridge pairs and asserts that they are valid
func BridgePairsInvariant(k Keeper) sdk.Invariant {
	broken := false
	message := sdk.FormatInvariant(types.ModuleName, "validate bridge pairs broken", "bridge pair invalid")

	return func(ctx sdk.Context) (string, bool) {
		k.IterateBridgePairs(ctx, func(pair types.ERC20BridgePair) bool {
			if err := pair.Validate(); err != nil {
				broken = true
				return true
			}
			return false
		})

		return message, broken
	}
}

// BridgePairsIndexInvariant iterates all bridge pairs and asserts the index for
// querying are all valid.
func BridgePairsIndexInvariant(k Keeper) sdk.Invariant {
	broken := false
	message := sdk.FormatInvariant(types.ModuleName, "validate bridge pairs broken", "bridge pair invalid")

	return func(ctx sdk.Context) (string, bool) {
		store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BridgedERC20PairKeyPrefix)

		byExternalIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BridgedERC20PairByExternalKeyPrefix)
		defer byExternalIterator.Close()

		var byExternalIndexLength int
		for ; byExternalIterator.Valid(); byExternalIterator.Next() {
			byExternalIndexLength++

			idBytes := byExternalIterator.Value()
			pairBytes := store.Get(idBytes)
			if pairBytes == nil {
				invariantMessage := sdk.FormatInvariant(
					types.ModuleName,
					"valid index",
					fmt.Sprintf(
						"\tbridge pair with ID '%s' found in external index but not in store",
						hex.EncodeToString(idBytes),
					),
				)
				return invariantMessage, true
			}
		}

		byInternalIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BridgedERC20PairByInternalKeyPrefix)
		defer byInternalIterator.Close()

		var byInternalIndexLength int
		for ; byInternalIterator.Valid(); byInternalIterator.Next() {
			byInternalIndexLength++

			idBytes := byInternalIterator.Value()
			pairBytes := store.Get(idBytes)
			if pairBytes == nil {
				invariantMessage := sdk.FormatInvariant(
					types.ModuleName,
					"valid index",
					fmt.Sprintf(
						"\tbridge pair with ID '%s' found in internal index but not in store",
						hex.EncodeToString(idBytes),
					),
				)
				return invariantMessage, true
			}
		}

		// Check length of bridge pairs store matches the length of both indices
		storeIterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BridgedERC20PairKeyPrefix)
		defer storeIterator.Close()
		var storeLength int
		for ; storeIterator.Valid(); storeIterator.Next() {
			storeLength++
		}

		if storeLength != byExternalIndexLength || storeLength != byInternalIndexLength {
			invariantMessage := sdk.FormatInvariant(
				types.ModuleName,
				"valid index",
				fmt.Sprintf(
					"\tmismatched number of items in bridge pair store (%d), internal index (%d), and external index (%d)",
					storeLength, byInternalIndexLength, byExternalIndexLength,
				),
			)
			return invariantMessage, true
		}

		return message, broken
	}
}

// BackedCoinsInvariant iterates all conversion pairs and asserts that the
// sdk.Coin balances are less than the module ERC20 balance.
// **Note:** This compares <= and not == as anyone can send tokens to the
// ERC20 contract address and break the invariant if a strict equal check.
func BackedCoinsInvariant(k Keeper) sdk.Invariant {
	broken := false
	message := sdk.FormatInvariant(
		types.ModuleName,
		"backed coins broken",
		"coin supply is greater than module account ERC20 tokens",
	)

	return func(ctx sdk.Context) (string, bool) {
		params := k.GetParams(ctx)
		for _, pair := range params.EnabledConversionPairs {
			erc20Balance, err := k.QueryERC20BalanceOf(
				ctx,
				pair.GetAddress(),
				types.NewInternalEVMAddress(types.ModuleEVMAddress),
			)
			// TODO: Panic or set broken here?
			if err != nil {
				panic(err)
			}

			supply := k.bankKeeper.GetSupply(ctx, pair.Denom)

			// Must be true: sdk.Coin supply < ERC20 balanceOf(module account)
			if supply.Amount.BigInt().Cmp(erc20Balance) > 0 {
				broken = true
				break
			}
		}

		return message, broken
	}
}
