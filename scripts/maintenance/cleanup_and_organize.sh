#!/bin/bash

# PhD Dissertation Planner - Cleanup and Organization Script
# This script helps maintain a clean, organized project structure

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${BLUE}🧹 PhD Dissertation Planner - Cleanup & Organization${NC}"
echo "=================================================="

# Function to create directories if they don't exist
create_dirs() {
    echo -e "${YELLOW}📁 Creating organization directories...${NC}"
    
    # Main organization directories
    mkdir -p .temp                    # Temporary files
    mkdir -p .build_artifacts         # Build artifacts
    mkdir -p docs/archive             # Archived documentation
    mkdir -p docs/examples            # Example configurations
    mkdir -p generated/preview        # Preview builds
    mkdir -p generated/logs           # Log files
    mkdir -p generated/tex            # LaTeX source files
    mkdir -p generated/pdfs           # Generated PDFs
    
    echo -e "${GREEN}✅ Directory structure created${NC}"
}

# Function to clean up scattered files
cleanup_scattered_files() {
    echo -e "${YELLOW}🧽 Cleaning up scattered files...${NC}"
    
    # Move scattered LaTeX files to temp
    find . -maxdepth 1 -name "*.aux" -exec mv {} .temp/ \; 2>/dev/null || true
    find . -maxdepth 1 -name "*.log" -exec mv {} .temp/ \; 2>/dev/null || true
    find . -maxdepth 1 -name "*.tmp" -exec mv {} .temp/ \; 2>/dev/null || true
    find . -maxdepth 1 -name "*.out" -exec mv {} .temp/ \; 2>/dev/null || true
    find . -maxdepth 1 -name "*.fdb_latexmk" -exec mv {} .temp/ \; 2>/dev/null || true
    find . -maxdepth 1 -name "*.fls" -exec mv {} .temp/ \; 2>/dev/null || true
    find . -maxdepth 1 -name "*.synctex.gz" -exec mv {} .temp/ \; 2>/dev/null || true
    
    # Move generated files to proper locations
    if [ -d "generated" ]; then
        find generated -name "*.pdf" -exec mv {} generated/pdfs/ \; 2>/dev/null || true
        find generated -name "*.tex" -exec mv {} generated/tex/ \; 2>/dev/null || true
        find generated -name "*.log" -exec mv {} generated/logs/ \; 2>/dev/null || true
        find generated -name "*.aux" -exec mv {} .temp/ \; 2>/dev/null || true
    fi
    
    echo -e "${GREEN}✅ Scattered files cleaned up${NC}"
}

# Function to organize documentation
organize_docs() {
    echo -e "${YELLOW}📚 Organizing documentation...${NC}"
    
    # Move PDF documents to archive
    find docs -name "*.pdf" -exec mv {} docs/archive/ \; 2>/dev/null || true
    
    # Create example configurations
    if [ -f "configs/base.yaml" ]; then
        cp configs/base.yaml docs/examples/ 2>/dev/null || true
    fi
    
    echo -e "${GREEN}✅ Documentation organized${NC}"
}

# Function to clean up test artifacts
cleanup_test_artifacts() {
    echo -e "${YELLOW}🧪 Cleaning up test artifacts...${NC}"
    
    # Clean test output directories
    find tests -name "generated" -type d -exec rm -rf {} + 2>/dev/null || true
    find tests -name "output" -type d -exec rm -rf {} + 2>/dev/null || true
    
    # Recreate empty directories
    mkdir -p tests/integration/generated
    mkdir -p tests/output
    
    echo -e "${GREEN}✅ Test artifacts cleaned up${NC}"
}

# Function to update .gitignore
update_gitignore() {
    echo -e "${YELLOW}📝 Updating .gitignore...${NC}"
    
    # Add new patterns to .gitignore
    cat >> .gitignore << 'EOF'

# Organization cleanup
.temp/
.build_artifacts/
docs/archive/*.pdf
generated/logs/
generated/tex/
generated/pdfs/

# Temporary files
*.aux
*.log
*.tmp
*.out
*.fdb_latexmk
*.fls
*.synctex.gz
EOF
    
    echo -e "${GREEN}✅ .gitignore updated${NC}"
}

# Function to create a project structure overview
create_structure_overview() {
    echo -e "${YELLOW}📋 Creating project structure overview...${NC}"
    
    cat > PROJECT_STRUCTURE.md << 'EOF'
# 📁 Project Structure Overview

## 🎯 Core Directories

### Source Code
- `src/` - Main application source code
  - `app/` - Application logic and CLI
  - `core/` - Core utilities and configuration
  - `calendar/` - Calendar generation and layout
  - `shared/` - Shared templates and utilities

### Configuration
- `configs/` - YAML configuration files
- `cmd/` - Application entry points

### Data
- `input_data/` - CSV input files and data
- `generated/` - Generated output files
  - `pdfs/` - Generated PDF files
  - `tex/` - LaTeX source files
  - `logs/` - Build and error logs
  - `preview/` - Preview builds

### Documentation
- `docs/` - Project documentation
  - `tasks/` - How-to guides
  - `fyi/` - Reference information
  - `archive/` - Archived documents
  - `examples/` - Example configurations

### Releases
- `releases/` - Versioned releases with timestamps

### Scripts
- `scripts/` - Build and utility scripts

### Tests
- `tests/` - Test files and test data
  - `integration/` - Integration tests
  - `unit/` - Unit tests

### Temporary
- `.temp/` - Temporary files (gitignored)
- `.build_artifacts/` - Build artifacts (gitignored)

## 🧹 Cleanup Commands

```bash
# Run full cleanup
./scripts/cleanup_and_organize.sh

# Clean only scattered files
./scripts/cleanup_and_organize.sh --scattered-only

# Clean only test artifacts
./scripts/cleanup_and_organize.sh --test-only
```

## 📝 File Organization Rules

1. **Generated files** go in `generated/` with appropriate subdirectories
2. **Temporary files** go in `.temp/` (gitignored)
3. **Documentation** goes in `docs/` with logical subdirectories
4. **Configuration** goes in `configs/`
5. **Source code** goes in `src/` with clear module separation
6. **Tests** go in `tests/` with appropriate subdirectories
EOF
    
    echo -e "${GREEN}✅ Project structure overview created${NC}"
}

# Function to show current status
show_status() {
    echo -e "${BLUE}📊 Current Project Status${NC}"
    echo "=========================="
    
    echo -e "\n${YELLOW}📁 Directory Structure:${NC}"
    tree -d -L 3 -I 'vendor|.git' 2>/dev/null || find . -type d -not -path './vendor*' -not -path './.git*' | head -20
    
    echo -e "\n${YELLOW}📄 File Counts:${NC}"
    echo "Go files: $(find . -name "*.go" | wc -l)"
    echo "YAML files: $(find . -name "*.yaml" -o -name "*.yml" | wc -l)"
    echo "Markdown files: $(find . -name "*.md" | wc -l)"
    echo "PDF files: $(find . -name "*.pdf" | wc -l)"
    echo "LaTeX files: $(find . -name "*.tex" | wc -l)"
    
    echo -e "\n${YELLOW}🧹 Cleanup Status:${NC}"
    echo "Scattered files in .temp/: $(ls .temp/ 2>/dev/null | wc -l)"
    echo "Generated PDFs: $(ls generated/pdfs/ 2>/dev/null | wc -l)"
    echo "Generated LaTeX: $(ls generated/tex/ 2>/dev/null | wc -l)"
}

# Main execution
main() {
    case "${1:-}" in
        --scattered-only)
            cleanup_scattered_files
            ;;
        --test-only)
            cleanup_test_artifacts
            ;;
        --status)
            show_status
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --scattered-only    Clean only scattered files"
            echo "  --test-only         Clean only test artifacts"
            echo "  --status            Show current project status"
            echo "  --help, -h          Show this help message"
            echo ""
            echo "Default: Run full cleanup and organization"
            ;;
        *)
            create_dirs
            cleanup_scattered_files
            organize_docs
            cleanup_test_artifacts
            update_gitignore
            create_structure_overview
            show_status
            echo -e "\n${GREEN}🎉 Cleanup and organization complete!${NC}"
            ;;
    esac
}

# Run main function with all arguments
main "$@"
