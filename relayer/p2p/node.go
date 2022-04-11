package p2p

import (
	"context"
	"fmt"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/kava-labs/kava-bridge/relayer/p2p/service"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/pnet"
	noise "github.com/libp2p/go-libp2p-noise"
	"github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("p2p")

type Node struct {
	Host        host.Host
	EchoService *service.EchoService
	done        chan bool
}

func NewNode(options NodeOptions, done chan bool) (*Node, error) {
	libp2pOpts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", options.Port)),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.PrivateNetwork(options.NetworkPrivateKey),
		libp2p.Identity(options.NodePrivateKey),
		libp2p.DisableRelay(),
		libp2p.Security(noise.ID, noise.New),
	}

	pnet.ForcePrivateNetwork = true

	host, err := libp2p.New(libp2pOpts...)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Host: host,
		// Sets stream handler
		EchoService: service.NewEchoService(host, done, options.EchoRequiredPeers),
		done:        done,
	}

	registerNotifiees(host)

	return node, nil
}

func (n Node) GetMultiAddress() ([]ma.Multiaddr, error) {
	peerInfo := peer.AddrInfo{
		ID:    n.Host.ID(),
		Addrs: n.Host.Addrs(),
	}

	return peer.AddrInfoToP2pAddrs(&peerInfo)
}

func (n Node) ConnectToPeers(ctx context.Context, peerAddrInfos []*peer.AddrInfo) error {
	for _, peer := range peerAddrInfos {
		// TODO: Determine TTL for peer
		n.Host.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.RecentlyConnectedAddrTTL)

		// Retry connection 10 times to account for peers starting later than others.
		// This is not entirely necessary as using a service will also connect
		// to the node, but connecting prior to making requests ensures that
		// this node can connect to the peer.
		err := retry(10, time.Second, func() error {
			return n.Host.Connect(ctx, *peer)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (n Node) EchoPeers(ctx context.Context) error {
	log.Debugf("sending echo to %d peers", len(n.Host.Peerstore().Peers())-1)

	for _, peerID := range n.Host.Peerstore().Peers() {
		// Skip self
		if n.Host.ID() == peerID {
			continue
		}

		res, err := n.EchoService.Echo(ctx, peerID, "hello world!\n")
		if err != nil {
			return err
		}

		log.Info("received echo response: ", res)
	}

	return nil
}

// Close cleans up and stops the node
func (n Node) Close() error {
	return n.Host.Close()
}

func registerNotifiees(host host.Host) {
	var notifee network.NotifyBundle
	notifee.ConnectedF = func(net network.Network, conn network.Conn) {
		log.Info("connected to peer: ", conn.RemotePeer())
	}

	notifee.DisconnectedF = func(net network.Network, conn network.Conn) {
		log.Info("disconnected from peer: ", conn.RemotePeer())
	}

	host.Network().Notify(&notifee)
}
