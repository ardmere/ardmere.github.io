#!/usr/bin/env python3
"""Batch-backfill schema-v7 reserve figures onto existing Sepolia composite rows.

For each docs/reports/*/*-assessment.json:
  - resolve composite snapshotId = expectedSnapshotId(exchange, auditId)
  - if no on-chain composite row: missing-anchor (do not invent an anchor)
  - if schemaVersion >= 7: already-done
  - else: encode USD/bps from assessment JSON and call owner-only
    backfillReserveFigures(...). Wait for success before the next send.

Default is dry-run. Pass --apply to broadcast (requires PRIVATE_KEY).
Never prints PRIVATE_KEY.
"""
from __future__ import annotations

import argparse
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

BACKFILL_SIG = (
    "backfillReserveFigures(string,string,uint128,uint128,uint16,uint16,uint16,uint16)"
)
ALREADY_APPLIED = "ReserveBackfillAlreadyApplied"
ALREADY_APPLIED_SEL = "0xa77dc6db"
SCHEMA_TARGET = 7
COINS = ("BTC", "ETH", "USDT", "USDC")


def assessment_paths(root: Path) -> list[Path]:
    return sorted((root / "docs" / "reports").glob("*/*-assessment.json"))


def on_chain_usd(value) -> int:
    """uint128 whole USD. null/absent/non-positive → 0; positive that rounds to 0 → 1."""
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
    """round(rate*10000); missing → 65535; saturate at 65534."""
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


def encode_figures(assessment: dict) -> tuple[int, int, int, int, int, int, bool]:
    liab = on_chain_usd(assessment.get("totalLiabilityUsd"))
    res = on_chain_usd(assessment.get("totalReserveUsd"))
    rates = assessment.get("majorCoinReserveRates") or []
    btc, eth, usdt, usdc = (on_chain_ratio_bps(rates, c) for c in COINS)
    null_usd = assessment.get("totalLiabilityUsd") is None and assessment.get("totalReserveUsd") is None
    return liab, res, btc, eth, usdt, usdc, null_usd


def composite_record(exchange: str, audit_id: str, contract: str, rpc: str) -> tuple[str, dict | None, str]:
    """Return (snapshotId, record|None, status) for the composite id only."""
    sid = por_anchor.composite_snapshot_id(audit_id, exchange)
    addr = por_anchor.contract_address(contract)
    for url in por_anchor.rpc_urls(rpc):
        try:
            record = por_anchor._fetch_snapshot_at_id(url, addr, sid, prefer_cast=False, timeout=25)
            if record is None:
                return sid, None, "not_found"
            return sid, record, "found"
        except por_anchor._Revert as exc:
            name = (exc.decoded or {}).get("name")
            if name == "SnapshotNotFound":
                return sid, None, "not_found"
            continue
        except por_anchor.AnchorError:
            continue
    return sid, None, "unavailable"


def redact_rpc(url: str) -> str:
    """Strip API-key path segments from Infura/Alchemy-style RPC URLs for logs."""
    if not url:
        return url
    for marker in ("/v3/", "/v2/"):
        if marker in url:
            head, _sep, _rest = url.partition(marker)
            return head + marker + "<redacted>"
    return url


def parse_tx_hash(stdout: str, stderr: str) -> str | None:
    """Prefer cast --json transactionHash; never grab log topic0 by accident."""
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
    if ALREADY_APPLIED in (stderr or ""):
        return ALREADY_APPLIED_SEL
    return None


def send_backfill(
    *,
    contract: str,
    rpc: str,
    pk: str,
    exchange: str,
    audit_id: str,
    liab: int,
    res: int,
    btc: int,
    eth: int,
    usdt: int,
    usdc: int,
    dry_run: bool,
) -> tuple[int, str | None, str]:
    cmd = [
        "cast",
        "send",
        contract,
        BACKFILL_SIG,
        exchange,
        audit_id,
        str(liab),
        str(res),
        str(btc),
        str(eth),
        str(usdt),
        str(usdc),
        "--rpc-url",
        rpc,
    ]
    if dry_run:
        # Omit private key; redact API keys in --rpc-url for any persisted log.
        printed = list(cmd)
        for i, part in enumerate(printed):
            if i > 0 and printed[i - 1] == "--rpc-url":
                printed[i] = redact_rpc(part)
        return 0, None, " ".join(printed)

    cmd_run = cmd + ["--private-key", pk, "--json"]
    proc = subprocess.run(cmd_run, cwd=ROOT, text=True, capture_output=True)
    tx = parse_tx_hash(proc.stdout, proc.stderr)
    detail = (proc.stderr or proc.stdout or "").strip()
    detail = detail.replace(pk, "<redacted>")
    if proc.returncode != 0:
        return proc.returncode, tx, detail
    return 0, tx, detail


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument(
        "--apply",
        action="store_true",
        help="broadcast backfillReserveFigures txs (requires PRIVATE_KEY; default dry-run)",
    )
    parser.add_argument("--limit", type=int, default=0, help="stop after N sends (0 = no limit)")
    parser.add_argument(
        "--log",
        type=Path,
        default=ROOT / "docs" / "v7-reserve-backfill-log.json",
        help="write JSON summary/log (default docs/v7-reserve-backfill-log.json)",
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

    counts = {
        "assessed": 0,
        "already_done": 0,
        "backfilled": 0,
        "missing_anchor": 0,
        "null_figures_but_written": 0,
        "failed": 0,
        "skipped_already_applied_revert": 0,
        "rpc_unavailable": 0,
    }
    rows: list[dict] = []
    failed_row: dict | None = None

    for path in assessment_paths(ROOT):
        rel = str(path.relative_to(ROOT))
        assessment = json.loads(path.read_text())
        exchange = str(assessment.get("exchange") or "")
        audit_id = str((assessment.get("snapshot") or {}).get("id") or "")
        counts["assessed"] += 1
        if not exchange or not audit_id:
            entry = {"path": rel, "status": "skip_incomplete"}
            rows.append(entry)
            print(f"SKIP {rel}: incomplete assessment", file=sys.stderr)
            continue

        sid, record, status = composite_record(exchange, audit_id, contract, rpc)
        if status == "unavailable":
            counts["rpc_unavailable"] += 1
            entry = {
                "path": rel,
                "exchange": exchange,
                "auditId": audit_id,
                "snapshotId": sid,
                "status": "rpc_unavailable",
            }
            rows.append(entry)
            print(f"ERROR {rel}: RPC unavailable for {sid}", file=sys.stderr)
            failed_row = entry
            counts["failed"] += 1
            break

        if status != "found" or record is None:
            counts["missing_anchor"] += 1
            entry = {
                "path": rel,
                "exchange": exchange,
                "auditId": audit_id,
                "snapshotId": sid,
                "status": "missing_anchor",
            }
            rows.append(entry)
            print(f"MISSING {rel}: no composite row {sid}", file=sys.stderr)
            continue

        schema = int(record.get("schemaVersion") or 0)
        if schema >= SCHEMA_TARGET:
            counts["already_done"] += 1
            entry = {
                "path": rel,
                "exchange": exchange,
                "auditId": audit_id,
                "snapshotId": sid,
                "schemaVersion": schema,
                "status": "already_done",
            }
            rows.append(entry)
            print(f"DONE {rel}: schemaVersion={schema}", file=sys.stderr)
            continue

        liab, res, btc, eth, usdt, usdc, null_usd = encode_figures(assessment)
        entry = {
            "path": rel,
            "exchange": exchange,
            "auditId": audit_id,
            "snapshotId": sid,
            "schemaVersionBefore": schema,
            "totalLiabilityUsd": liab,
            "totalReserveUsd": res,
            "btcReserveRatioBps": btc,
            "ethReserveRatioBps": eth,
            "usdtReserveRatioBps": usdt,
            "usdcReserveRatioBps": usdc,
            "nullUsdInAssessment": null_usd,
            "status": "planned" if not opts.apply else "sending",
        }

        print(
            f"{'PLAN' if not opts.apply else 'SEND'} {rel}: schema={schema} "
            f"usd={liab}/{res} bps={btc}/{eth}/{usdt}/{usdc}",
            file=sys.stderr,
        )
        code, tx, detail = send_backfill(
            contract=contract,
            rpc=rpc,
            pk=pk,
            exchange=exchange,
            audit_id=audit_id,
            liab=liab,
            res=res,
            btc=btc,
            eth=eth,
            usdt=usdt,
            usdc=usdc,
            dry_run=not opts.apply,
        )
        if not opts.apply:
            entry["status"] = "planned"
            entry["cast"] = detail
            rows.append(entry)
            counts["backfilled"] += 1  # planned count under dry-run "would backfill"
            if null_usd:
                counts["null_figures_but_written"] += 1
            if opts.limit and counts["backfilled"] >= opts.limit:
                break
            continue

        if code != 0:
            sel = revert_selector(detail)
            if sel == ALREADY_APPLIED_SEL or ALREADY_APPLIED in detail:
                counts["skipped_already_applied_revert"] += 1
                counts["already_done"] += 1
                entry["status"] = "already_done_revert"
                entry["tx"] = tx
                entry["revertSelector"] = sel
                rows.append(entry)
                print(f"SKIP {rel}: already schema-7 revert {sel}", file=sys.stderr)
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

        entry["status"] = "backfilled"
        entry["tx"] = tx
        rows.append(entry)
        counts["backfilled"] += 1
        if null_usd:
            counts["null_figures_but_written"] += 1
        print(f"OK {rel}: {tx}", file=sys.stderr)
        # Brief pause so the next eth_call sees the mined state.
        time.sleep(1.5)
        if opts.limit and counts["backfilled"] >= opts.limit:
            break

    summary = {
        "generatedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "mode": "apply" if opts.apply else "dry-run",
        "contract": contract,
        "schemaTarget": SCHEMA_TARGET,
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
