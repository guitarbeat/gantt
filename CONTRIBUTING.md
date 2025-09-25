# Contributing to PhD Dissertation Planner

Thank you for your interest in contributing to the PhD Dissertation Planner! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Development Setup

1. **Prerequisites**
   - Go 1.16 or later
   - XeLaTeX (for PDF generation)
   - Git

2. **Setup**
   ```bash
   # Clone the repository
   git clone <repository-url>
   cd gantt

   # Run the setup script
   ./scripts/setup.sh
   ```

3. **Build**
   ```bash
   # Simple build
   make

   # Or use the enhanced build script
   ./scripts/build.sh --clean
   ```

## 📋 Development Workflow

### 1. Choose an Issue
- Check the [issue tracker](../../issues) for open tasks
- Comment on the issue to indicate you're working on it
- Create a new branch for your work

### 2. Development
- Follow the existing code style and structure
- Write clear commit messages
- Test your changes thoroughly

### 3. Testing
```bash
# Run tests (when implemented)
go test ./...

# Build and verify PDF generation
make
```

### 4. Pull Request
- Ensure your branch is up-to-date with main
- Write a clear PR description
- Reference any related issues
- Request review from maintainers

## 🏗️ Project Structure

```
├── cmd/                    # Application entry points
│   └── planner/           # Main application
├── internal/              # Private application code
│   ├── application/       # Application logic
│   ├── common/           # Shared utilities
│   └── scheduler/        # Calendar/task logic
├── pkg/                   # Reusable libraries
│   └── templates/        # Template system
├── configs/               # Configuration files
├── data/                  # Input data files
├── docs/                  # Documentation
├── examples/              # Example files
├── assets/                # Static assets (PDFs, docs)
├── scripts/               # Build/deployment scripts
└── reference/            # Project-specific reference materials
```

## 📝 Code Guidelines

### Go Code
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Write meaningful variable and function names
- Add comments for complex logic

### Commit Messages
- Use present tense: "Add feature" not "Added feature"
- Start with a capital letter
- Keep the first line under 50 characters
- Add detailed description if needed

### Documentation
- Update README.md for user-facing changes
- Add code comments for complex functions
- Update this CONTRIBUTING.md for process changes

## 🐛 Reporting Issues

When reporting bugs, please include:
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Any error messages

## 💡 Feature Requests

For new features, please:
- Check if the feature is already requested
- Describe the problem you're trying to solve
- Provide examples of how it would work
- Consider alternative approaches

## 📞 Getting Help

- 📧 **Discussions**: Use [GitHub Discussions](../../discussions) for questions
- 🐛 **Bug Reports**: [Issues](../../issues) with the "bug" label
- ✨ **Feature Requests**: [Issues](../../issues) with the "enhancement" label

Thank you for contributing to the PhD Dissertation Planner! 🎓
