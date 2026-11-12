#!/usr/bin/env python3
"""Batch-anchor missing composite rows with live schema-v7 anchorSnapshot.

Scans docs/reports/*/*-assessment.json. For each assessment whose composite
snapshotId is absent on Sepolia, encodes the full v7 payload (USD totals,
four major-coin bps, custodyRiskScore from risk_score.py --no-anchor) and
calls signer-only anchorSnapshot. Rows that already exist are skipped
(including SnapshotAlreadyExists on send). periodSeq starts at
getSnapshotCount(exchange) and increments locally after each successful send
for that exchange.

Default is dry-run. Pass --apply to broadcast (requires PRIVATE_KEY).
Never prints PRIVATE_KEY or raw RPC API keys in the log.
"""
from __future__ import annotations

import argparse
import base64
import binascii
import json
import math
import os
import re
import subprocess
import sys
import time
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
LIB = ROOT / "skills" / "lib"
if str(LIB) not in sys.path:
    sys.path.insert(0, str(LIB))

import por_anchor  # noqa: E402

ANCHOR_SIG = (
    "anchorSnapshot(string,string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,"
    "uint8,uint16,uint8,bool,uint16,uint16,uint8,uint128,uint128,uint16,uint16,uint16,uint16)"
)
ALREADY_EXISTS = "SnapshotAlreadyExists"
ALREADY_EXISTS_SEL = "0xb296876b"
COINS = ("BTC", "ETH", "USDT", "USDC")
STAGE_MAP = {
    "Pre-Stage": 0,
    "Stage 0": 1,
    "Stage 1": 2,
    "Stage 2": 3,
}


def assessment_paths(root: Path) -> list[Path]:
    return sorted((root / "docs" / "reports").glob("*/*-assessment.json"))


def redact_rpc(url: str) -> str:
    if not url:
        return url
    for marker in ("/v3/", "/v2/"):
        if marker in url:
            head, _sep, _rest = url.partition(marker)
            return head + marker + "<redacted>"
    return url


def on_chain_usd(value) -> int:
    if value is None:
        return 0
    try:
        f = float(value)
    except (TypeError, ValueError):
        return 0
    if math.isnan(f) or math.isinf(f) or f <= 0:
        return 0
    rounded = int(round(f))
    if rounded < 1:
        return 1
    return rounded


def on_chain_ratio_bps(rates: list, coin: str) -> int:
    for row in rates or []:
        if not isinstance(row, dict):
            continue
        if str(row.get("coin") or "").strip().upper() != coin:
            continue
        raw = row.get("reserveRate")
        if raw is None:
            return por_anchor.RESERVE_RATIO_UNAVAILABLE
        try:
            f = float(raw)
        except (TypeError, ValueError):
            return por_anchor.RESERVE_RATIO_UNAVAILABLE
        if math.isnan(f) or math.isinf(f) or f < 0:
            return por_anchor.RESERVE_RATIO_UNAVAILABLE
        bps = int(round(f * 10000))
        if bps > 65534:
            return 65534
        return max(0, bps)
    return por_anchor.RESERVE_RATIO_UNAVAILABLE


def encode_figures(assessment: dict) -> tuple[int, int, int, int, int, int]:
    rates = assessment.get("majorCoinReserveRates") or []
    return (
        on_chain_usd(assessment.get("totalLiabilityUsd")),
        on_chain_usd(assessment.get("totalReserveUsd")),
        *(on_chain_ratio_bps(rates, c) for c in COINS),
    )


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
        [
            sys.executable,
            str(ROOT / "skills/por-custody-risk/scripts/risk_score.py"),
            "--no-anchor",
            str(assessment_path),
        ],
        cwd=ROOT,
        capture_output=True,
        text=True,
        check=True,
    )
    doc = json.loads(proc.stdout)
    score = doc.get("riskScore")
    return 255 if score is None else int(score)


def merkle_b32(value: str) -> str:
    v = (value or "").strip()
    zero = "0x" + "00" * 32
    if not v or v == zero:
        return zero
    if v.startswith("0x") and len(v) == 66:
        return v
    if len(v) == 64 and all(c in "0123456789abcdefABCDEF" for c in v):
        return "0x" + v
    try:
        raw = base64.b64decode(v, validate=True)
    except (binascii.Error, ValueError):
        raw = v.encode("utf-8")
    if len(raw) > 32:
        raw = raw[:32]
    raw = raw.ljust(32, b"\x00")
    return "0x" + raw.hex()


def composite_exists(exchange: str, audit_id: str, contract: str, rpc: str) -> tuple[str, str]:
    """Return (snapshotId, status) where status is found|not_found|unavailable."""
    sid = por_anchor.composite_snapshot_id(audit_id, exchange)
    addr = por_anchor.contract_address(contract)
    for url in por_anchor.rpc_urls(rpc):
        try:
            record = por_anchor._fetch_snapshot_at_id(url, addr, sid, prefer_cast=False, timeout=25)
            if record is None:
                return sid, "not_found"
            return sid, "found"
        except por_anchor._Revert as exc:
            name = (exc.decoded or {}).get("name")
            if name == "SnapshotNotFound":
                return sid, "not_found"
            continue
        except por_anchor.AnchorError:
            continue
    return sid, "unavailable"


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


def parse_tx_hash(stdout: str, stderr: str) -> str | None:
    for blob in (stdout or "", stderr or ""):
        blob = blob.strip()
        if not blob:
            continue
        try:
            doc = json.loads(blob)
        except json.JSONDecodeError:
            doc = None
        if isinstance(doc, dict):
            for key in ("transactionHash", "txHash", "hash"):
                val = doc.get(key)
                if isinstance(val, str) and re.fullmatch(r"0x[0-9a-fA-F]{64}", val):
                    return val
        m = re.search(
            r'"(?:transactionHash|txHash)"\s*:\s*"(0x[0-9a-fA-F]{64})"',
            blob,
        )
        if m:
            return m.group(1)
        m = re.search(
            r"(?:transactionHash|txHash)\s+[\"']?(0x[0-9a-fA-F]{64})",
            blob,
            flags=re.IGNORECASE,
        )
        if m:
            return m.group(1)
    return None


def revert_selector(stderr: str) -> str | None:
    m = re.search(r"data:\s*\"?(0x[0-9a-fA-F]{8})", stderr or "")
    if m:
        return m.group(1).lower()
    if ALREADY_EXISTS in (stderr or ""):
        return ALREADY_EXISTS_SEL
    return None


def build_cast_args(
    *,
    contract: str,
    exchange: str,
    audit_id: str,
    snap_b32: str,
    period: int,
    unix_t: int,
    btc_h: int,
    exch_root: str,
    art_root: str,
    ver_root: str,
    verdict: int,
    coverage_bps: int,
    por_stage: int,
    effective: bool,
    custody: int,
    liab: int,
    res: int,
    btc: int,
    eth: int,
    usdt: int,
    usdc: int,
) -> list[str]:
    return [
        contract,
        ANCHOR_SIG,
        exchange,
        audit_id,
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
        "true" if effective else "false",
        "0",  # walletCoverageBps — not present in assessment JSON
        "0",  # reserveCoverageBps
        str(custody),
        str(liab),
        str(res),
        str(btc),
        str(eth),
        str(usdt),
        str(usdc),
    ]


def send_anchor(
    *,
    contract: str,
    rpc: str,
    pk: str,
    args: list[str],
    dry_run: bool,
) -> tuple[int, str | None, str]:
    cmd = ["cast", "send", *args, "--rpc-url", rpc]
    if dry_run:
        printed = list(cmd)
        for i, part in enumerate(printed):
            if i > 0 and printed[i - 1] == "--rpc-url":
                printed[i] = redact_rpc(part)
        return 0, None, " ".join(printed)

    cmd_run = cmd + ["--private-key", pk, "--json"]
    proc = subprocess.run(cmd_run, cwd=ROOT, text=True, capture_output=True)
    tx = parse_tx_hash(proc.stdout, proc.stderr)
    detail = (proc.stderr or proc.stdout or "").strip().replace(pk, "<redacted>")
    if proc.returncode != 0:
        return proc.returncode, tx, detail
    return 0, tx, detail


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument(
        "--apply",
        action="store_true",
        help="broadcast anchorSnapshot txs (requires PRIVATE_KEY; default dry-run)",
    )
    parser.add_argument("--limit", type=int, default=0, help="stop after N anchors (0 = no limit)")
    parser.add_argument(
        "--only-missing-from-log",
        type=Path,
        default=None,
        help="optional path to v7-reserve-backfill-log.json; only process missing_anchor rows",
    )
    parser.add_argument(
        "--log",
        type=Path,
        default=ROOT / "docs" / "v7-anchor-missing-log.json",
        help="write JSON summary/log",
    )
    opts = parser.parse_args(argv)

    contract = os.environ.get("ANCHOR_CONTRACT") or por_anchor.DEFAULT_CONTRACT
    rpc = (
        os.environ.get("RPC_URL")
        or os.environ.get("SEPOLIA_RPC_URL")
        or por_anchor.DEFAULT_RPC_URL
    )
    pk = os.environ.get("PRIVATE_KEY") or ""
    if pk and not pk.startswith("0x"):
        pk = "0x" + pk
    if opts.apply and not pk:
        print("PRIVATE_KEY required for --apply", file=sys.stderr)
        return 1

    allow: set[tuple[str, str]] | None = None
    if opts.only_missing_from_log:
        log_doc = json.loads(opts.only_missing_from_log.read_text())
        allow = {
            (str(r.get("exchange")), str(r.get("auditId")))
            for r in log_doc.get("rows") or []
            if r.get("status") == "missing_anchor"
        }

    counts = {
        "assessed": 0,
        "skipped_already_exists": 0,
        "anchored": 0,
        "failed": 0,
        "rpc_unavailable": 0,
        "filtered_out": 0,
    }
    rows: list[dict] = []
    failed_row: dict | None = None
    next_period: dict[str, int] = {}

    for path in assessment_paths(ROOT):
        rel = str(path.relative_to(ROOT))
        assessment = json.loads(path.read_text())
        exchange = str(assessment.get("exchange") or "")
        audit_id = str((assessment.get("snapshot") or {}).get("id") or "")
        counts["assessed"] += 1
        if not exchange or not audit_id:
            rows.append({"path": rel, "status": "skip_incomplete"})
            continue
        if allow is not None and (exchange, audit_id) not in allow:
            counts["filtered_out"] += 1
            continue

        sid, status = composite_exists(exchange, audit_id, contract, rpc)
        if status == "unavailable":
            counts["rpc_unavailable"] += 1
            counts["failed"] += 1
            entry = {
                "path": rel,
                "exchange": exchange,
                "auditId": audit_id,
                "snapshotId": sid,
                "status": "rpc_unavailable",
            }
            rows.append(entry)
            failed_row = entry
            print(f"ERROR {rel}: RPC unavailable", file=sys.stderr)
            break
        if status == "found":
            counts["skipped_already_exists"] += 1
            rows.append(
                {
                    "path": rel,
                    "exchange": exchange,
                    "auditId": audit_id,
                    "snapshotId": sid,
                    "status": "skipped_already_exists",
                }
            )
            print(f"SKIP {rel}: composite exists {sid}", file=sys.stderr)
            continue

        snap = assessment.get("snapshot") or {}
        arts = assessment.get("artifacts") or {}
        snap_time = datetime.fromisoformat(str(snap["snapshotTime"]).replace("Z", "+00:00"))
        verdict, coverage_bps = summarize_verifications(assessment.get("verifierResults"))
        try:
            custody = custody_score(path)
        except (subprocess.CalledProcessError, json.JSONDecodeError, ValueError) as exc:
            counts["failed"] += 1
            entry = {
                "path": rel,
                "exchange": exchange,
                "auditId": audit_id,
                "snapshotId": sid,
                "status": "failed",
                "error": f"custody_score: {exc}"[:500],
            }
            rows.append(entry)
            failed_row = entry
            print(f"FAIL {rel}: custody_score — stopping", file=sys.stderr)
            break

        if exchange not in next_period:
            next_period[exchange] = get_snapshot_count(exchange, contract, rpc)
        period = next_period[exchange]

        liab, res, btc, eth, usdt, usdc = encode_figures(assessment)
        por_stage = encode_por_stage(str(assessment.get("porStage") or ""))
        effective = bool(assessment.get("effectivePoR"))
        exch_root = merkle_b32(str(snap.get("exchangeMerkleRoot") or ""))
        art_root = merkle_b32(str(arts.get("artifactBundleRoot") or ""))
        ver_root = merkle_b32(str(arts.get("verificationBundleRoot") or ""))
        btc_h = int(snap.get("btcBlockHeight") or 0)
        unix_t = int(snap_time.timestamp())

        args = build_cast_args(
            contract=contract,
            exchange=exchange,
            audit_id=audit_id,
            snap_b32=sid,
            period=period,
            unix_t=unix_t,
            btc_h=btc_h,
            exch_root=exch_root,
            art_root=art_root,
            ver_root=ver_root,
            verdict=verdict,
            coverage_bps=coverage_bps,
            por_stage=por_stage,
            effective=effective,
            custody=custody,
            liab=liab,
            res=res,
            btc=btc,
            eth=eth,
            usdt=usdt,
            usdc=usdc,
        )
        entry = {
            "path": rel,
            "exchange": exchange,
            "auditId": audit_id,
            "snapshotId": sid,
            "periodSeq": period,
            "custodyRiskScore": custody,
            "totalLiabilityUsd": liab,
            "totalReserveUsd": res,
            "btcReserveRatioBps": btc,
            "ethReserveRatioBps": eth,
            "usdtReserveRatioBps": usdt,
            "usdcReserveRatioBps": usdc,
            "porStage": por_stage,
            "effectivePoR": effective,
            "status": "planned" if not opts.apply else "sending",
        }
        print(
            f"{'PLAN' if not opts.apply else 'SEND'} {rel}: periodSeq={period} "
            f"custody={custody} usd={liab}/{res} bps={btc}/{eth}/{usdt}/{usdc}",
            file=sys.stderr,
        )
        code, tx, detail = send_anchor(
            contract=contract, rpc=rpc, pk=pk, args=args, dry_run=not opts.apply
        )
        if not opts.apply:
            entry["status"] = "planned"
            entry["cast"] = detail
            rows.append(entry)
            counts["anchored"] += 1
            next_period[exchange] = period + 1
            if opts.limit and counts["anchored"] >= opts.limit:
                break
            continue

        if code != 0:
            sel = revert_selector(detail)
            if sel == ALREADY_EXISTS_SEL or ALREADY_EXISTS in detail:
                counts["skipped_already_exists"] += 1
                entry["status"] = "skipped_already_exists_revert"
                entry["tx"] = tx
                entry["revertSelector"] = sel
                rows.append(entry)
                print(f"SKIP {rel}: SnapshotAlreadyExists {sel}", file=sys.stderr)
                # Do not bump periodSeq — the count slot was not consumed by us.
                continue

            counts["failed"] += 1
            entry["status"] = "failed"
            entry["tx"] = tx
            entry["revertSelector"] = sel
            entry["error"] = detail[:800]
            rows.append(entry)
            failed_row = entry
            print(
                f"FAIL {rel}: selector={sel or 'unknown'} — stopping\n{detail[:500]}",
                file=sys.stderr,
            )
            break

        entry["status"] = "anchored"
        entry["tx"] = tx
        rows.append(entry)
        counts["anchored"] += 1
        next_period[exchange] = period + 1
        print(f"OK {rel}: {tx}", file=sys.stderr)
        time.sleep(1.5)
        if opts.limit and counts["anchored"] >= opts.limit:
            break

    summary = {
        "generatedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "mode": "apply" if opts.apply else "dry-run",
        "contract": contract,
        "rpc": redact_rpc(rpc),
        "counts": counts,
        "failedRow": failed_row,
        "rows": rows,
    }
    opts.log.parent.mkdir(parents=True, exist_ok=True)
    opts.log.write_text(json.dumps(summary, indent=2) + "\n")
    print(json.dumps({"counts": counts, "log": str(opts.log)}, indent=2))
    return 1 if counts["failed"] else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
