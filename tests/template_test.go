package tests

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainTemplateParses(t *testing.T) {
	tmpl, err := template.New("main").Parse(getMainTemplate())
	require.NoError(t, err, "main template should parse without error")
	require.NotNil(t, tmpl, "main template should not be nil")
}

func TestMainTemplateExecutes(t *testing.T) {
	tmpl, err := template.New("main").Parse(getMainTemplate())
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{"ProjectName": "TestProject"})
	require.NoError(t, err, "main template should execute without error")

	output := buf.String()

	assert.Contains(t, output, "package main", "Should contain package main")
	assert.Contains(t, output, "github.com/charmbracelet/bubbletea", "Should reference bubbletea")
	assert.Contains(t, output, "TestProject", "Should contain project name")
	assert.Contains(t, output, "func main()", "Should have main function")
}

func TestMainTemplateProjectNameSubstitution(t *testing.T) {
	tmpl, err := template.New("main").Parse(getMainTemplate())
	require.NoError(t, err)

	tests := []struct {
		name        string
		projectName string
	}{
		{"SimpleProject", "myapp"},
		{"WithHyphens", "my-app"},
		{"WithNumbers", "app123"},
		{"CamelCase", "MyAwesomeApp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tmpl.Execute(&buf, map[string]string{"ProjectName": tt.projectName})
			require.NoError(t, err)

			output := buf.String()
			assert.Contains(t, output, tt.projectName, "Project name should appear in output")
		})
	}
}

func TestBubblesTemplateParses(t *testing.T) {
	tmpl, err := template.New("bubbles").Parse(getBubblesTemplate())
	require.NoError(t, err, "bubbles template should parse without error")
	require.NotNil(t, tmpl, "bubbles template should not be nil")
}

func TestBubblesTemplateExecutes(t *testing.T) {
	tmpl, err := template.New("bubbles").Parse(getBubblesTemplate())
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{"ProjectName": "BubblesProject"})
	require.NoError(t, err, "bubbles template should execute without error")

	output := buf.String()

	assert.Contains(t, output, "package main", "Should contain package main")
	assert.Contains(t, output, "github.com/charmbracelet/bubbletea", "Should reference bubbletea")
	assert.Contains(t, output, "github.com/charmbracelet/lipgloss", "Should reference lipgloss")
	assert.Contains(t, output, "BubblesProject", "Should contain project name")
	assert.Contains(t, output, "type spinner struct", "Should have spinner component")
	assert.Contains(t, output, "type textInput struct", "Should have textInput component")
	assert.Contains(t, output, "func main()", "Should have main function")
}

func TestBubblesTemplateProjectNameSubstitution(t *testing.T) {
	tmpl, err := template.New("bubbles").Parse(getBubblesTemplate())
	require.NoError(t, err)

	tests := []struct {
		name        string
		projectName string
	}{
		{"SimpleProject", "bubbles-app"},
		{"WithHyphens", "my-bubbles"},
		{"WithNumbers", "bubbles123"},
		{"CamelCase", "BubblesApp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tmpl.Execute(&buf, map[string]string{"ProjectName": tt.projectName})
			require.NoError(t, err)

			output := buf.String()
			assert.Contains(t, output, tt.projectName, "Project name should appear in output")
		})
	}
}

func TestMainTemplateStructure(t *testing.T) {
	tmpl, err := template.New("main").Parse(getMainTemplate())
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{"ProjectName": "ValidApp"})
	require.NoError(t, err)

	output := buf.String()

	assert.Contains(t, output, "package main", "Should declare package")
	assert.Contains(t, output, "import", "Should have imports")
	assert.Contains(t, output, "func main()", "Should have main function")
	assert.Contains(t, output, "model struct{}", "Should define model type")
	assert.Contains(t, output, "func (m model) Init()", "Should have Init method")
	assert.Contains(t, output, "func (m model) Update(msg tea.Msg)", "Should have Update method")
	assert.Contains(t, output, "func (m model) View() string", "Should have View method")
}

func TestBubblesTemplateStructure(t *testing.T) {
	tmpl, err := template.New("bubbles").Parse(getBubblesTemplate())
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{"ProjectName": "ComponentApp"})
	require.NoError(t, err)

	output := buf.String()

	assert.Contains(t, output, "package main", "Should declare package")
	assert.Contains(t, output, "import", "Should have imports")
	assert.Contains(t, output, "func main()", "Should have main function")
	assert.Contains(t, output, "type model struct", "Should define model type")
	assert.Contains(t, output, "type spinner struct", "Should define spinner component")
	assert.Contains(t, output, "type textInput struct", "Should define textInput component")
	assert.Contains(t, output, "func (m model) Init()", "Should have Init method")
	assert.Contains(t, output, "func (m model) Update(msg tea.Msg)", "Should have Update method")
	assert.Contains(t, output, "func (m model) View() string", "Should have View method")
	assert.Contains(t, output, "func (s spinner) init()", "Should have spinner init")
	assert.Contains(t, output, "func (t textInput) init()", "Should have textInput init")
}

func TestTemplateExecutionWithEdgeCaseNames(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
	}{
		{"SingleChar", "a"},
		{"Numbers", "123"},
		{"UnderscorePrefix", "_private"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("main").Parse(getMainTemplate())
			require.NoError(t, err)

			var buf bytes.Buffer
			err = tmpl.Execute(&buf, map[string]string{"ProjectName": tt.projectName})
			require.NoError(t, err, "Template should execute even with edge case names")

			output := buf.String()
			assert.NotEmpty(t, output, "Output should not be empty")
			assert.Contains(t, output, tt.projectName, "Project name should be in output")
		})
	}
}

func TestMainTemplateOutput(t *testing.T) {
	tmpl, err := template.New("main").Parse(getMainTemplate())
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{"ProjectName": "SyntaxTest"})
	require.NoError(t, err)

	output := buf.String()

	openBraces := bytes.Count([]byte(output), []byte("{"))
	closeBraces := bytes.Count([]byte(output), []byte("}"))
	assert.Equal(t, openBraces, closeBraces, "Braces should be balanced")

	openParens := bytes.Count([]byte(output), []byte("("))
	closeParens := bytes.Count([]byte(output), []byte(")"))
	assert.Equal(t, openParens, closeParens, "Parentheses should be balanced")
}

func TestBubblesTemplateOutput(t *testing.T) {
	tmpl, err := template.New("bubbles").Parse(getBubblesTemplate())
	require.NoError(t, err)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{"ProjectName": "BubblesSyntaxTest"})
	require.NoError(t, err)

	output := buf.String()

	openBraces := bytes.Count([]byte(output), []byte("{"))
	closeBraces := bytes.Count([]byte(output), []byte("}"))
	assert.Equal(t, openBraces, closeBraces, "Braces should be balanced")

	openParens := bytes.Count([]byte(output), []byte("("))
	closeParens := bytes.Count([]byte(output), []byte(")"))
	assert.Equal(t, openParens, closeParens, "Parentheses should be balanced")
}

func getMainTemplate() string {
	return `package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd {
	// Perform any initial setup here
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Hello from {{.ProjectName}}! Press q to quit.\n"
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}`
}

func getBubblesTemplate() string {
	return `package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(4).
	PaddingRight(4)

type model struct {
	spinner  spinner
	input    textInput
	loading  bool
	value    string
	quitting bool
}

func initialModel() model {
	return model{
		spinner: newSpinner(),
		input:   newTextInput(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.init(),
		m.input.init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.input.value != "" {
				m.loading = true
				m.value = m.input.value
				m.input.reset()
				return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
					return loadingFinishedMsg{}
				})
			}
		}
	case loadingFinishedMsg:
		m.loading = false
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.update(msg)
	if !m.loading {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.update(msg)
		return m, tea.Batch(cmd, inputCmd)
	}
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye! 👋\n"
	}

	var s strings.Builder

	s.WriteString(style.Render("{{.ProjectName}}") + "\n\n")

	if m.loading {
		s.WriteString(fmt.Sprintf("%s Loading: %s...\n", m.spinner.view(), m.value))
	} else if m.value != "" {
		s.WriteString(fmt.Sprintf("Last value: %s\n\n", m.value))
		s.WriteString(m.input.view())
	} else {
		s.WriteString(m.input.view())
	}

	s.WriteString("\n\nPress q to quit\n")

	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Spinner component
type spinner struct {
	frames  []string
	current int
}

func newSpinner() spinner {
	return spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

func (s spinner) init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

func (s spinner) update(msg tea.Msg) (spinner, tea.Cmd) {
	switch msg.(type) {
	case spinnerTickMsg:
		s.current = (s.current + 1) % len(s.frames)
		return s, s.init()
	default:
		return s, nil
	}
}

func (s spinner) view() string {
	return s.frames[s.current]
}

// Text input component
type textInput struct {
	value string
}

func newTextInput() textInput {
	return textInput{}
}

func (t textInput) init() tea.Cmd {
	return nil
}

func (t textInput) update(msg tea.Msg) (textInput, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyBackspace:
			if len(t.value) > 0 {
				t.value = t.value[:len(t.value)-1]
			}
		case tea.KeyRunes:
			t.value += string(msg.Runes)
		}
	}
	return t, nil
}

func (t textInput) view() string {
	return fmt.Sprintf("Enter some text: %s█", t.value)
}

func (t *textInput) reset() {
	t.value = ""
}

// Custom messages
type spinnerTickMsg struct{}
type loadingFinishedMsg struct{}`
}
