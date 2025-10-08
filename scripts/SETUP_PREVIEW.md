# Setup Guide: PDF Preview Generation

This guide will help you set up everything needed to generate PNG preview images from PDFs.

## Quick Setup (Windows)

### Step 1: Install Python Packages
```powershell
pip install pdf2image pillow
```

### Step 2: Install Poppler

**Option A: Using Chocolatey (Easiest)**
```powershell
# Install Chocolatey if you don't have it
# See: https://chocolatey.org/install

# Install poppler
choco install poppler -y
```

**Option B: Manual Installation**
1. Download poppler for Windows:
   - Go to: https://github.com/oschwartz10612/poppler-windows/releases/
   - Download the latest `Release-XX.XX.X-0.zip`
   - Extract to `C:\Program Files\poppler`

2. Add to PATH:
   ```powershell
   # Temporary (current session only)
   $env:PATH += ";C:\Program Files\poppler\Library\bin"
   
   # Permanent (requires admin)
   [Environment]::SetEnvironmentVariable(
       "Path",
       [Environment]::GetEnvironmentVariable("Path", "Machine") + ";C:\Program Files\poppler\Library\bin",
       "Machine"
   )
   ```

3. Restart PowerShell to apply changes

### Step 3: Verify Installation
```powershell
# Check Python
python --version
# Should show: Python 3.x.x

# Check pdf2image
python -c "import pdf2image; print('✓ pdf2image installed')"
# Should show: ✓ pdf2image installed

# Check poppler
pdfinfo -v
# Should show: pdfinfo version x.x.x
```

## Usage

Once everything is installed:

```powershell
# Build PDF and generate 3 preview images
.\scripts\build_and_preview.ps1 -Pages 3
```

Preview images will be created in `generated/preview/`:
- `page-01.png`
- `page-02.png`
- `page-03.png`

## What Each Tool Does

### Python
- **Purpose**: Runs the conversion script
- **Already installed**: Usually comes with Windows or can be installed from python.org

### pdf2image (Python package)
- **Purpose**: Python library that interfaces with poppler
- **Install**: `pip install pdf2image`

### Pillow (Python package)
- **Purpose**: Image processing library (saves PNG files)
- **Install**: `pip install pillow`

### Poppler
- **Purpose**: PDF rendering library (does the actual PDF→image conversion)
- **Install**: `choco install poppler` or manual download
- **Note**: This is the most important component!

## Alternative: Manual Conversion

If you can't install poppler, you can manually create preview images:

1. Open `generated/monthly_calendar.pdf` in a PDF viewer
2. Use Windows Snipping Tool (Win+Shift+S)
3. Capture the first few pages
4. Save as PNG files in `generated/preview/`
5. Name them: `page-01.png`, `page-02.png`, etc.

## Troubleshooting

### "poppler not found" or "Unable to get page count"

**Cause**: Poppler is not installed or not in PATH

**Fix**:
1. Install poppler: `choco install poppler`
2. Or add to PATH: `$env:PATH += ";C:\Program Files\poppler\Library\bin"`
3. Restart PowerShell
4. Verify: `pdfinfo -v`

### "ModuleNotFoundError: No module named 'pdf2image'"

**Cause**: Python package not installed

**Fix**:
```powershell
pip install pdf2image pillow
```

### "Python not found"

**Cause**: Python not installed or not in PATH

**Fix**:
1. Download from: https://www.python.org/downloads/
2. During installation, check "Add Python to PATH"
3. Restart PowerShell

### Images are blurry

**Fix**: Increase DPI in the script:
```powershell
python scripts\pdf_to_images.py "generated\monthly_calendar.pdf" "generated\preview" 3 300
```
(300 DPI instead of 150 DPI)

### Script runs but no images appear

**Check**:
1. Does `generated/preview/` directory exist?
2. Are there any error messages?
3. Try running Python script directly:
   ```powershell
   python scripts\pdf_to_images.py "generated\monthly_calendar.pdf" "generated\preview" 3 150
   ```

## Testing Your Setup

Run this test to verify everything works:

```powershell
# Test 1: Python
python --version

# Test 2: Python packages
python -c "import pdf2image, PIL; print('✓ All packages installed')"

# Test 3: Poppler
pdfinfo -v

# Test 4: Full conversion
python scripts\pdf_to_images.py "generated\monthly_calendar.pdf" "generated\preview" 1 150

# Test 5: Check output
Get-ChildItem "generated\preview\*.png"
```

If all tests pass, you're ready to go!

## Support

If you're still having issues:
1. Check the error messages carefully
2. Make sure you restarted PowerShell after installing poppler
3. Verify poppler is in your PATH: `$env:PATH -split ';' | Select-String poppler`
4. Try the manual conversion method as a fallback
