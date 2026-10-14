"""Commit a validated local diff and open one draft pull request.

The Cursor agent never receives the GitHub token. This module runs only after
the agent process has returned, and it refuses every ref outside agent/issue-*.
"""

from __future__ import annotations

import os
import subprocess
from pathlib import Path
from typing import Any, Protocol

from issue_agent.policy import PolicyError, diff_violations

SECRET_ENV = (
    "CURSOR_API_KEY",
    "GITHUB_TOKEN",
    "GH_TOKEN",
    "SPLUNK_PASSWORD",
    "SPLUNK_CI_PASSWORD",
)

COMMIT_ENV = {
    "GIT_AUTHOR_NAME": "github-actions[bot]",
    "GIT_AUTHOR_EMAIL": "41898282+github-actions[bot]@users.noreply.github.com",
    "GIT_COMMITTER_NAME": "github-actions[bot]",
    "GIT_COMMITTER_EMAIL": "41898282+github-actions[bot]@users.noreply.github.com",
}


class GitRunner(Protocol):
    def run(self, args: list[str], *, cwd: str, env: dict[str, str] | None = None) -> str:
        """Run one git command and return stdout."""


class SubprocessGit:
    def run(self, args: list[str], *, cwd: str, env: dict[str, str] | None = None) -> str:
        merged = os.environ.copy()
        for name in SECRET_ENV:
            merged.pop(name, None)
        merged["GIT_CONFIG_GLOBAL"] = os.devnull
        merged["GIT_CONFIG_NOSYSTEM"] = "1"
        if env:
            merged.update(env)
        completed = subprocess.run(
            args,
            cwd=cwd,
            env=merged,
            capture_output=True,
            text=True,
        )
        if completed.returncode != 0:
            detail = (completed.stderr or completed.stdout or "git command failed").strip()
            raise PolicyError(detail)
        return completed.stdout


class ChangePublisher:
    def __init__(self, git: GitRunner, github: Any, token: str) -> None:
        self.git = git
        self.github = github
        self.token = token
        self._git_dir: Path | None = None
        self._git_dir_identity: tuple[int, int] | None = None
        self._git_config: bytes | None = None
        self._git_attributes: bytes | None = None

    def prepare_from_master(self, cwd: str) -> None:
        self._run(["git", "fetch", "--no-tags", "origin", "master"], cwd)
        self._run(["git", "checkout", "--detach", "origin/master"], cwd)
        git_dir = Path(self._run(["git", "rev-parse", "--absolute-git-dir"], cwd).strip())
        stat = git_dir.stat()
        config = git_dir / "config"
        self._git_dir = git_dir
        self._git_dir_identity = (stat.st_dev, stat.st_ino)
        self._git_config = config.read_bytes()
        attributes = git_dir / "info" / "attributes"
        self._git_attributes = attributes.read_bytes() if attributes.exists() else None

    def inspect(self, cwd: str, allowed_paths: list[str]) -> dict[str, Any]:
        self._restore_git_control()
        files, lines, diff = self._collect(cwd)
        _reject_symlinks(cwd, files)
        violations = diff_violations(files, lines, allowed_paths)
        if violations:
            raise PolicyError("; ".join(violations))
        return {"files": files, "lines": lines, "diff": diff}

    def publish(
        self,
        cwd: str,
        *,
        issue_number: int,
        title: str,
        body: str,
        allowed_paths: list[str],
        base: str = "master",
    ) -> dict[str, Any]:
        if issue_number < 1:
            raise PolicyError("issue number must be positive")
        if not self.token:
            raise PolicyError("GITHUB_TOKEN is required to publish a draft pull request")
        branch = f"agent/issue-{issue_number}"
        ref = f"refs/heads/{branch}"
        if not ref.startswith("refs/heads/agent/issue-"):
            raise PolicyError("refusing to push a branch outside agent/issue-*")
        inspected = self.inspect(cwd, allowed_paths)
        no_hooks = ["git", "-c", f"core.hooksPath={os.devnull}"]
        self._run([*no_hooks, "checkout", "-B", branch], cwd)
        self._run(["git", "add", "--", *inspected["files"]], cwd, env=None)
        self._run([*no_hooks, "commit", "-m", title, "-m", body], cwd, env=COMMIT_ENV)
        sha = self._run(["git", "rev-parse", "HEAD"], cwd).strip()
        self._run(
            [
                *no_hooks,
                "-c",
                f"http.extraheader=AUTHORIZATION: bearer {self.token}",
                "push",
                "origin",
                f"HEAD:{ref}",
            ],
            cwd,
        )
        pull = self.github.create_draft_pull(title=title, body=body, head=branch, base=base)
        return {
            "sha": sha,
            "branch": branch,
            "html_url": pull.get("html_url", ""),
            "number": pull.get("number"),
            "draft": True,
        }

    def _collect(self, cwd: str) -> tuple[list[str], int, str]:
        name_status = self._run(["git", "diff", "--name-status", "HEAD"], cwd)
        numstat = self._run(["git", "diff", "--numstat", "HEAD"], cwd)
        untracked = self._run(["git", "ls-files", "--others", "--exclude-standard"], cwd)
        diff = self._run(["git", "diff", "--no-ext-diff", "--no-textconv", "HEAD"], cwd)
        files = _files_from_name_status(name_status)
        extra: list[str] = []
        for path in untracked.splitlines():
            cleaned = path.strip()
            if not cleaned or cleaned in files:
                continue
            files.append(cleaned)
            extra.append(cleaned)
            diff += _untracked_diff(cwd, cleaned)
        return files, _line_count(numstat) + _untracked_lines(cwd, extra), diff

    def _restore_git_control(self) -> None:
        if self._git_dir is None or self._git_dir_identity is None or self._git_config is None:
            raise PolicyError("repository control state was not captured")
        stat = self._git_dir.stat()
        if (stat.st_dev, stat.st_ino) != self._git_dir_identity:
            raise PolicyError("repository metadata changed during the agent run")
        config = self._git_dir / "config"
        if config.is_symlink():
            config.unlink()
        config.write_bytes(self._git_config)
        attributes = self._git_dir / "info" / "attributes"
        if attributes.is_symlink():
            attributes.unlink()
        if self._git_attributes is None:
            attributes.unlink(missing_ok=True)
        else:
            attributes.parent.mkdir(parents=True, exist_ok=True)
            attributes.write_bytes(self._git_attributes)

    def _run(self, args: list[str], cwd: str, env: dict[str, str] | None = None) -> str:
        if args[:2] == ["git", "config"]:
            raise PolicyError("refusing to change git config")
        try:
            return self.git.run(args, cwd=cwd, env=env)
        except Exception as exc:
            message = str(exc)
            if self.token:
                message = message.replace(self.token, "[redacted]")
            raise PolicyError(message) from None


def _files_from_name_status(text: str) -> list[str]:
    files: list[str] = []
    for line in text.splitlines():
        if not line.strip():
            continue
        parts = line.split("\t")
        if len(parts) < 2:
            raise PolicyError("git name-status output was incomplete")
        status = parts[0]
        if status.startswith(("R", "C")) and len(parts) >= 3:
            files.extend(parts[1:])
        else:
            files.append(parts[-1])
    return files


def _line_count(numstat: str) -> int:
    total = 0
    for line in numstat.splitlines():
        if not line.strip():
            continue
        added, removed, _path = line.split("\t", 2)
        if added != "-":
            total += int(added)
        if removed != "-":
            total += int(removed)
    return total


def _untracked_lines(cwd: str, paths: list[str]) -> int:
    total = 0
    for path in paths:
        candidate = _safe_file(cwd, path)
        if candidate is None:
            continue
        total += sum(1 for _ in candidate.open(encoding="utf-8", errors="replace"))
    return total


def _untracked_diff(cwd: str, path: str) -> str:
    candidate = _safe_file(cwd, path)
    content = "" if candidate is None else candidate.read_text(encoding="utf-8", errors="replace")
    return f"\ndiff --git a/{path} b/{path}\nnew file mode 100644\n{content}\n"


def _safe_file(cwd: str, path: str) -> Path | None:
    if path.startswith("/") or ".." in path.split("/"):
        return None
    candidate = Path(cwd, path)
    if candidate.is_symlink():
        raise PolicyError(f"{path} is a symbolic link")
    root = Path(cwd).resolve()
    try:
        candidate.resolve().relative_to(root)
    except ValueError:
        raise PolicyError(f"{path} resolves outside the workspace") from None
    if candidate.is_file():
        return candidate
    return None


def _reject_symlinks(cwd: str, paths: list[str]) -> None:
    for path in paths:
        candidate = Path(cwd, path)
        if candidate.is_symlink():
            raise PolicyError(f"{path} is a symbolic link")
