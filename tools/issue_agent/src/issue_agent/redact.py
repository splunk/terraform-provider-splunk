"""Remove credentials and certificate material before prompts, logs, or comments."""

from __future__ import annotations

import re

_PATTERNS = (
    re.compile(r"-----BEGIN [A-Z0-9 ]*-----.*?-----END [A-Z0-9 ]*-----", re.DOTALL),
    re.compile(r"\bAKIA[0-9A-Z]{16}\b"),
    re.compile(r"\bghp_[A-Za-z0-9]{20,}\b"),
    re.compile(r"\beyJ[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\.[A-Za-z0-9_-]{8,}\b"),
    re.compile(
        r"(?i)\b(password|passwd|secret|token|api[_-]?key|authorization)\b(\s*[:=]\s*)(\S+)"
    ),
)


def redact(text: str, extra: tuple[str, ...] = ()) -> str:
    cleaned = text or ""
    for secret in extra:
        if secret:
            cleaned = cleaned.replace(secret, "[REDACTED]")
    cleaned = _PATTERNS[0].sub("[REDACTED-CERTIFICATE]", cleaned)
    cleaned = _PATTERNS[1].sub("[REDACTED]", cleaned)
    cleaned = _PATTERNS[2].sub("[REDACTED]", cleaned)
    cleaned = _PATTERNS[3].sub("[REDACTED]", cleaned)

    def _assignment(match: re.Match[str]) -> str:
        return f"{match.group(1)}{match.group(2)}[REDACTED]"

    return _PATTERNS[4].sub(_assignment, cleaned)
