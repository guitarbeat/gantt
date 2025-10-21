# PhD Dissertation Planner

A Go-based tool for generating LaTeX/PDF calendars from CSV task data with hierarchical task organization and visual timeline management.

## Table of Contents

- [Quick Start](#quick-start)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Generate Calendar](#generate-calendar)
  - [Validate Data](#validate-data)
  - [Customize Layout](#customize-layout)
- [Directory Structure](#directory-structure)
- [Input Files](#input-files)
  - [CSV Files](#csv-files)
  - [Configuration](#configuration)
- [Output Files](#output-files)
- [Building](#building)
- [Make Commands](#make-commands)
- [CSV Format](#csv-format)
- [Tips & Best Practices](#tips--best-practices)
- [Troubleshooting](#troubleshooting)
- [Requirements](#requirements)
- [Additional Documentation](#additional-documentation)
- [License](#license)

---

## Quick Start

Get up and running in 60 seconds:

```bash
# 1. Build the binary
make build

# 2. Generate your calendar
make run

# 3. Find your PDF
open output_data/pdfs/config.pdf
```

That's it! Your calendar is ready.

---

## Features

- ✅ **Automatic CSV file merging** - Combine multiple task files seamlessly
- ✅ **In-memory processing** - No temporary files needed
- ✅ **Independent task files** - Edit files separately, merge automatically
- ✅ **Hierarchical task index** - 3-level organization (Sections → Phases → Tasks)
- ✅ **Color-coded phases** - Visual consistency across document
- ✅ **Clickable navigation** - Hyperlinks between index and calendar
- ✅ **Visual indicators** - Milestones (★) and completed tasks (✓)
- ✅ **Clean separation** - Input/output directories clearly organized
- ✅ **LaTeX/PDF generation** - Professional typesetting
- ✅ **Task validation** - Catch errors before generation
- ✅ **Configurable layouts** - Customize appearance and styling

---

## Installation

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **LaTeX distribution** with XeLaTeX:
  - macOS: [MacTeX](https://www.tug.org/mactex/)
  - Windows: [MiKTeX](https://miktex.org/)
  - Linux: [TeX Live](https://www.tug.org/texlive/)

### Build

```bash
# Clone or download the repository
cd phd-dissertation-planner

# Build the binary
make build

# Verify installation
./plannergen --help
```

---

## Usage

### Generate Calendar

Generate your calendar from CSV files:

```bash
# Using Make (recommended)
make run

# Or run directly
./plannergen

# Custom output directory
./plannergen --outdir custom_output
```

**Output location:** `output_data/pdfs/config.pdf`

### Validate Data

Check your CSV files for errors before generating:

```bash
# Validate all CSV files
make validate

# Or run directly
./plannergen --validate
```

Validation checks:
- CSV format and structure
- Required columns present
- Valid dates
- Unique task IDs
- No circular dependencies

### Customize Layout

Edit `input_data/config.yaml` to customize:

- **Paper size** - Letter, A4, custom dimensions
- **Fonts and colors** - Typography and color schemes
- **Task styling** - Box heights, spacing, opacity
- **Calendar layout** - Day cells, margins, spacing
- **Debug options** - Show frames, links for development

After editing, regenerate with `make run`.

---

## Directory Structure

```
phd-dissertation-planner/
├── input_data/              # All input files
│   ├── config.yaml          # Configuration
│   ├── proposal_and_setup.csv
│   ├── research_and_experiments.csv
│   ├── publications.csv
│   └── dissertation_and_defense.csv
│
├── output_data/             # All generated output
│   ├── latex/               # LaTeX source files
│   ├── pdfs/                # Compiled PDF files
│   ├── auxiliary/           # LaTeX auxiliary files
│   └── binaries/            # Binary outputs
│
├── internal/                # Application code
│   ├── app/                 # Application logic
│   ├── calendar/            # Calendar generation
│   ├── core/                # Core functionality
│   └── templates/           # LaTeX templates
│
├── main.go                  # Entry point
├── plannergen               # Compiled binary
├── Makefile                 # Build automation
└── README.md                # This file
```

**Design principles:**
- Clear input/output separation
- Minimal hierarchy (max 3 levels)
- Logical grouping of related files
- No redundancy or scattered configs

---

## Input Files

### CSV Files

Four independent CSV files in `input_data/`:

| File | Tasks | Description |
|------|-------|-------------|
| `proposal_and_setup.csv` | 40 | Proposal, setup, equipment |
| `research_and_experiments.csv` | 29 | Research aims and experiments |
| `publications.csv` | 13 | Publications and tools |
| `dissertation_and_defense.csv` | 26 | Writing and defense |

**Total: 108 tasks across 17 phases**

Each file is completely independent:
- No cross-file dependencies
- Edit separately
- Merge automatically at build time

### Configuration

`input_data/config.yaml` controls:

- **Core settings** - Week start, time format, output directory
- **Debug options** - Show frames, links
- **Layout** - Paper size, margins, fonts
- **Task styling** - Colors, heights, spacing, opacity
- **Calendar layout** - Day cells, task positioning

See the file for detailed comments and options.

---

## Output Files

All generated files go to `output_data/`:

```
output_data/
├── latex/          # LaTeX source files (.tex)
├── pdfs/           # Compiled PDF calendars
├── auxiliary/      # LaTeX auxiliary files (.aux, .log, .fls, .synctex.gz)
└── binaries/       # Binary outputs
```

**Main output:** `output_data/pdfs/config.pdf`

---

## Building

### Using Make (Recommended)

```bash
make build      # Build the binary
make run        # Build and run
make validate   # Build and validate CSV files
make clean      # Remove binary and output
make deps       # Install dependencies
make help       # Show all commands
```

### Using Go Directly

```bash
# Build
go build -o plannergen main.go

# Run
./plannergen

# With flags
./plannergen --config input_data/config.yaml --outdir output_data
```

---

## Make Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the binary |
| `make run` | Build and generate calendar |
| `make validate` | Build and validate CSV files |
| `make clean` | Remove binary and output files |
| `make deps` | Install Go dependencies |
| `make test` | Run tests |
| `make help` | Show all available commands |

---

## CSV Format

Each CSV file contains tasks with these columns:

| Column | Description | Example |
|--------|-------------|---------|
| **Phase** | Descriptive phase name | "PhD Proposal" |
| **Task ID** | Unique identifier | "T1.1" |
| **Dependencies** | Comma-separated task IDs | "T1.1,T1.2" |
| **Task** | Task name | "Write Proposal" |
| **Start Date** | YYYY-MM-DD format | "2025-09-01" |
| **End Date** | YYYY-MM-DD format | "2025-09-15" |
| **Objective** | Task description | "Complete proposal draft" |
| **Milestone** | true/false | "true" |
| **Status** | planned, in progress, completed | "in progress" |
| **Notes** | Additional notes | "Review with advisor" |
| **Category** | Task category | "PhD Proposal" |
| **Priority** | High, Medium, Low | "High" |
| **Assignee** | Person responsible | "Student" |
| **Resources** | Required resources | "Writing Tools" |

**Example row:**
```csv
PhD Proposal,T1.1,,Write Proposal,2025-09-01,2025-09-15,Complete proposal draft,true,in progress,Review with advisor,PhD Proposal,High,Student,Writing Tools
```

---

## Tips & Best Practices

### 1. Keep Files Independent
- Tasks should only depend on tasks in the same file
- Makes files easier to manage and edit
- Prevents complex cross-file dependencies

### 2. Use Descriptive Phase Names
- No numbers needed (e.g., "Proposal & Setup" not "1: Proposal")
- Makes the index more readable
- Easier to understand at a glance

### 3. Validate Often
- Run `make validate` after editing CSV files
- Catches errors early
- Saves time debugging LaTeX issues

### 4. Check Output
- PDFs are in `output_data/pdfs/`
- LaTeX source in `output_data/latex/` for debugging
- Logs in `output_data/auxiliary/` if compilation fails

### 5. Use Milestones
- Mark important tasks as milestones
- They appear with ★ in the index
- Helps track major achievements

### 6. Track Progress
- Update task status (planned → in progress → completed)
- Completed tasks show with ✓ in gray
- Provides visual progress tracking

---

## Troubleshooting

### "No CSV files found"

**Problem:** CSV files not detected

**Solutions:**
- Check that CSV files are in `input_data/`
- Ensure files have `.csv` extension
- Verify files are not hidden (don't start with `.`)

### "LaTeX compilation failed"

**Problem:** PDF generation fails

**Solutions:**
- Install XeLaTeX (part of TeX Live, MacTeX, or MiKTeX)
- Check `output_data/auxiliary/*.log` for specific errors
- Verify LaTeX is in your PATH: `which xelatex`
- Try compiling manually: `cd output_data/latex && xelatex config.tex`

### "Validation failed"

**Problem:** CSV data has errors

**Solutions:**
- Check error messages for specific issues
- Verify CSV format matches expected columns
- Ensure task IDs are unique within each file
- Check date formats (YYYY-MM-DD)
- Verify no circular dependencies

### "Configuration error"

**Problem:** Config file has issues

**Solutions:**
- Check YAML syntax in `input_data/config.yaml`
- Ensure proper indentation (use spaces, not tabs)
- Verify all required fields are present
- Check for typos in field names

### Build Issues

**Problem:** Compilation fails

**Solutions:**
- Update Go: `go version` (need 1.21+)
- Update dependencies: `make deps`
- Clean and rebuild: `make clean && make build`
- Check for vendor issues: `go mod vendor`

---

## Requirements

### Software

- **Go 1.21 or higher** - [Download](https://golang.org/dl/)
- **LaTeX distribution** with XeLaTeX compiler:
  - macOS: MacTeX
  - Windows: MiKTeX  
  - Linux: TeX Live

### System

- **Operating System:** macOS, Linux, or Windows
- **Disk Space:** ~500MB for LaTeX distribution
- **Memory:** 2GB RAM minimum

### Optional

- **Make** - For using Makefile commands (usually pre-installed on macOS/Linux)
- **Git** - For version control

---

## Additional Documentation

- **[DOCUMENTATION.md](DOCUMENTATION.md)** - Complete technical documentation
  - Project structure details
  - Task index implementation
  - Color coding system
  - CSV file organization
  - Phase ordering logic
  - Visual features and statistics

---

## License

See LICENSE file for details.

---

## Support

For issues, questions, or contributions:

1. Check this README and [DOCUMENTATION.md](DOCUMENTATION.md)
2. Review the code in `internal/` for implementation details
3. Check `output_data/auxiliary/*.log` for LaTeX errors
4. Validate your CSV files with `make validate`

---

**Happy planning! 📅**
