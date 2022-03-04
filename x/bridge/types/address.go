package types

import "github.com/ethereum/go-ethereum/common"

// ExternalEVMAddress is a type alias of common.Address to represent an address
// on an external EVM, e.g. Ethereum. This is used to make external / internal
// addresses type safe and un-assignable to each other. This also makes it more
// clear which address belongs where.
type ExternalEVMAddress struct {
	common.Address
}

// InternalEVMAddress is a type alias of common.Address to represent an address
// on the Kava EVM.
type InternalEVMAddress struct {
	common.Address
}
