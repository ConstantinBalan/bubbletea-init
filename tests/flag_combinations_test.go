package tests

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	initialize "github.com/ConstantinBalan/bubbletea-init/pkg/init"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func TestHelpFlag(t *testing.T) {
	_, cleanup := setupTest(t)
	defer cleanup()

	resetFlags()

	tests := []struct {
		name     string
		args     []string
		exitCode int
	}{
		{
			name:     "help flag short",
			args:     []string{"bubbletea-init", "-h"},
			exitCode: 0,
		},
		{
			name:     "help flag long",
			args:     []string{"bubbletea-init", "--help"},
			exitCode: 0,
		},
		{
			name:     "help with project name",
			args:     []string{"bubbletea-init", "-h", "myproject"},
			exitCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()
			// Capture stdout and stderr to verify help message
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			// Set up exit handler by overriding initialize.Exit
			var exitCode int
			oldExit := initialize.Exit
			initialize.Exit = func(code int) {
				exitCode = code
				panic("exit") // Use panic to stop execution but in a way we can recover
			}
			defer func() {
				initialize.Exit = oldExit
				os.Stdout = oldStdout
				os.Stderr = oldStderr
				_ = recover() // Recover from our panic
			}()

			os.Args = tt.args

			// Run initialize and expect it to "exit"
			func() {
				defer func() { _ = recover() }()
				initialize.Initialize()
			}()

			// Close writer and get output
			w.Close()
			outBytes, _ := io.ReadAll(r)
			out := string(outBytes)

			assert.Equal(t, tt.exitCode, exitCode)
			assert.Contains(t, out, "Usage: bubbletea-init [flags] <project-name>")
			assert.Contains(t, out, "--with-bubbles")
			assert.Contains(t, out, "--mod")
			assert.Contains(t, out, "--output-dir")
			assert.Contains(t, out, "--force")
		})
	}
}

func TestMultipleFlagCombinations(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	// Test with bubbles and custom module
	t.Run("bubbles and custom module", func(t *testing.T) {
		resetFlags()
		os.Args = []string{
			"bubbletea-init",
			"--with-bubbles",
			"--mod", "github.com/test/combo",
			"testcombo",
		}
		initialize.Initialize()

		mainFile := filepath.Join(projectDir, "testcombo", "main.go")
		modFile := filepath.Join(projectDir, "testcombo", "go.mod")

		// Check main.go content
		content, err := os.ReadFile(mainFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), "type spinner struct")
		assert.Contains(t, string(content), "github.com/charmbracelet/lipgloss")

		// Check go.mod content
		modContent, err := os.ReadFile(modFile)
		require.NoError(t, err)
		assert.Contains(t, string(modContent), "module github.com/test/combo")
		assert.Contains(t, string(modContent), "github.com/charmbracelet/lipgloss")
	})

	// Test with output directory and force
	t.Run("output dir and force", func(t *testing.T) {
		resetFlags()
		customDir := filepath.Join(projectDir, "custom")
		require.NoError(t, os.MkdirAll(customDir, 0755))

		targetDir := filepath.Join(customDir, "testoutputforce")
		require.NoError(t, os.MkdirAll(targetDir, 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(targetDir, "dummy.txt"),
			[]byte("test"),
			0644,
		))

		os.Args = []string{
			"bubbletea-init",
			"--output-dir", customDir,
			"--force",
			"testoutputforce",
		}
		initialize.Initialize()

		// Check if files were created in the correct location
		mainFile := filepath.Join(customDir, "testoutputforce", "main.go")
		modFile := filepath.Join(customDir, "testoutputforce", "go.mod")

		assert.FileExists(t, mainFile)
		assert.FileExists(t, modFile)
	})

	// Test all flags together
	t.Run("all flags", func(t *testing.T) {
		resetFlags()
		customDir := filepath.Join(projectDir, "custom")
		os.Args = []string{
			"bubbletea-init",
			"--with-bubbles",
			"--mod", "github.com/test/all",
			"--output-dir", customDir,
			"--force",
			"testall",
		}
		initialize.Initialize()

		mainFile := filepath.Join(customDir, "testall", "main.go")
		modFile := filepath.Join(customDir, "testall", "go.mod")

		// Check if files exist
		assert.FileExists(t, mainFile)
		assert.FileExists(t, modFile)

		// Check main.go content
		content, err := os.ReadFile(mainFile)
		require.NoError(t, err)
		assert.Contains(t, string(content), "type spinner struct")
		assert.Contains(t, string(content), "github.com/charmbracelet/lipgloss")

		// Check go.mod content
		modContent, err := os.ReadFile(modFile)
		require.NoError(t, err)
		assert.Contains(t, string(modContent), "module github.com/test/all")
		assert.Contains(t, string(modContent), "github.com/charmbracelet/lipgloss")
	})
}
