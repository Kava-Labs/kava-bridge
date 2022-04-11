package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
)

const START_PORT = 9000

type TestPeer struct {
	port       uint16
	keyPath    string
	peerNumber int

	cmd *exec.Cmd
}

func NewTestPeer(peerNumber int) *TestPeer {
	return &TestPeer{
		port:       START_PORT + uint16(peerNumber),
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

func startPeer(port uint16, keyPath string, targets string) *exec.Cmd {
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
