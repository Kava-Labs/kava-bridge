#! /bin/bash
set -e

DATADIR=./contrib/devnet/geth/data
GENESIS=./contrib/devnet/geth/genesis.json

# Delete state
rm -rf $DATADIR/geth

# To create the account, if more are needed
# geth account new --datadir $DATADIR --password ./contrib/devnet/eth-password

geth init --datadir $DATADIR $GENESIS

geth --datadir $DATADIR \
     --unlock 21e360e198cde35740e88572b59f2cade421e6b1 \
     --password ./contrib/devnet/eth-password \
     --mine \
     --allow-insecure-unlock \
     --http \
     --http.corsdomain '*' \
     --http.port 8555 \
     --ws.port 8556 &

# Deploy contracts after geth started
sleep 5 \
    && cd ./contract \
    && npx hardhat run scripts/init_dev_env.ts --network localhost &

wait
