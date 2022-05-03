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

```bash
$ kava-relayer network show-node-id  --p2p.private-key-path ./path/to/node.key
16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9
```

The [multiaddress](https://github.com/multiformats/multiaddr) of a node is used
by peers when connecting to this node via the `kava-relayer connect` command.
To calculate and display the multi-address for a given node (based off the
node id derived from its public key and the port you intend to have the node
listen on) use:

```bash
$ kava-relayer network show-node-multi-address  --p2p.private-key-path node.key --p2p.port 8000
/ip4/127.0.0.1/tcp/8000/p2p/16Uiu2HAm9h4W8Yjt1Z83znprgegVUa1D7UA3i7WYwN5xpGw6BzZ3
```

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
