package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

var log = logging.Logger(ProtocolID)

const (
	ProtocolID  = "/kava-relayer/echo/1.0.0"
	ServiceName = "kava-relayer.echo"
)

type EchoService struct {
	Host     host.Host
	peers    map[peer.ID]bool
	minPeers int
	done     chan bool
}

func NewEchoService(h host.Host, done chan bool, minPeers int) *EchoService {
	es := &EchoService{
		Host:     h,
		peers:    make(map[peer.ID]bool),
		minPeers: minPeers,
		done:     done,
	}

	// Labled as (Thread-safe), but double check if we need to use a sync.Mutex
	// in EchoService for state modifications in EchoHandler
	h.SetStreamHandler(ProtocolID, es.EchoHandler)
	return es
}

func (es *EchoService) EchoHandler(s network.Stream) {
	log.Info("listener received new echo stream", "peerID", s.Conn().RemotePeer())
	es.peers[s.Conn().RemotePeer()] = true
	log.Info(es.peers)

	if err := doEcho(s); err != nil {
		log.Error(err)
		s.Reset()
	} else {
		s.Close()
	}

	if len(es.peers) >= es.minPeers {
		es.done <- true
	}
}

func (es *EchoService) Echo(ctx context.Context, peerID peer.ID, payload string) (string, error) {
	stream, err := es.Host.NewStream(network.WithUseTransient(ctx, "echo"), peerID, ProtocolID)
	if err != nil {
		return "", err
	}

	_, err = stream.Write([]byte(payload))
	if err != nil {
		return "", fmt.Errorf("failed to write to stream: %s", err)
	}

	out, err := ioutil.ReadAll(stream)
	if err != nil {
		return "", fmt.Errorf("failed to read from stream: %s", err)
	}

	if !bytes.Equal(out, []byte(payload)) {
		return "", fmt.Errorf("echoed data was %s, expected %s", string(out), payload)
	}

	return fmt.Sprintf("%q", out), nil
}

// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Debug("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}
