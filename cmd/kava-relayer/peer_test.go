package main_test

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
)

type TestPeer struct {
	port       int
	keyPath    string
	peerNumber int

	cmd *exec.Cmd
}

func NewTestPeer(peerNumber int, port int) *TestPeer {
	return &TestPeer{
		port:       port,
		peerNumber: peerNumber,
		keyPath:    fmt.Sprintf("test-fixtures/pk%d.key", peerNumber),
		cmd:        nil,
	}
}

func (p *TestPeer) Start(targets []string) error {
	p.cmd = startPeer(p.port, p.keyPath, strings.Join(targets, ","))

	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr

	if err := p.cmd.Start(); err != nil {
		return err
	}

	return nil
}

func (p *TestPeer) Wait() error {
	return p.cmd.Wait()
}

func (p *TestPeer) GetPeerID() (peer.ID, error) {
	cmd := execRelayer(
		"network",
		"show-node-id",
		"--p2p.private-key-path", fmt.Sprintf("test-fixtures/pk%d.key", p.peerNumber),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return peer.Decode(string(out))
}

func (p *TestPeer) GetMultiAddr() string {
	id, err := p.GetPeerID()
	if err != nil {
		panic(fmt.Sprintf("failed to get peer ID: %s", err))
	}

	return fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/p2p/%s", p.port, id)
}

func startPeer(port int, keyPath string, targets string) *exec.Cmd {
	return execRelayer(
		"network",
		"connect",
		"--p2p.peer-multiaddrs", targets,
		"--p2p.port", fmt.Sprintf("%d", port),
		"--p2p.private-key-path", keyPath,
		"--p2p.shared-key-path", "test-fixtures/psk.key",
		"--log_level", "debug",
	)
}

// GetFreePort asks the kernel for free open ports that are ready to use.
func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}
