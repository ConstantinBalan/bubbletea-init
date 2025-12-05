package tests

import (
	"os"
	"path/filepath"
	"testing"

	initialize "github.com/ConstantinBalan/bubbletea-init/pkg/init"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomModulePath(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	os.Args = []string{"bubbletea-init", "--mod", "github.com/testuser/testmod", "testmod"}
	initialize.Initialize()

	modFile := filepath.Join(projectDir, "testmod", "go.mod")
	content, err := os.ReadFile(modFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "module github.com/testuser/testmod")
}

func TestForceOverwrite(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	targetDir := filepath.Join(projectDir, "testforce")
	require.NoError(t, os.MkdirAll(targetDir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(targetDir, "dummy.txt"),
		[]byte("test"),
		0644,
	))

	os.Args = []string{"bubbletea-init", "--force", "testforce"}
	initialize.Initialize()

	mainFile := filepath.Join(targetDir, "main.go")
	assert.FileExists(t, mainFile)
}
