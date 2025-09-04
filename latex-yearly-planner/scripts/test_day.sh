#!/bin/bash

# * Test script for day.go module
# * This script runs comprehensive tests for the calendar/day package

set -e  # Exit on any error

# * Configuration variables
PROJECT_ROOT="/Users/aaron/Downloads/gantt/latex-yearly-planner"
TEST_PACKAGE="./internal/calendar"
VERBOSE_FLAG="-v"
COVERAGE_FLAG="-cover"
BENCHMARK_FLAG="-bench=."

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

# * Function to check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_status $RED "Error: Go is not installed or not in PATH"
        exit 1
    fi
    
    local go_version=$(go version)
    print_status $GREEN "Found: $go_version"
}

# * Function to check if we're in the right directory
check_directory() {
    if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
        print_status $RED "Error: go.mod not found. Please run this script from the project root."
        exit 1
    fi
    
    print_status $GREEN "Project directory verified: $PROJECT_ROOT"
}

# * Function to run basic tests
run_tests() {
    print_section "Running Unit Tests"
    
    cd "$PROJECT_ROOT"
    
    # * Run tests with verbose output
    print_status $YELLOW "Running tests for $TEST_PACKAGE..."
    if go test $VERBOSE_FLAG $COVERAGE_FLAG $TEST_PACKAGE; then
        print_status $GREEN "✓ All tests passed!"
    else
        print_status $RED "✗ Some tests failed!"
        return 1
    fi
}

# * Function to run benchmarks
run_benchmarks() {
    print_section "Running Benchmarks"
    
    cd "$PROJECT_ROOT"
    
    print_status $YELLOW "Running benchmarks for $TEST_PACKAGE..."
    if go test $BENCHMARK_FLAG $TEST_PACKAGE; then
        print_status $GREEN "✓ Benchmarks completed!"
    else
        print_status $RED "✗ Benchmarks failed!"
        return 1
    fi
}

# * Function to generate test coverage report
generate_coverage() {
    print_section "Generating Coverage Report"
    
    cd "$PROJECT_ROOT"
    
    print_status $YELLOW "Generating coverage report..."
    
    # * Generate coverage profile
    go test -coverprofile=coverage.out $TEST_PACKAGE
    
    # * Display coverage in terminal
    go tool cover -func=coverage.out
    
    # * Generate HTML coverage report
    go tool cover -html=coverage.out -o coverage.html
    
    print_status $GREEN "✓ Coverage report generated: coverage.html"
    print_status $YELLOW "Open coverage.html in your browser to view detailed coverage"
}

# * Function to run specific test functions
run_specific_tests() {
    print_section "Running Specific Test Functions"
    
    cd "$PROJECT_ROOT"
    
    # * Test individual functions
    local test_functions=(
        "TestDay_TimeOperations"
        "TestDay_Ref"
        "TestDay_WeekLink"
        "TestDay_Quarter"
        "TestDay_Month"
        "TestDay_FormatHour"
        "TestDay_Hours"
        "TestDay_NextExists_PrevExists"
        "TestDay_TasksForDay"
        "TestDay_Day"
        "TestDay_ZeroTime"
        "TestDay_WithSpanningTasks"
        "TestDay_WithRegularTasks"
        "TestDay_Breadcrumb"
        "TestDay_PrevNext"
        "TestDay_HeadingMOS"
        "TestDay_LinkLeaf"
    )
    
    for test_func in "${test_functions[@]}"; do
        print_status $YELLOW "Running $test_func..."
        if go test -run $test_func $VERBOSE_FLAG $TEST_PACKAGE; then
            print_status $GREEN "✓ $test_func passed"
        else
            print_status $RED "✗ $test_func failed"
        fi
    done
}

# * Function to run tests with race detection
run_race_tests() {
    print_section "Running Tests with Race Detection"
    
    cd "$PROJECT_ROOT"
    
    print_status $YELLOW "Running tests with race detection..."
    if go test -race $VERBOSE_FLAG $TEST_PACKAGE; then
        print_status $GREEN "✓ No race conditions detected!"
    else
        print_status $RED "✗ Race conditions detected!"
        return 1
    fi
}

# * Function to clean up generated files
cleanup() {
    print_section "Cleaning Up"
    
    cd "$PROJECT_ROOT"
    
    print_status $YELLOW "Cleaning up generated files..."
    
    # * Remove coverage files
    rm -f coverage.out coverage.html
    
    # * Remove test binaries
    go clean -testcache
    
    print_status $GREEN "✓ Cleanup completed"
}

# * Function to show help
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -t, --test          Run unit tests only"
    echo "  -b, --benchmark     Run benchmarks only"
    echo "  -c, --coverage      Generate coverage report only"
    echo "  -s, --specific      Run specific test functions"
    echo "  -r, --race          Run tests with race detection"
    echo "  -a, --all           Run all tests (default)"
    echo "  --cleanup           Clean up generated files"
    echo ""
    echo "Examples:"
    echo "  $0                  # Run all tests"
    echo "  $0 -t               # Run unit tests only"
    echo "  $0 -c               # Generate coverage report"
    echo "  $0 -s               # Run specific test functions"
    echo "  $0 --cleanup        # Clean up generated files"
}

# * Main execution
main() {
    print_section "Day.go Test Suite"
    print_status $GREEN "Starting comprehensive testing of day.go module..."
    
    # * Check prerequisites
    check_go
    check_directory
    
    # * Parse command line arguments
    case "${1:-}" in
        -h|--help)
            show_help
            exit 0
            ;;
        -t|--test)
            run_tests
            ;;
        -b|--benchmark)
            run_benchmarks
            ;;
        -c|--coverage)
            generate_coverage
            ;;
        -s|--specific)
            run_specific_tests
            ;;
        -r|--race)
            run_race_tests
            ;;
        --cleanup)
            cleanup
            ;;
        -a|--all|"")
            run_tests
            run_benchmarks
            generate_coverage
            run_race_tests
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
