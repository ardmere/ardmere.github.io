#!/usr/bin/env python3
"""Broadcast anchorSnapshot from docs/reports assessment JSON (schema v5 roots + custody score)."""
import base64
import binascii
import json
import os
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent

STAGE_MAP = {
    "Pre-Stage": 0,
    "Stage 0": 1,
    "Stage 1": 2,
    "Stage 2": 3,
}


def encode_por_stage(stage: str) -> int:
    return STAGE_MAP.get((stage or "").strip(), 0)


def summarize_verifications(verifier_results) -> tuple[int, int]:
    summary = 0
    max_cov = 0.0
    for v in verifier_results or []:
        if not isinstance(v, dict):
            continue
        cov = float(v.get("coverage") or 0)
        if cov > max_cov:
            max_cov = cov
        vid = str(v.get("verifierId") or "")
        verdict = str(v.get("verdict") or "")
        if verdict != "PASS":
            continue
        if vid == "internal-consistency":
            summary |= 1 << 0
        elif vid == "onchain-balance-hot":
            summary |= 1 << 1
        elif vid == "solvency-claim":
            summary |= 1 << 2
    if summary & (1 << 0):
        summary |= 1 << 3
    return summary, int(max_cov * 10000)


def custody_score(assessment_path: Path) -> int:
    proc = subprocess.run(
        [sys.executable, str(ROOT / "skills/por-custody-risk/scripts/risk_score.py"), "--no-anchor", str(assessment_path)],
        cwd=ROOT,
        capture_output=True,
        text=True,
        check=True,
    )
    doc = json.loads(proc.stdout)
    score = doc.get("riskScore")
    return 255 if score is None else int(score)


def snapshot_id_b32(audit_id: str) -> str:
    out = subprocess.check_output(["cast", "keccak", audit_id], text=True).strip()
    return out if out.startswith("0x") else "0x" + out


def merkle_b32(value: str) -> str:
    v = (value or "").strip()
    if not v or v == "0x0000000000000000000000000000000000000000000000000000000000000000":
        return "0x0000000000000000000000000000000000000000000000000000000000000000"
    if v.startswith("0x") and len(v) == 66:
        return v
    if len(v) == 64 and all(c in "0123456789abcdefABCDEF" for c in v):
        return "0x" + v
    raw = None
    try:
        raw = base64.b64decode(v, validate=True)
    except (binascii.Error, ValueError):
        raw = v.encode("utf-8")
    if len(raw) > 32:
        raw = raw[:32]
    raw = raw.ljust(32, b"\x00")
    return "0x" + raw.hex()


def main(argv: list[str]) -> int:
    period_override = None
    args = argv[1:]
    if "--period-seq" in args:
        i = args.index("--period-seq")
        period_override = int(args[i + 1])
        del args[i : i + 2]
    if len(args) < 1:
        print("usage: anchor-assessment-v5.py [--period-seq N] docs/reports/<ex>/<id>-assessment.json", file=sys.stderr)
        return 2
    assessment_path = Path(args[0])
    if not assessment_path.is_file():
        print(f"missing {assessment_path}", file=sys.stderr)
        return 1
    assessment = json.loads(assessment_path.read_text())
    exchange = assessment["exchange"]
    snap = assessment["snapshot"]
    arts = assessment.get("artifacts") or {}
    audit_id = str(snap["id"])
    snap_time = datetime.fromisoformat(snap["snapshotTime"].replace("Z", "+00:00"))
    verdict, coverage_bps = summarize_verifications(assessment.get("verifierResults"))
    custody = custody_score(assessment_path)
    por_stage = encode_por_stage(assessment.get("porStage"))
    effective = "true" if assessment.get("effectivePoR") else "false"
    contract = os.environ.get("ANCHOR_CONTRACT", "0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7")
    rpc = os.environ.get("RPC_URL") or os.environ.get("SEPOLIA_RPC_URL") or "https://ethereum-sepolia-rpc.publicnode.com"
    pk = os.environ.get("PRIVATE_KEY")
    if not pk:
        print("PRIVATE_KEY required", file=sys.stderr)
        return 1
    if not pk.startswith("0x"):
        pk = "0x" + pk

    snap_b32 = snapshot_id_b32(audit_id)
    exch_root = merkle_b32(snap.get("exchangeMerkleRoot") or "")
    art_root = merkle_b32(arts.get("artifactBundleRoot") or "")
    ver_root = merkle_b32(arts.get("verificationBundleRoot") or "")
    btc_h = int(snap.get("btcBlockHeight") or 0)
    period = period_override if period_override is not None else int(snap.get("periodSeq") or 0)
    unix_t = int(snap_time.timestamp())

    sig = (
        "anchorSnapshot(string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,"
        "uint8,uint16,uint8,bool,uint16,uint16,uint8)"
    )
    cmd = [
        "cast",
        "send",
        contract,
        sig,
        exchange,
        snap_b32,
        str(period),
        str(unix_t),
        str(btc_h),
        exch_root,
        art_root,
        ver_root,
        str(verdict),
        str(coverage_bps),
        str(por_stage),
        effective,
        "0",
        "0",
        str(custody),
        "--rpc-url",
        rpc,
        "--private-key",
        pk,
    ]
    print(f"# {exchange} {audit_id} custodyRiskScore={custody}", file=sys.stderr)
    proc = subprocess.run(cmd, cwd=ROOT, text=True, capture_output=True)
    sys.stdout.write(proc.stdout)
    if proc.returncode != 0:
        sys.stderr.write(proc.stderr)
        return proc.returncode
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
