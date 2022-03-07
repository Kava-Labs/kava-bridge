package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName name used throughout module
	ModuleName = "bridge"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey Top level router key
	RouterKey = ModuleName

	// QuerierRoute should be set to module name
	QuerierRoute = ModuleName
)

// ModuleAddress is the native module address for EVM
var ModuleEVMAddress common.Address

func init() {
	ModuleEVMAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}
