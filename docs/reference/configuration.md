# ⚙️ Configuration Reference

Complete guide to configuring the PhD Dissertation Planner.

## 📖 Quick Links

- **[User Guide](../user/user-guide.md)** - How to use the planner
- **[Developer Guide](../developer/developer-guide.md)** - Development setup
- **[API Reference](api-reference.md)** - Technical API documentation
- **[Architecture](architecture.md)** - System design patterns
- **[Main README](../../README.md)** - Project overview

## Configuration Overview

The application uses YAML-based configuration with sensible defaults. Configuration can be overridden via:

1. **YAML files** - Custom config files
2. **Environment variables** - Runtime overrides
3. **Command line flags** - Direct overrides

## Configuration File Structure

```yaml
# Main configuration file (src/core/base.yaml)
debug:
  show_frame: false
  show_links: false

year: 2025
start_year: 2025
end_year: 2025

layout_engine:
  # Task positioning and visual settings
  initial_y_position_multiplier: 0.1
  task_height_multiplier: 0.6
  max_task_width_days: 7.0

  # Visual weight calculations
  duration_long_multiplier: 1.2
  duration_short_multiplier: 0.8
  milestone_weight_multiplier: 1.5
  category_weight_multiplier: 1.0

  # Urgency-based prominence
  urgency_multipliers:
    critical: 1.0
    high: 0.8
    medium: 0.6
    low: 0.4
    minimal: 0.2
    default: 0.5

  # Quality thresholds
  visual_quality_threshold: 0.3
  positioning_quality_threshold: 0.5

layout_engine.calendar_layout:
  # Day cell dimensions
  day_number_width: "6mm"
  day_content_margin: "8mm"
  task_cell_margin: "1mm"
  task_cell_spacing: "0.5mm"
  day_cell_minipage_width: "8mm"
  header_angle_size_offset: "2pt"

layout.latex:
  # Typography settings
  spacing:
    two_col: "5pt"
    tri_col: "5pt"
    five_col: "5pt"
    task_content_vspace: "0.2ex"
    task_overlay_arc: "2pt"

  typography:
    hyphen_penalty: 50
    tolerance: 1000
    emergency_stretch: "2em"
    sloppy_emergency_stretch: "2em"

# Task category colors (auto-generated)
layout.algorithmic_colors:
  proposal: "RGB(173,216,230)"
  laser: "RGB(144,238,144)"
  imaging: "RGB(255,182,193)"
  admin: "RGB(255,218,185)"
  dissertation: "RGB(221,160,221)"
  research: "RGB(176,224,230)"
  publication: "RGB(255,228,196)"
```

## Layout Engine Configuration

### Task Positioning (`layout_engine`)

Controls how tasks are positioned and sized in the calendar.

```yaml
layout_engine:
  # Starting position (percentage of day height from top)
  initial_y_position_multiplier: 0.1  # 10% from top

  # Task height (percentage of day height)
  task_height_multiplier: 0.6         # 60% of day height

  # Maximum task width in days
  max_task_width_days: 7.0            # Tasks won't span more than 1 week
```

### Visual Weight Calculation

Tasks are positioned based on calculated "visual weight":

```yaml
layout_engine:
  # Duration-based multipliers
  duration_long_multiplier: 1.2   # Tasks > 7 days get higher priority
  duration_short_multiplier: 0.8  # Tasks < 1 day get lower priority
  milestone_weight_multiplier: 1.5 # Milestone tasks get highest priority
  category_weight_multiplier: 1.0  # Base weight for categories
```

### Urgency Multipliers

Tasks can have different urgency levels affecting their prominence:

```yaml
urgency_multipliers:
  critical: 1.0  # Maximum prominence
  high: 0.8
  medium: 0.6
  low: 0.4
  minimal: 0.2
  default: 0.5  # Fallback value
```

### Quality Thresholds

Control when to apply different positioning strategies:

```yaml
# Thresholds for algorithm selection
visual_quality_threshold: 0.3      # Switch to quality positioning
positioning_quality_threshold: 0.5 # Minimum acceptable quality
```

## Calendar Layout Configuration

### Day Cell Dimensions

```yaml
calendar_layout:
  # Width settings
  day_number_width: "6mm"         # Day number column width
  day_cell_minipage_width: "8mm"  # Content area width

  # Spacing settings
  day_content_margin: "8mm"       # Margin around day content
  task_cell_margin: "1mm"         # Margin around individual tasks
  task_cell_spacing: "0.5mm"      # Spacing between stacked tasks

  # Header styling
  header_angle_size_offset: "2pt" # Angular header positioning
```

## LaTeX Typography Configuration

### Spacing Settings

```yaml
layout.latex.spacing:
  # Column spacing
  two_col: "5pt"    # 2-column layouts
  tri_col: "5pt"    # 3-column layouts
  five_col: "5pt"   # 5-column layouts

  # Task formatting
  task_content_vspace: "0.2ex"  # Vertical spacing in tasks
  task_overlay_arc: "2pt"       # Corner radius for task overlays
```

### Typography Controls

```yaml
layout.latex.typography:
  # Line breaking controls
  hyphen_penalty: 50      # Discourage hyphenation (0-1000)
  tolerance: 1000         # Allow loose line breaking (0-10000)

  # Stretch settings for tight layouts
  emergency_stretch: "2em"          # Maximum stretch for line breaking
  sloppy_emergency_stretch: "2em"   # Emergency stretch for problem lines
```

## Color Configuration

### Algorithmic Colors

Colors are automatically generated for task categories using a golden angle algorithm:

```yaml
layout.algorithmic_colors:
  proposal: "RGB(173,216,230)"    # Light blue
  laser: "RGB(144,238,144)"       # Light green
  imaging: "RGB(255,182,193)"     # Light pink
  admin: "RGB(255,218,185)"       # Peach
  dissertation: "RGB(221,160,221)" # Plum
  research: "RGB(176,224,230)"    # Powder blue
  publication: "RGB(255,228,196)" # Light orange
```

### Custom Colors

Override default colors by modifying the configuration:

```yaml
layout.algorithmic_colors:
  custom_category: "RGB(255,0,0)"  # Bright red
```

## Presets

The application includes three built-in presets:

### Academic Preset (`presets/academic.yaml`)

Optimized for detailed academic planning:
- Standard task heights and spacing
- Full typography settings
- Balanced color scheme

### Compact Preset (`presets/compact.yaml`)

Dense layout for space efficiency:
- Reduced task heights
- Tighter spacing
- Smaller fonts

### Presentation Preset (`presets/presentation.yaml`)

Designed for meetings and presentations:
- Larger fonts and spacing
- High contrast colors
- Clear visual hierarchy

## Environment Variables

Override configuration at runtime:

```bash
# Set year
export PLANNER_YEAR=2025

# Enable debug mode
export PLANNER_DEBUG=true

# Custom CSV file
export PLANNER_CSV_FILE="data/my_tasks.csv"

# Custom config file
export PLANNER_CONFIG_FILE="config/custom.yaml"
```

## Command Line Overrides

Direct overrides for common settings:

```bash
# Basic usage
go run ./cmd/planner --year 2025 --preset compact

# Custom files
go run ./cmd/planner --config custom.yaml --csv-file data/tasks.csv

# Debug mode
go run ./cmd/planner --debug --verbose
```

## Validation

Configuration is validated at startup:

- **Numeric ranges**: Values checked against valid ranges
- **Required fields**: Essential settings must be present
- **File paths**: Referenced files must exist (when specified)
- **YAML syntax**: Configuration files must be valid YAML

Invalid configurations will show detailed error messages with suggestions for fixes.

## Extending Configuration

### Adding Custom Presets

Create new preset files in `src/core/presets/`:

```yaml
# custom_preset.yaml
layout_engine:
  task_height_multiplier: 0.8
  initial_y_position_multiplier: 0.05

layout.latex.typography:
  hyphen_penalty: 25
```

Use with:
```bash
go run ./cmd/planner --preset custom_preset
```

### Custom Category Colors

Add new categories to the algorithmic colors:

```yaml
layout.algorithmic_colors:
  teaching: "RGB(255,255,153)"    # Light yellow
  collaboration: "RGB(204,204,255)" # Light purple
```

---

*Configuration reference last updated: October 2025*
