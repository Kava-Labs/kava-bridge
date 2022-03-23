package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava-bridge/contract"
)

func GetERC20Amount(
	ctx sdk.Context,
	k Keeper,
	eventName string,
	data []byte,
) (*big.Int, error) {
	ConvertToCoinEvent, err := contract.ERC20MintableBurnableContract.ABI.Unpack(eventName, data)
	if err != nil {
		k.Logger(ctx).Error("failed to unpack ConvertToCoin event", "error", err.Error())
		return nil, err
	}

	if len(ConvertToCoinEvent) == 0 {
		k.Logger(ctx).Error("ConvertToCoin event data is empty", "error", err.Error())
		return nil, err
	}

	// Data only contains non-indexed parameters, which is only the amount
	amount, ok := ConvertToCoinEvent[0].(*big.Int)
	// safety check and ignore if amount not positive
	if !ok || amount == nil || amount.Sign() != 1 {
		return nil, fmt.Errorf("invalid event amount: %v", amount)
	}

	return amount, nil
}
