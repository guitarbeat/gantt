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

.PHONY: help setup quality ci dev dev-air build clean install lint fmt run organize status test-coverage bench hooks build-latex build-pdf troubleshoot release release-dry release-snapshot completion completion-bash completion-zsh completion-fish completion-powershell docker-dev docker-build docker-run docker-clean

# Cross-platform file size function
define get_file_size
$(shell stat -c%s $(1) 2>/dev/null || stat -f%z $(1) 2>/dev/null || echo "0")
endef

# Default target
help:
	@echo "PhD Dissertation Planner - Available Commands:"
	@echo ""
	@echo "🚀 Main Commands:"
	@echo "  make setup          - Initialize project (deps, hooks, organize)"
	@echo "  make quality        - Code quality checks (fmt, lint, test + coverage)"
	@echo "  make build          - Build binary and generate PDF"
	@echo "  make run            - Build and run with default config (quiet)"
	@echo "  make run-verbose    - Build and run with verbose output"
	@echo "  make clean          - Clean all build artifacts"
	@echo ""
	@echo "🔧 Advanced Commands:"
	@echo "  make ci             - Full CI pipeline (clean + quality + build)"
	@echo "  make dev            - Development workflow (clean + quality + build + run)"
	@echo "  make dev-air        - Start hot-reloading development server with air"
	@echo "  make dev-verbose    - Development workflow with verbose output"
	@echo "  make bench          - Run performance benchmarks"
	@echo "  make status         - Show project organization status"
	@echo "  make troubleshoot   - Run build system diagnostics"
	@echo ""
	@echo "🚀 Release Commands:"
	@echo "  make release        - Create and publish a new release with goreleaser"
	@echo "  make release-dry    - Test release process without publishing"
	@echo "  make release-snapshot - Create snapshot release for testing"
	@echo ""
	@echo "🔧 Shell Completion:"
	@echo "  make completion     - Show available completion options"
	@echo "  make completion-bash - Generate and install bash completion"
	@echo "  make completion-zsh  - Generate and install zsh completion"
	@echo "  make completion-fish - Generate and install fish completion"
	@echo ""
	@echo "🐳 Docker Development:"
	@echo "  make docker-dev     - Start development environment with Docker"
	@echo "  make docker-build   - Build Docker development image"
	@echo "  make docker-run     - Run commands in Docker container"
	@echo "  make docker-clean   - Clean Docker containers and images"
	@echo ""

# ==================== Consolidated Commands ====================

# Initialize project - install dependencies, setup hooks, organize files
setup: install hooks organize
	@echo "🎯 Project setup complete! Ready for development."

# Run code quality checks - format, lint, test, and coverage
quality: fmt lint test-coverage
	@echo "✅ Code quality checks passed!"

# Full CI pipeline - clean, quality checks, and build
ci: clean quality build
	@echo "🚀 CI pipeline completed successfully!"

# Development workflow - clean, quality, build, and run (quiet)
dev: clean quality build run
	@echo "💻 Development workflow complete!"

# Development workflow with verbose output
dev-verbose: clean quality build run-verbose
	@echo "💻 Development workflow complete!"

# Start hot-reloading development server with air
dev-air:
	@echo "🔥 Starting hot-reloading development server..."
	@if ! which air > /dev/null 2>&1; then \
		echo "📦 Installing air for hot reloading..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@air

# ==================== Individual Commands ====================

# Build planner with optional PDF compilation and enhanced error handling
build: build-pdf

# Run tests with coverage (used by quality command)
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -race -coverprofile=coverage.txt ./... | grep -E "(PASS|FAIL|RUN|coverage:)" || true
	@go tool cover -html=coverage.txt -o coverage.html > /dev/null 2>&1
	@echo "✅ Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "📊 Running benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "✅ Benchmarks completed"

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
	@if ! which golangci-lint > /dev/null 2>&1; then \
		echo "📦 Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest > /dev/null 2>&1; \
	fi
	@PATH=$$PATH:$$(go env GOPATH)/bin golangci-lint run ./... --color=always
	@echo "✅ Lint complete!"

# Format code
fmt:
	@echo "✨ Formatting code..."
	@gofmt -w .
	@PATH=$$PATH:$$(go env GOPATH)/bin goimports -w src/
	@echo "✅ Format complete!"

# Build and run with default config (quiet)
run: build
	@echo "🚀 Running planner..."
	@PLANNER_SILENT=1 ./$(BINARY_PATH)

# Build and run with verbose output
run-verbose: build
	@echo "🚀 Running planner (verbose)..."
	@./$(BINARY_PATH)

# Install pre-commit hooks
hooks:
	@echo "🪝 Installing pre-commit hooks..."
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit install; \
	else \
		echo "📦 Installing pre-commit..."; \
		if command -v pip >/dev/null 2>&1; then \
			pip install pre-commit; \
		elif command -v pip3 >/dev/null 2>&1; then \
			pip3 install pre-commit; \
		else \
			echo "❌ pip not found. Please install pre-commit manually."; \
			exit 1; \
		fi; \
		pre-commit install; \
	fi
	@echo "✅ Hooks installed!"

# Organize project files
organize:
	@echo "🧹 Organizing project files..."
	@if [ -f "./scripts/maintenance/cleanup_and_organize.sh" ]; then \
		./scripts/maintenance/cleanup_and_organize.sh; \
	elif [ -f "./scripts/maintenance/cleanup_and_organize.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File ./scripts/maintenance/cleanup_and_organize.ps1; \
	else \
		echo "❌ No cleanup script found"; \
		exit 1; \
	fi
	@echo "✅ Organization complete!"

# Show project status
status:
	@echo "📊 Project Status:"
	@if [ -f "./scripts/maintenance/cleanup_and_organize.sh" ]; then \
		./scripts/maintenance/cleanup_and_organize.sh --status; \
	elif [ -f "./scripts/maintenance/cleanup_and_organize.ps1" ]; then \
		powershell -ExecutionPolicy Bypass -File ./scripts/maintenance/cleanup_and_organize.ps1 -Status; \
	else \
		echo "❌ No cleanup script found"; \
		exit 1; \
	fi

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
				PDF_SIZE=$$(call get_file_size,"$(OUTPUT_BASE_NAME).pdf"); \
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
		BINARY_SIZE=$$(call get_file_size,"$(BINARY_PATH)"); \
		echo "  - Binary: ✅ $(BINARY_PATH) ($$BINARY_SIZE bytes)"; \
	else \
		echo "  - Binary: ❌ Not found at $(BINARY_PATH)"; \
	fi
	@if [ -f "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex" ]; then \
		TEX_SIZE=$$(call get_file_size,"$(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex"); \
		echo "  - LaTeX: ✅ $(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex ($$TEX_SIZE bytes)"; \
	else \
		echo "  - LaTeX: ❌ Not found"; \
	fi
	@if [ -f "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).pdf" ]; then \
		PDF_SIZE=$$(call get_file_size,"$(BINARY_DIR)/$(OUTPUT_BASE_NAME).pdf"); \
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

# ==================== Release Commands ====================

# Create and publish a new release with goreleaser
release:
	@echo "🚀 Creating and publishing release..."
	@if ! which goreleaser > /dev/null 2>&1; then \
		echo "📦 Installing goreleaser..."; \
		go install github.com/goreleaser/goreleaser@latest; \
	fi
	@goreleaser release --clean

# Test release process without publishing
release-dry:
	@echo "🧪 Testing release process (dry run)..."
	@if ! which goreleaser > /dev/null 2>&1; then \
		echo "📦 Installing goreleaser..."; \
		go install github.com/goreleaser/goreleaser@latest; \
	fi
	@goreleaser release --clean --skip-publish

# Create snapshot release for testing
release-snapshot:
	@echo "📸 Creating snapshot release..."
	@if ! which goreleaser > /dev/null 2>&1; then \
		echo "📦 Installing goreleaser..."; \
		go install github.com/goreleaser/goreleaser@latest; \
	fi
	@goreleaser release --clean --snapshot

# ==================== Shell Completion Commands ====================

# Show available completion options
completion:
	@echo "🔧 Shell Completion Setup Instructions:"
	@echo ""
	@echo "Available shells: bash, zsh, fish, powershell"
	@echo ""
	@echo "To generate completion scripts:"
	@echo "  make completion-bash    # Generate bash completion"
	@echo "  make completion-zsh     # Generate zsh completion"
	@echo "  make completion-fish    # Generate fish completion"
	@echo "  make completion-powershell # Generate PowerShell completion"
	@echo ""
	@echo "Manual installation:"
	@echo "  ./plannergen completion <shell> > completion.<shell>"
	@echo "  Source the generated file in your shell profile"

# Generate and install bash completion
completion-bash:
	@echo "🐚 Generating bash completion..."
	@go build -o plannergen ./cmd/planner
	@./plannergen completion bash > plannergen.bash
	@echo "✅ Generated plannergen.bash"
	@echo "To install: source plannergen.bash >> ~/.bashrc"

# Generate and install zsh completion
completion-zsh:
	@echo "🐚 Generating zsh completion..."
	@go build -o plannergen ./cmd/planner
	@./plannergen completion zsh > plannergen.zsh
	@echo "✅ Generated plannergen.zsh"
	@echo "To install: source plannergen.zsh >> ~/.zshrc"

# Generate and install fish completion
completion-fish:
	@echo "🐚 Generating fish completion..."
	@go build -o plannergen ./cmd/planner
	@./plannergen completion fish > plannergen.fish
	@echo "✅ Generated plannergen.fish"
	@echo "To install: cp plannergen.fish ~/.config/fish/completions/"

# Generate PowerShell completion
completion-powershell:
	@echo "🐚 Generating PowerShell completion..."
	@go build -o plannergen ./cmd/planner
	@./plannergen completion powershell > plannergen.ps1
	@echo "✅ Generated plannergen.ps1"
	@echo "To install: Add to PowerShell profile"

# ==================== Docker Development Commands ====================

# Start development environment with Docker
docker-dev:
	@echo "🐳 Starting Docker development environment..."
	@docker-compose up dev

# Build Docker development image
docker-build:
	@echo "🔨 Building Docker development image..."
	@docker-compose build dev

# Run commands in Docker container
docker-run:
	@echo "🐳 Running command in Docker container..."
	@docker-compose run --rm dev $(CMD)

# Clean Docker containers and images
docker-clean:
	@echo "🧹 Cleaning Docker containers and images..."
	@docker-compose down --volumes --remove-orphans
	@docker system prune -f
	@docker image prune -f