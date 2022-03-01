# Kava Etherium Bridge

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

Ethermint has their own json format for a compiled contract

```bash
cat artifacts/contracts/ERC20MintableBurnable.sol/ERC20MintableBurnable.json | jq '.abi = (.abi | tostring) | {abi, bin: .bytecode[2:] }' > ethermint_json/ERC20MintableBurnable.json
```
