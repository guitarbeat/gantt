# Simple Makefile for latex-yearly-planner

GO ?= go
OUTDIR ?= build
BINARY ?= $(OUTDIR)/plannergen
CFG ?= configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml

.PHONY: build run clean preview fmt vet test list generate all clean-all

build:
	$(GO) build -o $(BINARY) ./cmd/plannergen

# Generate PDFs using the unified script
generate:
	./scripts/generate.sh

# List available CSV files
list:
	./scripts/generate.sh --list

# Generate PDFs for all CSV files
all:
	./scripts/generate.sh --all

# Generate PDF for specific CSV (usage: make csv CSV=1 or make csv CSV=../input/data.cleaned.csv)
csv:
	./scripts/generate.sh $(CSV)

# Legacy targets for backward compatibility
run: generate
preview:
	OUTDIR="$(OUTDIR)" CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/build.sh --preview

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

# Legacy test targets
test-single:
	./scripts/generate.sh ../input/test_single.csv

run-single: test-single

run-csv:
	./scripts/generate.sh ../input/test_triple.csv

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

