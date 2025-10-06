# Releases Directory

This directory contains timestamped releases of your PhD dissertation planner for progression tracking.

## Quick Start

```bash
# Create a new release (auto-detects latest version)
./scripts/build_release.sh

# Create a named release for important milestones
./scripts/build_release.sh --version v5.1 --name "Committee_Submission"
```

## Directory Structure

```
releases/
├── README.md          # This file
├── INDEX.md           # Complete release history (auto-updated)
├── v5.0/             # Version 5.0 releases
│   └── v5.0_YYYYMMDD_HHMMSS.*
└── v5.1/             # Version 5.1 releases
    └── v5.1_YYYYMMDD_HHMMSS.*
```

## Each Release Contains

- **PDF**: Compiled planner document (~430 KB)
- **TEX**: LaTeX source file (~7.5 KB)  
- **CSV**: Original data file (~17 KB)
- **JSON**: Build metadata (<1 KB)
- **MANIFEST.txt**: Human-readable summary

## Example Usage

### Track Weekly Progress
```bash
# Every Friday
./scripts/build_release.sh --name "Week_$(date +%U)"
```

### Before Important Events
```bash
# Before advisor meeting
./scripts/build_release.sh --name "Advisor_Meeting_2025-10-15"

# Committee submission
./scripts/build_release.sh --name "Committee_Final"
```

### Compare Versions
```bash
# View what changed
diff releases/v5.1/v5.1_20251003_100000_source.csv \
     releases/v5.1/v5.1_20251003_150000_source.csv
```

## More Information

- **Full Documentation**: See `../RELEASE_SYSTEM.md`
- **Release History**: See `INDEX.md`
- **Build Script**: See `../scripts/build_release.sh`

---

**Tip**: Each release is self-contained and reproducible. Keep them for historical reference!
