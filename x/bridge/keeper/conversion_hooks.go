package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// Hooks wrapper struct for bridge keeper
type ConversionHooks struct {
	k Keeper
}

var _ evmtypes.EvmHooks = ConversionHooks{}

// Return the wrapper struct
func (k Keeper) ConversionHooks() ConversionHooks {
	return ConversionHooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing. This handles minting
// sdk.Coin when ConvertToCoin() is called on an eligible Kava ERC20 contract.
func (h ConversionHooks) PostTxProcessing(
	ctx sdk.Context,
	from common.Address,
	to *common.Address,
	receipt *ethtypes.Receipt,
) error {
	erc20Abi := contract.ERC20MintableBurnableContract.ABI

	for _, log := range receipt.Logs {
		// ERC20MintableBurnableContract should contain 3 topics:
		// 0: Keccak-256 hash of ConvertToCoin(address,bytes32,uint256)
		// 1: address indexed sender
		// 2: address indexed toAddr
		if len(log.Topics) != 3 {
			continue
		}

		// event ID, e.g. Keccak-256 hash of ConvertToCoin(address,bytes32,uint256)
		eventID := log.Topics[0]

		event, err := erc20Abi.EventByID(eventID)
		if err != nil {
			// invalid event for ERC20
			continue
		}

		if event.Name != types.ContractEventTypeConvertToCoin {
			continue
		}

		ConvertToCoinEvent, err := erc20Abi.Unpack(event.Name, log.Data)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack ConvertToCoin event", "error", err.Error())
			continue
		}

		if len(ConvertToCoinEvent) == 0 {
			h.k.Logger(ctx).Error("ConvertToCoin event data is empty", "error", err.Error())
			continue
		}

		// Data only contains non-indexed parameters, which is only the amount
		amount, ok := ConvertToCoinEvent[0].(*big.Int)
		// safety check and ignore if amount not positive
		if !ok || amount == nil || amount.Sign() != 1 {
			continue
		}

		// Check that the contract is enabled to convert to coin
		contractAddr := types.NewInternalEVMAddress(log.Address)
		conversionPair, err := h.k.GetEnabledConversionPair(ctx, contractAddr)
		if err != nil {
			// Contract not a conversion pair in state
			continue
		}

		initiator := common.BytesToAddress(log.Topics[1].Bytes())

		// TODO: Must handle padded value
		receiver := sdk.AccAddress(log.Topics[2].Bytes())

		// Initiator is a **different** address from receiver
		coin, err := h.k.MintConversionPairCoin(ctx, conversionPair, amount, receiver)
		if err != nil {
			// Revert tx if conversion fails
			panic(err)
		}

		if err := ctx.EventManager().EmitTypedEvent(&types.EventConvertERC20ToCoin{
			ERC20Address: contractAddr.String(),
			Initiator:    initiator.String(),
			Receiver:     receiver.String(),
			Amount:       &coin,
		}); err != nil {
			panic(err)
		}
	}

	return nil
}
