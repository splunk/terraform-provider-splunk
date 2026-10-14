"""Shared disposition fixtures."""

from __future__ import annotations

import json
from types import SimpleNamespace

from issue_agent.runner import RunOutcome

BREAKING_FLAGS = {
    "changes_attribute_type": False,
    "changes_required_optional_computed": False,
    "changes_default": False,
    "changes_force_new": False,
    "tightens_validation": False,
    "changes_resource_identity": False,
    "changes_import_format": False,
    "requires_state_migration": False,
    "changes_update_vs_recreate": False,
    "changes_request_payload_or_splunk_defaults": False,
    "changes_older_splunk_behavior": False,
    "changes_successful_apply": False,
}


def contract(paths=None, files=2, lines=40):
    return {
        "summary": "Keep existing applies working and fix the reported read.",
        "allowed_paths": paths or ["splunk/resource_example.go", "splunk/resource_example_test.go"],
        "expected_files": files,
        "expected_lines": lines,
        "tests_required": ["existing apply", "reported failure"],
        "acceptance_criteria": ["plans that already apply stay unchanged"],
    }


def disposition(**overrides):
    raw = {
        "validity": "confirmed_defect",
        "confidence": "high",
        "compatibility": "compatible",
        "breaking": dict(BREAKING_FLAGS),
        "evidence": [
            {
                "source": "code",
                "ref": "splunk/resource_example.go",
                "note": "The read path drops a numeric value.",
            }
        ],
        "affected_users": "Users whose saved searches already apply are unchanged.",
        "summary": "The read path drops a numeric field that Splunk returns.",
        "model_decision": "proceed",
        "implementation_contract": contract(),
    }
    raw.update(overrides)
    return raw


def outcome(payload, *, status="finished", exit_code=0, started=True):
    text = payload if isinstance(payload, str) else json.dumps(payload)
    return RunOutcome(status, text, "agent-1", "run-1", exit_code, started)


class ScriptedRunner:
    def __init__(self, results):
        self.results = list(results)
        self.calls = []

    def complete(self, **kwargs):
        self.calls.append(kwargs)
        if not self.results:
            raise AssertionError("unexpected agent call")
        return self.results.pop(0)


def review_payload(**overrides):
    raw = {
        "within_contract": True,
        "compatibility": "compatible",
        "blocking": [],
        "non_blocking": [],
        "summary": "The diff stays inside the contract.",
    }
    raw.update(overrides)
    return raw


def finished_result(text="{}"):
    return SimpleNamespace(status="finished", result=text, id="run-1", agent_id="agent-1")
