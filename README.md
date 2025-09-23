# 📅 PhD Dissertation Planner

A Go-based tool for generating professional PDF planners from CSV task data using LaTeX.

## 🚀 Quick Start

```bash
# Generate PDF from CSV data
make test

# Generate PDF with custom data
make pdf CSV=../input/your_data.csv OUTPUT=your_planner
```

## ✅ Status

[![PDF Generation Test](https://github.com/your-username/phd-dissertation-planner/workflows/Test%20PDF%20Generation/badge.svg)](https://github.com/your-username/phd-dissertation-planner/actions)

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
# Build and test
make test

# Clean generated files
make clean

# Run Go tests only
cd src && go test ./internal/...
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
