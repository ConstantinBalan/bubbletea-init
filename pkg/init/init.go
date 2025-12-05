package init

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/pflag"
)

// Exit is used instead of direct calls to os.Exit so tests can override it.
var Exit = os.Exit

//go:embed templates/main.go.tmpl
var mainTemplate string

//go:embed templates/main_with_bubbles.go.tmpl
var bubblesTemplate string

type templateData struct {
	ProjectName string
}

var (
	style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(1).
		PaddingRight(1)
)

func Initialize() {
	withBubbles := pflag.Bool("with-bubbles", false, "Include example bubble components (spinner, textinput)")
	modPath := pflag.String("mod", "", "Custom Go module name")
	outputDir := pflag.StringP("output-dir", "o", "", "Directory where the project should be created (default: current directory)")
	force := pflag.Bool("force", false, "Overwrite existing files")
	help := pflag.BoolP("help", "h", false, "Show help message")

	pflag.Parse()

	if *help || pflag.NArg() < 1 {
		fmt.Println("Usage: bubbletea-init [flags] <project-name>")
		fmt.Println("\nFlags:")
		pflag.PrintDefaults()
		Exit(0)
	}

	projectName := pflag.Arg(0)
	var projectDir string

	if *outputDir != "" {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory '%s': %v\n", *outputDir, err)
			Exit(1)
		}
		projectDir = filepath.Join(*outputDir, projectName)
	} else {
		projectDir = filepath.Join(".", projectName)
	}

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) && !*force {
		fmt.Printf("Error: Directory '%s' already exists. Use --force to overwrite.\n", projectDir)
		Exit(1)
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		fmt.Printf("Error creating project directory '%s': %v\n", projectDir, err)
		Exit(1)
	}

	data := templateData{
		ProjectName: projectName,
	}

	// Choose template based on flags
	templateContent := mainTemplate
	if *withBubbles {
		templateContent = bubblesTemplate
	}

	tmpl, err := template.New("main").Parse(templateContent)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		Exit(1)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Println("Error executing template:", err)
		Exit(1)
	}

	mainFile := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(mainFile, buf.Bytes(), 0644); err != nil {
		fmt.Println("Error writing main.go:", err)
		Exit(1)
	}

	modName := *modPath
	if modName == "" {
		modName = fmt.Sprintf("github.com/%s/%s", "yourusername", projectName)
	}

	var goModContent string
	if *withBubbles {
		goModContent = fmt.Sprintf(`module %s

go 1.23

require (
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/charmbracelet/lipgloss v0.9.1
)
`, modName)
	} else {
		goModContent = fmt.Sprintf(`module %s

go 1.23

require github.com/charmbracelet/bubbletea v0.25.0
`, modName)
	}

	if err := os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		fmt.Println("Error writing go.mod:", err)
		Exit(1)
	}

	successMsg := style.Render("✅ Success!")
	fmt.Printf("\n%s Bubble Tea project '%s' created successfully!\n", successMsg, projectName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  go run .")
}
