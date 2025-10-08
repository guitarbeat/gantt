# User Guide

Complete guide to using the PhD Dissertation Planner.

## Table of Contents

1. [Quick Start](#quick-start)
2. [CSV Format](#csv-format)
3. [Configuration](#configuration)
4. [Generating Planners](#generating-planners)
5. [Customization](#customization)
6. [Common Workflows](#common-workflows)
7. [Tips & Best Practices](#tips--best-practices)

---

## Quick Start

### 1. Create Your Timeline CSV

Create a file `input_data/my_timeline.csv`:

```csv
Task,Start Date,End Date,Phase,Category,Priority
Literature Review,2025-01-15,2025-03-31,Phase 1: Foundation,Research,High
Research Design,2025-02-01,2025-03-15,Phase 1: Foundation,Planning,High
Data Collection,2025-03-15,2025-06-30,Phase 2: Execution,Research,High
Data Analysis,2025-06-01,2025-08-31,Phase 2: Execution,Analysis,Medium
Writing Chapter 1,2025-07-01,2025-08-15,Phase 3: Writing,Writing,Medium
Writing Chapter 2,2025-08-15,2025-09-30,Phase 3: Writing,Writing,Medium
Revisions,2025-10-01,2025-11-15,Phase 4: Finalization,Writing,High
Defense Preparation,2025-11-15,2025-12-15,Phase 4: Finalization,Presentation,High
```

### 2. Generate Your Planner

```bash
# Set your CSV file
set PLANNER_CSV_FILE=input_data/my_timeline.csv

# Generate
plannergen

# Output will be in generated/planner.pdf
```

### 3. View Your Planner

Open `generated/planner.pdf` in your PDF viewer.

---

## CSV Format

### Required Columns

| Column | Format | Example | Description |
|--------|--------|---------|-------------|
| `Task` | Text | "Literature Review" | Task name |
| `Start Date` | YYYY-MM-DD | 2025-01-15 | Task start date |
| `End Date` | YYYY-MM-DD | 2025-03-31 | Task end date |

### Optional Columns

| Column | Format | Example | Description |
|--------|--------|---------|-------------|
| `Phase` | Text | "Phase 1: Foundation" | Group tasks by phase |
| `Category` | Text | "Research" | Task category |
| `Priority` | Text | "High", "Medium", "Low" | Task priority |
| `Status` | Text | "Not Started", "In Progress", "Complete" | Current status |
| `Notes` | Text | "Requires IRB approval" | Additional notes |

### CSV Best Practices

**✅ Do:**
- Use YYYY-MM-DD date format
- Keep task names concise (< 50 characters)
- Use consistent phase naming
- Save as UTF-8 encoding

**❌ Don't:**
- Use special LaTeX characters without escaping (%, $, &, #, _, {, })
- Leave empty rows
- Use different date formats
- Include commas in task names (or quote them)

### Example CSV Templates

**Academic Research:**
```csv
Task,Start Date,End Date,Phase,Category
Literature Review,2025-01-01,2025-03-31,Phase 1,Research
Methodology Design,2025-02-01,2025-03-15,Phase 1,Planning
IRB Approval,2025-03-01,2025-04-01,Phase 1,Admin
Data Collection,2025-04-01,2025-07-31,Phase 2,Research
```

**Software Development:**
```csv
Task,Start Date,End Date,Phase,Priority
Requirements Analysis,2025-01-01,2025-01-31,Planning,High
System Design,2025-02-01,2025-02-28,Design,High
Implementation,2025-03-01,2025-06-30,Development,High
Testing,2025-06-01,2025-07-31,QA,Medium
```

---

## Configuration

### Using Presets

The planner includes three built-in presets:

#### Academic (Default)
```bash
plannergen --preset academic
```
- Full task index with phase grouping
- Detailed monthly calendars
- Progress tracking
- Best for: Comprehensive dissertation planning

#### Compact
```bash
plannergen --preset compact
```
- Minimal task index
- Condensed monthly views
- Space-efficient layout
- Best for: Quick reference, printing

#### Presentation
```bash
plannergen --preset presentation
```
- Clean, professional design
- Optimized for slides
- High-contrast colors
- Best for: Advisor meetings, presentations

### Environment Variables

Configure via environment variables:

```bash
# Windows (PowerShell)
$env:PLANNER_CSV_FILE = "input_data/my_timeline.csv"
$env:PLANNER_OUTPUT_DIR = "output"
$env:PLANNER_YEAR = "2025"

# Mac/Linux (Bash)
export PLANNER_CSV_FILE="input_data/my_timeline.csv"
export PLANNER_OUTPUT_DIR="output"
export PLANNER_YEAR="2025"
```

### Command Line Flags

Override settings with flags:

```bash
plannergen \
  --config configs/base.yaml \
  --outdir generated \
  --preset academic
```

Available flags:
- `--config` - Config file path
- `--outdir` - Output directory
- `--preset` - Configuration preset
- `--preview` - Generate preview only
- `--validate` - Validate CSV without generating

---

## Generating Planners

### Basic Generation

```bash
# Use default settings
plannergen

# Specify CSV file
plannergen --config configs/base.yaml

# Custom output directory
plannergen --outdir my_output
```

### Advanced Generation

```bash
# Validate CSV first
plannergen --validate

# Generate with specific preset
plannergen --preset compact --outdir compact_version

# Preview mode (faster, for testing)
plannergen --preview
```

### Batch Generation

Generate multiple versions:

```bash
# Generate all presets
plannergen --preset academic --outdir output/academic
plannergen --preset compact --outdir output/compact
plannergen --preset presentation --outdir output/presentation
```

---

## Customization

### Colors

Tasks are automatically color-coded by phase. Colors cycle through:
- Blue
- Green
- Orange
- Purple
- Red
- Teal

### Layout Adjustments

Edit `configs/base.yaml`:

```yaml
layout:
  week_column_width: 5mm    # Width of week number column
  day_height: 8mm           # Height of each day row
  month_spacing: 10mm       # Space between months
```

### Task Index

Customize the task index appearance:

```yaml
task_index:
  show_progress: true       # Show progress bars
  group_by_phase: true      # Group tasks by phase
  show_dates: true          # Show start/end dates
  show_duration: true       # Show task duration
```

---

## Common Workflows

### Workflow 1: Initial Planning

1. **Brainstorm tasks** - List all major tasks
2. **Estimate dates** - Add realistic start/end dates
3. **Group by phase** - Organize into logical phases
4. **Generate planner** - Create initial version
5. **Review & adjust** - Refine dates and tasks
6. **Regenerate** - Create final version

### Workflow 2: Progress Tracking

1. **Update CSV** - Mark completed tasks
2. **Adjust dates** - Update remaining tasks
3. **Regenerate** - Create updated planner
4. **Compare versions** - Track progress over time

### Workflow 3: Advisor Meetings

1. **Generate presentation version**
   ```bash
   plannergen --preset presentation --outdir meeting_2025_01
   ```
2. **Print or share PDF**
3. **Discuss timeline**
4. **Update based on feedback**
5. **Regenerate academic version**

### Workflow 4: Multiple Timelines

Manage different timeline versions:

```bash
# Optimistic timeline
set PLANNER_CSV_FILE=input_data/timeline_optimistic.csv
plannergen --outdir output/optimistic

# Realistic timeline
set PLANNER_CSV_FILE=input_data/timeline_realistic.csv
plannergen --outdir output/realistic

# Conservative timeline
set PLANNER_CSV_FILE=input_data/timeline_conservative.csv
plannergen --outdir output/conservative
```

---

## Tips & Best Practices

### Planning Tips

**Start with Major Milestones**
- Identify key deadlines (defense date, submission deadline)
- Work backwards from these dates
- Add buffer time for revisions

**Break Down Large Tasks**
- Split multi-month tasks into smaller chunks
- Makes progress tracking easier
- Helps identify dependencies

**Be Realistic**
- Add 20-30% buffer time
- Account for holidays and breaks
- Consider other commitments

**Review Regularly**
- Update timeline monthly
- Adjust based on actual progress
- Keep advisor informed

### CSV Management

**Version Control**
- Keep CSV files in git
- Track changes over time
- Easy to revert if needed

**Backup Regularly**
- Save copies before major changes
- Use descriptive filenames with dates
- Example: `timeline_2025_01_15.csv`

**Use Comments**
- Add notes in a Notes column
- Document assumptions
- Track dependencies

### Output Management

**Organize Output**
```
output/
├── 2025_01_15/
│   ├── planner.pdf
│   └── timeline.csv
├── 2025_02_01/
│   ├── planner.pdf
│   └── timeline.csv
└── current/
    ├── planner.pdf
    └── timeline.csv
```

**Name Versions Clearly**
```bash
# Include date in output directory
plannergen --outdir output/$(date +%Y_%m_%d)

# Or version number
plannergen --outdir output/v1.0
```

### Sharing

**For Advisors**
- Use presentation preset
- Include task index
- Highlight key milestones
- Print in color if possible

**For Committee**
- Use academic preset
- Show all details
- Include phase grouping
- Provide both PDF and CSV

**For Yourself**
- Use compact preset for daily reference
- Print monthly pages
- Keep digital copy accessible
- Update regularly

---

## Keyboard Shortcuts & Quick Commands

### Quick Generation

```bash
# Create alias (Mac/Linux)
alias genplan='plannergen --config configs/base.yaml'

# Use it
genplan

# Windows (PowerShell profile)
function genplan { plannergen --config configs/base.yaml }
```

### Batch Scripts

**Windows (batch file):**
```batch
@echo off
set PLANNER_CSV_FILE=input_data/my_timeline.csv
plannergen.exe --preset academic --outdir output
echo Done! Check output/planner.pdf
pause
```

**Mac/Linux (shell script):**
```bash
#!/bin/bash
export PLANNER_CSV_FILE="input_data/my_timeline.csv"
./plannergen --preset academic --outdir output
echo "Done! Check output/planner.pdf"
```

---

## FAQ

**Q: Can I use dates from different years?**
A: Yes! The planner automatically detects the date range from your CSV.

**Q: How many tasks can I include?**
A: Tested with 100+ tasks. Performance may vary with very large files.

**Q: Can I customize colors?**
A: Currently colors are automatic. Custom colors are planned for future versions.

**Q: Does it work offline?**
A: Yes! All generation happens locally.

**Q: Can I export to other formats?**
A: Currently PDF only. HTML export is planned.

**Q: How do I print the planner?**
A: Open the PDF and print. Use compact preset for better printing.

---

## Getting Help

- 🔧 [Troubleshooting Guide](TROUBLESHOOTING.md)
- 📖 [Setup Guide](SETUP.md)
- 💬 [GitHub Issues](https://github.com/yourusername/gantt/issues)
- 📧 Email support

---

## Next Steps

- ✅ Master the basics
- 🎨 Explore customization options
- 📊 Try different presets
- 🚀 Share with your advisor
- 💡 Provide feedback for improvements
