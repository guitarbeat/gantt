# Pre-commit Hooks Setup Guide

Pre-commit hooks help catch issues before they're committed to the repository.

## Installation

### 1. Install pre-commit

**Using pip (Python):**
```bash
pip install pre-commit
```

**Using Homebrew (Mac):**
```bash
brew install pre-commit
```

**Using Chocolatey (Windows):**
```powershell
choco install pre-commit
```

### 2. Install the hooks

From the repository root:
```bash
pre-commit install
```

This will install the hooks defined in `.pre-commit-config.yaml`.

## What Gets Checked

The pre-commit hooks will automatically check:

### On Every Commit:
- **Trailing whitespace** - Removes trailing spaces
- **End of file** - Ensures files end with a newline
- **YAML syntax** - Validates YAML files
- **Large files** - Prevents committing files >1MB
- **Merge conflicts** - Checks for unresolved conflicts
- **Line endings** - Ensures consistent line endings
- **Go formatting** - Runs `gofmt` on Go files
- **Go vet** - Runs static analysis on Go code

### On Push Only:
- **Go tests** - Runs the test suite

## Usage

Once installed, the hooks run automatically:

```bash
# Hooks run automatically on commit
git commit -m "your message"

# Run hooks manually on all files
pre-commit run --all-files

# Run a specific hook
pre-commit run go-fmt --all-files

# Skip hooks (not recommended)
git commit --no-verify -m "your message"
```

## Updating Hooks

To update to the latest hook versions:
```bash
pre-commit autoupdate
```

## Troubleshooting

### Hooks fail on first run
This is normal - the hooks will fix issues automatically. Just stage the changes and commit again:
```bash
git add -u
git commit -m "your message"
```

### Go commands not found
Ensure Go is installed and in your PATH:
```bash
go version
```

### Python/pip not found
Install Python 3.8+ from https://python.org

### Hooks are slow
You can disable the test hook for faster commits:
```bash
# Edit .pre-commit-config.yaml and comment out the go-test hook
```

## Uninstalling

To remove the hooks:
```bash
pre-commit uninstall
```

## CI/CD Integration

These same checks run in CI/CD, so running them locally helps catch issues early.
