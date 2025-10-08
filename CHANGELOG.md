# Changelog

All notable changes to the PhD Dissertation Planner will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [5.1.1] - 2025-10-08

### Added

- **Complete Professional Task Index**: Working PDF with clickable task index as the first page
- **Task Index Features**:
  - Tasks organized by research phases with clear headers
  - Clickable hyperlinks to jump directly to task locations in the timeline
  - Milestone highlighting with star symbols (★)
  - Simple, stable LaTeX formatting that compiles reliably
  - Phase-by-phase task organization

### Fixed

- **Week Label Positioning**: Restored original week label positioning as first table column
- **Table Structure**: Fixed calendar table layout to use proper 8-column structure
- **LaTeX Compilation**: Resolved complex LaTeX formatting issues that prevented PDF generation
- **Task Index Stability**: Simplified LaTeX generation to ensure reliable compilation

### Changed

- **Release System**: Enhanced build_release.sh to properly handle multi-file LaTeX projects with working PDF output
- **Task Index Generation**: Streamlined LaTeX formatting for stability and reliability

## [5.1.0] - 2025-10-03

### Added

- Timestamped release system with organized directory structure
- Per-release README.md and metadata.json files
- Version-specific INDEX.md files for tracking releases
- Automated release build script (`scripts/build_release.sh`)
- `.build_temp/` directory for temporary build artifacts

### Changed

- **Task Distribution**: Improved v5.1 timeline with better measurability (89% → 96%)
  - Removed 9 non-measurable/administrative tasks
  - Split 4 long tasks into 12 smaller milestones
  - Maintained 107 total tasks with better distribution
- **Task Rendering**: Fixed multi-day task spanning
  - Tasks now show text only on start day
  - Continuing days show only colored bar (no text duplication)
  - Much cleaner calendar appearance
- **Release Structure**: Organized releases in `releases/VERSION/TIMESTAMP_NAME/` format
  - Each release is self-contained
  - Simplified filenames (planner.pdf, planner.tex, source.csv)
  - Better navigation and comparison capabilities

### Removed

- `generated/` directory (replaced by release system)
- Old flat release structure with long filenames
- Non-measurable maintenance tasks from timeline:
  - T2.36: Maintain Automated Backups (608 days)
  - T2.37: Maintain Surgical Training (244 days)
  - T4.18: Maintain Lab Responsibilities (721 days)
- Administrative tasks without clear deliverables:
  - T4.11-T4.14: Semester registration maintenance
  - T4.17: SPIE Student Chapter Activities

### Fixed

- Task rendering duplication on multi-day spans
- Release file organization and naming
- Git configuration for new release structure

## [Unreleased]

### New Features

- `scripts/` directory with build and setup automation scripts
- `assets/` directory for static resources (PDFs, documents)
- `CONTRIBUTING.md` with development guidelines
- Standard Go project structure following community best practices

### Project Restructuring

- Reorganized project structure:
  - `cmd/planner/` for main application entry point
  - `pkg/templates/` for reusable template components
  - `configs/` for configuration files (renamed from `config/`)
  - `data/` for input data files (renamed from `input/`)
  - `docs/` at root level following Go standards
  - `examples/` at root level
- Enhanced build system with condensed output and better error handling
- Improved `.gitignore` with comprehensive build artifact exclusions

### Cleanup

- Build artifacts from version control
- Backup files (`.backup`, `.bak`) from repository
- Old `src/` directory structure

### Bug Fixes

- Import paths updated to reflect new package structure
- Makefile paths corrected for new directory layout
- Silent mode implementation for cleaner build output

## [1.0.0] - 2024-01-XX

### Features

- Initial release of PhD Dissertation Planner
- CSV-based timeline data processing
- LaTeX calendar generation with XeLaTeX
- YAML configuration system
- Template-based PDF generation
- Command-line interface with CLI library

- Academic timeline visualization
- Task stacking and layout management
- Configurable calendar rendering
- Multiple page support
- Customizable templates

---

## Guidelines for Changelog Updates

### Types of Changes

- **Added** for new features
- **Changed** for changes in existing functionality
- **Deprecated** for soon-to-be removed features
- **Removed** for now removed features
- **Fixed** for any bug fixes
- **Security** in case of vulnerabilities

### Version Numbering

This project uses [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

### Release Process

1. Update version in relevant files
2. Move unreleased changes to new version section
3. Create git tag
4. Update release notes on GitHub
