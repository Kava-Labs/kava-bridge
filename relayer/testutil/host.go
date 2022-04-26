package testutil

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// host.Host util functions

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

// AreAllConnected returns true if all hosts are connected to one another
func AreAllConnected(t *testing.T, hosts []host.Host) bool {
	for i, a := range hosts {
		// Only check hosts after i
		for _, b := range hosts[i+1:] {
			if !IsConnected(t, a, b) {
				t.Logf("%s and %s are not connected", a.ID(), b.ID())
				return false
			}
		}
	}

	return true
}

// IsConnected returns true if a and b hosts have a live, open connection in
// both directions
func IsConnected(t *testing.T, a, b host.Host) bool {
	return a.Network().Connectedness(b.ID()) == network.Connected &&
		b.Network().Connectedness(a.ID()) == network.Connected
}
