package allowlist_test

import (
	"context"
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/allowlist"
	"github.com/kava-labs/kava-bridge/relayer/testutil"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/suite"
)

type AllowListTestSuite struct {
	suite.Suite

	Ctx    context.Context
	Cancel context.CancelFunc

	Hosts []host.Host
}

func TestNodeTestSuite(t *testing.T) {
	suite.Run(t, new(AllowListTestSuite))
}

func (suite *AllowListTestSuite) SetupTest() {
	suite.Ctx, suite.Cancel = context.WithCancel(context.Background())
}

func (suite *AllowListTestSuite) TearDownTest() {
	suite.Cancel()

	for _, h := range suite.Hosts {
		h.Close()
	}
}

func (suite *AllowListTestSuite) TestPeerIDAllowList() {
	keys, err := testutil.GenerateNodeKeys(5)
	suite.Require().NoError(err)

	peerIDs := testutil.PeerIDsFromKeys(keys)

	suite.Hosts = testutil.CreateHostsWithKeys(
		suite.T(),
		suite.Ctx,
		keys,
		allowlist.PeerIDAllowList(peerIDs),
	)

	testutil.MustConnectAll(suite.T(), suite.Hosts)

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.True(allConnected)
}

func (suite *AllowListTestSuite) TestPeerIDAllowList_Missing() {
	keys, err := testutil.GenerateNodeKeys(5)
	suite.Require().NoError(err)

	peerIDs := testutil.PeerIDsFromKeys(keys)
	partialPeerIDs := peerIDs[:3]

	suite.Hosts = testutil.CreateHostsWithKeys(
		suite.T(),
		suite.Ctx,
		keys,
		allowlist.PeerIDAllowList(partialPeerIDs),
	)

	err = testutil.ConnectAll(suite.T(), suite.Hosts)
	suite.Require().Error(err)

	allConnected := testutil.AreAllConnected(suite.T(), suite.Hosts)
	suite.False(allConnected)
}
