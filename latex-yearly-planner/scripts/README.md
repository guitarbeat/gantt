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

## Test Scripts

### `test_single.sh`
Quick test with minimal data (single task).

### `test_three.sh`
Test with a small subset of tasks for faster iteration.

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