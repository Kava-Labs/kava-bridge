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
	BridgedERC20PairKeyPrefix           = []byte{0x01} // prefix for keys that store a ERC20 bridge pair
	BridgedERC20PairByExternalKeyPrefix = []byte{0x02} // prefix for keys that store the ID of a ERC20 bridge pair by external address
	BridgedERC20PairByInternalKeyPrefix = []byte{0x03} // prefix for keys that store the ID of a ERC20 bridge pair by internal address
	NextWithdrawSequenceKeyPrefix       = []byte{0x04} // prefix for key of next withdraw sequence
)
