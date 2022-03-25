#! /bin/bash
set -e

validatorMnemonic="equip town gesture square tomorrow volume nephew minute witness beef rich gadget actress egg sing secret pole winter alarm law today check violin uncover"
# kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c
# 0xeDAA1E944aeAC85a8b2aC41fA741B3b20B783136
faucetMnemonic="crash sort dwarf disease change advice attract clump avoid mobile clump right junior axis book fresh mask tube front require until face effort vault"
# kava1adkm6svtzjsxxvg7g6rshg6kj9qwej8gwqadqd
# 0xeb6dBd418B14A063311e46870BA3569140ECC8e8


DATA=~/.kava-bridged
# remove any old state and config
rm -rf $DATA

BINARY=kava-bridged

# Create new data directory, overwriting any that alread existed
chainID="kavabridgelocalnet_8888-1"
$BINARY init validator --chain-id $chainID
$BINARY config chain-id $chainID

# hacky enable of rest api
sed -in-place='' 's/enable = false/enable = true/g' $DATA/config/app.toml

# Set evm tracer to json
sed -in-place='' 's/tracer = ""/tracer = "json"/g' $DATA/config/app.toml

# avoid having to use password for keys
$BINARY config keyring-backend test

# Create validator keys and add account to genesis
validatorKeyName="validator"
printf "$validatorMnemonic\n" | $BINARY keys add $validatorKeyName --recover
$BINARY add-genesis-account $validatorKeyName 2000000000ukava,100000000000bnb

# Create faucet keys and add account to genesis
faucetKeyName="faucet"
printf "$faucetMnemonic\n" | $BINARY keys add $faucetKeyName --recover
$BINARY add-genesis-account $faucetKeyName 1000000000ukava,100000000000bnb

# Create a delegation tx for the validator and add to genesis
$BINARY gentx $validatorKeyName 1000000000ukava --keyring-backend test --chain-id $chainID
$BINARY collect-gentxs

# Replace stake with ukava
sed -in-place='' 's/stake/ukava/g' $DATA/config/genesis.json

# Replace the default evm denom of aphoton with ukava
sed -in-place='' 's/aphoton/ukava/g' $DATA/config/genesis.json

# Zero out the total supply so it gets recalculated during InitGenesis
jq '.app_state.bank.supply = []' $DATA/config/genesis.json | sponge $DATA/config/genesis.json

# Set relayer to facuet address
jq '.app_state.bridge.params.relayer = "kava1adkm6svtzjsxxvg7g6rshg6kj9qwej8gwqadqd"' $DATA/config/genesis.json | sponge $DATA/config/genesis.json

# Set enabled erc20 tokens to match ropsten testnet
jq '.app_state.bridge.params.enabled_erc20_tokens = [{address: "0xc778417e063141139fce010982780140aa0cd5ab", name: "Wrapped ETH", symbol: "WETH", decimals: 18},{address: "0x07865c6e87b9f70255377e024ace6630c1eaa37f", name: "USDC", symbol: "USDC", decimals: 6}]' $DATA/config/genesis.json | sponge $DATA/config/genesis.json
