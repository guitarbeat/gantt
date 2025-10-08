# Code Quality & Optimization

This document describes the automatic code quality tools and optimization workflows integrated into the PhD Dissertation Planner project.

## 🤖 Automatic Optimization Tools

The project includes several automatic optimization tools that run in CI/CD and can be used locally:

### Tools Used

1. **goimports** - Import optimization and cleanup
2. **gofmt -s** - Code simplification
3. **go mod tidy** - Dependency cleanup
4. **golangci-lint** - Comprehensive linting with auto-fix
5. **go build -ldflags="-s -w"** - Binary optimization

### What Gets Optimized

- **Import statements** - Organized and cleaned up
- **Code constructs** - Simplified where possible
- **Dependencies** - Unused dependencies removed
- **Linting issues** - Automatically fixed where possible
- **Binary size** - Optimized for smaller executable

## 🚀 Usage

### Local Development

#### Using Makefile (Recommended)
```bash
# Install optimization tools
make install-tools

# Run all optimizations
make optimize

# Fix only linting issues
make lint-fix

# Format code only
make format

# Run quality checks (no changes)
make test-quality
```

#### Using Pre-commit Hooks
```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run on all files
pre-commit run --all-files
```

#### Manual Commands
```bash
# Import optimization
goimports -w src/

# Code simplification
gofmt -s -w src/

# Dependency cleanup
go mod tidy

# Linting with auto-fix
golangci-lint run --fix src/...

# Build with optimizations
go build -ldflags="-s -w" -o plannergen cmd/planner/main.go
```

### CI/CD Integration

#### Automatic Workflows

1. **Main CI Pipeline** (`.github/workflows/ci.yml`)
   - Runs on every push and PR
   - Includes optimization in lint job
   - Fails if code needs optimization

2. **Code Quality Workflow** (`.github/workflows/code-quality.yml`)
   - Manual trigger via GitHub Actions
   - Can create PRs with optimizations
   - Comprehensive optimization suite

3. **Auto-Optimize Job**
   - Runs on PRs and pushes
   - Automatically creates PRs with optimizations
   - Ensures code stays optimized

## 📊 Quality Checks

### What Gets Checked

- **Code formatting** - gofmt compliance
- **Import organization** - goimports compliance
- **Linting issues** - golangci-lint rules
- **Dependency cleanliness** - go mod tidy compliance
- **Test coverage** - All tests must pass
- **Binary optimization** - Size and performance

### Quality Gates

The CI pipeline will fail if:
- Code is not properly formatted
- Imports are not optimized
- Linting issues exist
- Tests fail
- Dependencies are not clean

## 🔧 Configuration

### golangci-lint Configuration

The project uses golangci-lint with the following enabled linters:
- `errcheck` - Error handling
- `gosimple` - Code simplification
- `ineffassign` - Ineffectual assignments
- `staticcheck` - Static analysis
- `unused` - Unused code detection

### Pre-commit Configuration

Located in `.pre-commit-config.yaml`:
- goimports on Go files
- gofmt on Go files
- go mod tidy on go.mod/go.sum
- go vet on all packages
- go test on src packages

## 🎯 Benefits

### Automatic Benefits
- **Consistent code style** across the project
- **Reduced technical debt** through automatic fixes
- **Better performance** through optimizations
- **Cleaner dependencies** through automatic cleanup
- **Smaller binaries** through build optimizations

### Developer Benefits
- **No manual formatting** - tools handle it
- **Automatic PR creation** - CI creates optimization PRs
- **Quality gates** - prevents bad code from being merged
- **Easy local optimization** - simple make commands

## 🚨 Troubleshooting

### Common Issues

1. **Linting errors after optimization**
   ```bash
   # Run linter to see specific issues
   golangci-lint run src/...
   
   # Fix automatically where possible
   golangci-lint run --fix src/...
   ```

2. **Import issues**
   ```bash
   # Fix imports
   goimports -w src/
   ```

3. **Formatting issues**
   ```bash
   # Format code
   gofmt -s -w src/
   ```

4. **Dependency issues**
   ```bash
   # Clean dependencies
   go mod tidy
   ```

### Getting Help

- Check the CI logs for specific error messages
- Run `make test-quality` locally to see issues
- Use `golangci-lint run --help` for linter options
- Check the GitHub Actions logs for detailed output

## 📈 Metrics

The optimization tools help maintain:
- **Zero linting errors** in CI
- **Consistent formatting** across all files
- **Clean dependencies** in go.mod
- **Optimized binary size** (typically 4-5MB)
- **100% test coverage** for critical paths

## 🔄 Continuous Improvement

The optimization tools are continuously updated:
- New linters added as needed
- Rules refined based on project needs
- Performance optimizations applied
- Best practices enforced automatically

This ensures the codebase maintains high quality standards with minimal manual intervention.
