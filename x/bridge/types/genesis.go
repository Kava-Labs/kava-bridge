package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// DefaultNextWithdrawSequence is the starting point for withdraw sequence.
	DefaultNextWithdrawSequence sdk.Int = sdk.OneInt()

	// 1 << 256 == 257 bit int
	i = new(big.Int).Lsh(big.NewInt(1), 256)
	// 1 << 256 - 1 == max 256 bit int
	MaxWithdrawSequenceBigInt = new(big.Int).Sub(i, big.NewInt(1))

	// MaxWithdrawSequence is the maximum sdk.Int value a withdraw sequence can
	// be before it wraps.
	MaxWithdrawSequence = sdk.NewIntFromBigInt(MaxWithdrawSequenceBigInt)
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	erc20BridgePairs ERC20BridgePairs,
	nextSequence sdk.Int,
) GenesisState {
	return GenesisState{
		Params:               params,
		ERC20BridgePairs:     erc20BridgePairs,
		NextWithdrawSequence: nextSequence,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		ERC20BridgePairs:     ERC20BridgePairs{},
		NextWithdrawSequence: DefaultNextWithdrawSequence,
	}
}

// Validate validates genesis inputs. It returns error if validation of any input fails.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if err := gs.ERC20BridgePairs.Validate(); err != nil {
		return err
	}

	return nil
}
