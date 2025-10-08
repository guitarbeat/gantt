# PDF Preview Images - Quick Setup

## ✅ What's Working Now

The system can now generate PNG preview images from PDFs automatically!

## 🚀 Quick Start

### 1. Install Requirements (One-time setup)

```powershell
# Install Python packages
pip install pdf2image pillow

# Install poppler (PDF rendering library)
choco install poppler
```

### 2. Build and Generate Previews

```powershell
# Build PDF and create 3 preview images
.\scripts\build_and_preview.ps1 -Pages 3
```

### 3. Share Images

Preview images are created in `generated/preview/`:
- `page-01.png`
- `page-02.png`
- `page-03.png`

**Drag these files into the chat to share them!**

## 📋 What You Need

| Tool | Purpose | Install Command |
|------|---------|----------------|
| Python 3 | Run conversion script | Already installed ✓ |
| pdf2image | Python PDF library | `pip install pdf2image` |
| Pillow | Image processing | `pip install pillow` |
| Poppler | PDF rendering | `choco install poppler` |

## 🔧 Installation Details

### Python Packages (Already Done ✓)
```powershell
pip install pdf2image pillow
```

### Poppler (Still Needed)

**Option 1: Chocolatey (Recommended)**
```powershell
choco install poppler
```

**Option 2: Manual Download**
1. Download: https://github.com/oschwartz10612/poppler-windows/releases/
2. Extract to: `C:\Program Files\poppler`
3. Add to PATH: `C:\Program Files\poppler\Library\bin`

## ✅ Verify Installation

```powershell
# Check if everything is installed
python --version                    # Should show Python 3.x
python -c "import pdf2image"        # Should have no errors
pdfinfo -v                          # Should show poppler version
```

## 📖 Documentation

- **Quick Guide**: `scripts/README_PREVIEW.md`
- **Detailed Setup**: `scripts/SETUP_PREVIEW.md`
- **Python Script**: `scripts/pdf_to_images.py`
- **PowerShell Script**: `scripts/build_and_preview.ps1`

## 🎯 Current Status

- ✅ Python installed (3.13.0)
- ✅ pdf2image installed
- ✅ Pillow installed
- ✅ Python script created
- ✅ PowerShell script updated
- ⏳ **Poppler needs to be installed**

## 🚦 Next Steps

1. **Install poppler**: `choco install poppler`
2. **Restart PowerShell** (to load new PATH)
3. **Run the script**: `.\scripts\build_and_preview.ps1 -Pages 3`
4. **Drag images** from `generated/preview/` into chat

## 💡 Alternative (No Installation)

If you don't want to install poppler, you can manually create screenshots:
1. Open the PDF
2. Use Windows Snipping Tool (Win+Shift+S)
3. Save as PNG in `generated/preview/`

---

**Ready to test?** Run: `.\scripts\build_and_preview.ps1 -Pages 3`
