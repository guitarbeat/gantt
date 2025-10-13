# Configuration Examples

This directory contains example configuration files demonstrating different use cases and customization levels for the PhD Dissertation Planner.

## Available Examples

### 📚 [base.yaml](base.yaml)
**Default comprehensive configuration**
- Complete configuration with all available options
- Production-ready settings optimized for academic planning
- Includes advanced typography, task styling, and layout parameters
- Recommended starting point for most users

### 🎓 [academic-research.yaml](academic-research.yaml)
**Research-focused configuration**
- Optimized for PhD dissertation and research project planning
- Shows research objectives in task descriptions
- Academic typography and spacing
- Ideal for research timeline visualization

### 📊 [compact-presentation.yaml](compact-presentation.yaml)
**Presentation-optimized configuration**
- High information density for meetings and presentations
- Smaller fonts and tighter spacing
- Landscape orientation (A4)
- Minimal visual elements for clean presentation slides

### ⚡ [minimal.yaml](minimal.yaml)
**Bare-bones configuration**
- Essential settings only for quick testing
- Minimal configuration for development and debugging
- Fast generation with basic functionality
- Good starting point for customization

### 🔧 [advanced-customization.yaml](advanced-customization.yaml)
**Complete feature showcase**
- Demonstrates every available configuration option
- Includes multiple pages and advanced features
- Reference for understanding all capabilities
- Not recommended for production use

## Usage

Copy any example to your project root and modify as needed:

```bash
# Copy an example configuration
cp docs/examples/academic-research.yaml config.yaml

# Use with the planner
go run ./cmd/planner --config config.yaml
```

## Configuration Structure

All configurations follow this structure:

```yaml
# Basic settings
year: 2025
weekstart: 1  # Monday
ampmtime: false
dotted: true

# Debug options
debug:
  showframe: false
  showlinks: false

# Layout and styling
layout:
  paper: {...}
  latex: {...}
  task_styling: {...}

# Pages to generate
pages:
  - name: monthly
    renderblocks: [...]

# Data source
csv_file_path: input_data/your_data.csv
output_dir: generated
```

## Key Configuration Areas

### Basic Settings
- `year`: Calendar year to generate
- `weekstart`: First day of week (0=Sunday, 1=Monday)
- `ampmtime`: 12-hour vs 24-hour time format

### Layout Configuration
- `paper`: Page dimensions and margins
- `latex`: LaTeX rendering parameters
- `task_styling`: Task appearance and spacing

### Advanced Features
- `debug`: Enable debugging visualizations
- `constraints`: Algorithm tuning parameters
- `typography`: Font and text settings

## Customization Tips

1. **Start with minimal.yaml** for basic functionality
2. **Use base.yaml** as a comprehensive template
3. **Customize academic-research.yaml** for dissertation planning
4. **Reference advanced-customization.yaml** to understand all options

## Data Format

All configurations expect CSV data with these columns:
- Task Name
- Start Date (YYYY-MM-DD)
- End Date (YYYY-MM-DD)
- Phase (1, 2, 3, 4)
- Category
- Status (Planned, In Progress, Completed)
- Assignee (optional)
- Description (optional)

See `input_data/research_timeline_v5.1_comprehensive.csv` for an example.
