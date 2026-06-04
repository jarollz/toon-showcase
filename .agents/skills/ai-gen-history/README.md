# ai-gen-history skill usage

Canonical skill location:

- `.agents/skills/ai-gen-history/SKILL.md`

Generated outputs:

- `ai-gen-history/session-<YYYY-MM-DD_HH-MM-SS+ZZZZ>-<topic>.md`
- `ai-gen-history/session-<YYYY-MM-DD_HH-MM-SS+ZZZZ>-<topic>.json`

## Invocation examples

### OpenCode

- Via skill tool selection: choose `ai-gen-history`
- With topic argument: `ai-gen-history benchmark-readme-update`

### Claude

- Invoke skill: `/ai-gen-history`
- With topic argument: `/ai-gen-history benchmark-readme-update`

### Codex

- Mention skill: `$ai-gen-history`
- With topic argument: `$ai-gen-history benchmark-readme-update`
- Alternative: open `/skills`, select `ai-gen-history`, then provide topic text

## Notes

- Behavior source is shared skill content in `SKILL.md`.
- Templates used by skill:
  - `template.md`
  - `template.json`
- If topic omitted, skill infers concise topic slug from session context.
