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
	"github.com/libp2p/go-libp2p-core/peer"
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

	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, n)

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
	err := testutil.ConnectAll(suite.T(), suite.Hosts)
	suite.Require().NoError(err)

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

	hostCount := 5
	suite.CreateHostBroadcasters(hostCount, broadcast.WithHandler(handler))
	err = testutil.ConnectAll(suite.T(), suite.Hosts)
	suite.Require().NoError(err)

	time.Sleep(time.Second)

	for _, broadcaster := range suite.Broadcasters {
		// Peer count does not include self
		suite.Assert().Equal(hostCount-1, broadcaster.GetPeerCount())
	}

	// Send message to all peers. This includes broadcaster peer but is ok since
	// broadcaster ignores self node
	allPeerIDs := testutil.PeerIDsFromHosts(suite.Hosts)

	tests := []struct {
		name           string
		recipients     []peer.ID
		wantRawCount   int
		wantValidCount int
	}{
		{
			"all including broadcaster",
			allPeerIDs,
			20,
			5,
		},
		{
			"all excluding broadcaster",
			allPeerIDs[1:],
			20,
			5,
		},
		{
			"partial including broadcaster",
			allPeerIDs[:4],
			12,
			4,
		},
		{
			"partial excluding broadcaster",
			allPeerIDs[1:4],
			12,
			4,
		},
		{
			"single including broadcaster",
			allPeerIDs[:2],
			2,
			2,
		},
		{
			"single excluding broadcaster",
			allPeerIDs[1:2],
			2,
			2,
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Reset handler counts
			handler.mu.Lock()
			handler.rawCount = 0
			handler.validCount = 0
			handler.mu.Unlock()

			err = suite.Broadcasters[0].BroadcastMessage(
				context.Background(),
				&types.HelloRequest{
					Message: "hello world",
				},
				tc.recipients,
				1, // 1 second TTL
			)
			suite.Require().NoError(err)

			time.Sleep(time.Second * 3)

			handler.mu.Lock()
			defer handler.mu.Unlock()

			// A -> B, C, D, E (4) // initial receive
			// B, C, D, E rebroadcast to all other nodes (4 * 4)
			// 4 initial + 16 re-broadcast = 20

			// n * (n - 1) messages where n is number of recipients
			suite.Assert().Equal(tc.wantRawCount, handler.rawCount, "raw message count should be 20")
			suite.Assert().Equal(tc.wantValidCount, handler.validCount, "each recipient peer should get a validated message")
		})
	}
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

func (h *TestHandler) ValidatedMessage(msg types.BroadcastMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.validCount += 1
}

func (h *TestHandler) MismatchMessage(msg broadcast.MessageWithPeerMetadata) {}
