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
	oldCommittee map[*tss.PartyID]chan ReceivedPartyState
	newCommittee map[*tss.PartyID]chan ReceivedPartyState
}

var _ Transporter = (*MemoryTransporter)(nil)

func NewMemoryTransporter(partyID *tss.PartyID) *MemoryTransporter {
	ts := &MemoryTransporter{
		PartyID:      partyID,
		recvChan:     make(chan ReceivedPartyState, 1),
		sendChan:     make(map[*tss.PartyID]chan ReceivedPartyState),
		oldCommittee: make(map[*tss.PartyID]chan ReceivedPartyState),
		newCommittee: make(map[*tss.PartyID]chan ReceivedPartyState),
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

	dest := routing.To
	if dest == nil {
		return fmt.Errorf("resharing should not have a msg with nil destination")
	}

	if routing.IsToOldCommittee || routing.IsToOldAndNewCommittees {
		log.Debug("sending message to old committee")
		for _, partyID := range dest[:len(mt.oldCommittee)] {
			// Skip sending back to sender
			if partyID == routing.From {
				continue
			}

			ch := mt.oldCommittee[partyID]

			go func(partyID *tss.PartyID, ch chan ReceivedPartyState) {
				log.Debugw("sending message to party", "partyID", partyID, "len(ch)", len(ch))
				ch <- DataRoutingToMessage(data, routing)
				log.Debugw("sent message to party", "partyID", partyID, "len(ch)", len(ch))
			}(partyID, ch)
		}
	}

	if !routing.IsToOldCommittee || routing.IsToOldAndNewCommittees {
		log.Debug("sending message to new committee")
		for _, partyID := range dest {
			// Skip sending back to sender
			if partyID == routing.From {
				continue
			}

			ch := mt.newCommittee[partyID]

			go func(partyID *tss.PartyID, ch chan ReceivedPartyState) {
				log.Debugw("sending message to party", "partyID", partyID, "len(ch)", len(ch))
				ch <- DataRoutingToMessage(data, routing)
				log.Debugw("sent message to party", "partyID", partyID, "len(ch)", len(ch))
			}(partyID, ch)
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

		for partyID, ch := range mt.sendChan {
			// Skip send back to sender
			if partyID == routing.From {
				continue
			}

			go func(partyID *tss.PartyID, ch chan ReceivedPartyState) {
				log.Debugw("sending message to party", "partyID", partyID, "len(ch)", len(ch))
				ch <- DataRoutingToMessage(data, routing)
				log.Debugw("sent message to party", "partyID", partyID, "len(ch)", len(ch))
			}(partyID, ch)
		}

		log.Debug("done broadcast")

		return nil
	}

	for _, partyID := range routing.To {
		if partyID == routing.From {
			continue
		}

		ch, ok := mt.sendChan[partyID]
		if !ok {
			return fmt.Errorf("party %s not found", partyID)
		}

		go func(partyID *tss.PartyID, ch chan ReceivedPartyState) {
			log.Debugw("sending message to party", "partyID", partyID, "len(ch)", len(ch))
			ch <- DataRoutingToMessage(data, routing)
			log.Debugw("sent message to party", "partyID", partyID, "len(ch)", len(ch))
		}(partyID, ch)
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
func (mt *MemoryTransporter) Receive() chan ReceivedPartyState {
	return mt.recvChan
}

func (mt *MemoryTransporter) AddOldCommitteeTarget(partyID *tss.PartyID, ch chan ReceivedPartyState) {
	mt.oldCommittee[partyID] = ch
}

func (mt *MemoryTransporter) AddNewCommitteeTarget(partyID *tss.PartyID, ch chan ReceivedPartyState) {
	mt.newCommittee[partyID] = ch
}

func DataRoutingToMessage(data []byte, routing *tss.MessageRouting) ReceivedPartyState {
	return ReceivedPartyState{
		WireBytes:   data,
		From:        routing.From,
		IsBroadcast: routing.IsBroadcast,
	}
}
