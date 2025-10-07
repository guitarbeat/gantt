#!/bin/bash

# PhD Dissertation Planner - Quick Build Script
# Simplified build workflow for development and testing

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

# Help function
show_help() {
    echo "PhD Dissertation Planner - Quick Build Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -c, --csv FILE          CSV file to use (default: auto-detect)"
    echo "  -p, --preset PRESET     Configuration preset (academic, compact, presentation)"
    echo "  -n, --name NAME         Output name (default: planner)"
    echo "  -s, --skip-pdf          Skip PDF generation (LaTeX only)"
    echo "  -v, --validate          Validate CSV only (no build)"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Quick build with auto-detected CSV"
    echo "  $0 --preset compact                  # Build with compact preset"
    echo "  $0 --validate                        # Validate CSV only"
    echo "  $0 --csv custom.csv --name my_plan   # Build with custom CSV and name"
}

# Parse command line arguments
CSV_FILE=""
PRESET=""
OUTPUT_NAME="planner"
SKIP_PDF=false
VALIDATE_ONLY=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--csv)
            CSV_FILE="$2"
            shift 2
            ;;
        -p|--preset)
            PRESET="$2"
            shift 2
            ;;
        -n|--name)
            OUTPUT_NAME="$2"
            shift 2
            ;;
        -s|--skip-pdf)
            SKIP_PDF=true
            shift
            ;;
        -v|--validate)
            VALIDATE_ONLY=true
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

# Auto-detect CSV file if not specified
if [ -z "$CSV_FILE" ]; then
    if [ -f "input_data/research_timeline_v5.1_comprehensive.csv" ]; then
        CSV_FILE="input_data/research_timeline_v5.1_comprehensive.csv"
    else
        CSV_FILE="input_data/research_timeline_v5_comprehensive.csv"
    fi
fi

# Build output directory
OUTPUT_DIR="generated"
mkdir -p "$OUTPUT_DIR"

# Build command arguments
BUILD_ARGS=""
if [ -n "$PRESET" ]; then
    BUILD_ARGS="$BUILD_ARGS --preset $PRESET"
fi
# Note: --skip-pdf is handled by the script, not passed to Go application

# Show configuration
echo -e "${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║${NC}  PhD Dissertation Planner - Quick Build"
echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}"
echo ""
log_info "Configuration:"
echo "  • CSV File:   ${YELLOW}${CSV_FILE}${NC}"
echo "  • Preset:     ${YELLOW}${PRESET:-default}${NC}"
echo "  • Output:     ${YELLOW}${OUTPUT_NAME}${NC}"
echo "  • Directory:  ${YELLOW}${OUTPUT_DIR}/${NC}"
echo ""

# Validate only mode
if [ "$VALIDATE_ONLY" = true ]; then
    log_info "Validating CSV file..."
    if ! PLANNER_CSV_FILE="$CSV_FILE" go run ./cmd/planner --validate; then
        log_error "CSV validation failed"
        exit 1
    fi
    log_success "CSV validation passed"
    exit 0
fi

# Build the binary
log_info "Building planner binary..."
if ! go build -o "$OUTPUT_DIR/plannergen" ./cmd/planner; then
    log_error "Failed to build binary"
    exit 1
fi
log_success "Binary built successfully"

# Generate LaTeX
log_info "Generating LaTeX from CSV data..."
if ! PLANNER_CSV_FILE="$CSV_FILE" \
    "$OUTPUT_DIR/plannergen" --config "src/core/base.yaml,src/core/monthly_calendar.yaml" --outdir "$OUTPUT_DIR" $BUILD_ARGS; then
    log_error "Failed to generate LaTeX"
    exit 1
fi
log_success "LaTeX generated successfully"

# Check LaTeX file
if [ ! -f "$OUTPUT_DIR/monthly_calendar.tex" ]; then
    log_error "LaTeX file not found"
    exit 1
fi

# Rename output files
mv "$OUTPUT_DIR/monthly_calendar.tex" "$OUTPUT_DIR/${OUTPUT_NAME}.tex"
log_success "LaTeX saved: ${OUTPUT_DIR}/${OUTPUT_NAME}.tex"

# Generate PDF if not skipped
if [ "$SKIP_PDF" = false ]; then
    if command -v xelatex >/dev/null 2>&1; then
        log_info "Compiling PDF with XeLaTeX..."
        
        cd "$OUTPUT_DIR"
        # Run xelatex (allow warnings, just check if PDF is created)
        xelatex -file-line-error -interaction=nonstopmode "${OUTPUT_NAME}.tex" > "${OUTPUT_NAME}.tmp" 2>&1 || true
        
        if [ -f "${OUTPUT_NAME}.pdf" ]; then
            PDF_SIZE=$(stat -f%z "${OUTPUT_NAME}.pdf" 2>/dev/null || stat -c%s "${OUTPUT_NAME}.pdf")
            
            # Check if PDF is valid (>10KB)
            if [ "$PDF_SIZE" -gt 10000 ]; then
                log_success "PDF compiled successfully ($PDF_SIZE bytes)"
            else
                log_warning "PDF created but unusually small ($PDF_SIZE bytes)"
            fi
        else
            log_warning "PDF not created - check LaTeX errors in build log"
            log_info "LaTeX source saved successfully, continuing..."
        fi
        
        cd "$PROJECT_DIR"
    else
        log_warning "XeLaTeX not found - skipping PDF generation"
        log_info "To install: brew install --cask mactex (macOS)"
    fi
else
    log_info "PDF generation skipped (--skip-pdf flag)"
fi

# Summary
echo ""
log_success "Quick build complete!"
echo ""
echo "  Files created:"
if [ "$SKIP_PDF" = false ] && [ -f "$OUTPUT_DIR/${OUTPUT_NAME}.pdf" ]; then
    PDF_SIZE=$(stat -f%z "$OUTPUT_DIR/${OUTPUT_NAME}.pdf" 2>/dev/null || stat -c%s "$OUTPUT_DIR/${OUTPUT_NAME}.pdf")
    PDF_SIZE_KB=$((PDF_SIZE / 1024))
    echo -e "    ${GREEN}✓${NC} ${OUTPUT_NAME}.pdf (${PDF_SIZE_KB} KB)"
fi
echo -e "    ${GREEN}✓${NC} ${OUTPUT_NAME}.tex"
echo ""
echo "  Quick access:"
if [ -f "$OUTPUT_DIR/${OUTPUT_NAME}.pdf" ]; then
    echo "    ${CYAN}open $OUTPUT_DIR/${OUTPUT_NAME}.pdf${NC}"
fi
echo "    ${CYAN}cat $OUTPUT_DIR/${OUTPUT_NAME}.tex${NC}"
echo ""
