# 📚 Architecture Overview

System design and architectural patterns for the PhD Dissertation Planner.

## 📖 Quick Links

- **[API Reference](api-reference.md)** - Technical API documentation
- **[Configuration Reference](configuration.md)** - Configuration options
- **[Developer Guide](../developer/developer-guide.md)** - Development setup
- **[Code Quality](CODE_QUALITY.md)** - Code quality standards
- **[Main README](../../README.md)** - Project overview

## System Architecture

### High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Interface │    │   Application   │    │  LaTeX Engine   │
│                 │    │    Logic        │    │                 │
│ • Command line  │◄──►│ • Task processing│◄──►│ • PDF generation│
│ • Configuration │    │ • Template       │    │ • XeLaTeX       │
│ • Validation    │    │   rendering      │    │ • Font handling │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CSV Input     │    │  Core Business  │    │   Templates     │
│                 │    │    Logic        │    │                 │
│ • Data parsing  │    │ • Calendar gen  │    │ • LaTeX macros  │
│ • Validation    │    │ • Task stacking │    │ • Styling       │
│ • Error handling│    │ • Color mapping │    │ • Layouts       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Package Structure

### Core Package (`src/core`)

**Purpose**: Core business logic, configuration, and shared utilities

**Responsibilities**:
- Configuration management and validation
- Data models (Task, Category, etc.)
- CSV parsing and validation
- Error handling and logging
- Color generation algorithms
- Default value management

**Key Files**:
- `config.go` - Configuration structs and loading
- `task.go` - Data models and business logic
- `reader.go` - CSV processing and validation
- `errors.go` - Error types and handling
- `logger.go` - Logging utilities
- `defaults.go` - Default configuration values
- `color_utils.go` - Color conversion utilities

### App Package (`src/app`)

**Purpose**: Main application logic and CLI interface

**Responsibilities**:
- Command-line argument parsing
- Application orchestration
- Template rendering coordination
- File I/O operations
- Progress reporting

**Key Files**:
- `cli.go` - Command-line interface
- `generator.go` - Main PDF generation logic
- `template_funcs.go` - Template helper functions

### Calendar Package (`src/calendar`)

**Purpose**: Calendar layout and task positioning

**Responsibilities**:
- Calendar grid generation
- Task stacking and overlap detection
- Visual layout calculations
- LaTeX-compatible output formatting

**Key Files**:
- `calendar.go` - Calendar generation and rendering
- `task_stacker.go` - Task overlap detection and positioning
- `cell_builder.go` - Individual cell rendering

### Shared Package (`src/shared`)

**Purpose**: Reusable components and templates

**Responsibilities**:
- Embedded LaTeX templates
- Template rendering utilities
- Shared rendering functions

**Key Files**:
- `templates/embed.go` - Embedded template assets
- `templates/rendering.go` - Template rendering logic

## Data Flow

### 1. Input Processing

```
CSV File → Reader → Validation → Task Models → Business Logic
```

1. **CSV Reading**: Raw CSV data parsed into Task structs
2. **Validation**: Data quality checks and error aggregation
3. **Transformation**: Tasks converted to calendar-compatible format
4. **Processing**: Business rules applied (colors, dates, categories)

### 2. Calendar Generation

```
Tasks → Calendar Generator → Layout Engine → Task Stacker → LaTeX Output
```

1. **Month Creation**: Calendar grid generated for target period
2. **Task Assignment**: Tasks assigned to appropriate days
3. **Layout Calculation**: Visual positioning and sizing
4. **Overlap Resolution**: Stacking algorithm prevents overlaps
5. **LaTeX Rendering**: Calendar converted to LaTeX markup

### 3. PDF Generation

```
LaTeX Templates + Data → Template Engine → XeLaTeX → PDF Output
```

1. **Template Loading**: Embedded LaTeX templates retrieved
2. **Data Binding**: Calendar data injected into templates
3. **LaTeX Compilation**: XeLaTeX processes markup to PDF
4. **Font Rendering**: System fonts embedded in output

## Key Design Patterns

### 1. **Configuration as Code**

Configuration managed through structured Go structs with YAML binding:

```go
type Config struct {
    Layout LayoutConfig `yaml:"layout"`
    LayoutEngine LayoutEngine `yaml:"layout_engine"`
}

// With getter methods for defaults
func (c *Config) GetDayNumberWidth() string {
    return getStringWithDefault(c.Layout.DayNumberWidth, Defaults.DayNumberWidth)
}
```

**Benefits**:
- Type safety
- IDE support
- Validation at compile time
- Clear default handling

### 2. **Embedded Templates**

LaTeX templates embedded in binary using Go's `embed` package:

```go
//go:embed templates/*.tpl
var templateFiles embed.FS

func GetEmbeddedTemplates() map[string]string {
    // Return map of template name → content
}
```

**Benefits**:
- Single binary deployment
- No external template dependencies
- Version consistency
- Fast loading

### 3. **Error Aggregation**

Custom error types with context accumulation:

```go
type ErrorAggregator struct {
    Errors   []error
    Warnings []string
}

func (ea *ErrorAggregator) AddError(err error) {
    ea.Errors = append(ea.Errors, err)
}
```

**Benefits**:
- Comprehensive error reporting
- Warning vs error distinction
- Context preservation
- User-friendly messages

### 4. **Task Stacking Algorithm**

Intelligent overlap detection using vertical tracks:

```go
type TaskStacker struct {
    tasks []*SpanningTask
    dayStacks map[string]*DayTaskStack
}

func (ts *TaskStacker) findAvailableTrack(task *SpanningTask) int {
    // Find lowest available vertical position
}
```

**Benefits**:
- Automatic layout optimization
- Visual clarity maintenance
- Scalable to large task sets
- Deterministic results

## Component Interactions

### Dependency Injection

Components receive dependencies through constructors:

```go
func NewReader(config *Config) *Reader {
    return &Reader{config: config}
}

func NewGenerator(config *Config) *Generator {
    return &Generator{
        config: config,
        templateDir: "templates",
    }
}
```

### Interface Segregation

Small, focused interfaces:

```go
type Reader interface {
    ReadTasks(filename string) ([]Task, error)
}

type Generator interface {
    Generate(tasks []Task) error
}
```

### Single Responsibility

Each package has one primary concern:
- `core`: Business logic and data
- `app`: Application coordination
- `calendar`: Layout and positioning
- `shared`: Common utilities

## Performance Considerations

### 1. **Memory Management**

- **Streaming CSV processing**: Large files processed in chunks
- **Template caching**: Compiled templates reused
- **Lazy evaluation**: Calculations performed only when needed

### 2. **Algorithm Complexity**

- **Task stacking**: O(n²) acceptable for typical workloads (< 100 tasks/month)
- **Color generation**: O(1) hash-based consistent coloring
- **Calendar generation**: O(days × tasks) linear scaling

### 3. **I/O Optimization**

- **Buffered writing**: LaTeX output written in chunks
- **Embedded assets**: No filesystem access for templates
- **Parallel processing**: Independent operations can run concurrently

## Testing Strategy

### Unit Tests

- **Pure functions**: Color generation, date validation
- **Data transformation**: CSV parsing, task conversion
- **Error handling**: Invalid input, edge cases

### Integration Tests

- **End-to-end workflows**: CSV → PDF generation
- **Configuration loading**: YAML parsing and validation
- **Template rendering**: Data binding and output

### Test Coverage Goals

- **Core utilities**: 100% coverage
- **Business logic**: 80%+ coverage
- **Error paths**: All major error conditions tested
- **Integration paths**: Key user workflows tested

## Error Handling Strategy

### 1. **Error Types**

Custom error types for different failure modes:

```go
type ConfigError struct {
    Field   string
    Message string
    Value   interface{}
}

type DataError struct {
    Row     int
    Column  string
    Value   string
    Message string
}
```

### 2. **Error Propagation**

Errors wrapped with context:

```go
return fmt.Errorf("failed to parse CSV: %w", err)
```

### 3. **User-Friendly Messages**

Errors translated to actionable advice:

```go
if errors.Is(err, ErrMissingColumn) {
    return fmt.Errorf("CSV file missing required column. Add %s column with task data", missingCol)
}
```

## Configuration Management

### Hierarchical Overrides

Configuration loaded with precedence:

```
Defaults ← YAML File ← Environment Variables ← CLI Flags
```

### Validation

Configuration validated at startup:

- **Required fields**: Essential settings present
- **Value ranges**: Numeric values within bounds
- **File paths**: Referenced files exist
- **Cross-references**: Related settings consistent

## Extensibility Points

### 1. **Template System**

New output formats via additional templates:

```go
// Add HTML template alongside LaTeX
func renderHTML(tasks []Task) (string, error) {
    return RenderTemplate("calendar.html.tpl", data)
}
```

### 2. **Color Schemes**

Custom color algorithms:

```go
func CustomColorScheme(category string) string {
    // Implement custom color mapping
    return generateCustomColor(category)
}
```

### 3. **Layout Engines**

Alternative positioning algorithms:

```go
type LayoutEngine interface {
    PositionTasks(tasks []Task, calendar *Calendar) error
}
```

## Deployment Considerations

### Single Binary

Go compilation produces standalone executable:

```bash
# Build for current platform
go build -o planner ./cmd/planner

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o planner-linux ./cmd/planner
```

### Embedded Assets

Templates and defaults bundled in binary:

```go
//go:embed templates/*
//go:embed defaults.yaml
var embeddedFiles embed.FS
```

### Minimal Dependencies

Only requires:
- Go runtime (statically linked)
- XeLaTeX (system dependency)
- Input CSV file

---

*Architecture documentation last updated: October 2025*
