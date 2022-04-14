package broadcast_test

import (
	"context"
	"testing"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcast_Connect(t *testing.T) {
	logging.SetAllLoggers(logging.LevelDebug)

	ctx, cancel := context.WithCancel(context.Background())

	count := 5
	hosts := testutil.CreateHosts(t, ctx, count)

	broadcasters := CreateBroadcasters(t, ctx, hosts)

	testutil.ConnectAll(t, hosts)

	time.Sleep(time.Second)

	for _, broadcaster := range broadcasters {
		// Peer count does not include self
		assert.Equal(t, count-1, broadcaster.GetPeerCount())
	}

	cancel()
}

func CreateBroadcasters(
	t *testing.T,
	ctx context.Context,
	hosts []host.Host,
) []*broadcast.Broadcaster {
	var out []*broadcast.Broadcaster

	for _, h := range hosts {
		b, err := broadcast.NewBroadcaster(ctx, h)
		require.NoError(t, err)

		out = append(out, b)
	}

	return out
}
