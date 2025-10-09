# 🧭 Quick Navigation Guide

## 🎯 Most Important Directories

### 📁 **Start Here**
- `src/` - **Main source code** (Go application)
- `input_data/` - **Your CSV data files**
- `generated/` - **Generated PDFs and outputs**
- `releases/` - **Versioned releases** (organized by date)

### 📚 **Documentation**
- `docs/tasks/` - **How-to guides** (SETUP.md, USER_GUIDE.md)
- `docs/fyi/` - **Reference information**
- `PROJECT_STRUCTURE.md` - **Complete structure overview**

### ⚙️ **Configuration & Scripts**
- `configs/` - **YAML configuration files**
- `scripts/` - **Build and utility scripts**

## 🚀 Quick Commands

```bash
# Build a PDF from your data
./scripts/quick_build.sh

# Create a timestamped release
./scripts/build_release.sh

# Clean up and organize files
./scripts/cleanup_and_organize.sh

# View project status
./scripts/cleanup_and_organize.sh --status
```

## 📋 File Organization Rules

| File Type           | Location          | Purpose                     |
| ------------------- | ----------------- | --------------------------- |
| **Go source code**  | `src/`            | Main application            |
| **CSV data**        | `input_data/`     | Your research timeline data |
| **Generated PDFs**  | `generated/pdfs/` | Output PDFs                 |
| **LaTeX files**     | `generated/tex/`  | LaTeX source                |
| **Logs**            | `generated/logs/` | Build logs                  |
| **Temporary files** | `.temp/`          | Auto-cleaned                |
| **Documentation**   | `docs/`           | Guides and references       |
| **Releases**        | `releases/`       | Versioned outputs           |

## 🧹 Maintenance

The project now has automatic cleanup! Run this whenever things get messy:

```bash
./scripts/cleanup_and_organize.sh
```

This will:
- ✅ Move scattered files to proper locations
- ✅ Organize generated files by type
- ✅ Clean up temporary files
- ✅ Update documentation structure

## 🎯 What's Where

### Your Data
- **CSV files**: `input_data/research_timeline_v5.1_comprehensive.csv`
- **Generated PDFs**: `generated/pdfs/planner.pdf`
- **Latest release**: `releases/20251008_164428_Task-Index-Final/`

### Source Code
- **Main app**: `src/app/generator.go`
- **Calendar logic**: `src/calendar/calendar.go`
- **Configuration**: `src/core/config.go`
- **Templates**: `src/shared/templates/`

### Documentation
- **Setup guide**: `docs/tasks/SETUP.md`
- **User guide**: `docs/tasks/USER_GUIDE.md`
- **Developer guide**: `docs/tasks/DEVELOPER_GUIDE.md`

---

**💡 Pro Tip**: The project is now much more organized! Everything has its place, and the cleanup script will keep it that way.
