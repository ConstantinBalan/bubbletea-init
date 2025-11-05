package tests

import (
	"os"
	"path/filepath"
	"testing"

	initialize "github.com/ConstantinBalan/bubbletea-init/pkg/init"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicProject(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	// Run the test
	os.Args = []string{"bubbletea-init", "testproject"}
	initialize.Initialize()

	// Check if main.go exists
	mainFile := filepath.Join(projectDir, "testproject", "main.go")
	assert.FileExists(t, mainFile)

	// Check if go.mod exists
	modFile := filepath.Join(projectDir, "testproject", "go.mod")
	assert.FileExists(t, modFile)

	// Read main.go and verify it contains expected content
	content, err := os.ReadFile(mainFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "package main")
	assert.Contains(t, string(content), "github.com/charmbracelet/bubbletea")
}

func TestProjectWithBubbles(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	// Run the test
	os.Args = []string{"bubbletea-init", "--with-bubbles", "testbubbles"}
	initialize.Initialize()

	// Check if files exist
	mainFile := filepath.Join(projectDir, "testbubbles", "main.go")
	modFile := filepath.Join(projectDir, "testbubbles", "go.mod")
	assert.FileExists(t, mainFile)
	assert.FileExists(t, modFile)

	// Verify main.go contains bubble components
	content, err := os.ReadFile(mainFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "type spinner struct")
	assert.Contains(t, string(content), "type textInput struct")
	assert.Contains(t, string(content), "github.com/charmbracelet/lipgloss")
}
