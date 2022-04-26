package p2p_test

import (
	"context"
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/suite"
)

type NodeTestSuite struct {
	suite.Suite

	Ctx    context.Context
	Cancel context.CancelFunc

	Hosts []host.Host
}

func TestNodeTestSuite(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}

func (suite *NodeTestSuite) SetupTest() {
	suite.Ctx, suite.Cancel = context.WithCancel(context.Background())
}

func (suite *NodeTestSuite) TearDownTest() {
	suite.Cancel()

	for _, h := range suite.Hosts {
		h.Close()
	}
}

func (suite *NodeTestSuite) TestConnected() {
	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, 5)

	testutil.ConnectAll(suite.T(), suite.Hosts)

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.True(allConnected)
}

func (suite *NodeTestSuite) TestNotConnected() {
	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, 5)

	// testutil.ConnectAll(suite.T(), suite.Hosts)

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.False(allConnected)
}

func (suite *NodeTestSuite) TestNotConnected_Partial() {
	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, 5)

	testutil.ConnectAll(suite.T(), suite.Hosts)
	suite.Hosts[0].Close()

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.False(allConnected)
}
