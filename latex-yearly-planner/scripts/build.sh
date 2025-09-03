#!/usr/bin/env bash

set -eo pipefail

usage() {
  cat <<'USAGE'
Usage: scripts/build.sh [options]

Environment variables:
  CFG                  Comma-separated list of config files (required)
  PLANNERGEN_BINARY    Path to compiled generator (optional; if empty, runs `go run`)
  PREVIEW              If set (non-empty), passes --preview to the generator
  PASSES               Number of XeLaTeX passes (default: 1)
  NAME                 Output PDF name (default: based on last config)

Examples:
  CFG="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" scripts/build.sh
  PREVIEW=1 CFG="configs/base.yaml,configs/page_template.yaml,configs/planner_config.yaml" scripts/build.sh
USAGE
}

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
  usage
  exit 0
fi

if [ -z "${CFG}" ]; then
  echo "error: CFG must be set (comma-separated config files)" >&2
  exit 2
fi

if [ -z "$PLANNERGEN_BINARY" ]; then
  export GO_CMD="go run ./cmd/plannergen"
else
  export GO_CMD="$PLANNERGEN_BINARY"
  echo "Building using plannergen binary at \"${PLANNERGEN_BINARY}\""
fi

if [ -z "$PREVIEW" ]; then
  eval $GO_CMD --config "${CFG}"
else
  eval $GO_CMD --preview --config "${CFG}"
fi

nakedname=$(echo "${CFG}" | rev | cut -d, -f1 | cut -d'/' -f 1 | cut -d'.' -f 2-99 | rev)

_passes=(1)
if [[ -n "${PASSES}" ]]; then
  # shellcheck disable=SC2207
  _passes=($(seq 1 "${PASSES}"))
fi

for _ in "${_passes[@]}"; do
  xelatex \
    -file-line-error \
    -interaction=nonstopmode \
    -synctex=1 \
    -output-directory=./build \
    "build/${nakedname}.tex"
done

if [ -n "${NAME}" ]; then
  cp "build/${nakedname}.pdf" "${NAME}.pdf"
  echo "created ${NAME}.pdf"
else
  cp "build/${nakedname}.pdf" "${nakedname}.pdf"
  echo "created ${nakedname}.pdf"
fi
