# Reliable Broadcast

The base layer of communication relies on reliable broadcast to ensure all peers
not only receive a message, but also that the message is the same as all other
peers.

## Message

Each message contains the following data:

* Payload, protobuf Any data
* Signature of peer that initiated broadcast
* Initiating peer ID
* Unique ID (UUID or random 32 bytes?)
* Sequence
* Created timestamp to keep track of TTL / deadline to expire

Initiating peer pubkey is *not* included in the payload, each node can fetch it
directly with the peer ID to verify the message directly.

## Algorithm

Peer A wants to broadcast x.

1. A sends x to all other peers.
2. Every other peer re-sends x to everyone else including peer A.
3. Every peer checks that they received the same values.
4. If any inconsistent values, abort. Otherwise, x is output.

```go
var (
    peerMessages      = make(chan Message)
    confirmedMessages = make(chan Message)
)

// Peer A wants to broadcast "hello world"
for _, peer := range peerList {
    peer.send("hello world")
}

type Message struct {
    ID           string
    SourcePeerID string
    Content      []byte
    Created      time.Time
}

// message id -> peer id -> Message
var messages map[string]map[string]Message
var rejectedMessages map[string]bool

// All nodes including A
func main() {
    for msg := range peerMessages {
        // We can check the message signature against the SourcePeerID pubkey,
        // fetched from node peerstore.
        if err := msg.Validate(); err != nil {
            log.Errorf("invalid message: %v", err)
            continue
        }

        if _, rejected := rejectedMessages[msg.Id]; rejected {
            continue
        }

        peerMsgs, found := messages[msg.ID]
        if !found {
            // First time seeing this message, rebroadcast including to the peer
            // we got the message from. Peer A still needs a response, and
            // it's possible that we receive a message from a different peer
            // before peer A
            for _, peer := range peerList {
                peer.send(msg)
            }

            // Save message to keep track of it later
            messages[msg.ID] = map[string]Message{
                // Peer ID of message can be different from SourcePeerID
                msg.PeerID(): msg
            }
        }

        // Naive check, not entirely necessary to iterate over all messages if
        // they are the same.
        for _, peerMsg := range peerMsgs {
            if !peerMsg.Equal(msg) {
                log.Errorf("message from peer does not match: %v", msg)

                // Abort this message, cannot just delete it from messages map
                // since we want to reject any further same messages from other
                // peers.
                rejectedMessages[msg.ID] = true
                continue
            }
        }

        // Received a message from all other peers AND they all match
        if len(peerMsgs) == len(peerList) {
            confirmedMessages <- msg
        }
    }
}
```
