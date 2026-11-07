#!/usr/bin/env bash
# List exchanges whose PoR monitor window is open today (UTC).
# Config: config/por-monitor/schedule.yaml
# Docs:   docs/por-monitor-schedule.md
#
# Usage:
#   ./scripts/por-monitor-due.sh                 # due ids, one per line
#   ./scripts/por-monitor-due.sh --json          # JSON summary
#   ./scripts/por-monitor-due.sh --date 2026-11-01
#   ./scripts/por-monitor-due.sh --include-shortlist
#   ./scripts/por-monitor-due.sh --with-reasons
#
# Deep scans are NOT run here — caller / Cloud Agent should light-probe only the
# printed ids, then deep-run on new snapshot ids (SLA: ≤1 calendar day).

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG="${POR_MONITOR_SCHEDULE:-$ROOT/config/por-monitor/schedule.yaml}"
DATE_UTC=""
JSON=0
WITH_REASONS=0
INCLUDE_SHORTLIST=0

usage() {
  sed -n '2,16p' "$0" | sed 's/^# \{0,1\}//'
  exit "${1:-0}"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --date) DATE_UTC="${2:-}"; shift 2 ;;
    --json) JSON=1; shift ;;
    --with-reasons) WITH_REASONS=1; shift ;;
    --include-shortlist) INCLUDE_SHORTLIST=1; shift ;;
    -h|--help) usage 0 ;;
    *) echo "unknown arg: $1" >&2; usage 1 ;;
  esac
done

if [[ ! -f "$CONFIG" ]]; then
  echo "missing schedule config: $CONFIG" >&2
  exit 1
fi

# Prefer python3 + PyYAML when available; otherwise a tiny inline parser for this schema.
export POR_MONITOR_CONFIG="$CONFIG"
export POR_MONITOR_DATE="${DATE_UTC}"
export POR_MONITOR_JSON="$JSON"
export POR_MONITOR_REASONS="$WITH_REASONS"
export POR_MONITOR_SHORTLIST="$INCLUDE_SHORTLIST"
export POR_MONITOR_REPORTS="$ROOT/docs/reports"

python3 - <<'PY'
import json, os, re, sys
from datetime import date, datetime, timedelta
from pathlib import Path

config_path = Path(os.environ["POR_MONITOR_CONFIG"])
reports_root = Path(os.environ["POR_MONITOR_REPORTS"])
want_json = os.environ.get("POR_MONITOR_JSON") == "1"
with_reasons = os.environ.get("POR_MONITOR_REASONS") == "1"
include_shortlist = os.environ.get("POR_MONITOR_SHORTLIST") == "1"
date_s = os.environ.get("POR_MONITOR_DATE") or ""

try:
    import yaml  # type: ignore
except ImportError:
    yaml = None

def parse_simple_yaml(text: str) -> dict:
    """Minimal subset parser for schedule.yaml if PyYAML is absent."""
    # Fallback: extract exchange blocks with regex — good enough for this file.
    data = {"exchanges": [], "shortlist_preview": [], "defaults": {}, "timers": {}}
    # month_days lists
    blocks = re.split(r"\n  - id: ", text)
    header, *rest = blocks
    section = "exchanges"
    if "shortlist_preview:" in header:
        pass
    for raw in rest:
        body = "id: " + raw
        if body.strip().startswith("id: ") and "shortlist_preview" in text:
            # decide section by position
            pos = text.find(raw)
            sp = text.find("shortlist_preview:")
            section = "shortlist_preview" if sp != -1 and pos > sp else "exchanges"
        m_id = re.search(r"^id:\s*(\S+)", body, re.M)
        if not m_id:
            continue
        entry = {"id": m_id.group(1)}
        m_mode = re.search(r"^    mode:\s*(\S+)", body, re.M) or re.search(r"^  mode:\s*(\S+)", body, re.M)
        # modes are indented under list item
        m_mode = re.search(r"^\s+mode:\s*(\S+)", body, re.M)
        if m_mode:
            entry["mode"] = m_mode.group(1)
        m_days = re.search(r"month_days:\s*\[([^\]]+)\]", body)
        if m_days:
            entry["month_days"] = [int(x.strip()) for x in m_days.group(1).split(",") if x.strip()]
        m_gap = re.search(r"typical_gap_days:\s*(\d+)", body)
        if m_gap:
            entry["typical_gap_days"] = int(m_gap.group(1))
        m_rad = re.search(r"window_radius_days:\s*(\d+)", body)
        if m_rad:
            entry["window_radius_days"] = int(m_rad.group(1))
        m_fb = re.search(r"fallback_month_days:\s*\[([^\]]+)\]", body)
        if m_fb:
            entry["fallback_month_days"] = [int(x.strip()) for x in m_fb.group(1).split(",") if x.strip()]
        m_q = re.search(r"quarter_end_months:\s*\[([^\]]+)\]", body)
        if m_q:
            entry["quarter_end_months"] = [int(x.strip()) for x in m_q.group(1).split(",") if x.strip()]
        m_pri = re.search(r"priority:\s*(\S+)", body)
        if m_pri:
            entry["priority"] = m_pri.group(1)
        m_rat = re.search(r'rationale:\s*"([^"]*)"', body)
        if m_rat:
            entry["rationale"] = m_rat.group(1)
        data[section].append(entry)
    # defaults buffers
    m_bb = re.search(r"buffer_days_before:\s*(\d+)", text)
    m_ba = re.search(r"buffer_days_after:\s*(\d+)", text)
    if m_bb:
        data["defaults"]["buffer_days_before"] = int(m_bb.group(1))
    if m_ba:
        data["defaults"]["buffer_days_after"] = int(m_ba.group(1))
    return data

text = config_path.read_text(encoding="utf-8")
if yaml is not None:
    data = yaml.safe_load(text)
else:
    data = parse_simple_yaml(text)

if date_s:
    today = date.fromisoformat(date_s)
else:
    today = datetime.utcnow().date()

defaults = data.get("defaults") or {}
buf_before = int(defaults.get("buffer_days_before", 1))
buf_after = int(defaults.get("buffer_days_after", 2))

def latest_snapshot_date(ex_id: str):
    d = reports_root / ex_id
    if not d.is_dir():
        # gateio folder is gateio; some reports use aliases
        return None
    best = None
    for p in d.glob("*-assessment.json"):
        try:
            obj = json.loads(p.read_text(encoding="utf-8"))
        except Exception:
            continue
        snap = obj.get("snapshot") or {}
        st = snap.get("snapshotTime") or obj.get("snapshotTime") or ""
        if not st:
            continue
        try:
            dt = datetime.fromisoformat(st.replace("Z", "+00:00")).date()
        except Exception:
            continue
        if best is None or dt > best:
            best = dt
    return best

def in_month_days(days, d: date, before=0, after=0) -> bool:
    # Expand each target day with before/after across month boundaries.
    for day in days:
        try:
            target = date(d.year, d.month, day)
        except ValueError:
            continue
        for offset in range(-before, after + 1):
            if d == target + timedelta(days=offset):
                return True
        # also allow matching when today is day N and N is in expanded set relative to month_days list only
    # Simpler: today's day-of-month in expanded set including wrap from prev/next month edges
    expanded = set()
    for day in days:
        expanded.add(day)
        for o in range(1, before + 1):
            expanded.add((((day - 1 - o) % 31) + 1))  # rough; replaced below
    # Correct approach: check distance to nearest calendar occurrence this/prev/next month
    for month_delta in (-1, 0, 1):
        y = d.year
        m = d.month + month_delta
        while m < 1:
            m += 12
            y -= 1
        while m > 12:
            m -= 12
            y += 1
        for day in days:
            try:
                target = date(y, m, day)
            except ValueError:
                continue
            delta = (d - target).days
            if -before <= delta <= after:
                return True
    return False

def due_next_expected(entry, d: date):
    last = latest_snapshot_date(entry["id"])
    gap = int(entry.get("typical_gap_days", 30))
    radius = int(entry.get("window_radius_days", 4))
    if last is not None:
        expected = last + timedelta(days=gap)
        # Open around the predicted next publish (and a short late-publish tail).
        if -radius <= (d - expected).days <= radius + buf_after:
            kind = "next_expected" if (d - expected).days <= radius else "late_window"
            return True, f"{kind}={expected.isoformat()} last={last.isoformat()} gap={gap}d radius={radius}"
        # Same-day / day-after re-probe right after a known publish (merkle-in-progress etc.).
        if 0 <= (d - last).days <= buf_after:
            return True, f"post_publish_buffer last={last.isoformat()} +{buf_after}d"
        return False, f"not due (last={last.isoformat()} expected={expected.isoformat()})"
    # No assessed history yet: fall back to calendar days.
    fb = entry.get("fallback_month_days") or []
    if fb and in_month_days(fb, d, before=0, after=0):
        return True, f"fallback_month_days={fb} (no last snapshot)"
    return False, "not due (no last snapshot)"

def due_quarter_end(entry, d: date):
    months = entry.get("quarter_end_months") or [3, 6, 9, 12]
    days = entry.get("month_days") or [28, 29, 30, 31, 1, 2, 3]
    if d.month in months or (d.month in [1, 4, 7, 10] and d.day <= 5):
        # quarter-end month or first days of following month
        if d.month in months and in_month_days(days, d, before=0, after=0):
            return True, f"quarter_end month={d.month}"
        if d.month not in months and d.day <= 5:
            return True, f"quarter_end spillover day={d.day}"
    return False, "not quarter window"

def is_due(entry):
    mode = entry.get("mode") or "month_days"
    if mode == "month_days":
        days = entry.get("month_days") or []
        ok = in_month_days(days, today, before=buf_before, after=buf_after)
        return ok, f"month_days={days} buffer=-{buf_before}/+{buf_after}"
    if mode == "next_expected":
        return due_next_expected(entry, today)
    if mode == "quarter_end":
        return due_quarter_end(entry, today)
    return False, f"unknown mode {mode}"

entries = list(data.get("exchanges") or [])
if include_shortlist:
    for s in data.get("shortlist_preview") or []:
        s = dict(s)
        s["_shortlist"] = True
        entries.append(s)

due = []
for e in entries:
    ok, reason = is_due(e)
    if ok:
        item = {
            "id": e["id"],
            "priority": e.get("priority") or ("shortlist" if e.get("_shortlist") else ""),
            "mode": e.get("mode"),
            "reason": reason,
            "rationale": e.get("rationale") or "",
            "shortlist": bool(e.get("_shortlist")),
        }
        due.append(item)

if want_json:
    print(json.dumps({
        "date_utc": today.isoformat(),
        "config": str(config_path),
        "due_count": len(due),
        "due": due,
        "policy": "per-exchange windows; deep-scan only on new snapshot id; SLA <=1 calendar day",
    }, indent=2))
else:
    if not due:
        if with_reasons:
            print(f"# {today.isoformat()} UTC — none due", file=sys.stderr)
        # empty stdout = nothing to probe
        sys.exit(0)
    for item in due:
        if with_reasons:
            tag = item["priority"] or ("shortlist" if item["shortlist"] else "")
            print(f"{item['id']}\t{tag}\t{item['reason']}")
        else:
            print(item["id"])
PY
