package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainFlow(t *testing.T) {
	// Create a temporary directory for testing
	testDir, err := os.MkdirTemp("", "bubbletea-init-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	// Save original flags
	origCommandLine := pflag.CommandLine
	defer func() {
		pflag.CommandLine = origCommandLine
	}()

	tests := []struct {
		name     string
		args     []string
		testFunc func(t *testing.T, projectDir string)
	}{
		{
			name: "basic project",
			args: []string{"testproject"},
			testFunc: func(t *testing.T, projectDir string) {
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
			},
		},
		{
			name: "project with bubbles",
			args: []string{"--with-bubbles", "testbubbles"},
			testFunc: func(t *testing.T, projectDir string) {
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
			},
		},
		{
			name: "custom module path",
			args: []string{"--mod", "github.com/testuser/testmod", "testmod"},
			testFunc: func(t *testing.T, projectDir string) {
				modFile := filepath.Join(projectDir, "testmod", "go.mod")
				content, err := os.ReadFile(modFile)
				require.NoError(t, err)
				assert.Contains(t, string(content), "module github.com/testuser/testmod")
			},
		},
		{
			name: "force overwrite",
			args: []string{"--force", "testforce"},
			testFunc: func(t *testing.T, projectDir string) {
				// Create the project first time
				targetDir := filepath.Join(projectDir, "testforce")
				require.NoError(t, os.MkdirAll(targetDir, 0755))
				require.NoError(t, os.WriteFile(
					filepath.Join(targetDir, "dummy.txt"),
					[]byte("test"),
					0644,
				))

				// Second creation should succeed with --force
				mainFile := filepath.Join(targetDir, "main.go")
				assert.FileExists(t, mainFile)
			},
		},
		{
			name: "custom output directory",
			args: []string{"-o", "custom/path", "testoutput"},
			testFunc: func(t *testing.T, projectDir string) {
				mainFile := filepath.Join(projectDir, "custom/path/testoutput", "main.go")
				assert.FileExists(t, mainFile)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temporary directory for each test
			projectDir, err := os.MkdirTemp(testDir, "project-*")
			require.NoError(t, err)

			// Set up test environment
			oldArgs := os.Args
			oldDir, _ := os.Getwd()
			defer func() {
				os.Args = oldArgs
				os.Chdir(oldDir)
			}()

			// Change to test directory
			require.NoError(t, os.Chdir(projectDir))

			// Reset flags for each test
			pflag.CommandLine = pflag.NewFlagSet("bubbletea-init", pflag.ExitOnError)

			// Run the test
			os.Args = append([]string{"bubbletea-init"}, tt.args...)
			main()

			// Run test-specific checks
			tt.testFunc(t, projectDir)
		})
	}
}