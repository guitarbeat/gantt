# PhD Dissertation Planner

A Go-based application for generating LaTeX academic planners and calendars from CSV data.

## Features

- **CSV-based Task Management**: Import tasks from CSV files with support for dates, categories, and priorities
- **Monthly Calendar Generation**: Create detailed monthly calendars with task overlays
- **LaTeX Output**: Generate professional LaTeX documents for academic planning
- **Configurable Layouts**: Customize paper sizes, margins, and visual themes
- **Template System**: Flexible template-based rendering with embedded and filesystem support

## Installation

### Prerequisites

- Go 1.16 or later
- LaTeX distribution (for PDF generation)

### Build from Source

```bash
git clone <repository-url>
cd phd-dissertation-planner/src
go build -o plannergen
```

## Usage

### Basic Usage

```bash
# Generate planner with default configuration
./plannergen

# Use custom configuration file
./plannergen -config config/my-config.yaml

# Preview mode (one page per unique module)
./plannergen -preview

# Specify output directory
./plannergen -outdir /path/to/output
```

### Configuration

The application supports YAML configuration files with the following structure:

```yaml
# Basic settings
year: 2024
week_start: "Monday"
dotted: false
cal_after_schedule: true
clear_top_right_corner: false
ampm_time: false
add_last_half_hour: false

# Data source
csv_file: "data/tasks.csv"
start_year: 2024
end_year: 2025

# Layout settings
layout:
  paper:
    width: "8.5in"
    height: "11in"
    margin:
      top: "1in"
      bottom: "1in"
      left: "1in"
      right: "1in"
  
  colors:
    gray: "#CCCCCC"
    light_gray: "#F5F5F5"
  
  latex:
    tab_col_sep: "1.5pt"
    header_side_months_width: "2cm"
    task_border_width: "0.5pt"
    task_padding_h: "2pt"
    task_padding_v: "1pt"
    task_vertical_offset: "0pt"
    array_stretch: 1.2

# Pages configuration
pages:
  - name: "monthly"
    render_blocks:
      - func_name: "monthly"
        tpls: ["monthly.tpl"]
```

### CSV Data Format

The application expects CSV files with the following columns:

```csv
ID,Name,Start Date,End Date,Category,Description,Priority,Status,Assignee,Parent ID,Dependencies,Is Milestone
task-1,Research Proposal,2024-01-15,2024-03-15,PROPOSAL,Write PhD proposal,1,Planned,Student,,,false
task-2,Literature Review,2024-02-01,2024-04-30,RESEARCH,Review relevant papers,2,Planned,Student,,,false
```

### Environment Variables

- `PLANNER_YEAR`: Default year for calendar generation
- `PLANNER_CSV_FILE`: Path to CSV file with task data
- `PLANNER_START_YEAR`: Start year for date range
- `PLANNER_END_YEAR`: End year for date range
- `PLANNER_OUTPUT_DIR`: Output directory for generated files
- `PLANNER_LAYOUT_PAPER_WIDTH`: Paper width
- `PLANNER_LAYOUT_PAPER_HEIGHT`: Paper height
- `PLANNER_LAYOUT_PAPER_MARGIN_TOP`: Top margin
- `PLANNER_LAYOUT_PAPER_MARGIN_BOTTOM`: Bottom margin
- `PLANNER_LAYOUT_PAPER_MARGIN_LEFT`: Left margin
- `PLANNER_LAYOUT_PAPER_MARGIN_RIGHT`: Right margin
- `DEV_TEMPLATES`: Use filesystem templates instead of embedded ones

## Development

### Project Structure

```
src/
├── main.go                    # Application entry point
├── internal/
│   ├── application/          # CLI application logic
│   │   └── core.go
│   ├── common/              # Shared types and utilities
│   │   ├── config.go        # Configuration handling
│   │   ├── task.go          # Task data structures
│   │   ├── reader.go        # CSV reading functionality
│   │   └── errors.go        # Custom error types
│   └── scheduler/           # Calendar and layout logic
│       ├── calendar.go      # Calendar generation
│       └── layout_manager.go # Layout and positioning
└── templates/               # LaTeX templates
    ├── embed.go            # Embedded template files
    └── rendering.go         # Template helper functions
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Code Quality

```bash
# Format code
gofmt -w .

# Check for issues
go vet ./...

# Update dependencies
go get -u ./...
go mod tidy

# Check for security vulnerabilities
govulncheck ./...
```

### Building

```bash
# Build for current platform
go build -o plannergen

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o plannergen-linux-amd64
GOOS=windows GOARCH=amd64 go build -o plannergen-windows-amd64.exe
GOOS=darwin GOARCH=amd64 go build -o plannergen-darwin-amd64
```

## API Reference

### Core Types

#### Task
Represents a single task with timing, categorization, and metadata.

```go
type Task struct {
    ID           string
    Name         string
    StartDate    time.Time
    EndDate      time.Time
    Category     string
    Description  string
    Priority     int
    Status       string
    Assignee     string
    ParentID     string
    Dependencies []string
    IsMilestone  bool
}
```

#### Config
Application configuration with layout, data source, and output settings.

```go
type Config struct {
    Year                int
    WeekStart           time.Weekday
    CSVFilePath         string
    StartYear           int
    EndYear             int
    MonthsWithTasks     []MonthYear
    Pages               Pages
    Layout              Layout
    OutputDir           string
}
```

### Template Functions

The application provides several template functions for LaTeX generation:

- `dict(values...)`: Create key-value dictionaries
- `incr(i)`: Increment integer
- `dec(i)`: Decrement integer
- `is(value)`: Check if value is truthy
- `hasLayoutData(data)`: Check for layout data
- `getTaskBars(data)`: Get task bars from data
- `getLayoutStats(data)`: Get layout statistics
- `formatTaskBar(bar)`: Format task bar for LaTeX

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **Template not found**: Ensure templates are properly embedded or DEV_TEMPLATES is set
2. **CSV parsing errors**: Check CSV format and column headers
3. **LaTeX compilation errors**: Verify LaTeX installation and template syntax
4. **Permission errors**: Ensure output directory is writable

### Debug Mode

Enable debug output by setting environment variables:

```bash
export PLANNER_DEBUG=true
./plannergen
```

## Changelog

### v1.0.0
- Initial release
- CSV-based task management
- Monthly calendar generation
- LaTeX output support
- Configurable layouts