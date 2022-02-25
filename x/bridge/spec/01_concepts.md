# Concepts

The bridge module deploys an ERC-20 contract on the Kava EVM for cross-chain
asset ERC-20 transfers.

## ERC20

In the following documents, Ethereum ERC20 will refer to an ERC20 token deployed
on the Ethereum network. Kava ERC20 will refer to an ERC20 token deployed on the
Kava EVM.

## Requirements

Signer

* There must be trusted signer(s) that watch the bridge contract events for
  locked assets.
