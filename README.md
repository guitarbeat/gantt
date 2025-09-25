# 📅 PhD Dissertation Planner

A Go-based tool for generating professional PDF planners from CSV task data using LaTeX.

## 🚀 Quick Start

```bash
# Setup development environment (downloads dependencies locally for offline use)
./scripts/setup.sh

# Build the generator and produce the PDF
make clean-build

# Or run the generator manually
go build -mod=vendor -o generated/plannergen ./cmd/planner && \
  PLANNER_SILENT=1 PLANNER_CSV_FILE="input_data/Research Timeline v5 - Comprehensive.csv" \
  ./generated/plannergen --config "src/core/base.yaml,src/core/monthly_calendar.yaml" --outdir generated
```

**Note**: Dependencies are vendored locally, so the project works offline after initial setup.

## ✅ Status

[![CI](https://github.com/your-username/phd-dissertation-planner/actions/workflows/ci.yml/badge.svg)](https://github.com/your-username/phd-dissertation-planner/actions)

- ✅ **PDF Generation**: Working (generates ~116KB PDFs)
- ✅ **CSV Processing**: 84 tasks parsed successfully
- ✅ **LaTeX Compilation**: XeLaTeX integration working
- ✅ **Template System**: Go templates rendering correctly

## 📁 Project Structure

```
├── cmd/planner/           # Go application entry point
├── src/                   # Source code (beginner-friendly)
│   ├── app/              # Main application logic
│   ├── core/             # Core utilities and shared logic
│   ├── calendar/         # Calendar/scheduling functionality
│   ├── shared/           # Shared/reusable code
│   │   └── templates/    # LaTeX templates (embedded)
│   └── assets/           # Small runtime assets (embedded)
├── input_data/           # Input data files (CSV, etc.)
├── generated/            # Generated output files (PDFs, logs)
├── static_assets/        # Static files/assets
├── vendor/               # Vendored Go dependencies (for offline builds)
├── scripts/              # Setup and build scripts
└── docs/                 # Documentation
```

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

## 📚 Documentation

- [Developer Guide](docs/developer-guide/README.md)
- [User Guide](docs/user-guide/README.md)

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

*Last updated: September 2025*
