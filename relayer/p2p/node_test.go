package p2p_test

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/p2p"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	swarm "github.com/libp2p/go-libp2p-swarm"
	"github.com/libp2p/go-tcp-transport"
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

	testutil.MustConnectAll(suite.T(), suite.Hosts)

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.True(allConnected)
}

func (suite *NodeTestSuite) TestNotConnected() {
	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, 5)

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.False(allConnected, "should be false if no hosts are connected")
}

func (suite *NodeTestSuite) TestNotConnected_Partial() {
	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, 5)

	testutil.MustConnectAll(suite.T(), suite.Hosts)
	suite.Hosts[0].Close()

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.False(allConnected, "should not be true with one host disconnected")
}

func (suite *NodeTestSuite) TestPrivateNetwork_MissingSecret() {
	networkSecret := make([]byte, p2p.PreSharedNetworkKeyLengthBytes)
	_, err := rand.Read(networkSecret)
	suite.Require().NoError(err)

	suite.Hosts = testutil.CreateHosts(
		suite.T(),
		suite.Ctx,
		4,
		libp2p.DefaultListenAddrs,
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.PrivateNetwork(networkSecret),
	)

	testutil.MustConnectAll(suite.T(), suite.Hosts)

	// Node without private network set
	adversary := testutil.CreateHost(
		suite.T(),
		libp2p.DefaultListenAddrs,
		libp2p.Transport(tcp.NewTCPTransport),
		// No private network option
	)

	// Adversary connect to host
	err = testutil.Connect(suite.T(), suite.Hosts[0], adversary)
	suite.Error(err)
	suite.IsType(&swarm.DialError{}, err)
}

func (suite *NodeTestSuite) TestPrivateNetwork_IncorrectSecret() {
	networkSecret := make([]byte, p2p.PreSharedNetworkKeyLengthBytes)
	_, err := rand.Read(networkSecret)
	suite.Require().NoError(err)

	suite.Hosts = testutil.CreateHosts(
		suite.T(),
		suite.Ctx,
		4,
		libp2p.DefaultListenAddrs,
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.PrivateNetwork(networkSecret),
	)

	testutil.MustConnectAll(suite.T(), suite.Hosts)

	incorrectNetworkSecret := make([]byte, p2p.PreSharedNetworkKeyLengthBytes)
	_, err = rand.Read(incorrectNetworkSecret)
	suite.Require().NoError(err)

	// Node without private network set
	adversary := testutil.CreateHost(
		suite.T(),
		libp2p.DefaultListenAddrs,
		libp2p.Transport(tcp.NewTCPTransport),
		// Incorrect
		libp2p.PrivateNetwork(incorrectNetworkSecret),
	)

	// Adversary connect to host
	err = testutil.Connect(suite.T(), suite.Hosts[0], adversary)
	suite.Error(err)
	suite.IsType(&swarm.DialError{}, err)
}

// func (suite *NodeTestSuite) TestReconnect() {
// // TODO: Figure out how to restart a node without changing peer ID
// 	suite.Hosts = testutil.CreateHosts(suite.T(), suite.Ctx, 5)

// 	testutil.MustConnectAll(suite.T(), suite.Hosts)
// 	suite.T().Logf("conns %v", suite.Hosts[0].Network().ConnsToPeer(suite.Hosts[1].ID()))
// 	suite.Hosts[0].Network().ConnsToPeer(suite.Hosts[1].ID())[0].Close()

// 	suite.Require().Eventually(func() bool {
// 		return testutil.AreAllConnected(suite.T(), suite.Hosts)
// 	}, time.Second*5, 10*time.Second)
// }
