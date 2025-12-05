# Bubble Tea Project Initializer

A CLI tool to quickly scaffold new [Bubble Tea](https://github.com/charmbracelet/bubbletea) projects with optional components and styling.

[![Tests](https://github.com/ConstantinBalan/bubbletea-init/actions/workflows/tests.yml/badge.svg)](https://github.com/ConstantinBalan/bubbletea-init/actions/workflows/tests.yml)

[![Coverage Status](https://coveralls.io/repos/github/ConstantinBalan/bubbletea-init/badge.svg)](https://coveralls.io/github/ConstantinBalan/bubbletea-init)

## Features

- Create basic Bubble Tea projects
- Include example components (spinner, text input) with the `--with-bubbles` flag
- Custom module naming with `--mod` flag
- Force overwrite existing projects with `--force` flag
- Specify custom output directory with `--output-dir` or `-o` flag
- Pre-configured with [Lip Gloss](https://github.com/charmbracelet/lipgloss) styling

## Installation

```bash
go install github.com/ConstantinBalan/bubbletea-init@latest
```

## Usage

Basic usage:
```bash
bubbletea-init myproject
```

With example components:
```bash
bubbletea-init --with-bubbles myproject
```

With custom module path:
```bash
bubbletea-init --mod github.com/username/myproject myproject
```

In a specific directory:
```bash
bubbletea-init -o /path/to/projects myproject
```

Force overwrite existing project:
```bash
bubbletea-init --force myproject
```

## Development

To run tests:
```bash
go test -v ./...
```

To run tests with race detection:
```bash
go test -race -v ./...
```
