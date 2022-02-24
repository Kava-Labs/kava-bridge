package types

const (
	// ModuleName name used throughout module
	ModuleName = "bridge"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey Top level router key
	RouterKey = ModuleName

	// QuerierRoute should be set to module name
	QuerierRoute = ModuleName

	// QueryPath shortened name for public API (cli and REST)
	QueryPath = ModuleName
)
