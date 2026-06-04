# Session History: TOON Showcase Clean-Architecture + Coverage + CI
Date: 2026-06-04 17:46:56 +0700
Project: `jarollz/toon-showcase`
Branch: `main`

## Initial User Goals
- Improve program structure.
- Make code unit testable.
- Make architecture layered.
- Raise unit test coverage to >= 90% using `coverpkg`.
- Add CI workflow.
- Add README CI badge.
- Commit and push safely.
- Capture full session history.

## Starting State
- Almost all logic in `main.go` (monolithic, mixed concerns).
- No significant test coverage outside minimal seed tests.
- No `Makefile`.
- No GitHub Actions workflow.

## Architecture Decision
Adopt clean-architecture style, pragmatic version:
- **Core Entity**: data models and result structures.
- **Core Usecase**: orchestration + benchmark/charset/shape/dataset logic.
- **Adapters**:
  - codec adapters (JSON/TOON/YAML/TOML/XML/MessagePack),
  - presenter adapters (text report, csv, progress).
- **Infrastructure**: CLI flag parser.
- **Entrypoint**: thin `main.go` composition root.

## Major Refactor Outcomes
### Layered packages created
- `internal/core/entity`
- `internal/core/usecase`
- `internal/adapter/codec`
- `internal/adapter/presenter`
- `internal/infrastructure/cli`

### Entrypoint reshaped
- `main.go` reduced to orchestration.
- Added testable `run(args, stdout, stderr, now) int`.
- `main()` now only `os.Exit(run(...))`.

### Behavior preserved
- CLI behavior preserved (records/iters/warmup/seed/indent/sample-limit/no-progress/baseline/csv/csv-file).
- Benchmark reports, charset matrix, shape matrix, CSV modes preserved.
- Script compatibility preserved.

## Testing and Coverage Work
### Tests added across layers
- `main_test.go`
- `internal/core/usecase/*_test.go`:
  - run orchestration success/error branches
  - benchmark error/success/mismatch branches
  - charset matrix pass/fail/decode mismatch
  - shape matrix pass/fail/decode mismatch
  - options validation branches
  - dataset determinism + invariants
- `internal/adapter/codec/codecs_test.go`
  - codec contract roundtrip coverage
  - XML helper branch coverage
- `internal/adapter/presenter/*_test.go`
  - text renderer
  - csv builders + emit stdout/file/error
  - progress lifecycle
- `internal/infrastructure/cli/flags_test.go`
  - default values
  - custom parse path
  - parse error path

### Coverage method and result
- Required method: `go test -coverpkg=./... ./...`
- Final total coverage (from `go tool cover -func`): **94.6%**
- Meets/exceeds target >= 90%.

## Tooling/Automation Added
### Makefile
Added targets:
- `test`
- `test-race`
- `vet`
- `cover`
- `cover-check` (hard gate, compares total against `COVER_MIN`)
- `clean`
- `ci`

Coverage config:
- `COVERPKG := ./...`
- `COVERMODE := atomic`
- `COVERPROFILE := .coverage.out`
- `COVER_MIN := 90.0`

### GitHub Actions
Added `.github/workflows/ci.yml`:
- Triggers: `pull_request`, `push` to `main`
- Concurrency cancel-in-progress
- Setup Go via `go.mod`
- Runs `make ci`
- Always attempts `make cover || true`
- Uploads `.coverage.out` artifact

### README
Added CI badge:
- Workflow badge for `ci.yml` on `main`.

## Verification Performed
- `go test ./...` passed.
- quick runtime smoke passed:
  - `go run . -records 20 -iters 50 -warmup 5 -no-progress`
- `make cover-check` passed.
- coverage total verified 94.6%.

## Git Safety + Delivery
Pre-commit checks performed:
- `git status --short`
- `git log --oneline -10`
- `git diff` / `git diff --cached --stat`

Commit:
- SHA: `867f2a6`
- Message: `refactor: layer app and add 90% coverage CI`

Push:
- Remote: `origin`
- Branch: `main`
- Result: success (`36a004c -> 867f2a6`)
- Final working tree: clean.

## Key Decisions and Rationale
- Keep clean architecture pragmatic; avoid over-abstraction.
- Keep current functionality stable while restructuring.
- Add test seam in entrypoint (`run`) for reliable branch testing.
- Use strict `coverpkg` aggregate number, not package-local illusion.
- Enforce coverage via Makefile gate, then wire CI to gate.
- Upload coverage artifact for CI debugging.

## Knowledge Captured (Reusable)
- For this repo, `coverpkg` aggregate reached 94.6% with mixed unit + adapter contract tests.
- Strong ROI from testing error branches in orchestrator and CSV/file writers.
- CI gate chain `make ci` gives reproducible local/remote policy.
- `run(...)` extraction in main is minimal but high-impact for coverage/testability.

## Files Added/Changed (Session)
- `main.go`
- `main_test.go`
- `Makefile`
- `.github/workflows/ci.yml`
- `README.md`
- `internal/core/entity/model.go`
- `internal/core/entity/result.go`
- `internal/core/usecase/types.go`
- `internal/core/usecase/run.go`
- `internal/core/usecase/benchmark.go`
- `internal/core/usecase/dataset.go`
- `internal/core/usecase/charset.go`
- `internal/core/usecase/shape.go`
- `internal/core/usecase/types_test.go`
- `internal/core/usecase/run_test.go`
- `internal/core/usecase/benchmark_test.go`
- `internal/core/usecase/charset_shape_test.go`
- `internal/adapter/codec/codecs.go`
- `internal/adapter/codec/codecs_test.go`
- `internal/adapter/presenter/text.go`
- `internal/adapter/presenter/csv.go`
- `internal/adapter/presenter/progress.go`
- `internal/adapter/presenter/text_test.go`
- `internal/adapter/presenter/csv_test.go`
- `internal/adapter/presenter/progress_test.go`
- `internal/infrastructure/cli/flags.go`
- `internal/infrastructure/cli/flags_test.go`

## Suggested Next Steps
1. Add branch protection requiring `Go CI` check.
2. Optionally add CI badge status explanation section in README.
3. Optionally add per-package coverage thresholds (if stricter policy wanted).
