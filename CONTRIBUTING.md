# Contributing

Pull requests target `master`. A maintainer reviews and merges them. The issue agent can prepare a draft pull request, and it cannot approve or merge.

## Development

- Format with `make fmt`.
- On macOS, run unit tests with `make test`.
- On Linux, run unit tests with `make test-ci`. That command matches GitHub Actions and does not pass the macOS external linker flag.
- Acceptance tests need a Splunk Enterprise instance. `make testacc` reads `SPLUNK_HOME`, `SPLUNK_URL`, `SPLUNK_USERNAME`, and `SPLUNK_PASSWORD` from the environment. Do not commit those values.
- New resources need CRUD and import coverage in `resource_*_test.go` and a page under `docs/resources/`.
- Maintainers update `CHANGELOG.md` and push the `v*` tag after merge. Do not push release tags from a pull request.

## Issue agent

Run the agent for one open issue from the Actions tab. Choose **agent-issue**, enter the issue number, and choose a mode:

1. `review-only` posts the disposition and does not change code or open a pull request.
2. `review-and-fix` posts the disposition. For a `proceed` result, a local Cursor agent edits the checkout. Workflow code then checks the diff, runs `gofmt` and `go test`, runs a separate local review, and opens one draft pull request.

The Cursor agent runs on the GitHub Actions runner. It does not use the Cursor GitHub app, and it does not receive `GITHUB_TOKEN`. After the agent exits, the workflow commits branch `agent/issue-<number>` and opens the draft with the workflow token. One run handles one issue. A second run for the same issue waits until the first finishes.

Each issue gets a disposition before any code is written:

| Result | Meaning |
| --- | --- |
| `confirmed_defect` | A compatible fix can proceed |
| `docs_mismatch` | Only `docs/` may change |
| `user_error` | The reported result matches documented behavior |
| `unsupported_request`, `already_fixed`, `duplicate` | No change |
| `insufficient_evidence` | Not enough evidence to change code safely |
| `agent:breaking-risk` | A fix would change schema, defaults, identity, import format, or a configuration that already applies |

`proceed` requires a confirmed, high-confidence defect, or a docs-only correction, and a fix that preserves configurations that already apply. Tightening validation or changing a default is a breaking change.

`CURSOR_API_KEY` is a repository secret. Do not put it in source, prompts, logs, or comments. Fork pull requests cannot start the implementation agent.

The draft records the unit and format commands that ran in the agent workflow. A pull request opened with the workflow token does not start the repository test or lint workflows. Before merge, run **test** and **golangci-lint** from the Actions tab and select branch `agent/issue-<number>`. The agent cannot approve or merge.

One manual run sends at most the latest 12,000 characters of the issue body and 10 comments of 2,000 characters each. The model calls also stop the rest of the run when reported usage passes 80,000 tokens for disposition or review, or 200,000 tokens for an implementation attempt. Set a spend limit on the Cursor API key as the hard billing control. The SDK has no setting that stops a run before those tokens are consumed.

Labels record the outcome. They do not grant permission to start work:

- `agent:running`
- `agent:user-error`
- `agent:declined`
- `agent:breaking-risk`
- `agent:insufficient-evidence`
- `agent:review-failed`

Merging stays manual.

## Protected branch

Apply `.github/rulesets/master.json` only after `unit`, `acceptance`, and `lint` have succeeded on `master`. The ruleset requires one human approval, resolved conversations, dismissed stale reviews, and those three checks. It does not require an agent workflow. It has no bypass actors. Code-owner review stays off until `.github/CODEOWNERS` names a real maintainer team.

Remove older branch-protection bypasses when you apply the ruleset. The agent identity is not an allowed bypass.

```bash
# From master, after the checks above have run at least once:
APPLY_MASTER_RULESET=1 python tools/issue_agent/scripts/apply_ruleset.py
```

The script refuses to run without `APPLY_MASTER_RULESET=1`.
