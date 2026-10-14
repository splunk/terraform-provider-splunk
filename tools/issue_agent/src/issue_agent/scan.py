"""Reject workflow and orchestrator changes that call the GitHub merge API."""

from __future__ import annotations

import re
from pathlib import Path

_PATTERNS = (
    re.compile(r"gh\s+pr\s+merge"),
    re.compile(r"pulls/.*/merge"),
    re.compile(r"merge_pull_request"),
)

_SUFFIXES = {".py", ".yml", ".yaml"}


def find_merge_calls(root: Path) -> list[str]:
    bases = [
        root / ".github" / "workflows",
        root / "tools" / "issue_agent" / "src",
        root / "tools" / "issue_agent" / "scripts",
    ]
    hits: list[str] = []
    for base in bases:
        if not base.exists():
            continue
        for path in sorted(base.rglob("*")):
            if path.name == "scan.py" or path.suffix not in _SUFFIXES or not path.is_file():
                continue
            text = path.read_text(encoding="utf-8")
            for pattern in _PATTERNS:
                if pattern.search(text):
                    hits.append(f"{path.relative_to(root)} matches {pattern.pattern}")
    return hits
