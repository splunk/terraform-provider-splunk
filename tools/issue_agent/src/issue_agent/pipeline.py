"""Issue disposition, implementation, review, and readiness coordination."""

from __future__ import annotations

import json
import re
from typing import Any

from issue_agent.github import MemoryGitHub
from issue_agent.metrics import Metrics
from issue_agent.policy import (
    DISPOSITION_MARKER,
    MAX_COMMENT_CHARS,
    MAX_COMMENTS,
    MAX_DIFF_CHARS,
    MAX_ISSUE_BODY_CHARS,
    MAX_ISSUE_TEXT_CHARS,
    MODES,
    PolicyError,
    build_comment,
    contract_comment,
    normalize_disposition,
    normalize_review,
    readiness,
    terminal_label,
    validate_event,
)
from issue_agent.prompts import disposition_prompt, implementation_prompt, review_prompt
from issue_agent.redact import redact
from issue_agent.runner import AgentRunner, RunOutcome

_MENTION = re.compile(r"@[A-Za-z0-9-]+")


class Pipeline:
    def __init__(
        self,
        *,
        runner: AgentRunner,
        github: MemoryGitHub,
        mode: str,
        repository: str,
        metrics: Metrics | None = None,
        disposition_runtime: str = "local",
        implement_runtime: str = "local",
        secrets: tuple[str, ...] = (),
    ) -> None:
        if mode not in MODES:
            raise PolicyError(f"unknown mode {mode}")
        self.runner = runner
        self.github = github
        self.mode = mode
        self.repository = repository
        self.metrics = metrics or Metrics()
        self.disposition_runtime = disposition_runtime
        self.implement_runtime = implement_runtime
        self.secrets = secrets

    def handle_issue(self, issue: dict[str, Any], *, rerun: bool = False) -> dict[str, Any]:
        validate_event(
            {
                "repository": self.repository,
                "issue_state": issue.get("state", "open"),
                "is_fork": False,
            }
        )
        issue = _normalize_issue(issue)
        if not rerun and _already_dispositioned(issue):
            self.metrics.record(int(issue["number"]), "skipped")
            return {"issue_number": issue["number"], "skipped": True, "reason": "already dispositioned"}
        existing = None if rerun else terminal_label(list(issue.get("labels") or []))
        if existing:
            self.metrics.record(int(issue["number"]), "skipped")
            return {"issue_number": issue["number"], "skipped": True, "reason": existing}

        outcome = self.runner.complete(
            role="disposition",
            prompt=disposition_prompt(issue_text(issue, self.secrets)),
            runtime=self.disposition_runtime,
            idempotency_key=f"{self.repository}:{issue['number']}:disposition",
        )
        if outcome.exit_code != 0:
            self.metrics.record(int(issue["number"]), "errors")
            return {
                "issue_number": issue["number"],
                "error": outcome.status,
                "detail": outcome.result_text,
                "exit_code": outcome.exit_code,
                "agent_id": outcome.agent_id,
                "run_id": outcome.run_id,
            }

        try:
            disposition = normalize_disposition(extract_json(redact(outcome.result_text, self.secrets)))
        except (PolicyError, json.JSONDecodeError) as exc:
            self.metrics.record(int(issue["number"]), "errors")
            return {"issue_number": issue["number"], "error": str(exc), "exit_code": 2}

        kind = _metric_kind(disposition, self.mode, implemented=False)
        result: dict[str, Any] = {
            "issue_number": issue["number"],
            "decision": disposition["decision"],
            "label": disposition["label"],
            "validity": disposition["validity"],
            "disposition": disposition,
            "agent_id": outcome.agent_id,
            "run_id": outcome.run_id,
        }
        if self.mode == "dry-run":
            self.metrics.record(int(issue["number"]), kind)
            return result

        self._publish_disposition(int(issue["number"]), disposition)
        self.metrics.record(int(issue["number"]), kind)
        return result

    def implement(
        self,
        *,
        issue_number: int,
        disposition: dict[str, Any],
        revision: str = "",
        iteration: int = 1,
    ) -> dict[str, Any]:
        if self.mode == "dry-run":
            raise PolicyError("implementation is disabled in dry-run")
        if self.implement_runtime != "local":
            raise PolicyError("implementation runs only on the local runner")
        if disposition.get("decision") != "proceed" or not disposition.get("implementation_contract"):
            raise PolicyError("implementation requires a proceed contract")
        if iteration > 2:
            raise PolicyError("review correction limit reached")
        contract_text = json.dumps(disposition["implementation_contract"], indent=2)
        role = "revise" if revision else "implement"
        outcome = self.runner.complete(
            role=role,
            prompt=implementation_prompt(
                contract_text,
                issue_number=issue_number,
                revision=revision,
            ),
            runtime="local",
            idempotency_key=f"{self.repository}:{issue_number}:{role}:{iteration}",
        )
        return _outcome_dict(outcome)

    def review(
        self,
        *,
        contract: dict[str, Any],
        diff_text: str,
        changed_files: list[str],
        changed_lines: int,
        iteration: int,
        is_agent_pr: bool,
        is_fork: bool,
        workdir: str | None = None,
    ) -> dict[str, Any]:
        if not is_agent_pr:
            return {
                "decision": "pass",
                "correction_requested": False,
                "label": None,
                "summary": "not an agent pull request",
                "blocking": [],
                "non_blocking": [],
            }
        if len(diff_text) > MAX_DIFF_CHARS:
            self.metrics.record(0, "review_defects")
            return {
                "decision": "fail",
                "correction_requested": False,
                "label": "agent:review-failed",
                "summary": "diff exceeds the input cap",
                "blocking": [{"path": "", "summary": "diff exceeds the input cap"}],
                "non_blocking": [],
            }
        outcome = self.runner.complete(
            role="review",
            prompt=review_prompt(json.dumps(contract), diff_text),
            runtime="local",
            idempotency_key=f"{self.repository}:review:{iteration}:{','.join(changed_files)}",
            cwd=workdir,
        )
        if outcome.exit_code != 0:
            self.metrics.record(0, "errors")
            return {
                "decision": "fail",
                "correction_requested": False,
                "label": "agent:review-failed",
                "summary": outcome.status,
                "blocking": [{"path": "", "summary": "review agent did not finish"}],
                "non_blocking": [],
                "exit_code": outcome.exit_code,
            }
        reviewed = normalize_review(
            extract_json(redact(outcome.result_text, self.secrets)),
            changed_files=changed_files,
            changed_lines=changed_lines,
            allowed_paths=list(contract.get("allowed_paths") or []),
            iteration=iteration,
            is_agent_pr=is_agent_pr,
            is_fork=is_fork,
        )
        if reviewed["decision"] != "pass":
            self.metrics.record(0, "review_defects")
        return reviewed

    def judge_readiness(
        self,
        checks: dict[str, str],
        *,
        is_agent_pr: bool,
        review_passed: bool,
        issue_number: int | None = None,
        apply_label: bool = False,
    ) -> dict[str, Any]:
        result = readiness(checks, is_agent_pr=is_agent_pr, review_passed=review_passed)
        if apply_label and result["label"] and issue_number is not None and self.mode != "dry-run":
            self.github.add_labels(issue_number, [result["label"]])
        return result

    def _publish_disposition(self, number: int, disposition: dict[str, Any]) -> None:
        self.github.comment(number, build_comment(disposition))
        if disposition["label"]:
            self.github.add_labels(number, [disposition["label"]])
        if disposition["decision"] == "proceed" and self.mode != "dry-run":
            self.github.comment(number, contract_comment(number, disposition))


def issue_text(issue: dict[str, Any], secrets: tuple[str, ...] = ()) -> str:
    comments = list(issue.get("comments") or [])[-MAX_COMMENTS:]
    comment_text = "\n".join((item.get("body") or "")[:MAX_COMMENT_CHARS] for item in comments)
    raw = (
        f"number: {issue['number']}\n"
        f"title: {(issue.get('title') or '')[:300]}\n"
        f"body:\n{(issue.get('body') or '')[:MAX_ISSUE_BODY_CHARS]}\n"
        f"comments:\n{comment_text}\n"
    )
    text = _MENTION.sub("", redact(raw, extra=secrets))
    if len(text) > MAX_ISSUE_TEXT_CHARS:
        raise PolicyError("issue text exceeds the input cap")
    return text


def extract_json(text: str) -> dict[str, Any]:
    fenced = re.search(r"```(?:json)?\s*(\{.*\})\s*```", text, re.DOTALL)
    payload = fenced.group(1) if fenced else text
    start = payload.find("{")
    end = payload.rfind("}")
    if start < 0 or end <= start:
        raise PolicyError("agent did not return a JSON object")
    parsed = json.loads(payload[start : end + 1])
    if not isinstance(parsed, dict):
        raise PolicyError("agent JSON must be an object")
    return parsed


def _already_dispositioned(issue: dict[str, Any]) -> bool:
    for item in issue.get("comments") or []:
        if DISPOSITION_MARKER in (item.get("body") or ""):
            return True
    return False


def _normalize_issue(issue: dict[str, Any]) -> dict[str, Any]:
    labels = []
    for label in issue.get("labels") or []:
        labels.append(label["name"] if isinstance(label, dict) else str(label))
    normalized = dict(issue)
    normalized["labels"] = labels
    normalized["body"] = issue.get("body") or ""
    normalized["state"] = issue.get("state") or "open"
    return normalized


def _metric_kind(disposition: dict[str, Any], mode: str, *, implemented: bool) -> str:
    if disposition["decision"] == "proceed":
        if mode == "implement" and implemented:
            return "implemented"
        return "would_implement"
    label = disposition.get("label")
    if label == "agent:user-error":
        return "user_error"
    if label == "agent:breaking-risk":
        return "breaking_or_harmful"
    if label == "agent:insufficient-evidence":
        return "insufficient_evidence"
    return "declined_other"


def _outcome_dict(outcome: RunOutcome) -> dict[str, Any]:
    return {
        "status": outcome.status,
        "exit_code": outcome.exit_code,
        "agent_id": outcome.agent_id,
        "run_id": outcome.run_id,
        "started": outcome.started,
    }
