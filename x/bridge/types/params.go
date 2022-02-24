package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ParamKeyTable for bridge module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable() // TODO: .RegisterParamSet(&Params{})
}
