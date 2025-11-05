package tests

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// setupTest creates a temporary directory and returns cleanup function
func setupTest(t *testing.T) (string, func()) {
	testDir, err := os.MkdirTemp("", "bubbletea-init-test-*")
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(testDir)
	}

	return testDir, cleanup
}

// setupTestEnv prepares the test environment and returns cleanup function
func setupTestEnv(t *testing.T, projectDir string) func() {
	oldArgs := os.Args
	oldDir, _ := os.Getwd()
	oldCommandLine := pflag.CommandLine

	// Change to test directory
	require.NoError(t, os.Chdir(projectDir))

	// Reset flags
	pflag.CommandLine = pflag.NewFlagSet("bubbletea-init", pflag.ExitOnError)

	return func() {
		os.Args = oldArgs
		os.Chdir(oldDir)
		pflag.CommandLine = oldCommandLine
	}
}
