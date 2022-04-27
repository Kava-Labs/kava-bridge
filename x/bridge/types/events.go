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
	AttributeKeyReceiver = "receiver"
	AttributeKeyAmount   = "amount"

	// Event Attributes - Bridge
	AttributeKeyEthereumERC20Address = "ethereum_erc20_address"
	AttributeKeyKavaERC20Address     = "kava_erc20_address"
	AttributeKeyRelayer              = "relayer"
	AttributeKeySequence             = "sequence"
	AttributeKeyTxHash               = "tx_hash"

	// Event Attributes - Conversions
	AttributeKeyInitiator    = "initiator"
	AttributeKeyERC20Address = "erc20_address"
)
