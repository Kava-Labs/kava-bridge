# Kava Ethereum Bridge

A contract for cross-chain ERC20 token transfers.

## Setup

```
npm install
```

## Test

```
npm test
```

## Development

```
# Watch contract + tests
npm run dev

# Watch tests only
npm run test-watch
```

## Deploy with Hardhat

```
KAVA_BRIDGE_RELAYER_ADDRESS=0x6B1088f788b412Ad1280F95240d56B886A64bc05 npx hardhat run scripts/deploy.ts
```

## Compatibility with Ethermint

Ethermint has their own json format for a compiled contract. The following
converts the abi field to a stringified array, renames bytecode field name to
bin with the leading `0x` trimmed.

```bash
jq '.abi = (.abi | tostring) | {abi, bin: .bytecode[2:] }' < artifacts/contracts/ERC20MintableBurnable.sol/ERC20MintableBurnable.json > ethermint_json/ERC20MintableBurnable.json
```

This is performed by the root Makefile in the `make compile-contracts` command.
