# Generated Files

This directory contains the compiled binary and generated PDF outputs.

## Files

### Binary

- `plannergen` - The compiled planner generator binary

### PDF Outputs

- `phd_dissertation_planner.pdf` - **Current version** with fixed row heights and all 107 tasks
- `phd_dissertation_planner_old.pdf` - Previous version (for comparison)

## How to Generate

### Quick Start

```bash
# From repository root
make clean-build
```

### Manual Generation

```bash
# 1. Build the binary
go build -mod=vendor -o generated/plannergen ./cmd/planner

# 2. Generate LaTeX files
PLANNER_CSV_FILE="input_data/research_timeline_v5_comprehensive.csv" \
./generated/plannergen \
  --config src/core/base.yaml,src/core/monthly_calendar.yaml \
  --outdir output

# 3. Compile PDF
cd output
xelatex monthly_calendar.tex
xelatex monthly_calendar.tex  # Run twice for proper references
```

## Features

The generated calendar includes:

- ✅ Fixed row heights for consistent layout
- ✅ Single-day and multi-day task support
- ✅ Hyperlinked navigation between months
- ✅ Professional typography and spacing
- ✅ Automatic task layout and collision avoidance

## Customization

To adjust row height, edit `src/core/defaults.go`:

```go
MonthlyCellHeight: "5em"  // Try 3em, 4em, 5em, or 6em
```

Or in your YAML config:

```yaml
layout:
  latex:
    monthly_cell_height: "5em"
```
