# Build with Preview Images

## Quick Start

```powershell
# Build PDF and generate 3 preview images
.\scripts\build_and_preview.ps1 -Pages 3
```

## What It Does

1. **Builds** the Go binary
2. **Generates** LaTeX from CSV
3. **Compiles** PDF with pdflatex
4. **Creates** PNG preview images (if Ghostscript is installed)

## Preview Images

Preview images are saved to `generated/preview/`:
- `page-1.png` - First page
- `page-2.png` - Second page  
- `page-3.png` - Third page

## Requirements

### Required
- **Go** (for building)
- **pdflatex** (for PDF compilation)
- **Python 3** (for preview images)

### Python Packages
```powershell
# Install required Python packages
pip install pdf2image pillow
```

### Poppler (PDF rendering library)
The Python script needs poppler to render PDFs. Install it:

**Option 1: Chocolatey (Recommended)**
```powershell
choco install poppler
```

**Option 2: Manual Download**
1. Download from: https://github.com/oschwartz10612/poppler-windows/releases/
2. Extract to `C:\Program Files\poppler`
3. Add `C:\Program Files\poppler\Library\bin` to your PATH

**Option 3: Scoop**
```powershell
scoop install poppler
```

### Verify Installation
```powershell
# Check Python
python --version

# Check if pdf2image is installed
python -c "import pdf2image; print('pdf2image OK')"

# Check if poppler is accessible
pdfinfo -v
```

## Usage Examples

```powershell
# Generate 5 preview pages
.\scripts\build_and_preview.ps1 -Pages 5

# Default (3 pages)
.\scripts\build_and_preview.ps1
```

## Sharing Preview Images

After running the script, drag the PNG files from `generated/preview/` into the chat to share them.

The script will show you the full paths:
```
Drag these images into the chat to share them:
  C:\path\to\generated\preview\page-1.png
  C:\path\to\generated\preview\page-2.png
  C:\path\to\generated\preview\page-3.png
```

## Troubleshooting

### "Unable to get page count. Is poppler installed and in PATH?"
This means poppler is not installed or not in your PATH.

**Fix:**
```powershell
# Install with Chocolatey
choco install poppler

# Or add to PATH manually
$env:PATH += ";C:\Program Files\poppler\Library\bin"
```

### "ModuleNotFoundError: No module named 'pdf2image'"
Install the Python package:
```powershell
pip install pdf2image pillow
```

### "PDF not created"
- Check `generated/monthly_calendar.log` for LaTeX errors
- Ensure pdflatex is installed and in PATH
- Run `pdflatex --version` to verify

### Images not appearing
1. Check if `generated/preview/` directory exists
2. Look for error messages in the script output
3. Try running the Python script directly:
   ```powershell
   python scripts\pdf_to_images.py "generated\monthly_calendar.pdf" "generated\preview" 3 150
   ```

## Integration with Release Process

To create a release with preview images:

```powershell
# 1. Build with previews
.\scripts\build_and_preview.ps1 -Pages 3

# 2. Create release directory
$releaseDir = "releases\v5.1\$(Get-Date -Format 'yyyyMMdd_HHmmss')_MyRelease"
New-Item -ItemType Directory -Path $releaseDir -Force

# 3. Copy files
Copy-Item "generated\monthly_calendar.pdf" "$releaseDir\planner.pdf"
Copy-Item "generated\preview\*" "$releaseDir\"
```

Now your release will include both the PDF and preview images!
