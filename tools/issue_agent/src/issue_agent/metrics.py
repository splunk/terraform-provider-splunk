"""Counters for one issue run. Dry-run does not increment implemented."""

from __future__ import annotations

from dataclasses import dataclass, field


@dataclass
class Metrics:
    implemented: int = 0
    would_implement: int = 0
    user_error: int = 0
    breaking_or_harmful: int = 0
    insufficient_evidence: int = 0
    declined_other: int = 0
    review_defects: int = 0
    errors: int = 0
    skipped: int = 0
    by_issue: dict[str, str] = field(default_factory=dict)

    def record(self, issue_number: int, kind: str) -> None:
        self.by_issue[str(issue_number)] = kind
        if kind == "implemented":
            self.implemented += 1
        elif kind == "would_implement":
            self.would_implement += 1
        elif kind == "user_error":
            self.user_error += 1
        elif kind == "breaking_or_harmful":
            self.breaking_or_harmful += 1
        elif kind == "insufficient_evidence":
            self.insufficient_evidence += 1
        elif kind == "declined_other":
            self.declined_other += 1
        elif kind == "review_defects":
            self.review_defects += 1
        elif kind == "errors":
            self.errors += 1
        elif kind == "skipped":
            self.skipped += 1

    def as_dict(self) -> dict[str, object]:
        return {
            "implemented": self.implemented,
            "would_implement": self.would_implement,
            "user_error": self.user_error,
            "breaking_or_harmful": self.breaking_or_harmful,
            "insufficient_evidence": self.insufficient_evidence,
            "declined_other": self.declined_other,
            "review_defects": self.review_defects,
            "errors": self.errors,
            "skipped": self.skipped,
            "by_issue": self.by_issue,
        }
