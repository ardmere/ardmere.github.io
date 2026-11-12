"""Human-friendly display helpers for PoR reserve totals and major-coin ratios.

On-chain schema v7 stores whole USD as uint128 (`0` = unavailable; a positive
sub-dollar total that rounds to 0 is encoded as `1`) and major-coin reserve
ratios as uint16 bps (`10000` = 100%, `65535` = unavailable, `65534` = saturated
≥655.34%). Assessment JSON keeps float USD and ratio units where `1.0` = 100%.

These helpers turn those machine values into the lines skills should show
people. They do not judge solvency: every label is self-reported.
"""

from __future__ import annotations

USD_UNAVAILABLE = 0
RESERVE_RATIO_UNAVAILABLE = 65535
RESERVE_RATIO_BPS_MAX = 65534

SELF_REPORTED_NOTE = (
    "Self-reported exchange figures, not a solvency opinion or reserve adequacy claim."
)

MAJOR_COIN_BPS_FIELDS = (
    ("BTC", "btcReserveRatioBps"),
    ("ETH", "ethReserveRatioBps"),
    ("USDT", "usdtReserveRatioBps"),
    ("USDC", "usdcReserveRatioBps"),
)


def _is_number(value):
    return isinstance(value, (int, float)) and not isinstance(value, bool)


def format_usd_float(value):
    """Format an assessment float USD total with thousands separators."""
    if value is None:
        return "unavailable"
    if not _is_number(value):
        return "n/a"
    amount = float(value)
    if amount != amount:  # NaN
        return "unavailable"
    # Keep cents when the assessment published them; still use grouping.
    text = f"{amount:,.2f}"
    return f"${text} USD"


def format_usd_onchain(raw, assessment_float=None):
    """Format an on-chain uint128 whole-USD slot for people.

    Prefer the assessment float when the skill already has it (avoids showing
    the encode sentinel `1` as literal "$1" for a positive sub-dollar total).
    Without that float: `0` is unavailable; other integers are whole dollars,
    with an explicit whole-USD note when the stored value is `1`.
    """
    if assessment_float is not None and _is_number(assessment_float):
        return format_usd_float(assessment_float)
    if raw is None:
        return "unavailable"
    if not _is_number(raw):
        return "n/a"
    whole = int(raw)
    if whole == USD_UNAVAILABLE:
        return "unavailable"
    if whole == 1:
        return "$1 USD (on-chain storage is whole USD)"
    return f"${whole:,} USD (on-chain storage is whole USD)"


def format_reserve_ratio_bps(bps):
    """Format a uint16 reserve-ratio slot: 10000 → 100.00%, 65535 → unavailable."""
    if bps is None:
        return "unavailable"
    if not _is_number(bps):
        return "n/a"
    value = int(bps)
    if value == RESERVE_RATIO_UNAVAILABLE:
        return "unavailable"
    if value == RESERVE_RATIO_BPS_MAX:
        return "≥655.34%"
    if value < 0:
        return "n/a"
    return f"{value / 100:.2f}%"


def format_reserve_ratio_rate(rate):
    """Format an assessment reserveRate where 1.0 = 100%."""
    if rate is None:
        return "unavailable"
    if not _is_number(rate):
        return "n/a"
    value = float(rate)
    if value != value:  # NaN
        return "unavailable"
    # Match the on-chain saturate rule when the float itself exceeds the cap.
    if value * 10000 >= RESERVE_RATIO_BPS_MAX + 0.5:
        return "≥655.34%"
    return f"{value * 100:.2f}%"


def _assessment_rate_map(assessment):
    out = {}
    if not isinstance(assessment, dict):
        return out
    rows = assessment.get("majorCoinReserveRates") or []
    for row in rows:
        if not isinstance(row, dict):
            continue
        coin = str(row.get("coin") or "").strip().upper()
        if coin and _is_number(row.get("reserveRate")):
            out[coin] = float(row.get("reserveRate"))
    return out


def build_reserve_display(record=None, assessment=None, *, meaningful=None):
    """Build a display block from an on-chain record and/or assessment JSON.

    Returns raw machine values alongside friendly strings. `lines` is the
    primary human-facing form skills should quote or print.
    """
    record = record if isinstance(record, dict) else {}
    assessment = assessment if isinstance(assessment, dict) else {}

    if meaningful is None:
        if "reserveTotalsMeaningful" in record:
            meaningful = bool(record.get("reserveTotalsMeaningful"))
        elif record.get("schemaVersion") is not None:
            meaningful = int(record.get("schemaVersion") or 0) >= 7
        else:
            meaningful = None

    assessment_liab = assessment.get("totalLiabilityUsd")
    assessment_res = assessment.get("totalReserveUsd")
    onchain_liab = record.get("totalLiabilityUsd")
    onchain_res = record.get("totalReserveUsd")

    # Prefer assessment floats for the USD lines when present.
    if _is_number(assessment_liab):
        liability_display = format_usd_float(assessment_liab)
        liability_source = "assessment"
    elif meaningful is False:
        liability_display = "unavailable"
        liability_source = "onchain"
    else:
        liability_display = format_usd_onchain(onchain_liab)
        liability_source = "onchain" if onchain_liab is not None else "none"

    if _is_number(assessment_res):
        reserve_display = format_usd_float(assessment_res)
        reserve_source = "assessment"
    elif meaningful is False:
        reserve_display = "unavailable"
        reserve_source = "onchain"
    else:
        reserve_display = format_usd_onchain(onchain_res)
        reserve_source = "onchain" if onchain_res is not None else "none"

    rate_map = _assessment_rate_map(assessment)
    ratios = []
    for coin, field in MAJOR_COIN_BPS_FIELDS:
        raw_bps = record.get(field)
        if meaningful is False:
            # Pre-v7 storage reads 0 for these slots; that is not a 0% ratio.
            raw_bps = RESERVE_RATIO_UNAVAILABLE
        if coin in rate_map:
            display = format_reserve_ratio_rate(rate_map[coin])
            source = "assessment"
        elif meaningful is False or raw_bps is None:
            display = "unavailable"
            source = "onchain" if meaningful is False else "none"
        else:
            display = format_reserve_ratio_bps(raw_bps)
            source = "onchain"
        ratios.append({
            "coin": coin,
            "display": display,
            "rawBps": raw_bps,
            "assessmentRate": rate_map.get(coin),
            "source": source,
        })

    lines = [
        f"Total liabilities: {liability_display} (self-reported)",
        f"Total reserves: {reserve_display} (self-reported)",
    ]
    for row in ratios:
        lines.append(f"{row['coin']} reserve ratio: {row['display']} (self-reported)")
    lines.append(SELF_REPORTED_NOTE)

    return {
        "label": "self-reported",
        "note": SELF_REPORTED_NOTE,
        "reserveTotalsMeaningful": meaningful,
        "totalLiabilityUsd": {
            "display": liability_display,
            "rawOnchain": onchain_liab,
            "assessment": assessment_liab if _is_number(assessment_liab) else None,
            "source": liability_source,
        },
        "totalReserveUsd": {
            "display": reserve_display,
            "rawOnchain": onchain_res,
            "assessment": assessment_res if _is_number(assessment_res) else None,
            "source": reserve_source,
        },
        "majorCoinReserveRatios": ratios,
        "lines": lines,
    }


def print_reserve_lines(display, file=None):
    """Write the primary human-facing reserve lines (default: stderr)."""
    import sys

    handle = sys.stderr if file is None else file
    block = display if isinstance(display, dict) else {}
    for line in block.get("lines") or []:
        handle.write(f"{line}\n")


def attach_record_display(record, assessment=None):
    """Return a shallow copy of `record` with a `reserveDisplay` block added."""
    if not isinstance(record, dict):
        return record
    out = dict(record)
    out["reserveDisplay"] = build_reserve_display(record, assessment)
    return out
