# Scripts Directory

This directory contains automation scripts for development, building, and deployment of the PhD Dissertation Planner. Scripts are organized by category for better maintainability.

## Directory Structure

```
scripts/
├── build/           # Build and release scripts
├── dev/            # Development environment scripts
├── maintenance/    # Project maintenance and organization
└── docs/          # Documentation and guides
```

## Build Scripts (`build/`)

### `build_and_preview.sh` / `build_and_preview.ps1`

Cross-platform script for building PDFs and generating preview images.

```bash
# Linux/macOS
./scripts/build/build_and_preview.sh 3

# Windows PowerShell
.\scripts\build\build_and_preview.ps1 -Pages 3
```

### `build_release.sh`

Creates timestamped releases with full provenance tracking.

```bash
# Create a timestamped release
./scripts/build/build_release.sh

# Create a named release for milestones
./scripts/build/build_release.sh --name "Committee_Review"
```

## Development Scripts (`dev/`)

### `dev.sh` / `dev.ps1`

Development environment setup and command runner.

```bash
# Linux/macOS
./scripts/dev/dev.sh go run ./cmd/planner

# Windows PowerShell
.\scripts\dev\dev.ps1 go run ./cmd/planner
```

## Maintenance Scripts (`maintenance/`)

### `cleanup_and_organize.sh` / `cleanup_and_organize.ps1`

Project organization and cleanup utilities.

```bash
# Linux/macOS
./scripts/maintenance/cleanup_and_organize.sh --status

# Windows PowerShell
.\scripts\maintenance\cleanup_and_organize.ps1 -Status
```

### `setup.sh`

Development environment setup script.

```bash
# Run setup (checks Go, tools, installs deps, verifies build)
./scripts/maintenance/setup.sh
```

### `pdf_to_images.py`

Python utility for converting PDFs to preview images.

## Documentation (`docs/`)

- `README_CROSS_PLATFORM.md` - Cross-platform development guide
- `README_PREVIEW.md` - Preview system documentation
- `SETUP_PREVIEW.md` - Preview setup instructions

## Cross-Platform Support

All scripts support both Unix-like systems (Linux/macOS) and Windows:

- **Shell scripts** (`.sh`) work on Linux/macOS with bash/zsh
- **PowerShell scripts** (`.ps1`) work on Windows and cross-platform PowerShell
- Scripts automatically detect platform and use appropriate commands
- Consistent functionality across all platforms

## Adding New Scripts

When adding new scripts:

1. Choose appropriate category directory (`build/`, `dev/`, `maintenance/`)
2. Create both `.sh` and `.ps1` versions for cross-platform support
3. Include shebang (`#!/bin/bash`) for shell scripts
4. Make shell scripts executable with `chmod +x`
5. Add documentation to this README
6. Follow established logging patterns (log_info, log_success, log_warning, log_error)

## Script Guidelines

- Include error handling with `set -e` (shell) or try/catch (PowerShell)
- Use consistent color output for logging
- Provide help messages and usage examples
- Keep scripts focused on single responsibilities
- Document dependencies and requirements
- Support both interactive and automated usage
