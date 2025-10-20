# PhD Dissertation Planner

A Go-based tool for generating LaTeX calendar PDFs from CSV timeline data, designed for academic research planning and project management.

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- XeLaTeX (for PDF generation)
- Cursor CLI (optional, for AI-powered development)

### Installation
```bash
# Clone and setup
git clone <repository-url>
cd gantt
./scripts/unified.sh setup

# Build and run
./scripts/unified.sh build pdf
```

### Basic Usage
```bash
# Generate calendar from CSV
./generated/plannergen --config configs/base.yaml --outdir generated

# View generated PDF
open generated/monthly_calendar.pdf
```

## 📋 Features

- **CSV Timeline Processing**: Import research timelines from CSV files
- **LaTeX PDF Generation**: High-quality calendar PDFs with XeLaTeX
- **Multiple Layouts**: Academic, compact, and presentation formats
- **AI-Powered Development**: Cursor CLI integration for enhanced development
- **Comprehensive Testing**: Unit, integration, and performance tests
- **Cross-Platform**: Works on macOS, Linux, and Windows

## 🛠️ Development

### Unified Development Tool
The project includes a unified script that consolidates all development operations:

```bash
# Setup project
./scripts/unified.sh setup

# Development workflow
./scripts/unified.sh dev start

# Build and test
./scripts/unified.sh build pdf
./scripts/unified.sh test all

# Code quality
./scripts/unified.sh quality all

# Run CI pipeline
./scripts/unified.sh ci
```

### Available Commands
- `setup` - Setup dependencies and hooks
- `build [type]` - Build binary, LaTeX, or PDF
- `test [type]` - Run unit, integration, or coverage tests
- `quality [type]` - Format, lint, or security checks
- `dev [action]` - Development operations
- `maintenance [task]` - Clean, organize, or update dependencies
- `cursor [action]` - Cursor CLI operations
- `release [type]` - Build, package, or publish releases
- `ci` - Run full CI pipeline

### AI-Powered Development
With Cursor CLI installed:
```bash
# Install AI hooks
./scripts/unified.sh cursor hooks

# AI-powered development
./scripts/unified.sh cursor dev
./scripts/unified.sh cursor test

# Project statistics
./scripts/unified.sh cursor stats
```

## 📁 Project Structure

```
├── cmd/planner/           # Main application
├── internal/              # Internal packages
│   ├── app/              # Application logic
│   ├── calendar/         # Calendar generation
│   └── core/             # Core utilities
├── pkg/templates/        # LaTeX templates
├── configs/              # Configuration files
├── input_data/           # CSV input files
├── generated/            # Generated outputs
├── tests/                # Test files
├── scripts/              # Build and utility scripts
│   └── unified.sh        # Unified development tool
└── docs/                 # Documentation
```

## ⚙️ Configuration

### Configuration Files
- `configs/base.yaml` - Base configuration
- `configs/monthly_calendar.yaml` - Calendar-specific settings
- `configs/academic.yaml` - Academic layout preset
- `configs/compact.yaml` - Compact layout preset
- `configs/presentation.yaml` - Presentation layout preset

### Environment Variables
- `PLANNER_CSV_FILE` - Path to CSV input file
- `PLANNER_CONFIG_FILE` - Path to configuration file
- `PLANNER_OUTPUT_DIR` - Output directory for generated files
- `PLANNER_SILENT` - Silent mode (0/1)

## 🧪 Testing

### Test Types
- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end workflow testing
- **Performance Tests**: Benchmark testing
- **Coverage Tests**: Code coverage analysis

### Running Tests
```bash
# All tests
./scripts/unified.sh test all

# Specific test types
./scripts/unified.sh test unit
./scripts/unified.sh test integration
./scripts/unified.sh test coverage
./scripts/unified.sh test bench
```

## 📚 Documentation

### Key Documents
- [User Guide](docs/user/user-guide.md) - Complete user documentation
- [Developer Guide](docs/developer/developer-guide.md) - Development setup and workflows
- [API Reference](docs/reference/api-reference.md) - API documentation
- [Configuration Guide](docs/reference/configuration.md) - Configuration options
- [Cursor CLI Integration](docs/developer/cursor-cli-integration.md) - AI-powered development

### Examples
- [Academic Research](docs/examples/academic-research.yaml) - Academic timeline example
- [Advanced Customization](docs/examples/advanced-customization.yaml) - Advanced configuration
- [Minimal Setup](docs/examples/minimal.yaml) - Minimal configuration

## 🚀 CI/CD

### GitHub Actions
- **CI Pipeline**: Automated testing and building
- **Code Quality**: Linting, formatting, and security checks
- **AI Enhancement**: Automated code improvements with Cursor CLI
- **Release Management**: Automated releases with GoReleaser

### Pre-commit Hooks
- **Cursor CLI Hooks**: AI-powered pre-commit checks
- **Traditional Hooks**: Standard code quality checks
- **Automatic Fixes**: AI-powered code improvements

## 📊 Performance

### Benchmarks
- CSV processing: ~1000 rows/second
- LaTeX generation: ~2-5 seconds
- PDF compilation: ~3-8 seconds
- Memory usage: ~50-100MB

### Optimization
- Concurrent processing for large datasets
- Template caching for repeated operations
- Memory-efficient CSV parsing
- Optimized LaTeX generation

## 🤝 Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Make changes with AI assistance
4. Run tests and quality checks
5. Submit a pull request

### Code Standards
- Go 1.22+ with Google style guide
- Comprehensive test coverage
- AI-powered code review
- Automated formatting and linting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- LaTeX community for excellent document generation
- Go community for robust tooling
- Cursor team for AI-powered development tools
- Academic researchers for feedback and requirements

---

**Built with ❤️ for the academic community**