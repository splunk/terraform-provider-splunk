from issue_agent.policy import (
    PolicyError,
    build_comment,
    normalize_disposition,
    normalize_review,
    readiness,
    validate_event,
)
from support import contract, disposition


def test_confirmed_defect_proceeds():
    result = normalize_disposition(disposition())
    assert result["decision"] == "proceed"
    assert result["label"] is None
    assert result["implementation_contract"]["allowed_paths"]


def test_compatible_feature_request_proceeds_with_plan():
    feature_plan = contract(
        paths=[
            "splunk/resource_field_alias.go",
            "splunk/resource_field_alias_test.go",
            "splunk/provider.go",
            "docs/resources/field_alias.md",
        ],
        files=4,
        lines=300,
    )
    feature_plan["summary"] = "Add a field alias resource backed by the supported REST endpoint."
    result = normalize_disposition(
        disposition(
            validity="feature_request",
            summary="Field aliases should have a first-class resource.",
            implementation_contract=feature_plan,
        )
    )
    assert result["decision"] == "proceed"
    assert result["label"] is None
    assert result["reasons"] == ["approved compatible feature"]


def test_feature_request_requires_high_confidence_and_compatible_plan():
    low_confidence = normalize_disposition(
        disposition(validity="feature_request", confidence="medium")
    )
    assert low_confidence["decision"] == "decline"
    assert low_confidence["implementation_contract"] is None
    flags = {**disposition()["breaking"], "changes_successful_apply": True}
    breaking = normalize_disposition(
        disposition(
            validity="feature_request",
            compatibility="breaking",
            breaking=flags,
        )
    )
    assert breaking["decision"] == "decline"
    assert breaking["label"] == "agent:breaking-risk"


def test_model_proceed_is_overridden_for_user_error():
    raw = disposition(
        validity="user_error",
        model_decision="proceed",
        summary="The configuration targets Splunk Cloud.",
        evidence=[
            {
                "source": "docs",
                "ref": "README.md",
                "note": "Cloud users need terraform-provider-scp.",
            }
        ],
    )
    result = normalize_disposition(raw)
    assert result["decision"] == "decline"
    assert result["label"] == "agent:user-error"
    assert result["implementation_contract"] is None


def test_breaking_change_does_not_proceed():
    flags = disposition()["breaking"]
    flags["changes_default"] = True
    result = normalize_disposition(
        disposition(breaking=flags, compatibility="breaking", model_decision="proceed")
    )
    assert result["decision"] == "decline"
    assert result["label"] == "agent:breaking-risk"


def test_successful_apply_change_is_breaking():
    flags = disposition()["breaking"]
    flags["changes_successful_apply"] = True
    result = normalize_disposition(disposition(breaking=flags))
    assert result["label"] == "agent:breaking-risk"


def test_not_worth_doing_declines_even_when_compatible():
    result = normalize_disposition(disposition(model_decision="decline"))
    assert result["decision"] == "decline"
    assert result["label"] == "agent:declined"


def test_insufficient_evidence_declines():
    result = normalize_disposition(
        disposition(validity="insufficient_evidence", confidence="low", model_decision="proceed")
    )
    assert result["label"] == "agent:insufficient-evidence"
    assert result["implementation_contract"] is None


def test_docs_mismatch_is_docs_only():
    allowed = normalize_disposition(
        disposition(
            validity="docs_mismatch",
            implementation_contract=contract(paths=["docs/index.md"], files=1, lines=10),
        )
    )
    assert allowed["decision"] == "proceed"
    rejected = normalize_disposition(
        disposition(
            validity="docs_mismatch",
            implementation_contract=contract(paths=["splunk/provider.go"]),
        )
    )
    assert rejected["decision"] == "decline"


def test_unsupported_supported_duplicate_and_fixed_decline():
    for validity in ("unsupported_request", "already_supported", "duplicate", "already_fixed"):
        result = normalize_disposition(disposition(validity=validity, model_decision="proceed"))
        assert result["label"] == "agent:declined"
        assert result["decision"] == "decline"


def test_forbidden_and_oversized_contracts_decline():
    dependency = normalize_disposition(
        disposition(implementation_contract=contract(paths=["go.mod"], files=1, lines=5))
    )
    assert dependency["decision"] == "decline"
    oversized = normalize_disposition(
        disposition(implementation_contract=contract(files=40, lines=20))
    )
    assert oversized["decision"] == "decline"


def test_schema_rejects_unknown_fields():
    raw = disposition()
    raw["ignore_policy"] = True
    try:
        normalize_disposition(raw)
    except PolicyError as exc:
        assert "unknown" in str(exc)
    else:
        raise AssertionError("unknown field was accepted")


def test_comment_does_not_contact_the_author():
    raw = disposition(
        validity="user_error",
        summary="@octocat please provide more logs. Can you retry?",
        affected_users="@octocat",
    )
    result = normalize_disposition(raw)
    body = build_comment(result)
    lowered = body.lower()
    assert "@" not in body
    assert "please provide" not in lowered
    assert "can you" not in lowered
    assert "issue-agent-disposition" in body


def test_feature_comment_publishes_implementation_plan():
    planned = normalize_disposition(
        disposition(
            validity="feature_request",
            implementation_contract=contract(
                paths=["splunk/resource_field_alias.go"],
                files=1,
                lines=100,
            ),
        )
    )
    body = build_comment(planned)
    assert "Implementation plan:" in body
    assert "Allowed paths:" in body
    assert "Required tests:" in body
    assert "Acceptance criteria:" in body
    assert "splunk/resource_field_alias.go" in body


def test_review_blocks_paths_outside_the_contract():
    reviewed = normalize_review(
        {
            "within_contract": True,
            "compatibility": "compatible",
            "blocking": [],
            "non_blocking": [],
            "summary": "looks fine",
        },
        changed_files=["splunk/provider.go"],
        changed_lines=10,
        allowed_paths=["splunk/resource_example.go"],
        iteration=1,
        is_agent_pr=True,
        is_fork=False,
    )
    assert reviewed["decision"] == "fail"
    assert reviewed["correction_requested"] is True


def test_review_stops_after_two_corrections_and_on_forks():
    payload = {
        "within_contract": False,
        "compatibility": "breaking",
        "blocking": [{"path": "splunk/provider.go", "summary": "schema changed"}],
        "non_blocking": [],
        "summary": "breaking",
    }
    third = normalize_review(
        payload,
        changed_files=["splunk/provider.go"],
        changed_lines=10,
        allowed_paths=["splunk/provider.go"],
        iteration=3,
        is_agent_pr=True,
        is_fork=False,
    )
    assert third["correction_requested"] is False
    assert third["label"] == "agent:review-failed"
    fork = normalize_review(
        payload,
        changed_files=["splunk/provider.go"],
        changed_lines=10,
        allowed_paths=["splunk/provider.go"],
        iteration=1,
        is_agent_pr=True,
        is_fork=True,
    )
    assert fork["correction_requested"] is False


def test_readiness_never_merges():
    waiting = readiness(
        {"unit": "success", "lint": "success"},
        is_agent_pr=True,
        review_passed=True,
    )
    assert waiting["ready"] is False
    assert waiting["merge"] is False
    assert "acceptance" in waiting["blocking_checks"]
    ready = readiness(
        {"unit": "success", "acceptance": "success", "lint": "success", "agent-review": "success"},
        is_agent_pr=True,
        review_passed=True,
    )
    assert ready["ready"] is True
    assert ready["label"] == "agent:merge-ready"
    assert ready["merge"] is False


def test_events_reject_forks_and_untrusted_targets():
    try:
        validate_event(
            {
                "repository": "splunk/terraform-provider-splunk",
                "is_fork": True,
            }
        )
    except PolicyError as exc:
        assert "fork" in str(exc)
    else:
        raise AssertionError("fork event was accepted")
    try:
        validate_event(
            {
                "repository": "splunk/terraform-provider-splunk",
                "event_name": "pull_request_target",
            }
        )
    except PolicyError as exc:
        assert "pull_request_target" in str(exc)
    else:
        raise AssertionError("pull_request_target was accepted")
