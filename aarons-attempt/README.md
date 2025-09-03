# 🎯 LaTeX Project Timeline Generator

A LaTeX-first tool that turns your CSV into publication‑quality timelines. Perfect for PhD research, formal reports, and advisor meetings.

## ✨ Features

### 🎨 **Output**
- **Publication Quality** - Professional typography and page layout using LaTeX
- **Clean Vertical Timeline** - TikZ-based timeline with milestone and task markers
- **Formal Tables** - Longtable for multi-page task listings
- **Color Coding** - Consistent category/status colors for clarity
- **Portrait Orientation** - Readable layout designed for printing

### 🔧 **Functionality**
- **CSV → LaTeX** - Convert CSV into a complete .tex document
- **Automatic Categorization** - Derives categories from `Group` and `Deliverable Type`
- **Status & Priority** - Colorized status and clear priority cues in tables
- **No Truncation** - Full text preserved in tables and labels

### 🎯 **Timeline Design**
- **Enhanced Vertical Timeline** - Professional timeline with connection dots and visual hierarchy
- **Modern Task Cards** - Rich task information with shadows, better typography, and improved spacing
- **Category Color Coding** - 7 distinct colors for different task categories:
  - 🟣 **Purple** - Milestones
  - 🔵 **Blue** - Research Core (Proposals, Dissertation)
  - 🟢 **Green** - Experimental (Imaging, Laser work)
  - 🟠 **Orange** - Publications & Outputs
  - ⚫ **Gray** - Administrative tasks
  - 🟣 **Violet** - Meetings & Accountability
  - 🩷 **Pink** - Service & BOGO activities
- **Enhanced Status Indicators** - Thicker color-coded stripes (green=completed, orange=in progress, red=blocked)
- **Priority Markers** - Larger red triangles for high priority tasks
- **Dependency Indicators** - Enhanced red dots with white centers for tasks with dependencies
- **Milestone Diamonds** - Larger diamond shapes with better visual prominence
- **Today Marker** - Modern pink marker with enhanced styling and better visibility
- **Professional Layout** - Optimized spacing, proper margins, and visual hierarchy

### 🎨 **Task Information**
- **Smart Text Processing** - Removes redundant prefixes like "Milestone:", "Draft ", "Complete "
- **Enhanced CSV Usage** - Uses all CSV fields: groups, deliverable types, dependencies, notes
- **Intelligent Categorization** - Automatically categorizes tasks by research type and purpose
- **Enhanced Date Labels** - Bold, clear date stamps with better positioning
- **Group Information** - Task groups and deliverable types displayed with improved typography
- **Professional Legend** - Modern legend with background, better organization, and clear visual hierarchy

## 📊 **Analysis**
- **Critical Path Insights** - Highlights priority and overdue tasks
- **Progress Metrics** - Completion rates, spans, status breakdowns
- **Dependency Tracking** - Maps task relationships and bottlenecks

## 🚀 Quick Start

### Basic Usage
```bash
# Generate LaTeX from CSV with timestamp
make build

# The command will also run pdflatex and produce:
# output/pdf/Calendar_YYYYMMDD_HHMMSS.pdf

# Enhanced build system with multiple templates
make build-all

# Build for specific device (e.g., e-ink tablet)
make build-device

# List available configurations
make list
```

### Advanced Usage
```bash
# Generate with specific template and device
python main.py build single ../input/data.cleaned.csv -t monthly_calendar -d supernote_a5x

# Generate all templates
python main.py build all-templates ../input/data.cleaned.csv

# Generate for all devices
python main.py build all-devices ../input/data.cleaned.csv

# List available configurations
python main.py build list

# Use main application directly
python main.py --template monthly_calendar --device supernote_a5x
```

## 📁 Project Structure

```
/Users/aaron/Downloads/gantt/
├── src/                          # Source code package
│   ├── app.py                   # Main application
│   ├── build.py                 # Enhanced build system
│   ├── config_manager.py        # Configuration management
│   ├── template_generators.py   # Template generation system
│   ├── data_processor.py        # CSV processing
│   ├── latex_generator.py       # LaTeX generation
│   ├── models.py                # Data models
│   ├── utils.py                 # Shared utilities
│   └── config/                  # Configuration files
│       ├── templates.yaml       # Template definitions
│       └── device_profiles.yaml # Device profiles
├── ../input/                    # Input data (moved to root)
│   ├── data.csv                 # Source CSV
│   └── data.cleaned.csv         # Cleaned CSV
├── output/                      # Generated files
│   ├── pdf/                     # PDF outputs
│   └── tex/                     # LaTeX sources
├── main.py                      # Single entry point (app + build)
├── Makefile                     # Build automation
└── README.md                    # This file
```

## 📦 Output
- `output/Timeline_YYYYMMDD_HHMMSS.pdf` — Compiled PDF (timestamped)

## 🧠 Rendering

LaTeX + TikZ for vector graphics and `longtable` for multi-page tables. No HTML/JS.

Time Axis
- **Month Boundaries** - Clear month starts with enhanced ticks at top and bottom, bold month labels
- **Weekly Grid** - Light Monday grid lines with adaptive ISO week number labeling (W1-W53)
- **Quarter Bands** - Subtle alternating shading above chart with Q1-Q4 labels
- **Enhanced Today Marker** - Prominent red line with white "TODAY" label and background
- **Better Contrast** - Improved line weights and colors for optimal readability

## ⏱️ Date Range
The LaTeX output summarizes the full task range found in the CSV and places sample markers along a 12‑unit timeline axis for context.

## 🎯 Use Cases

### **PhD Research Planning**
- Track proposal milestones and deadlines
- Monitor experiment progress
- Prepare for committee meetings
- Export professional reports for advisors

### **Project Management**
- Visualize project timelines
- Identify critical paths
- Track dependencies
- Generate status reports

### **Team Collaboration**
- Share print-ready timelines
- Export single-file PDFs
- Focused views via time windows
- Professional presentations

## 🔧 Technical Details

### **Input Format**
CSV with columns: `Task Name`, `Start Date`, `Due Date`, `Duration (days)`, `Group`, `Deliverable Type`, `Owner`, `Status`, `Priority`, `Notes`, `Dependencies`, `Parent Index`

### **Dependencies**
- Python 3.8+
- LaTeX distribution with `pdflatex` (TeX Live, MiKTeX, MacTeX)

### **Performance**
- Handles 1000+ tasks efficiently
- Vector-style bars and clean typography
- Optimized for large timelines

## 🎨 Customization

### **Colors & Styling**
- Professional color scheme
- Lane-based color coding
- Hierarchical visual design
- Print-optimized CSS

### **Layout Options**
- Configurable dimensions
- Flexible lane ordering
- Custom title and branding
- Responsive design

## 📈 What's New

### **v3.0 - Simplified PDF Edition**
- PDF-only output
- Cleaner CLI and docs
- Same great visuals and insights

## 🤝 Contributing

This tool is designed for academic and professional use. Feel free to:
- Report bugs or issues
- Suggest new features
- Share your use cases
- Improve the documentation

## 📄 License

Open source for academic and professional use.

---

**Perfect for:** PhD students, researchers, project managers, team leads, and anyone who needs professional project visualization.
