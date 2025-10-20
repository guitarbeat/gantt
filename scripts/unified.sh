#!/bin/bash

# Unified Development Script for PhD Dissertation Planner
# This script consolidates all development, build, and maintenance operations

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

# Default values
GO_VERSION="1.22"
BINARY_NAME="plannergen"
BINARY_DIR="generated"
CSV_FILE="research_timeline_v5_comprehensive.csv"
CONFIG_FILES="configs/base.yaml,configs/monthly_calendar.yaml"

# Logging functions
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

# ==================== CORE FUNCTIONS ====================

# Setup and installation
setup() {
    log_info "Setting up project..."
    
    # Install dependencies
    go mod download
    go mod tidy
    
    # Install pre-commit hooks
    if command -v pre-commit >/dev/null 2>&1; then
        pre-commit install
    fi
    
    # Install Cursor CLI hooks if available
    if command -v cursor >/dev/null 2>&1; then
        make install-cursor-hooks 2>/dev/null || log_warning "Cursor CLI hooks not available"
    fi
    
    log_success "Project setup complete!"
}

# Build operations
build() {
    local build_type="${1:-full}"
    
    case "$build_type" in
        binary)
            log_info "Building binary..."
            go build -o "$BINARY_DIR/$BINARY_NAME" ./cmd/planner
            log_success "Binary built: $BINARY_DIR/$BINARY_NAME"
            ;;
        latex)
            log_info "Building LaTeX..."
            go build -o "$BINARY_DIR/$BINARY_NAME" ./cmd/planner
            PLANNER_SILENT=1 PLANNER_CSV_FILE="input_data/$CSV_FILE" \
            "$BINARY_DIR/$BINARY_NAME" --config "$CONFIG_FILES" --outdir "$BINARY_DIR"
            log_success "LaTeX generated: $BINARY_DIR/monthly_calendar.tex"
            ;;
        pdf)
            build latex
            log_info "Compiling PDF..."
            cd "$BINARY_DIR"
            if command -v xelatex >/dev/null 2>&1; then
                xelatex -interaction=nonstopmode monthly_calendar.tex >/dev/null 2>&1
                log_success "PDF compiled: $BINARY_DIR/monthly_calendar.pdf"
            else
                log_warning "XeLaTeX not found, PDF compilation skipped"
            fi
            cd ..
            ;;
        full)
            build pdf
            ;;
    esac
}

# Testing operations
test() {
    local test_type="${1:-all}"
    
    case "$test_type" in
        unit)
            log_info "Running unit tests..."
            go test ./tests/unit/... -v
            ;;
        integration)
            log_info "Running integration tests..."
            go test ./tests/integration/... -v
            ;;
        coverage)
            log_info "Running tests with coverage..."
            go test -coverprofile=coverage.out -covermode=atomic ./...
            go tool cover -html=coverage.out -o coverage.html
            log_success "Coverage report: coverage.html"
            ;;
        bench)
            log_info "Running benchmarks..."
            go test -bench=. -benchmem ./...
            ;;
        all)
            test unit
            test integration
            test coverage
            ;;
    esac
}

# Code quality operations
quality() {
    local quality_type="${1:-all}"
    
    case "$quality_type" in
        fmt)
            log_info "Formatting code..."
            gofmt -w .
            goimports -w internal/ pkg/ cmd/
            ;;
        lint)
            log_info "Running linters..."
            if command -v golangci-lint >/dev/null 2>&1; then
                golangci-lint run ./...
            else
                go vet ./...
            fi
            ;;
        security)
            log_info "Running security checks..."
            if command -v govulncheck >/dev/null 2>&1; then
                govulncheck ./...
            else
                log_warning "govulncheck not available"
            fi
            ;;
        all)
            quality fmt
            quality lint
            quality security
            ;;
    esac
}

# Development operations
dev() {
    local dev_type="${1:-start}"
    
    case "$dev_type" in
        start)
            log_info "Starting development environment..."
            setup
            build binary
            test all
            log_success "Development environment ready!"
            ;;
        watch)
            log_info "Starting file watcher..."
            if command -v air >/dev/null 2>&1; then
                air
            else
                log_warning "Air not installed, falling back to manual build"
                build binary
            fi
            ;;
        cursor)
            log_info "Opening in Cursor..."
            if command -v cursor >/dev/null 2>&1; then
                cursor .
            else
                log_error "Cursor not installed"
            fi
            ;;
    esac
}

# Maintenance operations
maintenance() {
    local maint_type="${1:-clean}"
    
    case "$maint_type" in
        clean)
            log_info "Cleaning build artifacts..."
            go clean -cache -testcache -modcache 2>/dev/null || true
            rm -rf "$BINARY_DIR"
            rm -f coverage.out coverage.html
            find . -name "plannergen" -type f -delete 2>/dev/null || true
            ;;
        organize)
            log_info "Organizing project files..."
            # Create organized directories
            mkdir -p .temp generated/{pdfs,tex,logs,preview}
            
            # Move scattered files
            find . -maxdepth 1 -name "*.aux" -exec mv {} .temp/ \; 2>/dev/null || true
            find . -maxdepth 1 -name "*.log" -exec mv {} .temp/ \; 2>/dev/null || true
            find . -maxdepth 1 -name "*.tmp" -exec mv {} .temp/ \; 2>/dev/null || true
            
            # Organize generated files
            find generated -name "*.pdf" -exec mv {} generated/pdfs/ \; 2>/dev/null || true
            find generated -name "*.tex" -exec mv {} generated/tex/ \; 2>/dev/null || true
            find generated -name "*.log" -exec mv {} generated/logs/ \; 2>/dev/null || true
            ;;
        deps)
            log_info "Updating dependencies..."
            go get -u ./...
            go mod tidy
            ;;
        all)
            maintenance clean
            maintenance organize
            maintenance deps
            ;;
    esac
}

# Cursor CLI operations
cursor() {
    local cursor_type="${1:-stats}"
    
    if ! command -v cursor >/dev/null 2>&1; then
        log_error "Cursor CLI not found"
        return 1
    fi
    
    case "$cursor_type" in
        stats)
            log_info "Project statistics:"
            echo "Go files: $(find . -name "*.go" -not -path "./vendor/*" | wc -l)"
            echo "Test files: $(find . -name "*_test.go" -not -path "./vendor/*" | wc -l)"
            echo "YAML files: $(find . -name "*.yaml" -o -name "*.yml" | wc -l)"
            echo "Markdown files: $(find . -name "*.md" | wc -l)"
            echo "Shell scripts: $(find . -name "*.sh" | wc -l)"
            ;;
        open)
            cursor .
            ;;
        hooks)
            make install-cursor-hooks
            ;;
        test)
            make cursor-test-enhance
            ;;
        dev)
            make cursor-dev-tools
            ;;
    esac
}

# Release operations
release() {
    local release_type="${1:-build}"
    
    case "$release_type" in
        build)
            log_info "Building release..."
            build full
            test all
            quality all
            ;;
        package)
            release build
            log_info "Packaging release..."
            if command -v goreleaser >/dev/null 2>&1; then
                goreleaser release --clean --snapshot
            else
                log_warning "GoReleaser not available"
            fi
            ;;
        publish)
            if command -v goreleaser >/dev/null 2>&1; then
                goreleaser release --clean
            else
                log_error "GoReleaser not available"
            fi
            ;;
    esac
}

# ==================== MAIN FUNCTION ====================

main() {
    local command="${1:-help}"
    local subcommand="${2:-}"
    
    echo -e "${CYAN}🛠️ PhD Dissertation Planner - Unified Development Tool${NC}"
    echo "=================================================="
    echo ""
    
    case "$command" in
        setup)
            setup
            ;;
        build)
            build "$subcommand"
            ;;
        test)
            test "$subcommand"
            ;;
        quality)
            quality "$subcommand"
            ;;
        dev)
            dev "$subcommand"
            ;;
        maintenance|maint)
            maintenance "$subcommand"
            ;;
        cursor)
            cursor "$subcommand"
            ;;
        release)
            release "$subcommand"
            ;;
        ci)
            log_info "Running CI pipeline..."
            maintenance clean
            quality all
            test all
            build full
            log_success "CI pipeline complete!"
            ;;
        help|--help|-h)
            echo "Usage: $0 <command> [subcommand]"
            echo ""
            echo "Commands:"
            echo "  setup                    - Setup project dependencies and hooks"
            echo "  build [binary|latex|pdf|full] - Build project components"
            echo "  test [unit|integration|coverage|bench|all] - Run tests"
            echo "  quality [fmt|lint|security|all] - Code quality checks"
            echo "  dev [start|watch|cursor] - Development operations"
            echo "  maintenance [clean|organize|deps|all] - Maintenance tasks"
            echo "  cursor [stats|open|hooks|test|dev] - Cursor CLI operations"
            echo "  release [build|package|publish] - Release operations"
            echo "  ci                       - Run full CI pipeline"
            echo ""
            echo "Examples:"
            echo "  $0 setup"
            echo "  $0 build pdf"
            echo "  $0 test coverage"
            echo "  $0 dev start"
            echo "  $0 cursor stats"
            echo "  $0 ci"
            ;;
        *)
            log_error "Unknown command: $command"
            echo "Use '$0 help' for usage information"
            exit 1
            ;;
    esac
    
    echo ""
    log_success "Operation complete! 🎉"
}

# Run main function
main "$@"
