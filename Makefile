# Simple Makefile for latex-yearly-planner

GO ?= go
BINARY ?= build/plannergen

.PHONY: build pdf clean fmt vet test

# Build the binary
build:
	cd src && $(GO) build -o build/plannergen .

# Generate PDF directly (no helper script)
pdf: build
	cd src && \
	echo "🎯 Generating PDF from: $(CSV)" && \
	echo "📄 Output: $(OUTPUT).pdf" && \
	PLANNER_CSV_FILE="$(CSV)" ./build/plannergen --config "internal/config/base.yaml" --outdir build && \
	echo "🔧 Fixing LaTeX comment issues..." && \
	sed -i '' 's/%\\ColorCircle{/\\ColorCircle{/g' build/monthly.tex || true && \
	sed -i '' 's/%\\hspace{/\\hspace{/g' build/monthly.tex || true && \
	sed -i '' 's/%\\end{center}/\\end{center}/g' build/monthly.tex || true && \
	echo "📚 Compiling PDF..." && \
	cd build && \
	if xelatex -file-line-error -interaction=nonstopmode page_template.tex > xelatex.log 2>&1; then \
		echo "✅ PDF compilation successful"; \
	else \
		echo "⚠️  PDF compilation completed with warnings (check xelatex.log for details)"; \
	fi && \
	cd .. && \
	if [ -f "build/page_template.pdf" ]; then \
		mkdir -p ../output/pdfs ../output/latex ../output/logs && \
		cp "build/page_template.pdf" "../output/pdfs/$(OUTPUT).pdf" && \
		cp "build/page_template.tex" "../output/latex/$(OUTPUT).tex" 2>/dev/null || true && \
		cp "build/page_template.log" "../output/logs/$(OUTPUT).log" 2>/dev/null || true && \
		echo "✅ Created: $(OUTPUT).pdf" && \
		echo "📁 Also saved to: ../output/pdfs/$(OUTPUT).pdf"; \
	else \
		echo "❌ PDF generation failed - check build/xelatex.log for details"; \
		exit 1; \
	fi

# Generate PDF with full dataset and run Go tests
test:
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
	PLANNER_CSV_FILE="../input/data.cleaned.csv" \
	./build/plannergen --config "config/base.yaml" --outdir build && \
	echo "🔧 Fixing LaTeX comment issues..." && \
	sed -i '' 's/%\\ColorCircle{/\\ColorCircle{/g' build/monthly.tex && \
	sed -i '' 's/%\\hspace{/\\hspace{/g' build/monthly.tex && \
	sed -i '' 's/%\\end{center}/\\end{center}/g' build/monthly.tex && \
	echo "📚 Compiling PDF..." && \
	cd build && \
	if xelatex -file-line-error -interaction=nonstopmode page_template.tex > xelatex.log 2>&1; then \
		echo "✅ PDF compilation successful"; \
	else \
		echo "⚠️  PDF compilation completed with warnings (check xelatex.log for details)"; \
	fi && \
	cd .. && \
	if [ -f "build/page_template.pdf" ]; then \
		mkdir -p ../output/pdfs ../output/latex ../output/logs && \
		cp "build/page_template.pdf" "test.pdf" && \
		cp "build/page_template.pdf" "../output/pdfs/test.pdf" && \
		cp "build/page_template.tex" "../output/latex/test.tex" 2>/dev/null || true && \
		cp "build/page_template.log" "../output/logs/test.log" 2>/dev/null || true && \
		echo "✅ Created: test.pdf" && \
		echo "📁 Also saved to: ../output/pdfs/test.pdf"; \
	else \
		echo "❌ PDF generation failed - check build/xelatex.log for details"; \
		exit 1; \
	fi

# Legacy targets for backward compatibility
run: test
run-csv: test
generate: test

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

# Clean output directory
clean-output:
	@echo "🧹 Cleaning output directory..."
	@rm -f output/pdfs/*.pdf
	@rm -f output/latex/*.tex
	@rm -f output/logs/*.log
	@echo "✅ Output directory cleaned"
	@echo "📁 Directory structure preserved"

# Clean both build and release directories
clean-all: clean clean-output

