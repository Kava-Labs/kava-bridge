# Reliable Broadcast

The base layer of communication relies on reliable broadcast to ensure all peers
not only receive a message, but also that the message is the same as all other
peers.

This communication layer also supports reliable broadcast for a **subset** of
all peers to selectively send messages for signing rounds.

## Message

Each message contains the following data:

* Unique ID
* Peer ID list of participating peers (e.g. subset of all connected peers, or all of them)
* Payload, protobuf Any data
* Initiating peer ID
* Created timestamp to keep track expire time

## Algorithm

Peer A wants to broadcast x.

1. A sends x to all other peers in the message peer list.
2. Every other peer re-sends x to all other peers in the message peer list
   including peer A.
3. Every peer checks that they received the same values.
4. If any inconsistent values, abort. Otherwise, x is output.
