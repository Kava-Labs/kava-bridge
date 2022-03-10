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

// Key prefixes
var (
	BridgedERC20KeyPrefix = []byte{0x01} // prefix for keys that store local ERC 20 address of bridged assets
)

// GetBridgedERC20Key returns the bytes of an BridgedERC20 key. This only
// accepts an ExternalEVMAddress and should not be used with
// InternalEVMAddresses.
func GetBridgedERC20Key(address ExternalEVMAddress) []byte {
	return append(BridgedERC20KeyPrefix, address.Bytes()...)
}
