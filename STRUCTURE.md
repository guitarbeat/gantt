# Project Structure

Clean, consolidated structure for the PhD Dissertation Planner.

## Overview

```
phd-dissertation-planner/
├── input_data/              # All inputs
├── output_data/             # All outputs  
├── internal/                # All code
├── main.go                  # Entry point
├── plannergen               # Binary
└── Makefile                 # Build automation
```

## Detailed Structure

### Input (input_data/)

All input files in one place:

```
input_data/
├── config.yaml                      # Configuration
├── proposal_and_setup.csv           # Phase 1 tasks
├── research_and_experiments.csv     # Phase 2 tasks
├── publications.csv                 # Phase 3 tasks
└── dissertation_and_defense.csv     # Phase 4 tasks
```

### Output (output_data/)

All generated files organized by type:

```
output_data/
├── latex/          # LaTeX source files
├── pdfs/           # Compiled PDFs
├── auxiliary/      # LaTeX auxiliary files (.aux, .log)
└── binaries/       # Binary outputs
```

### Code (internal/)

All application code:

```
internal/
├── app/            # Application logic
│   ├── cli.go                  # Command-line interface
│   ├── generator.go            # Generation orchestration
│   └── template_funcs.go       # Template functions
│
├── calendar/       # Calendar generation
│   ├── calendar.go             # Calendar structures
│   ├── cell_builder.go         # Cell rendering
│   └── task_stacker.go         # Task layout
│
├── core/           # Core functionality
│   ├── colors.go               # Color generation
│   ├── config.go               # Configuration
│   ├── config_manager.go       # Config management
│   ├── defaults.go             # Default values
│   ├── errors.go               # Error handling
│   ├── logger.go               # Logging
│   ├── reader.go               # CSV reading
│   ├── task.go                 # Task structures
│   └── validation.go           # Data validation
│
└── templates/      # LaTeX templates
    ├── monthly/                # Monthly calendar templates
    │   ├── body.tpl
    │   ├── calendar.tpl
    │   ├── document.tpl
    │   ├── header.tpl
    │   ├── macros.tpl
    │   ├── page.tpl
    │   └── toc.tpl
    └── rendering.go            # Template rendering
```

### Root Files

```
├── main.go              # Application entry point
├── plannergen           # Compiled binary (gitignored)
├── Makefile             # Build automation
├── README.md            # Full documentation
├── QUICKSTART.md        # Quick start guide
├── STRUCTURE.md         # This file
├── go.mod               # Go module definition
├── go.sum               # Go dependencies
├── .gitignore           # Git ignore rules
└── .gitattributes       # Git attributes
```

## Design Principles

### 1. Clear Separation

- **Input** - Everything you edit
- **Output** - Everything generated
- **Code** - Everything that runs

### 2. Minimal Hierarchy

- No deep nesting
- Maximum 3 levels deep
- Easy to navigate

### 3. Logical Grouping

- Related files together
- Clear naming
- Obvious purpose

### 4. No Redundancy

- Single source of truth
- No duplicate configs
- No scattered files

## Consolidation History

### Before

```
├── cmd/planner/main.go
├── configs/config.yaml
├── generated/plannergen
├── input_data/*.csv
├── internal/...
└── pkg/templates/...
```

### After

```
├── main.go
├── input_data/config.yaml
├── input_data/*.csv
├── plannergen
└── internal/templates/...
```

### Changes Made

1. **Moved config** - `configs/` → `input_data/`
2. **Moved binary** - `generated/` → root
3. **Moved main** - `cmd/planner/` → root
4. **Moved templates** - `pkg/` → `internal/`
5. **Removed directories** - `cmd/`, `configs/`, `generated/`, `pkg/`
6. **Created output_data/** - Centralized output location

## Benefits

### For Users

- ✅ Clear input/output separation
- ✅ All inputs in one place
- ✅ All outputs in one place
- ✅ Simple to understand

### For Developers

- ✅ Flat structure
- ✅ Easy navigation
- ✅ Logical organization
- ✅ Minimal complexity

### For Maintenance

- ✅ Fewer directories
- ✅ Clear responsibilities
- ✅ Easy to find files
- ✅ Simple to extend

## File Counts

- **Input files**: 5 (4 CSV + 1 config)
- **Output directories**: 4 (organized by type)
- **Code packages**: 4 (app, calendar, core, templates)
- **Root files**: 10 (including docs and build files)

## Navigation Tips

### Finding Files

```bash
# All inputs
ls input_data/

# All outputs
ls output_data/

# All code
ls internal/

# Build files
ls *.go Makefile
```

### Common Paths

```bash
# Configuration
input_data/config.yaml

# Task data
input_data/*.csv

# Generated PDFs
output_data/pdfs/

# LaTeX source
output_data/latex/

# Main code
internal/app/generator.go

# Templates
internal/templates/monthly/
```

## Future Growth

The structure supports easy extension:

- **New input types** → Add to `input_data/`
- **New output formats** → Add to `output_data/`
- **New features** → Add to `internal/`
- **New templates** → Add to `internal/templates/`

No restructuring needed!
