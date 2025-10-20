#!/bin/bash

# Simple Cursor CLI Integration
# This script provides basic Cursor CLI integration for common development tasks

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

# Logging functions
log_info() {
    echo -e "${BLUE}[CURSOR]${NC} $1"
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

# Open project in Cursor
open_project() {
    log_info "Opening project in Cursor..."
    "$CURSOR_CLI" "$PROJECT_ROOT"
    log_success "Project opened in Cursor"
}

# Open specific file in Cursor
open_file() {
    local file="$1"
    if [ -z "$file" ]; then
        log_error "File path required"
        return 1
    fi
    
    if [ ! -f "$file" ]; then
        log_error "File not found: $file"
        return 1
    fi
    
    log_info "Opening file in Cursor: $file"
    "$CURSOR_CLI" "$file"
    log_success "File opened in Cursor"
}

# Open file at specific line
open_file_at_line() {
    local file="$1"
    local line="$2"
    
    if [ -z "$file" ] || [ -z "$line" ]; then
        log_error "File path and line number required"
        return 1
    fi
    
    if [ ! -f "$file" ]; then
        log_error "File not found: $file"
        return 1
    fi
    
    log_info "Opening file at line $line in Cursor: $file"
    "$CURSOR_CLI" --goto "$file:$line"
    log_success "File opened at line $line in Cursor"
}

# Compare two files
compare_files() {
    local file1="$1"
    local file2="$2"
    
    if [ -z "$file1" ] || [ -z "$file2" ]; then
        log_error "Two file paths required"
        return 1
    fi
    
    if [ ! -f "$file1" ] || [ ! -f "$file2" ]; then
        log_error "One or both files not found"
        return 1
    fi
    
    log_info "Comparing files in Cursor: $file1 vs $file2"
    "$CURSOR_CLI" --diff "$file1" "$file2"
    log_success "File comparison opened in Cursor"
}

# Open test files
open_tests() {
    log_info "Opening test files in Cursor..."
    find tests -name "*_test.go" | head -10 | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Opening test file: $file"
            "$CURSOR_CLI" "$file"
        fi
    done
    log_success "Test files opened in Cursor"
}

# Open source files
open_source() {
    log_info "Opening source files in Cursor..."
    find internal pkg cmd -name "*.go" | head -10 | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Opening source file: $file"
            "$CURSOR_CLI" "$file"
        fi
    done
    log_success "Source files opened in Cursor"
}

# Open configuration files
open_configs() {
    log_info "Opening configuration files in Cursor..."
    find configs -name "*.yaml" -o -name "*.yml" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Opening config file: $file"
            "$CURSOR_CLI" "$file"
        fi
    done
    log_success "Configuration files opened in Cursor"
}

# Open documentation
open_docs() {
    log_info "Opening documentation in Cursor..."
    find docs -name "*.md" | head -10 | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Opening doc file: $file"
            "$CURSOR_CLI" "$file"
        fi
    done
    log_success "Documentation opened in Cursor"
}

# Show project structure
show_structure() {
    log_info "Project structure:"
    echo ""
    tree -d -L 3 -I 'vendor|.git|node_modules' 2>/dev/null || find . -type d -not -path './vendor*' -not -path './.git*' -not -path './node_modules*' | head -20
    echo ""
}

# Show file statistics
show_stats() {
    log_info "Project statistics:"
    echo ""
    echo "Go files: $(find . -name "*.go" -not -path "./vendor/*" | wc -l)"
    echo "Test files: $(find . -name "*_test.go" -not -path "./vendor/*" | wc -l)"
    echo "YAML files: $(find . -name "*.yaml" -o -name "*.yml" | wc -l)"
    echo "Markdown files: $(find . -name "*.md" | wc -l)"
    echo "Shell scripts: $(find . -name "*.sh" | wc -l)"
    echo ""
}

# Main function
main() {
    local command="${1:-help}"
    
    echo -e "${CYAN}🚀 Simple Cursor CLI Integration${NC}"
    echo "=================================="
    echo ""
    
    # Check Cursor CLI availability
    if ! check_cursor_cli; then
        exit 1
    fi
    
    case "$command" in
        open)
            open_project
            ;;
        file)
            open_file "$2"
            ;;
        line)
            open_file_at_line "$2" "$3"
            ;;
        diff)
            compare_files "$2" "$3"
            ;;
        tests)
            open_tests
            ;;
        source)
            open_source
            ;;
        configs)
            open_configs
            ;;
        docs)
            open_docs
            ;;
        structure)
            show_structure
            ;;
        stats)
            show_stats
            ;;
        all)
            log_info "Opening all project files in Cursor..."
            open_project
            ;;
        --help|-h|help)
            echo "Simple Cursor CLI Integration"
            echo ""
            echo "Usage: $0 [COMMAND] [OPTIONS]"
            echo ""
            echo "Commands:"
            echo "  open                    - Open entire project in Cursor"
            echo "  file <path>             - Open specific file in Cursor"
            echo "  line <path> <line>      - Open file at specific line"
            echo "  diff <file1> <file2>    - Compare two files in Cursor"
            echo "  tests                   - Open all test files in Cursor"
            echo "  source                  - Open all source files in Cursor"
            echo "  configs                 - Open all config files in Cursor"
            echo "  docs                    - Open all documentation in Cursor"
            echo "  structure               - Show project structure"
            echo "  stats                   - Show project statistics"
            echo "  all                     - Open entire project (same as open)"
            echo ""
            echo "Examples:"
            echo "  $0 open"
            echo "  $0 file internal/calendar/task_stacker.go"
            echo "  $0 line internal/calendar/task_stacker.go 42"
            echo "  $0 diff file1.go file2.go"
            echo "  $0 tests"
            echo "  $0 structure"
            ;;
        *)
            log_error "Unknown command: $command"
            echo "Use '$0 --help' for usage information"
            exit 1
            ;;
    esac
    
    echo ""
    log_success "Cursor CLI integration complete! 🎉"
}

# Run main function
main "$@"
