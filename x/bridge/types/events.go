package types

// Events for the module
const (
	AttributeValueCategory = ModuleName

	// ERC20MintableBurnable event names
	ContractEventTypeWithdraw      = "Withdraw"
	ContractEventTypeConvertToCoin = "ConvertToCoin"

	// Event Types
	EventTypeBridgeEthereumToKava = "bridge_ethereum_to_kava"
	EventTypeBridgeKavaToEthereum = "bridge_kava_to_ethereum"
	EventTypeConvertERC20ToCoin   = "convert_erc20_to_coin"
	EventTypeConvertCoinToERC20   = "convert_coin_to_erc20"

	// Event Attributes - Common
	AttributeKeyEthereumERC20Address = "ethereum_erc20_address"
	AttributeKeyReceiver             = "receiver"
	AttributeKeyAmount               = "amount"

	// Event Attributes - Bridge
	AttributeKeyRelayer  = "relayer"
	AttributeKeySequence = "sequence"

	// Event Attributes - Conversions
	AttributeKeyInitiator    = "initiator"
	AttributeKeyERC20Address = "erc20_address"
)
