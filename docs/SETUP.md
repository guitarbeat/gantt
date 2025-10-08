# Installation & Setup Guide

Complete guide to installing and setting up the PhD Dissertation Planner.

## Prerequisites

### Required Software

1. **Go 1.21 or higher**
   - Download: https://golang.org/dl/
   - Verify: `go version`

2. **LaTeX Distribution**
   - **Windows:** MiKTeX or TeX Live
   - **Mac:** MacTeX
   - **Linux:** TeX Live
   - Verify: `pdflatex --version`

### Optional Software

3. **Python 3.8+** (for preview images)
   - Download: https://python.org
   - Verify: `python --version`

4. **Git** (for version control)
   - Download: https://git-scm.com
   - Verify: `git --version`

---

## Installation Steps

### 1. Install Go

**Windows:**
```powershell
# Using winget
winget install GoLang.Go

# Or download installer from https://golang.org/dl/
```

**Mac:**
```bash
# Using Homebrew
brew install go

# Or download installer from https://golang.org/dl/
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang-go

# Or download from https://golang.org/dl/
```

Verify installation:
```bash
go version
# Should output: go version go1.21.x ...
```

### 2. Install LaTeX

**Windows (MiKTeX):**
```powershell
# Using winget
winget install MiKTeX.MiKTeX

# Or download from https://miktex.org/download
```

**Mac (MacTeX):**
```bash
# Using Homebrew
brew install --cask mactex

# Or download from https://www.tug.org/mactex/
```

**Linux (TeX Live):**
```bash
# Ubuntu/Debian
sudo apt-get install texlive-xetex texlive-fonts-recommended texlive-latex-extra

# Fedora
sudo dnf install texlive-scheme-medium
```

Verify installation:
```bash
pdflatex --version
# Should output: pdfTeX version information
```

### 3. Clone the Repository

```bash
# Clone the repository
git clone https://github.com/yourusername/gantt.git
cd gantt

# Or download ZIP from GitHub and extract
```

### 4. Install Go Dependencies

```bash
# Download all dependencies
go mod download

# Verify dependencies
go mod tidy
```

### 5. Build the Application

```bash
# Build the binary
go build -o plannergen.exe ./cmd/planner

# Or use make (if available)
make build
```

Verify build:
```bash
# Windows
.\plannergen.exe --help

# Mac/Linux
./plannergen --help
```

### 6. Optional: Install Python Dependencies (for preview images)

```bash
# Install required packages
pip install pdf2image Pillow

# Install poppler (required by pdf2image)
```

**Windows:**
- Download poppler from: https://github.com/oschwartz10612/poppler-windows/releases/
- Extract and add to PATH

**Mac:**
```bash
brew install poppler
```

**Linux:**
```bash
sudo apt-get install poppler-utils
```

---

## Configuration

### 1. Prepare Your CSV File

Create a CSV file in `input_data/` directory:

```csv
Task,Start Date,End Date,Phase,Category,Priority
Literature Review,2025-01-01,2025-03-31,Phase 1,Research,High
Data Collection,2025-02-01,2025-04-30,Phase 1,Research,High
Analysis,2025-04-01,2025-06-30,Phase 2,Analysis,Medium
```

**Required columns:**
- `Task` - Task name
- `Start Date` - Format: YYYY-MM-DD
- `End Date` - Format: YYYY-MM-DD

**Optional columns:**
- `Phase` - Group tasks by phase
- `Category` - Task category
- `Priority` - High, Medium, Low

### 2. Choose a Configuration Preset

The planner comes with three presets:

**Academic (default):**
```bash
plannergen --preset academic
```
- Full task index with progress tracking
- Detailed monthly views
- Best for comprehensive planning

**Compact:**
```bash
plannergen --preset compact
```
- Minimal task index
- Condensed monthly views
- Best for quick reference

**Presentation:**
```bash
plannergen --preset presentation
```
- Clean, professional layout
- Optimized for presentations
- Best for sharing with advisors

### 3. Custom Configuration (Advanced)

Create a custom YAML config file:

```yaml
# configs/custom.yaml
year: 2025
output_dir: "generated"
csv_file: "input_data/my_timeline.csv"

layout:
  week_column_width: 5mm
  day_height: 8mm
  
pages:
  - name: "monthly"
    render_blocks:
      - func_name: "monthly"
        templates: ["calendar.tpl"]
```

Use custom config:
```bash
plannergen --config configs/custom.yaml
```

---

## First Run

### Generate Your First Planner

```bash
# Using default settings
plannergen

# Or specify CSV file
set PLANNER_CSV_FILE=input_data/my_timeline.csv
plannergen

# Or use command line
plannergen --config configs/base.yaml --outdir generated
```

### Expected Output

```
🚀 Starting Planner Generation
═══════════════════════════════════════
📋 Loading configuration... ✅
📁 Setting up output directory... ✅
📄 Generating root document... ✅
📅 Generating calendar pages... ✅
═══════════════════════════════════════
✨ Generation complete!
📂 Output: generated
```

### Check the Output

```bash
# View generated files
ls generated/

# Files created:
# - planner.tex (LaTeX source)
# - planner.pdf (Final PDF)
# - monthly_*.tex (Monthly pages)
```

---

## Verification

### Test the Installation

Run the test suite:
```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.txt ./...

# Or use make
make test
```

### Generate Sample Planner

```bash
# Use the included sample CSV
set PLANNER_CSV_FILE=input_data/research_timeline_v5_comprehensive.csv
plannergen
```

---

## Troubleshooting

### Common Issues

**"go: command not found"**
- Go is not installed or not in PATH
- Solution: Install Go and restart terminal

**"pdflatex: command not found"**
- LaTeX is not installed or not in PATH
- Solution: Install LaTeX distribution

**"CSV file not found"**
- CSV file path is incorrect
- Solution: Check file exists in `input_data/` directory

**"Permission denied"**
- No write permissions for output directory
- Solution: Run with appropriate permissions or change output directory

For more issues, see [Troubleshooting Guide](TROUBLESHOOTING.md).

---

## Next Steps

1. ✅ Installation complete
2. 📖 Read the [User Guide](USER_GUIDE.md)
3. 🎨 Customize your planner
4. 📊 Generate your timeline
5. 🚀 Share with your advisor

---

## Updating

### Update the Application

```bash
# Pull latest changes
git pull origin main

# Rebuild
go build -o plannergen.exe ./cmd/planner

# Or use make
make build
```

### Update Dependencies

```bash
# Update Go modules
go get -u ./...
go mod tidy

# Update Python packages
pip install --upgrade pdf2image Pillow
```

---

## Uninstalling

To remove the planner:

```bash
# Remove the repository
cd ..
rm -rf gantt

# Optional: Remove Go (if not needed for other projects)
# Windows: Use "Add or Remove Programs"
# Mac: brew uninstall go
# Linux: sudo apt-get remove golang-go
```

---

## Getting Help

- 📖 [User Guide](USER_GUIDE.md)
- 🔧 [Troubleshooting](TROUBLESHOOTING.md)
- 💬 [GitHub Issues](https://github.com/yourusername/gantt/issues)
- 📧 Email: your-email@example.com
