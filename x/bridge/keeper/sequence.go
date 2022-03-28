package keeper

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// GetNextWithdrawSequence reads the next available global withdraw sequence
// from store.
func (k Keeper) GetNextWithdrawSequence(ctx sdk.Context) (sdk.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NextWithdrawSequenceKeyPrefix)
	bz := store.Get(types.NextWithdrawSequenceKeyPrefix)
	if bz == nil {
		return sdk.Int{}, types.ErrInvalidInitialWithdrawSequence
	}

	var seq sdk.Int
	if err := seq.Unmarshal(bz); err != nil {
		panic(err)
	}

	return seq, nil
}

// IncrementNextWithdrawSequence increments the next withdraw sequence in the
// store by 1.
func (k Keeper) IncrementNextWithdrawSequence(ctx sdk.Context) error {
	seq, err := k.GetNextWithdrawSequence(ctx)
	if err != nil {
		return err
	}

	k.SetNextWithdrawSequence(ctx, WrappingAddInt(seq, sdk.OneInt()))
	return nil
}

// SetNextWithdrawSequence stores a sequence to be used for the next withdraw.
func (k Keeper) SetNextWithdrawSequence(ctx sdk.Context, sequence sdk.Int) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NextWithdrawSequenceKeyPrefix)

	bz, err := sequence.Marshal()
	if err != nil {
		panic(err)
	}

	store.Set(types.NextWithdrawSequenceKeyPrefix, bz)
}

// WrappingAddInt adds two sdk.Int values with intentional wrapping.
func WrappingAddInt(i1 sdk.Int, i2 sdk.Int) sdk.Int {
	sum := new(big.Int).Add(i1.BigInt(), i2.BigInt())

	// Total is less than max, safe to convert to sdk.Int
	if sum.Cmp(types.MaxWithdrawSequenceBigInt) <= 0 {
		return sdk.NewIntFromBigInt(sum)
	}

	// Only keep the lower 256 bits to wrap
	wrappedSum := new(big.Int).And(sum, types.MaxWithdrawSequenceBigInt)
	return sdk.NewIntFromBigInt(wrappedSum)
}
