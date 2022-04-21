# Peer and Network Initialization

Before key generation and signing can take place, there are a few required
prerequisites:

1. Private / Public key pair for each peer.
2. 256-bit network secret

## 1. Peer Keys

To first initialize a peer, each node must first generate a private and public
key pair.

```bash
kava-relayer network generate-node-key > node.key
```

The peer ID can be then shared with other peers. This not only serves as a
unique identifier for each peer, but a verifiable link between peer and public
key as it is the hash of the peer public key.

```
$ kava-relayer network show-node-id  --p2p.private-key-path ./path/to/node.key
16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9
```

TODO: Easier way to show a full multiaddress to share with other peers.
Currently only shown when using `kava-relayer network connect`.

### Resources

* [Peer Identity][peer-id]
* [Peer ID spec][peer-id-spec]

## 2. Network Secret

The network secret is shared between **all** peers. One peer must generate this
and share it with the others over a secure channel.

```bash
$ kava-relayer network generate-network-secret
zGLVgXcdZK8urz8YGZskJByeciFQhWm4xG1XPBrSNrr6U
```

TODO: A command to send/receive the network secret to the list of other peers?

[peer-id]: https://docs.libp2p.io/concepts/peer-id/
[peer-id-spec]: https://github.com/libp2p/specs/blob/master/peer-ids/peer-ids.md
