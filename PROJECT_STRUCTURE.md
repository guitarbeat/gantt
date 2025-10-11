# 📁 Project Structure Overview

## 🎯 Core Directories

### Source Code
- `src/` - Main application source code
  - `app/` - Application logic and CLI
  - `core/` - Core utilities and configuration
  - `calendar/` - Calendar generation and layout
  - `shared/` - Shared templates and utilities

### Configuration
- `configs/` - YAML configuration files
- `cmd/` - Application entry points

### Data
- `input_data/` - CSV input files and data
- `generated/` - Generated output files
  - `pdfs/` - Generated PDF files
  - `tex/` - LaTeX source files
  - `logs/` - Build and error logs
  - `preview/` - Preview builds

### Documentation
- `docs/` - Project documentation
  - `tasks/` - How-to guides
  - `fyi/` - Reference information
  - `archive/` - Archived documents
  - `examples/` - Example configurations

### Releases
- `releases/` - Versioned releases with timestamps

### Scripts
- `scripts/` - Build and utility scripts

### Tests
- `tests/` - Test files and test data
  - `integration/` - Integration tests
  - `unit/` - Unit tests

### Temporary
- `.temp/` - Temporary files (gitignored)
- `.build_artifacts/` - Build artifacts (gitignored)

## 🧹 Cleanup Commands

```bash
# Run full cleanup
./scripts/cleanup_and_organize.sh

# Clean only scattered files
./scripts/cleanup_and_organize.sh --scattered-only

# Clean only test artifacts
./scripts/cleanup_and_organize.sh --test-only
```

## 📝 File Organization Rules

1. **Generated files** go in `generated/` with appropriate subdirectories
2. **Temporary files** go in `.temp/` (gitignored)
3. **Documentation** goes in `docs/` with logical subdirectories
4. **Configuration** goes in `configs/`
5. **Source code** goes in `src/` with clear module separation
6. **Tests** go in `tests/` with appropriate subdirectories
