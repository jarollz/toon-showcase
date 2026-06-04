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

## Benchmark snapshot (MacBook Pro M1 Pro)

Sample run on personal MacBook Pro M1 Pro using full preset (`records=300`, `iters=1200`, `warmup=120`, `seed=1780566620`). Source: `output/LASTRUN.txt`.

Environment note: personal laptop (`Apple Silicon M1 Pro`), single local run snapshot for quick comparison (not a controlled lab benchmark).

Reproduce this snapshot: `./run_showcase.sh --preset full -no-progress`

| Format | Marshal (ns/op) | Unmarshal (ns/op) | Bytes | Roundtrip |
| --- | ---: | ---: | ---: | :---: |
| JSON compact | 939,708 | 5,019,143 | 503,547 | ok |
| JSON pretty | 3,833,217 | 6,981,475 | 874,689 | ok |
| TOON | 5,802,221 | 7,128,738 | 487,180 | ok |
| YAML | 22,252,406 | 28,063,857 | 638,241 | ok |
| TOML | 17,243,525 | 36,782,208 | 754,286 | ok |
| XML compact | 4,295,982 | 21,441,007 | 724,441 | ok |
| XML pretty | 5,394,372 | 25,848,742 | 1,074,418 | ok |
| MessagePack | 992,791 | 2,081,718 | 419,879 | ok |

Percent delta vs TOON baseline (`-` means faster/smaller, `+` means slower/larger):

Emoji legend: `✅` better, `❌` worse, `➖` same.

| Format | Marshal delta | Unmarshal delta | Bytes delta |
| --- | ---: | ---: | ---: |
| JSON compact | -83.8% ✅ | -29.6% ✅ | +3.4% ❌ |
| JSON pretty | -33.9% ✅ | -2.1% ✅ | +79.5% ❌ |
| YAML | +283.5% ❌ | +293.7% ❌ | +31.0% ❌ |
| TOML | +197.2% ❌ | +416.0% ❌ | +54.8% ❌ |
| XML compact | -26.0% ✅ | +200.8% ❌ | +48.7% ❌ |
| XML pretty | -7.0% ✅ | +262.6% ❌ | +120.5% ❌ |
| MessagePack (binary) | -82.9% ✅ | -70.8% ✅ | -13.8% ✅ |

Quick read (all compared to TOON):

- Speed vs TOON: `JSON compact` (`-83.8%` marshal ✅, `-29.6%` unmarshal ✅), `MessagePack` (`-82.9%` marshal ✅, `-70.8%` unmarshal ✅; binary).
- Size vs TOON: `MessagePack` (`-13.8%` ✅), `JSON compact` (`+3.4%` ❌), `JSON pretty` (`+79.5%` ❌), `XML pretty` (`+120.5%` ❌).
- Unmarshal penalty vs TOON: `XML compact` (`+200.8%` ❌), `XML pretty` (`+262.6%` ❌), `YAML` (`+293.7%` ❌), `TOML` (`+416.0%` ❌); XML marshal still faster (`-26.0%` ✅ / `-7.0%` ✅).

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

### Encoded example output controls

- `-output-example` (default: `false`)
  - Enables writing full encoded shape-case examples to files.
  - `false`: skip file export.
  - `true`: write one file per `(format, shape case)` from the shape matrix.

- `-output-example-dir` (default: empty)
  - Base directory for encoded example files.
  - Used only when `-output-example=true`.
  - If empty, the program auto-generates: `output/examples_[YYYYMMDD]_[HHmmss]_[+0800]`.
  - Program output prints resolved directory path for quick lookup.

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

Output encoded examples to auto-generated directory:

```bash
go run . -output-example=true
```

Output encoded examples to custom directory:

```bash
go run . -output-example=true -output-example-dir ./output/examples_manual
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
  - `OUTPUT_EXAMPLE=false`
  - `OUTPUT_EXAMPLE_DIR=` (empty means app auto-generates timestamped output path)
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
- `OUTPUT_EXAMPLE` -> `-output-example`
- `OUTPUT_EXAMPLE_DIR` -> `-output-example-dir`

## Notes on comparability

- Baseline default is `TOON` (changeable via `-baseline`).
- MessagePack is binary format, so character-count metric is shown as `n/a`.
- XML has two lanes: compact and pretty.
- Benchmark dataset is intentionally TOML-safe for speed/size comparison.
- TOML/TOON/XML limitations on specific shapes are surfaced in the shape compatibility matrix section.

## For contributors

## Quality gate and tests

- Coverage floor is `90.0%` total (enforced by `Makefile` `COVER_MIN`).
- CI runs `make ci` (`vet + test-race + cover-check`).
- Any logic change should include matching unit test updates in affected layer.

Recommended local verification after changes:

```bash
make ci
go run . -records 20 -iters 50 -warmup 5 -no-progress
```

## Architecture

This repo follows layered architecture. New code should respect dependency direction and layer responsibility.

- `main.go`: composition root only (wire dependencies and execute flow)
- `internal/infrastructure/cli`: flag parsing and CLI option mapping
- `internal/adapter/codec`: format codec implementations
- `internal/adapter/presenter`: report/csv/progress output rendering
- `internal/core/usecase`: business flow, validation, orchestration contracts
- `internal/core/entity`: domain models and report DTOs

Dependency direction must stay inward:

- `infrastructure/adapter -> core/usecase -> core/entity`

`core/entity` and `core/usecase` must not import adapter or infrastructure packages.

## AI session history skill (client-agnostic)

This repo ships one canonical Agent Skill for session documentation:

- `.agents/skills/ai-gen-history/SKILL.md`

Supporting templates:

- `.agents/skills/ai-gen-history/template.md`
- `.agents/skills/ai-gen-history/template.json`
- Invocation examples: `.agents/skills/ai-gen-history/README.md`

The skill generates paired artifacts in `ai-gen-history/`:

- `session-<timestamp>-<topic>.md`
- `session-<timestamp>-<topic>.json`

Invocation differs by client, but uses the same skill definition:

- OpenCode: load skill `ai-gen-history` (or call through skill tool)
- Claude: run skill `ai-gen-history`
- Codex: `$ai-gen-history <optional-topic>` (or select from `/skills`)

No client-specific command files are required in this repository.
