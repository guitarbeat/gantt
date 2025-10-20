# Pre-commit Hooks Setup Guide

Pre-commit hooks help catch issues before they're committed to the repository. This project supports both traditional pre-commit hooks and Cursor CLI-powered hooks for AI-enhanced code quality checks.

## 🤖 Cursor CLI Hooks (Recommended)

The project now includes Cursor CLI-powered pre-commit hooks that provide AI-enhanced code quality checks and automated fixes.

### 1. Install Cursor CLI

**Using npm:**
```bash
npm install -g @cursor/cli
```

**Or visit:** https://cursor.com/docs/cli/headless

### 2. Install Cursor CLI Hooks

From the repository root:
```bash
make install-cursor-hooks
```

This will install AI-powered pre-commit hooks that automatically fix formatting, run tests, and provide intelligent code analysis.

## 📋 Traditional Pre-commit Hooks (Legacy)

If you prefer traditional pre-commit hooks or Cursor CLI is not available:

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

### 🤖 Cursor CLI Hooks (AI-Enhanced)

The Cursor CLI hooks provide intelligent code quality checks:

**On Every Commit:**
- **AI-Powered Formatting** - Automatically fixes Go code formatting using Cursor CLI
- **Smart Linting** - Uses Cursor CLI to analyze and suggest fixes for Go vet issues
- **Intelligent Testing** - Runs tests with AI-powered failure analysis
- **YAML Validation** - Validates YAML configuration files
- **Large File Detection** - Prevents committing files >1MB
- **Merge Conflict Detection** - Checks for unresolved conflicts
- **Code Quality Analysis** - AI-powered suggestions for code improvements

**Key Benefits:**
- Automatic code fixes without manual intervention
- AI-powered analysis of test failures and code issues
- Intelligent suggestions for code improvements
- Faster feedback loop with automated corrections

### 📋 Traditional Hooks (Legacy)

The traditional pre-commit hooks check:

**On Every Commit:**
- **Trailing whitespace** - Removes trailing spaces
- **End of file** - Ensures files end with a newline
- **YAML syntax** - Validates YAML files
- **Large files** - Prevents committing files >1MB
- **Merge conflicts** - Checks for unresolved conflicts
- **Line endings** - Ensures consistent line endings
- **Go formatting** - Runs `gofmt` on Go files
- **Go vet** - Runs static analysis on Go code

**On Push Only:**
- **Go tests** - Runs the test suite

## Usage

### 🤖 Cursor CLI Hooks

Once installed, the Cursor CLI hooks run automatically:

```bash
# Hooks run automatically on commit
git commit -m "your message"

# Test hooks without committing
make test-cursor-hooks

# Run Cursor CLI checks manually
make cursor-precommit

# Skip hooks (not recommended)
git commit --no-verify -m "your message"
```

### 📋 Traditional Hooks

Once installed, the traditional hooks run automatically:

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

### 🤖 Cursor CLI Hooks
To update Cursor CLI hooks:
```bash
# Reinstall to get latest version
make uninstall-cursor-hooks
make install-cursor-hooks
```

### 📋 Traditional Hooks
To update to the latest hook versions:
```bash
pre-commit autoupdate
```

## Troubleshooting

### 🤖 Cursor CLI Hooks

**Cursor CLI not found:**
```bash
# Install Cursor CLI
npm install -g @cursor/cli

# Or visit: https://cursor.com/docs/cli/headless
```

**Hooks fail on first run:**
This is normal - the hooks will fix issues automatically. Just stage the changes and commit again:
```bash
git add -u
git commit -m "your message"
```

**Test the installation:**
```bash
make test-cursor-hooks
```

### 📋 Traditional Hooks

**Hooks fail on first run:**
This is normal - the hooks will fix issues automatically. Just stage the changes and commit again:
```bash
git add -u
git commit -m "your message"
```

**Go commands not found:**
Ensure Go is installed and in your PATH:
```bash
go version
```

**Python/pip not found:**
Install Python 3.8+ from https://python.org

**Hooks are slow:**
You can disable the test hook for faster commits:
```bash
# Edit .pre-commit-config.yaml and comment out the go-test hook
```

## Uninstalling

### 🤖 Cursor CLI Hooks
To remove Cursor CLI hooks:
```bash
make uninstall-cursor-hooks
```

### 📋 Traditional Hooks
To remove traditional hooks:
```bash
pre-commit uninstall
```

## CI/CD Integration

These same checks run in CI/CD, so running them locally helps catch issues early. The Cursor CLI hooks provide additional AI-powered analysis that can help identify potential issues before they reach the CI pipeline.
