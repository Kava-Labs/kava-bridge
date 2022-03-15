# Concepts

The bridge module deploys ERC20 contracts and mints ERC20 tokens on the Kava EVM
for cross-chain ERC-20 token transfers.

## Requirements

Bridge Contract

* The bridge contract must be deployed on Ethereum.

Signer

* There must be trusted a signer that watches the bridge smart contract on
  Ethereum for locked asset events. This is a single signer for now.

## ERC20

In the following documents, Ethereum ERC20 will refer to an ERC20 token deployed
on the Ethereum network. Kava ERC20 will refer to an ERC20 token deployed on the
Kava EVM.

## Sequence

There are two different incrementing sequences, one on the Bridge contract and
`x/bridge` module. These are unique for each deposit (Ethereum to Kava) and
each withdraw (Kava to Ethereum), but are not unique in that a deposit sequence
value can be the same as a withdraw sequence. This is used by the relayer to
properly order transactions.

## Ethereum ERC20 to Kava Transfers

Before being able to bridge Ethereum ERC20 tokens, they need to be added to the
enabled ERC20 tokens in params.

In order to bridge an approved Ethereum ERC20 tokens to Kava, the following
steps are taken:

1. Account locks ERC20 tokens in the bridge contract on Ethereum. This emits an
   event with the Ethereum ERC20 address, Ethereum sender address, receiver Kava
   address, amount, and sequence.
2. After a reasonable number of confirmations, the relayer will sign and submit
   a `MsgERC20FromEthereum` message to the Kava chain.
3. The bridge module will verify the message for the following conditions. If
   any of these are false, the transaction will be rejected.
   * The signer address of the message matches the one set in params.
   * The Ethereum ERC20 token is contained in the enabled list set in params.
4. The target Kava ERC20 address is fetched in the module state. If it doesn't
   exist in state, i.e. the Kava ERC20 contract does not exist, it is deployed.
5. The bridge module mints Kava ERC20 tokens for the destination Kava address.

```mermaid
stateDiagram-v2
    state Ethereum {
        User --> Contract: Lock(Ethereum ERC20 Addr, receiver, amount)
    }
    
    Contract --> Relayer
    Relayer --> BridgeModule: MsgERC20FromEthereum

    state Kava {
        Reject: Reject TX
        BridgeModule: Bridge Module

        state if_has_permission <<choice>>
        state if_erc20_deployed <<choice>>
        state if_erc20_enabled <<choice>>

        BridgeModule --> if_has_permission
        if_has_permission --> Reject: Unknown Signer
        if_has_permission --> if_erc20_enabled: Signer in Params

        if_erc20_enabled --> Reject: Unknown Ethereum ERC20 Addr
        if_erc20_enabled --> if_erc20_deployed: Ethereum ERC20 Addr in Params

        DeployERC20: Deploy ERC20
        MintERC20: Mint ERC20 amount for receiver
        if_erc20_deployed --> DeployERC20: Kava ERC20 not deployed
        if_erc20_deployed --> MintERC20: Kava ERC20 exists
        DeployERC20 --> MintERC20
    }
```

## Kava ERC20 to Ethereum Transfers

Transferring from Kava to Ethereum follows a similar pattern. Of the following
steps, only step 1 is implemented in the bridge module and the subsequent steps
are done by the relayer.

1. Account calls `Withdraw(withdrawal Ethereum Address, amount)` on a Kava ERC20
   contract. This burns the account tokens and emits a `Withdraw` event
   containing the receiver Ethereum address and corresponding amount.
2. Module `PostTxProcessing` EVM hook scans for `Withdraw` events from enabled
   `ERC20BridgePair`s and emits a [Cosmos SDK Event][cosmos-event] with the
   withdrawal Ethereum address, corresponding ERC20 address to send funds to,
   amount, and unique withdraw sequence.

   **Note:** If the Withdraw event comes from a contract that isn't in bridge
   state (`EnabledERC20Tokens`), then it is ignored. These events may come from
   a contract that isn't deployed by the bridge module. Arbitrary contracts
   cannot maliciously try to get funds withdrawn this way as the withdrawal
   Ethereum ERC20 address is queried from from module params, not from contract
   events.
3. When Relayer queries a new Withdraw bridge module event, unlock funds on the
   Ethereum bridge contract.

```mermaid
sequenceDiagram
    participant acc as Kava Account
    participant KRC20 as Kava ERC20
    participant M as x/bridge Module
    participant R as Relayer
    participant B as Ethereum Bridge Contract

    acc->>KRC20: Withdraw(Ethereum toAddr, amount)
    KRC20->>KRC20: Burn ERC20 amount
    KRC20->>M: Emit Withdraw(Ethereum toAddr, amount)

    M->>M: Emit Cosmos WithdrawEvent(Ethereum ERC20 address, toAddr, amount, sequence)

    loop Monitor Bridge Events
        R->>+M: Get Cosmos WithdrawEvents
        M-->>-R: Cosmos WithdrawEvents
    end

    R->>B: Unlock(Ethereum ERC20 address, toAddr, amount)
```

[cosmos-event]: https://docs.cosmos.network/master/core/events.html
