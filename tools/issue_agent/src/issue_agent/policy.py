"""Decide whether an issue may be implemented, and whether a diff stays in bounds."""

from __future__ import annotations

import re
from typing import Any

from issue_agent.schema import SchemaError, load_schema, validate

ALLOWED_REPOSITORY = "splunk/terraform-provider-splunk"
MODES = ("dry-run", "comment", "implement")
MAX_FILES = 15
MAX_LINES = 800
MAX_REVIEW_ITERATIONS = 2
MAX_RETRIES = 2
TIMEOUT_SECONDS = 1800
MAX_ISSUE_BODY_CHARS = 12_000
MAX_COMMENT_CHARS = 2_000
MAX_COMMENTS = 10
MAX_ISSUE_TEXT_CHARS = 36_000
MAX_DIFF_CHARS = 60_000
TOKEN_BUDGET = {
    "disposition": 80_000,
    "implement": 200_000,
    "revise": 200_000,
    "review": 80_000,
}

CONTRACT_MARKER = "<!-- issue-agent-contract -->"
REVIEW_MARKER = "<!-- issue-agent-review -->"
DISPOSITION_MARKER = "<!-- issue-agent-disposition -->"

TERMINAL_LABELS = (
    "agent:user-error",
    "agent:declined",
    "agent:breaking-risk",
    "agent:insufficient-evidence",
    "agent:review-failed",
    "agent:merge-ready",
    "agent:running",
)

FORBIDDEN_EXACT = frozenset(
    {"go.mod", "go.sum", ".goreleaser.yml", ".gitattributes", ".gitmodules"}
)
FORBIDDEN_PREFIXES = (
    ".github/workflows/",
    "vendor/",
    "tools/issue_agent/",
)

REQUIRED_CHECKS = ("unit", "acceptance", "lint")

_COMMENT_BANS = (
    "please provide",
    "need more information",
    "waiting for",
    "can you",
    "@",
)


class PolicyError(ValueError):
    """Raised when an event or result violates pipeline policy."""


def validate_event(event: dict[str, Any]) -> None:
    repository = event.get("repository")
    if repository != ALLOWED_REPOSITORY:
        raise PolicyError(f"repository {repository} is outside the allowlist")
    if event.get("pull_request_target") or event.get("event_name") == "pull_request_target":
        raise PolicyError("pull_request_target cannot start the pipeline")
    if event.get("is_fork"):
        raise PolicyError("fork pull requests cannot start write-capable work")
    if event.get("issue_state") not in (None, "open"):
        raise PolicyError("only open issues are in scope")


def normalize_disposition(raw: dict[str, Any]) -> dict[str, Any]:
    try:
        validate(raw, load_schema("disposition.schema.json"))
    except SchemaError as exc:
        raise PolicyError(str(exc)) from exc

    result = dict(raw)
    reasons: list[str] = []
    validity = result["validity"]
    contract = result["implementation_contract"]
    breaking = _breaking(result["breaking"])
    compatible = result["compatibility"] == "compatible" and not breaking
    high = result["confidence"] == "high"
    worth_it = result["model_decision"] == "proceed"
    code_evidence = any(
        item["source"] in {"code", "history", "docs"} for item in result["evidence"]
    )
    decision = "decline"
    label: str | None = "agent:declined"

    if validity == "user_error":
        label = "agent:user-error"
        reasons.append("reported behavior matches documented usage")
    elif validity in {
        "unsupported_request",
        "already_supported",
        "already_fixed",
        "duplicate",
    }:
        label = "agent:declined"
        reasons.append(validity)
    elif validity == "insufficient_evidence":
        label = "agent:insufficient-evidence"
        reasons.append("evidence does not separate a defect from user error")
    elif breaking or result["compatibility"] != "compatible":
        label = "agent:breaking-risk"
        reasons.append("fix would change current successful behavior or schema")
    elif validity == "docs_mismatch":
        docs_only = bool(contract) and _docs_only(contract["allowed_paths"])
        if worth_it and compatible and high and docs_only:
            decision = "proceed"
            label = None
            reasons.append("docs-only correction")
        else:
            label = "agent:declined"
            reasons.append("docs mismatch cannot change provider behavior")
    elif validity == "confirmed_defect" and worth_it and high and code_evidence and contract:
        path_error = _paths_allowed(contract["allowed_paths"])
        size_error = _size_allowed(contract["expected_files"], contract["expected_lines"])
        if path_error or size_error:
            label = "agent:declined"
            reasons.append(path_error or size_error or "out of bounds")
        else:
            decision = "proceed"
            label = None
            reasons.append("confirmed compatible defect")
    elif validity == "feature_request" and worth_it and compatible and high and code_evidence and contract:
        path_error = _paths_allowed(contract["allowed_paths"])
        size_error = _size_allowed(contract["expected_files"], contract["expected_lines"])
        if path_error or size_error:
            label = "agent:declined"
            reasons.append(path_error or size_error or "out of bounds")
        else:
            decision = "proceed"
            label = None
            reasons.append("approved compatible feature")
    else:
        label = "agent:declined"
        reasons.append("the change is not worth an automated fix")

    if decision != "proceed":
        contract = None
    result["decision"] = decision
    result["label"] = label
    result["reasons"] = reasons
    result["implementation_contract"] = contract
    return result


def normalize_review(
    raw: dict[str, Any],
    *,
    changed_files: list[str],
    changed_lines: int,
    allowed_paths: list[str],
    iteration: int,
    is_agent_pr: bool,
    is_fork: bool,
) -> dict[str, Any]:
    try:
        validate(raw, load_schema("review.schema.json"))
    except SchemaError as exc:
        raise PolicyError(str(exc)) from exc

    blocking = list(raw["blocking"])
    if changed_lines > MAX_LINES or len(changed_files) > MAX_FILES:
        blocking.append({"path": "", "summary": "diff exceeds the file or line cap"})
    outside = [path for path in changed_files if not _path_allowed(path, allowed_paths)]
    for path in outside:
        blocking.append({"path": path, "summary": "path is outside the approved contract"})
    forbidden = [path for path in changed_files if _forbidden(path)]
    for path in forbidden:
        blocking.append({"path": path, "summary": "path is forbidden for the agent"})

    model_pass = (
        raw["within_contract"]
        and raw["compatibility"] == "compatible"
        and not raw["blocking"]
    )
    passed = model_pass and not blocking
    correction = (
        is_agent_pr
        and not is_fork
        and not passed
        and iteration <= MAX_REVIEW_ITERATIONS
    )
    return {
        "decision": "pass" if passed else "fail",
        "within_contract": passed,
        "compatibility": "compatible" if passed else "breaking",
        "blocking": blocking,
        "non_blocking": raw["non_blocking"],
        "summary": raw["summary"],
        "correction_requested": correction,
        "label": None if passed or correction or not is_agent_pr else "agent:review-failed",
    }


def readiness(
    checks: dict[str, str],
    *,
    is_agent_pr: bool,
    review_passed: bool,
) -> dict[str, Any]:
    blocking = [name for name in REQUIRED_CHECKS if checks.get(name) != "success"]
    if is_agent_pr and not review_passed:
        blocking.append("agent-review-findings")
    ready = not blocking
    label = "agent:merge-ready" if ready and is_agent_pr else None
    result = {
        "ready": ready,
        "label": label,
        "merge": False,
        "blocking_checks": blocking,
    }
    validate(result, load_schema("readiness.schema.json"))
    if result["merge"] is not False:
        raise PolicyError("readiness cannot authorize a merge")
    return result


def build_comment(disposition: dict[str, Any]) -> str:
    lines = [
        f"Disposition: {disposition['validity']}.",
        f"Decision: {disposition['decision']}.",
        _public_text(disposition["summary"]),
        "",
        "Affected users:",
        _public_text(disposition["affected_users"]),
        "",
        "Reasons:",
    ]
    lines.extend(f"- {_public_text(reason)}" for reason in disposition["reasons"])
    lines.append("")
    lines.append("Evidence:")
    for item in disposition["evidence"]:
        lines.append(
            f"- {item['source']}: {_public_text(item['ref'])} — {_public_text(item['note'])}"
        )
    contract = disposition.get("implementation_contract")
    if disposition["decision"] == "proceed" and contract:
        lines.extend(
            [
                "",
                "Implementation plan:",
                _public_text(contract["summary"]),
                (
                    f"Estimated scope: {contract['expected_files']} files, "
                    f"{contract['expected_lines']} changed lines."
                ),
                "",
                "Allowed paths:",
            ]
        )
        lines.extend(f"- {_public_text(path)}" for path in contract["allowed_paths"])
        lines.extend(["", "Required tests:"])
        lines.extend(f"- {_public_text(test)}" for test in contract["tests_required"])
        lines.extend(["", "Acceptance criteria:"])
        lines.extend(
            f"- {_public_text(criterion)}" for criterion in contract["acceptance_criteria"]
        )
    if disposition["label"] == "agent:breaking-risk":
        lines.extend(
            ["", "A maintainer has to decide whether this breaking change is worth making."]
        )
    body = "\n".join([DISPOSITION_MARKER, *lines])
    lowered = body.lower()
    if any(phrase in lowered for phrase in _COMMENT_BANS):
        raise PolicyError("disposition comment would contact the author or ask for more input")
    return body


def _public_text(text: str) -> str:
    cleaned = text.replace("@", "")
    for phrase in _COMMENT_BANS:
        if phrase == "@":
            continue
        cleaned = re.sub(re.escape(phrase), "", cleaned, flags=re.IGNORECASE)
    return " ".join(cleaned.split())


def contract_comment(issue_number: int, disposition: dict[str, Any]) -> str:
    import json

    payload = {
        "issue_number": issue_number,
        "decision": disposition["decision"],
        "implementation_contract": disposition["implementation_contract"],
    }
    return f"{CONTRACT_MARKER}\n```json\n{json.dumps(payload, indent=2)}\n```\n"


def diff_violations(
    changed_files: list[str],
    changed_lines: int,
    allowed_paths: list[str],
) -> list[str]:
    reasons: list[str] = []
    if not changed_files:
        reasons.append("agent made no changes")
    if changed_lines > MAX_LINES or len(changed_files) > MAX_FILES:
        reasons.append("diff exceeds the file or line cap")
    for path in changed_files:
        if path.startswith("/") or ".." in path.split("/"):
            reasons.append(f"{path} is not a relative path")
        elif _forbidden(path):
            reasons.append(f"{path} is forbidden")
        elif not _path_allowed(path, allowed_paths):
            reasons.append(f"{path} is outside the approved contract")
    return reasons


def terminal_label(labels: list[str]) -> str | None:
    for label in labels:
        if label in TERMINAL_LABELS:
            return label
    return None


def _breaking(flags: dict[str, bool]) -> bool:
    return any(flags.values())


def _docs_only(paths: list[str]) -> bool:
    return bool(paths) and all(path == "docs" or path.startswith("docs/") for path in paths)


def _paths_allowed(paths: list[str]) -> str | None:
    if not paths:
        return "contract has no allowed paths"
    for path in paths:
        if path.startswith("/") or ".." in path.split("/"):
            return f"path {path} is not relative"
        if _forbidden(path):
            return f"path {path} is forbidden"
    return None


def _size_allowed(files: int, lines: int) -> str | None:
    if files < 1 or files > MAX_FILES:
        return "expected file count is outside the cap"
    if lines < 1 or lines > MAX_LINES:
        return "expected line count is outside the cap"
    return None


def _forbidden(path: str) -> bool:
    if path in FORBIDDEN_EXACT:
        return True
    return any(path.startswith(prefix) for prefix in FORBIDDEN_PREFIXES)


def _path_allowed(path: str, allowed: list[str]) -> bool:
    for entry in allowed:
        if path == entry or path.startswith(entry.rstrip("/") + "/"):
            return True
    return False
