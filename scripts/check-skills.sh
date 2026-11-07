#!/usr/bin/env bash
# Run the PoR skill tests: inclusion-proof scheme parity and the custody risk table.
#
#   ./scripts/check-skills.sh          # tests only, no network
#   ./scripts/check-skills.sh --live   # also read the Sepolia anchor for one snapshot
#
# This repository has no .github/workflows directory, so this script is the PR-time
# check. Wire it into CI wherever the rest of the checks run.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

fail() { printf 'check-skills: %s\n' "$*" >&2; exit 1; }

command -v python3 >/dev/null 2>&1 || fail "python3 is required"
PY="$(command -v python3)"
"$PY" - <<'PY' || fail "python3 is older than 3.8; the skills use only the standard library but need 3.8+"
import sys
sys.exit(0 if sys.version_info >= (3, 8) else 1)
PY

echo "== syntax"
"$PY" -m py_compile \
  skills/lib/por_anchor.py \
  skills/por-custody-risk/scripts/risk_score.py \
  skills/por-inclusion-proof/scripts/inclusion_check.py \
  tests/test_inclusion.py
for script in skills/por-custody-risk/scripts/por-check.sh scripts/check-skills.sh; do
  bash -n "$script" || fail "$script has a syntax error"
done

echo "== schemes"
SCHEMES="$("$PY" skills/por-inclusion-proof/scripts/inclusion_check.py --list-schemes)"
COUNT="$(printf '%s\n' "$SCHEMES" | grep -c . || true)"
[[ "$COUNT" -ge 1 ]] || fail "--list-schemes printed nothing"
for exchange in $SCHEMES; do
  [[ -d "internal/exchanges/$exchange/merkle" ]] \
    || fail "scheme $exchange has no Go package in internal/exchanges/$exchange/merkle"
done
printf '  %s schemes: %s\n' "$COUNT" "$(printf '%s' "$SCHEMES" | tr '\n' ' ')"

echo "== tests"
"$PY" -m unittest discover -s tests -v

if [[ "${1-}" == "--live" ]]; then
  echo "== live anchor read"
  # One anchored snapshot and one unanchored one, so both paths are exercised.
  "$PY" skills/por-custody-risk/scripts/risk_score.py \
    docs/reports/bitget/20261021-assessment.json >/dev/null \
    || fail "live anchor read of bitget/20261021 failed; pass an RPC via POR_RPC_URL"
  echo "  bitget 20261021: anchored"

  # A not-found answer is a conclusion of its own, so this must exit 0 while
  # reporting ANCHOR_MISSING rather than a reserve judgement.
  "$PY" skills/por-custody-risk/scripts/risk_score.py \
    docs/reports/okx/508756245-assessment.json >"${TMPDIR:-/tmp}/por-not-found.json" \
    || fail "okx/508756245 is expected to answer not-found with a conclusion, not an error"
  "$PY" - "${TMPDIR:-/tmp}/por-not-found.json" <<'PY'
import json, sys
doc = json.load(open(sys.argv[1]))
anchor = doc.get("anchor") or {}
if anchor.get("status") != "not_found":
    sys.exit(f"expected not_found for okx/508756245, got {anchor.get('status')!r}")
if not doc.get("notAReserveJudgement"):
    sys.exit("not_found must be marked notAReserveJudgement")
print("  okx 508756245: not anchored (expected)")
PY
  rm -f "${TMPDIR:-/tmp}/por-not-found.json"
fi

echo "check-skills: ok"
