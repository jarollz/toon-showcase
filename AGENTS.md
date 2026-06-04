# Agent Notes

## Repo shape
- Single Go module binary with layered packages under `internal/`.
- Entrypoint orchestration lives in `main.go`; business logic does not.
- Current layers:
  - `internal/core/entity`: domain models and result DTOs.
  - `internal/core/usecase`: application use cases, interfaces, validation, runner flow.
  - `internal/adapter/codec`: format codec adapters implementing use case contracts.
  - `internal/adapter/presenter`: text/csv/progress output adapters.
  - `internal/infrastructure/cli`: CLI flag parsing and option wiring.

## Architecture rules (must follow)
- Keep dependency direction inward: `infrastructure/adapter -> core/usecase -> core/entity`.
- `core/entity` must not import adapter/infrastructure packages.
- `core/usecase` must not import adapter/infrastructure packages.
- Put new behavior in layer by responsibility, not in `main.go`.
- `main.go` stays composition root: parse options, wire dependencies, execute runner, emit output.

## Test and coverage requirements
- Unit tests exist across layers and are part of normal workflow.
- Coverage gate is enforced by `make ci` via `make cover-check`.
- Minimum total coverage is `90.0%` (`Makefile` `COVER_MIN`).
- Any logic change should include/adjust tests in same affected layer.

## Fast commands
- Run full local quality gate: `make ci`.
- Compile + unit tests only: `go test ./...`.
- Coverage report with profile: `make cover`.
- Run benchmark/report directly: `go run . [flags]`.
- Preferred scripted run: `./run_showcase.sh [--preset quick|dev|full] [extra go run flags]`.

## Runtime gotchas
- Script default preset is `full` (`records=300`, `iters=1200`, `warmup=120`), slow for iteration.
- Use `--preset quick` (or `-records 20 -iters 50 -warmup 5`) for fast local checks.
- Use `-no-progress` for non-interactive runs to avoid spinner noise in captured logs.
- CSV files write only when `-csv file` (or `CSV_MODE=file` in script). Default is `off`.

## Verification flow for code changes
- 1) `make ci`
- 2) `go run . -records 20 -iters 50 -warmup 5 -no-progress`
- Use larger presets only when validating benchmark stability.

## AI session history skill
- Canonical skill lives at `.agents/skills/ai-gen-history/SKILL.md`.
- Skill templates live at `.agents/skills/ai-gen-history/template.md` and `.agents/skills/ai-gen-history/template.json`.
- Generated artifacts should be written to `ai-gen-history/session-<timestamp>-<topic>.md` and `.json`.
- For session recap/handoff/history generation, prefer using this shared skill instead of ad-hoc formats.
- Keep skill workflow client-agnostic in repo source; do not add repo command logic under `.opencode/` or `.claude/` unless explicitly requested.
- Never include secrets/credentials/private URLs or hidden chain-of-thought in history artifacts.

## Ignore unless task asks
- Generated output folder is gitignored: `output/` (`.gitignore`).
