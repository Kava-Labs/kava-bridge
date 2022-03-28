#! /bin/bash
set -e

DATADIR=./contrib/devnet/geth/data
BRIDGE_ARTIFACT=./contract/artifacts/contracts/Bridge.sol/Bridge.json
ERC20_ARTIFACT=./contract/artifacts/@openzeppelin/contracts/token/ERC20/ERC20.sol/ERC20.json
GENESIS=./contrib/devnet/geth/genesis.json

# Delete state
rm -rf $DATADIR/geth

# rm -rf $DATADIR
# geth account new --datadir $DATADIR --password ./contrib/devnet/eth-password

# Add bridge contract bytecode to fixed address
# 000000000000000000000000000000000b121d6e
jq -s '.[1].alloc."000000000000000000000000000000000b121d6e" = { code: .[0].deployedBytecode, balance: "0"} |
       .[1]' \
    $BRIDGE_ARTIFACT \
    $GENESIS | sponge $GENESIS

# Add ERC20 contracts

# WETH with same mainnet address
jq -s '.[1].alloc."0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" = { code: .[0].deployedBytecode, balance: "0"} |
       .[1]' \
    $ERC20_ARTIFACT \
    $GENESIS | sponge $GENESIS

# Kava to hex "0x000000000000000000000000000000006b617661"

geth init --datadir $DATADIR $GENESIS

geth --datadir $DATADIR \
     --unlock 21e360e198cde35740e88572b59f2cade421e6b1 \
     --password ./contrib/devnet/eth-password \
     --mine \
     --allow-insecure-unlock \
     --http \
     --http.corsdomain '*'
