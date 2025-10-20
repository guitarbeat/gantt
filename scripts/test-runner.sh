#!/bin/bash

# Consolidated Test Runner for PhD Dissertation Planner
# This script consolidates all testing operations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Project configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Test configuration
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
BENCHMARK_FILE="benchmark.txt"

# Logging functions
log_info() {
    echo -e "${BLUE}[TEST]${NC} $1"
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

# Run unit tests
run_unit_tests() {
    log_info "Running unit tests..."
    
    local test_packages=(
        "./internal/app"
        "./internal/calendar"
        "./internal/core"
        "./pkg/templates"
        "./cmd/planner"
    )
    
    local failed_packages=()
    
    for package in "${test_packages[@]}"; do
        if [ -d "$package" ]; then
            log_info "Testing package: $package"
            if ! go test -v "$package"; then
                failed_packages+=("$package")
            fi
        else
            log_warning "Package not found: $package"
        fi
    done
    
    if [ ${#failed_packages[@]} -eq 0 ]; then
        log_success "All unit tests passed!"
        return 0
    else
        log_error "Unit tests failed in packages: ${failed_packages[*]}"
        return 1
    fi
}

# Run integration tests
run_integration_tests() {
    log_info "Running integration tests..."
    
    if [ -d "tests/integration" ]; then
        if go test -v "./tests/integration/..."; then
            log_success "Integration tests passed!"
            return 0
        else
            log_error "Integration tests failed!"
            return 1
        fi
    else
        log_warning "No integration tests found"
        return 0
    fi
}

# Run tests with coverage
run_coverage_tests() {
    log_info "Running tests with coverage..."
    
    # Clean previous coverage files
    rm -f "$COVERAGE_FILE" "$COVERAGE_HTML"
    
    # Run tests with coverage
    if go test -coverprofile="$COVERAGE_FILE" -covermode=atomic ./...; then
        # Generate HTML coverage report
        go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
        
        # Show coverage summary
        local coverage_percent
        coverage_percent=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}')
        
        log_success "Coverage tests completed!"
        log_info "Coverage: $coverage_percent"
        log_info "HTML report: $COVERAGE_HTML"
        
        return 0
    else
        log_error "Coverage tests failed!"
        return 1
    fi
}

# Run benchmarks
run_benchmarks() {
    log_info "Running benchmarks..."
    
    # Clean previous benchmark file
    rm -f "$BENCHMARK_FILE"
    
    # Run benchmarks and save output
    if go test -bench=. -benchmem ./... > "$BENCHMARK_FILE" 2>&1; then
        log_success "Benchmarks completed!"
        log_info "Benchmark results saved to: $BENCHMARK_FILE"
        
        # Show summary
        echo ""
        log_info "Benchmark Summary:"
        grep -E "Benchmark|PASS|FAIL" "$BENCHMARK_FILE" | head -20
        
        return 0
    else
        log_error "Benchmarks failed!"
        return 1
    fi
}

# Run specific test package
run_package_tests() {
    local package="$1"
    
    if [ -z "$package" ]; then
        log_error "Package name required"
        return 1
    fi
    
    log_info "Running tests for package: $package"
    
    if go test -v "$package"; then
        log_success "Package tests passed: $package"
        return 0
    else
        log_error "Package tests failed: $package"
        return 1
    fi
}

# Run tests with race detection
run_race_tests() {
    log_info "Running tests with race detection..."
    
    if go test -race ./...; then
        log_success "Race detection tests passed!"
        return 0
    else
        log_error "Race detection tests failed!"
        return 1
    fi
}

# Clean test artifacts
clean_test_artifacts() {
    log_info "Cleaning test artifacts..."
    
    # Remove test output directories
    find tests -name "generated" -type d -exec rm -rf {} + 2>/dev/null || true
    find tests -name "output" -type d -exec rm -rf {} + 2>/dev/null || true
    
    # Remove coverage files
    rm -f "$COVERAGE_FILE" "$COVERAGE_HTML"
    rm -f "$BENCHMARK_FILE"
    
    # Clean Go test cache
    go clean -testcache
    
    # Recreate empty directories
    mkdir -p tests/integration/generated
    mkdir -p tests/output
    
    log_success "Test artifacts cleaned!"
}

# Show test statistics
show_test_stats() {
    log_info "Test Statistics:"
    echo ""
    
    # Count test files
    local test_files
    test_files=$(find . -name "*_test.go" -not -path "./vendor/*" | wc -l)
    echo "Test files: $test_files"
    
    # Count test functions
    local test_functions
    test_functions=$(find . -name "*_test.go" -not -path "./vendor/*" -exec grep -c "func Test" {} + | awk '{sum += $1} END {print sum}')
    echo "Test functions: $test_functions"
    
    # Count benchmark functions
    local benchmark_functions
    benchmark_functions=$(find . -name "*_test.go" -not -path "./vendor/*" -exec grep -c "func Benchmark" {} + | awk '{sum += $1} END {print sum}')
    echo "Benchmark functions: $benchmark_functions"
    
    # Show package structure
    echo ""
    log_info "Test Package Structure:"
    find . -name "*_test.go" -not -path "./vendor/*" | sed 's|/[^/]*_test.go||' | sort | uniq -c | sort -nr
}

# Main function
main() {
    local command="${1:-all}"
    local package="${2:-}"
    
    echo -e "${CYAN}🧪 PhD Dissertation Planner - Test Runner${NC}"
    echo "=============================================="
    echo ""
    
    case "$command" in
        unit)
            run_unit_tests
            ;;
        integration)
            run_integration_tests
            ;;
        coverage)
            run_coverage_tests
            ;;
        bench|benchmark)
            run_benchmarks
            ;;
        race)
            run_race_tests
            ;;
        package)
            run_package_tests "$package"
            ;;
        clean)
            clean_test_artifacts
            ;;
        stats)
            show_test_stats
            ;;
        all)
            log_info "Running comprehensive test suite..."
            run_unit_tests
            run_integration_tests
            run_coverage_tests
            run_benchmarks
            run_race_tests
            log_success "All tests completed!"
            ;;
        --help|-h|help)
            echo "Usage: $0 [command] [package]"
            echo ""
            echo "Commands:"
            echo "  unit                    - Run unit tests"
            echo "  integration            - Run integration tests"
            echo "  coverage               - Run tests with coverage analysis"
            echo "  bench|benchmark        - Run performance benchmarks"
            echo "  race                   - Run tests with race detection"
            echo "  package <name>         - Run tests for specific package"
            echo "  clean                  - Clean test artifacts"
            echo "  stats                  - Show test statistics"
            echo "  all                    - Run all tests (default)"
            echo ""
            echo "Examples:"
            echo "  $0 unit"
            echo "  $0 coverage"
            echo "  $0 package ./internal/core"
            echo "  $0 clean"
            echo "  $0 all"
            ;;
        *)
            log_error "Unknown command: $command"
            echo "Use '$0 help' for usage information"
            exit 1
            ;;
    esac
    
    echo ""
    log_success "Test runner complete! 🎉"
}

# Run main function
main "$@"
