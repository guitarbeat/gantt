# Build Scripts

This directory contains the simplified build script for the LaTeX Yearly Planner.

## Main Script

### `simple.sh`
The single, simple script that handles all PDF generation:
- Builds the Go application if needed
- Generates LaTeX files from templates
- Compiles LaTeX to PDF using XeLaTeX
- Creates clean, named output files

**Usage:**
```bash
# Basic usage
./scripts/simple.sh ../input/test_single.csv my_planner

# With custom CSV and output name
./scripts/simple.sh ../input/your_data.csv custom_name
```

## Makefile Integration

For the simplest usage, use the Makefile:

```bash
# Quick test with single task
make test

# Demo with multiple tasks  
make demo

# Custom CSV file
make pdf CSV=../input/your_file.csv OUTPUT=my_planner

# Build the binary
make build
```

## What Was Removed

The following complex scripts were removed for simplicity:
- `generate.sh` - Complex multi-layer script
- `build.sh` - Environment variable wrapper
- `single.sh` - Core build with hardcoded filenames
- `build_release.sh` - Release builder
- `build_with_data.sh` - Data-specific builder
- `test_day.sh` - Day module tests
- `test_triple_csv.sh` - CSV-specific tests

## Why Simplified?

The old system had:
- ❌ 4 script layers deep
- ❌ 8+ environment variables
- ❌ Hardcoded filename mismatches
- ❌ Complex file management
- ❌ Multiple failure points

The new system has:
- ✅ 1 simple script (25 lines)
- ✅ 0 environment variables needed
- ✅ Consistent filename handling
- ✅ Clear error messages
- ✅ Easy to debug and modify

## Output

Generated files are placed in the current directory:
- `{output_name}.pdf` - Final PDF output
- `build/` - Temporary build files (LaTeX artifacts)