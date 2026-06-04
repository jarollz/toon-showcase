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

Sample run on personal MacBook Pro M1 Pro using full preset (`records=300`, `iters=1200`, `warmup=120`, `seed=1780601429`). Source: `benchmark-snapshot/snapshot_20260605_023027.md`.

Environment note: personal laptop (`Apple Silicon M1 Pro`), single local run snapshot for quick comparison (not a controlled lab benchmark).

Reproduce this snapshot: `./run_showcase.sh --preset full -no-progress`

For complete details (including charset and shape sections), see [`benchmark-snapshot/snapshot_20260605_023027.md`](benchmark-snapshot/snapshot_20260605_023027.md).

| Format | Marshal (ns/op) | Unmarshal (ns/op) | Bytes | Roundtrip |
| --- | ---: | ---: | ---: | :---: |
| JSON compact | 941,834 | 5,128,791 | 503,572 | ok |
| JSON pretty | 3,704,891 | 6,926,182 | 874,714 | ok |
| TOON | 5,848,249 | 7,025,284 | 487,205 | ok |
| YAML | 25,221,971 | 27,765,547 | 638,266 | ok |
| TOML | 17,435,765 | 35,610,408 | 754,281 | ok |
| XML compact | 4,159,604 | 20,866,120 | 724,466 | ok |
| XML pretty | 5,425,438 | 25,723,025 | 1,074,443 | ok |
| MessagePack | 1,030,919 | 2,062,222 | 419,879 | ok |

Percent delta vs TOON baseline (`-` means faster/smaller, `+` means slower/larger):

Emoji legend: `✅` better, `❌` worse, `➖` same.

| Format | Marshal delta | Unmarshal delta | Bytes delta |
| --- | ---: | ---: | ---: |
| JSON compact | -83.9% ✅ | -27.0% ✅ | +3.4% ❌ |
| JSON pretty | -36.6% ✅ | -1.4% ✅ | +79.5% ❌ |
| YAML | +331.3% ❌ | +295.2% ❌ | +31.0% ❌ |
| TOML | +198.1% ❌ | +406.9% ❌ | +54.8% ❌ |
| XML compact | -28.9% ✅ | +197.0% ❌ | +48.7% ❌ |
| XML pretty | -7.2% ✅ | +266.1% ❌ | +120.5% ❌ |
| MessagePack (binary) | -82.4% ✅ | -70.6% ✅ | -13.8% ✅ |

Quick read (all compared to TOON):

- Speed vs TOON: `JSON compact` (`-83.9%` marshal ✅, `-27.0%` unmarshal ✅), `MessagePack` (`-82.4%` marshal ✅, `-70.6%` unmarshal ✅; binary).
- Size vs TOON: `MessagePack` (`-13.8%` ✅), `JSON compact` (`+3.4%` ❌), `JSON pretty` (`+79.5%` ❌), `XML pretty` (`+120.5%` ❌).
- Unmarshal penalty vs TOON: `XML compact` (`+197.0%` ❌), `XML pretty` (`+266.1%` ❌), `YAML` (`+295.2%` ❌), `TOML` (`+406.9%` ❌); XML marshal still faster (`-28.9%` ✅ / `-7.2%` ✅).

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

- `-encoded-examples-dir` (default: empty)
  - Directory for encoded shape-case example files.
  - Empty value disables example export.
  - Non-empty value must be a directory path (created automatically if missing).
  - If the path exists and is a file, the program fails with a clear error.
  - Files are written under `<dir>/shape-encoded/`.

### CSV output controls

- `-csv-dir` (default: empty)
  - Empty value disables CSV export.
  - `stdout` prints benchmark/charset/shape CSV to terminal.
  - Any other non-empty value must be a directory path (created automatically if missing).
  - If the path exists and is a file, the program fails with a clear error.
  - Directory mode writes fixed files:
    - `benchmark_report.csv`
    - `benchmark_report.charset.csv`
    - `benchmark_report.shape.csv`
  - `stdout` is a reserved sentinel. If you want a literal directory named `stdout`, use `./stdout`.

### Markdown output controls

- `-md-file` (default: empty)
  - Empty value disables markdown report output.
  - `stdout` renders markdown report to terminal with Glamour `dark` theme.
  - Any other non-empty value writes raw markdown report to that file.
  - `stdout` is a reserved sentinel. If you want a literal file named `stdout`, use `./stdout`.

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
go run . -records 300 -iters 1200 -warmup 120 -seed 20260604 -csv-dir ./output
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
go run . -csv-dir stdout
```

Markdown to file:

```bash
go run . -md-file ./output/report.md
```

Markdown rendered in terminal (Glamour dark theme):

```bash
go run . -md-file stdout
```

Disable progress animation:

```bash
go run . -no-progress
```

Output encoded examples to custom directory:

```bash
go run . -encoded-examples-dir ./output/examples_manual
```

## Checked-in encoded examples (`encoded-examples/`)

This repo includes committed sample output at `encoded-examples/` for quick inspection of full shape-case payloads without running the tool.

- One directory per format: `json-compact`, `json-pretty`, `toon`, `yaml`, `toml`, `xml-compact`, `xml-pretty`, `messagepack`.
- Each format contains the same 7 shape-case filenames from the shape compatibility matrix.
- File suffix meaning:
  - `.txt`: encode succeeded (text formats)
  - `.bin`: encode succeeded (binary format, currently MessagePack)
  - `.error.txt`: encode failed for that `(format, shape case)`; file content is encoder error message.

Expected failure artifacts in this snapshot:

- `toml`: `map_non_string_key`, `slice_pointer_with_nil`
- `toon`: `map_non_string_key`
- `xml-compact` and `xml-pretty`: `map_non_string_key`, `map_string_to_slice_struct`, `nested_slice_map`, `slice_map_string`

To regenerate fresh examples from current code:

```bash
go run . -encoded-examples-dir ./output/examples_manual -no-progress
```

Generated files are written under `./output/examples_manual/shape-encoded/`.
`encoded-examples/` is checked-in sample copy for reference.

## CSV output mode

Write CSV files:

```bash
go run . -csv-dir ./output
```

This writes:

- `benchmark_report.csv` (speed + length comparison)
- `benchmark_report.charset.csv` (character set support matrix)
- `benchmark_report.shape.csv` (shape compatibility matrix)

Print CSV to stdout instead:

```bash
go run . -csv-dir stdout
```

Note: `stdout` is a reserved sentinel. If you want a literal output directory named `stdout`, use `./stdout`.

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
  - `CSV_DIR=` (empty means off)
  - `MD_FILE=` (empty means off)
  - `ENCODED_EXAMPLES_DIR=` (empty means off)

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
- `CSV_DIR` -> `-csv-dir`
- `MD_FILE` -> `-md-file`
- `ENCODED_EXAMPLES_DIR` -> `-encoded-examples-dir`

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
