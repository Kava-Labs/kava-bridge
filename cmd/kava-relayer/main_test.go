package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
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
	cmd := exec.Command(fmt.Sprintf("./%s", binName))
	out, err := cmd.CombinedOutput()

	assert.Contains(t, string(out), "The kava relayer processes ethereum and kava blocks to transfer ERC20 tokens between chains.")
	assert.NoError(t, err)
}
