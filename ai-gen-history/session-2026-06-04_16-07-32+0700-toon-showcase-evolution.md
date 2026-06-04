# Session History: TOON Showcase Evolution
Date: 2026-06-04 16:07:32 +0700
Project: `try-toon-format/toon-showcase`

## Initial Goal
- Build Go showcase for `toon-go` marshalling/unmarshalling.
- Compare against JSON on:
  - speed
  - encoded length

## Scope Expansion Timeline
1. Added 3-way benchmark: JSON compact, JSON pretty, TOON.
2. Made indent explicit and equalized for JSON pretty + TOON.
3. Added UTF-8 dataset content (English/Japanese/Chinese/emoji).
4. Added charset support comparison matrix.
5. Added CSV export mode.
6. Added helper script `run_showcase.sh`.
7. Aligned seed behavior: `-1` => current unix timestamp.
8. Added human-friendly delta report (`faster/slower`, `smaller/larger`).
9. Expanded formats: YAML, TOML, XML, MessagePack.
10. Added shape compatibility matrix for complex object structures.
11. Changed baseline logic:
    - default baseline = TOON
    - configurable via `-baseline`
    - alias support (`json`, `yml`, `msgpack`, etc.)
12. Changed script default `CSV_MODE=off`.
13. Added animated progress bar with percentage/stage.
14. Added `-no-progress` flag.
15. Added XML split lanes:
    - `XML compact`
    - `XML pretty`

## Final Program Capabilities
- Benchmark formats:
  - `JSON compact`
  - `JSON pretty`
  - `TOON`
  - `YAML`
  - `TOML`
  - `XML compact`
  - `XML pretty`
  - `MessagePack`
- Reports:
  - Benchmark table (marshal/unmarshal + bytes/chars + roundtrip)
  - Ratio table vs selected baseline
  - Human-friendly delta table
  - Charset support matrix
  - Shape compatibility matrix
- CSV outputs:
  - benchmark CSV
  - charset CSV
  - shape CSV

## Key Design Decisions
- Keep benchmark dataset TOML-safe so all major lanes produce performance numbers.
- Surface unsupported structural cases in dedicated shape matrix instead of breaking benchmark.
- Keep binary format (`MessagePack`) char metric as `n/a`.
- Use baseline as first-class flag; default TOON to match showcase theme.
- Keep script convenience defaults for long run preset, but CSV off unless requested.

## Current CLI Flags (Program)
- `-records`
- `-iters`
- `-warmup`
- `-seed` (`-1` => unix timestamp)
- `-indent`
- `-sample-limit`
- `-baseline`
- `-csv`
- `-csv-file`
- `-no-progress`

## Baseline Behavior
- Default baseline: `toon`.
- Supported canonical values:
  - `toon`
  - `json-compact`
  - `json-pretty`
  - `yaml`
  - `toml`
  - `xml-compact`
  - `xml-pretty`
  - `messagepack`
- Alias examples:
  - `json` -> `json-compact`
  - `yml` -> `yaml`
  - `xml` -> `xml-compact`
  - `msgpack` -> `messagepack`

## Known Compatibility Findings
- TOON:
  - strong structured support
  - unsupported control chars in tested cases (`U+0000`, `U+0001`)
  - non-string map key unsupported
- TOML:
  - supports many nested struct/slice shapes
  - rejects nil element in arrays
  - rejects non-string map keys
- XML:
  - works on projected benchmark model
  - weak for map-heavy/mixed dynamic probe shapes
- MessagePack:
  - broad shape support in probes
  - binary payload (char metric not meaningful)

## Script Behavior (`run_showcase.sh`)
- Presets:
  - `quick`: records=20, iters=50, warmup=5
  - `dev`: records=120, iters=500, warmup=50
  - `full`: records=300, iters=1200, warmup=120 (default)
- Defaults:
  - `SEED=-1`
  - `CSV_MODE=off`
- Environment overrides supported.

## Files Added/Changed During Session
- `main.go`
- `go.mod`
- `go.sum`
- `README.md`
- `run_showcase.sh`

## Validation Summary
- Multiple `go run` checks across default/explicit baselines.
- Script checks for preset and defaults.
- Verified progress bar and `-no-progress` behavior.
- Verified CSV generation modes and output naming.

## Suggested Next Steps
- Add optional markdown report export (`-report-md`).
- Add allocation stats (`testing.Benchmark` style or `runtime` snapshots).
- Add optional per-format timeout guard for very large runs.
