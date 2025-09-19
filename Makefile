# Simple Makefile for latex-yearly-planner

GO ?= go
BINARY ?= build/plannergen

.PHONY: build pdf clean fmt vet test

# Build the binary
build:
	cd src && $(GO) build -o build/plannergen ./cmd/plannergen

# Generate PDF (simple approach)
pdf:
	cd src && ./scripts/simple.sh $(CSV) $(OUTPUT)

# Generate PDF with full dataset
test:
	cd src && ./scripts/simple.sh ../input/data.cleaned.csv test

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
	./scripts/clean_output.sh

# Clean both build and release directories
clean-all:
	./scripts/generate.sh --clean

