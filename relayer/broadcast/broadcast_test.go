package broadcast_test

import (
	"context"
	"sync"
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
	mu sync.Mutex

	rawCount   int
	validCount int
}

func (h *TestHandler) HandleRawMessage(msg *broadcast.MessageWithPeerMetadata) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.rawCount += 1
}

func (h *TestHandler) HandleValidatedMessage(msg *types.MessageData) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.validCount += 1
}

func TestBroadcast_Responses(t *testing.T) {
	// This is really noisy but useful for... debugging
	logging.SetAllLoggers(logging.LevelDebug)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	count := 5
	hosts := testutil.CreateHosts(t, ctx, count)

	for i, h := range hosts {
		t.Logf("peer index %v id: %v", i, h.ID())
	}

	handler := &TestHandler{
		rawCount:   0,
		validCount: 0,
	}

	broadcasters := CreateBroadcasters(t, ctx, hosts, broadcast.WithHandler(handler))

	testutil.ConnectAll(t, hosts)

	time.Sleep(time.Second)

	err := broadcasters[0].BroadcastMessage(
		context.Background(),
		"1234 message id",
		&types.HelloRequest{
			Message: "hello world",
		},
	)
	require.NoError(t, err)

	time.Sleep(time.Second * 10)

	for _, broadcaster := range broadcasters {
		// Peer count does not include self
		assert.Equal(t, count-1, broadcaster.GetPeerCount())
	}

	handler.mu.Lock()
	defer handler.mu.Unlock()

	// A -> B, C, D, E (4) // initial receive
	// B, C, D, E rebroadcast to all other nodes (4 * 4)
	// 4 initial + 16 re-broadcast = 20
	assert.Equal(t, 20, handler.rawCount, "raw message count should be 20")
	assert.Equal(t, count, handler.validCount, "each peer should get a validated message")
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
