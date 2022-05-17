package mp_tss_test

import (
	"testing"

	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
)

func CreateAndConnectTransports(
	t *testing.T,
	partyIDs []*tss.PartyID,
) []*mp_tss.MemoryTransporter {
	// Create transport between peers
	var transports []*mp_tss.MemoryTransporter
	for _, partyID := range partyIDs {
		transports = append(transports, mp_tss.NewMemoryTransporter(partyID, len(partyIDs)))
	}

	t.Logf("transports: %+v", transports)

	// Add transport receivers to each other
	for _, transport := range transports {
		for _, otherTransport := range transports {
			// Skip self
			// if transport.PartyID.KeyInt().Cmp(otherTransport.PartyID.KeyInt()) == 0 {
			// 	t.Logf(
			// 		"skipping self for transport: %v == %v",
			// 		transport.PartyID.KeyInt(), otherTransport.PartyID.KeyInt(),
			// 	)
			// 	continue
			// }

			transport.AddTarget(otherTransport.PartyID, otherTransport.GetReceiver())
		}
	}

	t.Logf("transports connected: %+v", transports)

	return transports
}

func CreateAndConnectReSharingTransports(
	t *testing.T,
	oldCommittee []*tss.PartyID,
	newCommittee []*tss.PartyID,
) ([]*mp_tss.MemoryTransporter, []*mp_tss.MemoryTransporter) {
	// Create transport between peers
	var oldTransports []*mp_tss.MemoryTransporter
	for _, partyID := range oldCommittee {
		mt := mp_tss.NewMemoryTransporter(partyID, len(oldCommittee)+len(newCommittee))
		oldTransports = append(oldTransports, mt)
	}

	var newTransports []*mp_tss.MemoryTransporter
	for _, partyID := range newCommittee {
		mt := mp_tss.NewMemoryTransporter(partyID, len(oldCommittee)+len(newCommittee))
		newTransports = append(newTransports, mt)
	}

	t.Logf("old transports: %+v", oldTransports)
	t.Logf("new transports: %+v", newTransports)

	// Add old transport receivers to each other
	for _, transport := range oldTransports {
		for _, otherTransport := range oldTransports {
			transport.AddOldCommitteeTarget(otherTransport.PartyID, otherTransport.GetReceiver())
		}

		for _, otherTransport := range newTransports {
			transport.AddNewCommitteeTarget(otherTransport.PartyID, otherTransport.GetReceiver())
		}
	}

	t.Logf("old transports connected")

	// Add new transport receivers to each other
	for _, transport := range newTransports {
		for _, otherTransport := range oldTransports {
			transport.AddOldCommitteeTarget(otherTransport.PartyID, otherTransport.GetReceiver())
		}

		for _, otherTransport := range newTransports {
			transport.AddNewCommitteeTarget(otherTransport.PartyID, otherTransport.GetReceiver())
		}
	}

	t.Logf("new transports connected")

	return oldTransports, newTransports
}
