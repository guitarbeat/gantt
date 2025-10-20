#!/bin/bash

# Cursor CLI Test Enhancement Script
# This script uses Cursor CLI to enhance testing workflows with AI-powered analysis

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
CURSOR_CLI="cursor"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

# Test directories
UNIT_TEST_DIR="tests/unit"
INTEGRATION_TEST_DIR="tests/integration"
SOURCE_DIR="internal"
PKG_DIR="pkg"
CMD_DIR="cmd"

# Logging functions
log_info() {
    echo -e "${BLUE}[CURSOR-TEST]${NC} $1"
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

# Check if Cursor CLI is available
check_cursor_cli() {
    if ! command -v "$CURSOR_CLI" >/dev/null 2>&1; then
        log_error "Cursor CLI not found. Please install it first:"
        echo "   Visit: https://cursor.com/docs/cli/headless"
        echo "   Or run: npm install -g @cursor/cli"
        return 1
    fi
    return 0
}

# Analyze test coverage and suggest improvements
analyze_test_coverage() {
    log_info "Analyzing test coverage with Cursor CLI..."
    
    # Generate coverage report
    if ! go test -coverprofile=coverage.out -covermode=atomic ./... >/dev/null 2>&1; then
        log_warning "Could not generate coverage report, running tests first..."
        go test ./...
        go test -coverprofile=coverage.out -covermode=atomic ./...
    fi
    
    # Convert coverage to HTML for analysis
    go tool cover -html=coverage.out -o coverage.html >/dev/null 2>&1
    
    # Use Cursor CLI to analyze coverage
    log_info "Using Cursor CLI to analyze test coverage..."
    "$CURSOR_CLI" --headless --analyze-coverage --context "Go test coverage analysis" --input coverage.html
    
    log_success "Test coverage analysis complete"
}

# Generate missing tests using Cursor CLI
generate_missing_tests() {
    local target_dir="$1"
    local package_name="$2"
    
    if [ -z "$target_dir" ] || [ -z "$package_name" ]; then
        log_error "Target directory and package name required"
        return 1
    fi
    
    log_info "Generating missing tests for $package_name using Cursor CLI..."
    
    # Find Go files without corresponding test files
    local go_files
    go_files=$(find "$target_dir" -name "*.go" -not -name "*_test.go" | head -5)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target_dir"
        return 0
    fi
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            local test_file="${file%.go}_test.go"
            if [ ! -f "$test_file" ]; then
                log_info "Generating test file for: $file"
                "$CURSOR_CLI" --headless --generate-tests "$file" --output "$test_file" --context "Generate comprehensive Go tests for $package_name package"
            fi
        fi
    done
    
    log_success "Test generation complete for $package_name"
}

# Analyze test failures with AI
analyze_test_failures() {
    log_info "Running tests and analyzing failures with Cursor CLI..."
    
    # Run tests and capture output
    local test_output
    if ! test_output=$(go test -v ./... 2>&1); then
        log_warning "Tests failed, analyzing with Cursor CLI..."
        
        # Save test output to file for analysis
        echo "$test_output" > test_failures.log
        
        # Use Cursor CLI to analyze failures
        "$CURSOR_CLI" --headless --analyze-test-failures --input test_failures.log --context "Go test failure analysis and suggestions"
        
        # Clean up
        rm -f test_failures.log
        
        return 1
    else
        log_success "All tests passed!"
        return 0
    fi
}

# Enhance existing tests with Cursor CLI
enhance_existing_tests() {
    local test_dir="$1"
    
    if [ -z "$test_dir" ]; then
        log_error "Test directory required"
        return 1
    fi
    
    log_info "Enhancing existing tests in $test_dir with Cursor CLI..."
    
    # Find test files
    local test_files
    test_files=$(find "$test_dir" -name "*_test.go" | head -10)
    
    if [ -z "$test_files" ]; then
        log_info "No test files found in $test_dir"
        return 0
    fi
    
    echo "$test_files" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Enhancing test file: $file"
            "$CURSOR_CLI" --headless --enhance-tests "$file" --context "Improve Go test quality, add edge cases, and enhance test coverage"
        fi
    done
    
    log_success "Test enhancement complete for $test_dir"
}

# Run performance analysis with Cursor CLI
analyze_performance() {
    log_info "Running performance analysis with Cursor CLI..."
    
    # Run benchmarks
    if ! go test -bench=. -benchmem ./... > benchmark_results.txt 2>&1; then
        log_warning "Could not run benchmarks, some tests may be failing"
    fi
    
    # Use Cursor CLI to analyze performance
    if [ -f "benchmark_results.txt" ]; then
        "$CURSOR_CLI" --headless --analyze-performance --input benchmark_results.txt --context "Go benchmark analysis and optimization suggestions"
        rm -f benchmark_results.txt
    fi
    
    log_success "Performance analysis complete"
}

# Generate test documentation
generate_test_docs() {
    log_info "Generating test documentation with Cursor CLI..."
    
    # Create test documentation directory
    mkdir -p docs/testing
    
    # Generate test strategy document
    "$CURSOR_CLI" --headless --generate-docs --output docs/testing/test-strategy.md --context "Generate comprehensive test strategy documentation for Go project"
    
    # Generate test coverage report
    if [ -f "coverage.html" ]; then
        cp coverage.html docs/testing/
    fi
    
    log_success "Test documentation generated in docs/testing/"
}

# Main function
main() {
    local command="${1:-all}"
    
    echo -e "${CYAN}🧪 Cursor CLI Test Enhancement${NC}"
    echo "=================================="
    echo ""
    
    # Check Cursor CLI availability
    if ! check_cursor_cli; then
        exit 1
    fi
    
    case "$command" in
        analyze-coverage)
            analyze_test_coverage
            ;;
        generate-tests)
            local package="${2:-all}"
            case "$package" in
                unit)
                    generate_missing_tests "$UNIT_TEST_DIR" "unit"
                    ;;
                integration)
                    generate_missing_tests "$INTEGRATION_TEST_DIR" "integration"
                    ;;
                src)
                    generate_missing_tests "$SOURCE_DIR" "source"
                    ;;
                all)
                    generate_missing_tests "$UNIT_TEST_DIR" "unit"
                    generate_missing_tests "$INTEGRATION_TEST_DIR" "integration"
                    generate_missing_tests "$SOURCE_DIR" "internal"
                    generate_missing_tests "$PKG_DIR" "pkg"
                    generate_missing_tests "$CMD_DIR" "cmd"
                    ;;
                *)
                    log_error "Unknown package: $package"
                    echo "Available packages: unit, integration, internal, pkg, cmd, all"
                    exit 1
                    ;;
            esac
            ;;
        analyze-failures)
            analyze_test_failures
            ;;
        enhance-tests)
            local test_dir="${2:-$UNIT_TEST_DIR}"
            enhance_existing_tests "$test_dir"
            ;;
        performance)
            analyze_performance
            ;;
        docs)
            generate_test_docs
            ;;
        all)
            log_info "Running comprehensive test enhancement..."
            analyze_test_coverage
            generate_missing_tests "$UNIT_TEST_DIR" "unit"
            generate_missing_tests "$INTEGRATION_TEST_DIR" "integration"
            generate_missing_tests "$SOURCE_DIR" "internal"
            generate_missing_tests "$PKG_DIR" "pkg"
            generate_missing_tests "$CMD_DIR" "cmd"
            enhance_existing_tests "$UNIT_TEST_DIR"
            analyze_performance
            generate_test_docs
            ;;
        --help|-h)
            echo "Cursor CLI Test Enhancement Script"
            echo ""
            echo "Usage: $0 [COMMAND] [OPTIONS]"
            echo ""
            echo "Commands:"
            echo "  analyze-coverage    - Analyze test coverage with AI suggestions"
            echo "  generate-tests      - Generate missing tests for packages"
            echo "  analyze-failures    - Analyze test failures with AI"
            echo "  enhance-tests       - Enhance existing tests with AI"
            echo "  performance         - Run performance analysis"
            echo "  docs                - Generate test documentation"
            echo "  all                 - Run all enhancements (default)"
            echo ""
            echo "Options for generate-tests:"
            echo "  unit                - Generate unit tests"
            echo "  integration         - Generate integration tests"
            echo "  src                 - Generate tests for source code"
            echo "  all                 - Generate all tests"
            echo ""
            echo "Examples:"
            echo "  $0 analyze-coverage"
            echo "  $0 generate-tests unit"
            echo "  $0 enhance-tests tests/unit"
            echo "  $0 all"
            ;;
        *)
            log_error "Unknown command: $command"
            echo "Use '$0 --help' for usage information"
            exit 1
            ;;
    esac
    
    echo ""
    log_success "Test enhancement complete! 🎉"
}

# Run main function
main "$@"
