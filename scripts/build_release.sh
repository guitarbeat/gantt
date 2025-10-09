#!/bin/bash

# PhD Dissertation Planner - Release Build Script
# This script builds timestamped releases for version tracking

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_header() {
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  $1"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}"
}

# Help function
show_help() {
    echo "PhD Dissertation Planner - Release Build Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -v, --version VERSION    Version identifier (e.g., v5.1, default: auto-detect)"
    echo "  -c, --csv FILE          CSV file to use (default: latest v5.1)"
    echo "  -n, --name NAME         Custom release name (default: based on version)"
    echo "  -s, --skip-pdf          Skip PDF generation (LaTeX only)"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Build with auto-detected version"
    echo "  $0 --version v5.1                    # Build v5.1 release"
    echo "  $0 --version v5.1 --name 'Final'     # Build with custom name"
    echo "  $0 --csv input_data/custom.csv       # Build with custom CSV"
}

# Parse command line arguments
VERSION=""
CSV_FILE=""
CUSTOM_NAME=""
SKIP_PDF=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -c|--csv)
            CSV_FILE="$2"
            shift 2
            ;;
        -n|--name)
            CUSTOM_NAME="$2"
            shift 2
            ;;
        -s|--skip-pdf)
            SKIP_PDF=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Auto-detect version from CSV file if not specified
if [ -z "$CSV_FILE" ]; then
    if [ -f "input_data/research_timeline_v5.1_comprehensive.csv" ]; then
        CSV_FILE="research_timeline_v5.1_comprehensive.csv"
        [ -z "$VERSION" ] && VERSION="v5.1"
    else
        CSV_FILE="research_timeline_v5_comprehensive.csv"
        [ -z "$VERSION" ] && VERSION="v5.0"
    fi
else
    # Extract version from CSV filename if not specified
    if [ -z "$VERSION" ]; then
        VERSION=$(echo "$CSV_FILE" | grep -oP 'v\d+\.\d+' | head -1)
        [ -z "$VERSION" ] && VERSION="custom"
    fi
fi

# Generate timestamp
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
DATE_HUMAN=$(date +"%Y-%m-%d %H:%M:%S")

# Build release name
if [ -n "$CUSTOM_NAME" ]; then
    RELEASE_NAME="${CUSTOM_NAME}"
else
    RELEASE_NAME="release"
fi

# Create timestamped release directory
RELEASE_DIR="releases/${TIMESTAMP}_${RELEASE_NAME}"
mkdir -p "$RELEASE_DIR"

log_header "PhD Dissertation Planner - Release Build"
echo ""
log_info "Configuration:"
echo "  • Version:    ${YELLOW}${VERSION}${NC}"
echo "  • CSV File:   ${YELLOW}${CSV_FILE}${NC}"
echo "  • Timestamp:  ${YELLOW}${DATE_HUMAN}${NC}"
echo "  • Name:       ${YELLOW}${RELEASE_NAME}${NC}"
echo "  • Output:     ${YELLOW}${RELEASE_DIR}/${NC}"
echo ""

# Build using the main Makefile
log_info "Building project using Makefile..."
make clean
if [ "$SKIP_PDF" = true ]; then
    make build-latex CSV_FILE="$CSV_FILE"
else
    make build CSV_FILE="$CSV_FILE"
fi

# Copy artifacts to release directory
log_info "Copying artifacts to release directory..."
cp "generated/monthly_calendar.pdf" "$RELEASE_DIR/planner.pdf" 2>/dev/null || true
cp "generated/monthly_calendar.tex" "$RELEASE_DIR/planner.tex" 2>/dev/null || true
cp "generated/monthly_calendar.log" "$RELEASE_DIR/xelatex.log" 2>/dev/null || true
cp "input_data/$CSV_FILE" "$RELEASE_DIR/source.csv" 2>/dev/null || true

# Create release metadata
METADATA_FILE="$RELEASE_DIR/metadata.json"
cat > "$METADATA_FILE" << EOF
{
  "version": "$VERSION",
  "release_name": "$RELEASE_NAME",
  "timestamp": "$TIMESTAMP",
  "date": "$DATE_HUMAN",
  "csv_file": "$CSV_FILE",
  "csv_basename": "$(basename $CSV_FILE)",
  "files": {
    "latex": "planner.tex",
    "pdf": "planner.pdf",
    "log": "xelatex.log",
    "csv": "source.csv"
  },
  "build_info": {
    "go_version": "$(go version | cut -d' ' -f3)",
    "hostname": "$(hostname)",
    "user": "$(whoami)"
  }
}
EOF
log_success "Metadata saved: $METADATA_FILE"

# Create release README
README_FILE="$RELEASE_DIR/README.md"
cat > "$README_FILE" << EOF
# Release: $RELEASE_NAME

**Generated:** $DATE_HUMAN  
**Version:** $VERSION  
**Timestamp:** $TIMESTAMP

## Files

- **planner.pdf** - Compiled PDF planner
- **planner.tex** - LaTeX source (main document)
- **xelatex.log** - LaTeX compilation log
- **source.csv** - Original CSV data
- **metadata.json** - Build metadata

## Source

**CSV File:** \`$(basename $CSV_FILE)\`

## Build Info

- Go: $(go version | cut -d' ' -f3)
- User: $(whoami)
- Host: $(hostname)

## Usage

\`\`\`bash
# View PDF
open planner.pdf

# Recompile LaTeX if needed
xelatex planner.tex

# View source data
cat source.csv
\`\`\`

---
*Generated by PhD Dissertation Planner build system*
EOF
log_success "README saved: $README_FILE"

# Update releases index
INDEX_FILE="releases/INDEX.md"

# Update main index
if [ ! -f "$INDEX_FILE" ]; then
    cat > "$INDEX_FILE" << 'EOF'
# PhD Dissertation Planner - Releases

This directory contains timestamped releases organized by version and timestamp.

## Quick Access

See the main INDEX.md file for all releases.

## Release History

EOF
fi

# Add to main index
echo "### ${TIMESTAMP}_${RELEASE_NAME}" >> "$INDEX_FILE"
echo "" >> "$INDEX_FILE"
echo "- **Date:** $DATE_HUMAN" >> "$INDEX_FILE"
echo "- **Version:** $VERSION" >> "$INDEX_FILE"
echo "- **CSV:** $(basename $CSV_FILE)" >> "$INDEX_FILE"
echo "- **Location:** \`$RELEASE_DIR/\`" >> "$INDEX_FILE"
echo "" >> "$INDEX_FILE"


log_success "Release indexes updated"

# Summary
echo ""
log_header "Release Build Complete"
echo ""
echo -e "  📦 ${GREEN}Release:${NC}  ${YELLOW}${RELEASE_NAME}${NC}"
echo -e "  ⏰ ${GREEN}Timestamp:${NC} ${YELLOW}${TIMESTAMP}${NC}"
echo -e "  📁 ${GREEN}Location:${NC} ${YELLOW}${RELEASE_DIR}/${NC}"
echo ""
echo "  Files created:"
if [ "$SKIP_PDF" = false ] && [ -f "$RELEASE_DIR/planner.pdf" ]; then
    PDF_SIZE=$(stat -f%z "$RELEASE_DIR/planner.pdf" 2>/dev/null || stat -c%s "$RELEASE_DIR/planner.pdf")
    PDF_SIZE_KB=$((PDF_SIZE / 1024))
    echo -e "    ${GREEN}✓${NC} planner.pdf (${PDF_SIZE_KB} KB)"
fi
if [ -f "$RELEASE_DIR/planner.tex" ]; then
    TEX_SIZE=$(stat -f%z "$RELEASE_DIR/planner.tex" 2>/dev/null || stat -c%s "$RELEASE_DIR/planner.tex")
    TEX_SIZE_KB=$((TEX_SIZE / 1024))
    echo -e "    ${GREEN}✓${NC} planner.tex (${TEX_SIZE_KB} KB)"
fi
if [ -f "$RELEASE_DIR/xelatex.log" ]; then
    LOG_SIZE=$(stat -f%z "$RELEASE_DIR/xelatex.log" 2>/dev/null || stat -c%s "$RELEASE_DIR/xelatex.log")
    LOG_SIZE_KB=$((LOG_SIZE / 1024))
    echo -e "    ${GREEN}✓${NC} xelatex.log (${LOG_SIZE_KB} KB)"
fi
echo -e "    ${GREEN}✓${NC} source.csv"
echo -e "    ${GREEN}✓${NC} metadata.json"
echo -e "    ${GREEN}✓${NC} README.md"
echo ""
echo "  Quick access:"
echo "    ${CYAN}open $RELEASE_DIR/planner.pdf${NC}"
echo "    ${CYAN}cat $RELEASE_DIR/README.md${NC}"
echo ""
echo "  View all releases:"
echo "    ${CYAN}cat releases/INDEX.md${NC}"
echo ""