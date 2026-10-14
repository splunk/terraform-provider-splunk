import pytest
import tempfile
from pathlib import Path

from issue_agent.github import MemoryGitHub
from issue_agent.local_flow import VerifyResult, run_issue
from issue_agent.policy import ALLOWED_REPOSITORY, PolicyError
from issue_agent.publisher import ChangePublisher
from issue_agent.runner import AgentRunner, RunnerError
from support import disposition, outcome, review_payload, ScriptedRunner
from test_pipeline import _issue, _pipeline


class FakeGit:
    def __init__(self, name="splunk/resource_example.go", fail_push=False):
        self.calls = []
        self.name = name
        self.fail_push = fail_push
        self.control_dir = tempfile.TemporaryDirectory()
        self.git_dir = Path(self.control_dir.name)
        (self.git_dir / "config").write_text("[core]\n", encoding="utf-8")

    def run(self, args, *, cwd, env=None):
        self.calls.append(list(args))
        if "push" in args:
            if self.fail_push:
                raise RuntimeError("remote rejected test-github-token")
            return ""
        if args[:4] == ["git", "diff", "--name-status", "HEAD"]:
            return f"M\t{self.name}\n"
        if args[:4] == ["git", "diff", "--numstat", "HEAD"]:
            return f"3\t1\t{self.name}\n"
        if args[:3] == ["git", "diff", "HEAD"]:
            return "diff --git a/splunk/resource_example.go b/splunk/resource_example.go\n"
        if args[:2] == ["git", "ls-files"]:
            return ""
        if args[:3] == ["git", "rev-parse", "--absolute-git-dir"]:
            return f"{self.git_dir}\n"
        if args[:3] == ["git", "rev-parse", "HEAD"]:
            return "abc123\n"
        return ""


def _verifier(cwd, files):
    return VerifyResult(True, "passed", [("unit", "go test ./... passed with TF_ACC empty")])


def _publisher(git, github, token="test-github-token"):
    return ChangePublisher(git, github, token)


def test_review_only_does_not_publish():
    git = FakeGit()
    github = MemoryGitHub()
    runner = ScriptedRunner([outcome(disposition())])
    result = run_issue(
        issue=_issue(12),
        mode="review-only",
        pipeline=_pipeline(runner, github=github),
        workspace="/work",
        publisher=_publisher(git, github),
        verifier=_verifier,
    )
    assert result["decision"] == "proceed"
    assert github.pulls == []
    assert git.calls == []
    assert [call["role"] for call in runner.calls] == ["disposition"]


def test_review_and_fix_opens_one_local_draft():
    git = FakeGit()
    github = MemoryGitHub()
    runner = ScriptedRunner(
        [
            outcome(disposition()),
            outcome("{}"),
            outcome(review_payload()),
        ]
    )
    result = run_issue(
        issue=_issue(12),
        mode="review-and-fix",
        pipeline=_pipeline(runner, github=github),
        workspace="/work",
        publisher=_publisher(git, github),
        verifier=_verifier,
    )
    assert result["pull_request"]["draft"] is True
    assert result["pull_request"]["branch"] == "agent/issue-12"
    assert len(github.pulls) == 1
    assert github.pulls[0]["head"] == "agent/issue-12"
    assert github.pulls[0]["base"] == "master"
    assert "Fixes #12" in github.pulls[0]["body"]
    assert "test-github-token" not in github.pulls[0]["body"]
    roles = [call["role"] for call in runner.calls]
    assert roles == ["disposition", "implement", "review"]
    assert runner.calls[1]["runtime"] == "local"
    assert "Issue number: 12" in runner.calls[1]["prompt"]
    assert "Do not run git commit" in runner.calls[1]["prompt"]
    assert runner.calls[2]["cwd"] != "/work"
    pushed = [call for call in git.calls if "push" in call]
    assert pushed == [
        [
            "git",
            "-c",
            "core.hooksPath=/dev/null",
            "-c",
            "http.extraheader=AUTHORIZATION: bearer test-github-token",
            "push",
            "origin",
            "HEAD:refs/heads/agent/issue-12",
        ]
    ]
    assert not any(call[:2] == ["git", "config"] for call in git.calls)
    assert any(labels == ["agent:running"] for _, labels in github.labels)


def test_forbidden_diff_is_not_pushed():
    git = FakeGit(name=".github/workflows/test.yml")
    github = MemoryGitHub()
    runner = ScriptedRunner([outcome(disposition()), outcome("{}")])
    result = run_issue(
        issue=_issue(4),
        mode="review-and-fix",
        pipeline=_pipeline(runner, github=github),
        workspace="/work",
        publisher=_publisher(git, github),
        verifier=_verifier,
    )
    assert result["decision"] == "decline"
    assert github.pulls == []
    assert not any("push" in call for call in git.calls)
    assert not any(call[:3] == ["git", "checkout", "-B"] for call in git.calls)
    assert any(labels == ["agent:review-failed"] for _, labels in github.labels)


def test_review_correction_publishes_once():
    git = FakeGit()
    github = MemoryGitHub()
    failing = review_payload(
        within_contract=False,
        compatibility="breaking",
        blocking=[{"path": "splunk/resource_example.go", "summary": "missed a test"}],
    )
    runner = ScriptedRunner(
        [
            outcome(disposition()),
            outcome("{}"),
            outcome(failing),
            outcome("{}"),
            outcome(review_payload()),
        ]
    )
    result = run_issue(
        issue=_issue(8),
        mode="review-and-fix",
        pipeline=_pipeline(runner, github=github),
        workspace="/work",
        publisher=_publisher(git, github),
        verifier=_verifier,
    )
    assert result["pull_request"]["branch"] == "agent/issue-8"
    assert len(github.pulls) == 1
    assert [call["role"] for call in runner.calls] == [
        "disposition",
        "implement",
        "review",
        "revise",
        "review",
    ]
    assert sum(1 for call in git.calls if call[:3] == ["git", "fetch", "--no-tags"]) == 1


def test_push_error_redacts_the_token():
    git = FakeGit(fail_push=True)
    github = MemoryGitHub()
    publisher = _publisher(git, github)
    publisher.prepare_from_master("/work")
    with pytest.raises(PolicyError) as caught:
        publisher.publish(
            "/work",
            issue_number=3,
            title="Fix #3",
            body="Fixes #3\n",
            allowed_paths=["splunk/resource_example.go", "splunk/resource_example_test.go"],
        )
    assert "test-github-token" not in str(caught.value)
    assert "[redacted]" in str(caught.value)


def test_runner_rejects_cloud_runtime():
    runner = AgentRunner(
        api_key="cursor_test_key",
        repository=ALLOWED_REPOSITORY,
        cwd=".",
        sdk_factory=lambda: {},
    )
    with pytest.raises(RunnerError, match="local"):
        runner.complete(
            role="implement",
            prompt="implement",
            runtime="cloud",
            idempotency_key="issue-1",
        )
