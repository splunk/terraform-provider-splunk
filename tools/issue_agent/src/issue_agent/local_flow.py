"""Run one explicit issue through local disposition, fix, review, and draft PR."""

from __future__ import annotations

import json
import os
import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Callable

from issue_agent.pipeline import Pipeline
from issue_agent.policy import REVIEW_MARKER, PolicyError
from issue_agent.publisher import ChangePublisher
from issue_agent.redact import redact

Verifier = Callable[[str, list[str]], "VerifyResult"]


@dataclass
class VerifyResult:
    ok: bool
    detail: str
    checks: list[tuple[str, str]]


def run_issue(
    *,
    issue: dict[str, Any],
    mode: str,
    pipeline: Pipeline,
    workspace: str,
    publisher: ChangePublisher,
    verifier: Verifier,
) -> dict[str, Any]:
    if mode not in {"review-only", "review-and-fix"}:
        raise PolicyError("mode must be review-only or review-and-fix")
    pipeline.mode = "comment"
    result = pipeline.handle_issue(issue, rerun=True)
    if result.get("skipped") or result.get("error") or result.get("decision") != "proceed":
        return result
    if mode == "review-only":
        return result

    disposition = result["disposition"]
    contract = disposition["implementation_contract"]
    issue_number = int(issue["number"])
    revision = ""
    last_review: dict[str, Any] | None = None
    for iteration in (1, 2):
        pipeline.mode = "implement"
        if iteration == 1:
            publisher.prepare_from_master(workspace)
        outcome = pipeline.implement(
            issue_number=issue_number,
            disposition=disposition,
            revision=revision,
            iteration=iteration,
        )
        if outcome["exit_code"] != 0:
            _fail(pipeline, issue_number, f"implementation agent status: {outcome['status']}")
            return {**result, "implementation": outcome, "decision": "decline"}
        try:
            inspected = publisher.inspect(workspace, list(contract["allowed_paths"]))
        except PolicyError as exc:
            _fail(pipeline, issue_number, str(exc))
            return {**result, "decision": "decline", "error": str(exc)}
        verified = verifier(workspace, inspected["files"])
        if not verified.ok:
            _fail(pipeline, issue_number, verified.detail)
            return {**result, "decision": "decline", "error": "verification failed"}
        with tempfile.TemporaryDirectory(prefix="issue-agent-review-") as review_dir:
            Path(review_dir, "diff.txt").write_text(inspected["diff"], encoding="utf-8")
            last_review = pipeline.review(
                contract=contract,
                diff_text=inspected["diff"],
                changed_files=inspected["files"],
                changed_lines=inspected["lines"],
                iteration=iteration,
                is_agent_pr=True,
                is_fork=False,
                workdir=review_dir,
            )
        if last_review["decision"] == "pass":
            try:
                published = publisher.publish(
                    workspace,
                    issue_number=issue_number,
                    title=_title(issue_number, disposition),
                    body=_pull_body(issue_number, disposition, verified.checks),
                    allowed_paths=list(contract["allowed_paths"]),
                )
            except PolicyError as exc:
                _fail(pipeline, issue_number, str(exc))
                return {**result, "decision": "decline", "error": str(exc)}
            pipeline.github.add_labels(issue_number, ["agent:running"])
            pipeline.metrics.record(issue_number, "implemented")
            return {**result, "pull_request": published, "review": last_review}
        if not last_review.get("correction_requested"):
            _fail(pipeline, issue_number, json.dumps(last_review["blocking"]))
            return {**result, "review": last_review, "decision": "decline"}
        revision = json.dumps(last_review)
    _fail(pipeline, issue_number, "review correction limit reached")
    return {**result, "review": last_review, "decision": "decline"}


def verify_workspace(cwd: str, files: list[str]) -> VerifyResult:
    go_files = [path for path in files if path.endswith(".go")]
    if go_files:
        formatted = subprocess.run(
            ["gofmt", "-l", *go_files],
            cwd=cwd,
            capture_output=True,
            text=True,
        )
        if formatted.returncode != 0:
            return VerifyResult(False, formatted.stderr[-2000:], [])
        if formatted.stdout.strip():
            return VerifyResult(False, "gofmt would change files:\n" + formatted.stdout[:1000], [])
    env = os.environ.copy()
    env["TF_ACC"] = ""
    for name in ("CURSOR_API_KEY", "GITHUB_TOKEN", "GH_TOKEN", "SPLUNK_PASSWORD", "SPLUNK_CI_PASSWORD"):
        env.pop(name, None)
    tested = subprocess.run(
        ["go", "test", "./..."],
        cwd=cwd,
        env=env,
        capture_output=True,
        text=True,
    )
    if tested.returncode != 0:
        detail = redact((tested.stdout + "\n" + tested.stderr)[-4000:])
        return VerifyResult(False, detail, [])
    checks = [("unit", "go test ./... passed with TF_ACC empty")]
    if go_files:
        checks.append(("lint", "gofmt -l reported no changed Go files"))
    return VerifyResult(True, "passed", checks)


def _fail(pipeline: Pipeline, issue_number: int, detail: str) -> None:
    pipeline.github.comment(issue_number, f"{REVIEW_MARKER}\n{_public(detail)}")
    pipeline.github.add_labels(issue_number, ["agent:review-failed"])


def _public(text: str) -> str:
    return " ".join(redact(text).replace("@", "").split())


def _title(issue_number: int, disposition: dict[str, Any]) -> str:
    compact = " ".join(_public(disposition["summary"]).split())
    action = "Implement" if disposition["validity"] == "feature_request" else "Fix"
    return f"{action} #{issue_number}: {compact[:80]}".strip()


def _pull_body(issue_number: int, disposition: dict[str, Any], checks: list[tuple[str, str]]) -> str:
    closing = "Closes" if disposition["validity"] == "feature_request" else "Fixes"
    lines = [
        f"{closing} #{issue_number}",
        "",
        _public(disposition["summary"]),
        "",
        "The agent workflow ran these checks before opening this draft:",
    ]
    lines.extend(f"- {name}: {_public(summary)}" for name, summary in checks)
    lines.extend(
        [
            "",
            "A maintainer runs the repository test workflow on this branch before merge.",
            "The workflow token that opened this draft does not start other workflows.",
            "",
            f"Issue-Agent: {issue_number}",
        ]
    )
    return "\n".join(lines) + "\n"
