package keeper

import (
	"encoding/json"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/kava-labs/kava-bridge/x/bridge/types"

	"github.com/tharsis/ethermint/server/config"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// CallEVM performs a smart contract method call using given args
func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from common.Address,
	contract types.InternalEVMAddress,
	method string,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrABIPack,
			sdkerrors.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	resp, err := k.CallEVMWithData(ctx, from, &contract, data)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}

// CallEVMWithData performs a smart contract method call using contract data
// Derived from tharsis/evmos
// https://github.com/tharsis/evmos/blob/ee54f496551df937915ff6f74a94732a35abc505/x/erc20/keeper/evm.go
func (k Keeper) CallEVMWithData(
	ctx sdk.Context,
	from common.Address,
	contract *types.InternalEVMAddress,
	data []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	// To param needs to be nil to correctly apply txs to create contracts
	// Default common.Address value is 0x0000000000000000000000000000000000000000, not nil
	// which Ethermint handles differently -- erc20_test will fail
	// https://github.com/tharsis/ethermint/blob/caa1c5a6c6b7ed8ba4aaf6e0b0848f6be5ba6668/x/evm/keeper/state_transition.go#L357
	var to *common.Address
	if contract != nil {
		to = &contract.Address
	}

	transactionArgs := evmtypes.TransactionArgs{
		From: &from,
		To:   to,
		Data: (*hexutil.Bytes)(&data),
	}

	args, err := json.Marshal(transactionArgs)
	if err != nil {
		return nil, err
	}

	gasRes, err := k.evmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evmtypes.EthCallRequest{
		Args:   args,
		GasCap: config.DefaultGasCap,
	})
	if err != nil {
		return nil, err
	}

	msg := ethtypes.NewMessage(
		from,
		to,
		nonce,
		big.NewInt(0), // amount
		gasRes.Gas,    // gasLimit
		big.NewInt(0), // gasFeeCap
		big.NewInt(0), // gasTipCap
		big.NewInt(0), // gasPrice
		data,
		ethtypes.AccessList{}, // AccessList
		true,                  // checkNonce
	)

	res, err := k.evmKeeper.ApplyMessage(ctx, msg, evmtypes.NewNoOpTracer(), true)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return res, nil
}
