package types

import (
	bytes "bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// ERC20BridgePairs defines a slice of ERC20BridgePair
type ERC20BridgePairs []ERC20BridgePair

func NewERC20BridgePairs(pairs ...ERC20BridgePair) ERC20BridgePairs {
	return pairs
}

func (pairs ERC20BridgePairs) Validate() error {
	// Check for duplicates for both internal/external addrs
	internalAddrs := map[string]bool{}
	externalAddrs := map[string]bool{}

	for _, pair := range pairs {
		if err := pair.Validate(); err != nil {
			return err
		}

		if internalAddrs[pair.GetInternalAddress().String()] {
			return fmt.Errorf(
				"found duplicate enabled bridge pair internal address %s",
				pair.GetInternalAddress().String(),
			)
		}

		if externalAddrs[pair.GetExternalAddress().String()] {
			return fmt.Errorf(
				"found duplicate enabled bridge pair external address %s",
				pair.GetExternalAddress().String(),
			)
		}

		internalAddrs[pair.GetInternalAddress().String()] = true
		externalAddrs[pair.GetExternalAddress().String()] = true
	}

	return nil
}

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
	if len(pair.ExternalERC20Address) != common.AddressLength {
		return fmt.Errorf(
			"external address length is %v but expected %v",
			len(pair.ExternalERC20Address),
			common.AddressLength,
		)
	}

	if len(pair.InternalERC20Address) != common.AddressLength {
		return fmt.Errorf(
			"internal address length is %v but expected %v",
			len(pair.InternalERC20Address),
			common.AddressLength,
		)
	}

	if bytes.Equal(pair.ExternalERC20Address, common.Address{}.Bytes()) {
		return fmt.Errorf("external address cannot be zero value %v", pair.ExternalERC20Address)
	}

	if bytes.Equal(pair.InternalERC20Address, common.Address{}.Bytes()) {
		return fmt.Errorf("internal address cannot be zero value %v", pair.InternalERC20Address)
	}

	if bytes.Equal(pair.ExternalERC20Address, pair.InternalERC20Address) {
		return fmt.Errorf("external and internal bytes are same: %x", pair.ExternalERC20Address)
	}

	return nil
}

// GetID returns the SHA256 hash of the external and internal address
func (pair *ERC20BridgePair) GetID() []byte {
	return tmhash.Sum(append(pair.ExternalERC20Address, pair.InternalERC20Address...))
}

// GetExternal returns the typed ExternalAddress.
func (pair *ERC20BridgePair) GetExternalAddress() ExternalEVMAddress {
	return NewExternalEVMAddress(
		common.BytesToAddress(pair.ExternalERC20Address),
	)
}

// GetInternalAddress returns the typed InternalAddress.
func (pair *ERC20BridgePair) GetInternalAddress() InternalEVMAddress {
	return NewInternalEVMAddress(
		common.BytesToAddress(pair.InternalERC20Address),
	)
}
