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

func NewNode(options NodeOptions) (*Node, error) {
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

	// Need to be buffered by 1 to not block
	done := make(chan bool, 1)

	node := &Node{
		Host: host,
		// Sets stream handler
		EchoService: service.NewEchoService(host, done, 1),
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

func (n Node) Connect(ctx context.Context, addr ma.Multiaddr) error {
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return err
	}

	// TODO: Determine TTL for peer
	n.Host.Peerstore().AddAddrs(peerAddrInfo.ID, peerAddrInfo.Addrs, peerstore.RecentlyConnectedAddrTTL)

	// Retry connection 10 times to account for peers starting later than others
	err = retry(10, time.Second, func() error {
		return n.Host.Connect(ctx, *peerAddrInfo)
	})
	if err != nil {
		return err
	}

	log.Info("success")

	res, err := n.EchoService.Echo(ctx, peerAddrInfo.ID, "hello world!\n")
	if err != nil {
		return err
	}

	log.Info("received echo response: ", res)

	log.Info("waiting for all echo requests")

	select {
	case <-n.done:
	case <-ctx.Done():
		return ctx.Err()
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
