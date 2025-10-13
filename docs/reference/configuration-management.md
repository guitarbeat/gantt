# Configuration Management

This document describes the centralized configuration management system for the PhD dissertation planner.

## Overview

The configuration system provides:
- **Centralized Management**: All configuration loading, validation, and management in one place
- **Environment Variable Consolidation**: Comprehensive environment variable handling with validation
- **Startup Validation**: Automatic validation at application startup with detailed error reporting
- **Hot Reloading**: Configuration file watching for development (when `DEV_TEMPLATES` is set)
- **Preset System**: Predefined configuration profiles for different use cases

## Configuration Files

### Base Configuration (`configs/base.yaml`)

The base configuration file contains all default settings. It includes:
- Layout settings (paper size, margins, colors)
- Typography settings (fonts, spacing)
- Task styling (appearance, spacing)
- Calendar layout parameters
- Stacking and rendering algorithms

### Specialized Configurations

- `configs/monthly_calendar.yaml`: Monthly calendar layout with task support
- `src/core/presets/academic.yaml`: Default academic planning preset
- `src/core/presets/compact.yaml`: Dense information display preset
- `src/core/presets/presentation.yaml`: Large text for presentations

## Environment Variables

All environment variables are centrally managed and validated:

### Core Variables
- `PLANNER_YEAR`: Academic year (YYYY format, 2000-2100)
- `PLANNER_START_YEAR`: Multi-year planning start year
- `PLANNER_END_YEAR`: Multi-year planning end year
- `PLANNER_CSV_FILE`: Path to CSV task data file
- `PLANNER_OUTPUT_DIR`: Output directory (default: "build")

### Layout Variables
- `PLANNER_LAYOUT_PAPER_WIDTH`: Paper width (with units: pt, mm, cm, in, em, ex)
- `PLANNER_LAYOUT_PAPER_HEIGHT`: Paper height
- `PLANNER_LAYOUT_PAPER_MARGIN_TOP`: Top margin
- `PLANNER_LAYOUT_PAPER_MARGIN_BOTTOM`: Bottom margin
- `PLANNER_LAYOUT_PAPER_MARGIN_LEFT`: Left margin
- `PLANNER_LAYOUT_PAPER_MARGIN_RIGHT`: Right margin

### Control Variables
- `PLANNER_PRESET`: Configuration preset (academic, compact, presentation)
- `PLANNER_SILENT`: Suppress log output (true/false)
- `PLANNER_LOG_LEVEL`: Logging level (silent, info, debug)
- `DEV_TEMPLATES`: Use filesystem templates for development (true/false)

## Presets

### Academic (Default)
Balanced configuration for academic planning:
- Standard font sizes and spacing
- Task descriptions enabled
- Balanced density and readability

### Compact
Dense information display:
- Smaller fonts and tighter spacing
- Reduced cell heights
- More tasks per page

### Presentation
Optimized for presentations:
- Larger fonts for readability
- Increased spacing
- Enhanced visual clarity

## Validation

### Startup Validation

The system performs comprehensive validation at startup:

1. **Required Fields**: Ensures all required configuration fields are present
2. **Structure Validation**: Validates YAML structure and value ranges
3. **Environment Variables**: Validates all environment variable values
4. **File Paths**: Checks file existence and permissions
5. **CSV Data**: Validates CSV file format and content (if specified)

### CLI Validation

Run configuration validation without generating output:

```bash
# Validate configuration files and environment
./plannergen --validate-config

# Validate with specific preset
./plannergen --validate-config --preset academic

# Validate with custom config files
./plannergen --validate-config --config custom.yaml
```

### Validation Output

Validation provides detailed error messages:

```
🔍 Configuration Validation
═══════════════════════════════════════
Configuration files: [src/core/base.yaml]
Using preset: academic

✅ All configuration validation passed!
═══════════════════════════════════════
```

## Hot Reloading

When `DEV_TEMPLATES=true` is set, the system automatically watches configuration files for changes and reloads them:

```bash
export DEV_TEMPLATES=true
./plannergen --config base.yaml
# Edit base.yaml in another terminal
# Configuration automatically reloads
```

Hot reload events are logged:

```
🔄 Configuration reloaded at 15:30:22
```

## Configuration Loading Order

Configuration is loaded in this order (later sources override earlier ones):

1. **Base Configuration Files**: YAML files specified with `--config`
2. **Preset Configuration**: Applied if `--preset` is specified
3. **Environment Variables**: Override configuration file values
4. **CLI Overrides**: Output directory from `--outdir` flag

## Examples

### Basic Usage
```bash
# Use default configuration
./plannergen

# Use specific config file
./plannergen --config configs/base.yaml

# Use multiple config files (later files override earlier)
./plannergen --config configs/base.yaml,configs/custom.yaml
```

### With Presets
```bash
# Academic preset (default detailed view)
./plannergen --preset academic

# Compact preset (dense display)
./plannergen --preset compact

# Presentation preset (large text)
./plannergen --preset presentation
```

### Environment Variables
```bash
# Set year and output directory
export PLANNER_YEAR=2024
export PLANNER_OUTPUT_DIR=./output
./plannergen

# Use custom CSV file
export PLANNER_CSV_FILE=./data/tasks.csv
./plannergen
```

### Development Mode
```bash
# Enable hot reloading and filesystem templates
export DEV_TEMPLATES=true
export PLANNER_LOG_LEVEL=debug
./plannergen --config src/core/base.yaml
```

## Error Handling

Configuration errors provide detailed context:

```
❌ Startup validation failed: configuration validation failed:
layout.paper.width is required
environment variable 'PLANNER_YEAR': year must be 4 digits
```

## Security

The system includes security measures:
- File paths cannot contain `..` for directory traversal protection
- Output directories are validated for write permissions
- Environment variable values are validated for proper formats

## Migration from Old System

The new configuration manager is backward compatible with existing configurations. The main changes:

- **Centralized Loading**: `core.NewConfig()` replaced with `core.NewConfigManager().Load()`
- **Validation**: Automatic validation at startup (can be disabled if needed)
- **Hot Reloading**: Automatic when `DEV_TEMPLATES=true`
- **Preset Support**: New `--preset` flag and `PLANNER_PRESET` environment variable

Existing scripts and configurations continue to work without changes.
