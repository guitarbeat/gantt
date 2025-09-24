# Simple Makefile for latex-yearly-planner

.DEFAULT_GOAL := build

GO ?= go
BINARY ?= build/plannergen
OUTDIR ?= src/build

.PHONY: build clean fmt vet

# Build planner PDF (runs tests, generates LaTeX, compiles PDF)
build:
	@echo "🧪 Running Go tests..."
	cd src && unset PLANNER_CSV_FILE && go test ./tests/unit/...
	@echo "📄 Generating PDF test..."
	@echo "🎯 Generating PDF from: ../input/Research Timeline v5 - Comprehensive.csv"
	@echo "📄 Output: test.pdf"
	@cd src && \
	if [ ! -f "build/plannergen" ]; then \
		echo "🔨 Building plannergen..."; \
		go build -o build/plannergen .; \
	fi && \
	echo "📝 Generating LaTeX..." && \
	PLANNER_SILENT=1 PLANNER_CSV_FILE="../input/Research Timeline v5 - Comprehensive.csv" \
	./build/plannergen --config "config/base.yaml,config/page_template.yaml" --outdir build && \
	echo "📚 Compiling PDF..." && \
	cd build && \
	if xelatex -file-line-error -interaction=nonstopmode page_template.tex > xelatex.log 2>&1; then \
		echo "✅ PDF compilation successful"; \
	else \
		echo "⚠️  PDF compilation completed with warnings (check xelatex.log for details)"; \
	fi && \
	cd .. && \
	if [ -f "build/page_template.pdf" ]; then \
		mkdir -p ../output && \
		cp "build/page_template.pdf" "test.pdf" && \
		cp "build/page_template.pdf" "../output/test.pdf" && \
		cp "build/page_template.tex" "../output/test.tex" 2>/dev/null || true && \
		cp "build/page_template.log" "../output/test.log" 2>/dev/null || true && \
		echo "🧹 Cleaning up auxiliary files from output..." && \
		cd ../output && rm -f *.aux *.fdb_latexmk *.fls *.out *.synctex.gz 2>/dev/null || true && \
		cd ../src && \
		echo "✅ Created: test.pdf" && \
		echo "📁 Also saved to: ../output/test.pdf"; \
	else \
		echo "❌ PDF generation failed - check build/xelatex.log for details"; \
		exit 1; \
	fi

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

