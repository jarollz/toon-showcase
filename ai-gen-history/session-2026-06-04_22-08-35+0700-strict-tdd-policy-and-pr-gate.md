# Session History: Strict TDD Policy and PR Gate
Date: 2026-06-04 22:08:35 +0700
Project: `jarollz/toon-showcase`
Branch: `main`

## Initial User Goals
- Expand strict TDD governance beyond recommendations.
- Enforce TDD evidence through strict PR text pattern, not committed evidence files.
- Implement changes directly in repo and then capture session via `ai-gen-history`.

## Starting State
- `AGENTS.md` had architecture and coverage rules but no strict Red/Green/Refactor enforcement contract.
- CI workflow existed in `.github/workflows/ci.yml` with only Go quality gate.
- No PR template and no PR-body validation for TDD evidence.

## Collaboration Timeline
- User asked what to add for TDD, then asked for stricter mechanisms and CI implications.
- AI proposed policy + template + PR-body validation approach to avoid repo artifact bloat.
- User selected strict text pattern validation and requested execution.
- AI updated policy docs, added PR template, added validation script, and wired PR-only CI gate.
- User requested `ai-gen-history` invocation, then commit/push.

## Architecture / Design Decisions
- Use PR body validation instead of committed evidence files to keep repository clean.
- Keep validation strict but simple: required sections and non-empty field checks with placeholder rejection.
- Gate PR test job behind evidence check while preserving push-to-main CI behavior.

## Implementation Work
### Files added/changed
- `AGENTS.md`
- `.github/pull_request_template.md`
- `.github/scripts/validate_tdd_pr_body.py`
- `.github/workflows/ci.yml`
- `ai-gen-history/session-2026-06-04_22-08-35+0700-strict-tdd-policy-and-pr-gate.md`
- `ai-gen-history/session-2026-06-04_22-08-35+0700-strict-tdd-policy-and-pr-gate.json`

### Key behavior changes
- Added mandatory TDD contract and completion evidence requirements to `AGENTS.md`.
- Added standard PR template requiring Red/Green/Refactor and final-check fields.
- Added CI-enforced PR body validator to fail PRs missing required TDD evidence structure.
- CI now runs `tdd-evidence` on pull requests before Go CI test job.

## Verification Performed
- `python3 .github/scripts/validate_tdd_pr_body.py` -> failed as expected with empty `PR_BODY`.
- `PR_BODY='<sample valid body>' python3 .github/scripts/validate_tdd_pr_body.py` -> passed.
- `git status --short` -> confirmed expected modified and added files.

## Artifacts Generated
- `ai-gen-history/session-2026-06-04_22-08-35+0700-strict-tdd-policy-and-pr-gate.md`
- `ai-gen-history/session-2026-06-04_22-08-35+0700-strict-tdd-policy-and-pr-gate.json`

## Risks / Tradeoffs
- Regex-based strict text validation may reject valid but differently formatted PR descriptions.
- PR authors must fill template carefully; initial friction may increase.
- Validator checks structure and declared evidence, not actual command execution truthfulness.

## Reusable Knowledge
- Strict TDD adherence improves when policy, template, and CI gate align on same required fields.
- PR-body enforcement avoids per-feature evidence file sprawl while still creating merge gate.
- Section-scoped validation prevents duplicate field-name ambiguity across Red and Green sections.
- Keep push CI path unblocked by PR-only evidence job conditions.

## Suggested Next Steps
1. Consider extending validator with optional semantic checks (for example fail/pass token heuristics).
2. Share one sample filled PR in team docs to reduce onboarding friction.
