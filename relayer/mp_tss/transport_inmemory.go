package mp_tss

import (
	"fmt"

	"github.com/binance-chain/tss-lib/tss"
)

type MemoryTransporter struct {
	PartyID *tss.PartyID
	// incoming messages from other parties.
	recvChan chan ReceivedPartyState
	// outgoing messages to other parties
	sendChan map[*tss.PartyID]chan ReceivedPartyState
}

var _ Transporter = (*MemoryTransporter)(nil)

func NewMemoryTransporter(partyID *tss.PartyID) *MemoryTransporter {
	return &MemoryTransporter{
		PartyID:  partyID,
		recvChan: make(chan ReceivedPartyState, 1),
		sendChan: make(map[*tss.PartyID]chan ReceivedPartyState),
	}
}

func (mt *MemoryTransporter) Send(data []byte, routing *tss.MessageRouting) error {
	log.Debugw("sending message", "to", routing.To, "isBroadcast", routing.IsBroadcast)

	if routing.IsBroadcast && len(routing.To) != 0 {
		return fmt.Errorf("cannot send broadcast message to a specific party")
	}

	if routing.IsBroadcast && len(routing.To) == 0 {
		log.Debug("broadcast message to all peers")

		for party, ch := range mt.sendChan {
			log.Debugw("sending message to party", "partyID", party.String())
			ch <- ReceivedPartyState{
				wireBytes:   data,
				from:        routing.From,
				isBroadcast: true,
			}
			log.Debugw("sent message to party", "partyID", party.String())
		}

		log.Debug("done broadcast")

		return nil
	}

	for _, party := range routing.To {
		ch, ok := mt.sendChan[party]
		if !ok {
			return fmt.Errorf("party %s not found", party.String())
		}

		log.Debugw("sending message to party", "partyID", party.String())
		ch <- ReceivedPartyState{
			wireBytes:   data,
			from:        routing.From,
			isBroadcast: true,
		}
		log.Debugw("sent message to party", "partyID", party.String())
	}

	return nil
}

func (mt *MemoryTransporter) AddTarget(partyID *tss.PartyID, ch chan ReceivedPartyState) {
	mt.sendChan[partyID] = ch
}

// GetReceiver returns a channel for other peer to send messages to.
func (mt *MemoryTransporter) GetReceiver() chan ReceivedPartyState {
	return mt.recvChan
}

// Receive returns a channel for the current peer to receive messages from
// other peers.
func (mt *MemoryTransporter) Receive() <-chan ReceivedPartyState {
	return mt.recvChan
}
