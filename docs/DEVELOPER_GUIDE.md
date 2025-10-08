# Developer Guide

Guide for contributing to and developing the PhD Dissertation Planner.

## Table of Contents

1. [Development Setup](#development-setup)
2. [Project Structure](#project-structure)
3. [Architecture](#architecture)
4. [Development Workflow](#development-workflow)
5. [Testing](#testing)
6. [Contributing](#contributing)
7. [Release Process](#release-process)

---

## Development Setup

### Prerequisites

- Go 1.21+
- LaTeX distribution (MiKTeX, MacTeX, or TeX Live)
- Git
- Python 3.8+ (optional, for preview system)
- Make (optional, for convenience commands)

### Clone and Setup

```bash
# Clone repository
git clone https://github.com/yourusername/gantt.git
cd gantt

# Install dependencies
go mod download
go mod tidy

# Install pre-commit hooks (recommended)
pip install pre-commit
pre-commit install

# Build
make build
# or
go build -o plannergen.exe ./cmd/planner
```

### Development Tools

**Recommended:**
- VS Code with Go extension
- GoLand
- Git GUI (GitKraken, SourceTree, or GitHub Desktop)

**Optional:**
- golangci-lint for linting
- gopls for language server
- delve for debugging

---

## Project Structure

```
gantt/
├── cmd/
│   └── planner/          # Main application entry point
│       └── main.go
├── src/
│   ├── app/              # Application logic
│   │   ├── cli.go        # CLI setup
│   │   └── generator.go  # Document generation
│   ├── calendar/         # Calendar generation
│   │   └── calendar.go
│   ├── core/             # Core types and logic
│   │   ├── config.go     # Configuration
│   │   ├── reader.go     # CSV reading
│   │   └── types.go      # Data types
│   └── shared/
│       └── templates/    # LaTeX templates
├── configs/              # Configuration files
│   ├── base.yaml
│   ├── academic.yaml
│   ├── compact.yaml
│   └── presentation.yaml
├── input_data/           # Sample CSV files
├── docs/                 # Documentation
├── tests/                # Test files
├── scripts/              # Utility scripts
├── go.mod                # Go dependencies
├── go.sum                # Dependency checksums
├── Makefile              # Build automation
└── README.md             # Main documentation
```

### Key Directories

**`cmd/planner/`**
- Application entry point
- Minimal code, delegates to `src/app`

**`src/app/`**
- CLI application logic
- Template loading and rendering
- Document generation orchestration

**`src/core/`**
- Core data types (Task, Config, etc.)
- CSV reading and parsing
- Configuration management

**`src/calendar/`**
- Calendar generation logic
- Month/week/day calculations
- LaTeX table generation

**`src/shared/templates/`**
- LaTeX templates
- Embedded in binary using go:embed

---

## Architecture

### Data Flow

```
CSV File → Reader → Tasks → Generator → LaTeX → PDF
                                ↓
                          Templates
```

### Component Diagram

```
┌─────────────┐
│   CLI App   │
└──────┬──────┘
       │
       ├─→ Config Loader
       │
       ├─→ CSV Reader ──→ Tasks
       │
       ├─→ Generator
       │   ├─→ Template Engine
       │   ├─→ Calendar Builder
       │   └─→ LaTeX Compiler
       │
       └─→ Output Writer
```

### Key Components

**1. CLI Application (`src/app/cli.go`)**
- Parses command-line arguments
- Loads configuration
- Orchestrates generation process

**2. Configuration (`src/core/config.go`)**
- YAML-based configuration
- Environment variable support
- Preset management

**3. CSV Reader (`src/core/reader.go`)**
- Parses CSV files
- Validates data
- Error handling and reporting

**4. Generator (`src/app/generator.go`)**
- Template rendering
- Document composition
- LaTeX compilation

**5. Calendar (`src/calendar/calendar.go`)**
- Month/week calculations
- Calendar table generation
- Task placement

**6. Templates (`src/shared/templates/`)**
- LaTeX templates
- Embedded at compile time
- Customizable via filesystem (dev mode)

---

## Development Workflow

### Making Changes

1. **Create a branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes**
   - Edit code
   - Add tests
   - Update documentation

3. **Test locally**
   ```bash
   make test
   make build
   ./plannergen.exe
   ```

4. **Commit**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **Push and create PR**
   ```bash
   git push origin feature/my-feature
   ```

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

**Examples:**
```
feat(calendar): add week number highlighting
fix(reader): handle empty CSV rows gracefully
docs: update installation guide
chore: update dependencies
```

### Code Style

**Go Code:**
- Follow `gofmt` formatting
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

**Example:**
```go
// ReadTasks reads and parses tasks from the CSV file.
// It returns a slice of tasks and any error encountered.
func (r *Reader) ReadTasks() ([]Task, error) {
    // Implementation
}
```

**YAML Configuration:**
- Use 2-space indentation
- Keep keys lowercase with underscores
- Add comments for complex settings

**LaTeX Templates:**
- Use consistent indentation
- Comment complex macros
- Keep templates modular

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -v -race -coverprofile=coverage.txt ./...

# View coverage in browser
go tool cover -html=coverage.txt

# Or use make
make test
make test-coverage
```

### Writing Tests

**Unit Test Example:**
```go
// src/core/reader_test.go
package core

import (
    "testing"
)

func TestReadTasks(t *testing.T) {
    reader := NewReader("testdata/sample.csv")
    tasks, err := reader.ReadTasks()
    
    if err != nil {
        t.Fatalf("ReadTasks failed: %v", err)
    }
    
    if len(tasks) != 5 {
        t.Errorf("Expected 5 tasks, got %d", len(tasks))
    }
}
```

**Table-Driven Test Example:**
```go
func TestParseDate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    time.Time
        wantErr bool
    }{
        {"valid date", "2025-01-15", time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), false},
        {"invalid format", "01/15/2025", time.Time{}, true},
        {"empty string", "", time.Time{}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseDate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !got.Equal(tt.want) {
                t.Errorf("ParseDate() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Coverage Goals

- **Unit tests:** 80%+ coverage
- **Integration tests:** Key workflows
- **Edge cases:** Error conditions, boundary values

---

## Contributing

### Before Contributing

1. Check existing issues
2. Discuss major changes first
3. Read this guide
4. Set up development environment

### Contribution Process

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Add tests**
5. **Update documentation**
6. **Submit a pull request**

### Pull Request Guidelines

**PR Title:**
- Use conventional commit format
- Be descriptive

**PR Description:**
- Explain what and why
- Link related issues
- Include screenshots if UI changes
- List breaking changes

**PR Checklist:**
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Code formatted (`gofmt`)
- [ ] Tests pass locally
- [ ] No merge conflicts
- [ ] Commit messages follow convention

### Code Review Process

1. Automated checks run (tests, linting)
2. Maintainer reviews code
3. Address feedback
4. Approval and merge

---

## Debugging

### Debug Mode

Enable debug output:
```bash
# Set environment variable
export DEBUG=1
./plannergen

# Or in code
cfg.Debug.ShowFrame = true
cfg.Debug.ShowLinks = true
```

### Common Debugging Tasks

**Debug CSV Parsing:**
```go
// Add logging in reader.go
log.Printf("Parsing row %d: %v", lineNum, record)
```

**Debug Template Rendering:**
```go
// Add template debugging
t.Execute(os.Stdout, data) // Print to stdout
```

**Debug LaTeX Compilation:**
```bash
# Compile manually to see errors
cd generated
pdflatex -interaction=nonstopmode planner.tex
# Check planner.log for errors
```

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug ./cmd/planner -- --config configs/base.yaml

# Set breakpoint
(dlv) break src/app/generator.go:100
(dlv) continue
```

---

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- Example: `1.2.3`

**Increment:**
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

### Creating a Release

1. **Update version**
   ```bash
   # Update CHANGELOG.md
   # Update version in code if applicable
   ```

2. **Create tag**
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```

3. **Build binaries**
   ```bash
   # Windows
   GOOS=windows GOARCH=amd64 go build -o plannergen-windows-amd64.exe ./cmd/planner
   
   # Mac
   GOOS=darwin GOARCH=amd64 go build -o plannergen-darwin-amd64 ./cmd/planner
   
   # Linux
   GOOS=linux GOARCH=amd64 go build -o plannergen-linux-amd64 ./cmd/planner
   ```

4. **Create GitHub release**
   - Go to GitHub Releases
   - Create new release from tag
   - Upload binaries
   - Add release notes

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version number updated
- [ ] Tag created
- [ ] Binaries built for all platforms
- [ ] GitHub release created
- [ ] Release notes written

---

## Advanced Topics

### Adding New Templates

1. Create template file in `src/shared/templates/monthly/`
2. Use go:embed to include it
3. Reference in configuration

### Adding New Configuration Options

1. Add field to `Config` struct in `src/core/config.go`
2. Update YAML parsing
3. Add validation
4. Update documentation

### Adding New Output Formats

1. Create new generator in `src/app/`
2. Implement generation logic
3. Add CLI flag
4. Update documentation

---

## Resources

### Go Resources
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

### LaTeX Resources
- [LaTeX Documentation](https://www.latex-project.org/help/documentation/)
- [Overleaf Guides](https://www.overleaf.com/learn)
- [TeX Stack Exchange](https://tex.stackexchange.com/)

### Project Resources
- [GitHub Repository](https://github.com/yourusername/gantt)
- [Issue Tracker](https://github.com/yourusername/gantt/issues)
- [Discussions](https://github.com/yourusername/gantt/discussions)

---

## Getting Help

- 📖 Read the documentation
- 💬 Ask in GitHub Discussions
- 🐛 Report bugs in Issues
- 📧 Email maintainers

---

## License

This project is licensed under the MIT License. See LICENSE file for details.
