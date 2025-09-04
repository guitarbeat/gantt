# 🎯 LaTeX Project Timeline Generator

A comprehensive LaTeX-first tool that transforms CSV data into publication-quality timelines, calendars, and Gantt charts. Perfect for PhD research, formal reports, advisor meetings, and professional project management.

## ✨ Features

### 🎨 **Output Quality**
- **Publication Quality** - Professional typography and page layout using LaTeX
- **Multiple Templates** - Gantt timeline, monthly calendar, weekly planner
- **Device Optimization** - E-ink tablets, print formats, digital viewing
- **Professional Styling** - Enhanced TikZ graphics with shadows, rounded corners, and modern design

### 🔧 **Functionality**
- **CSV → LaTeX** - Convert CSV into complete .tex documents
- **Automatic Categorization** - Derives categories from `Group` and `Deliverable Type`
- **Status & Priority** - Colorized status and clear priority cues
- **No Truncation** - Full text preserved in tables and labels
- **Enhanced Build System** - Multiple build modes and automated compilation

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

### 🎨 **Advanced TikZ Features**
- **13 TikZ Libraries** automatically loaded for enhanced graphics:
  - `arrows.meta`, `shapes.geometric`, `positioning`, `calc`
  - `decorations.pathmorphing`, `patterns`, `shadows`, `fit`
  - `backgrounds`, `matrix`, `chains`, `scopes`, `pgfgantt`
- **Professional TikZ Styles** for task nodes, milestones, timeline axes, and calendar styling
- **Enhanced Generators** with modern timeline views, professional Gantt charts, and enhanced calendar grids

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
├── aarons-attempt/                    # Main Python application
│   ├── src/                          # Source code package
│   │   ├── __init__.py              # Package initialization
│   │   ├── app.py                   # Main application
│   │   ├── build.py                 # Enhanced build system
│   │   ├── config.py                # Core configuration
│   │   ├── config_manager.py        # YAML-based configuration
│   │   ├── data_processor.py        # CSV processing
│   │   ├── export_system.py         # Multi-format export
│   │   ├── interactive_generator.py # Interactive features
│   │   ├── latex_generator.py       # LaTeX generation
│   │   ├── models.py                # Data models
│   │   ├── template_generators.py   # Template generation
│   │   ├── utils.py                 # Shared utilities
│   │   └── config/                  # YAML configuration files
│   │       ├── templates.yaml       # Template definitions
│   │       └── device_profiles.yaml # Device profiles
│   ├── output/                      # Generated files
│   │   ├── pdf/                     # PDF outputs
│   │   └── tex/                     # LaTeX sources
│   ├── main.py                      # Single entry point (app + build)
│   ├── Makefile                     # Build automation
│   └── README.md                    # This file
├── latex-yearly-planner/             # Go-based LaTeX planner (reference)
│   ├── cmd/                         # Application entry points
│   ├── internal/                    # Private application code
│   ├── pkg/                         # Reusable components
│   ├── configs/                     # Configuration files
│   ├── templates/                   # LaTeX templates
│   └── scripts/                     # Build scripts
└── input/                           # Shared input data
    ├── data.csv                     # Source CSV
    └── data.cleaned.csv             # Cleaned CSV
```

## 📋 Available Templates

### 1. **Gantt Timeline** (Default)
- **Description**: Vertical timeline with task bars and dependencies
- **Layout**: Portrait orientation, A4 paper
- **Features**: Task bars, dependencies, milestones, status indicators
- **Best For**: Project planning, research timelines, formal reports

### 2. **Monthly Calendar**
- **Description**: Monthly grid view with task overlays
- **Layout**: Landscape orientation, A4 paper
- **Features**: Monthly grid, task overlays, milestone markers, week numbers
- **Best For**: Monthly planning, calendar integration, overview visualization

### 3. **Weekly Planner**
- **Description**: Weekly view with detailed task scheduling
- **Layout**: Landscape orientation, A4 paper
- **Features**: Weekly grid, time slots, task scheduling, notes sections
- **Best For**: Weekly planning, detailed scheduling, time management

## 📱 Device Profiles

### E-ink Devices
- **Supernote A5X**: Optimized for Supernote A5X e-ink tablet
- **ReMarkable 2**: Optimized for ReMarkable 2 e-ink tablet
- **Boox Max Lumi**: Optimized for Boox Max Lumi large e-ink tablet

### Print Devices
- **Standard Print**: Optimized for standard office/home printing
- **Professional Print**: Optimized for professional printing and binding
- **Large Format Print**: Optimized for large format printing (A3, A2)

### Digital Devices
- **Desktop PDF**: Optimized for desktop PDF viewing and annotation
- **Tablet PDF**: Optimized for tablet PDF viewing and annotation
- **Mobile PDF**: Optimized for mobile PDF viewing

## 🎨 Color Schemes

### Academic
- **Primary**: Blue (#3B82F6)
- **Secondary**: Green (#10B981)
- **Accent**: Orange (#F59E0B)
- **Best For**: Research, academic papers, formal reports

### Corporate
- **Primary**: Dark Gray (#1F2937)
- **Secondary**: Medium Gray (#4B5563)
- **Accent**: Blue (#3B82F6)
- **Best For**: Business presentations, corporate reports

### Vibrant
- **Primary**: Purple (#9333EA)
- **Secondary**: Pink (#EC4899)
- **Accent**: Green (#10B981)
- **Best For**: Creative projects, presentations, visual reports

## 🔧 Configuration

### Template Configuration (`src/config/templates.yaml`)
```yaml
templates:
  gantt_timeline:
    name: "Gantt Timeline"
    description: "Vertical timeline with task bars and dependencies"
    layout: "vertical"
    orientation: "portrait"
    page_size: "a4paper"
    margin: "0.5in"
    features:
      - "task_bars"
      - "dependencies"
      - "milestones"
```

### Device Profile Configuration (`src/config/device_profiles.yaml`)
```yaml
profiles:
  supernote_a5x:
    name: "Supernote A5X"
    description: "Optimized for Supernote A5X e-ink tablet"
    device_type: "eink"
    optimizations:
      - "high_contrast_colors"
      - "thick_lines"
      - "large_fonts"
    layout:
      page_size: "a5paper"
      orientation: "portrait"
      margin: "0.3in"
```

## 📊 Command Line Options

### Main Application (`main.py`)
```bash
python main.py [OPTIONS]

Options:
  --input, -i FILE          Input CSV file
  --output, -o FILE         Output LaTeX file
  --title, -t TITLE         Document title
  --template TEMPLATE       Template type (gantt_timeline, monthly_calendar, weekly_planner)
  --device, -d DEVICE       Device profile
  --color-scheme, -c SCHEME Color scheme
  --list-templates          List available templates
  --list-devices            List available device profiles
  --list-color-schemes      List available color schemes
  --verbose, -v             Enable verbose logging
  --quiet, -q               Suppress all output except errors
```

### Build System (`main.py build`)
```bash
python main.py build COMMAND [OPTIONS]

Commands:
  single INPUT              Build single document
  multiple INPUT            Build multiple documents
  all-templates INPUT       Build all templates
  all-devices INPUT         Build for all devices
  clean                     Clean build artifacts
  list                      List available configurations

Options:
  -t, --template TEMPLATE   Template type to use
  -d, --device DEVICE       Device profile to use
  -c, --color-scheme SCHEME Color scheme to use
  --title TITLE             Document title
  -o, --output OUTPUT       Output filename
  --verbose, -v             Enable verbose logging
  --quiet, -q               Suppress all output except errors
```

## 🎯 Use Cases

### Academic Research
```bash
# Generate academic timeline for PhD proposal
python main.py build single ../input/data.csv -t gantt_timeline -c academic --title "PhD Research Timeline"

# Generate monthly calendar for research planning
python main.py build single ../input/data.csv -t monthly_calendar -c academic --title "Research Calendar 2025"
```

### E-ink Device Usage
```bash
# Generate for Supernote A5X
python main.py build single ../input/data.csv -d supernote_a5x --title "Project Planner"

# Generate for ReMarkable 2
python main.py build single ../input/data.csv -d remarkable_2 --title "Weekly Planner"
```

### Professional Presentations
```bash
# Generate corporate-style timeline
python main.py build single ../input/data.csv -t gantt_timeline -c corporate --title "Project Timeline"

# Generate for large format printing
python main.py build single ../input/data.csv -d large_format_print --title "Project Overview"
```

### Batch Generation
```bash
# Generate all templates for review
python main.py build all-templates ../input/data.csv --title "Project Analysis"

# Generate for all devices
python main.py build all-devices ../input/data.csv -t gantt_timeline --title "Multi-Device Planner"
```

## 🧠 Rendering

LaTeX + TikZ for vector graphics and `longtable` for multi-page tables. No HTML/JS.

### Time Axis Features
- **Month Boundaries** - Clear month starts with enhanced ticks at top and bottom, bold month labels
- **Weekly Grid** - Light Monday grid lines with adaptive ISO week number labeling (W1-W53)
- **Quarter Bands** - Subtle alternating shading above chart with Q1-Q4 labels
- **Enhanced Today Marker** - Prominent red line with white "TODAY" label and background
- **Better Contrast** - Improved line weights and colors for optimal readability

## ⏱️ Date Range
The LaTeX output summarizes the full task range found in the CSV and places sample markers along a 12‑unit timeline axis for context.

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

## 🎨 Enhanced Visual Design

### **Modern TikZ Styling**
- **Enhanced Task Nodes** - Rounded corners, drop shadows, and modern typography
- **Status-Based Colors** - Green (completed), orange (in progress), red (blocked)
- **Progress Indicators** - Visual progress bars within task elements
- **Interactive Elements** - Clickable links and navigation
- **Professional Shadows** - Subtle drop shadows for depth and modern appearance

### **Advanced Color Schemes**
- **Academic** - Professional blue/green/orange palette for research
- **Corporate** - Dark gray/blue scheme for business presentations  
- **Vibrant** - Purple/pink/green for creative projects
- **Custom Colors** - Easily configurable color palettes

### **Enhanced Typography**
- **Modern Fonts** - Helvetica-based sans-serif for clean readability
- **Hierarchical Text** - Different font weights and sizes for information hierarchy
- **Better Spacing** - Improved line heights and margins for readability

## 🚀 New Export Formats

### **Multiple Output Formats**
- **PDF** - High-quality vector graphics for printing and sharing
- **SVG** - Scalable vector graphics for web and presentations
- **HTML** - Interactive web-based timelines with clickable elements
- **PNG** - High-resolution raster images for digital use

### **Export Commands**
```bash
# Export to multiple formats
make build-formats

# Export specific formats
python main.py build multiple-formats input/data.csv --formats pdf html

# Export with custom settings
python main.py build multiple-formats input/data.csv \
    --template gantt_timeline \
    --device supernote_a5x \
    --color-scheme academic \
    --formats pdf svg html png
```

## 🔧 Enhanced Data Processing

### **Improved CSV Validation**
- **Automatic Delimiter Detection** - Supports comma, semicolon, tab, and pipe delimiters
- **Multiple Date Formats** - Handles various date formats automatically
- **Enhanced Error Reporting** - Detailed validation messages with row numbers
- **Data Type Validation** - Ensures data integrity and consistency

### **Better Error Handling**
- **Graceful Degradation** - Continues processing even with some invalid rows
- **Detailed Logging** - Comprehensive error messages and warnings
- **Validation Reports** - Summary of data quality issues
- **Recovery Options** - Suggestions for fixing common data problems

## 🎯 Interactive Features

### **Enhanced Navigation**
- **Table of Contents** - Clickable navigation with hyperlinks
- **Progress Summary** - Visual progress indicators and statistics
- **Task Quick Links** - Jump to specific tasks in the document
- **Interactive Elements** - Hover effects and clickable components

### **Modern UI Elements**
- **Progress Bars** - Visual representation of project completion
- **Status Indicators** - Color-coded task status with icons
- **Interactive Legends** - Clickable legend items for filtering
- **Responsive Design** - Adapts to different screen sizes and devices

## 📈 What's New

### **v3.0 - Enhanced TikZ Edition**
- **Enhanced TikZ Libraries** - 13 powerful libraries for professional graphics
- **Modern Styling** - Shadows, rounded corners, and improved typography
- **Multiple Templates** - Gantt timeline, monthly calendar, weekly planner
- **Device Optimization** - E-ink tablets, print formats, digital viewing
- **Simplified Architecture** - Enhanced features are now default, no special configuration needed

## 🔄 Migration from Original

The enhanced system is fully backward compatible with the original codebase:

1. **Existing Commands**: All original commands continue to work
2. **Configuration**: Original configuration is preserved and enhanced
3. **Output Format**: Same LaTeX and PDF output format
4. **Data Format**: Same CSV input format

### Migration Steps
1. **No Changes Required**: Existing workflows continue to work
2. **Optional Enhancement**: Use new features as needed
3. **Gradual Adoption**: Migrate to new features over time

## 🛠️ Development

### Adding New Templates
1. Create template class in `src/template_generators.py`
2. Add template configuration in `src/config/templates.yaml`
3. Register template in `TemplateGeneratorFactory`

### Adding New Device Profiles
1. Add device profile in `src/config/device_profiles.yaml`
2. Configure optimizations and layout settings
3. Test with different templates

### Adding New Color Schemes
1. Add color scheme in `src/config/templates.yaml`
2. Define color palette and usage guidelines
3. Test with different templates and devices

## 📈 Performance

- **Efficient Processing**: Optimized for large datasets (1000+ tasks)
- **Memory Management**: Streamlined processing pipeline
- **Parallel Building**: Support for batch operations
- **Caching**: Configuration caching for faster subsequent builds

## 🔍 Troubleshooting

### Common Issues
1. **LaTeX Not Found**: Install LaTeX distribution (TeX Live, MiKTeX, MacTeX)
2. **Configuration Errors**: Check YAML syntax in config files
3. **Template Not Found**: Verify template name in available templates list
4. **Device Profile Issues**: Check device profile configuration

### Debug Mode
```bash
# Enable verbose logging
python main.py --verbose

# Enable debug logging in build system
python main.py build single ../input/data.csv --verbose
```

## 📚 References

- **Original Template**: [latex-yearly-planner](https://github.com/kudrykv/latex-yearly-planner)
- **LaTeX Documentation**: [LaTeX Project](https://www.latex-project.org/)
- **TikZ Documentation**: [TikZ & PGF](https://tikz.dev/)
- **Awesome TikZ Repository**: [awesome-tikz](https://github.com/xiaohanyu/awesome-tikz)

## 🤝 Contributing

Contributions are welcome! Areas for improvement:
- New template types
- Additional device profiles
- Enhanced color schemes
- Performance optimizations
- Documentation improvements

## 📄 License

Open source for academic and professional use.

---

**Perfect for**: PhD students, researchers, project managers, team leads, and anyone who needs professional project visualization with device-specific optimization.