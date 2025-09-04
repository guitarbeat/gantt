# Build Scripts

This directory contains build and utility scripts for the LaTeX Yearly Planner.

## Main Scripts

### `single.sh`
Core build script that:
- Compiles the Go application
- Generates LaTeX files from templates
- Compiles LaTeX to PDF using XeLaTeX
- Supports preview mode and multiple passes

**Usage:**
```bash
# Basic usage
CFG="configs/base.yaml,configs/page_template.yaml" ./scripts/single.sh

# With environment variables
PLANNER_CSV_FILE="examples/sample_project_data.csv" \
PLANNER_YEAR=2025 \
CFG="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" \
./scripts/single.sh
```

### `build_with_data.sh`
Convenience script for building with CSV data. Pre-configured to use the sample project data.

**Usage:**
```bash
./scripts/build_with_data.sh
```

### `build_release.sh`
Release builder for generating timestamped planner PDFs. Creates clean builds with timestamped filenames in the release/ directory.

**Usage:**
```bash
./scripts/build_release.sh -c "configs/base.yaml,configs/page_template.yaml,configs/csv_config.yaml" -n "overlap_test"
```

## Test Scripts

### `test_day.sh`
Comprehensive test script for the `day.go` module:
- Unit tests for all Day struct methods
- Coverage reports and benchmarks
- Race condition detection
- Multiple test modes with colored output

**Usage:**
```bash
./scripts/test_day.sh -t    # Run unit tests
./scripts/test_day.sh -c    # Run with coverage
./scripts/test_day.sh -b    # Run benchmarks
./scripts/test_day.sh -r    # Run race detection
./scripts/test_day.sh --help # Show all options
```

### `test_triple_csv.sh`
Test script specifically for `test_triple.csv` file processing:
- CSV structure validation
- Date parsing verification
- Go parsing functionality
- Planner binary testing

**Usage:**
```bash
./scripts/test_triple_csv.sh -v    # Validate CSV structure
./scripts/test_triple_csv.sh -p    # Test CSV parsing
./scripts/test_triple_csv.sh -d    # Test date parsing
./scripts/test_triple_csv.sh --help # Show all options
```

## Makefile Integration

For common development tasks, use the Makefile instead of scripts:

```bash
make build          # Build the Go binary
make clean          # Clean build artifacts  
make fmt            # Format Go code
make vet            # Lint Go code
make test-single    # Run single task test (replaces scripts/test_single.sh)
make run-single     # Run single task (replaces scripts/run_single.sh)
make run-csv        # Run with CSV data (replaces scripts/run_with_csv.sh)
make run            # Run with default config
make preview        # Run in preview mode
```

**Use scripts for:**
- Complex builds with specific options
- Release builds with timestamps
- Comprehensive testing
- CSV processing and validation

**Use Makefile for:**
- Standard Go development workflow
- Quick builds and tests
- Code formatting and linting

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PLANNER_CSV_FILE` | Path to CSV data file | None |
| `PLANNER_YEAR` | Base year for planner | Current year |
| `PASSES` | Number of LaTeX compilation passes | 1 |
| `CFG` | Comma-separated list of config files | Required |
| `NAME` | Output PDF filename | Based on config name |
| `PREVIEW` | Generate preview (unique pages only) | false |

## Configuration Files

Scripts typically use these configuration combinations:

1. **Basic**: `configs/base.yaml,configs/page_template.yaml`
2. **Full**: `configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml`
3. **CSV-based**: Automatically includes CSV file path and date range detection

## Output

Generated files are placed in the `build/` directory:
- `*.tex` - Generated LaTeX files
- `*.pdf` - Compiled PDF output
- `*.aux`, `*.log`, etc. - LaTeX compilation artifacts