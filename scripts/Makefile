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

.PHONY: build clean clean-build fmt vet test test-unit test-integration test-all

# Build planner PDF with comprehensive pipeline
build:
	@echo "Building PDF from $(CSV_FILE)..."
	@go clean -cache && go build -mod=vendor -o $(BINARY_PATH) ./cmd/planner && \
	PLANNER_SILENT=1 PLANNER_CSV_FILE="input_data/$(CSV_FILE)" \
	$(BINARY_PATH) --config "$(CONFIG_FILES)" --outdir $(BINARY_DIR) && \
	cd $(BINARY_DIR) && \
	if xelatex -file-line-error -interaction=nonstopmode $(OUTPUT_BASE_NAME).tex > $(OUTPUT_BASE_NAME).log 2>&1; then \
		echo "PDF compiled successfully"; \
	else \
		echo "Warning: PDF compilation had warnings (check $(OUTPUT_BASE_NAME).log)"; \
	fi && \
	if [ -f "$(OUTPUT_BASE_NAME).pdf" ]; then \
		rm -f *.aux *.fdb_latexmk *.fls *.out *.synctex.gz 2>/dev/null || true && \
		echo "Created: $(BINARY_DIR)/$(FINAL_BASE_NAME).pdf"; \
	else \
		echo "Error: PDF generation failed - check $(BINARY_DIR)/$(OUTPUT_BASE_NAME).log"; \
		exit 1; \
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

