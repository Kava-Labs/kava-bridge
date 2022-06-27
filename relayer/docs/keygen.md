# Keygen

The relayer requires a TSS 

## Peer and network setup

Before we can run the keygen process, we need the following:

1. Network secret
2. Peer multiaddresses.

### Generate and Distribute Network Secret

1. One peer is designated to generate and distribute the network secret (likely
   Kava Labs). This peer will be referenced as the "dealer."
2. Each peer should generate a new GPG public and private key pair and share the
   public key with the dealer if they are not the dealer.
   1. Generate a GPG key pair.

      ```bash
      gpg --full-generate-key
      ```
   2. Specify which type of key you want (Recommended: `ECC` + `Curve 25519` or `RSA` + `4096`)
   3. Enter your user ID information and a secure passphrase.
   4. List the gpg keys for which you have both a public and private key, then
      copy the long form of the GPG key ID you just generated. In this example,
      the GPG key ID is `DD4617D066BCD508`.

      ```bash
      $ gpg --list-secret-keys --keyid-format=long
      /Users/you/.gnupg/pubring.kbx
      -----------------------------
      sec   ed25519/DD4617D066BCD508 2022-06-24 [SC]
            FE760D3EEFD2116D33832AB0DD4617D066BCD508
      uid                 [ultimate] bridge_test
      ssb   cv25519/DD68E29FC03783D9 2022-06-24 [E]
      ```
   5. Output your GPG public key and share it with the dealer via email, slack, etc.
      This should begin with ` -----BEGIN PGP PUBLIC KEY BLOCK-----` and end with
      `-----END PGP PUBLIC KEY BLOCK-----`.

      ```bash
      $ gpg --armor --export DD4617D066BCD508
      -----BEGIN PGP PUBLIC KEY BLOCK-----

      mDMEYrYRchYJKwYBBAHaRw8BAQdAk0/iFLmMWJIiP63gE9QuolR4TSuGGTraBhNw
      CrT5DNu0C2JyaWRnZV90ZXN0iJQEExYKADwWIQT+dg0+79IRbTODKrDdRhfQZrzV
      CAUCYrYRcgIbAwULCQgHAgMiAgEGFQoJCAsCBBYCAwECHgcCF4AACgkQ3UYX0Ga8
      1Qg7zgD/SeBsrZs29so838eNpGYN//rVXf5bBqzrNaFuUn7E/8wA/AozY1el2OUo
      lSE6gZDNP6tkWO59vjlXnYpRXPEPjCkCuDgEYrYRchIKKwYBBAGXVQEFAQEHQA/b
      mKl4wB9lfcYUrtF6axwL1wGd5oOCGyBEL1LI9tYoAwEIB4h4BBgWCgAgFiEE/nYN
      Pu/SEW0zgyqw3UYX0Ga81QgFAmK2EXICGwwACgkQ3UYX0Ga81QiwZgD+LWPcyYDV
      0M7V0rJRXbdeeSSUvEPImI7eCJOG0AP3gacBAOgxW7eyIMSYrkV0q7uxRL42oIzu
      wu3BecvKlTxM+/EJ
      =pfC6
      -----END PGP PUBLIC KEY BLOCK-----
      ```
3. **Dealer only:** Generate a network secret.
   ```bash
   # Generate network secret.
   kava-relayer network generate-network-secret
   ```

4. **Dealer only:** Encrypt the secret with each other peers' public key and
   share the network secret.
   1. Import all public keys from other peers
      ```bash
      echo "-----BEGIN PGP PUBLIC KEY BLOCK-----

      mDMEYrYRchYJKwYBBAHaRw8BAQdAk0/iFLmMWJIiP63gE9QuolR4TSuGGTraBhNw
      CrT5DNu0C2JyaWRnZV90ZXN0iJQEExYKADwWIQT+dg0+79IRbTODKrDdRhfQZrzV
      CAUCYrYRcgIbAwULCQgHAgMiAgEGFQoJCAsCBBYCAwECHgcCF4AACgkQ3UYX0Ga8
      1Qg7zgD/SeBsrZs29so838eNpGYN//rVXf5bBqzrNaFuUn7E/8wA/AozY1el2OUo
      lSE6gZDNP6tkWO59vjlXnYpRXPEPjCkCuDgEYrYRchIKKwYBBAGXVQEFAQEHQA/b
      mKl4wB9lfcYUrtF6axwL1wGd5oOCGyBEL1LI9tYoAwEIB4h4BBgWCgAgFiEE/nYN
      Pu/SEW0zgyqw3UYX0Ga81QgFAmK2EXICGwwACgkQ3UYX0Ga81QiwZgD+LWPcyYDV
      0M7V0rJRXbdeeSSUvEPImI7eCJOG0AP3gacBAOgxW7eyIMSYrkV0q7uxRL42oIzu
      wu3BecvKlTxM+/EJ
      =pfC6
      -----END PGP PUBLIC KEY BLOCK-----" | gpg --import
   2. List the public key IDs to encrypt with.
      ```bash
      gpg --list-keys
      ```
   3. Encrypt the network secret with each public key, using the corresponding
      key IDs in the list in the previous step. The recipient flag can be used
      multiple times to specify all key IDs at once.
      ```bash
      gpg --encrypt --sign --armor -r FE760D3EEFD2116D33832AB0DD4617D066BCD508 -r 8DEA6D963741D36156C09B80D6D07E98B2C50F41 -r ... psk.key
      ```
   4. Share the encrypted `psk.key.asc` file with the corresponding peers.

5. **Non-dealers only**: Decrypt network secret.
   ```bash
   gpg --decrypt ./path/to/psk.key.asc
   ```

### Distribute multi-addresses

Similar process to distributing the network secret, all peers should exchange
public keys then encrypt and share their full multiaddress containing the
address and peer ID.

Multiaddresses are not sensitive data, but is still best to limit information on
where peers are located. e.g. an external malicious actor knowing peer
multiaddress(es) allows them to have specific servers to target.

## Key Generation

Requirements of key generation are:

1. libp2p node key
2. Pre-computed pre-parameters
3. List of peer multiaddresses

### Generate node key for each peer.

```bash
kava-relayer network generate-node-key
kava-relayer network show-node-id
```

### Pre-compute keygen pre-parameters.

Pre-compute 2 safe primes and Paillier secret required for the protocol. This is
done separately before connecting to other peers as it may take some time and
does not require any communication between peers. All available CPU cores will
be used. This is saved to disk and will be used in the next keygen phase.

```bash
kava-relayer key precompute-preparams
```

### Generate key

TODO:

```bash
kava-relayer join keygen
```