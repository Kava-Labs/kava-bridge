# Events

The `x/bridge` module emits the following events:

## Handlers

### MsgERC20FromEthereum

Note that the receiver address is in hex format, not bech32.

| Type         | Attribute Key          | Attribute Value               |
| ------------ | ---------------------- | ----------------------------- |
| bridge_erc20 | ethereum_erc20_address | `{ethereum erc20 address}`    |
| bridge_erc20 | amount                 | `{amount}`                    |
| bridge_erc20 | receiver               | `{kava hex receiver address}` |
| bridge_erc20 | sequence               | `{lock sequence}`             |
| message      | module                 | bridge                        |
| message      | sender                 | `{relayer address}`           |
