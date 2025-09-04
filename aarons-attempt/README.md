# 🎯 LaTeX Project Timeline Generator

A LaTeX-first tool that transforms CSV data into publication-quality timelines and Gantt charts. Perfect for PhD research, formal reports, and professional project management.

## ✨ Features

- **CSV → LaTeX** - Convert CSV into complete .tex documents
- **Timeline View** - Professional timeline with task bars and milestones
- **Task List** - Detailed table with all task information
- **PRISMA Diagram** - Optional PRISMA flow diagram for systematic reviews
- **Category Color Coding** - 7 distinct colors for different task categories
- **Professional Styling** - Clean typography and modern design

## 🚀 Quick Start

### Basic Usage
```bash
# Generate LaTeX from CSV with timestamp
make build

# Build and compile to PDF
make build-pdf

# Clean build artifacts
make clean
```

### Advanced Usage
```bash
# Generate with custom title
python main.py --input ../input/data.cleaned.csv --title "My Project" --output output.tex

# Include PRISMA diagram
python main.py --input ../input/data.cleaned.csv --title "Research Timeline" --prisma --output output.tex
```

## 📁 Project Structure

```
aarons-attempt/
├── src/                          # Source code
│   ├── __init__.py              # Package initialization
│   ├── config.py                # Configuration
│   ├── core.py                  # Main application logic
│   └── prisma_generator.py      # PRISMA diagram generation
├── output/tex/                  # Generated LaTeX files
├── main.py                      # Entry point
├── Makefile                     # Build automation
└── README.md                    # This file
```

## 🔧 Command Line Options

```bash
python main.py [OPTIONS]

Options:
  --input FILE          Input CSV file (default: ../input/data.cleaned.csv)
  --output FILE         Output LaTeX file (default: output/tex/Timeline_template.tex)
  --title TITLE         Document title (default: Project Timeline)
  --prisma              Include PRISMA flow diagram
  --verbose, -v         Enable verbose logging
  --quiet, -q           Suppress all output except errors
```

## 🔧 Technical Details

### **Input Format**
CSV with columns: `Task ID`, `Task Name`, `Start Date`, `Due Date`, `Category`, `Dependencies`, `Description`

### **Dependencies**
- Python 3.8+
- LaTeX distribution with `pdflatex` (TeX Live, MiKTeX, MacTeX)

### **Performance**
- Handles 1000+ tasks efficiently
- Vector-style bars and clean typography
- Optimized for large timelines