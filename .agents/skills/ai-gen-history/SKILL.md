---
name: ai-gen-history
description: Generate ai-gen-history session artifacts (.md + .json) that capture goals, decisions, collaboration flow, changed files, validation, and reusable knowledge for future human/AI contributors.
---

# AI Gen History

Generate session history artifacts for current collaboration.

Topic input: `$ARGUMENTS`

Use local templates:
- `template.md`
- `template.json`

Output location and naming:
- Directory: `ai-gen-history/`
- Base name: `session-<YYYY-MM-DD_HH-MM-SS+ZZZZ>-<topic-slug>`
- Files:
  - `ai-gen-history/<base>.md`
  - `ai-gen-history/<base>.json`

Rules:
1. Infer `topic-slug` from `$ARGUMENTS`; if empty, infer from session goals and keep concise kebab-case.
2. Reuse style from existing `ai-gen-history/session-*.md` and `session-*.json` if present.
3. Capture: user goals, starting state, decisions, rationale, key commands, artifacts, changed files, validation, and next steps.
4. Include collaboration flow (human asks, AI actions, decision pivots), but summarize reasoning. Do not expose hidden chain-of-thought.
5. If data is unknown, write `unknown` (or empty array) instead of inventing facts.
6. Exclude secrets, tokens, credentials, and private URLs.

Content requirements:
- Markdown report must be human-readable and sufficient for future AI agents to continue work.
- JSON report must be machine-friendly and follow template key structure.
- Keep md/json aligned on major facts (topic, goals, decisions, files, outcomes).

After writing files:
- Print created file paths.
- Print 3-6 bullet recap of most reusable knowledge from session.
