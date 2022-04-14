package broadcast_test

import (
	"context"
	"testing"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/kava-labs/kava-bridge/relayer/types"
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

type TestHandler struct {
	rawCount   int
	validCount int
}

func (h *TestHandler) HandleRawMessage(msg *broadcast.MessageWithPeerMetadata) {
	h.rawCount += 1
}

func (h *TestHandler) HandleValidatedMessage(msg *types.MessageData) {
	h.validCount += 1
}

func TestBroadcast_Responses(t *testing.T) {
	logging.SetAllLoggers(logging.LevelDebug)

	ctx, cancel := context.WithCancel(context.Background())

	count := 5
	hosts := testutil.CreateHosts(t, ctx, count)

	handler := &TestHandler{
		rawCount:   0,
		validCount: 0,
	}

	broadcasters := CreateBroadcasters(t, ctx, hosts, broadcast.WithHandler(handler))

	testutil.ConnectAll(t, hosts)

	time.Sleep(time.Second)

	err := broadcasters[0].SendProtoMessage(&types.HelloRequest{
		Message: "hello world",
	})
	require.NoError(t, err)

	time.Sleep(time.Second * 5)

	for _, broadcaster := range broadcasters {
		// Peer count does not include self
		assert.Equal(t, count-1, broadcaster.GetPeerCount())
	}

	require.Equal(t, count-1, handler.rawCount)

	cancel()
}

func CreateBroadcasters(
	t *testing.T,
	ctx context.Context,
	hosts []host.Host,
	options ...broadcast.BroadcasterOption,
) []*broadcast.Broadcaster {
	var out []*broadcast.Broadcaster

	for _, h := range hosts {
		b, err := broadcast.NewBroadcaster(ctx, h, options...)
		require.NoError(t, err)

		out = append(out, b)
	}

	return out
}
