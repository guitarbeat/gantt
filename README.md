# 📅 PhD Dissertation Planner

A Go-based tool for generating professional PDF planners from CSV task data using LaTeX.

## 🚀 Quick Start

```bash
# Build the generator and produce the PDF (auto-detects first CSV in input/)
make clean-build

# Or run the generator manually from src/
cd src && go build -o ../output/plannergen . && \
  PLANNER_SILENT=1 PLANNER_CSV_FILE="../input/your_data.csv" \
  ../output/plannergen --config "config/base.yaml,config/monthly_calendar.yaml" --outdir ../output
```

## ✅ Status

[![CI](https://github.com/your-username/phd-dissertation-planner/actions/workflows/ci.yml/badge.svg)](https://github.com/your-username/phd-dissertation-planner/actions)

- ✅ **PDF Generation**: Working (generates ~116KB PDFs)
- ✅ **CSV Processing**: 84 tasks parsed successfully
- ✅ **LaTeX Compilation**: XeLaTeX integration working
- ✅ **Template System**: Go templates rendering correctly

## 📁 Project Structure

```
├── src/                    # Go source code
│   ├── internal/          # Core application logic
│   ├── templates/         # LaTeX templates
│   └── config/            # Configuration files
├── input/                 # CSV input data
├── output/               # Generated PDFs and logs
├── reference/            # Documentation
└── Makefile             # Build automation
```

## 🔧 Development

```bash
# Format and vet
make fmt
make vet

# Clean generated files
make clean

# Build only (without PDF generation)
cd src && go build ./...

# Run Go tests (if present)
cd src && go test ./...
```

## 📚 Documentation

- [Developer Guide](reference/docs/developer-guide/README.md)
- [User Guide](reference/docs/user-guide/README.md)
- [Examples](reference/examples/README.md)

## 🐛 Troubleshooting

If PDF generation fails:
1. Check that XeLaTeX is installed: `xelatex --version`
2. Verify CSV data format in `input/data.cleaned.csv`
3. Check generated files: `ls -la src/build/`
4. Review logs: `cat src/build/*.log`

---

*Last updated: September 2024*
