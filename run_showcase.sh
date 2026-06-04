#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PRESET="full"
PASSTHROUGH_ARGS=()

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --preset)
      if [[ "$#" -lt 2 ]]; then
        echo "Error: --preset requires value: quick|dev|full" >&2
        exit 1
      fi
      PRESET="$2"
      shift 2
      ;;
    --preset=*)
      PRESET="${1#*=}"
      shift
      ;;
    -h|--help)
      cat <<'EOF'
Usage: ./run_showcase.sh [--preset quick|dev|full] [extra go-run flags]

Presets:
  quick  -> records=20,  iters=50,   warmup=5
  dev    -> records=120, iters=500,  warmup=50
  full   -> records=300, iters=1200, warmup=120 (default)

Environment overrides:
  RECORDS ITERS WARMUP INDENT SEED CSV_DIR MD_FILE ENCODED_EXAMPLES_DIR

Notes:
  - Default seed is -1 (program resolves to current unix timestamp).
  - Empty CSV_DIR disables CSV output.
  - CSV_DIR=stdout prints CSV tables to terminal.
  - Extra args are forwarded to `go run .`.
EOF
      exit 0
      ;;
    *)
      PASSTHROUGH_ARGS+=("$1")
      shift
      ;;
  esac
done

case "${PRESET}" in
  quick)
    DEFAULT_RECORDS=20
    DEFAULT_ITERS=50
    DEFAULT_WARMUP=5
    ;;
  dev)
    DEFAULT_RECORDS=120
    DEFAULT_ITERS=500
    DEFAULT_WARMUP=50
    ;;
  full)
    DEFAULT_RECORDS=300
    DEFAULT_ITERS=1200
    DEFAULT_WARMUP=120
    ;;
  *)
    echo "Error: unknown preset '${PRESET}'. Use quick|dev|full." >&2
    exit 1
    ;;
esac

RECORDS="${RECORDS:-${DEFAULT_RECORDS}}"
ITERS="${ITERS:-${DEFAULT_ITERS}}"
WARMUP="${WARMUP:-${DEFAULT_WARMUP}}"
INDENT="${INDENT:-2}"
SEED="${SEED:--1}"
CSV_DIR="${CSV_DIR:-}"
MD_FILE="${MD_FILE:-}"
ENCODED_EXAMPLES_DIR="${ENCODED_EXAMPLES_DIR:-}"

echo "Running TOON showcase with preset: ${PRESET}"
echo "- records=${RECORDS}"
echo "- iters=${ITERS}"
echo "- warmup=${WARMUP}"
echo "- indent=${INDENT}"
echo "- seed=${SEED}"
if [[ -n "${CSV_DIR}" ]]; then
  echo "- csv_dir=${CSV_DIR}"
fi
if [[ -n "${MD_FILE}" ]]; then
  echo "- md_file=${MD_FILE}"
fi
if [[ -n "${ENCODED_EXAMPLES_DIR}" ]]; then
  echo "- encoded_examples_dir=${ENCODED_EXAMPLES_DIR}"
fi

cmd=(
  go run .
  -records "${RECORDS}"
  -iters "${ITERS}"
  -warmup "${WARMUP}"
  -indent "${INDENT}"
  -seed "${SEED}"
)

if [[ -n "${CSV_DIR}" ]]; then
  cmd+=( -csv-dir "${CSV_DIR}" )
fi

if [[ -n "${MD_FILE}" ]]; then
  cmd+=( -md-file "${MD_FILE}" )
fi

if [[ -n "${ENCODED_EXAMPLES_DIR}" ]]; then
  cmd+=( -encoded-examples-dir "${ENCODED_EXAMPLES_DIR}" )
fi

if [[ "${#PASSTHROUGH_ARGS[@]}" -gt 0 ]]; then
  cmd+=( "${PASSTHROUGH_ARGS[@]}" )
fi

(
  cd "${SCRIPT_DIR}"
  "${cmd[@]}"
)
