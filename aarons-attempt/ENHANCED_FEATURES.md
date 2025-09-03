# 🚀 Enhanced LaTeX Gantt Chart Generator

## Overview

This enhanced version of the LaTeX Gantt Chart Generator incorporates the best features from the [latex-yearly-planner](https://github.com/kudrykv/latex-yearly-planner) template system, providing a comprehensive, flexible, and professional document generation system.

## ✨ New Features

### 🎯 **Template System**
- **Multiple Template Types**: Gantt timeline, monthly calendar, weekly planner
- **Modular Design**: Easy to add new templates and layouts
- **YAML Configuration**: Flexible configuration management
- **Template Inheritance**: Reusable components and layouts

### 📱 **Device Optimization**
- **E-ink Devices**: Optimized for Supernote A5X, ReMarkable 2, Boox Max Lumi
- **Print Optimization**: Standard, professional, and large format printing
- **Digital Viewing**: Desktop, tablet, and mobile PDF optimization
- **Auto-Detection**: Automatic device profile selection

### 🎨 **Advanced Customization**
- **Color Schemes**: Academic, corporate, and vibrant color palettes
- **Font Options**: Serif, sans-serif, and modern font configurations
- **Layout Control**: Flexible page layouts and orientations
- **Device-Specific Styling**: Optimized for different output devices

### 🔧 **Enhanced Build System**
- **Multiple Build Modes**: Single, multiple, all templates, all devices
- **Automated Compilation**: LaTeX compilation with error handling
- **Batch Processing**: Build multiple configurations at once
- **Clean Management**: Automated cleanup of build artifacts

## 📁 Project Structure

```
/Users/aaron/Downloads/gantt/
├── config/                          # Configuration files
│   ├── templates.yaml              # Template definitions
│   └── device_profiles.yaml        # Device-specific profiles
├── src/                            # Source code package
│   ├── config_manager.py           # Enhanced configuration management
│   ├── template_generators.py      # Template generation system
│   ├── app.py                      # Enhanced main application
│   └── [existing files...]
├── build.py                        # Enhanced build system
├── main.py                         # Primary entry point
└── [existing files...]
```

## 🚀 Quick Start

### Basic Usage

```bash
# Generate with default settings
python main.py

# Generate with specific template
python main.py --template monthly_calendar

# Generate for specific device
python main.py --device supernote_a5x

# Generate with specific color scheme
python main.py --color-scheme corporate
```

### Enhanced Build System

```bash
# Build single document
python build.py single input/data.csv

# Build with specific template and device
python build.py single input/data.csv -t monthly_calendar -d supernote_a5x

# Build all templates
python build.py all-templates input/data.csv

# Build for all devices
python build.py all-devices input/data.csv

# Build multiple configurations
python build.py multiple input/data.csv -t gantt_timeline monthly_calendar -d supernote_a5x remarkable_2

# Clean build artifacts
python build.py clean

# List available configurations
python build.py list
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

### Template Configuration (`config/templates.yaml`)
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

### Device Profile Configuration (`config/device_profiles.yaml`)
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

### Build System (`build.py`)
```bash
python build.py COMMAND [OPTIONS]

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
python build.py single input/data.csv -t gantt_timeline -c academic --title "PhD Research Timeline"

# Generate monthly calendar for research planning
python build.py single input/data.csv -t monthly_calendar -c academic --title "Research Calendar 2025"
```

### E-ink Device Usage
```bash
# Generate for Supernote A5X
python build.py single input/data.csv -d supernote_a5x --title "Project Planner"

# Generate for ReMarkable 2
python build.py single input/data.csv -d remarkable_2 --title "Weekly Planner"
```

### Professional Presentations
```bash
# Generate corporate-style timeline
python build.py single input/data.csv -t gantt_timeline -c corporate --title "Project Timeline"

# Generate for large format printing
python build.py single input/data.csv -d large_format_print --title "Project Overview"
```

### Batch Generation
```bash
# Generate all templates for review
python build.py all-templates input/data.csv --title "Project Analysis"

# Generate for all devices
python build.py all-devices input/data.csv -t gantt_timeline --title "Multi-Device Planner"
```

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
2. Add template configuration in `config/templates.yaml`
3. Register template in `TemplateGeneratorFactory`

### Adding New Device Profiles
1. Add device profile in `config/device_profiles.yaml`
2. Configure optimizations and layout settings
3. Test with different templates

### Adding New Color Schemes
1. Add color scheme in `config/templates.yaml`
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
python build.py single input/data.csv --verbose
```

## 📚 References

- **Original Template**: [latex-yearly-planner](https://github.com/kudrykv/latex-yearly-planner)
- **LaTeX Documentation**: [LaTeX Project](https://www.latex-project.org/)
- **TikZ Documentation**: [TikZ & PGF](https://tikz.dev/)

## 🤝 Contributing

Contributions are welcome! Areas for improvement:
- New template types
- Additional device profiles
- Enhanced color schemes
- Performance optimizations
- Documentation improvements

---

**Perfect for**: PhD students, researchers, project managers, team leads, and anyone who needs professional project visualization with device-specific optimization.
