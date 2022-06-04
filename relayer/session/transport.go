package session

import (
	"context"
	"fmt"

	"github.com/binance-chain/tss-lib/tss"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/mp_tss"
	mp_tss_types "github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/sync/errgroup"
)

// SessionTransport is a transport for a specific session.
type SessionTransport struct {
	broadcaster  *broadcast.Broadcaster
	partyIDStore *mp_tss.PartyIDStore
	sessionID    mp_tss_types.AggregateSigningSessionID
	participants []peer.ID

	recvChan chan mp_tss.ReceivedPartyState
}

var _ mp_tss.Transporter = (*SessionTransport)(nil)

func (mt *SessionTransport) Send(data []byte, routing *tss.MessageRouting, isResharing bool) error {
	// TODO: Implement broadcast resharing
	if isResharing {
		return nil
	}

	msg := mp_tss_types.NewSigningPartMessage(mt.sessionID, data, routing.IsBroadcast)

	if routing.IsBroadcast {
		return mt.broadcaster.BroadcastMessage(
			context.Background(),
			&msg,
			mt.participants,
			30,
		)
	}

	// Point to point concurrently
	// TODO: Might not be necessary, routing.To may only consist of one peer.
	g, ctx := errgroup.WithContext(context.Background())

	for _, to := range routing.To {
		peerID, found := mt.partyIDStore.GetPeerID(to)
		if !found {
			return fmt.Errorf("peer %s not found", to)
		}

		g.Go(func() error {
			err := mt.broadcaster.BroadcastMessage(
				ctx,
				&msg,
				[]peer.ID{peerID},
				30,
			)

			if err != nil {
				return fmt.Errorf("failed to send message to peer: %w", err)
			}

			return nil
		})
	}

	return g.Wait()
}

func (mt *SessionTransport) Receive() chan mp_tss.ReceivedPartyState {
	return mt.recvChan
}

func NewSessionTransport(
	broadcaster *broadcast.Broadcaster,
	sessionID mp_tss_types.AggregateSigningSessionID,
) mp_tss.Transporter {
	return &SessionTransport{
		broadcaster: broadcaster,
		sessionID:   sessionID,
		recvChan:    make(chan mp_tss.ReceivedPartyState, 1),
	}
}
