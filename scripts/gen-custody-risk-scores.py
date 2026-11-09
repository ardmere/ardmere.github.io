#!/usr/bin/env python3
"""Batch-compute custody risk scores for all public *-assessment.json files."""
import json
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SCRIPT = ROOT / "skills" / "por-custody-risk" / "scripts" / "risk_score.py"
OUT = ROOT / "docs" / "reports" / "custody-risk-scores.json"


def score_file(path: Path) -> dict:
    proc = subprocess.run(
        [sys.executable, str(SCRIPT), "--no-anchor", str(path)],
        cwd=ROOT,
        capture_output=True,
        text=True,
        check=False,
    )
    if proc.returncode not in (0, 1):
        raise RuntimeError(f"{path}: {proc.stderr or proc.stdout}")
    doc = json.loads(proc.stdout)
    return {
        "exchange": doc.get("exchange"),
        "snapshot": doc.get("snapshot"),
        "riskScore": doc.get("riskScore"),
        "riskLevel": doc.get("riskLevel"),
        "conclusive": doc.get("conclusive"),
        "assessmentPath": str(path.relative_to(ROOT)).replace("\\", "/"),
        "assessmentSha256": doc.get("assessmentSha256"),
    }


def main() -> int:
    rows = []
    for path in sorted(ROOT.glob("docs/reports/*/*-assessment.json")):
        rows.append(score_file(path))
    doc = {
        "schema": "ardmere/custody-risk-scores@1",
        "generatedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "scoringMode": "file-only (--no-anchor)",
        "decisionTable": "skills/por-custody-risk/references/decision-table.md",
        "scores": rows,
    }
    OUT.write_text(json.dumps(doc, indent=2) + "\n")
    print(f"wrote {len(rows)} scores to {OUT.relative_to(ROOT)}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
