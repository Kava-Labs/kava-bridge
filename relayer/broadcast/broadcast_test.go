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
	"github.com/stretchr/testify/suite"
)

type BroadcasterTestSuite struct {
	testutil.Suite
}

func TestBroadcasterTestSuite(t *testing.T) {
	suite.Run(t, new(BroadcasterTestSuite))
}

func (suite *BroadcasterTestSuite) TestBroadcast_ConnectPeers() {
	count := 5
	suite.CreateHostBroadcasters(count)
	testutil.ConnectAll(suite.T(), suite.Hosts)

	time.Sleep(time.Second)

	for _, broadcaster := range suite.Broadcasters {
		// Peer count does not include self
		suite.Assert().Equal(count-1, broadcaster.GetPeerCount())
	}
}

func (suite *BroadcasterTestSuite) TestBroadcast_Responses() {
	logging.SetLogLevelRegex("broadcast", "debug")

	handler := &TestHandler{
		rawCount:   0,
		validCount: 0,
	}

	count := 5
	suite.CreateHostBroadcasters(count, broadcast.WithHandler(handler))
	testutil.ConnectAll(suite.T(), suite.Hosts)

	time.Sleep(time.Second)

	err := suite.Broadcasters[0].BroadcastMessage(
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
