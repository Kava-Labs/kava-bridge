package broadcast_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/pending_store"
	"github.com/kava-labs/kava-bridge/relayer/broadcast/types"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/suite"
)

type TestBroadcaster struct {
	*broadcast.Broadcaster
	handler *TestHandler
}

func NewTestBroadcaster(
	ctx context.Context,
	h host.Host,
	opts ...broadcast.BroadcasterOption,
) (*TestBroadcaster, error) {
	handler := &TestHandler{
		mu:         &sync.Mutex{},
		rawCount:   0,
		validCount: 0,
	}

	opts = append(opts, broadcast.WithHandler(handler))

	b, err := broadcast.NewBroadcaster(ctx, h, opts...)
	if err != nil {
		return nil, err
	}

	return &TestBroadcaster{
		Broadcaster: b,
		handler:     handler,
	}, nil
}

func (tb *TestBroadcaster) GetValidCount() int {
	tb.handler.mu.Lock()
	defer tb.handler.mu.Unlock()

	return tb.handler.validCount
}

func (tb *TestBroadcaster) GetRawCount() int {
	tb.handler.mu.Lock()
	defer tb.handler.mu.Unlock()

	return tb.handler.rawCount
}

func (tb *TestBroadcaster) ResetCounts() {
	tb.handler.mu.Lock()
	defer tb.handler.mu.Unlock()

	tb.handler.rawCount = 0
	tb.handler.validCount = 0
}

type BroadcasterTestSuite struct {
	suite.Suite

	Ctx    context.Context
	Cancel context.CancelFunc

	Hosts        []host.Host
	Broadcasters []*TestBroadcaster
}

func TestBroadcasterTestSuite(t *testing.T) {
	suite.Run(t, new(BroadcasterTestSuite))
}

func (suite *BroadcasterTestSuite) SetupTest() {
	suite.Ctx, suite.Cancel = context.WithCancel(context.Background())
}

func (suite *BroadcasterTestSuite) TearDownTest() {
	suite.Cancel()

	for _, h := range suite.Hosts {
		h.Close()
	}
}

func (suite *BroadcasterTestSuite) CreateHostBroadcasters(n int, options ...broadcast.BroadcasterOption) {
	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, n)

	suite.CreateBroadcasterWithHosts(suite.Hosts, options...)
}

func (suite *BroadcasterTestSuite) CreateBroadcasterWithHosts(
	hosts []host.Host,
	options ...broadcast.BroadcasterOption,
) {
	// Without setting to nil first, suite tests will connect to peers
	// on a different suite test for some reason.
	suite.Broadcasters = nil

	for i, h := range hosts {
		suite.T().Logf("peer index %v id: %v", i, h.ID())

		b, err := NewTestBroadcaster(suite.Ctx, h, options...)
		suite.Require().NoError(err)

		suite.Broadcasters = append(suite.Broadcasters, b)
	}
}

func (suite *BroadcasterTestSuite) WaitForAllBroadcastersConnected() {
	timeout := time.Second * 10
	start := time.Now()

	for {
		time.Sleep(time.Millisecond * 300)

		completeCount := 0
		for _, broadcaster := range suite.Broadcasters {
			// Peer count does not include self
			if broadcaster.GetPeerCount() == len(suite.Hosts)-1 {
				completeCount += 1
			}
		}

		if completeCount == len(suite.Broadcasters) {
			return
		}

		if start.Add(timeout).Before(time.Now()) {
			suite.T().Fatal("timeout waiting for hosts to connect")
		}
	}
}

func (suite *BroadcasterTestSuite) RequireHandlersRawCounts(expectedRawCounts []int) {
	if len(expectedRawCounts) != len(suite.Broadcasters) {
		suite.Fail("expectedRawCounts and Broadcasters are not the same length")
	}

	for i, b := range suite.Broadcasters {
		suite.Equalf(expectedRawCounts[i], b.GetRawCount(), "invalid raw message count for peer %v", i)
	}
}

func (suite *BroadcasterTestSuite) RequireHandlersValidCounts(expectedValidCounts []int) {
	if len(expectedValidCounts) != len(suite.Broadcasters) {
		suite.Fail("expectedValidCounts and Broadcasters are not the same length")
	}

	for i, b := range suite.Broadcasters {
		suite.Equalf(expectedValidCounts[i], b.GetValidCount(), "invalid valid message count for peer %v", i)
	}
}

// ResetBroadcasterCounts resets the raw and valid message counts for all
// broadcasters, should be called between test cases
func (suite *BroadcasterTestSuite) ResetBroadcasterCounts() {
	for _, b := range suite.Broadcasters {
		b.ResetCounts()
	}
}

func (suite *BroadcasterTestSuite) TestBroadcast_ConnectPeers() {
	count := 5
	suite.CreateHostBroadcasters(count)
	err := testutil.ConnectAll(suite.T(), suite.Hosts)
	suite.Require().NoError(err)
	suite.WaitForAllBroadcastersConnected()

	for _, broadcaster := range suite.Broadcasters {
		// Peer count does not include self
		suite.Assert().Equal(count-1, broadcaster.GetPeerCount())
	}
}

func (suite *BroadcasterTestSuite) TestBroadcast_Responses() {
	err := logging.SetLogLevelRegex("broadcast", "debug")
	suite.Require().NoError(err)

	hostCount := 5
	suite.CreateHostBroadcasters(hostCount)

	err = testutil.ConnectAll(suite.T(), suite.Hosts)
	suite.Require().NoError(err)

	suite.WaitForAllBroadcastersConnected()

	// Send message to all peers. This includes broadcaster peer but is ok since
	// broadcaster ignores self node
	allPeerIDs := testutil.PeerIDsFromHosts(suite.Hosts)

	tests := []struct {
		name            string
		recipients      []peer.ID
		wantRawCounts   []int
		wantValidCounts []int
	}{
		{
			"all including broadcaster",
			allPeerIDs,
			[]int{4, 4, 4, 4, 4},
			[]int{1, 1, 1, 1, 1},
		},
		{
			"all excluding broadcaster",
			allPeerIDs[1:],
			[]int{4, 4, 4, 4, 4},
			[]int{1, 1, 1, 1, 1},
		},
		{
			"partial including broadcaster",
			allPeerIDs[:4],
			[]int{3, 3, 3, 3, 0},
			[]int{1, 1, 1, 1, 0},
		},
		{
			"partial excluding broadcaster",
			allPeerIDs[1:4],
			[]int{3, 3, 3, 3, 0},
			[]int{1, 1, 1, 1, 0},
		},
		{
			"single including broadcaster",
			allPeerIDs[:2],
			[]int{1, 1, 0, 0, 0},
			[]int{1, 1, 0, 0, 0},
		},
		{
			"single excluding broadcaster",
			allPeerIDs[1:2],
			[]int{1, 1, 0, 0, 0},
			[]int{1, 1, 0, 0, 0},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.ResetBroadcasterCounts()

			err = suite.Broadcasters[0].BroadcastMessage(
				context.Background(),
				&types.HelloRequest{
					PeerID:      suite.Hosts[0].ID(),
					NodeMoniker: "hello world",
				},
				tc.recipients,
				8, // TTL
			)
			suite.Require().NoError(err)

			time.Sleep(time.Second * 3)

			suite.RequireHandlersRawCounts(tc.wantRawCounts)
			suite.RequireHandlersValidCounts(tc.wantValidCounts)
		})
	}
}

func (suite *BroadcasterTestSuite) TestBroadcast_TTL() {
	err := logging.SetLogLevelRegex("broadcast", "debug")
	suite.Require().NoError(err)

	keys, err := testutil.GenerateNodeKeys(5)
	suite.Require().NoError(err)

	peerIDs := testutil.PeerIDsFromKeys(keys)

	suite.Hosts = testutil.CreateHostsWithKeys(
		suite.T(),
		suite.Ctx,
		keys,
	)

	sleepTimeFn := func(
		b *broadcast.Broadcaster,
		target peer.ID,
		pb *proto.Message,
	) time.Duration {
		// Delay sending to last peer
		if target == peerIDs[4] {
			suite.T().Logf("Delaying message to %v", target)
			return time.Second * 7
		}

		return 0
	}

	suite.CreateBroadcasterWithHosts(suite.Hosts, broadcast.WithHook(&SleepyBroadcasterHook{
		sleepTimeFn,
	}))

	err = testutil.ConnectAll(suite.T(), suite.Hosts)
	suite.Require().NoError(err)
	suite.WaitForAllBroadcastersConnected()

	// Send message to all peers. This includes broadcaster peer but is ok since
	// broadcaster ignores self node
	allPeerIDs := testutil.PeerIDsFromHosts(suite.Hosts)

	err = suite.Broadcasters[0].BroadcastMessage(
		context.Background(),
		&types.HelloRequest{
			PeerID:      suite.Hosts[0].ID(),
			NodeMoniker: "hello world",
		},
		allPeerIDs,
		// Takes a few seconds for other peers to receive the message
		// Should be high enough for round trip
		7,
	)
	suite.Require().NoError(err)

	time.Sleep(time.Second * 12)

	// All peers delay send to peer 4, peer 4 will consider them expired and won't re-broadcast any
	suite.RequireHandlersRawCounts([]int{3, 3, 3, 3, 0})
	suite.RequireHandlersValidCounts([]int{0, 0, 0, 0, 0})
}

// ----------------------------------------------------------------------------
// test handler
type TestHandler struct {
	// single mutex for everything just for simplicity.
	mu *sync.Mutex

	rawCount   int
	validCount int
}

func (h *TestHandler) RawMessage(msg pending_store.MessageWithPeerMetadata) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.rawCount += 1
}

func (h *TestHandler) ValidatedMessage(msg types.BroadcastMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.validCount += 1
}

func (h *TestHandler) MismatchMessage(msg pending_store.MessageWithPeerMetadata) {}

// ----------------------------------------------------------------------------
// delay broadcast hook

// SleepyBroadcasterHook is a broadcasterHook that conditionally delays broadcasting raw messages
type SleepyBroadcasterHook struct {
	sleepTimeFn func(
		b *broadcast.Broadcaster,
		target peer.ID,
		pb *proto.Message,
	) time.Duration
}

var _ broadcast.BroadcasterHook = (*SleepyBroadcasterHook)(nil)

func (h *SleepyBroadcasterHook) BeforeBroadcastRawMessage(
	b *broadcast.Broadcaster,
	target peer.ID,
	pb *proto.Message,
) {
	time.Sleep(h.sleepTimeFn(b, target, pb))
}
