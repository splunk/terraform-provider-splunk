import json
import subprocess
import sys
import urllib.error
from pathlib import Path

from issue_agent.__main__ import main
from issue_agent.github import GitHubClient, MemoryGitHub
from issue_agent.scan import find_merge_calls

ROOT = Path(__file__).resolve().parents[3]


def test_repository_has_no_merge_api():
    assert find_merge_calls(ROOT) == []


def test_merge_scan_detects_a_call(tmp_path):
    workflow = tmp_path / ".github" / "workflows"
    workflow.mkdir(parents=True)
    (workflow / "bad.yml").write_text("run: gh pr merge --auto\n", encoding="utf-8")
    hits = find_merge_calls(tmp_path)
    assert hits


def test_clients_cannot_merge():
    assert not hasattr(GitHubClient, "merge")
    assert not hasattr(MemoryGitHub, "merge")


def test_github_requests_do_not_target_merge():
    seen = []

    class Response:
        def __init__(self, payload):
            self.payload = payload

        def read(self):
            return json.dumps(self.payload).encode("utf-8")

        def __enter__(self):
            return self

        def __exit__(self, *args):
            return False

    def opener(request, timeout):
        seen.append(request.full_url)
        if request.method == "GET" and request.full_url.endswith("/comments?per_page=100"):
            return Response([])
        if request.method == "GET":
            return Response(
                {
                    "number": 1,
                    "title": "Example",
                    "body": "",
                    "state": "open",
                    "labels": [{"name": "bug"}],
                }
            )
        return Response({})

    client = GitHubClient("splunk/terraform-provider-splunk", "token", opener=opener)
    issue = client.get_issue(1)
    client.comment(1, "Disposition: user_error.")
    assert issue["labels"] == ["bug"]
    assert all("/merge" not in url for url in seen)


def test_missing_agent_label_is_created_before_it_is_applied():
    seen = []

    class Response:
        def read(self):
            return b"{}"

        def __enter__(self):
            return self

        def __exit__(self, *args):
            return False

    def opener(request, timeout):
        seen.append((request.method, request.full_url))
        if request.method == "GET" and "/labels/agent%3Auser-error" in request.full_url:
            raise urllib.error.HTTPError(request.full_url, 404, "missing", {}, None)
        return Response()

    client = GitHubClient("splunk/terraform-provider-splunk", "token", opener=opener)
    client.add_labels(1, ["agent:user-error"])
    assert [method for method, _ in seen] == ["GET", "POST", "POST"]


def test_readiness_command_blocks_without_merging(tmp_path):
    checks = tmp_path / "checks.json"
    output = tmp_path / "ready.json"
    checks.write_text(
        json.dumps({"unit": "success", "lint": "success", "agent-review": "success"}),
        encoding="utf-8",
    )
    status = main(
        [
            "readiness",
            "--checks-file",
            str(checks),
            "--output",
            str(output),
            "--agent-pr",
            "--review-passed",
        ]
    )
    payload = json.loads(output.read_text(encoding="utf-8"))
    assert status == 1
    assert payload["merge"] is False
    assert payload["ready"] is False


def test_non_agent_review_passes_without_a_model(tmp_path):
    numstat = tmp_path / "numstat.txt"
    diff = tmp_path / "diff.txt"
    output = tmp_path / "review.json"
    numstat.write_text("1\t1\tsplunk/provider.go\n", encoding="utf-8")
    diff.write_text("diff\n", encoding="utf-8")
    status = main(
        [
            "review",
            "--numstat-file",
            str(numstat),
            "--diff-file",
            str(diff),
            "--output",
            str(output),
        ]
    )
    assert status == 0
    assert json.loads(output.read_text(encoding="utf-8"))["correction_requested"] is False


def test_verify_command_passes_for_this_repository():
    assert main(["verify-no-merge", "--root", str(ROOT)]) == 0


def test_ruleset_application_refuses_by_default():
    script = ROOT / "tools" / "issue_agent" / "scripts" / "apply_ruleset.py"
    completed = subprocess.run(
        [sys.executable, str(script)],
        check=False,
        capture_output=True,
        text=True,
    )
    assert completed.returncode == 2
    assert "Refusing" in completed.stderr


def test_ci_password_is_not_committed_and_image_is_pinned():
    workflow = (ROOT / ".github" / "workflows" / "test.yml").read_text(encoding="utf-8")
    terraform = (ROOT / "terraform" / "main.tf").read_text(encoding="utf-8")
    assert 'SPLUNK_PASSWORD: password' not in workflow
    assert "password = \"password\"" not in terraform
    assert "splunk/splunk:9.4.14" in workflow
    assert "services/server/info" in workflow
    assert "name: unit" in workflow
    assert "name: acceptance" in workflow
