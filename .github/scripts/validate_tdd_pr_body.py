#!/usr/bin/env python3

import os
import re
import sys


def fail(message: str) -> None:
    print(f"TDD evidence check failed: {message}")
    sys.exit(1)


pr_body = os.getenv("PR_BODY", "")

if not pr_body.strip():
    fail("PR description is empty")

required_headings = [
    "## TDD Evidence",
    "### Red",
    "### Green",
    "### Refactor",
    "### Final Checks",
]

missing_headings = [h for h in required_headings if h not in pr_body]
if missing_headings:
    fail(f"missing required heading(s): {', '.join(missing_headings)}")

def get_section_body(title: str) -> str:
    pattern = re.compile(rf"(?ms)^### {re.escape(title)}\n(.*?)(?=^### |\Z)")
    match = pattern.search(pr_body)
    if not match:
        fail(f"missing section: {title}")
    assert match is not None
    return match.group(1)


def validate_field(section_name: str, section_body: str, field: str) -> None:
    pattern = re.compile(rf"(?m)^- {re.escape(field)}:\s*(.+)$")
    match = pattern.search(section_body)
    if not match:
        fail(f"missing field in {section_name}: {field}")
    assert match is not None

    value = match.group(1).strip()
    if not value:
        fail(f"field is empty in {section_name}: {field}")

    lowered = value.lower()
    if lowered in {"tbd", "todo", "n/a", "na", "none"}:
        fail(f"field has placeholder value in {section_name}: {field}")


section_fields = {
    "Red": ["Test", "Command", "Expected fail", "Observed fail"],
    "Green": ["Code change", "Command", "Observed pass"],
    "Refactor": ["Refactor done", "Behavior safety"],
    "Final Checks": [
        "go test ./...",
        "make ci",
        "go run . -records 20 -iters 50 -warmup 5 -no-progress",
    ],
}

for section_name, fields in section_fields.items():
    section_body = get_section_body(section_name)
    for field in fields:
        validate_field(section_name, section_body, field)

print("TDD evidence check passed")
