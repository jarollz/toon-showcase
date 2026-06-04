# Agent Notes

## Repo shape
- Single Go module/binary (`go.mod`, `main.go`); no subpackages.
- Main behavior lives in `main.go`; CLI flags are source of truth for runtime options.

## Fast commands
- Compile check: `go test ./...` (currently returns `[no test files]`; still best syntax/build sanity check).
- Run benchmark/report directly: `go run . [flags]`.
- Preferred scripted run: `./run_showcase.sh [--preset quick|dev|full] [extra go run flags]`.

## Runtime gotchas
- Script default preset is `full` (`records=300`, `iters=1200`, `warmup=120`), which is slow; use `--preset quick` for iteration.
- Use `-no-progress` for non-interactive runs to avoid spinner noise in captured logs.
- CSV export only writes files when `-csv file` (or `CSV_MODE=file` in script). Default is `off`.

## Verification flow for code changes
- 1) `go test ./...`
- 2) `go run . -records 20 -iters 50 -warmup 5 -no-progress`
- Use larger presets only when validating benchmark stability.

## Ignore unless task asks
- Generated output folder is gitignored: `output/` (`.gitignore`).
