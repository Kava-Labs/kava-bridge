package mp_tss

import (
	"fmt"

	"github.com/binance-chain/tss-lib/tss"
)

// RunParty starts the local party in the background and handles incoming and
// outgoing messages. Does **not** block.
func RunParty(
	party tss.Party,
	errCh chan *tss.Error,
	outCh chan tss.Message,
	transport Transporter,
) {
	// Start party in goroutine
	go func() {
		log.Debug("Starting keygen party")
		if err := party.Start(); err != nil {
			errCh <- err
		}
	}()

	// Process outgoing and incoming messages
	go func() {
		incomingMsgCh := transport.Receive()

		log.Debug("Starting out/in message loop")
		for {
			log.Debugw("waiting for next message...", "party index", party.PartyID().Index)
			select {
			case outgoingMsg := <-outCh:
				log.Debugw("outgoing message", "GetTo()", outgoingMsg.GetTo())

				data, routing, err := outgoingMsg.WireBytes()
				log.Debugw(
					"keygen outgoing msg write bytes",
					"party index", party.PartyID().Index,
					"routing", routing,
				)

				if err != nil {
					errCh <- party.WrapError(err)
					return
				}

				// send to other parties
				if err := transport.Send(data, routing); err != nil {
					log.Errorw("failed to send output message", "err", err)
					errCh <- party.WrapError(err)
					return
				}

				log.Debugw("outgoing message done", "party index", party.PartyID().Index)
			case incomingMsg := <-incomingMsgCh:
				log.Debugw(
					"received message",
					"party index", party.PartyID().Index,
					"from index", incomingMsg.from.Index,
					"isBroadcast", incomingMsg.isBroadcast,
					"len(bytes)", len(incomingMsg.wireBytes),
				)

				ok, err := party.UpdateFromBytes(incomingMsg.wireBytes, incomingMsg.from, incomingMsg.isBroadcast)
				if err != nil {
					log.Errorw("failed to update from bytes", "err", err)
					errCh <- party.WrapError(err)
					return
				}

				log.Debugw(
					"updated party from bytes",
					"party index", party.PartyID().Index,
					"ok", ok,
				)

				// TODO: What does ok mean?
				if !ok {
					log.Errorw("failed to update party from bytes")
					errCh <- party.WrapError(fmt.Errorf("keygen update returned not ok"))
					return
				}
			}
		}
	}()
}
