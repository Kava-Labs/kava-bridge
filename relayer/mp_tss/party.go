package mp_tss

import (
	"github.com/binance-chain/tss-lib/tss"
)

// RunParty starts the local party in the background and handles incoming and
// outgoing messages. Does **not** block.
func RunParty(
	party tss.Party,
	errCh chan<- *tss.Error,
	outCh <-chan tss.Message,
	transport Transporter,
	isReSharing bool,
) {
	// Start party in goroutine
	go func() {
		log.Debug("Starting party")
		if err := party.Start(); err != nil {
			errCh <- err
		}
	}()

	// Process outgoing and incoming messages
	go func() {
		incomingMsgCh := transport.Receive()

		log.Debug("Starting out/in message loop")
		for {
			log.Debugw("waiting for next message...", "partyID", party.PartyID())
			select {
			case outgoingMsg := <-outCh:
				log.Debugw("outgoing message", "GetTo()", outgoingMsg.GetTo())

				data, routing, err := outgoingMsg.WireBytes()
				log.Debugw(
					"party outgoing msg write bytes",
					"partyID", party.PartyID(),
					"routing", routing,
				)

				if err != nil {
					errCh <- party.WrapError(err)
					return
				}

				// Prevent blocking goroutine to receive messages, may deadlock
				// if receive channels are full if not in goroutine.
				go func() {
					// send to other parties
					if err := transport.Send(data, routing, isReSharing); err != nil {
						log.Errorw(
							"failed to send output message",
							"from PartyID", party.PartyID(),
							"err", err,
						)
						errCh <- party.WrapError(err)
						return
					}

					log.Debugw("done sending outgoing message", "partyID", party.PartyID())
				}()
			case incomingMsg := <-incomingMsgCh:
				// Running in goroutine prevents blocking when channels get
				// filled up. This may deadlock if not run in a goroutine and
				// blocks receiving messages.
				go func() {
					log.Debugw(
						"received message",
						"partyID", party.PartyID(),
						"from partyID", incomingMsg.from,
						"isBroadcast", incomingMsg.isBroadcast,
						"len(bytes)", len(incomingMsg.wireBytes),
					)

					// The first return value `ok` is false only when there is
					// an error. This should be fine to ignore as we handle err
					// instead.
					_, err := party.UpdateFromBytes(
						incomingMsg.wireBytes,
						incomingMsg.from,
						incomingMsg.isBroadcast,
					)
					if err != nil {
						log.Errorw("failed to update from bytes", "err", err)
						errCh <- party.WrapError(err)
						return
					}

					log.Debugw(
						"updated party from bytes",
						"partyID", party.PartyID(),
					)
				}()
			}
		}
	}()
}
