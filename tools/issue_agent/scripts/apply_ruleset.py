"""Apply the master ruleset. Refuses unless a maintainer opts in after the checks exist."""

from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

REPOSITORY = "splunk/terraform-provider-splunk"


def main() -> int:
    if os.environ.get("APPLY_MASTER_RULESET") != "1":
        print(
            "Refusing to change branch protection. "
            "Set APPLY_MASTER_RULESET=1 only after this change is on master "
            "and the required checks have run.",
            file=sys.stderr,
        )
        return 2
    ruleset = Path(__file__).resolve().parents[3] / ".github" / "rulesets" / "master.json"
    if not ruleset.is_file():
        print(f"Missing ruleset file {ruleset}", file=sys.stderr)
        return 1
    command = [
        "gh",
        "api",
        "--method",
        "POST",
        f"repos/{REPOSITORY}/rulesets",
        "--input",
        str(ruleset),
    ]
    completed = subprocess.run(command, check=False)
    return completed.returncode


if __name__ == "__main__":
    raise SystemExit(main())
