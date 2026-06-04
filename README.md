# TOON vs Multi-Format Showcase (Go)

[![CI](https://github.com/jarollz/toon-showcase/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/jarollz/toon-showcase/actions/workflows/ci.yml)

This program compares:

- JSON compact (`encoding/json` + `json.Marshal`)
- JSON pretty (`encoding/json` + `json.MarshalIndent`)
- TOON (`github.com/toon-format/toon-go` + `toon.Marshal`)
- YAML (`go.yaml.in/yaml/v3`)
- TOML (`github.com/BurntSushi/toml`)
- XML compact (`encoding/xml` + `xml.Marshal`)
- XML pretty (`encoding/xml` + `xml.MarshalIndent`)
- MessagePack (`github.com/vmihailenco/msgpack/v5`)

Comparison aspects:

- Marshal / unmarshal speed
- Encoded output length (bytes + character count)
- Character set support (English / Japanese / Chinese / emoji / escapes / control chars)
- Shape compatibility matrix for complex structures (e.g. mixed arrays, non-string map keys, nil elements)

## Run

```bash
go run .
```

Or use script with good defaults (full preset):

```bash
./run_showcase.sh
```

## Program parameters

All parameters are standard Go flags passed to `go run .`.

### Benchmark and dataset controls

- `-records` (default: `120`)
  - Number of generated orders in the synthetic dataset.
  - Bigger value increases realism and output size, but also increases runtime.

- `-iters` (default: `500`)
  - Number of measured benchmark iterations for each operation.
  - Applied to both marshal and unmarshal loops for each format.

- `-warmup` (default: `50`)
  - Warmup iterations executed before timing starts.
  - Helps reduce one-time startup effects in benchmark numbers.

- `-seed` (default: `-1`)
  - Seed used for deterministic dataset generation.
  - Use a fixed positive seed for reproducible runs.
  - Special behavior: `-seed -1` means "use current unix timestamp".

### Formatting controls

- `-indent` (default: `2`)
  - Explicit indent width in spaces.
  - Used by JSON pretty, TOON, YAML, TOML, and XML pretty for fair comparison.

- `-sample-limit` (default: `180`)
  - Maximum number of characters shown in the "encoded sample" preview for each format.
  - Display-only setting; does not affect benchmark logic.

- `-no-progress` (default: `false`)
  - Disables animated progress bar output.
  - Useful for CI logs or non-interactive terminals.

- `-baseline` (default: `toon`)
  - Baseline format used by ratio and human-friendly delta sections.
  - Supported canonical values: `toon`, `json-compact`, `json-pretty`, `yaml`, `toml`, `xml-compact`, `xml-pretty`, `messagepack`.
  - Alias examples accepted: `json`, `yml`, `xml`, `msgpack`.

### CSV output controls

- `-csv` (default: `off`)
  - CSV mode selector: `off`, `stdout`, or `file`.
  - `off`: no CSV export.
  - `stdout`: print benchmark CSV, charset CSV, and shape-compatibility CSV to terminal.
  - `file`: write benchmark CSV, charset CSV, and shape-compatibility CSV files.

- `-csv-file` (default: `benchmark_report.csv`)
  - Used when `-csv file` is selected.
  - Primary benchmark CSV writes to this path.
  - Charset CSV writes to derived path with `.charset.csv` suffix.
  - Shape compatibility CSV writes to derived path with `.shape.csv` suffix.

### Example commands

Recommended presets:

- Quick smoke check (fast local feedback):

```bash
go run . -records 20 -iters 50 -warmup 5 -seed -1
```

- Dev comparison (balanced speed/stability):

```bash
go run . -records 120 -iters 500 -warmup 50 -seed 20260604
```

- Full run (more stable benchmark numbers):

```bash
go run . -records 300 -iters 1200 -warmup 120 -seed 20260604 -csv file -csv-file benchmark_full.csv
```

Reproducible run:

```bash
go run . -records 200 -iters 800 -warmup 80 -indent 2 -seed 20260604
```

Timestamp-based seed:

```bash
go run . -seed -1
```

Compare everything against JSON compact instead of TOON:

```bash
go run . -baseline json-compact
```

CSV to stdout:

```bash
go run . -csv stdout
```

Disable progress animation:

```bash
go run . -no-progress
```

## CSV output mode

Write CSV files:

```bash
go run . -csv file -csv-file benchmark_report.csv
```

This writes:

- `benchmark_report.csv` (speed + length comparison)
- `benchmark_report.charset.csv` (character set support matrix)
- `benchmark_report.shape.csv` (shape compatibility matrix)

Print CSV to stdout instead:

```bash
go run . -csv stdout
```

Script supports env override and extra args:

```bash
RECORDS=200 ITERS=800 ./run_showcase.sh
./run_showcase.sh -sample-limit 300
./run_showcase.sh --preset quick
./run_showcase.sh --preset dev
./run_showcase.sh --preset full
```

- Script defaults:
  - `preset=full`
  - `RECORDS=300`
  - `ITERS=1200`
  - `WARMUP=120`
  - `INDENT=2`
  - `SEED=-1` (current unix timestamp)
  - `CSV_MODE=off`
  - `OUTPUT_DIR=./output`

- Script presets:
  - `quick`: `RECORDS=20`, `ITERS=50`, `WARMUP=5`
  - `dev`: `RECORDS=120`, `ITERS=500`, `WARMUP=50`
  - `full`: `RECORDS=300`, `ITERS=1200`, `WARMUP=120` (default)

- Environment variables override preset values.

Script env-to-flag mapping:

- `RECORDS` -> `-records`
- `ITERS` -> `-iters`
- `WARMUP` -> `-warmup`
- `INDENT` -> `-indent`
- `SEED` -> `-seed`
- `CSV_MODE` -> `-csv`
- `CSV_FILE` -> `-csv-file`

## Notes on comparability

- Baseline default is `TOON` (changeable via `-baseline`).
- MessagePack is binary format, so character-count metric is shown as `n/a`.
- XML has two lanes: compact and pretty.
- Benchmark dataset is intentionally TOML-safe for speed/size comparison.
- TOML/TOON/XML limitations on specific shapes are surfaced in the shape compatibility matrix section.
