package main

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
	if err := build.Run(); err != nil {
		fmt.Errorf("Failed to build %s: %s", binName, err)
		os.Exit(1)
	}

	r := m.Run()

	os.Remove(binName)
	os.Exit(r)
}

func TestNoArgs(t *testing.T) {
	cmd := exec.Command(fmt.Sprintf("./%s", binName))
	out, err := cmd.CombinedOutput()

	usageMsg := fmt.Sprintf("usage: %s \n", cmd.String())

	assert.Equal(t, usageMsg, string(out))
	assert.EqualError(t, err, "exit status 1")
}

func TestUnkownCommand(t *testing.T) {
	cmd := exec.Command(fmt.Sprintf("./%s", binName), "some-command")
	out, err := cmd.CombinedOutput()
	assert.Equal(t, "unknown command: some-command\n", string(out))
	assert.EqualError(t, err, "exit status 1")
}
