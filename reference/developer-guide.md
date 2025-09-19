# 🔧 Developer Guide - PhD Dissertation Planner

Complete guide for developers contributing to the PhD Dissertation Planner project.

## 📋 Table of Contents

1. [Development Setup](#development-setup)
2. [Project Architecture](#project-architecture)
3. [Code Organization](#code-organization)
4. [Building and Testing](#building-and-testing)
5. [Contributing Guidelines](#contributing-guidelines)
6. [API Documentation](#api-documentation)
7. [Troubleshooting](#troubleshooting)

## 🚀 Development Setup

### Prerequisites
- **Go 1.16+**
- **XeLaTeX** (for PDF generation)
- **Git** (for version control)
- **Make** (for build automation)

### Environment Setup
```bash
# Clone the repository
git clone <repository-url>
cd phd-dissertation-planner

# Install dependencies
go mod tidy

# Build the application
make build

# Run tests
make test
```

### Development Tools
```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest

# Run linter
golangci-lint run

# Run with hot reload (if using air)
air
```

## 🏗️ Project Architecture

### High-Level Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CSV Input     │───▶│  Data Parser    │───▶│  Calendar       │
│                 │    │                 │    │  Generator      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PDF Output    │◀───│  LaTeX          │◀───│  Template       │
│                 │    │  Compiler       │    │  Engine         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Component Overview

#### Core Components
- **`cmd/plannergen/`**: CLI application entry point
- **`internal/app/`**: Application setup and CLI handling
- **`internal/config/`**: Configuration management
- **`internal/data/`**: CSV data parsing and validation
- **`internal/generator/`**: LaTeX generation logic
- **`internal/calendar/`**: Calendar data structures and algorithms

#### Supporting Components
- **`internal/header/`**: Header components
- **`internal/latex/`**: LaTeX utilities
- **`internal/layout/`**: Layout measurements
- **`templates/`**: LaTeX templates
- **`configs/`**: Configuration files

## 📁 Code Organization

### Directory Structure
```
src/
├── cmd/plannergen/          # CLI application
│   └── main.go
├── internal/                # Private application code
│   ├── app/                # Application setup
│   ├── calendar/           # Calendar functionality
│   ├── config/             # Configuration
│   ├── data/               # Data processing
│   ├── generator/          # PDF generation
│   ├── header/             # Header components
│   ├── latex/              # LaTeX utilities
│   └── layout/             # Layout measurements
├── templates/              # LaTeX templates
├── configs/                # Configuration files
├── scripts/                # Build scripts
└── tests/                  # Test files
```

### Package Responsibilities

#### `internal/app/`
- CLI application setup
- Command-line argument parsing
- Application lifecycle management

#### `internal/data/`
- CSV file parsing
- Data validation
- Task data structures
- Error handling and reporting

#### `internal/calendar/`
- Calendar data structures
- Task layout algorithms
- Conflict resolution
- Multi-day task handling

#### `internal/generator/`
- LaTeX template processing
- PDF generation coordination
- Output file management

## 🔨 Building and Testing

### Build Commands
```bash
# Build the application
make build

# Build with debug information
go build -ldflags="-X main.version=dev" -o build/plannergen ./cmd/plannergen

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o build/plannergen-linux ./cmd/plannergen
```

### Testing
```bash
# Run all tests
make test

# Run specific test package
go test ./internal/data/...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./internal/calendar/...

# Run integration tests
go test -tags=integration ./tests/integration/...
```

### Code Quality
```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...

# Run go vet
go vet ./...

# Check for security issues
gosec ./...
```

## 📝 Contributing Guidelines

### Code Style
- Follow Go standard formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for public APIs
- Keep functions small and focused
- Use interfaces for testability

### Commit Messages
```
type(scope): brief description

Detailed description of changes.

Fixes #123
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Pull Request Process
1. Create a feature branch
2. Make your changes
3. Add tests for new functionality
4. Update documentation
5. Run all tests and linters
6. Submit pull request

### Testing Requirements
- All new code must have tests
- Maintain or improve test coverage
- Include integration tests for new features
- Test error conditions and edge cases

## 📚 API Documentation

### Core APIs

#### Data Package
```go
// Parse CSV data
func (r *Reader) ReadTasks(filename string) ([]*Task, error)

// Validate task data
func (v *Validator) ValidateTask(task *Task) []error

// Categorize tasks
func (m *TaskCategoryManager) GetCategory(name string) (TaskCategory, bool)
```

#### Calendar Package
```go
// Generate calendar
func (g *Generator) GenerateCalendar(tasks []*Task) (*Calendar, error)

// Layout tasks
func (l *LayoutEngine) LayoutTasks(tasks []*Task) (*Layout, error)

// Resolve conflicts
func (r *ConflictResolver) ResolveConflicts(conflicts []Conflict) []Resolution
```

#### Generator Package
```go
// Generate LaTeX
func (g *Generator) GenerateLaTeX(calendar *Calendar) (string, error)

// Compile PDF
func (c *Compiler) CompilePDF(latex string) error
```

## 🔧 Troubleshooting

### Common Development Issues

#### Build Failures
```bash
# Check Go version
go version

# Clean module cache
go clean -modcache

# Reinstall dependencies
go mod download
go mod tidy
```

#### Test Failures
```bash
# Run tests with verbose output
go test -v ./...

# Check test data
ls -la tests/data/

# Run specific test
go test -run TestSpecificFunction ./internal/data/
```

#### LaTeX Compilation Issues
```bash
# Check XeLaTeX installation
xelatex --version

# Test LaTeX compilation
cd build && xelatex test.tex

# Check LaTeX logs
cat build/*.log
```

### Debugging
```bash
# Run with debug logging
DEBUG=1 ./build/plannergen --config configs/base.yaml

# Generate debug output
./scripts/simple.sh ../input/test.csv debug_output

# Check intermediate files
ls -la build/
```

## 📖 Additional Resources

- [API Reference](../api-reference/README.md) - Detailed API documentation
- [Architecture Overview](architecture.md) - System design details
- [Testing Guide](testing.md) - Testing strategies and best practices
- [Performance Guide](performance.md) - Performance optimization tips

---

*For questions about development, check the [API Reference](../api-reference/README.md) or review existing code examples.*
