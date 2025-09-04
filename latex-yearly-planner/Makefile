# Simple Makefile for latex-yearly-planner

GO ?= go
BINARY ?= build/plannergen

.PHONY: build pdf clean fmt vet test

# Build the binary
build:
	$(GO) build -o $(BINARY) ./cmd/plannergen

# Generate PDF (simple approach)
pdf:
	./scripts/simple.sh $(CSV) $(OUTPUT)

# Quick test with single task
test:
	./scripts/simple.sh ../input/test_single.csv test

# Generate with triple CSV
demo:
	./scripts/simple.sh ../input/test_triple.csv demo

# Legacy targets for backward compatibility
run: test
run-single: test
run-csv: demo
test-single: test
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

# Clean both build and release directories
clean-all:
	./scripts/generate.sh --clean

