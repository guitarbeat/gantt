# 📚 PhD Dissertation Planner - Complete Reference Guide

## Table of Contents

1. [Project Overview](#-project-overview)
2. [Quick Start](#-quick-start)
3. [Project Status](#-project-status)
4. [Directory Structure & Organization](#-directory-structure--organization)
5. [Go Project Structure Guide](#-go-project-structure-guide)
6. [Lessons Learned from aarons-attempt](#-lessons-learned-from-aarons-attempt)
7. [Architecture Patterns](#-architecture-patterns)
8. [Implementation Strategies](#-implementation-strategies)
9. [Design Patterns](#-design-patterns)
10. [Key Features to Adopt](#-key-features-to-adopt)
11. [Migration Strategy](#-migration-strategy)
12. [Code Quality Lessons](#-code-quality-lessons)
13. [Success Metrics](#-success-metrics)
14. [Development](#-development)
15. [Troubleshooting](#-troubleshooting)

---

## 📋 Project Overview

Welcome to the comprehensive reference documentation for the PhD Dissertation Planner project. This document combines project overview, lessons learned, and directory organization plans into a complete guide for understanding, using, and contributing to the project.

### 🎯 Project Mission

The PhD Dissertation Planner is a Go-based application that transforms CSV data into professional LaTeX-generated PDF planners and Gantt charts for academic project management.

### 🔗 External Resources

- **Project Repository**: Main source code and issue tracking
- **LaTeX Documentation**: [LaTeX Project](https://www.latex-project.org/)
- **Go Documentation**: [Go Programming Language](https://golang.org/doc/)

---

## 🚀 Quick Start

```bash
# Setup development environment (downloads dependencies locally for offline use)
./scripts/setup.sh

# Build the generator and produce the PDF
make clean-build

# Or run the generator manually
go build -mod=vendor -o generated/plannergen ./cmd/planner && \
  PLANNER_SILENT=1 PLANNER_CSV_FILE="input_data/Research Timeline v5 - Comprehensive.csv" \
  ./generated/plannergen --config "src/core/base.yaml" --outdir generated
```

**Note**: Dependencies are vendored locally, so the project works offline after initial setup.

---

## ⚙️ Configuration

The application uses a comprehensive YAML-based configuration system that allows you to customize all aspects of the calendar generation without modifying code.

### Configuration File

The main configuration file is located at `src/core/base.yaml` and contains settings for:

- **Layout Engine**: Task positioning, visual weights, quality thresholds
- **Calendar Layout**: Day dimensions, spacing, typography
- **LaTeX Spacing**: Column widths, task spacing, visual elements
- **Stacking Engine**: Visual constraints and overflow handling

### Key Configuration Sections

#### Layout Engine (`layout_engine`)
Controls task positioning and visual appearance:

```yaml
layout_engine:
  # Task positioning multipliers (relative to day dimensions)
  initial_y_position_multiplier: 0.1  # 10% from top of day
  task_height_multiplier: 0.6         # 60% of day height
  max_task_width_days: 7.0            # Maximum task width in days
  
  # Visual weight calculation multipliers
  duration_long_multiplier: 1.2       # For tasks > 7 days
  duration_short_multiplier: 0.8      # For tasks < 1 day
  milestone_weight_multiplier: 1.5    # For milestone tasks
  category_weight_multiplier: 1.0     # Base category weight
  
  # Urgency multipliers for prominence calculation
  urgency_multipliers:
    critical: 1.0
    high: 0.8
    medium: 0.6
    low: 0.4
    minimal: 0.2
    default: 0.5
```

#### Calendar Layout (`layout_engine.calendar_layout`)
Controls day cell dimensions and spacing:

```yaml
calendar_layout:
  day_number_width: "6mm"           # Width of day number cells
  day_content_margin: "8mm"         # Margin around day content
  task_cell_margin: "1mm"           # Margin around task cells
  task_cell_spacing: "0.5mm"        # Spacing between task cells
  day_cell_minipage_width: "8mm"    # Width of day cell minipages
  header_angle_size_offset: "2pt"   # Offset for header angle calculations
```

#### LaTeX Spacing (`layout.latex.spacing`)
Controls LaTeX-specific spacing and visual elements:

```yaml
layout:
  latex:
    spacing:
      two_col: "5pt"                # Two-column spacing
      tri_col: "5pt"                # Three-column spacing
      five_col: "5pt"               # Five-column spacing
      task_content_vspace: "0.2ex"  # Vertical spacing for task content
      task_overlay_arc: "2pt"       # Arc radius for task overlays
```

### Default Values and Validation

The application provides sensible defaults for all configuration values and validates them at startup:

- **Multipliers**: Must be positive numbers
- **Thresholds**: Must be between 0 and 1
- **Dimensions**: Must be positive LaTeX units (pt, mm, em, ex)
- **Ranges**: Automatically clamped to valid ranges

### Overriding Configuration

You can override any configuration value by:

1. **Environment Variables**: Set `PLANNER_CONFIG_FILE` to point to a custom YAML file
2. **Command Line**: Use `--config` flag to specify a different configuration file
3. **YAML Override**: Create a custom YAML file that extends the base configuration

### Example Custom Configuration

```yaml
# custom_config.yaml
layout_engine:
  task_height_multiplier: 0.8  # Make tasks taller
  max_task_width_days: 10.0    # Allow wider tasks
  
layout:
  latex:
    spacing:
      task_content_vspace: "0.3ex"  # More vertical spacing
```

Run with custom config:
```bash
go run ./src/cmd/planner --config custom_config.yaml
```

---

## ✅ Project Status

[![CI](https://github.com/your-username/phd-dissertation-planner/actions/workflows/ci.yml/badge.svg)](https://github.com/your-username/phd-dissertation-planner/actions)

- ✅ **PDF Generation**: Working (generates ~116KB PDFs)
- ✅ **CSV Processing**: 84 tasks parsed successfully
- ✅ **LaTeX Compilation**: XeLaTeX integration working
- ✅ **Template System**: Go templates rendering correctly

---

## 📁 Current Project Structure

```
├── cmd/planner/           # Go application entry point
├── src/                   # Source code (beginner-friendly)
│   ├── app/              # Main application logic
│   ├── core/             # Core utilities and shared logic
│   ├── calendar/         # Calendar/scheduling functionality
│   └── shared/           # Shared/reusable code
│       └── templates/    # LaTeX templates (embedded)
├── input_data/           # Input data files (CSV, etc.)
├── generated/            # Generated output files (PDFs, logs)
├── static_assets/        # Static files/assets
├── vendor/               # Vendored Go dependencies (for offline builds)
├── scripts/              # Setup and build scripts
└── docs/                 # Documentation
```

---

## 🏗️ Go Project Structure Guide

This section explains the beginner-friendly directory structure used in this Go project.

### Core Directories

#### `src/` - Source Code

**Purpose**: Contains all the main application source code
**Why this name**: "src" is universally understood as "source code" - much clearer than "internal" for beginners
**Contents**:

- `src/app/` - Main application logic and CLI handling
- `src/core/` - Core utilities, configuration, and shared business logic
- `src/calendar/` - Calendar generation, task scheduling, and layout management

#### `shared/` - Shared/Reusable Code

**Purpose**: Contains code that can be reused across different parts of the application
**Why this name**: "shared" clearly indicates reusable code - more intuitive than "pkg"
**Contents**:

- `shared/templates/` - LaTeX templates and rendering utilities

#### `input_data/` - Input Data Files

**Purpose**: Contains all input data files (CSV, configuration, etc.)
**Why this name**: "input_data" explicitly shows this is where data goes in - clearer than just "data"
**Contents**:

- CSV files with task data
- Markdown files with project plans
- Other input files

#### `generated/` - Generated Output Files

**Purpose**: Contains all files generated by the application
**Why this name**: "generated" clearly indicates these are output files created by the program
**Contents**:

- PDF files
- LaTeX source files
- Log files
- Compiled binaries

#### `static_assets/` - Static Files/Assets

**Purpose**: Contains static files that don't change during execution
**Why this name**: "static_assets" is explicit about containing unchanging files
**Contents**:

- PDF documents
- Images
- Other static resources

### Supporting Directories

#### `cmd/` - Command Line Applications

**Purpose**: Contains the main entry points for executable programs
**Why this name**: Standard Go convention for application entry points
**Contents**:

- `cmd/planner/` - Main CLI application

#### `vendor/` - Dependencies

**Purpose**: Contains vendored (local copies) of external dependencies
**Why this name**: Standard Go convention for dependency management
**Contents**:

- Local copies of external Go packages

#### `scripts/` - Build and Setup Scripts

**Purpose**: Contains shell scripts for building, testing, and setup
**Why this name**: Self-explanatory - contains scripts
**Contents**:

- Build scripts
- Setup scripts
- Development utilities

#### `docs/` - Documentation

**Purpose**: Contains all project documentation
**Why this name**: Self-explanatory - contains docs
**Contents**:

- User guides
- Developer guides
- API documentation

### Why This Structure is Beginner-Friendly

#### 1. **Clear Naming**

- `src/` instead of `internal/` - everyone knows what "src" means
- `shared/` instead of `pkg/` - "shared" is more descriptive
- `input_data/` instead of `data/` - explicit about direction of data flow
- `generated/` instead of `output/` - clear that these are created files

#### 2. **Logical Grouping**

- Related functionality is grouped together
- Clear separation between input, processing, and output
- Easy to understand the data flow through the project

#### 3. **Standard Go Conventions**

- Still follows Go best practices where it makes sense
- Uses `cmd/` for executables (standard Go)
- Uses `vendor/` for dependencies (standard Go)
- Maintains proper package structure within `src/`

#### 4. **Self-Documenting**

- Directory names explain their purpose
- Easy to guess where to find specific types of files
- Reduces cognitive load for new developers

### Data Flow

``` text
input_data/ → src/core/ → src/app/ → generated/
     ↓              ↓         ↓
static_assets/ → shared/ → cmd/planner/
```

1. **Input**: Data starts in `input_data/` and `static_assets/`
2. **Processing**: Code in `src/` processes the data
3. **Templates**: `shared/templates/` provides rendering templates
4. **Output**: Results are written to `generated/`

### Getting Started

1. **Look at `src/app/`** - Start here to understand the main application
2. **Check `src/core/`** - Core utilities and configuration
3. **Explore `src/calendar/`** - Calendar generation logic
4. **Review `input_data/`** - See what data the app processes
5. **Check `generated/`** - See what the app produces

### Further Reading

- [Go Project Layout](https://github.com/golang-standards/project-layout) - Official Go project layout recommendations
- [Effective Go](https://golang.org/doc/effective_go.html) - Go best practices
- [Go Modules](https://golang.org/ref/mod) - Go dependency management

---

## 📁 Directory Structure & Organization

### Current State Analysis

- **Total Files:** 149 files
- **Go Files:** 124 files (83.2%)
- **Test Files:** 18 files scattered in root
- **Backup Files:** 4 files (to be removed)
- **Temp Files:** 2 files (to be removed)
- **Test Output:** 3 directories with generated files

### Proposed Clean Directory Structure

```
latex-yearly-planner/
├── cmd/                                    # Application entry points
│   └── plannergen/
│       └── main.go
├── internal/                               # Core application logic
│   ├── app/
│   │   └── app.go
│   ├── calendar/                           # Calendar functionality
│   │   ├── *.go (core files)
│   │   └── *_test.go (test files)
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── data/
│   │   ├── *.go (core files)
│   │   └── *_test.go (test files)
│   ├── generator/                          # PDF generation
│   │   ├── *.go (core files)
│   │   └── *_test.go (test files)
│   ├── header/
│   │   └── *.go
│   ├── latex/
│   │   └── *.go
│   └── layout/
│       └── lengths.go
├── configs/                                # Configuration files
│   ├── base.yaml
│   ├── csv_config.yaml
│   ├── page_template.yaml
│   ├── planner_config.yaml
│   └── README.md
├── templates/                              # Template files
│   ├── embed.go
│   ├── monthly/
│   │   └── *.tpl
│   └── README.md
├── scripts/                                # Build and utility scripts
│   ├── README.md
│   └── simple.sh
├── tests/                                  # Organized test files
│   ├── integration/
│   │   ├── test_feedback_integration.go
│   │   ├── test_performance_integration.go
│   │   └── test_visual_integration.go
│   ├── quality/
│   │   ├── test_quality_assurance.go
│   │   ├── test_quality_issue_resolver.go
│   │   ├── test_quality_system.go
│   │   └── test_quality_validator.go
│   ├── validation/
│   │   ├── test_user_coordination.go
│   │   ├── test_user_validation.go
│   │   ├── test_final_approval.go
│   │   └── test_final_assessment.go
│   ├── performance/
│   │   ├── test_performance_optimization.go
│   │   ├── test_performance_simple.go
│   │   └── test_visual_spacing.go
│   ├── feedback/
│   │   ├── test_feedback_system.go
│   │   └── test_improvement_logic.go
│   └── parser/
│       ├── test_enhanced_parser.go
│       ├── test_enhanced_visual.go
│       ├── test_multi_format.go
│       └── test_parser.go
├── docs/                                   # Consolidated documentation
│   ├── README.md                           # Main project documentation
│   ├── reports/
│   │   ├── TASK_3_3_COMPLETION_REPORT.md
│   │   ├── TASK_3_4_COMPLETION_REPORT.md
│   │   └── PERFORMANCE_OPTIMIZATION_REPORT.md
│   ├── lessons/
│   │   └── LESSONS_LEARNED_FROM_AARONS_ATTEMPT.md
│   └── api/                                # API documentation
│       └── README.md
├── examples/                               # Example configurations
│   ├── sample_batch_config.json
│   └── README.md
├── build/                                  # Build artifacts (gitignored)
│   └── .gitkeep
├── dist/                                   # Distribution files (gitignored)
│   └── .gitkeep
├── go.mod
├── go.sum
├── Makefile
└── .gitignore                              # Updated to exclude build artifacts
```

### File Organization Actions

#### 1. Files to Remove (Cleanup)

**Backup Files (4 files):**
- `internal/calendar/integration_test.go.bak`
- `internal/calendar/complex_overlap_test.go.bak`
- `internal/calendar/integration_test.go.temp.bak`
- `internal/calendar/complex_overlap_test.go.temp.bak`

**Temporary Files (2 files):**
- `internal/calendar/complex_overlap_test.go.temp`
- `internal/calendar/integration_test.go.temp`

**Temporary Test Files (5 files):**
- `simple_quality_test.go`
- `simple_quality_validation.go`
- `simple_test.go`
- `simple_validation.go`
- `simple_visual_test.go`
- `simple_visual_validation.go`
- `validate_integration.go`

**Test Output Directories (3 directories):**
- `test_output/` (remove contents, keep directory)
- `test_single_output/` (remove contents, keep directory)
- `test_triple_output/` (remove contents, keep directory)

#### 2. Files to Move (Reorganization)

**Test Files to Move to `tests/` directory:**
- Integration, Quality, Validation, Performance, Feedback, and Parser tests moved to organized subdirectories

**Documentation Files to Move to `docs/` directory:**
- Reports and lessons moved to appropriate subdirectories

**Example Files to Move to `examples/` directory:**
- Sample configurations moved to examples directory

#### 3. Implementation Benefits

1. **Clean Root Directory:** Only essential files in root
2. **Organized Tests:** Tests grouped by functionality
3. **Consolidated Documentation:** All docs in one place
4. **Clear Separation:** Core code, tests, docs, examples separated
5. **Professional Structure:** Industry-standard Go project layout
6. **Maintainable:** Easy to find and manage files
7. **Scalable:** Structure supports future growth

---

## 🎯 Lessons Learned from aarons-attempt

This section captures key architectural patterns, design decisions, and implementation strategies from the `aarons-attempt` project that can be applied to improve the `latex-yearly-planner` Go application.

### Project Background

The `aarons-attempt` project is a Python-based LaTeX timeline generator that transforms CSV data into publication-quality timelines and Gantt charts. It demonstrates several excellent patterns for building robust, maintainable document generation tools.

---

## 🏗️ Architecture Patterns

### 1. **Modular Package Structure**
```
src/
├── __init__.py              # Clean package exports
├── app.py                   # Main application logic
├── config.py                # Centralized configuration
├── generator.py             # LaTeX generation
├── models.py                # Data models
├── processor.py             # CSV processing
├── prisma_generator.py      # Specialized generators
└── utils.py                 # Utility functions
```

**Key Lessons:**
- Clear separation of concerns
- Single responsibility per module
- Clean import structure with `__all__` exports
- Logical grouping of related functionality

### 2. **Configuration Management**
```python
@dataclass
class ColorScheme:
    milestone: Tuple[int, int, int] = (147, 51, 234)
    researchcore: Tuple[int, int, int] = (59, 130, 246)
    # ... more colors

@dataclass
class TaskConfig:
    category_keywords: Dict[str, List[str]] = field(default_factory=lambda: {
        "researchcore": ["PROPOSAL"],
        "researchexp": ["LASER", "EXPERIMENTAL"],
        # ... more mappings
    })
```

**Key Lessons:**
- Type-safe configuration with dataclasses
- Hierarchical configuration structure
- Sensible defaults with field factories
- Environment variable support
- Category-to-color mapping system

### 3. **Data Model Design**
```python
@dataclass
class Task:
    id: str
    name: str
    start_date: date
    due_date: date
    category: str
    dependencies: str = ""
    notes: str = ""
    is_milestone: bool = False

    def __post_init__(self):
        self._validate_dates()
        self._determine_milestone()

    @property
    def duration_days(self) -> int:
        return (self.due_date - self.start_date).days + 1
```

**Key Lessons:**
- Immutable data structures with validation
- Computed properties for derived data
- Automatic milestone detection
- Input validation in `__post_init__`
- Clear, descriptive field names

---

## 🔧 Implementation Strategies

### 1. **Robust CSV Processing**
```python
def _parse_date(self, date_str: str) -> Optional[date]:
    formats = ['%Y-%m-%d', '%m/%d/%Y', '%d/%m/%Y', '%B %d, %Y']
    for fmt in formats:
        try:
            return datetime.strptime(date_str, fmt).date()
        except ValueError:
            continue
    return None
```

**Key Lessons:**
- Multiple date format support
- Graceful error handling
- Optional return types for failed parsing
- Comprehensive format coverage

### 2. **Professional LaTeX Generation**
```python
def _generate_color_definitions(self) -> str:
    colors = []
    for attr_name in dir(config.colors):
        if not attr_name.startswith('_') and not callable(getattr(config.colors, attr_name)):
            rgb = getattr(config.colors, attr_name)
            colors.append(f"\\definecolor{{{attr_name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
    return '\n'.join(colors)
```

**Key Lessons:**
- Dynamic color definition generation
- Professional typography with Helvetica
- Comprehensive LaTeX package usage
- Clean, readable LaTeX output
- Proper escaping of special characters

### 3. **Timeline View Generation**
```python
def _generate_timeline_view(self, timeline: ProjectTimeline) -> str:
    content = """
\\section*{Project Timeline}
\\begin{enumerate}[leftmargin=0pt, itemindent=0pt, labelsep=0pt, labelwidth=0pt]
"""
    for i, task in enumerate(timeline.tasks, 1):
        if task.is_milestone:
            milestone_indicator = "\\textcolor{red}{\\textbf{★}}"
            content += f"\\item[{milestone_indicator}] \\textbf{{{task_name}}}\n"
        else:
            content += f"\\item[{i:02d}] \\textbf{{{task_name}}}\n"
```

**Key Lessons:**
- Beautiful enumerated timeline layout
- Milestone indicators with special symbols
- Consistent formatting and spacing
- Professional typography choices

### 4. **Error Handling and Logging**
```python
def run(self, args: argparse.Namespace) -> int:
    try:
        self.logger.info("Starting LaTeX Gantt Chart Generator")
        # ... processing
        return 0
    except KeyboardInterrupt:
        self.logger.info("Application interrupted by user")
        return 130
    except Exception as e:
        self.logger.error(f"Unexpected error: {e}")
        return 1
```

**Key Lessons:**
- Comprehensive error handling
- Proper exit codes
- User-friendly error messages
- Graceful handling of interruptions

---

## 🎨 Design Patterns

### 1. **Builder Pattern for LaTeX Generation**
```python
def generate_document(self, timeline: ProjectTimeline, include_prisma: bool = False) -> str:
    content = self._generate_header()
    content += self._generate_title_page(timeline)
    content += self._generate_timeline_view(timeline)
    content += self._generate_task_list(timeline)

    if include_prisma:
        content += self._generate_prisma_section()

    content += self._generate_footer()
    return content
```

**Key Lessons:**
- Modular document generation
- Optional sections with flags
- Clean separation of concerns
- Easy to extend with new sections

### 2. **Factory Pattern for Data Processing**
```python
def process_csv_to_timeline(self, csv_file_path: str, title: str = None) -> ProjectTimeline:
    tasks = self._read_csv_data(csv_file_path)
    timeline_title = title or config.latex.default_title
    return ProjectTimeline(
        tasks=tasks,
        title=timeline_title,
        start_date=date.today(),
        end_date=date.today()
    )
```

**Key Lessons:**
- High-level factory methods
- Sensible defaults
- Clear data flow
- Easy to test and mock

### 3. **Strategy Pattern for Color Mapping**
```python
def get_category_color(category: str) -> str:
    category_upper = category.upper()
    return next(
        (color_name for color_name, keywords in config.tasks.category_keywords.items()
         if any(keyword in category_upper for keyword in keywords)),
        "other"
    )
```

**Key Lessons:**
- Flexible category mapping
- Fallback to default color
- Case-insensitive matching
- Easy to extend with new categories

---

## 🚀 Key Features to Adopt

### 1. **Comprehensive Color Scheme System**
- 7 distinct colors for different task categories
- Automatic color mapping based on keywords
- Professional RGB color definitions
- Easy to extend and customize

### 2. **Milestone Detection and Rendering**
- Automatic detection of single-day tasks as milestones
- Special visual indicators (★)
- Different rendering for milestones vs. tasks
- Clear visual hierarchy

### 3. **Professional LaTeX Styling**
- Modern typography with Helvetica font
- Comprehensive package usage
- Professional table styling
- Clean, readable output

### 4. **Robust CSV Processing**
- Multiple date format support
- Graceful error handling
- Input validation
- Memory-efficient processing

### 5. **Timeline and List Views**
- Beautiful enumerated timeline
- Detailed task list tables
- Category color coding
- Professional formatting

---

## 🔄 Migration Strategy for latex-yearly-planner

### Phase 1: Foundation
1. **Add comprehensive color scheme system**
   - Create `ColorScheme` struct with category colors
   - Add color mapping configuration
   - Implement dynamic color definition generation

2. **Enhance task model**
   - Add milestone detection
   - Add computed properties (duration, color)
   - Improve validation

### Phase 2: Processing
3. **Improve CSV parsing**
   - Add multiple date format support
   - Enhance error handling
   - Add input validation

4. **Add timeline view generation**
   - Create timeline list view
   - Add milestone indicators
   - Implement professional formatting

### Phase 3: Polish
5. **Enhance LaTeX generation**
   - Improve typography
   - Add professional styling
   - Enhance table formatting

6. **Add task list generation**
   - Create detailed task tables
   - Add category color coding
   - Implement professional layout

---

## 📝 Code Quality Lessons

### 1. **Type Safety**
- Use strong typing throughout
- Leverage Go's type system effectively
- Avoid `interface{}` where possible
- Use custom types for domain concepts

### 2. **Error Handling**
- Return errors explicitly
- Provide helpful error messages
- Use error wrapping for context
- Handle edge cases gracefully

### 3. **Configuration**
- Use structured configuration
- Provide sensible defaults
- Support environment variables
- Make configuration testable

### 4. **Testing**
- Write comprehensive tests
- Test error conditions
- Mock external dependencies
- Use table-driven tests

### 5. **Documentation**
- Clear function and type documentation
- README with usage examples
- Code comments for complex logic
- Architecture documentation

---

## 🎯 Success Metrics

After implementing these lessons, the `latex-yearly-planner` should have:

- ✅ **Professional output** with modern typography and styling
- ✅ **Robust processing** with comprehensive error handling
- ✅ **Flexible configuration** with easy customization
- ✅ **Clear architecture** with separation of concerns
- ✅ **Comprehensive features** including timeline and list views
- ✅ **Type safety** with strong typing throughout
- ✅ **Easy maintenance** with clean, documented code

---

## 🔗 Key Files to Reference

- `src/config.py` - Configuration management patterns
- `src/models.py` - Data model design
- `src/processor.py` - CSV processing strategies
- `src/generator.py` - LaTeX generation patterns
- `src/app.py` - Application architecture
- `src/utils.py` - Utility function patterns

---

## 🔧 Development

```bash
# Setup environment (downloads and vendors dependencies)
./scripts/setup.sh

# Format and vet
make fmt
make vet

# Clean generated files
make clean

# Build only (without PDF generation)
go build -mod=vendor -o generated/plannergen ./cmd/planner

# Run Go tests
go test -mod=vendor ./...
```

---

## 🐛 Troubleshooting

If PDF generation fails:

1. Check that XeLaTeX is installed: `xelatex --version`
2. Verify CSV data format in `input_data/Research Timeline v5 - Comprehensive.csv`
3. Check generated files: `ls -la generated/`
4. Review logs: `cat generated/*.log`

For offline builds:

- Run `./scripts/setup.sh` first to download and vendor dependencies
- The `vendor/` directory contains all dependencies for offline use

---

## 📚 Documentation

- [Developer Guide](docs/developer-guide.md)
- [User Guide](docs/user-guide.md)

---

## 📞 Support & Contributing

### For Users
1. **Start here**: Check the user guide for complete feature documentation
2. **Examples**: Review example configurations for real-world usage
3. **Templates**: Use pre-built templates for common scenarios

### For Developers
1. **Development Setup**: Follow the developer guide for environment setup
2. **API Reference**: Consult technical documentation for all APIs
3. **Contributing**: Follow contribution guidelines for code changes

### Getting Help
1. Check documentation first
2. Look at examples for similar use cases
3. Review developer guide for technical details
4. Check lessons learned for common issues

To contribute to this project:
1. Follow existing structure and naming conventions
2. Use clear, concise language
3. Include examples where helpful
4. Update documentation when adding new features

---

*This comprehensive guide serves as the complete reference for the PhD Dissertation Planner project, combining project overview, architectural lessons, and organizational planning into a single authoritative document.*

*Last updated: September 2025*
