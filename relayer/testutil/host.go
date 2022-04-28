package testutil

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// host.Host util functions

func Connect(t *testing.T, a, b host.Host) error {
	pinfo := a.Peerstore().PeerInfo(a.ID())

	return b.Connect(context.Background(), pinfo)
}

func MustConnectAll(t *testing.T, hosts []host.Host) {
	err := ConnectAll(t, hosts)
	require.NoError(t, err)
}

func ConnectAll(t *testing.T, hosts []host.Host) error {
	for i, a := range hosts {
		for j, b := range hosts {
			if i == j {
				continue
			}

			if err := Connect(t, a, b); err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateHosts(
	t *testing.T,
	ctx context.Context,
	n int,
	options ...libp2p.Option,
) []host.Host {
	var out []host.Host

	for i := 0; i < n; i++ {
		h := CreateHost(t, options...)

		out = append(out, h)
	}

	return out
}

func CreateHostsWithKeys(
	t *testing.T,
	ctx context.Context,
	keys []crypto.PrivKey,
	options ...libp2p.Option,
) []host.Host {
	var out []host.Host

	for _, k := range keys {
		opts := append([]libp2p.Option{libp2p.Identity(k)}, options...)

		h := CreateHost(t, opts...)

		out = append(out, h)
	}

	return out
}

func CreateHost(t *testing.T, options ...libp2p.Option) host.Host {
	h, err := libp2p.New(options...)
	require.NoError(t, err)

	return h
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
