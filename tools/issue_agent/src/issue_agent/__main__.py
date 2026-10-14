"""Command line for review, readiness, and implementation checks."""

from __future__ import annotations

import argparse
import json
import os
import sys
from pathlib import Path
from typing import Any, Sequence

from issue_agent.github import GitHubClient, MemoryGitHub
from issue_agent.pipeline import Pipeline
from issue_agent.policy import ALLOWED_REPOSITORY, PolicyError, readiness
from issue_agent.runner import DEFAULT_MODEL, AgentRunner, RunnerError
from issue_agent.scan import find_merge_calls


def main(argv: Sequence[str] | None = None) -> int:
    parser = argparse.ArgumentParser(prog="issue-agent")
    sub = parser.add_subparsers(dest="command", required=True)

    review = sub.add_parser("review")
    review.add_argument("--contract-file")
    review.add_argument("--numstat-file", required=True)
    review.add_argument("--diff-file", required=True)
    review.add_argument("--iteration", type=int, default=1)
    review.add_argument("--agent-pr", action="store_true")
    review.add_argument("--fork", action="store_true")
    review.add_argument("--output", required=True)
    review.add_argument("--cwd", default=os.getcwd())
    review.add_argument("--repository", default=ALLOWED_REPOSITORY)
    review.set_defaults(func=cmd_review)

    ready = sub.add_parser("readiness")
    ready.add_argument("--checks-file", required=True)
    ready.add_argument("--output", required=True)
    ready.add_argument("--agent-pr", action="store_true")
    ready.add_argument("--review-passed", action="store_true")
    ready.add_argument("--issue-number", type=int)
    ready.add_argument("--apply-label", action="store_true")
    ready.add_argument("--repository", default=ALLOWED_REPOSITORY)
    ready.set_defaults(func=cmd_readiness)

    verify = sub.add_parser("verify-no-merge")
    verify.add_argument("--root", default=".")
    verify.set_defaults(func=cmd_verify)

    implement = sub.add_parser("implement")
    implement.add_argument("--bundle-file", required=True)
    implement.add_argument("--repository", default=ALLOWED_REPOSITORY)
    implement.add_argument("--cwd", default=os.getcwd())
    implement.add_argument("--revision-file")
    implement.add_argument("--iteration", type=int, default=1)
    implement.add_argument("--output", required=True)
    implement.set_defaults(func=cmd_implement)

    args = parser.parse_args(list(argv) if argv is not None else None)
    try:
        return int(args.func(args))
    except (PolicyError, RunnerError, json.JSONDecodeError) as exc:
        print(str(exc), file=sys.stderr)
        return 1


def cmd_review(args: argparse.Namespace) -> int:
    files, lines = parse_numstat(Path(args.numstat_file).read_text(encoding="utf-8"))
    if not args.agent_pr:
        payload = {
            "decision": "pass",
            "correction_requested": False,
            "label": None,
            "summary": "not an agent pull request",
            "blocking": [],
            "non_blocking": [],
        }
        _write(args.output, payload)
        return 0
    if not args.contract_file:
        raise PolicyError("agent pull requests require a contract file")
    contract = json.loads(Path(args.contract_file).read_text(encoding="utf-8"))
    pipeline = Pipeline(
        runner=_runner(args.repository, args.cwd),
        github=MemoryGitHub(),
        mode="dry-run",
        repository=args.repository,
        secrets=_secrets(),
    )
    result = pipeline.review(
        contract=contract,
        diff_text=Path(args.diff_file).read_text(encoding="utf-8"),
        changed_files=files,
        changed_lines=lines,
        iteration=args.iteration,
        is_agent_pr=True,
        is_fork=args.fork,
    )
    _write(args.output, result)
    return 0 if result["decision"] == "pass" else 2


def cmd_readiness(args: argparse.Namespace) -> int:
    checks = json.loads(Path(args.checks_file).read_text(encoding="utf-8"))
    result = readiness(
        checks,
        is_agent_pr=args.agent_pr,
        review_passed=args.review_passed,
    )
    if args.apply_label and result["label"] and args.issue_number is not None:
        _live_github(args.repository).add_labels(args.issue_number, [result["label"]])
    _write(args.output, result)
    if result["merge"] is not False:
        raise PolicyError("readiness cannot authorize a merge")
    return 0 if result["ready"] else 1


def cmd_implement(args: argparse.Namespace) -> int:
    bundle = json.loads(Path(args.bundle_file).read_text(encoding="utf-8"))
    disposition = bundle["disposition"]
    revision = ""
    if args.revision_file:
        revision = Path(args.revision_file).read_text(encoding="utf-8")
    pipeline = Pipeline(
        runner=_runner(args.repository, args.cwd),
        github=_live_github(args.repository),
        mode="implement",
        repository=args.repository,
        implement_runtime="local",
        secrets=_secrets(),
    )
    result = pipeline.implement(
        issue_number=int(bundle["issue_number"]),
        disposition=disposition,
        revision=revision,
        iteration=args.iteration,
    )
    _write(args.output, result)
    return 0 if result["exit_code"] == 0 else int(result["exit_code"])


def cmd_verify(args: argparse.Namespace) -> int:
    hits = find_merge_calls(Path(args.root).resolve())
    if hits:
        print("\n".join(hits), file=sys.stderr)
        return 1
    print("no merge API calls")
    return 0


def _runner(repository: str, cwd: str) -> AgentRunner:
    api_key = os.environ.get("CURSOR_API_KEY", "")
    record = Path(os.environ.get("AGENT_RUN_RECORD", "agent-run.json"))
    return AgentRunner(
        api_key=api_key,
        model=os.environ.get("AGENT_MODEL", DEFAULT_MODEL),
        repository=repository,
        cwd=cwd,
        run_record_path=record,
    )


def _live_github(repository: str) -> GitHubClient:
    token = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN") or ""
    return GitHubClient(repository, token)


def _secrets() -> tuple[str, ...]:
    values = []
    for name in ("CURSOR_API_KEY", "GITHUB_TOKEN", "GH_TOKEN", "SPLUNK_PASSWORD"):
        value = os.environ.get(name)
        if value:
            values.append(value)
    return tuple(values)


def _write(path: str, payload: Any) -> None:
    Path(path).write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")


def parse_numstat(text: str) -> tuple[list[str], int]:
    files: list[str] = []
    total = 0
    for line in text.splitlines():
        if not line.strip():
            continue
        added, removed, path = line.split("\t")
        files.append(path)
        if added != "-":
            total += int(added)
        if removed != "-":
            total += int(removed)
    return files, total
