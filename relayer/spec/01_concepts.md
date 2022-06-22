# Concepts

## Multi-Party Threshold Signature Scheme

We use [tss-lib] for key generation, signing using secret shares, and key
re-sharing.

This threshold signature scheme enables multi-party signing among n peers such
that any subset of size t + 1 can sign, while any group with t or fewer cannot.

[tss-lib]: https://github.com/bnb-chain/tss-lib
