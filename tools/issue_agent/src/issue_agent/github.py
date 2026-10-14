"""Minimal GitHub client. It can comment and label, and it cannot merge."""

from __future__ import annotations

import json
import urllib.error
import urllib.parse
import urllib.request
from typing import Any, Callable

Opener = Callable[..., Any]


class GitHubError(RuntimeError):
    """Raised when the GitHub API rejects a request."""

    def __init__(self, message: str, *, status: int | None = None) -> None:
        super().__init__(message)
        self.status = status


AGENT_LABELS = {
    "agent:running": ("1d76db", "The issue agent opened a draft pull request."),
    "agent:user-error": ("d4c5f9", "The report matches documented behavior."),
    "agent:declined": ("d4c5f9", "The issue agent declined an automated change."),
    "agent:breaking-risk": ("b60205", "The proposed fix could be breaking."),
    "agent:insufficient-evidence": ("fbca04", "More evidence is needed for a safe change."),
    "agent:review-failed": ("b60205", "The issue-agent review or verification failed."),
}


class GitHubClient:
    def __init__(self, repository: str, token: str, opener: Opener | None = None) -> None:
        if not token:
            raise GitHubError("GITHUB_TOKEN is required for GitHub writes")
        self.repository = repository
        self._token = token
        self._opener = opener or urllib.request.urlopen

    def get_issue(self, number: int) -> dict[str, Any]:
        item = self._request("GET", f"/repos/{self.repository}/issues/{number}")
        comments = self._request(
            "GET",
            f"/repos/{self.repository}/issues/{number}/comments?per_page=100",
        )
        if "pull_request" in item:
            raise GitHubError("issue number points at a pull request")
        issue = _issue(item)
        issue["comments"] = [{"body": comment.get("body") or ""} for comment in comments]
        return issue

    def comment(self, number: int, body: str) -> None:
        self._request(
            "POST",
            f"/repos/{self.repository}/issues/{number}/comments",
            {"body": body},
        )

    def add_labels(self, number: int, labels: list[str]) -> None:
        for label in labels:
            self._ensure_label(label)
        self._request(
            "POST",
            f"/repos/{self.repository}/issues/{number}/labels",
            {"labels": labels},
        )

    def _ensure_label(self, name: str) -> None:
        metadata = AGENT_LABELS.get(name)
        if metadata is None:
            return
        encoded = urllib.parse.quote(name, safe="")
        try:
            self._request("GET", f"/repos/{self.repository}/labels/{encoded}")
            return
        except GitHubError as exc:
            if exc.status != 404:
                raise
        color, description = metadata
        try:
            self._request(
                "POST",
                f"/repos/{self.repository}/labels",
                {"name": name, "color": color, "description": description},
            )
        except GitHubError as exc:
            if exc.status != 422:
                raise
            self._request("GET", f"/repos/{self.repository}/labels/{encoded}")

    def create_draft_pull(self, *, title: str, body: str, head: str, base: str) -> dict[str, Any]:
        if not head.startswith("agent/issue-"):
            raise GitHubError("draft pull requests stay on agent/issue-* branches")
        owner = self.repository.split("/")[0]
        query = urllib.parse.urlencode({"head": f"{owner}:{head}", "state": "open"})
        existing = self._request("GET", f"/repos/{self.repository}/pulls?{query}")
        if existing:
            return existing[0]
        return self._request(
            "POST",
            f"/repos/{self.repository}/pulls",
            {"title": title, "body": body, "head": head, "base": base, "draft": True},
        )

    def _request(self, method: str, path: str, payload: dict[str, Any] | None = None) -> Any:
        data = None if payload is None else json.dumps(payload).encode("utf-8")
        request = urllib.request.Request(
            f"https://api.github.com{path}",
            data=data,
            method=method,
        )
        request.add_header("Authorization", f"Bearer {self._token}")
        request.add_header("Accept", "application/vnd.github+json")
        request.add_header("X-GitHub-Api-Version", "2022-11-28")
        if payload is not None:
            request.add_header("Content-Type", "application/json")
        try:
            with self._opener(request, timeout=30) as response:
                raw = response.read().decode("utf-8")
        except urllib.error.HTTPError as exc:
            raise GitHubError(
                f"GitHub API {method} {path} failed with {exc.code}",
                status=exc.code,
            ) from exc
        if not raw:
            return None
        return json.loads(raw)


class MemoryGitHub:
    """Test double that records writes."""

    def __init__(self, issues: list[dict[str, Any]] | None = None) -> None:
        self.repository = "splunk/terraform-provider-splunk"
        self.issues = list(issues or [])
        self.comments: list[tuple[int, str]] = []
        self.labels: list[tuple[int, list[str]]] = []
        self.pulls: list[dict[str, Any]] = []
        self.writes = 0

    def comment(self, number: int, body: str) -> None:
        self.writes += 1
        self.comments.append((number, body))

    def add_labels(self, number: int, labels: list[str]) -> None:
        self.writes += 1
        self.labels.append((number, list(labels)))

    def create_draft_pull(self, *, title: str, body: str, head: str, base: str) -> dict[str, Any]:
        if not head.startswith("agent/issue-"):
            raise GitHubError("draft pull requests stay on agent/issue-* branches")
        for pull in self.pulls:
            if pull["head"] == head and pull["state"] == "open":
                return pull
        pull = {
            "number": len(self.pulls) + 1,
            "html_url": f"https://github.com/{self.repository}/pull/{len(self.pulls) + 1}",
            "title": title,
            "body": body,
            "head": head,
            "base": base,
            "draft": True,
            "state": "open",
        }
        self.pulls.append(pull)
        self.writes += 1
        return pull


def _issue(item: dict[str, Any]) -> dict[str, Any]:
    return {
        "number": item["number"],
        "title": item.get("title") or "",
        "body": item.get("body") or "",
        "state": item.get("state") or "open",
        "labels": [label["name"] for label in item.get("labels", [])],
        "comments": [],
    }
