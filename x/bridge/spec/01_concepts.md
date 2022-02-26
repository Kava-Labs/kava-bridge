# Concepts

The bridge module deploys an ERC-20 contract on the Kava EVM for cross-chain
asset ERC-20 transfers.

## ERC20

In the following documents, Ethereum ERC20 will refer to an ERC20 token deployed
on the Ethereum network. Kava ERC20 will refer to an ERC20 token deployed on the
Kava EVM.

## Requirements

Signer

* There must be trusted a signer that watches the bridge smart contract on
  Ethereum for locked asset events. This is a single signer for now.

## Ethereum ERC20 to Kava Transfers

```mermaid
stateDiagram-v2
    state Ethereum {
        User --> Contract: Lock(Ethereum ERC20 Addr, toAddr, amount)
    }
    
    Contract --> Relayer
    Relayer --> BridgeModule: MsgERC20FromEthereum

    state Kava {
        Reject: Reject TX

        state if_has_permission <<choice>>
        state if_erc20_deployed <<choice>>
        state if_erc20_whitelisted <<choice>>

        BridgeModule --> if_has_permission
        if_has_permission --> Reject: Unknown Signer
        if_has_permission --> if_erc20_whitelisted: Signer in Params

        if_erc20_whitelisted --> Reject: Unknown Ethereum ERC20 Addr
        if_erc20_whitelisted --> if_erc20_deployed: Ethereum ERC20 Addr in Params

        DeployERC20: Deploy ERC20
        MintERC20: Mint ERC20 amount for toAddr
        if_erc20_deployed --> DeployERC20: Kava ERC20 not deployed
        if_erc20_deployed --> MintERC20: Kava ERC20 exists
        DeployERC20 --> MintERC20
    }
```

