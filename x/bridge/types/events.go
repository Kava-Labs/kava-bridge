package types

// Events for the module
const (
	EventTypeDeposit  = "bridge_deposit"
	EventTypeWithdraw = "bridge_withdraw"

	// ERC20MintableBurnable event names
	ContractEventTypeWithdraw = "Withdraw"

	AttributeValueCategory           = ModuleName
	AttributeKeyEthereumERC20Address = "ethereum_erc20_address"
	AttributeKeyRelayer              = "relayer"
	AttributeKeyReceiver             = "receiver"
	AttributeKeySequence             = "sequence"
)
