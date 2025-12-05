package tests

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (string, func()) {
	testDir, err := os.MkdirTemp("", "bubbletea-init-test-*")
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(testDir)
	}

	return testDir, cleanup
}

func setupTestEnv(t *testing.T, projectDir string) func() {
	oldArgs := os.Args
	oldDir, _ := os.Getwd()
	oldCommandLine := pflag.CommandLine

	require.NoError(t, os.Chdir(projectDir))

	pflag.CommandLine = pflag.NewFlagSet("bubbletea-init", pflag.ExitOnError)

	return func() {
		os.Args = oldArgs
		os.Chdir(oldDir)
		pflag.CommandLine = oldCommandLine
	}
}
