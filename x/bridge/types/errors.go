package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrERC20NotEnabled = sdkerrors.Register(ModuleName, 2, "erc20 token not enabled")
	ErrABIPack         = sdkerrors.Register(ModuleName, 3, "contract ABI pack failed")
	ErrABIUnpack       = sdkerrors.Register(ModuleName, 4, "contract ABI unpack failed")
)
