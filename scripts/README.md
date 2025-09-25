# Scripts Directory

This directory contains automation scripts for development, building, and deployment of the PhD Dissertation Planner.

## Available Scripts

### `build.sh`
Enhanced build script with additional options beyond the basic `make` command.

```bash
# Simple build (same as make)
./scripts/build.sh

# Clean and build
./scripts/build.sh --clean

# Clean, build, and run tests
./scripts/build.sh --clean --test

# Run code quality checks only
./scripts/build.sh --lint

# Verbose output
./scripts/build.sh --verbose

# Show help
./scripts/build.sh --help
```

### `setup.sh`
Development environment setup script that checks dependencies and verifies the build.

```bash
# Run setup (checks Go, tools, installs deps, verifies build)
./scripts/setup.sh
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
