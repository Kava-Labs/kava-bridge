# Relayer

## Usage

1. Generate network secret. This should be only done by one peer and shared with
   the other peers over a secure medium.

```bash
# Generate network secret.
kava-relayer network generate-network-secret
```

2. Generate node key for each peer.

```bash
kava-relayer network generate-node-key
kava-relayer network show-node-id
```

3. Pre-compute keygen pre-parameters.

Pre-compute 2 safe primes and Paillier secret required for the protocol. This is
done separately before connecting to other peers as it may take some time and
does not require any communication between peers. All available CPU cores will
be used. This is saved to disk and will be used in the next keygen phase.

```bash
kava-relayer key precompute-preparams
```

4. Start the relayer, connect to other peers.

```bash
# Connect to p2p network
kava-relayer network connect

# Start single signer relayer without P2P network
kava-relayer start 
```
