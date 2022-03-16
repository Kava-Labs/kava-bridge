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
type Hooks struct {
	k Keeper
}

var _ evmtypes.EvmHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (h Hooks) PostTxProcessing(
	ctx sdk.Context,
	from common.Address,
	to *common.Address,
	receipt *ethtypes.Receipt,
) error {
	erc20Abi := contract.ERC20MintableBurnableContract.ABI

	for _, log := range receipt.Logs {
		if len(log.Topics) < 3 {
			continue
		}

		// event ID, e.g. Keccak-256 hash of Withdraw(address,address,uint256)
		eventID := log.Topics[0]

		event, err := erc20Abi.EventByID(eventID)
		if err != nil {
			// invalid event for ERC20
			continue
		}

		if event.Name != types.ContractEventTypeWithdraw {
			continue
		}

		withdrawEvent, err := erc20Abi.Unpack(event.Name, log.Data)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack withdraw event", "error", err.Error())
			continue
		}

		if len(withdrawEvent) == 0 {
			h.k.Logger(ctx).Error("withdraw event data is empty", "error", err.Error())
			continue
		}

		// Data only contains non-indexed parameters, which is only the amount
		amount, ok := withdrawEvent[0].(*big.Int)
		// safety check and ignore if amount not positive
		if !ok || amount == nil || amount.Sign() != 1 {
			continue
		}

		// Check that the contract is an enabled token pair
		contractAddr := types.NewInternalEVMAddress(log.Address)
		pair, found := h.k.GetBridgePairFromInternal(ctx, contractAddr)
		if !found {
			// Contract not a bridge pair in state
			continue
		}

		externalERC20Addr := pair.GetExternalAddress()
		toAddr := common.BytesToAddress(log.Topics[2].Bytes())

		sequence, err := h.k.GetNextWithdrawSequence(ctx)
		if err != nil {
			// Panic since we actually want to revert the entire TX if this
			// fails otherwise funds would be burned without event emitted for
			// relayer to unlock.
			panic(err)
		}

		if err := h.k.IncrementNextWithdrawSequence(ctx); err != nil {
			panic(err)
		}

		// TODO: Replace with EmitTypedEvent + EventWithdraw proto type
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeySequence, sequence.String()),
			sdk.NewAttribute(types.AttributeKeyEthereumERC20Address, externalERC20Addr.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, toAddr.String()),
		))
	}

	return nil
}
