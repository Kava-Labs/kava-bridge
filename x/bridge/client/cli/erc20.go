package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/tharsis/ethermint/server/config"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// ERC20Query queries a contract with a specific method and input parameters
func ERC20Query(
	ctx client.Context,
	contractAddr common.Address,
	method string,
	args ...interface{},
) ([]interface{}, error) {
	data, err := PackContractCallData(
		contract.ERC20MintableBurnableContract.ABI,
		method,
		args...,
	)
	if err != nil {
		return nil, err
	}

	transactionArgs := evmtypes.TransactionArgs{
		To:   &contractAddr,
		Data: (*hexutil.Bytes)(&data),
	}

	ethCalArgs, err := json.Marshal(transactionArgs)
	if err != nil {
		return nil, err
	}

	evmQueryClient := evmtypes.NewQueryClient(ctx)
	res, err := evmQueryClient.EthCall(context.Background(), &evmtypes.EthCallRequest{
		Args:   ethCalArgs,
		GasCap: uint64(config.DefaultGasCap),
	})
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		if res.VmError == vm.ErrExecutionReverted.Error() {
			// Unpacks revert
			return nil, evmtypes.NewExecErrorWithReason(res.Ret)
		}

		return nil, fmt.Errorf(res.VmError)
	}

	anyOutput, err := contract.ERC20MintableBurnableContract.ABI.Unpack(method, res.Ret)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unpack method %v response: %w",
			method,
			err,
		)
	}

	return anyOutput, nil
}
