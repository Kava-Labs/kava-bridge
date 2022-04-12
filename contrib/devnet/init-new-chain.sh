#! /bin/bash
set -e

validatorMnemonic="equip town gesture square tomorrow volume nephew minute witness beef rich gadget actress egg sing secret pole winter alarm law today check violin uncover"
# kava1ffv7nhd3z6sych2qpqkk03ec6hzkmufy0r2s4c
# 0xeDAA1E944aeAC85a8b2aC41fA741B3b20B783136
faucetMnemonic="crash sort dwarf disease change advice attract clump avoid mobile clump right junior axis book fresh mask tube front require until face effort vault"
# kava1adkm6svtzjsxxvg7g6rshg6kj9qwej8gwqadqd
# 0xeb6dBd418B14A063311e46870BA3569140ECC8e8
userMnemonic="news tornado sponsor drastic dolphin awful plastic select true lizard width idle ability pigeon runway lift oppose isolate maple aspect safe jungle author hole"
# --eth --coin-type 60
# kava10wlnqzyss4accfqmyxwx5jy5x9nfkwh6qm7n4t
# 0x7Bbf300890857b8c241b219C6a489431669b3aFA
relayerMnemonic="never reject sniff east arctic funny twin feed upper series stay shoot vivid adapt defense economy pledge fetch invite approve ceiling admit gloom exit"
# --eth --coin-type 60
# kava15tmj37vh7ch504px9fcfglmvx6y9m70646ev8t
# 0xa2F728F997f62F47D4262a70947F6c36885dF9fa


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

userKeyName="user"
printf "$userMnemonic\n" | $BINARY keys add $userKeyName --recover --eth --coin-type 60
$BINARY add-genesis-account $userKeyName 1000000000ukava

relayerKeyName="relayer"
printf "$relayerMnemonic\n" | $BINARY keys add $relayerKeyName --recover --eth --coin-type 60
$BINARY add-genesis-account $relayerKeyName 1000000000ukava

# Create a delegation tx for the validator and add to genesis
$BINARY gentx $validatorKeyName 1000000000ukava --keyring-backend test --chain-id $chainID
$BINARY collect-gentxs

# Replace stake with ukava
sed -in-place='' 's/stake/ukava/g' $DATA/config/genesis.json

# Replace the default evm denom of aphoton with ukava
sed -in-place='' 's/aphoton/ukava/g' $DATA/config/genesis.json

# Zero out the total supply so it gets recalculated during InitGenesis
jq '.app_state.bank.supply = []' $DATA/config/genesis.json | sponge $DATA/config/genesis.json

# Set relayer to devnet relayer address
jq '.app_state.bridge.params.relayer = "kava15tmj37vh7ch504px9fcfglmvx6y9m70646ev8t"' $DATA/config/genesis.json | sponge $DATA/config/genesis.json

# Set enabled erc20 tokens to match local geth testnet
jq '.app_state.bridge.params.enabled_erc20_tokens = [
    {
        address: "0x6098c27D41ec6dc280c2200A737D443b0AaA2E8F",
        name: "Wrapped ETH",
        symbol: "WETH",
        decimals: 18
    },
    {
        address: "0x4Fb48E68842bb59f07569c623ACa5826b600F8F7",
        name: "USDC",
        symbol: "USDC",
        decimals: 6
    }]' $DATA/config/genesis.json | sponge $DATA/config/genesis.json

# Set enabled conversion pairs - weth address is the first contract bridge module
# deploys
jq '.app_state.bridge.params.enabled_conversion_pairs = [
    {
        kava_erc20_address: "0x404F9466d758eA33eA84CeBE9E444b06533b369e",
        denom: "erc20/weth",
    }]' $DATA/config/genesis.json | sponge $DATA/config/genesis.json
