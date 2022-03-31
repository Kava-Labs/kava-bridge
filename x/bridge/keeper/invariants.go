package keeper

import (
	"github.com/kava-labs/kava-bridge/x/bridge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants registers the bridge module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "bridge-pairs", BridgePairsInvariant(k))
}

// AllInvariants runs all invariants of the bridge module
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := BridgePairsInvariant(k)(ctx)
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
