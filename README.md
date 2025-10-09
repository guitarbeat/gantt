# 📚 PhD Dissertation Planner

A Go-based application that transforms CSV data into professional LaTeX-generated PDF planners and Gantt charts for academic project management.

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd phd-dissertation-planner

# Install dependencies and build
make install
make build

# Generate your first planner
make run
```

## 📖 Documentation

| Document                                                   | Description                        |
| ---------------------------------------------------------- | ---------------------------------- |
| **[📋 User Guide](docs/user/user-guide.md)**                | How to use the planner effectively |
| **[🛠️ Developer Guide](docs/developer/developer-guide.md)** | Development setup and contributing |
| **[⚙️ Configuration](docs/reference/configuration.md)**     | Configuration options and presets  |
| **[🔧 API Reference](docs/reference/api-reference.md)**     | Technical API documentation        |
| **[📚 Architecture](docs/reference/architecture.md)**       | System design and patterns         |

## 🎯 Key Features

- **CSV to LaTeX**: Transform spreadsheet data into professional PDFs
- **Task Stacking**: Intelligent overlap detection and visual layering
- **Multiple Layouts**: Academic, compact, and presentation presets
- **LaTeX Integration**: Professional typography with XeLaTeX
- **Release Management**: Timestamped releases with full provenance

## 🚀 Usage

### Basic Usage

```bash
# Generate PDF from CSV
go run ./cmd/planner

# Use different layout preset
go run ./cmd/planner --preset compact

# Validate CSV without generating PDF
go run ./cmd/planner --validate
```

### Advanced Usage

```bash
# Custom configuration
go run ./cmd/planner --config my-config.yaml

# Generate release with timestamp
./scripts/build_release.sh --name "Committee_Review"

# Preview mode (development)
go run ./cmd/planner --preview
```

---

## 📁 Project Structure

```
├── cmd/planner/           # Application entry point
├── src/                   # Source code
│   ├── app/              # Main application logic
│   ├── core/             # Core utilities & config
│   ├── calendar/         # Calendar generation
│   └── shared/           # Shared templates & utils
├── docs/                 # Documentation
├── scripts/              # Build & utility scripts
├── input_data/           # CSV input files
├── generated/            # Generated PDFs & LaTeX
├── releases/             # Release archives
└── tests/                # Test suites
```

---

## 🔧 Development

```bash
# Setup development environment
./scripts/setup.sh

# Run tests
make test

# Format and lint
make fmt && make lint

# Clean build
make clean-build
```

See [Developer Guide](docs/developer/developer-guide.md) for detailed development instructions.

---

## 📞 Support

- **📖 [User Guide](docs/user/user-guide.md)** - Complete usage documentation
- **🐛 [Troubleshooting](docs/user/troubleshooting.md)** - Common issues and solutions
- **🔧 [Developer Guide](docs/developer/developer-guide.md)** - Contributing guidelines

---

*Last updated: October 2025*
