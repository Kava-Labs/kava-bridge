package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, erc20BridgePairs ERC20BridgePairs) GenesisState {
	return GenesisState{
		Params:           params,
		ERC20BridgePairs: erc20BridgePairs,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:           DefaultParams(),
		ERC20BridgePairs: ERC20BridgePairs{},
	}
}

// Validate validates genesis inputs. It returns error if validation of any input fails.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if err := gs.ERC20BridgePairs.Validate(); err != nil {
		return err
	}

	return nil
}
