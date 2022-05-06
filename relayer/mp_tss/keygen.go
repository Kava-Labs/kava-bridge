package mp_tss

import (
	"fmt"

	logging "github.com/ipfs/go-log/v2"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
)

var log = logging.Logger("mp_tss")

func RunKeyGen(params SetupOutput, transport Transporter) (chan keygen.LocalPartySaveData, chan *tss.Error) {
	// outgoing messages to other peers
	outCh := make(chan tss.Message, 10)
	// output data when keygen finished
	endCh := make(chan keygen.LocalPartySaveData, 10)
	// error if keygen fails, contains culprits to blame
	errCh := make(chan *tss.Error, 10)

	log.Debugw("creating new local party")
	party := keygen.NewLocalParty(params.params, outCh, endCh, *params.preParams)
	log.Debugw("local party created", "partyID", party.PartyID().String())

	// Start party in goroutine
	go func() {
		log.Debug("Starting keygen party")
		if err := party.Start(); err != nil {
			errCh <- err
		}
	}()

	// Process outgoing and incoming messages
	go func() {
		log.Debug("Starting out/in message loop")
		for {
			select {
			case outgoingMsg := <-outCh:
				outgoingMsg.GetTo()

				data, routing, err := outgoingMsg.WireBytes()
				log.Debugw("keygen output message", "routing", routing)

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

				log.Debugw("outgoing message done")
			case incomingMsg := <-transport.Receive():
				log.Debugw("received message", "from", incomingMsg.from.String(), "isBroadcast", incomingMsg.isBroadcast)
				ok, err := party.UpdateFromBytes(incomingMsg.wireBytes, incomingMsg.from, incomingMsg.isBroadcast)
				if err != nil {
					log.Errorw("failed to update from bytes", "err", err)
					errCh <- party.WrapError(err)
					return
				}

				log.Debugw("updated party from bytes", "ok", ok)

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
