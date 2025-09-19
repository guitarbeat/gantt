# PhD Dissertation Planner - Output Directory

This directory contains all generated outputs from the PhD dissertation planner system.

## Directory Structure

```
output/
├── pdfs/          # Generated PDF files
├── latex/         # LaTeX source files
├── logs/          # Compilation logs
└── README.md      # This file
```

## Files

### PDFs (`pdfs/`)
Generated PDF files from the PhD dissertation planner system.

### LaTeX (`latex/`)
LaTeX source files used to generate the PDFs.

### Logs (`logs/`)
LaTeX compilation logs for debugging and troubleshooting.

## Usage

To generate new outputs, run from the `latex-yearly-planner` directory:

```bash
# Generate test PDF
make test

# Generate demo PDF
make demo

# Generate PDF from custom CSV
make pdf CSV=../input/your_file.csv OUTPUT=your_name

# Or use the script directly
./scripts/simple.sh ../input/your_file.csv your_name
```

All outputs will be automatically organized into the appropriate subdirectories.
