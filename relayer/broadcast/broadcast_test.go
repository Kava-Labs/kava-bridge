package broadcast_test

import (
	"context"
	"testing"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/require"
)

func TestBroadcast_Connect(t *testing.T) {
	logging.SetAllLoggers(logging.LevelDebug)

	ctx, cancel := context.WithCancel(context.Background())

	hosts := testutil.CreateHosts(t, ctx, 2)

	b0 := newBroadcast(ctx, hosts[0])
	b1 := newBroadcast(ctx, hosts[1])

	testutil.Connect(t, hosts[0], hosts[1])

	time.Sleep(time.Second)

	b0PeerCount := b0.GetPeerCount()
	b1PeerCount := b1.GetPeerCount()

	require.Equal(t, 1, b0PeerCount)
	require.Equal(t, 1, b1PeerCount)

	cancel()
}

func newBroadcast(ctx context.Context, h host.Host) *broadcast.Broadcast {
	return broadcast.NewBroadcast(ctx, h)
}
