#!/usr/bin/env bash

set -eo pipefail

if [ -z "$PLANNERGEN_BINARY" ]; then
  export GO_CMD="go run ./cmd/plannergen"
else
  export GO_CMD="$PLANNERGEN_BINARY"
  echo "Building using plannergen binary at \"${PLANNERGEN_BINARY}\""
fi

OUTDIR=${OUTDIR:-build}

if [ -z "$PREVIEW" ]; then
  eval $GO_CMD --config "${CFG}" --outdir "$OUTDIR"
else
  eval $GO_CMD --preview --config "${CFG}" --outdir "$OUTDIR"
fi



# Use the fixed filename we generate in app.go
nakedname="proposal-timeline"

if [ -n "${TRANSLATION}" ]; then
  python3 translate.py ${TRANSLATION}
fi

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
    -output-directory="./${OUTDIR}" \
    "${OUTDIR}/${nakedname}.tex"
done

if [ -n "${NAME}" ]; then
  cp "${OUTDIR}/${nakedname}.pdf" "${NAME}.pdf"
  echo "created ${NAME}.pdf"
else
  cp "${OUTDIR}/${nakedname}.pdf" "${nakedname}.pdf"
  echo "created ${nakedname}.pdf"
fi
