# Key Generation

To generate the key, all peers can connect to each other with the keygen
command, verify that all correct peers have been connected (denying peers
without an included peer id), and then start the keygen process. This process
will wait until all peers are connected, as  **all** peers are required to
participate in this process.

The keygen process must have the threshold t defined, in addition to the number
of peers n. This means the same threshold value should be shared with all peers
so that they use the same value, and each peer should have the same list of
party IDs that should participate and receive a key part.

At the end of the keygen process, each node will log the associated
public key, ethereum address, and kava address associated with the private key.
Each peer should verify that all other peers have the correct data. If there
is a discrepancy, the process is restarted.

