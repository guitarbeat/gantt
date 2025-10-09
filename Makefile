# PhD Dissertation Planner - Makefile
#
# This Makefile orchestrates the complete build process for generating LaTeX-based
# calendar PDFs from CSV timeline data.

.DEFAULT_GOAL := help

GO ?= go
BINARY_DIR ?= generated
BINARY_NAME ?= plannergen
BINARY_PATH ?= $(BINARY_DIR)/$(BINARY_NAME)

# Configurable paths with defaults
CONFIG_BASE ?= src/core/base.yaml
CONFIG_PAGE ?= src/core/monthly_calendar.yaml
CONFIG_FILES ?= $(CONFIG_BASE),$(CONFIG_PAGE)

# Configurable output file names with defaults
OUTPUT_BASE_NAME ?= monthly_calendar
FINAL_BASE_NAME ?= monthly_calendar

# Use the most comprehensive CSV file
CSV_FILE := research_timeline_v5_comprehensive.csv

.PHONY: help build test clean install lint fmt run organize status test-coverage lint-basic hooks check build-latex build-pdf troubleshoot

# Default target
help:
	@echo "PhD Dissertation Planner - Available Commands:"
	@echo ""
	@echo "  make build          - Build the planner binary and generate PDF"
	@echo "  make test           - Run all tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install dependencies"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo "  make run            - Build and run with default config"
	@echo "  make hooks          - Install pre-commit hooks"
	@echo "  make organize       - Clean up and organize project files"
	@echo "  make status         - Show project organization status"
	@echo "  make build-latex    - Build LaTeX only"
	@echo "  make build-pdf      - Build PDF from LaTeX"
	@echo "  make troubleshoot   - Run build system diagnostics"
	@echo ""

# Build planner with optional PDF compilation and enhanced error handling
build: build-pdf

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...
	@echo "✅ Tests complete!"

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@go clean -cache -testcache -modcache 2>/dev/null || true
	@rm -rf "$(BINARY_DIR)"
	@rm -f plannergen.exe plannergen
	@rm -f coverage.txt coverage.html
	@find . -name "plannergen" -o -name "phd-dissertation-planner" -type f -delete 2>/dev/null || true
	@echo "✅ Clean complete!"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed!"

# Run linters
lint:
	@echo "🔍 Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run ./...
	@echo "✅ Lint complete!"

# Run basic linters (without golangci-lint)
lint-basic:
	@echo "🔍 Running basic linters..."
	@go vet ./...
	@gofmt -l .
	@echo "✅ Basic lint complete!"

# Format code
fmt:
	@echo "✨ Formatting code..."
	@gofmt -w .
	@goimports -w src/
	@echo "✅ Format complete!"

# Build and run with default config
run: build
	@echo "🚀 Running planner..."
	@./$(BINARY_PATH)

# Install pre-commit hooks
hooks:
	@echo "🪝 Installing pre-commit hooks..."
	@pre-commit install
	@echo "✅ Hooks installed!"

# Run pre-commit on all files
check:
	@echo "✅ Running pre-commit checks..."
	@pre-commit run --all-files

# Organize project files
organize:
	@echo "🧹 Organizing project files..."
	@./scripts/cleanup_and_organize.sh
	@echo "✅ Organization complete!"

# Show project status
status:
	@echo "📊 Project Status:"
	@./scripts/cleanup_and_organize.sh --status

# Build LaTeX only (without PDF compilation)
build-latex:
	@echo "Building LaTeX from $(CSV_FILE)..."
	@go clean -cache && go build -mod=vendor -o $(BINARY_PATH) ./cmd/planner && \
	PLANNER_SILENT=1 PLANNER_CSV_FILE="input_data/$(CSV_FILE)" \
	$(BINARY_PATH) --config "$(CONFIG_FILES)" --outdir $(BINARY_DIR) && \
	echo "LaTeX file generated: $(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex"

# Force PDF compilation with enhanced error reporting (fails if XeLaTeX not available)
build-pdf: build-latex
	@echo "🔧 Attempting PDF compilation..."
	@cd $(BINARY_DIR) && \
	if command -v xelatex >/dev/null 2>&1; then \
		echo "✅ XeLaTeX found - compiling PDF..."; \
		if xelatex -file-line-error -interaction=batchmode -halt-on-error $(OUTPUT_BASE_NAME).tex > $(OUTPUT_BASE_NAME).tmp 2>&1; then \
			if grep -q "Error\\|Fatal\\|!" $(OUTPUT_BASE_NAME).tmp; then \
				echo "❌ PDF compilation failed with errors:"; \
				grep -A3 -B1 "Error\\|Fatal\\|!" $(OUTPUT_BASE_NAME).tmp; \
				exit 1; \
			else \
				echo "✅ PDF compiled successfully"; \
				PDF_SIZE=$$(stat -c%s "$(OUTPUT_BASE_NAME).pdf" 2>/dev/null || echo "0"); \
				echo "✅ Created: $(BINARY_DIR)/$(FINAL_BASE_NAME).pdf ($$PDF_SIZE bytes)"; \
			fi; \
		else \
			echo "❌ PDF compilation failed - LaTeX errors:"; \
			grep -A3 -B1 "Error\\|Fatal\\|!" $(OUTPUT_BASE_NAME).tmp || cat $(OUTPUT_BASE_NAME).tmp; \
			exit 1; \
		fi; \
	else \
		echo "❌ ERROR: XeLaTeX not found. Install with:"; \
		echo "   Ubuntu/Debian: sudo apt-get install texlive-xetex texlive-latex-extra"; \
		echo "   macOS: brew install --cask mactex"; \
		echo "   Windows: Install MiKTeX or TeX Live"; \
		exit 1; \
	fi

# Troubleshooting and diagnostics
troubleshoot:
	@echo "🔍 PhD Dissertation Planner - Build System Diagnostics"
	@echo "========================================================"
	@echo ""
	@echo "📋 Environment Information:"
	@echo "  - Go version: $$(go version 2>/dev/null || echo 'Go not found')"
	@echo "  - XeLaTeX: $$(command -v xelatex >/dev/null 2>&1 && echo 'Available' || echo 'Not found')"
	@echo "  - CSV file: $(CSV_FILE) $$([ -f 'input_data/$(CSV_FILE)' ] && echo '✅' || echo '❌ Missing')"
	@echo "  - Config files: $(CONFIG_FILES)"
	@echo ""
	@echo "📁 Directory Structure:"
	@ls -la $(BINARY_DIR)/ 2>/dev/null | head -10 || echo "  $(BINARY_DIR)/ directory not found"
	@echo ""
	@echo "🔧 Build Status:"
	@if [ -f "$(BINARY_PATH)" ]; then \
		echo "  - Binary: ✅ $(BINARY_PATH) ($$(stat -c%s $(BINARY_PATH) 2>/dev/null || echo '0') bytes)"; \
	else \
		echo "  - Binary: ❌ Not found at $(BINARY_PATH)"; \
	fi
	@if [ -f "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex" ]; then \
		TEX_SIZE=$$(stat -c%s "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex" 2>/dev/null || echo "0"); \
		echo "  - LaTeX: ✅ $(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex ($$TEX_SIZE bytes)"; \
	else \
		echo "  - LaTeX: ❌ Not found"; \
	fi
	@if [ -f "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).pdf" ]; then \
		PDF_SIZE=$$(stat -c%s "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).pdf" 2>/dev/null || echo "0"); \
		echo "  - PDF: ✅ $(BINARY_DIR)/$(OUTPUT_BASE_NAME).pdf ($$PDF_SIZE bytes)"; \
	else \
		echo "  - PDF: ❌ Not found"; \
	fi
	@echo ""
	@echo "🚀 Quick Actions:"
	@echo "  - Clean and rebuild: make clean build"
	@echo "  - LaTeX only: make build-latex"
	@echo "  - Run tests: make test"
	@if [ -f "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).log" ]; then \
		echo "  - View build log: cat $(BINARY_DIR)/$(OUTPUT_BASE_NAME).log"; \
	fi