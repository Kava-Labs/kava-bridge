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
	// old/new committee only for resigning
	oldCommittee map[*tss.PartyID]struct{}
	newCommittee map[*tss.PartyID]struct{}
}

var _ Transporter = (*MemoryTransporter)(nil)

func NewMemoryTransporter(partyID *tss.PartyID, bufSize int) *MemoryTransporter {
	ts := &MemoryTransporter{
		PartyID: partyID,
		// Much faster with buf size more than 1
		recvChan:     make(chan ReceivedPartyState, bufSize),
		sendChan:     make(map[*tss.PartyID]chan ReceivedPartyState),
		oldCommittee: make(map[*tss.PartyID]struct{}),
		newCommittee: make(map[*tss.PartyID]struct{}),
	}

	return ts
}

func (mt *MemoryTransporter) Send(data []byte, routing *tss.MessageRouting, isResharing bool) error {
	if isResharing {
		return mt.sendReSharing(data, routing)
	}

	return mt.sendKeygenOrSigning(data, routing)
}

func (mt *MemoryTransporter) sendReSharing(data []byte, routing *tss.MessageRouting) error {
	log.Debugw(
		"sending resharing message",
		"to", routing.To,
		"isBroadcast", routing.IsBroadcast,
		"IsToOldCommittee", routing.IsToOldCommittee,
		"IsToOldAndNewCommittees", routing.IsToOldAndNewCommittees,
	)

	if routing.IsToOldCommittee || routing.IsToOldAndNewCommittees {
		log.Debug("sending message to old committee")
		for partyID := range mt.oldCommittee {
			ch, ok := mt.sendChan[partyID]
			if !ok {
				return fmt.Errorf("old committee party %s not found", partyID)
			}
			ch <- DataRoutingToMessage(data, routing)
		}
	}

	if !routing.IsToOldCommittee || routing.IsToOldAndNewCommittees {
		log.Debug("sending message to new committee")
		for partyID := range mt.newCommittee {
			ch, ok := mt.sendChan[partyID]
			if !ok {
				return fmt.Errorf("new committee party %s not found", partyID)
			}
			ch <- DataRoutingToMessage(data, routing)
		}
	}

	return nil
}

func (mt *MemoryTransporter) sendKeygenOrSigning(data []byte, routing *tss.MessageRouting) error {
	log.Debugw(
		"sending message",
		"to", routing.To,
		"isBroadcast", routing.IsBroadcast,
	)

	if routing.IsBroadcast && len(routing.To) == 0 {
		log.Debug("broadcast message to all peers")

		for party, ch := range mt.sendChan {
			log.Debugw("sending message to party", "partyID", party)
			ch <- DataRoutingToMessage(data, routing)
			log.Debugw("sent message to party", "partyID", party)
		}

		log.Debug("done broadcast")

		return nil
	}

	for _, partyID := range routing.To {
		ch, ok := mt.sendChan[partyID]
		if !ok {
			return fmt.Errorf("party %s not found", partyID)
		}

		log.Debugw("sending message to party", "partyID", partyID, "len(ch)", len(ch))
		ch <- DataRoutingToMessage(data, routing)
		log.Debugw("sent message to party", "partyID", partyID)
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

func (mt *MemoryTransporter) AddOldCommittee(partyIDs ...*tss.PartyID) {
	for _, partyID := range partyIDs {
		mt.oldCommittee[partyID] = struct{}{}
	}
}

func (mt *MemoryTransporter) AddNewCommittee(partyIDs ...*tss.PartyID) {
	for _, partyID := range partyIDs {
		mt.newCommittee[partyID] = struct{}{}
	}
}

func DataRoutingToMessage(data []byte, routing *tss.MessageRouting) ReceivedPartyState {
	return ReceivedPartyState{
		wireBytes:   data,
		from:        routing.From,
		isBroadcast: routing.IsBroadcast,
	}
}
