# State

## State Objects

The `x/bridge` module keeps the following in state:

| State Object | Description                             | Key                                          | Value                        |
| ------------ | --------------------------------------- | -------------------------------------------- | ---------------------------- |
| ERC20        | Mapping of Ethereum ERC20 to Kava ERC20 | `[]byte{1} + []byte(Ethereum ERC20 address)` | `[]byte{Kava ERC20 address}` |

