package main_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrecomputeParams(t *testing.T) {
	t.Skip("skipping intensive precompute params test")

	// Unique per call, not necessary to delete file beforehand
	testHomeDir := t.TempDir()

	cmd := execRelayer("key", "precompute-preparams", "--home", testHomeDir)
	out, err := cmd.CombinedOutput()

	t.Log(string(out))
	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))
}

func TestPrecomputeParams_NoOverwrite(t *testing.T) {
	t.Skip("skipping intensive precompute params test")

	testHomeDir := t.TempDir()

	cmd := execRelayer("key", "precompute-preparams", "--home", testHomeDir)
	out, err := cmd.CombinedOutput()

	t.Log(string(out))
	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))

	// Run second time and check that it doesn't overwrite

	cmd2 := execRelayer("key", "precompute-preparams", "--home", testHomeDir)
	out2, err := cmd2.CombinedOutput()

	assert.Contains(t, string(out2), "Error: pre-params file already exists:")
	assert.EqualError(t, err, "exit status 1")
}

func TestPrecomputeParams_Force(t *testing.T) {
	t.Skip("skipping intensive precompute params test")

	testHomeDir := t.TempDir()

	cmd := execRelayer("key", "precompute-preparams", "--home", testHomeDir)
	out, err := cmd.CombinedOutput()

	t.Log(string(out))
	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))

	// Run second time and check that it **does** overwrite

	cmd2 := execRelayer("key", "precompute-preparams", "--home", testHomeDir, "--force")
	_, err = cmd2.CombinedOutput()

	require.NoError(t, err, fmt.Sprintf("expected '%s' to return successful status code", cmd.String()))
}
