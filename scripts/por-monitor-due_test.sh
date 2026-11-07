#!/usr/bin/env bash
# Regression tests for scripts/por-monitor-due.sh.
#
# Guards the properties that made ten adapted venues invisible to every timer run:
#   1. every exchange with committed reports has a schedule entry (--check)
#   2. newly added entries actually appear on their real publish dates
#   3. per-entry buffers, overdue catch-up and publisher lag behave as configured
#   4. no entry is silently unreachable (every entry is due on >=1 date in a year)
#   5. WAF/browser venues carry their blocking note rather than looking healthy
#
# Usage: ./scripts/por-monitor-due_test.sh        (exits non-zero on failure)

set -uo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DUE="$ROOT/scripts/por-monitor-due.sh"
FAIL=0

pass() { printf 'ok   %s\n' "$1"; }
fail() { printf 'FAIL %s\n' "$1"; FAIL=1; }

# due_on <date> -> ids due that day, one per line
due_on() { "$DUE" --date "$1"; }

expect_contains() { # expect_contains <date> <id> <label>
  if due_on "$1" | grep -qx "$2"; then
    pass "$3"
  else
    fail "$3 (expected $2 due on $1, got: $(due_on "$1" | tr '\n' ' '))"
  fi
}

expect_missing() { # expect_missing <date> <id> <label>
  if due_on "$1" | grep -qx "$2"; then
    fail "$3 (expected $2 NOT due on $1)"
  else
    pass "$3"
  fi
}

check_field() { # check_field <date> <id> <column> <field=value> <label>
  local row
  row="$("$DUE" --date "$1" --with-reasons | awk -F'\t' -v id="$2" -v c="$3" '$1==id {print $c}')"
  if [[ "$row" == *"$4"* ]]; then
    pass "$5"
  else
    fail "$5 (expected $4 in column $3, got: ${row:-<no row>})"
  fi
}

echo "== config validation"
if "$DUE" --check >/dev/null; then
  pass "--check reports every exchange with reports has a schedule entry"
else
  fail "--check reported coverage gaps: $("$DUE" --check)"
fi

echo "== newly added venues appear on their observed publish dates"
expect_contains 2026-11-01 ourbit    "ourbit due on month-start window"
expect_contains 2026-10-15 coinex    "coinex due inside next_expected window"
expect_contains 2026-11-03 bitunix   "bitunix due on its first-days window"
expect_contains 2026-10-15 woox      "woox due on day 15"
expect_contains 2026-10-01 toobit    "toobit due on day 1"
expect_contains 2026-11-01 phemex    "phemex due on day 1"
expect_contains 2026-11-20 coinw     "coinw due on its late-month window"
expect_contains 2026-10-15 btcc      "btcc due on day 15"
expect_contains 2026-11-07 kraken    "kraken due inside its publish-lag window"
expect_contains 2026-10-01 poloniex  "poloniex due on day 1"

echo "== these venues were provably absent before the fix (regression guard)"
# 2026-11-01 previously produced only binance, okx, kucoin, htx.
expect_missing 2026-11-01 mexc "2026-11-01 still excludes mexc (window not open)"

echo "== buffers are per-entry, not global"
# coinw uses -1d before / +3d after around a late-month (20–31) window.
expect_missing 2026-11-18 coinw "coinw not due before its -1d buffer edge"
expect_contains 2026-11-19 coinw "coinw -1d buffer opens day 19"
expect_contains 2026-12-01 coinw "coinw +3d buffer spills into the next month"
expect_missing 2026-12-04 coinw "coinw +3d buffer closes after day 3 of the next month"

echo "== overdue catch-up keeps a stalled venue visible"
expect_contains 2026-11-07 coinex "coinex still due after 54d of silence"
check_field 2026-11-07 coinex 3 "overdue_catchup" "coinex reason is overdue_catchup"
# Before the window opens it must still be quiet.
expect_missing 2026-09-20 coinex "coinex quiet outside window and before silence cap"

echo "== kraken quarterly lag window"
expect_contains 2026-10-15 kraken "kraken due 15d after Q3 quarter end"
expect_contains 2026-11-21 kraken "kraken due 52d after Q3 quarter end (late publish)"
expect_missing 2026-11-25 kraken "kraken stops once the lag window closes"
expect_missing 2026-10-05 kraken "kraken quiet before the minimum publish lag"

echo "== WAF / browser venues are marked, not silently healthy"
check_field 2026-11-20 coinw 4 "capture=browser" "coinw entry is browser capture"
check_field 2026-11-20 coinw 4 "waf=" "coinw entry carries a WAF note"
check_field 2026-11-01 ourbit 4 "confidence=low" "ourbit is marked low confidence (single snapshot)"
check_field 2026-10-15 btcc 4 "confidence=high" "btcc is marked high confidence"

echo "== shortlist flag still works and no longer double-lists adapted venues"
sl=$("$DUE" --date 2026-10-15 --include-shortlist)
if [[ "$(echo "$sl" | sort | uniq -d)" == "" ]]; then
  pass "--include-shortlist produces no duplicate ids"
else
  fail "--include-shortlist duplicates: $(echo "$sl" | sort | uniq -d | tr '\n' ' ')"
fi
if "$DUE" --date 2026-10-15 --include-shortlist --with-reasons | awk -F'\t' '$2=="shortlist"' | grep -q .; then
  fail "--include-shortlist still lists adapted venues as shortlist previews"
else
  pass "--include-shortlist has no stale shortlist_preview members"
fi

echo "== every entry is reachable at least once per year"
python3 - "$ROOT" <<'PY'
import subprocess, sys, datetime, yaml
root = sys.argv[1]
due = root + "/scripts/por-monitor-due.sh"
cfg = yaml.safe_load(open(root + "/config/por-monitor/schedule.yaml"))
ids = [e["id"] for e in cfg["exchanges"]]
seen = set()
d = datetime.date(2027, 1, 1)
while d < datetime.date(2028, 1, 1):
    out = subprocess.run([due, "--date", d.isoformat()], capture_output=True, text=True).stdout
    seen.update(x for x in out.split() if x)
    d += datetime.timedelta(days=1)
missing = [i for i in ids if i not in seen]
for i in ids:
    if i in missing:
        print(f"FAIL entry {i} never due in 2027")
    else:
        print(f"ok   entry {i} is due on at least one 2027 date")
sys.exit(1 if missing else 0)
PY
if [[ $? -ne 0 ]]; then FAIL=1; fi

echo
if [[ $FAIL -ne 0 ]]; then
  echo "FAILED"
  exit 1
fi
echo "all schedule checks passed"
