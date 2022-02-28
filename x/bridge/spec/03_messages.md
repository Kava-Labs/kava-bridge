# Messages

## Ethereum ERC20 to Kava Transfers

Ethereum ERC20 tokens are transferred with the `MsgERC20FromEthereum` message
type. Only addresses assigned as a permissioned relayer set in [params] may
submit this message, otherwise the transaction will be rejected.

```go
// MsgERC20FromEthereum defines a cross-chain transfer of ERC20 tokens from Ethereum
type MsgERC20FromEthereum struct {
    Relayer              sdk.AccAddress
    EthereumERC20Address []byte
    Amount               sdk.Int
    // **Hex** Kava address, not bech32
    Receiver             []byte
    // Unique sequence per bridge event
    Sequence             sdk.Int
}
```

[params]: 05_params.md
