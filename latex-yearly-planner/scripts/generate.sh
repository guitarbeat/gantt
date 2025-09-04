#!/usr/bin/env bash
set -euo pipefail

# Unified script for generating proposal timeline PDFs
# Supports all CSV files in input/ directory and automatically creates releases

# Configuration
CFG_DEFAULT="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml"
OUTDIR="${OUTDIR:-build}"
RELEASE_DIR="${RELEASE_DIR:-release}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print section headers
print_section() {
    echo ""
    print_status $BLUE "=========================================="
    print_status $BLUE "$1"
    print_status $BLUE "=========================================="
    echo ""
}

# Function to generate timestamp
get_timestamp() {
    date +"%Y%m%d_%H%M%S"
}

# Function to create release directory
setup_release_dir() {
    if [ ! -d "$RELEASE_DIR" ]; then
        mkdir -p "$RELEASE_DIR"
        print_status $GREEN "Created release directory: $RELEASE_DIR"
    fi
}

# Function to generate PDF for a specific CSV file
generate_pdf() {
    local csv_file="$1"
    local csv_name=$(basename "$csv_file" .csv)
    local timestamp=$(get_timestamp)
    local output_name="proposal-timeline_${csv_name}_${timestamp}"
    
    print_section "Generating PDF for $csv_name"
    
    # Build the binary if it doesn't exist
    if [ ! -f "build/plannergen" ]; then
        print_status $YELLOW "Building plannergen binary..."
        go build -o build/plannergen ./cmd/plannergen
    fi
    
    # Generate the PDF
    print_status $YELLOW "Generating PDF with CSV: $csv_file"
    PLANNER_CSV_FILE="$csv_file" \
    OUTDIR="$OUTDIR" \
    CFG="$CFG_DEFAULT" \
    PLANNERGEN_BINARY="build/plannergen" \
    ./scripts/build.sh
    
    # Copy to release directory with timestamp
    setup_release_dir
    cp "$OUTDIR/proposal-timeline.pdf" "$RELEASE_DIR/${output_name}.pdf"
    print_status $GREEN "✓ Created: $RELEASE_DIR/${output_name}.pdf"
    
    # Also create a symlink to the latest version
    ln -sf "${output_name}.pdf" "$RELEASE_DIR/proposal-timeline_${csv_name}_latest.pdf"
    print_status $GREEN "✓ Latest symlink: $RELEASE_DIR/proposal-timeline_${csv_name}_latest.pdf"
}

# Function to list available CSV files
list_csv_files() {
    local input_dir="../input"
    if [ ! -d "$input_dir" ]; then
        print_status $RED "Error: Input directory $input_dir not found"
        exit 1
    fi
    
    local csv_files=($(find "$input_dir" -name "*.csv" -type f | sort))
    
    if [ ${#csv_files[@]} -eq 0 ]; then
        print_status $RED "No CSV files found in $input_dir"
        exit 1
    fi
    
    echo "Available CSV files:"
    for i in "${!csv_files[@]}"; do
        local csv_name=$(basename "${csv_files[$i]}" .csv)
        echo "  $((i+1)). $csv_name (${csv_files[$i]})"
    done
    echo ""
}

# Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS] [CSV_FILE|NUMBER]"
    echo ""
    echo "Generate proposal timeline PDFs from CSV files"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -l, --list          List available CSV files"
    echo "  -a, --all           Generate PDFs for all CSV files"
    echo "  -c, --clean         Clean build and release directories"
    echo ""
    echo "Arguments:"
    echo "  CSV_FILE            Path to specific CSV file"
    echo "  NUMBER              Number from --list output"
    echo ""
    echo "Examples:"
    echo "  $0 --list                    # List available CSV files"
    echo "  $0 1                         # Generate PDF for first CSV file"
    echo "  $0 ../input/data.cleaned.csv # Generate PDF for specific file"
    echo "  $0 --all                     # Generate PDFs for all CSV files"
    echo ""
}

# Function to clean directories
clean_dirs() {
    print_section "Cleaning Directories"
    rm -rf "$OUTDIR"/*.pdf "$OUTDIR"/*.aux "$OUTDIR"/*.log "$OUTDIR"/*.out "$OUTDIR"/*.tex "$OUTDIR"/*.synctex.gz
    rm -f build/plannergen
    print_status $GREEN "✓ Cleaned build directory"
    
    if [ -d "$RELEASE_DIR" ]; then
        rm -rf "$RELEASE_DIR"/*
        print_status $GREEN "✓ Cleaned release directory"
    fi
}

# Main execution
main() {
    print_section "Proposal Timeline Generator"
    
    # Parse arguments
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
        -l|--list)
            list_csv_files
            exit 0
            ;;
        -c|--clean)
            clean_dirs
            exit 0
            ;;
        -a|--all)
            list_csv_files
            local input_dir="../input"
            local csv_files=($(find "$input_dir" -name "*.csv" -type f | sort))
            
            for csv_file in "${csv_files[@]}"; do
                generate_pdf "$csv_file"
            done
            
            print_section "All PDFs Generated"
            print_status $GREEN "All proposal timeline PDFs have been generated and saved to $RELEASE_DIR/"
            ;;
        "")
            # No arguments - show help
            show_help
            exit 0
            ;;
        *)
            # Check if it's a number
            if [[ "$1" =~ ^[0-9]+$ ]]; then
                local input_dir="../input"
                local csv_files=($(find "$input_dir" -name "*.csv" -type f | sort))
                local index=$((1-1))
                
                if [ "$1" -ge 1 ] && [ "$1" -le ${#csv_files[@]} ]; then
                    generate_pdf "${csv_files[$index]}"
                else
                    print_status $RED "Error: Invalid number. Use --list to see available options."
                    exit 1
                fi
            else
                # Treat as file path
                if [ -f "$1" ]; then
                    generate_pdf "$1"
                else
                    print_status $RED "Error: File not found: $1"
                    exit 1
                fi
            fi
            ;;
    esac
}

# Run main function with all arguments
main "$@"
