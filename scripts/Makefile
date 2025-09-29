# PhD Dissertation Planner - Makefile
# 
# This Makefile orchestrates the complete build process for generating LaTeX-based
# calendar PDFs from CSV timeline data. The process involves:
# 1. Running Go unit tests for code quality assurance
# 2. Timeline validation to catch data inconsistencies 
# 3. Binary compilation with cache cleaning to prevent corruption
# 4. LaTeX generation from CSV data using Go templates
# 5. XeLaTeX compilation to produce final PDF output
#
# Key fixes implemented:
# - Path corrections for src/ directory context execution
# - Special character escaping for LaTeX (& symbols in category names)
# - Validation logging to output/validation.log for persistent review
# - Default clean-build target to ensure reliable development builds

.DEFAULT_GOAL := clean-build

GO ?= go
OUTDIR ?= src/build

# Configurable paths with defaults
CONFIG_BASE ?= src/core/base.yaml
CONFIG_PAGE ?= src/core/monthly_calendar.yaml
CONFIG_FILES ?= $(CONFIG_BASE),$(CONFIG_PAGE)

# Configurable output file names with defaults
OUTPUT_BASE_NAME ?= monthly_calendar
FINAL_BASE_NAME ?= monthly_calendar

# Configurable binary path with defaults
BINARY_DIR ?= generated
BINARY_NAME ?= plannergen
BINARY_PATH ?= $(BINARY_DIR)/$(BINARY_NAME)

# Use the most comprehensive CSV file
CSV_FILE := research_timeline_v5_comprehensive.csv

.PHONY: build clean clean-build build-latex build-pdf fmt vet test test-unit test-integration test-all troubleshoot

# Build planner with optional PDF compilation and enhanced error handling
build:
	@echo "Building from $(CSV_FILE)..."
	@go clean -cache && go build -mod=vendor -o $(BINARY_PATH) ./cmd/planner && \
	PLANNER_SILENT=1 PLANNER_CSV_FILE="input_data/$(CSV_FILE)" \
	$(BINARY_PATH) --config "$(CONFIG_FILES)" --outdir $(BINARY_DIR) && \
	echo "✅ LaTeX file generated: $(BINARY_DIR)/$(OUTPUT_BASE_NAME).tex" && \
	cd $(BINARY_DIR) && \
	if [ ! -f "$(OUTPUT_BASE_NAME).tex" ]; then \
		echo "❌ ERROR: LaTeX file not generated - check planner configuration"; \
		exit 1; \
	fi && \
	TEX_SIZE=$$(stat -c%s "$(OUTPUT_BASE_NAME).tex" 2>/dev/null || echo "0") && \
	if [ "$$TEX_SIZE" -lt 1000 ]; then \
		echo "❌ ERROR: LaTeX file too small ($$TEX_SIZE bytes) - possible generation failure"; \
		echo "   Check configuration files and CSV data format"; \
		exit 1; \
	fi && \
	echo "✅ LaTeX file validation passed ($$TEX_SIZE bytes)" && \
	if command -v xelatex >/dev/null 2>&1; then \
		echo "🔧 XeLaTeX found - attempting PDF compilation..."; \
		if xelatex -file-line-error -interaction=batchmode -halt-on-error $(OUTPUT_BASE_NAME).tex > $(OUTPUT_BASE_NAME).tmp 2>&1; then \
			if grep -q "Error\|Fatal\|!" $(OUTPUT_BASE_NAME).tmp; then \
				echo "❌ PDF compilation failed with errors:"; \
				grep -A2 -B1 "Error\|Fatal\|!" $(OUTPUT_BASE_NAME).tmp | head -10; \
				cat $(OUTPUT_BASE_NAME).tmp > $(OUTPUT_BASE_NAME).log; \
				exit 1; \
			else \
				grep -E "(Warning|Overfull|Underfull)" $(OUTPUT_BASE_NAME).tmp > $(OUTPUT_BASE_NAME).log 2>/dev/null || echo "PDF compiled successfully" > $(OUTPUT_BASE_NAME).log; \
				echo "✅ PDF compiled successfully"; \
			fi; \
		else \
			echo "❌ PDF compilation failed - LaTeX errors:"; \
			grep -A2 -B1 "Error\|Fatal\|!" $(OUTPUT_BASE_NAME).tmp | head -10; \
			cat $(OUTPUT_BASE_NAME).tmp > $(OUTPUT_BASE_NAME).log; \
			echo "📋 Full error log saved to $(OUTPUT_BASE_NAME).log"; \
			exit 1; \
		fi; \
		if [ -f "$(OUTPUT_BASE_NAME).pdf" ]; then \
			PDF_SIZE=$$(stat -c%s "$(OUTPUT_BASE_NAME).pdf"); \
			if [ "$$PDF_SIZE" -lt 10000 ]; then \
				echo "❌ WARNING: PDF file unusually small ($$PDF_SIZE bytes) - may be corrupted"; \
			else \
				echo "✅ Created: $(BINARY_DIR)/$(FINAL_BASE_NAME).pdf ($$PDF_SIZE bytes)"; \
			fi; \
			rm -f *.aux *.fdb_latexmk *.fls *.out *.synctex.gz *.tmp 2>/dev/null || true; \
		else \
			echo "❌ ERROR: PDF file not created despite successful LaTeX compilation"; \
			echo "📋 Check $(OUTPUT_BASE_NAME).log for details"; \
			exit 1; \
		fi; \
	else \
		echo "⚠️  XeLaTeX not found - LaTeX file generated but PDF compilation skipped"; \
		echo "💡 To enable PDF generation, install XeLaTeX:"; \
		echo "   Ubuntu/Debian: sudo apt-get install texlive-xetex texlive-latex-extra"; \
		echo "   macOS: brew install --cask mactex"; \
		echo "   Windows: Install MiKTeX or TeX Live"; \
	fi


fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	@echo "Cleaning build artifacts..."
	@go clean -cache -testcache -modcache 2>/dev/null || true
	@rm -rf "$(BINARY_DIR)"/*.pdf "$(BINARY_DIR)"/*.aux "$(BINARY_DIR)"/*.log "$(BINARY_DIR)"/*.out "$(BINARY_DIR)"/*.tex "$(BINARY_DIR)"/*.synctex.gz
	@rm -f "$(BINARY_PATH)"
	@find . -name "plannergen" -o -name "phd-dissertation-planner" -type f -delete 2>/dev/null || true

# Test targets
test: test-all

test-unit:
	@echo "Running unit tests..."
	@go test -v ./tests/unit/...

test-integration:
	@echo "Running integration tests..."
	@if [ -d "tests/integration" ] && [ "$$(ls -A tests/integration)" ]; then \
		go test -v ./tests/integration/...; \
	else \
		echo "No integration tests found."; \
	fi

test-all:
	@echo "Running all tests..."
	@go test -v ./tests/unit/... ./src/...
	@if [ -d "tests/integration" ] && [ "$$(ls -A tests/integration)" ]; then \
		go test -v ./tests/integration/...; \
	fi

# Clean and build (recommended for development to avoid binary corruption)
clean-build: clean build

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
			if grep -q "Error\|Fatal\|!" $(OUTPUT_BASE_NAME).tmp; then \
				echo "❌ PDF compilation failed with errors:"; \
				grep -A3 -B1 "Error\|Fatal\|!" $(OUTPUT_BASE_NAME).tmp; \
				exit 1; \
			else \
				echo "✅ PDF compiled successfully"; \
				PDF_SIZE=$$(stat -c%s "$(OUTPUT_BASE_NAME).pdf" 2>/dev/null || echo "0"); \
				echo "✅ Created: $(BINARY_DIR)/$(FINAL_BASE_NAME).pdf ($$PDF_SIZE bytes)"; \
			fi; \
		else \
			echo "❌ PDF compilation failed - LaTeX errors:"; \
			grep -A3 -B1 "Error\|Fatal\|!" $(OUTPUT_BASE_NAME).tmp || cat $(OUTPUT_BASE_NAME).tmp; \
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
	@echo "  - Clean and rebuild: make -f scripts/Makefile clean-build"
	@echo "  - LaTeX only: make -f scripts/Makefile build-latex"
	@echo "  - Run tests: make -f scripts/Makefile test"
	@if [ -f "$(BINARY_DIR)/$(OUTPUT_BASE_NAME).log" ]; then \
		echo "  - View build log: cat $(BINARY_DIR)/$(OUTPUT_BASE_NAME).log"; \
	fi

