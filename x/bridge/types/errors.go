package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrERC20NotEnabled                = sdkerrors.Register(ModuleName, 2, "ERC20 token not enabled")
	ErrABIPack                        = sdkerrors.Register(ModuleName, 3, "contract ABI pack failed")
	ErrABIUnpack                      = sdkerrors.Register(ModuleName, 4, "contract ABI unpack failed")
	ErrInvalidInitialWithdrawSequence = sdkerrors.Register(ModuleName, 5, "initial withdraw sequence hasn't been set")
	ErrConversionNotEnabled           = sdkerrors.Register(ModuleName, 6, "ERC20 token not enabled to convert to sdk.Coin")
	ErrBridgeDisabled                 = sdkerrors.Register(ModuleName, 7, "Bridge transactions and conversions are disabled")
	ErrNoRelayer                      = sdkerrors.Register(ModuleName, 8, "There is no relayer set")
)
