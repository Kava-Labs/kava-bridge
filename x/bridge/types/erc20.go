package types

import (
	bytes "bytes"
	"fmt"
)

// ERC20BridgePairs defines a slice of ERC20BridgePair
type ERC20BridgePairs []ERC20BridgePair

// NewERC20BridgePair returns a new ERC20BridgePair from an external and
// internal address.
func NewERC20BridgePair(
	externalAddress ExternalEVMAddress,
	internalAddress InternalEVMAddress,
) ERC20BridgePair {
	return ERC20BridgePair{
		ExternalERC20Address: externalAddress.Address.Bytes(),
		InternalERC20Address: internalAddress.Address.Bytes(),
	}
}

// Validate returns an error if a ERC20BridgePair contains the same address.
func (pair *ERC20BridgePair) Validate() error {
	if bytes.Equal(pair.ExternalERC20Address, pair.InternalERC20Address) {
		return fmt.Errorf("external and internal bytes are same: %x", pair.ExternalERC20Address)
	}

	return nil
}