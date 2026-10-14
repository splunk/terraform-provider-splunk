"""Cursor SDK runner. Runtime is always explicit, and runs are always closed."""

from __future__ import annotations

import json
import os
import sys
from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import TimeoutError as FuturesTimeout
from contextlib import contextmanager
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Callable, Iterator

from issue_agent.policy import MAX_RETRIES, TIMEOUT_SECONDS, TOKEN_BUDGET

# composer-2.5 is absent from this API key's catalog. `default` is the
# server-selected model that catalog does list.
DEFAULT_MODEL = "default"

SECRET_ENV = (
    "CURSOR_API_KEY",
    "GITHUB_TOKEN",
    "GH_TOKEN",
    "SPLUNK_PASSWORD",
    "SPLUNK_CI_PASSWORD",
)

Sleeper = Callable[[float], None]


@dataclass
class RunOutcome:
    status: str
    result_text: str
    agent_id: str
    run_id: str
    exit_code: int
    started: bool


class RunnerError(RuntimeError):
    """Raised when the runner is configured incorrectly."""


def default_sdk_factory() -> dict[str, Any]:
    from cursor_sdk import Agent, AgentOptions, CursorAgentError, LocalAgentOptions

    return {
        "Agent": Agent,
        "AgentOptions": AgentOptions,
        "CursorAgentError": CursorAgentError,
        "LocalAgentOptions": LocalAgentOptions,
    }


class AgentRunner:
    def __init__(
        self,
        *,
        api_key: str,
        model: str = DEFAULT_MODEL,
        repository: str,
        cwd: str | None = None,
        sdk_factory: Callable[[], dict[str, Any]] | None = None,
        sleeper: Sleeper | None = None,
        timeout_seconds: int = TIMEOUT_SECONDS,
        max_retries: int = MAX_RETRIES,
        run_record_path: Path | None = None,
    ) -> None:
        if not api_key:
            raise RunnerError("CURSOR_API_KEY is required")
        if not repository:
            raise RunnerError("repository is required")
        self.api_key = api_key
        self.model = model
        self.repository = repository
        self.cwd = cwd
        self.sdk_factory = sdk_factory or default_sdk_factory
        self.sleeper = sleeper or _sleep
        self.timeout_seconds = timeout_seconds
        self.max_retries = max_retries
        self.run_record_path = run_record_path

    def complete(
        self,
        *,
        role: str,
        prompt: str,
        runtime: str,
        idempotency_key: str,
        cwd: str | None = None,
    ) -> RunOutcome:
        if runtime != "local":
            raise RunnerError("runtime must be local")
        error_type: type[BaseException] = RunnerError
        last_error: Exception | None = None
        workdir = cwd or self.cwd
        for attempt in range(self.max_retries + 1):
            agent = None
            try:
                with hidden_secrets():
                    sdk = self.sdk_factory()
                    error_type = sdk["CursorAgentError"]
                    # Local create and send both reject a non-empty idempotency key.
                    # The key is only used in the agent name.
                    options = {
                        "model": self.model,
                        "api_key": self.api_key,
                        "name": f"{role}-{idempotency_key}"[:80],
                        "local": _local(sdk, workdir),
                        "cloud": None,
                        "tools": _tools(role),
                    }
                    option_type = sdk.get("AgentOptions")
                    if option_type is None:
                        agent = sdk["Agent"].create(**options)
                    else:
                        agent = sdk["Agent"].create(option_type(**options))
                    run = agent.send(prompt)
                    agent_id = getattr(agent, "agent_id", "")
                    run_id = getattr(run, "run_id", "")
                    self._record(role, agent_id, run_id)
                    result = self._wait(run)
                status = getattr(result, "status", "error")
                text = getattr(result, "result", "") or ""
                total_tokens = _total_tokens(result)
                budget = TOKEN_BUDGET.get(role)
                if budget is not None and total_tokens > budget:
                    return RunOutcome("token_budget", text, agent_id, run_id, 2, True)
                if status != "finished":
                    return RunOutcome(status, text, agent_id, run_id, 2, True)
                return RunOutcome(status, text, agent_id, run_id, 0, True)
            except error_type as exc:
                last_error = exc
                retryable = bool(getattr(exc, "is_retryable", False))
                if not retryable or attempt >= self.max_retries:
                    return RunOutcome("startup_error", str(exc), "", "", 1, False)
                self.sleeper(_delay(attempt, exc))
            finally:
                if agent is not None:
                    agent.close()
        return RunOutcome("startup_error", str(last_error or "startup failed"), "", "", 1, False)

    def _wait(self, run: Any) -> Any:
        if not run.supports("wait") and hasattr(run, "supports"):
            raise RunnerError("run does not support wait")

        def _call() -> Any:
            return run.wait()

        with ThreadPoolExecutor(max_workers=1) as pool:
            future = pool.submit(_call)
            try:
                return future.result(timeout=self.timeout_seconds)
            except FuturesTimeout:
                if run.supports("cancel"):
                    run.cancel()
                return _Expired()

    def _record(self, role: str, agent_id: str, run_id: str) -> None:
        payload = {"role": role, "agent_id": agent_id, "run_id": run_id}
        print(json.dumps(payload), file=sys.stderr, flush=True)
        if self.run_record_path is not None:
            self.run_record_path.write_text(json.dumps(payload), encoding="utf-8")


class _Expired:
    status = "expired"
    result = ""


def _total_tokens(result: Any) -> int:
    usage = getattr(result, "usage", None)
    total = getattr(usage, "total_tokens", 0) or 0
    try:
        return int(total)
    except (TypeError, ValueError):
        return 0


def _local(sdk: dict[str, Any], cwd: str | None) -> Any:
    if not cwd:
        raise RunnerError("local runtime requires cwd")
    return sdk["LocalAgentOptions"](cwd=cwd, setting_sources=[])


def _tools(role: str) -> list[str]:
    if role in {"disposition", "review"}:
        return ["read", "grep", "glob", "ls"]
    return ["read", "grep", "glob", "ls", "edit"]


def _delay(attempt: int, exc: Exception) -> float:
    raw = getattr(exc, "retry_after", None)
    if raw is not None:
        try:
            return min(max(float(raw), 0), 30)
        except (TypeError, ValueError):
            pass
    return float(min(2**attempt, 30))


def _sleep(seconds: float) -> None:
    import time

    time.sleep(seconds)


@contextmanager
def hidden_secrets() -> Iterator[None]:
    saved = {key: os.environ.pop(key) for key in SECRET_ENV if key in os.environ}
    try:
        yield
    finally:
        os.environ.update(saved)
