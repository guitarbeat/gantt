# 👤 User Guide - PhD Dissertation Planner

Complete guide for using the PhD Dissertation Planner to create professional timeline visualizations from your CSV data.

## 📋 Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [CSV Data Format](#csv-data-format)
4. [Generating PDFs](#generating-pdfs)
5. [Configuration](#configuration)
6. [Templates](#templates)
7. [Troubleshooting](#troubleshooting)
8. [Advanced Usage](#advanced-usage)

## 🚀 Installation

### Prerequisites
- **Go 1.16+** (for building from source)
- **XeLaTeX** (for PDF generation)
- **Git** (for cloning the repository)

### Quick Installation
```bash
# Clone the repository
git clone <repository-url>
cd phd-dissertation-planner

# Build the application
make build

# Test the installation
make test
```

## ⚡ Quick Start

### 1. Prepare Your Data
Create a CSV file with your task data (see [CSV Data Format](#csv-data-format) below).

### 2. Generate a PDF
```bash
# From the project root
cd src

# Generate PDF from your CSV
./scripts/simple.sh ../input/your_data.csv my_planner

# Or use the Makefile
make pdf CSV=../input/your_data.csv OUTPUT=my_planner
```

### 3. View Results
- **Main output**: `my_planner.pdf` in the `src/` directory
- **Organized output**: `output/pdfs/my_planner.pdf` in the parent directory

## 📊 CSV Data Format

Your CSV file should contain the following columns:

| Column         | Required | Description              | Example                  |
| -------------- | -------- | ------------------------ | ------------------------ |
| Task Name      | ✅        | Name of the task         | "Literature Review"      |
| Start Date     | ✅        | Task start date          | "2024-01-15"             |
| Due Date       | ✅        | Task end date            | "2024-03-15"             |
| Category       | ✅        | Task category            | "RESEARCH"               |
| Description    | ❌        | Task description         | "Review relevant papers" |
| Priority       | ❌        | Task priority (1-5)      | "3"                      |
| Status         | ❌        | Task status              | "Planned"                |
| Assignee       | ❌        | Task assignee            | "John Doe"               |
| Parent Task ID | ❌        | Parent task ID           | "TASK-001"               |
| Dependencies   | ❌        | Comma-separated task IDs | "TASK-001,TASK-002"      |

### Example CSV
```csv
Task Name,Start Date,Due Date,Category,Description,Priority,Status,Assignee
Literature Review,2024-01-15,2024-03-15,RESEARCH,Review relevant papers,3,Planned,John Doe
Data Collection,2024-02-01,2024-04-01,RESEARCH,Collect experimental data,2,Planned,John Doe
Analysis,2024-03-01,2024-05-01,ANALYSIS,Analyze collected data,1,Planned,John Doe
```

## 📄 Generating PDFs

### Basic Usage
```bash
# Generate PDF with default settings
./scripts/simple.sh ../input/data.csv output_name

# Generate PDF with custom configuration
make pdf CSV=../input/data.csv OUTPUT=custom_name
```

### Available Commands
```bash
# Quick test with sample data
make test

# Demo with multiple tasks
make demo

# Custom CSV file
make pdf CSV=../input/your_file.csv OUTPUT=your_name

# Build the application
make build

# Clean generated files
make clean
```

### Output Files
- **PDF**: `{output_name}.pdf` - Final generated PDF
- **LaTeX**: `output/latex/{output_name}.tex` - LaTeX source (for debugging)
- **Logs**: `output/logs/{output_name}.log` - Compilation logs

## ⚙️ Configuration

The application uses YAML configuration files in the `configs/` directory:

### Main Configuration Files
- **`base.yaml`**: Base application settings
- **`planner_config.yaml`**: Planner-specific settings
- **`page_template.yaml`**: Page template configuration

### Key Settings
```yaml
# Example configuration
output:
  format: "pdf"
  quality: "high"
  
calendar:
  start_date: "2024-01-01"
  end_date: "2024-12-31"
  
tasks:
  max_per_day: 5
  show_dependencies: true
  color_by_category: true
```

## 🎨 Templates

The application uses LaTeX templates for PDF generation:

### Template Structure
```
templates/
├── monthly/                 # Monthly calendar templates
│   ├── main_document.tpl   # Main document template
│   ├── monthly_page.tpl    # Monthly page template
│   └── macros.tpl          # LaTeX macros
```

### Customizing Templates
1. Edit template files in `templates/monthly/`
2. Rebuild the application: `make build`
3. Generate your PDF: `make pdf`

## 🔧 Troubleshooting

### Common Issues

#### PDF Generation Fails
```bash
# Check LaTeX installation
xelatex --version

# Check compilation logs
cat output/logs/your_file.log

# Try with verbose output
DEV_TEMPLATES=1 make pdf CSV=../input/data.csv OUTPUT=debug
```

#### CSV Parsing Errors
```bash
# Check CSV format
head -5 ../input/your_file.csv

# Validate required columns
# Ensure dates are in YYYY-MM-DD format
```

#### Build Errors
```bash
# Clean and rebuild
make clean
make build

# Check Go installation
go version
```

### Getting Help
1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Review [Examples](examples/README.md) for similar use cases
3. Check compilation logs in `output/logs/`

## 🚀 Advanced Usage

### Batch Processing
```bash
# Process multiple CSV files
for file in ../input/*.csv; do
    ./scripts/simple.sh "$file" "$(basename "$file" .csv)"
done
```

### Custom Output Directory
```bash
# Generate to custom directory
OUTPUT_DIR=/path/to/output ./scripts/simple.sh ../input/data.csv my_planner
```

### Development Mode
```bash
# Load templates from filesystem (for development)
DEV_TEMPLATES=1 make pdf CSV=../input/data.csv OUTPUT=dev_test
```

## 📚 Additional Resources

- [Examples](examples/README.md) - Real-world usage examples
- [Templates](templates/README.md) - Template customization guide
- [Developer Guide](../developer-guide/README.md) - For contributors
- [API Reference](../api-reference/README.md) - Technical documentation

---

*Need help? Check the [Troubleshooting Guide](troubleshooting.md) or review the [Examples](examples/README.md).*
