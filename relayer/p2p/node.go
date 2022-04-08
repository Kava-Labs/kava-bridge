package p2p

import (
	"bufio"
	"context"
	"fmt"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/libp2p/go-libp2p-core/protocol"
	noise "github.com/libp2p/go-libp2p-noise"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/libp2p/go-tcp-transport"
	ma "github.com/multiformats/go-multiaddr"
)

var log = logging.Logger("p2p")

const RelayerProtocolID protocol.ID = "/kava-relayer/base/1.0.0"

type Node struct {
	Host host.Host
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

	host.SetStreamHandler(RelayerProtocolID, handleStream)

	return &Node{
		Host: host,
	}, nil
}

func (n Node) GetMultiAddress() ([]ma.Multiaddr, error) {
	peerInfo := peer.AddrInfo{
		ID:    n.Host.ID(),
		Addrs: n.Host.Addrs(),
	}
	return peer.AddrInfoToP2pAddrs(&peerInfo)
}

func (n Node) Connect(ctx context.Context, addr ma.Multiaddr) error {
	peer, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return err
	}

	if err := n.Host.Connect(context.Background(), *peer); err != nil {
		return err
	}

	res := <-ping.Ping(ctx, n.Host, peer.ID)
	if res.Error != nil {
		return fmt.Errorf("failed to ping peer: %w", res.Error)
	}

	log.Infof("ping took: %s", res.RTT)

	return nil
}

func handleStream(stream network.Stream) {
	log.Debug("Received a new stream")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	_, err := rw.WriteString("hi")
	if err != nil {
		log.Error("error writing to stream: ", err)
		stream.Reset()
	}

	stream.Close()
}
