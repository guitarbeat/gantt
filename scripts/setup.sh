#!/bin/bash

# PhD Dissertation Planner - Development Setup Script
# This script sets up the development environment

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

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go first."
        log_info "Visit: https://golang.org/dl/"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go version: $GO_VERSION"
}

# Check if required tools are installed
check_tools() {
    local missing_tools=()

    # Check for LaTeX (required for PDF generation)
    if ! command -v xelatex &> /dev/null; then
        missing_tools+=("XeLaTeX (for PDF generation)")
    fi

    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_warning "Some optional tools are missing:"
        for tool in "${missing_tools[@]}"; do
            log_warning "  - $tool"
        done
        log_info "PDF generation may not work without these tools."
    fi
}

# Install Go dependencies
install_deps() {
    log_info "Installing Go dependencies..."
    go mod download
    go mod tidy
    log_success "Dependencies installed"
}

# Verify the setup
verify_setup() {
    log_info "Verifying setup..."

    # Try to build the project
    if make > /dev/null 2>&1; then
        log_success "Build verification passed"
    else
        log_error "Build verification failed"
        exit 1
    fi
}

# Main setup process
main() {
    echo "🔧 PhD Dissertation Planner - Development Setup"
    echo "=============================================="

    log_info "Checking Go installation..."
    check_go

    log_info "Checking for required tools..."
    check_tools

    log_info "Installing dependencies..."
    install_deps

    log_info "Verifying setup..."
    verify_setup

    echo ""
    log_success "Setup completed successfully!"
    echo ""
    echo "🚀 You can now:"
    echo "   make          # Build the project"
    echo "   make clean    # Clean build artifacts"
    echo "   ./scripts/build.sh --help  # See build options"
}

# Run main function
main "$@"
