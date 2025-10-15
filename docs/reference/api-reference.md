# 🔧 API Reference

Technical API documentation for the PhD Dissertation Planner.

## 📖 Quick Links

- **[Configuration Reference](configuration.md)** - Configuration options
- **[Architecture](architecture.md)** - System design patterns
- **[Developer Guide](../developer/developer-guide.md)** - Development setup
- **[Code Quality](CODE_QUALITY.md)** - Code quality standards
- **[Main README](../../README.md)** - Project overview

## Package Overview

### Core Packages

- **`core`** - Configuration, data models, and utilities
- **`app`** - Main application logic and CLI interface
- **`calendar`** - Calendar generation and task positioning
- **`shared`** - Shared templates and rendering utilities

## Core Package (`src/core`)

### Configuration (`config.go`)

#### Types

```go
type Config struct {
    Debug     DebugConfig     `yaml:"debug"`
    Year      int             `yaml:"year"`
    StartYear int             `yaml:"start_year"`
    EndYear   int             `yaml:"end_year"`
    OutputDir string          `yaml:"output_dir"`

    Layout LayoutConfig       `yaml:"layout"`
    LayoutEngine LayoutEngine `yaml:"layout_engine"`
}
```

#### Functions

```go
// NewConfig creates a new configuration with defaults
func NewConfig() (*Config, error)

// LoadConfig loads configuration from YAML file
func LoadConfig(filename string) (*Config, error)

// GetYears returns slice of years to generate
func (c *Config) GetYears() []int

// Configuration getters (return config value or default)
func (c *Config) GetDayNumberWidth() string
func (c *Config) GetDayContentMargin() string
func (c *Config) GetTaskCellMargin() string
func (c *Config) GetTaskCellSpacing() string
func (c *Config) GetHeaderAngleSizeOffset() string
func (c *Config) GetHyphenPenalty() int
func (c *Config) GetTolerance() int
func (c *Config) GetEmergencyStretch() string
func (c *Config) GetOutputDir() string
func (c *Config) GetYear() int
func (c *Config) IsDebugMode() bool
```

### Data Models (`task.go`)

#### Types

```go
type Task struct {
    ID          string    `csv:"ID"`
    Name        string    `csv:"Task"`
    StartDate   time.Time `csv:"StartDate"`
    EndDate     time.Time `csv:"EndDate"`
    Phase       string    `csv:"Phase"`
    SubPhase    string    `csv:"SubPhase"`
    Category    string    `csv:"Category"`
    Description string    `csv:"Description"`
    Status      string    `csv:"Status"`
    Assignee    string    `csv:"Assignee"`
}

type TaskCategory struct {
    Name        string
    DisplayName string
    Color       string
    Description string
}
```

#### Functions

```go
// GenerateCategoryColor creates consistent color for category
func GenerateCategoryColor(category string) string
```

### CSV Reader (`reader.go`)

#### Types

```go
type Reader struct {
    config *Config
}

type ErrorAggregator struct {
    Errors   []error
    Warnings []string
}
```

#### Functions

```go
// NewReader creates new CSV reader
func NewReader(config *Config) *Reader

// ReadTasks reads and validates CSV data
func (r *Reader) ReadTasks(filename string) ([]Task, error)

// Validation functions
func (r *Reader) ValidateTask(task Task) error
func (r *Reader) ValidateDateRange(start, end time.Time) error
```

### Error Handling (`errors.go`)

#### Types

```go
type ConfigError struct {
    Field   string
    Message string
    Value   interface{}
}

type FileError struct {
    Path    string
    Op      string
    Err     error
}

type TemplateError struct {
    Template string
    Err      error
}

type DataError struct {
    Row     int
    Column  string
    Value   string
    Message string
}
```

#### Functions

```go
// Error constructors
func NewConfigError(field, message string, value interface{}) *ConfigError
func NewFileError(path, op string, err error) *FileError
func NewTemplateError(template string, err error) *TemplateError
func NewDataError(row int, column, value, message string) *DataError
```

### Logging (`logger.go`)

#### Functions

```go
// Logging functions
func Info(format string, args ...interface{})
func Error(format string, args ...interface{})
func Debug(format string, args ...interface{})
func IsSilent() bool
```

### Defaults (`defaults.go`)

#### Functions

```go
// Get default configurations
func DefaultConfig() *Config
func DefaultLayoutConfig() LayoutConfig
func DefaultLayoutEngine() LayoutEngine
```

### Color Utilities (`color_utils.go`)

#### Functions

```go
// HexToRGB converts hex color to RGB format for LaTeX
func HexToRGB(hex string) string
```

## App Package (`src/app`)

### CLI Interface (`cli.go`)

#### Types

```go
type CLI struct {
    config *core.Config
}
```

#### Functions

```go
// NewCLI creates new CLI interface
func NewCLI(config *core.Config) *CLI

// Run executes the CLI application
func (c *CLI) Run(args []string) error
```

### Generator (`generator.go`)

#### Types

```go
type Generator struct {
    config      *core.Config
    templateDir string
}
```

#### Functions

```go
// NewGenerator creates new PDF generator
func NewGenerator(config *core.Config) *Generator

// Generate creates PDF from tasks
func (g *Generator) Generate(tasks []core.Task) error

// GenerateWithPreview generates with preview mode
func (g *Generator) GenerateWithPreview(tasks []core.Task, preview bool) error
```

### Template Functions (`template_funcs.go`)

#### Functions

```go
// TemplateFuncs returns FuncMap for templates
func TemplateFuncs() template.FuncMap

// Individual template functions
func dictFunc(values ...interface{}) map[string]interface{}
func incrFunc(x int) int
func decrFunc(x int) int
func isFunc(a, b interface{}) bool
```

## Calendar Package (`src/calendar`)

### Calendar Generation (`calendar.go`)

#### Types

```go
type Month struct {
    Year      int
    Month     time.Month
    Weeks     []Week
    TaskCount int
}

type Week struct {
    WeekNumber int
    Days       []Day
}

type Day struct {
    Time        time.Time
    Tasks       []SpanningTask
    IsCurrentMonth bool
}

type SpanningTask struct {
    ID          string
    Name        string
    Description string
    Phase       string
    SubPhase    string
    Category    string
    StartDate   time.Time
    EndDate     time.Time
    Color       string
    Progress    int
}
```

#### Functions

```go
// NewMonth creates calendar month
func NewMonth(year int, month time.Month) *Month

// GenerateMonth generates month with tasks
func GenerateMonth(year int, month time.Month, tasks []SpanningTask) *Month

// GetTaskColors returns color map for tasks
func (m *Month) GetTaskColors() map[string]string

// GetTaskColorsByPhase returns colors grouped by phase
func (m *Month) GetTaskColorsByPhase() []PhaseGroup

// CreateSpanningTask converts Task to SpanningTask
func CreateSpanningTask(task core.Task, startDate, endDate time.Time) SpanningTask
```

### Task Stacker (`task_stacker.go`)

#### Types

```go
type TaskStacker struct {
    tasks        []*SpanningTask
    dayStacks    map[string]*DayTaskStack
    maxTracks    int
    weekStartDay time.Weekday
}

type TaskStack struct {
    Track    int
    Task     *SpanningTask
    StartCol int
    EndCol   int
}

type DayTaskStack struct {
    Date   time.Time
    Stacks []TaskStack
}
```

#### Functions

```go
// NewTaskStacker creates task stacker
func NewTaskStacker(tasks []*SpanningTask, weekStartDay time.Weekday) *TaskStacker

// ComputeStacks calculates task positions
func (ts *TaskStacker) ComputeStacks()

// GetMaxTracks returns maximum tracks needed
func (ts *TaskStacker) GetMaxTracks() int

// GetDayStack returns tasks for specific day
func (ts *TaskStacker) GetDayStack(date time.Time) *DayTaskStack
```

## Shared Package (`src/shared`)

### Templates (`templates/embed.go`, `rendering.go`)

#### Functions

```go
// GetEmbeddedTemplates returns embedded template files
func GetEmbeddedTemplates() map[string]string

// RenderTemplate renders template with data
func RenderTemplate(name string, data interface{}) (string, error)

// RenderTemplateWithFuncs renders with custom functions
func RenderTemplateWithFuncs(name string, data interface{}, funcs template.FuncMap) (string, error)
```

## Command Line Interface

### Global Flags

```bash
# Configuration
--config string     Configuration file path
--year int         Year to generate (default current year)
--preset string    Layout preset (academic|compact|presentation)

# Input/Output
--csv-file string   CSV input file path
--output-dir string Output directory

# Debug/Development
--debug            Enable debug output
--validate         Validate CSV without generating PDF
--preview          Generate preview mode
--verbose          Verbose output
```

### Examples

```bash
# Basic usage
go run ./cmd/planner

# With custom config
go run ./cmd/planner --config custom.yaml --year 2025

# Validation only
go run ./cmd/planner --validate --csv-file data/tasks.csv

# Debug mode
go run ./cmd/planner --debug --verbose
```

## Performance & Testing

### Benchmarks

The project includes comprehensive performance benchmarks:

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmarks
go test -bench=BenchmarkCSVReading ./src/core
go test -bench=BenchmarkConfigurationLoading ./src/app
```

#### Current Performance Metrics

- **CSV Reading**: ~5.6 µs/op, 208 KB/op, 3405 allocs/op
- **Configuration Loading**: ~86 µs/op, 100 KB/op, 1949 allocs/op
- **Template Rendering**: ~50-200 µs/op (varies by complexity)

### Test Coverage

Current test coverage by package:
- **App**: 7.7% (utility functions and template helpers)
- **Core**: 26.2% (configuration, validation, color utilities)
- **Calendar**: 16.1% (LaTeX escaping, phase descriptions, spanning tasks)
- **Overall**: ~16% (excluding generated code and templates)

## Error Codes

- `0` - Success
- `1` - General error
- `2` - Configuration error
- `3` - File I/O error
- `4` - Template error
- `5` - Data validation error
- `130` - Interrupted (Ctrl+C)

## Environment Variables

```bash
PLANNER_CONFIG_FILE    # Path to config file
PLANNER_CSV_FILE       # Path to CSV file
PLANNER_YEAR          # Year to generate
PLANNER_DEBUG         # Enable debug mode
PLANNER_SILENT        # Suppress output
```

---

*API reference last updated: October 2025*
