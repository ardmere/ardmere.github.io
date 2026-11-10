#!/usr/bin/env python3
"""Scan manifest assessments and anchor missing v6 composite snapshotIds.

For each docs/reports/*/*-assessment.json, when the composite snapshotId is absent
on Sepolia or live custody scoring disagrees with the on-chain custodyRiskScore,
invokes scripts/anchor-assessment-v6.py with periodSeq = getSnapshotCount(exchangeTag).

Default is dry-run (prints cast send). Pass --apply to broadcast (requires PRIVATE_KEY).
"""
from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
LIB = ROOT / "skills" / "lib"
if str(LIB) not in sys.path:
    sys.path.insert(0, str(LIB))

import por_anchor  # noqa: E402

RISK = ROOT / "skills/por-custody-risk/scripts/risk_score.py"
ANCHOR_V6 = ROOT / "scripts/anchor-assessment-v6.py"


def assessment_paths(root: Path) -> list[Path]:
    reports = root / "docs" / "reports"
    return sorted(reports.glob("*/*-assessment.json"))


def composite_exists(exchange: str, audit_id: str, contract: str, rpc: str) -> bool | None:
    """snapshotExists for the v6 composite id only (not legacy)."""
    sid = por_anchor.composite_snapshot_id(audit_id, exchange)
    addr = por_anchor.contract_address(contract)
    sel = por_anchor.selector("snapshotExists(bytes32)") + sid[2:]
    for url in por_anchor.rpc_urls(rpc):
        try:
            payload = por_anchor.rpc_call(url, addr, sel)
            return por_anchor.decode_bool(payload)
        except por_anchor.AnchorError:
            continue
    return None


def live_risk_doc(path: Path, contract: str | None, rpc: str | None) -> dict:
    cmd = [sys.executable, str(RISK), str(path)]
    if contract:
        cmd.extend(["--anchor-contract", contract])
    if rpc:
        cmd.extend(["--rpc-url", rpc])
    proc = subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True)
    if proc.returncode not in (0, 1) or not proc.stdout.strip():
        raise RuntimeError((proc.stderr or proc.stdout or "risk_score failed").strip()[:500])
    return json.loads(proc.stdout)


def already_good(path: Path, contract: str, rpc: str) -> tuple[bool, str]:
    assessment = json.loads(path.read_text())
    exchange = str(assessment.get("exchange") or "")
    audit_id = str((assessment.get("snapshot") or {}).get("id") or "")
    if not exchange or not audit_id:
        return True, "skip: incomplete assessment"

    exists = composite_exists(exchange, audit_id, contract, rpc)
    if exists is None:
        return False, "rpc unavailable (will retry when anchoring)"

    try:
        doc = live_risk_doc(path, contract, rpc)
    except RuntimeError as exc:
        return False, f"risk_score: {exc}"

    anchor = doc.get("anchor") or {}
    composite = por_anchor.composite_snapshot_id(audit_id, exchange)
    matches = anchor.get("custodyRiskScoreMatchesComputed")
    on_chain = anchor.get("custodyRiskScoreOnChain")

    if exists and anchor.get("status") == "found" and anchor.get("snapshotId") == composite:
        if matches is True:
            return True, "composite anchored; custody score matches"
        if matches is None and on_chain == por_anchor.CUSTODY_RISK_UNAVAILABLE:
            return True, "composite anchored; on-chain score unavailable (255)"

    if not exists:
        return False, "composite snapshotId missing on-chain"
    if anchor.get("status") != "found" or anchor.get("snapshotId") != composite:
        return False, "only legacy or foreign snapshotId on-chain"
    if matches is False and exists:
        return True, "skip: composite exists; custody score drift is immutable on-chain"
    if matches is False:
        return False, "custodyRiskScore mismatch vs legacy row"
    return False, "needs v6 composite anchor"


def get_snapshot_count(exchange: str, contract: str, rpc: str) -> int:
    tag = por_anchor.exchange_tag(exchange)
    addr = por_anchor.contract_address(contract)
    sel = por_anchor.selector("getSnapshotCount(bytes32)") + tag[2:]
    for url in por_anchor.rpc_urls(rpc):
        try:
            payload = por_anchor.rpc_call(url, addr, sel)
            return int(payload, 16)
        except por_anchor.AnchorError:
            continue
    if por_anchor.cast_available():
        out = subprocess.check_output(
            [
                "cast",
                "call",
                addr,
                "getSnapshotCount(bytes32)(uint256)",
                tag,
                "--rpc-url",
                rpc,
            ],
            text=True,
            cwd=ROOT,
        ).strip()
        return int(out, 0)
    raise RuntimeError("getSnapshotCount: no RPC endpoint answered")


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument(
        "--apply",
        action="store_true",
        help="broadcast anchorSnapshot txs (requires PRIVATE_KEY; default is dry-run)",
    )
    parser.add_argument("--limit", type=int, default=0, help="stop after N anchors (0 = no limit)")
    opts = parser.parse_args(argv)

    contract = os.environ.get("ANCHOR_CONTRACT") or por_anchor.DEFAULT_CONTRACT
    rpc = os.environ.get("RPC_URL") or os.environ.get("SEPOLIA_RPC_URL") or por_anchor.DEFAULT_RPC_URL

    if opts.apply and not os.environ.get("PRIVATE_KEY"):
        print("PRIVATE_KEY required for --apply", file=sys.stderr)
        return 1

    dry_args = [] if opts.apply else ["--dry-run"]
    anchored = 0
    skipped = 0
    errors = 0

    for path in assessment_paths(ROOT):
        rel = path.relative_to(ROOT)
        good, reason = already_good(path, contract, rpc)
        if good:
            print(f"SKIP {rel}: {reason}", file=sys.stderr)
            skipped += 1
            continue

        assessment = json.loads(path.read_text())
        exchange = str(assessment.get("exchange") or "")
        try:
            period = get_snapshot_count(exchange, contract, rpc)
        except (subprocess.CalledProcessError, ValueError, FileNotFoundError) as exc:
            print(f"ERROR {rel}: getSnapshotCount failed: {exc}", file=sys.stderr)
            errors += 1
            continue

        print(f"PLAN {rel}: {reason}; periodSeq={period}", file=sys.stderr)
        cmd = [
            sys.executable,
            str(ANCHOR_V6),
            *dry_args,
            "--period-seq",
            str(period),
            str(path),
        ]
        proc = subprocess.run(cmd, cwd=ROOT, text=True)
        if proc.returncode != 0:
            errors += 1
            continue
        anchored += 1
        if opts.limit and anchored >= opts.limit:
            break

    print(
        f"# done: planned_or_sent={anchored} skipped={skipped} errors={errors} "
        f"mode={'apply' if opts.apply else 'dry-run'}",
        file=sys.stderr,
    )
    return 1 if errors else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
