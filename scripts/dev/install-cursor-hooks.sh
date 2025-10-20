#!/bin/bash

# Cursor CLI Hooks Installation Script
# This script installs Cursor CLI-based pre-commit hooks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Get project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"

# Configuration
GIT_HOOKS_DIR=".git/hooks"
CURSOR_SCRIPT="scripts/dev/cursor-precommit.sh"
GIT_HOOK_SCRIPT="scripts/dev/git-hook-pre-commit"
BACKUP_DIR=".git/hooks/backup"

# Logging functions
log_info() {
    echo -e "${BLUE}[INSTALL]${NC} $1"
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

# Check if we're in a git repository
check_git_repo() {
    if [ ! -d ".git" ]; then
        log_error "Not in a git repository. Please run this script from the project root."
        exit 1
    fi
    log_success "Git repository detected"
}

# Check if Cursor CLI is available
check_cursor_cli() {
    if ! command -v cursor >/dev/null 2>&1; then
        log_warning "Cursor CLI not found in PATH"
        echo -e "${YELLOW}Please install Cursor CLI first:${NC}"
        echo "   Visit: https://cursor.com/docs/cli/headless"
        echo "   Or run: npm install -g @cursor/cli"
        echo ""
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log_success "Cursor CLI found: $(which cursor)"
    fi
}

# Check if required scripts exist
check_required_files() {
    local missing_files=0
    
    if [ ! -f "$CURSOR_SCRIPT" ]; then
        log_error "Cursor pre-commit script not found: $CURSOR_SCRIPT"
        missing_files=1
    fi
    
    if [ ! -f "$GIT_HOOK_SCRIPT" ]; then
        log_error "Git hook script not found: $GIT_HOOK_SCRIPT"
        missing_files=1
    fi
    
    if [ $missing_files -eq 1 ]; then
        log_error "Required files are missing. Please ensure all scripts are present."
        exit 1
    fi
    
    log_success "All required scripts found"
}

# Create backup of existing hooks
backup_existing_hooks() {
    log_info "Creating backup of existing hooks..."
    
    # Create backup directory
    mkdir -p "$BACKUP_DIR"
    
    # Backup existing pre-commit hook if it exists
    if [ -f "$GIT_HOOKS_DIR/pre-commit" ]; then
        local backup_file="$BACKUP_DIR/pre-commit.$(date +%Y%m%d_%H%M%S).bak"
        cp "$GIT_HOOKS_DIR/pre-commit" "$backup_file"
        log_success "Backed up existing pre-commit hook to: $backup_file"
    else
        log_info "No existing pre-commit hook found"
    fi
}

# Install the Cursor CLI hooks
install_cursor_hooks() {
    log_info "Installing Cursor CLI hooks..."
    
    # Ensure git hooks directory exists
    mkdir -p "$GIT_HOOKS_DIR"
    
    # Copy the git hook wrapper
    cp "$GIT_HOOK_SCRIPT" "$GIT_HOOKS_DIR/pre-commit"
    chmod +x "$GIT_HOOKS_DIR/pre-commit"
    
    log_success "Cursor CLI pre-commit hook installed"
}

# Test the installation
test_installation() {
    log_info "Testing installation..."
    
    # Test if the Cursor script is executable
    if [ ! -x "$CURSOR_SCRIPT" ]; then
        log_error "Cursor pre-commit script is not executable"
        return 1
    fi
    
    # Test if the hook script exists and is executable (if installed)
    if [ -f "$GIT_HOOKS_DIR/pre-commit" ]; then
        if [ ! -x "$GIT_HOOKS_DIR/pre-commit" ]; then
            log_warning "Pre-commit hook exists but is not executable, fixing..."
            chmod +x "$GIT_HOOKS_DIR/pre-commit"
        fi
        
        # Test the hook with a dry run (if possible)
        log_info "Running dry test of pre-commit hook..."
        if "$GIT_HOOKS_DIR/pre-commit" --help >/dev/null 2>&1; then
            log_success "Hook installation test passed"
        else
            log_warning "Could not test hook execution, but files are in place"
        fi
    else
        log_info "Pre-commit hook not installed yet, testing Cursor script only..."
        # Test the Cursor script directly
        if "$CURSOR_SCRIPT" --help >/dev/null 2>&1; then
            log_success "Cursor script test passed"
        else
            log_warning "Could not test Cursor script execution, but file is in place"
        fi
    fi
}

# Show installation summary
show_summary() {
    echo ""
    echo -e "${CYAN}🎉 Cursor CLI Hooks Installation Complete!${NC}"
    echo "=============================================="
    echo ""
    echo -e "${GREEN}What was installed:${NC}"
    echo "  • Cursor CLI pre-commit hook at: $GIT_HOOKS_DIR/pre-commit"
    echo "  • Cursor CLI script at: $CURSOR_SCRIPT"
    echo ""
    echo -e "${GREEN}What happens now:${NC}"
    echo "  • Every commit will run Cursor CLI checks"
    echo "  • Code formatting will be automatically fixed"
    echo "  • Tests will run before commits"
    echo "  • AI-powered analysis will suggest improvements"
    echo ""
    echo -e "${YELLOW}Commands:${NC}"
    echo "  • Test hooks: make test-cursor-hooks"
    echo "  • Uninstall: make uninstall-cursor-hooks"
    echo "  • Skip hooks: git commit --no-verify"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo "  1. Make a test commit to see the hooks in action"
    echo "  2. Check the logs if any issues occur"
    echo "  3. Customize the Cursor CLI script as needed"
}

# Uninstall function
uninstall_hooks() {
    log_info "Uninstalling Cursor CLI hooks..."
    
    # Remove the pre-commit hook
    if [ -f "$GIT_HOOKS_DIR/pre-commit" ]; then
        rm "$GIT_HOOKS_DIR/pre-commit"
        log_success "Removed pre-commit hook"
    else
        log_info "No pre-commit hook found to remove"
    fi
    
    # Restore backup if available
    local latest_backup
    latest_backup=$(ls -t "$BACKUP_DIR"/pre-commit.*.bak 2>/dev/null | head -n1)
    if [ -n "$latest_backup" ]; then
        log_info "Restoring previous pre-commit hook..."
        cp "$latest_backup" "$GIT_HOOKS_DIR/pre-commit"
        chmod +x "$GIT_HOOKS_DIR/pre-commit"
        log_success "Restored previous pre-commit hook"
    fi
    
    echo ""
    echo -e "${GREEN}✅ Cursor CLI hooks uninstalled${NC}"
    echo "Previous hooks have been restored if available."
}

# Main installation function
install() {
    echo -e "${CYAN}🔧 Installing Cursor CLI Pre-commit Hooks${NC}"
    echo "=============================================="
    echo ""
    
    check_git_repo
    check_cursor_cli
    check_required_files
    backup_existing_hooks
    install_cursor_hooks
    test_installation
    show_summary
}

# Main execution
case "${1:-install}" in
    install)
        install
        ;;
    uninstall)
        uninstall_hooks
        ;;
    test)
        check_git_repo
        check_required_files
        test_installation
        ;;
    --help|-h)
        echo "Cursor CLI Hooks Installation Script"
        echo ""
        echo "Usage: $0 [COMMAND]"
        echo ""
        echo "Commands:"
        echo "  install     Install Cursor CLI hooks (default)"
        echo "  uninstall   Remove Cursor CLI hooks and restore previous"
        echo "  test        Test the installation"
        echo "  --help, -h  Show this help message"
        echo ""
        echo "This script replaces traditional pre-commit hooks with Cursor CLI"
        echo "for AI-powered code quality checks and automated fixes."
        ;;
    *)
        log_error "Unknown command: $1"
        echo "Use '$0 --help' for usage information"
        exit 1
        ;;
esac
