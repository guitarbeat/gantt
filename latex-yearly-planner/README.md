# LaTeX Yearly Planner

A Go application for generating LaTeX yearly planners from CSV data.

## Project Structure

This project follows the standard Go project layout:

```text
latex-yearly-planner/
├── cmd/                    # Application entry points
│   └── plannergen/         # Main CLI application
├── internal/               # Private application code
│   ├── app/               # CLI application setup
│   ├── config/            # Configuration management
│   ├── data/              # CSV data reading
│   ├── generator/         # LaTeX generation logic
│   └── layout/            # Layout measurements
├── pkg/                   # Reusable components
│   ├── calendar/          # Calendar data structures
│   ├── header/            # Header components
│   └── latex/             # LaTeX utilities
├── configs/               # Configuration files
├── templates/             # LaTeX templates
│   ├── common/           # Common templates
│   ├── breadcrumbs/      # Breadcrumb templates
│   └── layouts/          # Layout templates
├── scripts/              # Build scripts
└── build/                # Build output (gitignored)
```

## Building

```bash
go build -o build/plannergen ./cmd/plannergen
```

## Usage

```bash
./build/plannergen --config configs/planner_config.yaml
```

## Configuration

The application uses YAML configuration files in the `configs/` directory:

- `base.yaml` - Base configuration
- `planner_config.yaml` - Planner-specific settings
- `page_template.yaml` - Page template configuration

## Templates

LaTeX templates are organized in the `templates/` directory:

- `common/` - Shared template components
- `breadcrumbs/` - Navigation breadcrumb templates
- `layouts/` - Page layout templates
- `document.tpl` - Main document template
- `macro.tpl` - LaTeX macro definitions

## Development

The project is organized to follow Go best practices:

- **cmd/**: Contains the main application entry point
- **internal/**: Private application code that cannot be imported by other projects
- **pkg/**: Public packages that can be imported by other projects
- **configs/**: Configuration files separate from code
- **templates/**: Template files organized by purpose
- **build/**: Build artifacts (gitignored)
