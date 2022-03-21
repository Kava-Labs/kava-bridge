package types

import (
	bytes "bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
)

// Parameter keys and default values
var (
	KeyEnabledERC20Tokens     = []byte("EnabledERC20Tokens")
	KeyRelayer                = []byte("Relayer")
	KeyEnabledConversionPairs = []byte("EnabledConversionPairs")
	DefaultEnabledERC20Tokens = EnabledERC20Tokens{}
	DefaultRelayer            = sdk.AccAddress{}
	DefaultConversionPairs    = ConversionPairs{}
)

// ParamKeyTable for bridge module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value
// pairs pairs of the bridge module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnabledERC20Tokens, &p.EnabledERC20Tokens, validateEnabledERC20Tokens),
		paramtypes.NewParamSetPair(KeyRelayer, &p.Relayer, validateRelayer),
		paramtypes.NewParamSetPair(KeyEnabledConversionPairs, &p.EnabledConversionPairs, validateConversionPairs),
	}
}

// NewParams returns new bridge module Params.
func NewParams(enabledERC20Tokens EnabledERC20Tokens, relayer sdk.AccAddress, conversionPairs ConversionPairs) Params {
	return Params{
		EnabledERC20Tokens:     enabledERC20Tokens,
		Relayer:                relayer,
		EnabledConversionPairs: conversionPairs,
	}
}

// DefaultParams returns the default parameters for bridge.
func DefaultParams() Params {
	return NewParams(DefaultEnabledERC20Tokens, DefaultRelayer, DefaultConversionPairs)
}

func (p *Params) Validate() error {
	if err := p.EnabledERC20Tokens.Validate(); err != nil {
		return err
	}

	if p.Relayer == nil {
		return errors.New("relayer cannot be nil")
	}

	if err := p.EnabledConversionPairs.Validate(); err != nil {
		return err
	}

	return nil
}

// validateEnabledERC20Tokens validates an interface as EnabledERC20Tokens
func validateEnabledERC20Tokens(i interface{}) error {
	enabledERC20Tokens, ok := i.(EnabledERC20Tokens)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return enabledERC20Tokens.Validate()
}

// EnabledERC20Tokens defines a slice of EnabledERC20Token
type EnabledERC20Tokens []EnabledERC20Token

// NewEnabledERC20Tokens returns EnabledERC20Tokens from the provided values
func NewEnabledERC20Tokens(enabledERC20Tokens ...EnabledERC20Token) EnabledERC20Tokens {
	return EnabledERC20Tokens(enabledERC20Tokens)
}

// Validate returns an error if any token in a slice of EnabledERC20Tokens is
// invalid.
func (tokens EnabledERC20Tokens) Validate() error {
	// Check for duplicates
	addrs := map[string]bool{}

	for _, token := range tokens {
		if addrs[hex.EncodeToString(token.Address)] {
			return fmt.Errorf(
				"found duplicate enabled ERC20 token address %s",
				hex.EncodeToString(token.Address),
			)
		}

		if err := token.Validate(); err != nil {
			return err
		}

		addrs[hex.EncodeToString(token.Address)] = true
	}

	return nil
}

// NewEnabledERC20Token returns a new EnabledERC20Token.
func NewEnabledERC20Token(address ExternalEVMAddress, name string, symbol string, decimals uint32) EnabledERC20Token {
	return EnabledERC20Token{
		Address:  address.Bytes(),
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
}

// Validate returns an error if the EnabledERC20Token is invalid.
func (e EnabledERC20Token) Validate() error {
	if len(e.Address) != common.AddressLength {
		return fmt.Errorf("address length is %v but expected %v", len(e.Address), common.AddressLength)
	}

	if bytes.Equal(e.Address, common.Address{}.Bytes()) {
		return fmt.Errorf("address cannot be zero value %v", hex.EncodeToString(e.Address))
	}

	if e.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if e.Symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}

	// Decimals being 0 is technically allowed in ERC20, but unlikely
	if e.Decimals == 0 {
		return fmt.Errorf("decimals cannot be 0")
	}

	// Check size since the go type is uint32 from proto int, but the actual
	// size of erc20 decimals should be an uint8
	if e.Decimals > math.MaxUint8 {
		return fmt.Errorf("decimals is too large, max %v", math.MaxUint8)
	}

	return nil
}

// validateRelayer validates a relayer address
func validateRelayer(i interface{}) error {
	relayerAddr, ok := i.(sdk.AccAddress)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if relayerAddr == nil {
		return errors.New("relayer address cannot be nil")
	}

	return nil
}
