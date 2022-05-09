package mp_tss

import (
	"fmt"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
)

var log = logging.Logger("mp_tss")

func RunKeyGen(
	preParams *keygen.LocalPreParams,
	params *tss.Parameters,
	transport Transporter,
) (chan keygen.LocalPartySaveData, chan *tss.Error) {
	// outgoing messages to other peers
	outCh := make(chan tss.Message, 10)
	// output data when keygen finished
	endCh := make(chan keygen.LocalPartySaveData, 10)
	// error if keygen fails, contains culprits to blame
	errCh := make(chan *tss.Error, 10)

	log.Debugw("creating new local party")
	party := keygen.NewLocalParty(params, outCh, endCh, *preParams)

	log := log.Named(party.PartyID().String())

	log.Debugw("local party created", "partyID", party.PartyID().String())

	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Debugw(
				"party waiting for",
				"party index", party.PartyID().Index,
				"waitingFor()", party.WaitingFor(),
			)
		}
	}()

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

	return endCh, errCh
}
