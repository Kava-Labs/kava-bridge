// Derived from https://github.com/tharsis/evmos/blob/0bfaf0db7be47153bc651e663176ba8deca960b5/x/erc20/keeper/evm_hooks.go
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// Hooks wrapper struct for bridge keeper
type WithdrawHook struct {
	k Keeper
}

var _ evmtypes.EvmHooks = WithdrawHook{}

// Return the wrapper struct
func (k Keeper) WithdrawHooks() WithdrawHook {
	return WithdrawHook{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (h WithdrawHook) PostTxProcessing(
	ctx sdk.Context,
	from common.Address,
	to *common.Address,
	receipt *ethtypes.Receipt,
) error {
	erc20Abi := contract.ERC20MintableBurnableContract.ABI
	params := h.k.GetParams(ctx)

	for _, log := range receipt.Logs {
		// ERC20MintableBurnableContract should contain 3 topics:
		// 0: Keccak-256 hash of Withdraw(address,address,uint256)
		// 1: address indexed sender
		// 2: address indexed toAddr
		if len(log.Topics) != 3 {
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

		// Only check if the bridge is enabled for contracts that are enabled.
		if !params.BridgeEnabled {
			return types.ErrBridgeDisabled
		}

		enabledERC20, err := h.k.GetEnabledERC20TokenFromExternal(ctx, pair.GetExternalAddress())
		if err != nil {
			// This will error only if an existing erc20 was enabled in params,
			// had the Kava erc20 contract deployed, then later removed. Doing
			// so does *not* remove the ERC20 contract from erc20 bridge pairs
			// state.

			// Error is not user facing, but is logged.
			return fmt.Errorf("failed to get enabled ERC20 token: %w", err)
		}

		if amount.Cmp(enabledERC20.MinimumWithdrawAmount.BigInt()) < 0 {
			// Return error to revert transaction and avoid loss of funds.
			// Only revert when this is an enabled ERC20 token.
			return fmt.Errorf(
				"withdraw amount is less than minimum withdraw amount: %v < %v",
				amount,
				enabledERC20.MinimumWithdrawAmount,
			)
		}

		externalERC20Addr := pair.GetExternalAddress()
		toAddr := common.BytesToAddress(log.Topics[2].Bytes())

		// Panics since we actually want to revert the entire TX if any of these
		// fail otherwise funds would be burned without event emitted for
		// relayer to unlock.

		sequence, err := h.k.GetNextWithdrawSequence(ctx)
		if err != nil {
			panic(err)
		}

		if err := h.k.IncrementNextWithdrawSequence(ctx); err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeBridgeKavaToEthereum,
			sdk.NewAttribute(types.AttributeKeyEthereumERC20Address, externalERC20Addr.String()),
			sdk.NewAttribute(types.AttributeKeyKavaERC20Address, contractAddr.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, toAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeySequence, sequence.String()),
		))
	}

	return nil
}
