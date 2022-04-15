package testutil

import (
	"context"
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/broadcast"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Connect(t *testing.T, a, b host.Host) {
	pinfo := a.Peerstore().PeerInfo(a.ID())
	err := b.Connect(context.Background(), pinfo)
	if err != nil {
		t.Fatal(err)
	}
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

		t.Cleanup(func() { h.Close() })
		out = append(out, h)
	}

	return out
}

type Suite struct {
	suite.Suite

	Ctx    context.Context
	Cancel context.CancelFunc

	Hosts        []host.Host
	Broadcasters []*broadcast.Broadcaster
}

func (suite *Suite) CreateHostBroadcasters(n int, options ...broadcast.BroadcasterOption) {
	suite.Ctx, suite.Cancel = context.WithCancel(context.Background())

	count := 5
	suite.Hosts = CreateHosts(suite.T(), suite.Ctx, count)

	for i, h := range suite.Hosts {
		suite.T().Logf("peer index %v id: %v", i, h.ID())
	}

	for _, h := range suite.Hosts {
		b, err := broadcast.NewBroadcaster(suite.Ctx, h, options...)
		suite.Require().NoError(err)

		suite.Broadcasters = append(suite.Broadcasters, b)
	}
}

func (suite *Suite) ConnectAllHosts() {
	ConnectAll(suite.T(), suite.Hosts)
}

func (suite *Suite) SetupTest() {
}

func (suite *Suite) TearDownTest() {
	suite.Cancel()

	for _, h := range suite.Hosts {
		h.Close()
	}
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
