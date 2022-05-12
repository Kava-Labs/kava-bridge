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
		transports = append(transports, mp_tss.NewMemoryTransporter(partyID))
	}

	t.Logf("transports: %+v", transports)

	// Add transport receivers to each other
	for _, transport := range transports {
		for _, otherTransport := range transports {
			// Skip self
			if transport.PartyID.KeyInt().Cmp(otherTransport.PartyID.KeyInt()) == 0 {
				t.Logf(
					"skipping self for transport: %v == %v",
					transport.PartyID.KeyInt(), otherTransport.PartyID.KeyInt(),
				)
				continue
			}

			transport.AddTarget(otherTransport.PartyID, otherTransport.GetReceiver())
		}
	}

	t.Logf("transports connected: %+v", transports)

	return transports
}
