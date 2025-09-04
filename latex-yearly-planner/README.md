# LaTeX Yearly Planner

A Go application for generating LaTeX yearly planners from CSV data.

## Project Structure

This project follows the standard Go project layout:

```text
latex-yearly-planner/
├── cmd/                    # Application entry points
│   └── plannergen/         # Main CLI application
├── internal/               # Private application code
│   ├── app/               # CLI application setup
│   ├── config/            # Configuration management
│   ├── data/              # CSV data reading
│   ├── generator/         # LaTeX generation logic
│   └── layout/            # Layout measurements
├── pkg/                   # Reusable components
│   ├── calendar/          # Calendar data structures
│   ├── header/            # Header components
│   └── latex/             # LaTeX utilities
├── configs/               # Configuration files
├── templates/             # LaTeX templates (embedded)
│   └── monthly/          # Flattened monthly templates (*.tpl)
├── scripts/              # Build scripts
└── build/                # Build output (gitignored)
```

## Quick Start

```bash
# Generate PDF from test data
make test

# Generate PDF from demo data  
make demo

# Generate PDF from custom CSV
make pdf CSV=../input/your_file.csv OUTPUT=my_planner
```

## Building

```bash
make build
```

## Usage

The simplest way to generate a PDF:

```bash
# Quick test with single task
make test

# Demo with multiple tasks
make demo

# Custom CSV file
make pdf CSV=../input/your_data.csv OUTPUT=my_planner
```

## Direct Script Usage

For more control, use the simple script directly:

```bash
# Basic usage
./scripts/simple.sh ../input/test_single.csv my_output

# With custom CSV
./scripts/simple.sh ../input/your_data.csv custom_name
```

## Legacy Support

The old complex script system still works for backward compatibility:

```bash
# Old way (still works)
make run-single
make run-csv
```
- `--preview` (optional): render only one page per unique module
- `--outdir` (optional): output directory for generated files. Overrides the config `OutputDir` and defaults to `build`.

By default, templates are embedded into the binary. For development, you can override to load templates from disk by setting `DEV_TEMPLATES=1`:

```bash
DEV_TEMPLATES=1 make run
```

## Configuration

The application uses YAML configuration files in the `configs/` directory:

- `base.yaml` - Base configuration
- `planner_config.yaml` - Planner-specific settings
- `page_template.yaml` - Page template configuration

## Templates

- All LaTeX templates used at runtime live under `templates/monthly/*.tpl` and are embedded into the binary.
- During development, set `DEV_TEMPLATES=1` to load from the filesystem instead of the embedded FS.
- The main entry is `main_document.tpl`, which includes other templates in that directory.

See `templates/README.md` for an overview of the monthly layout and iteration tips.

## Development

The project is organized to follow Go best practices:

- **cmd/**: Contains the main application entry point
- **internal/**: Private application code that cannot be imported by other projects
- **pkg/**: Public packages that can be imported by other projects
- **configs/**: Configuration files separate from code
- **templates/**: Template files organized by purpose
- **build/**: Build artifacts (gitignored)

Scripts:

- `scripts/build.sh` — unified runner (wraps `scripts/single.sh`)
- `scripts/run_single.sh` — convenience wrapper for the single-sample run
- `scripts/run_with_csv.sh` — convenience wrapper for CSV-driven runs
