# Changelog

All notable changes to the PhD Dissertation Planner will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `scripts/` directory with build and setup automation scripts
- `assets/` directory for static resources (PDFs, documents)
- `CONTRIBUTING.md` with development guidelines
- Standard Go project structure following community best practices

### Changed

- Reorganized project structure:
  - `cmd/planner/` for main application entry point
  - `pkg/templates/` for reusable template components
  - `configs/` for configuration files (renamed from `config/`)
  - `data/` for input data files (renamed from `input/`)
  - `docs/` at root level following Go standards
  - `examples/` at root level
- Enhanced build system with condensed output and better error handling
- Improved `.gitignore` with comprehensive build artifact exclusions

### Removed

- Build artifacts from version control
- Backup files (`.backup`, `.bak`) from repository
- Old `src/` directory structure

### Fixed

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
