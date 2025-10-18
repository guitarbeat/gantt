# Cross-Platform Development Guide

This project now supports cross-platform development on Windows, macOS, and Linux.

## Available Scripts

### Build and Preview Scripts

**PowerShell (Windows):**
```powershell
.\scripts\build_and_preview.ps1 [pages]
```

**Bash/Shell (Linux/macOS):**
```bash
./scripts/build_and_preview.sh [pages]
```

Both scripts perform the same function:
- Build the Go binary
- Generate LaTeX from CSV data
- Compile PDF using XeLaTeX
- Generate preview images (requires Python with pdf2image)

### Development Environment Scripts

**PowerShell (Windows):**
```powershell
.\scripts\dev.ps1 <command> [args...]
```

**Bash/Shell (Linux/macOS):**
```bash
./scripts/dev.sh <command> [args...]
```

Both scripts:
- Load development environment variables from `.config/.env.dev` (if present)
- Set up development defaults
- Execute the provided command

### Cleanup and Organization Scripts

**PowerShell (Windows):**
```powershell
.\scripts\cleanup_and_organize.ps1 [-ScatteredOnly] [-TestOnly] [-Status] [-Help]
```

**Bash/Shell (Linux/macOS):**
```bash
./scripts/cleanup_and_organize.sh [--scattered-only] [--test-only] [--status] [--help]
```

Both scripts:
- Organize project files into proper directories
- Clean up temporary and scattered files
- Update .gitignore
- Create project structure documentation

## Platform Detection

The scripts automatically detect the platform and use appropriate:
- Path separators (`\` on Windows, `/` on Unix)
- Executable extensions (`.exe` on Windows, none on Unix)
- Commands and tools

## Makefile Improvements

The Makefile has been updated for cross-platform compatibility:
- Cross-platform file size detection using `stat` commands
- Automatic script selection (bash vs PowerShell)
- Improved pre-commit hook installation

## Dependencies

### Required for all platforms:
- Go 1.19+
- XeLaTeX (TeX distribution)

### For PDF preview generation:
- Python 3
- pdf2image package
- Poppler (system library)

### Installation by platform:

**Windows:**
```powershell
# Using Chocolatey
choco install golang mingw latex poppler

# Python dependencies
pip install pdf2image
```

**macOS:**
```bash
# Using Homebrew
brew install go mactex poppler

# Python dependencies
pip3 install pdf2image
```

**Ubuntu/Debian:**
```bash
# System dependencies
sudo apt-get update
sudo apt-get install golang-go texlive-xetex texlive-latex-extra python3-pdf2image poppler-utils

# Python dependencies (if not installed by apt)
pip3 install pdf2image
```

## Usage Examples

### Quick build and preview (3 pages):
```bash
# Windows PowerShell
.\scripts\build_and_preview.ps1 3

# Linux/macOS
./scripts/build_and_preview.sh 3
```

### Development workflow:
```bash
# Windows PowerShell
.\scripts\dev.ps1 go run ./cmd/planner

# Linux/macOS
./scripts/dev.sh go run ./cmd/planner
```

### Project cleanup:
```bash
# Windows PowerShell
.\scripts\cleanup_and_organize.ps1

# Linux/macOS
./scripts/cleanup_and_organize.sh
```

### Using Make (works on all platforms with make installed):
```bash
make build          # Build PDF
make dev           # Development workflow
make clean         # Clean artifacts
make status        # Show project status
```

## Troubleshooting

### Common Issues:

1. **"pdflatex not found"**
   - Install XeLaTeX/TeX distribution for your platform

2. **"python/pdf2image not working"**
   - Install Python 3 and pdf2image package
   - Install Poppler system library

3. **Permission denied on scripts**
   - Run `chmod +x scripts/*.sh` on Unix systems
   - Or use `icacls scripts\*.ps1 /grant Everyone:RX` on Windows

4. **Path issues**
   - Scripts automatically detect platform and use correct path separators
   - If issues persist, check that you're in the project root directory
