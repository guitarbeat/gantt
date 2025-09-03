# Simple Makefile for latex-yearly-planner

GO ?= go
OUTDIR ?= build
BINARY ?= $(OUTDIR)/plannergen
CFG ?= configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml

.PHONY: build run clean preview fmt vet test test-single run-csv run-single

build:
	$(GO) build -o $(BINARY) ./cmd/plannergen

run:
	OUTDIR="$(OUTDIR)" CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/build.sh

preview:
	OUTDIR="$(OUTDIR)" CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/build.sh --preview

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test-single:
	OUTDIR="$(OUTDIR)" CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/run_single.sh

run-single: test-single

run-csv:
	OUTDIR="$(OUTDIR)" CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/run_with_csv.sh

clean:
	rm -rf "$(OUTDIR)"/*.pdf "$(OUTDIR)"/*.aux "$(OUTDIR)"/*.log "$(OUTDIR)"/*.out "$(OUTDIR)"/*.tex "$(OUTDIR)"/*.synctex.gz
	rm -f "$(BINARY)"

