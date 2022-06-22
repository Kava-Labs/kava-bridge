package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/kava-labs/kava-bridge/cmd/kava-relayer/cmd/network"

	"github.com/libp2p/go-libp2p-core/crypto"
	crypto_pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
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

func TestShowNodeMultiAddress(t *testing.T) {
	cmd := execRelayer("network", "show-node-multi-address", "--p2p.private-key-path", "test-fixtures/pk1.key", "--p2p.port", "8000")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))

	actualPeerMultiAddress := string(out)
	_, err = ma.NewMultiaddr(actualPeerMultiAddress)
	require.NoError(t, err)

	testHostIPv4, err := network.GetHostIPv4()
	require.NoError(t, err)
	expectedPeerMultiAddress := fmt.Sprintf("/ip4/%s/tcp/8000/p2p/16Uiu2HAm9z3t15JpqBbPQJ1ZLHm6w1AXD6M2FXdCG3GLoY4iDcD9", testHostIPv4)
	assert.Equal(t, actualPeerMultiAddress, expectedPeerMultiAddress)
}

func TestConnectPeers(t *testing.T) {
	var peers []*TestPeer
	var peerAddrs []string

	ports, err := GetFreePorts(5)
	require.NoError(t, err)

	// Start at 1 for pk1.key
	for i := 1; i <= 5; i++ {
		peer := NewTestPeer(i, ports[i-1])
		peers = append(peers, peer)
		peerAddrs = append(peerAddrs, peer.GetMultiAddr())
	}

	for i, peer := range peers {
		// Remove current peer id from peerAddrs. Need to copy to not modify original slice.
		peerAddrsWithoutSelf := make([]string, len(peerAddrs))
		copy(peerAddrsWithoutSelf, peerAddrs)
		peerAddrsWithoutSelf = append(peerAddrsWithoutSelf[:i], peerAddrsWithoutSelf[i+1:]...)

		err := peer.Start(peerAddrsWithoutSelf)
		require.NoErrorf(t, err, "expected peer %d (%s) to start successfully", i, peer.cmd.String())
	}

	for i, peer := range peers {
		err := peer.Wait()
		require.NoErrorf(t, err, "expected peer %d (%s) to return with exit code 0", i, peer.cmd.String())
	}
}

func execRelayer(args ...string) *exec.Cmd {
	return exec.Command(fmt.Sprintf("./%s", binName), args...)
}
