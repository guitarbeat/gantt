# Contributing to PhD Dissertation Planner

Thank you for your interest in contributing to the PhD Dissertation Planner! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Documentation](#documentation)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project follows a code of conduct to ensure a welcoming environment for all contributors. By participating, you agree to:

- Be respectful and inclusive
- Focus on constructive feedback
- Accept responsibility for mistakes
- Show empathy towards other contributors
- Help create a positive community

## Getting Started

### Prerequisites

- **Go**: Version 1.19 or later
- **XeLaTeX**: For PDF generation (TeX Live or MiKTeX)
- **Git**: Version control system
- **Python 3**: For PDF preview generation (optional)

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/phd-dissertation-planner.git
cd phd-dissertation-planner

# Run setup script
./scripts/maintenance/setup.sh

# Verify installation
make test
```

## Development Setup

### Environment Setup

The project includes automated setup scripts for different platforms:

```bash
# Linux/macOS
./scripts/maintenance/setup.sh

# Windows PowerShell
.\scripts\maintenance\setup.ps1
```

### Development Environment

```bash
# Start development environment with hot reloading
make dev-air

# Or run manually
make dev
```

### IDE Configuration

The project works with any Go-compatible IDE. Recommended settings:

- Go modules enabled
- Format on save
- Run goimports on save
- Use golangci-lint for linting

## Project Structure

```
├── cmd/planner/          # Application entry point
├── src/
│   ├── app/             # CLI application logic
│   ├── core/            # Core business logic and configuration
│   ├── calendar/        # Calendar generation and layout
│   └── shared/          # Shared utilities and templates
├── configs/             # YAML configuration files
├── scripts/             # Build and development scripts
│   ├── build/          # Build scripts
│   ├── dev/            # Development scripts
│   ├── maintenance/    # Maintenance scripts
│   └── docs/          # Script documentation
├── tests/               # Test files
├── docs/                # Documentation
└── generated/           # Generated output files
```

## Development Workflow

### 1. Choose an Issue

- Check the [TODO.md](TODO.md) for current priorities
- Look for issues labeled `good first issue` or `help wanted`
- Comment on the issue to indicate you're working on it

### 2. Create a Branch

```bash
# Create and switch to a feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/issue-number-description
```

### 3. Make Changes

```bash
# Run tests frequently
make test

# Format code
make fmt

# Run linter
make lint

# Build and test
make build
```

### 4. Test Your Changes

```bash
# Run all tests
make test-coverage

# Run specific tests
go test ./src/core/... -v

# Build and preview
./scripts/build/build_and_preview.sh 3
```

## Coding Standards

### Go Code Style

This project follows the Google Go Style Guide with some modifications:

- Use `gofmt` for formatting
- Use `goimports` for import organization
- Maximum line length: 120 characters
- Use meaningful variable names
- Add comments for exported functions and types

### Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

Examples:
```
feat: add progress tracking for tasks
fix: correct week column width on Windows
docs: update cross-platform setup guide
```

### Code Review Guidelines

- Ensure tests pass
- Code is well-documented
- Follows existing patterns
- No breaking changes without discussion
- Performance considerations addressed

## Testing

### Test Structure

Tests are organized by package:

```
tests/
├── unit/                 # Unit tests
│   ├── reader_test.go   # Core reader tests
│   └── validation_test.go # Validation tests
└── integration/          # Integration tests
    └── build_process_test.go
```

### Running Tests

```bash
# All tests
make test-coverage

# Unit tests only
go test ./tests/unit/... -v

# Integration tests only
go test ./tests/integration/... -v

# Benchmarks
make bench
```

### Writing Tests

- Use table-driven tests for multiple test cases
- Test both success and error cases
- Include edge cases
- Use descriptive test names
- Add comments explaining complex test scenarios

Example:

```go
func TestReader_ReadTasks(t *testing.T) {
    tests := []struct {
        name     string
        csvData  string
        wantErr  bool
        wantLen  int
    }{
        {
            name: "valid csv with tasks",
            csvData: `Task Name,Start Date,End Date,Category,Description
Test Task,2024-01-15,2024-01-20,Test,A test task`,
            wantErr: false,
            wantLen: 1,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation...
        })
    }
}
```

## Submitting Changes

### Pull Request Process

1. **Ensure tests pass**: All tests must pass before submitting
2. **Update documentation**: Update relevant docs if needed
3. **Write clear description**: Explain what changes and why
4. **Reference issues**: Link to related issues with `#issue-number`
5. **Request review**: Ask specific contributors for review

### Pull Request Template

Please use this template for pull requests:

```markdown
## Description
Brief description of the changes made.

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change
- [ ] Documentation update
- [ ] Refactoring

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Cross-platform testing done

## Screenshots (if applicable)
Add screenshots of UI changes or output examples.

## Checklist
- [ ] Code follows project style guidelines
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes
```

## Documentation

### Documentation Structure

```
docs/
├── user/                 # User-facing documentation
│   ├── user-guide.md    # Main user guide
│   ├── setup.md         # Setup instructions
│   └── troubleshooting.md # Common issues
├── developer/            # Developer documentation
│   ├── developer-guide.md # Development guide
│   └── pre-commit-setup.md # Pre-commit setup
├── reference/            # Reference documentation
│   ├── api-reference.md # API documentation
│   ├── architecture.md  # System architecture
│   └── configuration.md # Configuration reference
└── examples/             # Example configurations
```

### Updating Documentation

- Keep documentation current with code changes
- Use clear, concise language
- Include code examples where helpful
- Test instructions on multiple platforms
- Update screenshots when UI changes

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

1. **Clear title**: Describe the issue concisely
2. **Steps to reproduce**: Detailed steps to reproduce the issue
3. **Expected behavior**: What should happen
4. **Actual behavior**: What actually happens
5. **Environment**: OS, Go version, XeLaTeX version
6. **Logs**: Relevant error messages or log output
7. **Screenshots**: If applicable

### Feature Requests

For feature requests, please include:

1. **Clear description**: What feature you want
2. **Use case**: Why you need this feature
3. **Proposed solution**: How you think it should work
4. **Alternatives**: Other approaches considered

### Issue Labels

Common labels used in this project:

- `bug`: Something isn't working
- `enhancement`: New feature or improvement
- `documentation`: Documentation updates
- `good first issue`: Suitable for new contributors
- `help wanted`: Community help needed
- `question`: Question or discussion
- `wontfix`: Will not be implemented

## Recognition

Contributors are recognized in several ways:

- GitHub contributor statistics
- Mention in release notes
- Attribution in documentation
- Community recognition

Thank you for contributing to the PhD Dissertation Planner! Your efforts help make this tool better for the academic community.
