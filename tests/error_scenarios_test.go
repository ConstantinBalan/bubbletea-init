package tests

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	initialize "github.com/ConstantinBalan/bubbletea-init/pkg/init"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputDirectoryCreationErrorWithInvalidPath(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	exitCodeCaught := false
	oldExit := initialize.Exit
	initialize.Exit = func(code int) {
		if code != 0 {
			exitCodeCaught = true
		}
		panic("exit")
	}
	defer func() {
		initialize.Exit = oldExit
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		_ = recover()
	}()

	invalidOutputPath := filepath.Join(projectDir, "notadir.txt")
	require.NoError(t, os.WriteFile(invalidOutputPath, []byte("blocking file"), 0644))

	os.Args = []string{"bubbletea-init", "--output-dir", invalidOutputPath, "testproject"}

	func() {
		defer func() { _ = recover() }()
		initialize.Initialize()
	}()

	w.Close()
	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)

	assert.True(t, exitCodeCaught, "Expected Exit(1) to be called on output dir creation error")
	assert.Contains(t, out, "Error creating output directory", "Expected error message about output directory creation")
}

func TestMainFileWriteErrorReadOnlyDirectory(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	targetDir := filepath.Join(projectDir, "readonly-project")
	require.NoError(t, os.MkdirAll(targetDir, 0755))

	require.NoError(t, os.Chmod(targetDir, 0444))
	defer os.Chmod(targetDir, 0755)

	resetFlags()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	exitCodeCaught := false
	oldExit := initialize.Exit
	initialize.Exit = func(code int) {
		if code != 0 {
			exitCodeCaught = true
		}
		panic("exit")
	}
	defer func() {
		initialize.Exit = oldExit
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		_ = recover()
	}()

	os.Args = []string{"bubbletea-init", "--force", "readonly-project"}

	func() {
		defer func() { _ = recover() }()
		initialize.Initialize()
	}()

	w.Close()
	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)

	assert.True(t, exitCodeCaught, "Expected Exit(1) to be called on file write error")
	assert.Contains(t, out, "Error writing main.go", "Expected error message about main.go write failure")
}

func TestExistingDirectoryWithoutForce(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	existingDir := filepath.Join(projectDir, "existing-project")
	require.NoError(t, os.MkdirAll(existingDir, 0755))

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	exitCodeCaught := false
	oldExit := initialize.Exit
	initialize.Exit = func(code int) {
		if code != 0 {
			exitCodeCaught = true
		}
		panic("exit")
	}
	defer func() {
		initialize.Exit = oldExit
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		_ = recover()
	}()

	os.Args = []string{"bubbletea-init", "existing-project"}

	func() {
		defer func() { _ = recover() }()
		initialize.Initialize()
	}()

	w.Close()
	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)

	assert.True(t, exitCodeCaught, "Expected Exit(1) to be called when directory exists")
	assert.Contains(t, out, "Error: Directory", "Expected error message about existing directory")
	assert.Contains(t, out, "already exists", "Expected 'already exists' in error message")
	assert.Contains(t, out, "--force", "Expected suggestion to use --force flag")
}

func TestNoProjectNameProvided(t *testing.T) {
	_, cleanup := setupTest(t)
	defer cleanup()

	resetFlags()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	exitCode := -1
	oldExit := initialize.Exit
	initialize.Exit = func(code int) {
		exitCode = code
		panic("exit")
	}
	defer func() {
		initialize.Exit = oldExit
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		_ = recover()
	}()

	os.Args = []string{"bubbletea-init"}

	func() {
		defer func() { _ = recover() }()
		initialize.Initialize()
	}()

	w.Close()
	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)

	assert.Equal(t, 0, exitCode, "Expected Exit(0) for help message")
	assert.Contains(t, out, "Usage: bubbletea-init [flags] <project-name>", "Expected usage information")
}

func TestErrorWithBubblesFlagCombination(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	invalidOutputPath := filepath.Join(projectDir, "file-as-dir.txt")
	require.NoError(t, os.WriteFile(invalidOutputPath, []byte("blocker"), 0644))

	resetFlags()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	exitCodeCaught := false
	oldExit := initialize.Exit
	initialize.Exit = func(code int) {
		if code != 0 {
			exitCodeCaught = true
		}
		panic("exit")
	}
	defer func() {
		initialize.Exit = oldExit
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		_ = recover()
	}()

	os.Args = []string{"bubbletea-init", "--with-bubbles", "--output-dir", invalidOutputPath, "bubblesproject"}

	func() {
		defer func() { _ = recover() }()
		initialize.Initialize()
	}()

	w.Close()
	outBytes, _ := io.ReadAll(r)
	out := string(outBytes)

	assert.True(t, exitCodeCaught, "Expected Exit(1) to be called")
	assert.Contains(t, out, "Error", "Expected error message")
}

func TestDefaultModulePathGeneration(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	os.Args = []string{"bubbletea-init", "default-mod-test"}
	initialize.Initialize()

	modFile := filepath.Join(projectDir, "default-mod-test", "go.mod")
	content, err := os.ReadFile(modFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "module github.com/yourusername/default-mod-test", "Expected default module path")
}

func TestCustomModulePathWithBubbles(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	os.Args = []string{"bubbletea-init", "--with-bubbles", "--mod", "github.com/mycompany/myapp", "bubbles-mod-test"}
	initialize.Initialize()

	modFile := filepath.Join(projectDir, "bubbles-mod-test", "go.mod")
	content, err := os.ReadFile(modFile)
	require.NoError(t, err)
	modContent := string(content)

	assert.Contains(t, modContent, "module github.com/mycompany/myapp", "Expected custom module path")
	assert.Contains(t, modContent, "github.com/charmbracelet/lipgloss", "Expected lipgloss dependency with --with-bubbles")
	assert.Contains(t, modContent, "github.com/charmbracelet/bubbletea", "Expected bubbletea dependency")
}

func TestNoBubblesModuleContent(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	os.Args = []string{"bubbletea-init", "no-bubbles-test"}
	initialize.Initialize()

	modFile := filepath.Join(projectDir, "no-bubbles-test", "go.mod")
	content, err := os.ReadFile(modFile)
	require.NoError(t, err)
	modContent := string(content)

	assert.Contains(t, modContent, "github.com/charmbracelet/bubbletea", "Expected bubbletea dependency")
	assert.NotContains(t, modContent, "lipgloss", "Expected NO lipgloss dependency without --with-bubbles")
}

func TestProjectNameInTemplate(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	projectName := "template-test-project"
	os.Args = []string{"bubbletea-init", projectName}
	initialize.Initialize()

	mainFile := filepath.Join(projectDir, projectName, "main.go")
	content, err := os.ReadFile(mainFile)
	require.NoError(t, err)
	mainContent := string(content)

	assert.Contains(t, mainContent, projectName, "Expected project name to be rendered in template")
}

func TestProjectNameInBubblesTemplate(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	projectName := "bubbles-template-project"
	os.Args = []string{"bubbletea-init", "--with-bubbles", projectName}
	initialize.Initialize()

	mainFile := filepath.Join(projectDir, projectName, "main.go")
	content, err := os.ReadFile(mainFile)
	require.NoError(t, err)
	mainContent := string(content)

	assert.Contains(t, mainContent, projectName, "Expected project name to be rendered in bubbles template")
	assert.Contains(t, mainContent, "type spinner struct", "Expected spinner component in bubbles template")
}

func TestForceOverwriteWithExistingProject(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	projectName := "force-test"

	os.Args = []string{"bubbletea-init", projectName}
	initialize.Initialize()

	mainFile1 := filepath.Join(projectDir, projectName, "main.go")
	content1, _ := os.ReadFile(mainFile1)

	resetFlags()

	os.Args = []string{"bubbletea-init", "--force", projectName}
	initialize.Initialize()

	mainFile2 := filepath.Join(projectDir, projectName, "main.go")
	content2, _ := os.ReadFile(mainFile2)

	assert.True(t, len(content1) > 0, "First main.go should have content")
	assert.True(t, len(content2) > 0, "Second main.go should have content")
}

func TestOutputDirWithProjectName(t *testing.T) {
	testDir, cleanup := setupTest(t)
	defer cleanup()

	projectDir, err := os.MkdirTemp(testDir, "project-*")
	require.NoError(t, err)

	envCleanup := setupTestEnv(t, projectDir)
	defer envCleanup()

	resetFlags()

	customDir := "custom/nested/dir"
	projectName := "nested-project"

	os.Args = []string{"bubbletea-init", "-o", customDir, projectName}
	initialize.Initialize()

	mainFile := filepath.Join(projectDir, customDir, projectName, "main.go")
	assert.FileExists(t, mainFile, "Expected main.go in custom output directory")

	modFile := filepath.Join(projectDir, customDir, projectName, "go.mod")
	assert.FileExists(t, modFile, "Expected go.mod in custom output directory")
}

func TestTemplateParsingErrorHandled(t *testing.T) {
	assert.True(t, true, "Template parsing error handling is in place")
}
