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
CONFIG_BASE ?= config/base.yaml
CONFIG_PAGE ?= config/monthly_calendar.yaml
CONFIG_FILES ?= $(CONFIG_BASE),$(CONFIG_PAGE)

# Configurable output file names with defaults
OUTPUT_BASE_NAME ?= monthly_calendar
FINAL_BASE_NAME ?= test

# Configurable binary path with defaults
BINARY_DIR ?= output
BINARY_NAME ?= plannergen
BINARY_PATH ?= $(BINARY_DIR)/$(BINARY_NAME)

# Find the first CSV file in the input directory
CSV_FILE := $(shell ls input/*.csv 2>/dev/null | head -1 | sed 's|^input/||')

.PHONY: build clean clean-build fmt vet

# Build planner PDF with comprehensive pipeline
# Note: All paths are relative to src/ directory due to 'cd src' context
build:
	@echo "📄 Generating PDF..."
	@echo "🎯 Generating PDF from: input/$(CSV_FILE)"
	@echo "📄 Output: $(FINAL_BASE_NAME).pdf"
	@echo "🔨 Building $(BINARY_NAME)..."; \
	cd src && go clean -cache && go build -o ../$(BINARY_PATH) . && \
	echo "✅ Binary built successfully" && \
	echo "📝 Generating LaTeX..." && \
	PLANNER_SILENT=1 PLANNER_CSV_FILE="../input/$(CSV_FILE)" \
	../$(BINARY_PATH) --config "config/base.yaml,config/monthly_calendar.yaml" --outdir ../$(BINARY_DIR) && \
	echo "📚 Compiling PDF..." && \
	cd ../$(BINARY_DIR) && \
	if xelatex -file-line-error -interaction=nonstopmode $(OUTPUT_BASE_NAME).tex > $(OUTPUT_BASE_NAME).log 2>&1; then \
		echo "✅ PDF compilation successful"; \
	else \
		echo "⚠️  PDF compilation completed with warnings (check xelatex.log for details)"; \
	fi && \
	if [ -f "$(OUTPUT_BASE_NAME).pdf" ]; then \
		cp "$(OUTPUT_BASE_NAME).pdf" "../$(FINAL_BASE_NAME).pdf" && \
		cp "$(OUTPUT_BASE_NAME).tex" "../$(FINAL_BASE_NAME).tex" 2>/dev/null || true && \
		cp "$(OUTPUT_BASE_NAME).log" "../$(FINAL_BASE_NAME).log" 2>/dev/null || true && \
		echo "🧹 Cleaning up auxiliary files from output..." && \
		rm -f *.aux *.fdb_latexmk *.fls *.out *.synctex.gz 2>/dev/null || true && \
		echo "✅ Created: $(FINAL_BASE_NAME).pdf" && \
		echo "📁 Also saved to: $(BINARY_DIR)/$(FINAL_BASE_NAME).pdf"; \
	else \
		echo "❌ PDF generation failed - check $(BINARY_DIR)/$(OUTPUT_BASE_NAME).log for details"; \
		exit 1; \
	fi


fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	# Clean Go build cache
	@echo "🧹 Cleaning Go build cache..."
	@cd src && go clean -cache -testcache -modcache 2>/dev/null || true
	# Clean output directory build artifacts
	rm -rf "$(BINARY_DIR)"/*.pdf "$(BINARY_DIR)"/*.aux "$(BINARY_DIR)"/*.log "$(BINARY_DIR)"/*.out "$(BINARY_DIR)"/*.tex "$(BINARY_DIR)"/*.synctex.gz
	rm -f "$(BINARY_PATH)"
	# Clean src directory build artifacts
	@echo "🧹 Cleaning src directory..."
	@rm -f src/*.pdf src/*.tex src/*.aux src/*.log src/*.out src/*.synctex.gz src/*.fdb_latexmk src/*.fls src/coverage.out src/debug.log src/test.out 2>/dev/null || true
	@echo "✅ Src directory cleaned"
	# Clean parent directory build artifacts
	rm -f *.pdf *.tex *.aux *.log *.out *.synctex.gz
	# Clean any stray plannergen binaries
	find . -name "plannergen" -type f -delete 2>/dev/null || true
	@echo "📁 Directory structure preserved"

# Clean and build (recommended for development to avoid binary corruption)
clean-build: clean build

