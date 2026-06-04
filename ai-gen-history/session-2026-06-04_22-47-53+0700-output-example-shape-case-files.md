# Session History: Output Example Shape-Case Files
Date: 2026-06-04 22:47:53 +0700
Project: `jarollz/toon-showcase`
Branch: `main`

## Initial User Goals
- Add feature to export full encoded payload files per format and per shape case.
- Control feature with `-output-example` and `-output-example-dir`.
- Keep default disabled, and when enabled with empty dir use timestamped default folder.
- Print resolved output directory in program output.

## Starting State
- Program only printed `Encoded output sample (truncated)` in terminal.
- No exporter for full shape-case payload files.
- CLI had no output-example flags.

## Collaboration Timeline
- User requested non-truncated encoded outputs by format and shape case.
- User refined API to `-output-example` + `-output-example-dir` and set default behavior.
- User requested timestamp default directory `output/examples_[YYYYMMDD]_[HHmmss]_[z]` and timezone token `+0800` style.
- AI implemented feature with TDD slices, validated by targeted and full test runs.
- User manually verified behavior and requested history + commit + push.

## Architecture / Design Decisions
- Keep benchmark terminal sample truncated; add separate file export path to avoid noisy stdout.
- Capture shape artifacts in core usecase/entity layer so adapter writer stays I/O-only.
- Exporter writes `.txt` for text codecs, `.bin` for binary codecs, and `.error.txt` for encode failures.
- Resolve default directory via run clock for deterministic tests: `output/examples_<date>_<time>_<-0700>`.

## Implementation Work
### Files added/changed
- `internal/infrastructure/cli/flags.go`
- `internal/infrastructure/cli/flags_test.go`
- `internal/core/usecase/types.go`
- `internal/core/usecase/types_test.go`
- `internal/core/entity/result.go`
- `internal/core/usecase/shape.go`
- `internal/core/usecase/charset_shape_test.go`
- `internal/core/usecase/run.go`
- `internal/adapter/presenter/encoded_examples.go`
- `internal/adapter/presenter/encoded_examples_test.go`
- `main.go`
- `main_test.go`
- `run_showcase.sh`
- `README.md`
- `encoded-examples/` (manual sample output folder from user verification)
- `ai-gen-history/session-2026-06-04_22-47-53+0700-output-example-shape-case-files.md`
- `ai-gen-history/session-2026-06-04_22-47-53+0700-output-example-shape-case-files.json`

### Key behavior changes
- New flags:
  - `-output-example` default `false`
  - `-output-example-dir` default empty
- When `-output-example=true` and dir empty, output dir auto-resolves to timestamped path under `output/`.
- Program prints `Encoded examples output dir: <resolved-path>` when feature enabled.
- Full shape-case encoded artifacts are written to disk per format and case.

## Verification Performed
- `go test ./internal/infrastructure/cli ./internal/core/usecase` -> red then green for new options + resolver.
- `go test ./internal/core/usecase -run TestRunShapeMatrix` -> red for missing shape artifacts.
- `go test ./internal/core/usecase` -> green after shape artifact implementation.
- `go test ./internal/adapter/presenter ./... -run 'TestWriteShapeEncodedExamples|TestRunOutputExample|TestRunShapeMatrix'` -> red then green for writer + main integration.
- `go test ./...` -> pass.
- `make ci && go run . -records 20 -iters 50 -warmup 5 -no-progress` -> pass, coverage 94.6%.

## Artifacts Generated
- `encoded-examples/` (manual sample output folder included by user request)
- `ai-gen-history/session-2026-06-04_22-47-53+0700-output-example-shape-case-files.md`
- `ai-gen-history/session-2026-06-04_22-47-53+0700-output-example-shape-case-files.json`

## Risks / Tradeoffs
- Timestamped default path can produce many output folders across repeated runs.
- Exporting all shape cases increases filesystem output volume.
- `.error.txt` files represent encode failures and are expected for unsupported codec/shape combinations.

## Reusable Knowledge
- Keep export feature behind explicit flag to avoid default disk churn.
- Use run-injected clock for deterministic timestamp-path tests.
- Store export payloads in core shape report DTO to keep adapter implementation simple.
- Emit resolved output path in stdout to reduce user confusion about destination.

## Suggested Next Steps
1. Optionally add cleanup/retention utility for timestamped `output/examples_*` folders.
2. Optionally add README snippet showing expected file tree for one run.
