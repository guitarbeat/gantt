# Scripts Directory

This directory contains automation scripts for development, building, and deployment of the PhD Dissertation Planner.

## Available Scripts

### `build_release.sh`

Creates timestamped releases with full provenance tracking.

```bash
# Create a timestamped release
./scripts/build_release.sh

# Create a named release for milestones
./scripts/build_release.sh --name "Committee_Review"

# Create weekly releases
./scripts/build_release.sh --name "Week_$(date +%U)"
```

### `setup.sh`

Development environment setup script that checks dependencies and verifies the build.

```bash
# Run setup (checks Go, tools, installs deps, verifies build)
./scripts/setup.sh
```

### `build_and_preview.ps1`

PowerShell script for building PDFs and generating preview images (Windows).

```powershell
# Build PDF and generate 3 preview images
.\scripts\build_and_preview.ps1 -Pages 3
```

### `dev.sh`

Quick development script for common tasks during development.

```bash
# Run development setup
./scripts/dev.sh
```

## Adding New Scripts

When adding new scripts to this directory:

1. Use the `.sh` extension for shell scripts
2. Include a shebang (`#!/bin/bash`) at the top
3. Make scripts executable with `chmod +x`
4. Add documentation to this README
5. Follow the established logging patterns (log_info, log_success, log_warning, log_error)

## Script Guidelines

- Include error handling with `set -e`
- Use consistent color output for logging
- Provide help messages and usage examples
- Keep scripts focused on single responsibilities
- Document dependencies and requirements
