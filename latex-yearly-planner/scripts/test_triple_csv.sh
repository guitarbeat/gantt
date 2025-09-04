#!/bin/bash

# * Test script specifically for test_triple.csv
# * This script tests the CSV processing functionality without LaTeX generation issues

set -e  # Exit on any error

# * Configuration variables
PROJECT_ROOT="/Users/aaron/Downloads/gantt/latex-yearly-planner"
CSV_FILE="../input/test_triple.csv"
TEST_OUTPUT_DIR="test_output"
BINARY="./build/plannergen"

# * Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# * Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# * Function to print section headers
print_section() {
    echo ""
    print_status $BLUE "=========================================="
    print_status $BLUE "$1"
    print_status $BLUE "=========================================="
    echo ""
}

# * Function to check if CSV file exists
check_csv_file() {
    if [ ! -f "$CSV_FILE" ]; then
        print_status $RED "Error: CSV file not found: $CSV_FILE"
        exit 1
    fi
    
    print_status $GREEN "✓ CSV file found: $CSV_FILE"
    
    # * Show CSV file contents
    print_status $YELLOW "CSV file contents:"
    cat "$CSV_FILE"
    echo ""
}

# * Function to test CSV parsing with Go
test_csv_parsing() {
    print_section "Testing CSV Parsing with Go"
    
    cd "$PROJECT_ROOT"
    
    # * Create a simple Go test program
    cat > test_csv_parser.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"
    
    "latex-yearly-planner/internal/data"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run test_csv_parser.go <csv_file>")
    }
    
    csvFile := os.Args[1]
    
    // * Create reader
    reader := data.NewReader(csvFile)
    
    // * Read tasks
    tasks, err := reader.ReadTasks()
    if err != nil {
        log.Fatalf("Failed to read tasks: %v", err)
    }
    
    fmt.Printf("✓ Successfully parsed %d tasks from CSV\n\n", len(tasks))
    
    // * Display task details
    for i, task := range tasks {
        fmt.Printf("Task %d:\n", i+1)
        fmt.Printf("  ID: %s\n", task.ID)
        fmt.Printf("  Name: %s\n", task.Name)
        fmt.Printf("  Category: %s\n", task.Priority)
        fmt.Printf("  Start Date: %s\n", task.StartDate.Format("2006-01-02"))
        fmt.Printf("  End Date: %s\n", task.EndDate.Format("2006-01-02"))
        fmt.Printf("  Description: %s\n", task.Description)
        fmt.Printf("  Status: Planned\n")
        fmt.Println()
    }
    
    // * Test date range
    dateRange, err := reader.GetDateRange()
    if err != nil {
        log.Fatalf("Failed to get date range: %v", err)
    }
    
    fmt.Printf("Date Range:\n")
    fmt.Printf("  Earliest: %s\n", dateRange.Earliest.Format("2006-01-02"))
    fmt.Printf("  Latest: %s\n", dateRange.Latest.Format("2006-01-02"))
    fmt.Println()
    
    // * Test months with tasks
    months, err := reader.GetMonthsWithTasks()
    if err != nil {
        log.Fatalf("Failed to get months with tasks: %v", err)
    }
    
    fmt.Printf("Months with tasks:\n")
    for _, month := range months {
        fmt.Printf("  %s\n", month.String())
    }
}
EOF
    
    # * Run the test
    print_status $YELLOW "Running CSV parsing test..."
    if go run test_csv_parser.go "$CSV_FILE"; then
        print_status $GREEN "✓ CSV parsing test passed!"
    else
        print_status $RED "✗ CSV parsing test failed!"
        return 1
    fi
    
    # * Clean up
    rm -f test_csv_parser.go
}

# * Function to test the planner binary with CSV
test_planner_binary() {
    print_section "Testing Planner Binary with CSV"
    
    cd "$PROJECT_ROOT"
    
    # * Check if binary exists
    if [ ! -f "$BINARY" ]; then
        print_status $YELLOW "Binary not found, building..."
        make build
    fi
    
    # * Create test output directory
    mkdir -p "$TEST_OUTPUT_DIR"
    
    # * Test with minimal configuration
    print_status $YELLOW "Testing planner binary with test_triple.csv..."
    
    # * Run with CSV file
    if PLANNER_CSV_FILE="$CSV_FILE" \
       PLANNER_YEAR=2025 \
       PASSES=1 \
       CFG="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" \
       NAME="test-triple-csv" \
       OUTDIR="$TEST_OUTPUT_DIR" \
       PLANNERGEN_BINARY="$BINARY" \
       ./scripts/build.sh; then
        print_status $GREEN "✓ Planner binary test passed!"
        
        # * Check if PDF was generated
        if [ -f "$TEST_OUTPUT_DIR/test-triple-csv.pdf" ]; then
            print_status $GREEN "✓ PDF generated successfully: $TEST_OUTPUT_DIR/test-triple-csv.pdf"
            ls -la "$TEST_OUTPUT_DIR/test-triple-csv.pdf"
        else
            print_status $YELLOW "⚠ PDF not found, but build completed"
        fi
    else
        print_status $RED "✗ Planner binary test failed!"
        return 1
    fi
}

# * Function to validate CSV structure
validate_csv_structure() {
    print_section "Validating CSV Structure"
    
    cd "$PROJECT_ROOT"
    
    print_status $YELLOW "Validating CSV file structure..."
    
    # * Check header
    header=$(head -n 1 "$CSV_FILE")
    expected_header="Task ID,Task Name,Parent Task ID,Category,Start Date,Due Date,Dependencies,Description"
    
    if [ "$header" = "$expected_header" ]; then
        print_status $GREEN "✓ CSV header is correct"
    else
        print_status $RED "✗ CSV header mismatch"
        print_status $YELLOW "Expected: $expected_header"
        print_status $YELLOW "Found: $header"
        return 1
    fi
    
    # * Count data rows
    data_rows=$(tail -n +2 "$CSV_FILE" | wc -l)
    print_status $GREEN "✓ Found $data_rows data rows"
    
    # * Check for required fields
    print_status $YELLOW "Checking for required fields in each row..."
    
    local row_num=2
    while IFS= read -r line; do
        if [ -n "$line" ]; then
            # * Count commas to ensure we have enough fields
            comma_count=$(echo "$line" | tr -cd ',' | wc -c)
            if [ "$comma_count" -ge 7 ]; then
                print_status $GREEN "✓ Row $row_num has sufficient fields"
            else
                print_status $RED "✗ Row $row_num has insufficient fields (only $comma_count commas)"
                return 1
            fi
            ((row_num++))
        fi
    done < <(tail -n +2 "$CSV_FILE")
}

# * Function to test date parsing
test_date_parsing() {
    print_section "Testing Date Parsing"
    
    cd "$PROJECT_ROOT"
    
    print_status $YELLOW "Testing date parsing from CSV..."
    
    # * Extract dates from CSV
    dates=$(tail -n +2 "$CSV_FILE" | cut -d',' -f5,6 | tr ',' '\n' | sort -u)
    
    print_status $YELLOW "Found dates in CSV:"
    echo "$dates"
    echo ""
    
    # * Validate date format
    for date in $dates; do
        if [[ $date =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}$ ]]; then
            print_status $GREEN "✓ Valid date format: $date"
        else
            print_status $RED "✗ Invalid date format: $date"
            return 1
        fi
    done
}

# * Function to clean up test files
cleanup() {
    print_section "Cleaning Up"
    
    cd "$PROJECT_ROOT"
    
    print_status $YELLOW "Cleaning up test files..."
    
    # * Remove test output directory
    if [ -d "$TEST_OUTPUT_DIR" ]; then
        rm -rf "$TEST_OUTPUT_DIR"
        print_status $GREEN "✓ Removed test output directory"
    fi
    
    # * Remove any temporary files
    rm -f test_csv_parser.go
    
    print_status $GREEN "✓ Cleanup completed"
}

# * Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -p, --parse         Test CSV parsing only"
    echo "  -b, --binary        Test planner binary only"
    echo "  -v, --validate      Validate CSV structure only"
    echo "  -d, --dates         Test date parsing only"
    echo "  -a, --all           Run all tests (default)"
    echo "  --cleanup           Clean up test files"
    echo ""
    echo "Examples:"
    echo "  $0                  # Run all tests"
    echo "  $0 -p               # Test CSV parsing only"
    echo "  $0 -v               # Validate CSV structure"
    echo "  $0 --cleanup        # Clean up test files"
}

# * Main execution
main() {
    print_section "Test Triple CSV Test Suite"
    print_status $GREEN "Testing test_triple.csv file processing..."
    
    # * Check prerequisites
    check_csv_file
    
    # * Parse command line arguments
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
        -p|--parse)
            test_csv_parsing
            ;;
        -b|--binary)
            test_planner_binary
            ;;
        -v|--validate)
            validate_csv_structure
            ;;
        -d|--dates)
            test_date_parsing
            ;;
        --cleanup)
            cleanup
            ;;
        -a|--all|"")
            validate_csv_structure
            test_date_parsing
            test_csv_parsing
            test_planner_binary
            ;;
        *)
            print_status $RED "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
    
    print_section "Test Suite Complete"
    print_status $GREEN "All requested tests completed successfully!"
}

# * Run main function with all arguments
main "$@"
