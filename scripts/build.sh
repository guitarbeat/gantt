#!/bin/bash

# PhD Dissertation Planner - Build Script
# This script provides a convenient way to build the project with various options

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Help function
show_help() {
    echo "PhD Dissertation Planner Build Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -c, --clean       Clean build artifacts before building"
    echo "  -t, --test        Run tests after building"
    echo "  -l, --lint        Run linting and code quality checks"
    echo "  -v, --verbose     Enable verbose output"
    echo "  -h, --help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Simple build"
    echo "  $0 --clean           # Clean and build"
    echo "  $0 --clean --test    # Clean, build and test"
    echo "  $0 --lint            # Run code quality checks"
}

# Parse command line arguments
CLEAN=false
TEST=false
VERBOSE=false
LINT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--clean)
            CLEAN=true
            shift
            ;;
        -t|--test)
            TEST=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -l|--lint)
            LINT=true
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

# Main build process
log_info "Starting build process..."

if [ "$LINT" = true ]; then
    log_info "Running code quality checks..."

    # Check Go formatting
    log_info "Checking Go formatting..."
    if ! gofmt -d . | tee /tmp/gofmt.out && [ -s /tmp/gofmt.out ]; then
        log_error "Code is not properly formatted. Run 'gofmt -w .' to fix."
        cat /tmp/gofmt.out
        exit 1
    fi
    log_success "Code formatting is correct"

    # Run go vet
    log_info "Running go vet..."
    if ! go vet ./...; then
        log_error "go vet found issues"
        exit 1
    fi
    log_success "go vet passed"

    # Check for common issues
    log_info "Checking for common issues..."
    # Check for TODO comments (warning only)
    if grep -r "TODO\|FIXME\|XXX" --include="*.go" . > /dev/null 2>&1; then
        log_warning "Found TODO/FIXME/XXX comments in code"
    fi

    log_success "Code quality checks completed!"
    exit 0
fi

if [ "$CLEAN" = true ]; then
    log_info "Cleaning build artifacts..."
    make -f scripts/Makefile clean
fi

log_info "Building project..."
if [ "$VERBOSE" = true ]; then
    make -f scripts/Makefile
else
    make -f scripts/Makefile > /dev/null 2>&1
fi

log_success "Build completed successfully!"

if [ "$TEST" = true ]; then
    log_info "Running tests..."
    if ! go test ./...; then
        log_error "Tests failed"
        exit 1
    fi
    log_success "Tests passed"
fi

log_success "Build script completed!"
