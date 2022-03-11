# State

## State Objects

The `x/bridge` module keeps the following in state:

| State Object      | Description                             | Key                                          | Value                     |
| ----------------- | --------------------------------------- | -------------------------------------------- | ------------------------- |
| ERC20 Bridge Pair | Mapping of Ethereum ERC20 to Kava ERC20 | `[]byte{1} + []byte(Ethereum ERC20 address)` | `[]byte{ERC20BridgePair}` |

## ERC20 Bridge Pair

One-to-one mapping of bridged Ethereum ERC20 and Kava ERC20 tokens. Tokens
bridged from the pair `ExternalERC20Address` will be minted on the ERC20 at
`InternalERC20Address` on the Kava EVM.

```go
// ERC20BridgePair defines an ERC20 token bridged between external and Kava EVM
type ERC20BridgePair struct {
	// external_erc20_address represents the external EVM ERC20 address
	ExternalERC20Address []byte `protobuf:"bytes,1,opt,name=external_erc20_address,json=externalErc20Address,proto3" json:"external_erc20_address,omitempty"`
	// internal_erc20_address represents the corresponding internal Kava EVM ERC20 address
	InternalERC20Address []byte `protobuf:"bytes,2,opt,name=internal_erc20_address,json=internalErc20Address,proto3" json:"internal_erc20_address,omitempty"`
}
```
