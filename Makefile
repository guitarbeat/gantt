# Simple Makefile for latex-yearly-planner

.DEFAULT_GOAL := build

GO ?= go
BINARY ?= build/plannergen
OUTDIR ?= src/build

.PHONY: build clean fmt vet

# Build planner PDF (runs tests, generates LaTeX, compiles PDF)
build:
	@echo "🧪 Running Go tests..."
	cd src && go test ./internal/...
	@echo "📄 Generating PDF test..."
	@echo "🎯 Generating PDF from: ../input/data.cleaned.csv"
	@echo "📄 Output: test.pdf"
	@cd src && \
	if [ ! -f "build/plannergen" ]; then \
		echo "🔨 Building plannergen..."; \
		go build -o build/plannergen .; \
	fi && \
	echo "📝 Generating LaTeX..." && \
	PLANNER_SILENT=1 PLANNER_CSV_FILE="../input/data.cleaned.csv" \
	./build/plannergen --config "config/base.yaml,config/page_template.yaml" --outdir build && \
	echo "🔧 Fixing LaTeX comment issues..." && \
	sed -i '' 's/%\\ColorCircle{/\\ColorCircle{/g' build/page_template.tex && \
	sed -i '' 's/%\\hspace{/\\hspace{/g' build/page_template.tex && \
	sed -i '' 's/%\\end{center}/\\end{center}/g' build/page_template.tex && \
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
	# Clean parent directory build artifacts
	rm -f ../*.pdf ../*.tex ../*.aux ../*.log ../*.out ../*.synctex.gz
	# Clean any stray plannergen binaries
	find .. -name "plannergen" -type f -delete 2>/dev/null || true
	# Clean flat output directory
	@echo "🧹 Cleaning output directory..."
	@rm -f output/*.pdf 2>/dev/null || true
	@rm -f output/*.tex 2>/dev/null || true
	@rm -f output/*.log 2>/dev/null || true
	# Also clean legacy subfolders if present
	@rm -f output/pdfs/*.pdf 2>/dev/null || true
	@rm -f output/latex/*.tex 2>/dev/null || true
	@rm -f output/logs/*.log 2>/dev/null || true
	@echo "✅ Output directory cleaned"
	@echo "📁 Directory structure preserved"

