//go:build integration
// +build integration

package testing

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var (
	bridgeBinName  = "kava-bridged"
	relayerBinName = "kava-relayer"
)

func buildBin(binName string) (func(), error) {
	build := exec.Command("go", "build", "-o", binName, fmt.Sprintf("../cmd/%s", binName))
	if err := build.Run(); err != nil {
		return nil, fmt.Errorf("Failed to build %s: %s", binName, err)
	}

	return func() { os.Remove(binName) }, nil
}

func TestMain(m *testing.M) {
	// build the kava test chain w/
	cleanupBridge, err := buildBin(bridgeBinName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// build the kava relayer
	cleanupRelayer, err := buildBin(relayerBinName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// run test suite
	r := m.Run()

	// clean up binaries
	cleanupBridge()
	cleanupRelayer()

	// exit
	os.Exit(r)
}
