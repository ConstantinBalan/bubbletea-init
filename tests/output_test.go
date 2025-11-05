package tests

import (
	"os"
	"path/filepath"
	"testing"

	initialize "github.com/ConstantinBalan/bubbletea-init/pkg/init"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomOutputDirectory(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	// Run the test
	os.Args = []string{"bubbletea-init", "-o", "custom/path", "testoutput"}
	initialize.Initialize()

	mainFile := filepath.Join(projectDir, "custom/path/testoutput", "main.go")
	assert.FileExists(t, mainFile)
}
