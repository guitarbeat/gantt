#!/bin/bash

# Cursor CLI Pre-commit Hook Integration
# This script replaces traditional pre-commit hooks with Cursor's headless CLI
# for AI-powered code quality checks and automated fixes.

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

# Check if Cursor CLI is available
check_cursor_cli() {
    if ! command -v "$CURSOR_CLI" >/dev/null 2>&1; then
        echo -e "${RED}❌ Cursor CLI not found. Please install it first:${NC}"
        echo "   Visit: https://cursor.com/docs/cli/headless"
        echo "   Or run: npm install -g @cursor/cli"
        return 1
    fi
    return 0
}

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

# Get staged Go files
get_staged_go_files() {
    git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true
}

# Get staged files for general checks
get_staged_files() {
    git diff --cached --name-only --diff-filter=ACM || true
}

# Run Go formatting check and fix
run_go_formatting() {
    local staged_files="$1"
    
    if [ -z "$staged_files" ]; then
        log_info "No Go files staged for formatting"
        return 0
    fi

    log_info "Running Go formatting check..."
    
    # Check if files need formatting
    local unformatted_files
    unformatted_files=$(echo "$staged_files" | xargs gofmt -l 2>/dev/null || true)
    
    if [ -n "$unformatted_files" ]; then
        log_warning "Found unformatted Go files:"
        echo "$unformatted_files" | sed 's/^/  /'
        
        # Use Cursor CLI to fix formatting
        log_info "Using Cursor CLI to fix formatting..."
        echo "$unformatted_files" | while read -r file; do
            if [ -n "$file" ]; then
                log_info "Fixing formatting for: $file"
                if ! "$CURSOR_CLI" --headless --fix-formatting "$file"; then
                    log_warning "Cursor CLI formatting failed for $file, falling back to gofmt"
                    gofmt -w "$file"
                fi
            fi
        done
        
        # Stage the formatted files
        echo "$unformatted_files" | xargs git add
        log_success "Go files formatted and staged"
    else
        log_success "All Go files are properly formatted"
    fi
}

# Run Go vet
run_go_vet() {
    local staged_files="$1"
    
    if [ -z "$staged_files" ]; then
        log_info "No Go files staged for vetting"
        return 0
    fi

    log_info "Running Go vet..."
    
    # Get unique packages from staged files
    local packages
    packages=$(echo "$staged_files" | grep '\.go$' | sed 's|/[^/]*\.go$||' | sort -u || true)
    
    if [ -z "$packages" ]; then
        log_info "No Go packages to vet"
        return 0
    fi
    
    local vet_errors=0
    echo "$packages" | while read -r pkg; do
        if [ -n "$pkg" ]; then
            log_info "Vetting package: $pkg"
            if ! go vet "./$pkg" 2>/dev/null; then
                log_error "Go vet found issues in package: $pkg"
                vet_errors=1
                
                # Use Cursor CLI to analyze and suggest fixes
                log_info "Using Cursor CLI to analyze vet issues in: $pkg"
                "$CURSOR_CLI" --headless --analyze-issues "$pkg" --context "go vet issues" 2>/dev/null || true
            fi
        fi
    done
    
    if [ $vet_errors -eq 1 ]; then
        return 1
    else
        log_success "Go vet passed"
        return 0
    fi
}

# Run tests
run_tests() {
    log_info "Running tests..."
    
    # Run tests with coverage
    if ! go test -race -coverprofile=coverage.out -covermode=atomic ./...; then
        log_error "Tests failed"
        
        # Use Cursor CLI to analyze test failures
        log_info "Using Cursor CLI to analyze test failures..."
        "$CURSOR_CLI" --headless --analyze-test-failures --context "Go test failures"
        
        return 1
    else
        log_success "All tests passed"
        return 0
    fi
}

# Run YAML validation
run_yaml_validation() {
    local staged_files="$1"
    local yaml_files
    
    yaml_files=$(echo "$staged_files" | grep -E '\.(yaml|yml)$' || true)
    
    if [ -z "$yaml_files" ]; then
        log_info "No YAML files staged for validation"
        return 0
    fi

    log_info "Validating YAML files..."
    
    local yaml_errors=0
    echo "$yaml_files" | while read -r file; do
        if [ -n "$file" ]; then
            if ! python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null; then
                log_error "YAML syntax error in: $file"
                yaml_errors=1
            fi
        fi
    done
    
    if [ $yaml_errors -eq 1 ]; then
        return 1
    else
        log_success "All YAML files are valid"
        return 0
    fi
}

# Check for large files
check_large_files() {
    local staged_files="$1"
    local max_size_kb=1000
    
    log_info "Checking for large files..."
    
    local large_files=0
    echo "$staged_files" | while read -r file; do
        if [ -n "$file" ] && [ -f "$file" ]; then
            local size_kb
            size_kb=$(du -k "$file" | cut -f1)
            if [ "$size_kb" -gt "$max_size_kb" ]; then
                log_error "Large file detected: $file (${size_kb}KB > ${max_size_kb}KB)"
                large_files=1
            fi
        fi
    done
    
    if [ $large_files -eq 1 ]; then
        return 1
    else
        log_success "No large files detected"
        return 0
    fi
}

# Check for merge conflicts
check_merge_conflicts() {
    local staged_files="$1"
    
    log_info "Checking for merge conflicts..."
    
    local conflict_files=0
    echo "$staged_files" | while read -r file; do
        if [ -n "$file" ] && [ -f "$file" ]; then
            if grep -q '^<<<<<<< \|^======= \|^>>>>>>> ' "$file"; then
                log_error "Merge conflict markers found in: $file"
                conflict_files=1
            fi
        fi
    done
    
    if [ $conflict_files -eq 1 ]; then
        return 1
    else
        log_success "No merge conflicts detected"
        return 0
    fi
}

# Main pre-commit function
main() {
    echo -e "${CYAN}🔍 Cursor CLI Pre-commit Checks${NC}"
    echo "=================================="
    
    # Check Cursor CLI availability
    if ! check_cursor_cli; then
        exit 1
    fi
    
    # Get staged files
    local staged_files
    staged_files=$(get_staged_files)
    
    if [ -z "$staged_files" ]; then
        log_info "No files staged for commit"
        exit 0
    fi
    
    local staged_go_files
    staged_go_files=$(get_staged_go_files)
    
    # Run all checks
    local exit_code=0
    
    # Basic file checks
    if ! check_large_files "$staged_files"; then
        exit_code=1
    fi
    
    if ! check_merge_conflicts "$staged_files"; then
        exit_code=1
    fi
    
    if ! run_yaml_validation "$staged_files"; then
        exit_code=1
    fi
    
    # Go-specific checks
    if [ -n "$staged_go_files" ]; then
        if ! run_go_formatting "$staged_go_files"; then
            exit_code=1
        fi
        
        if ! run_go_vet "$staged_go_files"; then
            exit_code=1
        fi
        
        # Only run tests if there are Go files and not in a merge/rebase
        if [ -z "$GIT_INDEX_FILE" ] && [ -z "$GIT_DIR" ]; then
            if ! run_tests; then
                exit_code=1
            fi
        else
            log_info "Skipping tests during merge/rebase"
        fi
    fi
    
    # Summary
    echo ""
    if [ $exit_code -eq 0 ]; then
        log_success "All pre-commit checks passed! 🎉"
        echo -e "${GREEN}Ready to commit.${NC}"
    else
        log_error "Pre-commit checks failed! ❌"
        echo -e "${RED}Please fix the issues above before committing.${NC}"
        echo -e "${YELLOW}Use 'git commit --no-verify' to skip checks (not recommended).${NC}"
    fi
    
    exit $exit_code
}

# Run main function
main "$@"
