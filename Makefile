# Simple Makefile for latex-yearly-planner

.DEFAULT_GOAL := build

GO ?= go
BINARY ?= build/plannergen
OUTDIR ?= src/build

# Configurable paths with defaults
CONFIG_BASE ?= config/base.yaml
CONFIG_PAGE ?= config/page_template.yaml
CONFIG_FILES ?= $(CONFIG_BASE),$(CONFIG_PAGE)

# Configurable output file names with defaults
OUTPUT_BASE_NAME ?= page_template
FINAL_BASE_NAME ?= test

# Find the first CSV file in the input directory
CSV_FILE := $(shell ls input/*.csv 2>/dev/null | head -1 | xargs basename)

.PHONY: build clean fmt vet test

# Build planner PDF (runs tests, generates LaTeX, compiles PDF)
build:
	@echo "🧪 Running Go tests..."
	cd src && unset PLANNER_CSV_FILE && go test ./tests/unit/...
	@echo "📄 Generating PDF test..."
	@echo "🎯 Generating PDF from: input/$(CSV_FILE)"
	@echo "📄 Output: $(FINAL_BASE_NAME).pdf"
	@cd src && \
	if [ ! -f "build/plannergen" ]; then \
		echo "🔨 Building plannergen..."; \
		go build -o build/plannergen .; \
	fi && \
	echo "📝 Generating LaTeX..." && \
	PLANNER_SILENT=1 PLANNER_CSV_FILE="../input/$(CSV_FILE)" \
	./build/plannergen --config "$(CONFIG_FILES)" --outdir build && \
	echo "📚 Compiling PDF..." && \
	cd build && \
	if xelatex -file-line-error -interaction=nonstopmode $(OUTPUT_BASE_NAME).tex > $(OUTPUT_BASE_NAME).log 2>&1; then \
		echo "✅ PDF compilation successful"; \
	else \
		echo "⚠️  PDF compilation completed with warnings (check xelatex.log for details)"; \
	fi && \
	cd .. && \
	if [ -f "build/$(OUTPUT_BASE_NAME).pdf" ]; then \
		mkdir -p ../output && \
		cp "build/$(OUTPUT_BASE_NAME).pdf" "$(FINAL_BASE_NAME).pdf" && \
		cp "build/$(OUTPUT_BASE_NAME).pdf" "../output/$(FINAL_BASE_NAME).pdf" && \
		cp "build/$(OUTPUT_BASE_NAME).tex" "../output/$(FINAL_BASE_NAME).tex" 2>/dev/null || true && \
		cp "build/$(OUTPUT_BASE_NAME).log" "../output/$(FINAL_BASE_NAME).log" 2>/dev/null || true && \
		echo "🧹 Cleaning up auxiliary files from output..." && \
		cd ../output && rm -f *.aux *.fdb_latexmk *.fls *.out *.synctex.gz 2>/dev/null || true && \
		cd ../src && \
		echo "✅ Created: $(FINAL_BASE_NAME).pdf" && \
		echo "📁 Also saved to: ../output/$(FINAL_BASE_NAME).pdf"; \
	else \
		echo "❌ PDF generation failed - check build/$(OUTPUT_BASE_NAME).log for details"; \
		exit 1; \
	fi

test:
	@echo "🧪 Running Go tests..."
	cd src && unset PLANNER_CSV_FILE && go test ./tests/unit/...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	# Clean build directory
	rm -rf "$(OUTDIR)"/*.pdf "$(OUTDIR)"/*.aux "$(OUTDIR)"/*.log "$(OUTDIR)"/*.out "$(OUTDIR)"/*.tex "$(OUTDIR)"/*.synctex.gz
	rm -f "$(BINARY)"
	# Clean src directory build artifacts
	@echo "🧹 Cleaning src directory..."
	@rm -f src/*.pdf src/*.tex src/*.aux src/*.log src/*.out src/*.synctex.gz src/*.fdb_latexmk src/*.fls src/coverage.out src/debug.log src/test.out 2>/dev/null || true
	@echo "✅ Src directory cleaned"
	# Clean parent directory build artifacts
	rm -f ../*.pdf ../*.tex ../*.aux ../*.log ../*.out ../*.synctex.gz
	# Clean any stray plannergen binaries
	find .. -name "plannergen" -type f -delete 2>/dev/null || true
	# Clean flat output directory
	@echo "🧹 Cleaning output directory..."
	@rm -f output/*.pdf output/*.tex output/*.log output/*.aux output/*.fdb_latexmk output/*.fls output/*.out output/*.synctex.gz 2>/dev/null || true
	@echo "✅ Output directory cleaned"
	@echo "📁 Directory structure preserved"

