package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	erc20BalanceOfMethod = "balanceOf"
)

func (k Keeper) QueryERC20BalanceOf(
	ctx sdk.Context,
	contractAddr types.InternalEVMAddress,
	account types.InternalEVMAddress,
) (*big.Int, error) {
	res, err := k.CallEVM(
		ctx,
		contract.ERC20MintableBurnableContract.ABI,
		types.ModuleEVMAddress,
		contractAddr,
		erc20BalanceOfMethod,
		// balanceOf ERC20 args
		account.Address,
	)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		if res.VmError == vm.ErrExecutionReverted.Error() {
			// Unpacks revert
			return nil, evmtypes.NewExecErrorWithReason(res.Ret)
		}

		return nil, status.Error(codes.Internal, res.VmError)
	}

	anyOutput, err := contract.ERC20MintableBurnableContract.ABI.Unpack(erc20BalanceOfMethod, res.Ret)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unpack method %v response: %w",
			erc20BalanceOfMethod,
			err,
		)
	}

	if len(anyOutput) != 1 {
		return nil, fmt.Errorf(
			"invalid ERC20 %v call return outputs %v, expected %v",
			erc20BalanceOfMethod,
			len(anyOutput),
			1,
		)
	}

	bal, ok := anyOutput[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf(
			"invalid ERC20 return type %T, expected %T",
			anyOutput[0],
			&big.Int{},
		)
	}

	return bal, nil
}
