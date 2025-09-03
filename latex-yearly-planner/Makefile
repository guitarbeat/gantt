# Simple Makefile for latex-yearly-planner

GO ?= go
BINARY ?= build/plannergen
CFG ?= configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml

.PHONY: build run clean preview fmt vet test

build:
	$(GO) build -o $(BINARY) ./cmd/plannergen

run:
	CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" scripts/single.sh

preview:
	CFG="$(CFG)" PLANNERGEN_BINARY="$(BINARY)" PREVIEW=1 scripts/single.sh

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	rm -rf build/*.pdf build/*.aux build/*.log build/*.out build/*.tex
	rm -f $(BINARY)
	rm -rf build/plannergen

