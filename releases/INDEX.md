# PhD Dissertation Planner - Releases

This directory contains timestamped releases for version tracking and progression evaluation.

## Directory Structure

```
releases/
├── INDEX.md                 # This file - release history
├── v5.0/                   # Version 5.0 releases
│   ├── v5.0_YYYYMMDD_HHMMSS.pdf
│   ├── v5.0_YYYYMMDD_HHMMSS.tex
│   ├── v5.0_YYYYMMDD_HHMMSS_source.csv
│   ├── v5.0_YYYYMMDD_HHMMSS_metadata.json
│   └── MANIFEST.txt
└── v5.1/                   # Version 5.1 releases
    ├── v5.1_YYYYMMDD_HHMMSS.pdf
    ├── v5.1_YYYYMMDD_HHMMSS.tex
    ├── v5.1_YYYYMMDD_HHMMSS_source.csv
    ├── v5.1_YYYYMMDD_HHMMSS_metadata.json
    └── MANIFEST.txt
```

## How to Use

### Building a Release

```bash
# Build latest version (auto-detected)
./scripts/build_release.sh

# Build specific version
./scripts/build_release.sh --version v5.1

# Build with custom name
./scripts/build_release.sh --version v5.1 --name "Final_Review"

# Build with custom CSV
./scripts/build_release.sh --csv input_data/custom.csv --version v5.2
```

### Release Files

Each release includes:

1. **PDF** (`vX.X_YYYYMMDD_HHMMSS.pdf`)
   - Final compiled planner document
   - Ready for printing or distribution

2. **LaTeX** (`vX.X_YYYYMMDD_HHMMSS.tex`)
   - Source LaTeX file
   - Can be manually edited or recompiled

3. **CSV** (`vX.X_YYYYMMDD_HHMMSS_source.csv`)
   - Original CSV data used for generation
   - Allows exact reproduction of release

4. **Metadata** (`vX.X_YYYYMMDD_HHMMSS_metadata.json`)
   - Build information (date, version, environment)
   - Useful for troubleshooting and tracking

5. **Manifest** (`MANIFEST.txt`)
   - Human-readable release summary
   - Lists all files and their sizes

## Comparing Releases

### View Release History
```bash
cat releases/INDEX.md
```

### Compare Two Releases
```bash
# Compare PDFs visually
open releases/v5.0/v5.0_20251003_120000.pdf
open releases/v5.1/v5.1_20251003_130000.pdf

# Compare CSV changes
diff releases/v5.0/v5.0_20251003_120000_source.csv \
     releases/v5.1/v5.1_20251003_130000_source.csv

# Compare LaTeX output
diff releases/v5.0/v5.0_20251003_120000.tex \
     releases/v5.1/v5.1_20251003_130000.tex
```

### View Release Metadata
```bash
# Pretty-print JSON metadata
cat releases/v5.1/v5.1_20251003_130000_metadata.json | python3 -m json.tool

# Or use jq (if installed)
jq . releases/v5.1/v5.1_20251003_130000_metadata.json
```

## Progression Tracking

Use timestamped releases to:

1. **Track improvements** over time
2. **Compare different versions** side-by-side
3. **Roll back** to previous versions if needed
4. **Document evolution** of your timeline
5. **Archive milestones** for historical reference

### Example Workflow

```bash
# Monday: Initial v5.1 release
./scripts/build_release.sh --version v5.1

# Wednesday: After task adjustments
./scripts/build_release.sh --version v5.1 --name "Revised"

# Friday: Final version for committee
./scripts/build_release.sh --version v5.1 --name "Committee_Review"

# Compare Wednesday vs Friday
diff releases/v5.1/v5.1_*_Revised_source.csv \
     releases/v5.1/v5.1_*_Committee_Review_source.csv
```

## Best Practices

1. **Regular Builds**: Create releases regularly (weekly/monthly)
2. **Meaningful Names**: Use custom names for important milestones
3. **Keep Manifests**: Don't delete old releases - they document progression
4. **Version Control**: Commit `INDEX.md` to track release history
5. **Backup**: Periodically backup the entire `releases/` directory

## Cleanup

To remove old releases while keeping important ones:

```bash
# Keep only last 5 releases per version
cd releases/v5.1
ls -t v5.1_*.pdf | tail -n +6 | xargs rm -f

# Remove all releases older than 30 days
find releases/ -name "v*.pdf" -mtime +30 -delete
```

## Troubleshooting

### Release build fails
```bash
# Check build logs
cat generated/monthly_calendar.log

# Try manual build first
make -f scripts/Makefile clean-build

# Then create release
./scripts/build_release.sh
```

### Missing XeLaTeX
```bash
# Install on macOS
brew install --cask mactex

# Install on Linux
sudo apt-get install texlive-xetex texlive-latex-extra

# Build without PDF
./scripts/build_release.sh --skip-pdf
```

## Release History

<!-- Releases are automatically appended below by build_release.sh -->

### v5.1_20251003_153351

- **Date:** 2025-10-03 15:33:51
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/`

### v5.1_No_Continuing_Tasks_20251003_155125

- **Date:** 2025-10-03 15:51:25
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/`

### 20251003_155402_Test_New_Structure

- **Date:** 2025-10-03 15:54:02
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251003_155402_Test_New_Structure/`

### 20251003_155558_Final_Clean_Version

- **Date:** 2025-10-03 15:55:58
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251003_155558_Final_Clean_Version/`

### 20251006_150623_release

- **Date:** 2025-10-06 15:06:23
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_150623_release/`

### 20251006_151224_release

- **Date:** 2025-10-06 15:12:24
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_151224_release/`

### 20251006_152516_release

- **Date:** 2025-10-06 15:25:16
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_152516_release/`

### 20251006_154044_release

- **Date:** 2025-10-06 15:40:44
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_154044_release/`

### 20251006_154931_release

- **Date:** 2025-10-06 15:49:31
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_154931_release/`

### 20251006_162130_release

- **Date:** 2025-10-06 16:21:30
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_162130_release/`

### 20251006_162818_release

- **Date:** 2025-10-06 16:28:18
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_162818_release/`

### 20251006_163721_release

- **Date:** 2025-10-06 16:37:21
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_163721_release/`

### 20251006_164102_release

- **Date:** 2025-10-06 16:41:02
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_164102_release/`

### 20251006_164354_release

- **Date:** 2025-10-06 16:43:54
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_164354_release/`

### 20251006_164725_release

- **Date:** 2025-10-06 16:47:25
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_164725_release/`

### 20251006_171853_release

- **Date:** 2025-10-06 17:18:53
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_171853_release/`

### 20251006_172229_release

- **Date:** 2025-10-06 17:22:29
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_172229_release/`

### 20251006_172500_release

- **Date:** 2025-10-06 17:25:00
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_172500_release/`

### 20251006_172644_release

- **Date:** 2025-10-06 17:26:44
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_172644_release/`

### 20251006_172811_release

- **Date:** 2025-10-06 17:28:11
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_172811_release/`

### 20251006_172908_release

- **Date:** 2025-10-06 17:29:08
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/v5.1/20251006_172908_release/`

