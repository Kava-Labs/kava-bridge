package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	crypto_pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	binName = "kava-relayer"
)

func TestMain(m *testing.M) {
	build := exec.Command("go", "build", "-o", binName)
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build %s: %s\n\n", binName, err)
		fmt.Fprintf(os.Stderr, "%s\n\n", string(out))
		os.Exit(1)
	}

	r := m.Run()

	os.Remove(binName)
	os.Exit(r)
}

func TestNoArgs(t *testing.T) {
	cmd := execRelayer()
	out, err := cmd.CombinedOutput()

	assert.Contains(t, string(out), "The kava relayer processes ethereum and kava blocks to transfer ERC20 tokens between chains.")
	assert.NoError(t, err)
}

func TestUnknownCommand(t *testing.T) {
	cmd := execRelayer("some-command")
	out, err := cmd.CombinedOutput()
	assert.Contains(t, string(out), "Error: unknown command \"some-command\" for \"kava-relayer")
	assert.EqualError(t, err, "exit status 1")
}

func TestPrivateNetworkSecretGeneration(t *testing.T) {
	cmd := execRelayer("network", "generate-network-secret")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))

	encoding, secret, err := multibase.Decode(string(out))
	require.NoError(t, err, "expected secret to successfully decode")

	assert.Equal(t, multibase.Encoding(multibase.Base58BTC), encoding, "expected secret to be base 58 (btc) encoded")
	assert.Equal(t, 32, len(secret), "expected secret to be 256-bits / 32 bytes")
}

func TestNodePrivateKeyGeneration(t *testing.T) {
	cmd := execRelayer("network", "generate-node-key")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))

	privKey, err := crypto.UnmarshalPrivateKey(out)
	require.NoError(t, err, "expected private key to successfully decode")

	assert.Equal(t, crypto_pb.KeyType_Secp256k1, privKey.Type(), "expected private key to be a secp256k1 key")

	rawKey, err := privKey.Raw()
	require.NoError(t, err, "expected private key to successful decode to raw bytes")
	assert.Equal(t, 32, len(rawKey), "expected private key to be 256-bits / 32 bytes")
}

func TestShowNodeID(t *testing.T) {
	cmd := execRelayer("network", "show-node-id", "--p2p.private-key-path", "test-fixtures/pk1.key")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))

	peerID, err := peer.Decode(string(out))
	require.NoError(t, err)

	pkData, err := os.ReadFile("test-fixtures/pk1.key")
	require.NoError(t, err)
	privKey, err := crypto.UnmarshalPrivateKey(pkData)
	require.NoError(t, err)

	assert.True(t, peerID.MatchesPrivateKey(privKey))
}

func TestConnectPeers(t *testing.T) {
	peer1 := startPeer(8765, "test-fixtures/pk1.key", "/ip4/127.0.0.1/tcp/8764/p2p/16Uiu2HAkwC5w1fC4xLL3hWjD6PGuk2qzGgsWdXfNeqMi8xDn2AT7")
	peer2 := startPeer(8764, "test-fixtures/pk2.key", "/ip4/127.0.0.1/tcp/8765/p2p/16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9")

	peer1.Stderr = os.Stderr
	peer2.Stderr = os.Stderr

	err := peer1.Start()
	require.NoErrorf(t, err, "expected peer 1 (%s) to start successfully", peer1.String())

	err = peer2.Start()
	require.NoErrorf(t, err, "expected peer 2 (%s) to start successfully", peer2.String())

	err = peer1.Wait()
	require.NoError(t, err)

	err = peer2.Wait()
	require.NoError(t, err)
}

func startPeer(port uint16, key string, target string) *exec.Cmd {
	return execRelayer(
		"network",
		"connect",
		target,
		"--p2p.port", fmt.Sprintf("%d", port),
		"--p2p.private-key-path", key,
		"--p2p.shared-key-path", "test-fixtures/psk.key",
	)
}

func execRelayer(args ...string) *exec.Cmd {
	return exec.Command(fmt.Sprintf("./%s", binName), args...)
}
