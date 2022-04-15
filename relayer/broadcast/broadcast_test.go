package broadcast_test

import (
	"context"
	"sync"
	"testing"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/types"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BroadcasterTestSuite struct {
	suite.Suite

	Ctx    context.Context
	Cancel context.CancelFunc

	Hosts        []host.Host
	Broadcasters []*broadcast.Broadcaster
}

func TestBroadcasterTestSuite(t *testing.T) {
	suite.Run(t, new(BroadcasterTestSuite))
}

func (suite *BroadcasterTestSuite) TearDownTest() {
	suite.Cancel()

	for _, h := range suite.Hosts {
		h.Close()
	}
}

func (suite *BroadcasterTestSuite) CreateHostBroadcasters(n int, options ...broadcast.BroadcasterOption) {
	suite.Ctx, suite.Cancel = context.WithCancel(context.Background())

	suite.Hosts = CreateHosts(suite.T(), suite.Ctx, n)

	// Without setting to nil first, suite tests will connect to peers
	// on a different suite test for some reason.
	suite.Broadcasters = nil

	for i, h := range suite.Hosts {
		suite.T().Logf("peer index %v id: %v", i, h.ID())

		b, err := broadcast.NewBroadcaster(suite.Ctx, h, options...)
		suite.Require().NoError(err)

		suite.Broadcasters = append(suite.Broadcasters, b)
	}
}

func (suite *BroadcasterTestSuite) TestBroadcast_ConnectPeers() {
	count := 5
	suite.CreateHostBroadcasters(count)
	ConnectAll(suite.T(), suite.Hosts)

	time.Sleep(time.Second)

	for _, broadcaster := range suite.Broadcasters {
		// Peer count does not include self
		suite.Assert().Equal(count-1, broadcaster.GetPeerCount())
	}
}

func (suite *BroadcasterTestSuite) TestBroadcast_Responses() {
	err := logging.SetLogLevelRegex("broadcast", "debug")
	suite.Require().NoError(err)

	handler := &TestHandler{
		rawCount:   0,
		validCount: 0,
	}

	count := 5
	suite.CreateHostBroadcasters(count, broadcast.WithHandler(handler))
	ConnectAll(suite.T(), suite.Hosts)

	time.Sleep(time.Second)

	err = suite.Broadcasters[0].BroadcastMessage(
		context.Background(),
		"1234 message id",
		&types.HelloRequest{
			Message: "hello world",
		},
	)
	suite.Require().NoError(err)

	time.Sleep(time.Second * 10)

	for _, broadcaster := range suite.Broadcasters {
		// Peer count does not include self
		suite.Assert().Equal(count-1, broadcaster.GetPeerCount())
	}

	handler.mu.Lock()
	defer handler.mu.Unlock()

	// A -> B, C, D, E (4) // initial receive
	// B, C, D, E rebroadcast to all other nodes (4 * 4)
	// 4 initial + 16 re-broadcast = 20
	suite.Assert().Equal(20, handler.rawCount, "raw message count should be 20")
	suite.Assert().Equal(count, handler.validCount, "each peer should get a validated message")
}

// ----------------------------------------------------------------------------
// test handler
type TestHandler struct {
	mu sync.Mutex

	rawCount   int
	validCount int
}

func (h *TestHandler) RawMessage(msg broadcast.MessageWithPeerMetadata) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.rawCount += 1
}

func (h *TestHandler) ValidatedMessage(msg types.MessageData) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.validCount += 1
}

func (h *TestHandler) MismatchMessage(msg broadcast.MessageWithPeerMetadata) {}

// ----------------------------------------------------------------------------
// util functions

func Connect(t *testing.T, a, b host.Host) {
	pinfo := a.Peerstore().PeerInfo(a.ID())
	err := b.Connect(context.Background(), pinfo)
	require.NoError(t, err)
}

func ConnectAll(t *testing.T, hosts []host.Host) {
	for i, a := range hosts {
		for j, b := range hosts {
			if i == j {
				continue
			}

			Connect(t, a, b)
		}
	}
}

func CreateHosts(t *testing.T, ctx context.Context, n int) []host.Host {
	var out []host.Host

	for i := 0; i < n; i++ {
		h, err := libp2p.New()
		require.NoError(t, err)

		out = append(out, h)
	}

	return out
}
