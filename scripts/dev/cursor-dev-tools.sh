#!/bin/bash

# Cursor CLI Development Tools
# This script provides various development utilities powered by Cursor CLI

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

# Directories
SOURCE_DIR="internal"
PKG_DIR="pkg"
CMD_DIR="cmd"
CONFIG_DIR="configs"
DOCS_DIR="docs"
TESTS_DIR="tests"

# Logging functions
log_info() {
    echo -e "${BLUE}[CURSOR-DEV]${NC} $1"
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

# Code refactoring with Cursor CLI
refactor_code() {
    local target="${1:-$SOURCE_DIR}"
    
    log_info "Refactoring code in $target with Cursor CLI..."
    
    # Find Go files
    local go_files
    go_files=$(find "$target" -name "*.go" -not -name "*_test.go" | head -10)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target"
        return 0
    fi
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Refactoring: $file"
            "$CURSOR_CLI" --headless --refactor "$file" --context "Refactor Go code for better readability, performance, and maintainability"
        fi
    done
    
    log_success "Code refactoring complete for $target"
}

# Code review with AI analysis
review_code() {
    local target="${1:-$SOURCE_DIR}"
    
    log_info "Performing AI code review for $target..."
    
    # Find Go files
    local go_files
    go_files=$(find "$target" -name "*.go" | head -10)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target"
        return 0
    fi
    
    # Create review report
    local review_file="code_review_$(date +%Y%m%d_%H%M%S).md"
    
    echo "# Code Review Report - $(date)" > "$review_file"
    echo "" >> "$review_file"
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Reviewing: $file"
            echo "## $file" >> "$review_file"
            echo "" >> "$review_file"
            
            # Use Cursor CLI to review the file
            "$CURSOR_CLI" --headless --review "$file" --output "${file}.review" --context "Comprehensive Go code review focusing on best practices, performance, and maintainability"
            
            # Append review to main report
            if [ -f "${file}.review" ]; then
                cat "${file}.review" >> "$review_file"
                rm "${file}.review"
            fi
            
            echo "" >> "$review_file"
        fi
    done
    
    log_success "Code review complete. Report saved to: $review_file"
}

# Optimize code performance
optimize_performance() {
    local target="${1:-$SOURCE_DIR}"
    
    log_info "Optimizing performance for $target with Cursor CLI..."
    
    # Find Go files
    local go_files
    go_files=$(find "$target" -name "*.go" -not -name "*_test.go" | head -10)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target"
        return 0
    fi
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Optimizing: $file"
            "$CURSOR_CLI" --headless --optimize "$file" --context "Optimize Go code for better performance, memory usage, and efficiency"
        fi
    done
    
    log_success "Performance optimization complete for $target"
}

# Generate documentation
generate_docs() {
    local target="${1:-$SOURCE_DIR}"
    
    log_info "Generating documentation for $target with Cursor CLI..."
    
    # Find Go files
    local go_files
    go_files=$(find "$target" -name "*.go" -not -name "*_test.go" | head -10)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target"
        return 0
    fi
    
    # Create docs directory
    mkdir -p "$DOCS_DIR/generated"
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            local doc_file="$DOCS_DIR/generated/$(basename "$file" .go).md"
            log_info "Generating docs for: $file"
            "$CURSOR_CLI" --headless --generate-docs "$file" --output "$doc_file" --context "Generate comprehensive documentation for Go code including examples and usage"
        fi
    done
    
    log_success "Documentation generation complete in $DOCS_DIR/generated/"
}

# Fix code issues automatically
fix_issues() {
    local target="${1:-$SOURCE_DIR}"
    
    log_info "Fixing code issues in $target with Cursor CLI..."
    
    # Find Go files
    local go_files
    go_files=$(find "$target" -name "*.go" | head -10)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target"
        return 0
    fi
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Fixing issues in: $file"
            "$CURSOR_CLI" --headless --fix-issues "$file" --context "Fix common Go issues including linting errors, style violations, and best practices"
        fi
    done
    
    log_success "Code issue fixing complete for $target"
}

# Analyze code complexity
analyze_complexity() {
    local target="${1:-$SOURCE_DIR}"
    
    log_info "Analyzing code complexity for $target with Cursor CLI..."
    
    # Find Go files
    local go_files
    go_files=$(find "$target" -name "*.go" -not -name "*_test.go" | head -10)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target"
        return 0
    fi
    
    # Create complexity report
    local complexity_file="complexity_analysis_$(date +%Y%m%d_%H%M%S).md"
    
    echo "# Code Complexity Analysis - $(date)" > "$complexity_file"
    echo "" >> "$complexity_file"
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Analyzing complexity: $file"
            echo "## $file" >> "$complexity_file"
            echo "" >> "$complexity_file"
            
            # Use Cursor CLI to analyze complexity
            "$CURSOR_CLI" --headless --analyze-complexity "$file" --output "${file}.complexity" --context "Analyze Go code complexity and suggest simplifications"
            
            # Append analysis to main report
            if [ -f "${file}.complexity" ]; then
                cat "${file}.complexity" >> "$complexity_file"
                rm "${file}.complexity"
            fi
            
            echo "" >> "$complexity_file"
        fi
    done
    
    log_success "Complexity analysis complete. Report saved to: $complexity_file"
}

# Security analysis
analyze_security() {
    local target="${1:-$SOURCE_DIR}"
    
    log_info "Performing security analysis for $target with Cursor CLI..."
    
    # Find Go files
    local go_files
    go_files=$(find "$target" -name "*.go" | head -10)
    
    if [ -z "$go_files" ]; then
        log_info "No Go files found in $target"
        return 0
    fi
    
    # Create security report
    local security_file="security_analysis_$(date +%Y%m%d_%H%M%S).md"
    
    echo "# Security Analysis Report - $(date)" > "$security_file"
    echo "" >> "$security_file"
    
    echo "$go_files" | while read -r file; do
        if [ -n "$file" ]; then
            log_info "Analyzing security: $file"
            echo "## $file" >> "$security_file"
            echo "" >> "$security_file"
            
            # Use Cursor CLI to analyze security
            "$CURSOR_CLI" --headless --analyze-security "$file" --output "${file}.security" --context "Analyze Go code for security vulnerabilities and best practices"
            
            # Append analysis to main report
            if [ -f "${file}.security" ]; then
                cat "${file}.security" >> "$security_file"
                rm "${file}.security"
            fi
            
            echo "" >> "$security_file"
        fi
    done
    
    log_success "Security analysis complete. Report saved to: $security_file"
}

# Generate API documentation
generate_api_docs() {
    log_info "Generating API documentation with Cursor CLI..."
    
    # Find main packages
    local main_files
    main_files=$(find "$SOURCE_DIR" -name "main.go" -o -name "*.go" | grep -E "(main|api|server)" | head -5)
    
    if [ -z "$main_files" ]; then
        log_info "No main/API files found"
        return 0
    fi
    
    # Create API docs directory
    mkdir -p "$DOCS_DIR/api"
    
    echo "$main_files" | while read -r file; do
        if [ -n "$file" ]; then
            local api_doc="$DOCS_DIR/api/$(basename "$(dirname "$file")").md"
            log_info "Generating API docs for: $file"
            "$CURSOR_CLI" --headless --generate-api-docs "$file" --output "$api_doc" --context "Generate comprehensive API documentation for Go code"
        fi
    done
    
    log_success "API documentation generated in $DOCS_DIR/api/"
}

# Main function
main() {
    local command="${1:-help}"
    local target="${2:-$SOURCE_DIR}"
    
    echo -e "${CYAN}🛠️ Cursor CLI Development Tools${NC}"
    echo "=================================="
    echo ""
    
    # Check Cursor CLI availability
    if ! check_cursor_cli; then
        exit 1
    fi
    
    case "$command" in
        refactor)
            refactor_code "$target"
            ;;
        review)
            review_code "$target"
            ;;
        optimize)
            optimize_performance "$target"
            ;;
        docs)
            generate_docs "$target"
            ;;
        fix)
            fix_issues "$target"
            ;;
        complexity)
            analyze_complexity "$target"
            ;;
        security)
            analyze_security "$target"
            ;;
        api-docs)
            generate_api_docs
            ;;
        all)
            log_info "Running comprehensive development analysis..."
            refactor_code "$target"
            review_code "$target"
            optimize_performance "$target"
            generate_docs "$target"
            fix_issues "$target"
            analyze_complexity "$target"
            analyze_security "$target"
            generate_api_docs
            ;;
        --help|-h|help)
            echo "Cursor CLI Development Tools"
            echo ""
            echo "Usage: $0 [COMMAND] [TARGET_DIR]"
            echo ""
            echo "Commands:"
            echo "  refactor    - Refactor code for better readability and maintainability"
            echo "  review      - Perform AI-powered code review"
            echo "  optimize    - Optimize code for better performance"
            echo "  docs        - Generate comprehensive documentation"
            echo "  fix         - Fix common code issues automatically"
            echo "  complexity  - Analyze code complexity and suggest simplifications"
            echo "  security    - Perform security analysis"
            echo "  api-docs    - Generate API documentation"
            echo "  all         - Run all development tools"
            echo ""
            echo "Target directories:"
            echo "  src         - Source code (default)"
            echo "  tests       - Test code"
            echo "  configs     - Configuration files"
            echo "  docs        - Documentation"
            echo ""
            echo "Examples:"
            echo "  $0 refactor src/"
            echo "  $0 review tests/"
            echo "  $0 all src/"
            echo "  $0 api-docs"
            ;;
        *)
            log_error "Unknown command: $command"
            echo "Use '$0 --help' for usage information"
            exit 1
            ;;
    esac
    
    echo ""
    log_success "Development tools execution complete! 🎉"
}

# Run main function
main "$@"
