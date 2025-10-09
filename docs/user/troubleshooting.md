# 🐛 Troubleshooting Guide

Common issues and solutions for the PhD Dissertation Planner.

## PDF Generation Issues

### "xelatex not found" Error

**Problem:** PDF generation fails with "xelatex: command not found"

**Solution:** Install XeLaTeX:

```bash
# macOS (with Homebrew)
brew install mactex

# Ubuntu/Debian
sudo apt-get install texlive-xetex texlive-latex-extra

# CentOS/RHEL/Fedora
sudo yum install texlive-xetex
# or
sudo dnf install texlive-xetex
```

### LaTeX Compilation Errors

**Problem:** LaTeX compilation fails with package errors

**Solution:** Install additional LaTeX packages:

```bash
# Install comprehensive LaTeX distribution
# macOS
brew install mactex

# Ubuntu
sudo apt-get install texlive-full

# Or install specific packages
sudo apt-get install texlive-fonts-recommended texlive-fonts-extra
```

### Font Issues

**Problem:** PDF generates but fonts look wrong

**Solution:** XeLaTeX uses system fonts. Ensure Helvetica is available:

```bash
# Check available fonts
fc-list | grep -i helvetica

# Install if missing (Ubuntu)
sudo apt-get install fonts-liberation
```

## CSV Data Issues

### "No CSV file found" Error

**Problem:** Application can't find input data

**Solutions:**

1. **Check file location:**
   ```bash
   ls -la input_data/
   ```

2. **Set environment variable:**
   ```bash
   export PLANNER_CSV_FILE="input_data/your_file.csv"
   ```

3. **Use command line flag:**
   ```bash
   go run ./cmd/planner --csv-file input_data/your_file.csv
   ```

### CSV Format Errors

**Problem:** CSV parsing fails with format errors

**Required CSV columns:**
- `Task` - Task name/description
- `StartDate` - Start date (YYYY-MM-DD)
- `EndDate` - End date (YYYY-MM-DD)
- `Phase` - Project phase number
- `SubPhase` - Phase description
- `Category` - Task category

**Common issues:**
- Date format must be YYYY-MM-DD
- Missing required columns
- Extra whitespace in headers

### Data Validation Errors

**Problem:** CSV validates but data seems wrong

**Check data quality:**
```bash
# Validate CSV without generating PDF
go run ./cmd/planner --validate
```

This will show:
- Total task count
- Date range coverage
- Phase distribution
- Data quality metrics

## Configuration Issues

### Configuration Not Loading

**Problem:** Custom configuration file ignored

**Solutions:**

1. **Check file path:**
   ```bash
   go run ./cmd/planner --config path/to/config.yaml
   ```

2. **Validate YAML syntax:**
   ```bash
   # Use online YAML validator or
   python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"
   ```

3. **Check file permissions:**
   ```bash
   ls -la config.yaml
   ```

### Preset Issues

**Problem:** Preset layouts not working

**Available presets:** `academic`, `compact`, `presentation`

```bash
# Use preset
go run ./cmd/planner --preset compact

# List available presets
go run ./cmd/planner --help
```

## Build Issues

### Go Build Errors

**Problem:** `go build` fails

**Solutions:**

1. **Check Go version:**
   ```bash
   go version  # Should be 1.22+
   ```

2. **Clean and rebuild:**
   ```bash
   make clean
   make build
   ```

3. **Check dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

### Test Failures

**Problem:** Tests fail unexpectedly

**Debug steps:**
```bash
# Run specific test
go test -v ./src/core -run TestConfig

# Run with race detector
go test -race ./...

# Check test coverage
go test -cover ./src/core
```

## Runtime Issues

### Memory Issues

**Problem:** Application uses too much memory or crashes

**Solutions:**

1. **Check input data size:** Large CSV files may cause issues
2. **Use smaller date ranges:** Limit to 1-2 years of data
3. **Monitor memory usage:**
   ```bash
   go run ./cmd/planner --debug-memory
   ```

### Performance Issues

**Problem:** PDF generation is slow

**Optimization tips:**

1. **Use compact preset:** Faster rendering
   ```bash
   go run ./cmd/planner --preset compact
   ```

2. **Reduce date range:** Smaller time periods generate faster

3. **Check system resources:** Ensure adequate RAM and CPU

## Development Issues

### Git Hook Issues

**Problem:** Pre-commit hooks fail

**Solutions:**

1. **Reinstall hooks:**
   ```bash
   ./scripts/setup.sh
   ```

2. **Check hook permissions:**
   ```bash
   ls -la .git/hooks/
   chmod +x .git/hooks/*
   ```

3. **Manual testing:**
   ```bash
   make lint
   make test
   ```

### IDE Issues

**Problem:** Go extension or IDE not working properly

**Solutions:**

1. **Go modules:**
   ```bash
   go mod tidy
   go mod download
   ```

2. **IDE settings:** Ensure Go tools are in PATH

3. **Restart IDE** and reload Go modules

## Getting Help

If these solutions don't resolve your issue:

1. **Check existing issues** on GitHub
2. **Create a new issue** with:
   - Error messages (full output)
   - Your environment (OS, Go version, LaTeX version)
   - Steps to reproduce
   - Sample input data (anonymized)

3. **Include debug information:**
   ```bash
   go run ./cmd/planner --debug
   ```

---

*Last updated: October 2025*