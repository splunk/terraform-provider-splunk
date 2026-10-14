"""Ad-hoc entry point for one issue. The agent runs locally and does not get a GitHub token."""

from __future__ import annotations

import os
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "src"))

from issue_agent.github import GitHubClient, GitHubError  # noqa: E402
from issue_agent.local_flow import run_issue, verify_workspace  # noqa: E402
from issue_agent.pipeline import Pipeline  # noqa: E402
from issue_agent.policy import ALLOWED_REPOSITORY, PolicyError  # noqa: E402
from issue_agent.publisher import ChangePublisher, SubprocessGit  # noqa: E402
from issue_agent.runner import AgentRunner, RunnerError  # noqa: E402


def main() -> int:
    number_raw = os.environ.get("ISSUE_NUMBER", "")
    mode = os.environ.get("ISSUE_AGENT_MODE", "")
    if not number_raw.isdigit() or int(number_raw) < 1:
        print("ISSUE_NUMBER must be a positive integer", file=sys.stderr)
        return 1
    repository = os.environ.get("GITHUB_REPOSITORY", ALLOWED_REPOSITORY)
    if repository != ALLOWED_REPOSITORY:
        print("repository is outside the allowlist", file=sys.stderr)
        return 1
    workspace = os.environ.get("GITHUB_WORKSPACE", os.getcwd())
    token = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN") or ""
    api_key = os.environ.get("CURSOR_API_KEY", "")
    try:
        runner = AgentRunner(
            api_key=api_key,
            model=os.environ.get("AGENT_MODEL", "composer-2.5"),
            repository=repository,
            cwd=workspace,
        )
        github = GitHubClient(repository, token)
        pipeline = Pipeline(
            runner=runner,
            github=github,
            mode="comment",
            repository=repository,
            implement_runtime="local",
            secrets=_secrets(),
        )
        result = run_issue(
            issue=github.get_issue(int(number_raw)),
            mode=mode,
            pipeline=pipeline,
            workspace=workspace,
            publisher=ChangePublisher(SubprocessGit(), github, token),
            verifier=verify_workspace,
        )
    except (PolicyError, RunnerError, GitHubError) as exc:
        print(str(exc).replace(token, "[redacted]") if token else str(exc), file=sys.stderr)
        return 1
    print(result.get("decision"), result.get("pull_request", {}).get("html_url", ""))
    disposition = result.get("disposition") or {}
    fix_failed = (
        mode == "review-and-fix"
        and disposition.get("decision") == "proceed"
        and "pull_request" not in result
    )
    if result.get("error") or fix_failed:
        return 1
    return 0


def _secrets() -> tuple[str, ...]:
    values = []
    for name in ("CURSOR_API_KEY", "GITHUB_TOKEN", "GH_TOKEN", "SPLUNK_PASSWORD"):
        value = os.environ.get(name)
        if value:
            values.append(value)
    return tuple(values)


if __name__ == "__main__":
    raise SystemExit(main())
