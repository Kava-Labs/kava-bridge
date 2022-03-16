package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// GetOrDeployInternalERC20 returns the internal EVM address
// mapped to the provided ExternalEVMAddress. This will either return from the
// store if it is already deployed, or will first deploy the internal ERC20
// contract and return the new address.
func (k Keeper) GetOrDeployInternalERC20(
	ctx sdk.Context,
	externalAddress types.ExternalEVMAddress,
) (types.InternalEVMAddress, error) {
	pair, found := k.GetBridgePairFromExternal(ctx, externalAddress)
	if found {
		// If external ERC20 address is already mapped in store, there is
		// already a ERC20 deployed on Kava EVM.
		return pair.GetInternalAddress(), nil
	}

	// The first time this external ERC20 is being bridged.
	// Check params for enabled ERC20. This both ensures the ERC20 is
	// whitelisted and fetches required ERC20 metadata: name, symbol,
	// decimals.
	enabledToken, err := k.GetEnabledERC20Token(ctx, externalAddress)
	if err != nil {
		return types.InternalEVMAddress{}, err
	}

	// Deploy the ERC20 contract on the Kava EVM
	internalAddress, err := k.DeployMintableERC20Contract(ctx, enabledToken)
	if err != nil {
		return types.InternalEVMAddress{}, err
	}

	addrPair := types.NewERC20BridgePair(externalAddress, internalAddress)
	if err := addrPair.Validate(); err != nil {
		return types.InternalEVMAddress{}, err
	}

	// Save the internal ERC20 address to state in all indices.
	k.RegisterBridgePair(ctx, addrPair)

	return internalAddress, nil
}

// DeployMintableERC20Contract deploys an ERC20 contract on the EVM as the
// module account and returns the address of the contract. This contract has
// minting permissions for the module account.
// Derived from tharsis/evmos
// https://github.com/tharsis/evmos/blob/ee54f496551df937915ff6f74a94732a35abc505/x/erc20/keeper/evm.go
func (k Keeper) DeployMintableERC20Contract(
	ctx sdk.Context,
	token types.EnabledERC20Token,
) (types.InternalEVMAddress, error) {
	ctorArgs, err := contract.ERC20MintableBurnableContract.ABI.Pack(
		"", // Empty string for contract constructor
		token.Name,
		token.Symbol,
		uint8(token.Decimals),
	)
	if err != nil {
		return types.InternalEVMAddress{}, sdkerrors.Wrapf(err, "token %v is invalid", token.Name)
	}

	data := make([]byte, len(contract.ERC20MintableBurnableContract.Bin)+len(ctorArgs))
	copy(
		data[:len(contract.ERC20MintableBurnableContract.Bin)],
		contract.ERC20MintableBurnableContract.Bin,
	)
	copy(
		data[len(contract.ERC20MintableBurnableContract.Bin):],
		ctorArgs,
	)

	nonce, err := k.accountKeeper.GetSequence(ctx, types.ModuleEVMAddress.Bytes())
	if err != nil {
		return types.InternalEVMAddress{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleEVMAddress, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleEVMAddress, nil, data)
	if err != nil {
		return types.InternalEVMAddress{}, fmt.Errorf("failed to deploy ERC20 for %s: %w", token.Name, err)
	}

	return types.NewInternalEVMAddress(contractAddr), nil
}

// MintERC20 mints the given amount of an ERC20 token to an address. This is
// unchecked and should only be called after permission and enabled ERC20 checks.
func (k Keeper) MintERC20(
	ctx sdk.Context,
	contractAddr types.InternalEVMAddress,
	receiver common.Address,
	amount *big.Int,
) error {
	_, err := k.CallEVM(
		ctx,
		contract.ERC20MintableBurnableContract.ABI,
		types.ModuleEVMAddress,
		contractAddr,
		"mint",
		// Mint ERC20 args
		receiver,
		amount,
	)

	return err
}
