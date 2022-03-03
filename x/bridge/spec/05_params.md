# Parameters

The bridge module contains the following parameters:

| Key                | Type               | Example                                         | Description                             |
| ------------------ | ------------------ | ----------------------------------------------- | --------------------------------------- |
| EnabledERC20Tokens | EnabledERC20Tokens | `[]EnabledERC20Token`                           | array of ERC20 tokens allowed to bridge |
| Relayer            | sdk.AccAddress     | `"kava123fxg0l602etulhhcdm0vt7l57qya5wjcrwhzz"` | bridge relayer address                  |

Each EnabledERC20Token has the following parameters:

| Key                        | Type   | Example                                      | Description               |
| -------------------------- | ------ | -------------------------------------------- | ------------------------- |
| EnabledERC20Token.Address  | bytes  | `0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2` | ERC20 address on Ethereum |
| EnabledERC20Token.Name     | string | `"Wrapped Ether"`                            | ERC20 token name          |
| EnabledERC20Token.Symbol   | string | `"WETH"`                                     | ERC20 token symbol        |
| EnabledERC20Token.Decimals | uint8  | `18`                                         | ERC20 token decimals      |

Governance param change proposals are used to add new Ethereum ERC20s to the
enabled list. Ethereum ERC20s that are not in the list are rejected from
being bridged to Kava.
