# Quick Start Guide

Get up and running in 60 seconds.

## Prerequisites

- Go 1.21+
- LaTeX distribution (XeLaTeX)

## Installation

```bash
# Clone or download the repository
cd phd-dissertation-planner

# Build the binary
make build
```

## Usage

### Generate Your Calendar

```bash
# Generate calendar from CSV files
make run

# Or run directly
./plannergen
```

Output will be in `output_data/pdfs/`

### Validate Your Data

```bash
# Check CSV files for errors
make validate
```

### Edit Your Tasks

1. Open any CSV file in `input_data/`:
   - `proposal_and_setup.csv`
   - `research_and_experiments.csv`
   - `publications.csv`
   - `dissertation_and_defense.csv`

2. Edit tasks (add, remove, modify)

3. Run `make run` to regenerate

### Customize Layout

Edit `input_data/config.yaml` to change:
- Paper size
- Fonts and colors
- Task styling
- Calendar layout

## File Structure

```
input_data/     → Your CSV files and config
output_data/    → Generated PDFs and LaTeX
```

## Common Commands

```bash
make build      # Build binary
make run        # Generate calendar
make validate   # Check CSV files
make clean      # Remove output
make help       # Show all commands
```

## CSV Format

Each CSV file contains tasks with these columns:

- **Phase** - Descriptive phase name
- **Task ID** - Unique identifier (e.g., T1.1)
- **Dependencies** - Comma-separated task IDs
- **Task** - Task name
- **Start Date** - YYYY-MM-DD
- **End Date** - YYYY-MM-DD
- **Objective** - Task description
- **Milestone** - true/false
- **Status** - planned, in progress, completed
- **Notes** - Additional notes
- **Category** - Task category
- **Priority** - High, Medium, Low
- **Assignee** - Person responsible
- **Resources** - Required resources

## Tips

1. **Keep files independent** - Tasks should only depend on tasks in the same file
2. **Use descriptive phase names** - No numbers needed (e.g., "Proposal & Setup")
3. **Validate often** - Run `make validate` after editing
4. **Check output** - PDFs are in `output_data/pdfs/`

## Troubleshooting

### "No CSV files found"
- Check that CSV files are in `input_data/`
- Ensure files have `.csv` extension

### "LaTeX compilation failed"
- Install XeLaTeX (part of TeX Live, MacTeX, or MiKTeX)
- Check `output_data/latex/*.log` for errors

### "Validation failed"
- Check error messages for specific issues
- Verify CSV format matches expected columns
- Ensure task IDs are unique

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Customize `input_data/config.yaml` for your needs
- Add your own tasks to the CSV files
- Generate your personalized calendar!

## Support

For issues or questions, check the documentation or review the code in `internal/`.
