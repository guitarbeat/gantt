# PhD Dissertation Planner

A Go-based tool for generating LaTeX/PDF calendars from CSV task data.

## Quick Start

```bash
# Generate calendar
./plannergen

# Validate CSV files
./plannergen --validate
```

## Directory Structure

```
├── input_data/          # All input files
│   ├── config.yaml      # Configuration
│   └── *.csv            # Task data files
│
├── output_data/         # All generated output
│   ├── latex/           # Generated LaTeX files
│   ├── pdfs/            # Generated PDF files
│   ├── auxiliary/       # LaTeX auxiliary files
│   └── binaries/        # Binary outputs
│
├── internal/            # Internal packages
│   ├── app/             # Application logic
│   ├── calendar/        # Calendar generation
│   ├── core/            # Core functionality
│   └── templates/       # LaTeX templates
│
├── main.go              # Entry point
├── plannergen           # Compiled binary
└── Makefile             # Build automation
```

## Input Files

### CSV Files (input_data/)

- `proposal_and_setup.csv` - Proposal & Setup (40 tasks)
- `research_and_experiments.csv` - Research & Experiments (29 tasks)
- `publications.csv` - Publications (13 tasks)
- `dissertation_and_defense.csv` - Dissertation & Defense (26 tasks)

**Total: 108 tasks**

Each file is independent with no cross-file dependencies.

### Configuration (input_data/config.yaml)

Controls calendar layout, styling, and output settings.

## Output Files

All generated files go to `output_data/`:

- **latex/** - LaTeX source files
- **pdfs/** - Compiled PDF calendars
- **auxiliary/** - LaTeX auxiliary files (.aux, .log, etc.)
- **binaries/** - Binary outputs

## Usage

### Generate Calendar

```bash
./plannergen
```

Output will be in `output_data/pdfs/`

### Validate Input

```bash
./plannergen --validate
```

### Custom Output Directory

```bash
./plannergen --outdir custom_output
```

## Building

```bash
# Using Make (recommended)
make build

# Or directly with Go
go build -o plannergen main.go
```

## Make Commands

```bash
make build      # Build the binary
make run        # Build and run
make validate   # Build and validate CSV files
make clean      # Remove binary and output
make deps       # Install dependencies
make help       # Show all commands
```

## Features

- ✅ Automatic CSV file merging
- ✅ In-memory processing (no temporary files)
- ✅ Independent task files
- ✅ Descriptive phase names
- ✅ Clean input/output separation
- ✅ LaTeX/PDF generation
- ✅ Task validation
- ✅ Configurable layouts

## Requirements

- Go 1.21+
- LaTeX distribution (MiKTeX/MacTeX/TeX Live)
- XeLaTeX compiler

## License

See LICENSE file for details.
