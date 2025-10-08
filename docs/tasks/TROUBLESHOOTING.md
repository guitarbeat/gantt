# Troubleshooting Guide

**Last Updated:** October 7, 2025

This guide covers common issues and their solutions when using the PhD Dissertation Planner.

---

## Table of Contents

1. [Installation Issues](#installation-issues)
2. [Build Errors](#build-errors)
3. [LaTeX Compilation Errors](#latex-compilation-errors)
4. [CSV Format Issues](#csv-format-issues)
5. [Preview System Issues](#preview-system-issues)
6. [Platform-Specific Issues](#platform-specific-issues)
7. [Performance Issues](#performance-issues)

---

## Installation Issues

### Go Not Found

**Error:**
```
'go' is not recognized as an internal or external command
```

**Solution:**
1. Install Go from https://golang.org/dl/
2. Verify installation: `go version`
3. Restart your terminal/IDE

### LaTeX Not Found

**Error:**
```
pdflatex: command not found
```

**Solution:**

**Windows:**
```powershell
# Install MiKTeX
winget install MiKTeX.MiKTeX
```

**Mac:**
```bash
# Install MacTeX
brew install --cask mactex
```

**Linux:**
```bash
# Install TeX Live
sudo apt-get install texlive-xetex texlive-fonts-recommended
```

### Missing LaTeX Packages

**Error:**
```
! LaTeX Error: File `xcolor.sty' not found.
```

**Solution:**

**MiKTeX (Windows):**
- Packages install automatically on first use
- Or use: `mpm --install=xcolor`

**TeX Live (Mac/Linux):**
```bash
sudo tlmgr install xcolor
```

---

## Build Errors

### Module Not Found

**Error:**
```
package github.com/user/repo: cannot find module
```

**Solution:**
```bash
go mod download
go mod tidy
```

### Build Fails with "undefined"

**Error:**
```
undefined: SomeFunction
```

**Solution:**
1. Check imports in the file
2. Run `go mod tidy`
3. Rebuild: `go build ./cmd/planner`

### Permission Denied

**Error:**
```
permission denied: ./plannergen
```

**Solution:**

**Mac/Linux:**
```bash
chmod +x plannergen
```

**Windows:**
- Run PowerShell as Administrator
- Or check antivirus settings

---

## LaTeX Compilation Errors

### PDF Generation Fails

**Error:**
```
Error: PDF generation failed
```

**Diagnosis Steps:**
1. Check if LaTeX is installed: `pdflatex --version`
2. Look at the generated `.tex` file in `generated/`
3. Try compiling manually: `pdflatex generated/planner.tex`
4. Check the `.log` file for detailed errors

### Special Characters in Task Names

**Error:**
```
! Missing $ inserted.
```

**Cause:** Special LaTeX characters in CSV data (%, $, &, #, _, {, })

**Solution:**
Escape special characters in your CSV:
- `%` → `\%`
- `$` → `\$`
- `&` → `\&`
- `#` → `\#`
- `_` → `\_`
- `{` → `\{`
- `}` → `\}`

Or use the `--escape` flag (if implemented)

### Font Not Found

**Error:**
```
! Font \TU/lmr/m/n/10=lmroman10-regular at 10.0pt not loadable
```

**Solution:**
```bash
# Update font cache
fc-cache -fv

# Or install Latin Modern fonts
sudo apt-get install lmodern  # Linux
brew install --cask font-latin-modern  # Mac
```

### Week Column Width Issue

**Known Issue:** Week columns may appear too wide on some systems.

**Current Status:** Under investigation (5 attempts made)

**Workaround:**
1. Use the compact preset: `--preset compact`
2. Or manually edit `generated/planner.tex` and adjust column widths

**Tracking:** See REPOSITORY_IMPROVEMENTS.md for detailed investigation notes

---

## CSV Format Issues

### Date Parse Error

**Error:**
```
Error parsing date: invalid format
```

**Solution:**
Dates must be in `YYYY-MM-DD` format:
- ✅ Correct: `2025-03-15`
- ❌ Wrong: `03/15/2025`, `15-03-2025`, `March 15, 2025`

### Missing Required Columns

**Error:**
```
Error: missing required column: Task
```

**Solution:**
Your CSV must have these columns:
- `Task` (required)
- `Start Date` (required)
- `End Date` (required)
- `Phase` (optional)
- `Category` (optional)
- `Priority` (optional)

**Example:**
```csv
Task,Start Date,End Date,Phase,Category
Literature Review,2025-01-01,2025-03-31,Phase 1,Research
```

### Empty Rows

**Error:**
```
Error: invalid task on line 15
```

**Solution:**
Remove empty rows from your CSV file. Each row must have at least Task, Start Date, and End Date.

### Encoding Issues

**Error:**
```
Error: invalid UTF-8 encoding
```

**Solution:**
Save your CSV as UTF-8:
- Excel: Save As → CSV UTF-8
- Google Sheets: Download → CSV
- Text Editor: Save with UTF-8 encoding

---

## Preview System Issues

### Python Not Found

**Error:**
```
'python' is not recognized as an internal or external command
```

**Solution:**
1. Install Python 3.8+ from https://python.org
2. Verify: `python --version`
3. On some systems, use `python3` instead of `python`

### pdf2image Module Not Found

**Error:**
```
ModuleNotFoundError: No module named 'pdf2image'
```

**Solution:**
```bash
pip install pdf2image Pillow
```

### Poppler Not Found

**Error:**
```
PDFInfoNotInstalledError: Unable to get page count. Is poppler installed?
```

**Solution:**

**Windows:**
```powershell
# Download from: https://github.com/oschwartz10612/poppler-windows/releases/
# Extract and add to PATH
```

**Mac:**
```bash
brew install poppler
```

**Linux:**
```bash
sudo apt-get install poppler-utils
```

### Preview Images Not Generated

**Diagnosis:**
1. Check if PDF was generated: `ls generated/*.pdf`
2. Check if Python script exists: `ls scripts/pdf_to_images.py`
3. Run manually: `python scripts/pdf_to_images.py generated/planner.pdf generated/preview`
4. Check for errors in output

---

## Platform-Specific Issues

### Windows: PowerShell Execution Policy

**Error:**
```
cannot be loaded because running scripts is disabled
```

**Solution:**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Windows: Path Too Long

**Error:**
```
The specified path, file name, or both are too long
```

**Solution:**
1. Enable long paths in Windows:
   - Run `gpedit.msc`
   - Navigate to: Computer Configuration → Administrative Templates → System → Filesystem
   - Enable "Enable Win32 long paths"
2. Or move project to shorter path (e.g., `C:\gantt`)

### Mac: Gatekeeper Blocking Binary

**Error:**
```
"plannergen" cannot be opened because it is from an unidentified developer
```

**Solution:**
```bash
xattr -d com.apple.quarantine plannergen
```

Or: System Preferences → Security & Privacy → Allow

### Linux: Missing Dependencies

**Error:**
```
error while loading shared libraries
```

**Solution:**
```bash
sudo apt-get update
sudo apt-get install build-essential
```

---

## Performance Issues

### Slow PDF Generation

**Symptoms:** Build takes >2 minutes

**Solutions:**
1. **Reduce task count:** Split into multiple planners
2. **Use SSD:** Move project to SSD if on HDD
3. **Close other apps:** Free up system resources
4. **Disable preview:** Skip image generation if not needed

### High Memory Usage

**Symptoms:** System becomes slow during build

**Solutions:**
1. **Increase swap space** (Linux/Mac)
2. **Close browser tabs** and other memory-intensive apps
3. **Process in batches:** Split large CSV files

### LaTeX Compilation Hangs

**Symptoms:** Build appears stuck

**Solutions:**
1. **Kill process:** `Ctrl+C` or Task Manager
2. **Check for interactive prompts:** LaTeX may be waiting for input
3. **Run in batch mode:** Add `-interaction=nonstopmode` to pdflatex command
4. **Clear cache:** Delete `generated/` folder and rebuild

---

## Getting More Help

### Check Logs

1. **Build logs:** Check terminal output
2. **LaTeX logs:** Check `generated/*.log`
3. **Go logs:** Run with verbose flag: `go run ./cmd/planner -v`

### Debug Mode

Enable debug output:
```bash
export DEBUG=1
go run ./cmd/planner
```

### Report an Issue

If you can't resolve the issue:

1. **Check existing issues:** https://github.com/yourusername/gantt/issues
2. **Create new issue** with:
   - Error message (full text)
   - Steps to reproduce
   - System info: OS, Go version, LaTeX version
   - Sample CSV (if relevant)
   - Generated `.tex` file (if relevant)

### Community Support

- **GitHub Discussions:** Ask questions and share tips
- **Stack Overflow:** Tag with `latex`, `go`, `pdf-generation`

---

## Quick Diagnostic Checklist

Run through this checklist before reporting issues:

- [ ] Go is installed: `go version`
- [ ] LaTeX is installed: `pdflatex --version`
- [ ] Dependencies are current: `go mod download`
- [ ] CSV format is correct (YYYY-MM-DD dates)
- [ ] No special characters in task names
- [ ] Generated folder exists and is writable
- [ ] Enough disk space (>100MB free)
- [ ] No antivirus blocking the build
- [ ] Latest version of code: `git pull`

---

## Common Error Messages Reference

| Error | Likely Cause | Quick Fix |
|-------|-------------|-----------|
| `command not found` | Missing installation | Install the tool |
| `permission denied` | File permissions | `chmod +x` or run as admin |
| `module not found` | Missing Go dependency | `go mod download` |
| `invalid date format` | Wrong date format | Use YYYY-MM-DD |
| `! LaTeX Error` | LaTeX syntax issue | Check .log file |
| `PDF generation failed` | LaTeX not installed | Install LaTeX |
| `cannot find package` | Import path issue | `go mod tidy` |

---

**Still stuck?** Open an issue with the "help wanted" label and we'll assist you!
