# PhD Dissertation Planner - Releases

This directory contains timestamped releases for version tracking and progression evaluation.

## Directory Structure

```
releases/
├── INDEX.md                           # This file - release history
├── YYYYMMDD_HHMMSS_NAME/              # Individual release directories
│   ├── planner.pdf                    # Final compiled planner document
│   ├── planner.tex                    # Source LaTeX file
│   ├── source.csv                     # Original CSV data used for generation
│   ├── metadata.json                  # Build information (date, version, environment)
│   └── README.md                      # Human-readable release summary
└── ...
```

## How to Use

### Building a Release

```bash
# Build latest version (auto-detected)
./scripts/build_release.sh

# Build with custom name
./scripts/build_release.sh --name "Final_Review"

# Build with custom CSV
./scripts/build_release.sh --csv input_data/custom.csv --name "Custom_Data"

# Build with preset
./scripts/build_release.sh --preset compact --name "Compact_View"
```

### Release Files

Each release includes:

1. **PDF** (`planner.pdf`)
   - Final compiled planner document
   - Ready for printing or distribution

2. **LaTeX** (`planner.tex`)
   - Source LaTeX file
   - Can be manually edited or recompiled

3. **CSV** (`source.csv`)
   - Original CSV data used for generation
   - Allows exact reproduction of release

4. **Metadata** (`metadata.json`)
   - Build information (date, version, environment)
   - Useful for troubleshooting and tracking

5. **README** (`README.md`)
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
open releases/20251006_190535_SubPhase_Headers_Test/planner.pdf
open releases/20251006_184526_QuickWins_Implementation/planner.pdf

# Compare CSV changes
diff releases/20251006_190535_SubPhase_Headers_Test/source.csv \
     releases/20251006_184526_QuickWins_Implementation/source.csv

# Compare LaTeX output
diff releases/20251006_190535_SubPhase_Headers_Test/planner.tex \
     releases/20251006_184526_QuickWins_Implementation/planner.tex
```

### View Release Metadata
```bash
# Pretty-print JSON metadata
cat releases/20251006_190535_SubPhase_Headers_Test/metadata.json | python3 -m json.tool

# Or use jq (if installed)
jq . releases/20251006_190535_SubPhase_Headers_Test/metadata.json
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
# Monday: Initial release
./scripts/build_release.sh --name "Initial"

# Wednesday: After task adjustments
./scripts/build_release.sh --name "Revised"

# Friday: Final version for committee
./scripts/build_release.sh --name "Committee_Review"

# Compare Wednesday vs Friday
diff releases/*_Revised/source.csv \
     releases/*_Committee_Review/source.csv
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
# Keep only last 10 releases
cd releases
ls -t */planner.pdf | tail -n +11 | xargs -I {} rm -rf {}

# Remove all releases older than 30 days
find releases/ -name "planner.pdf" -mtime +30 -exec dirname {} \; | xargs rm -rf
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

### 20251006_190535_SubPhase_Headers_Test

- **Date:** 2025-10-06 19:05:35
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_190535_SubPhase_Headers_Test/`
### 20251006_194103_Flat_Structure_Test

- **Date:** 2025-10-06 19:41:03
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_194103_Flat_Structure_Test/`

### 20251006_194228_Compact_Task_Index

- **Date:** 2025-10-06 19:42:28
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_194228_Compact_Task_Index/`

### 20251006_194350_Compact_No_M_Indicator

- **Date:** 2025-10-06 19:43:50
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_194350_Compact_No_M_Indicator/`

### 20251006_194532_Bidirectional_Hyperlinks

- **Date:** 2025-10-06 19:45:32
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_194532_Bidirectional_Hyperlinks/`

### 20251006_194740_Fixed_Hyperlinks_Chronological

- **Date:** 2025-10-06 19:47:40
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_194740_Fixed_Hyperlinks_Chronological/`

### 20251006_195040_Fixed_Bidirectional_Hyperlinks

- **Date:** 2025-10-06 19:50:40
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_195040_Fixed_Bidirectional_Hyperlinks/`

### 20251006_195213_No_Color_Legend

- **Date:** 2025-10-06 19:52:13
- **Version:** v5.1
- **CSV:** research_timeline_v5.1_comprehensive.csv
- **Location:** `releases/20251006_195213_No_Color_Legend/`

