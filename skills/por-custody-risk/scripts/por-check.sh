#!/usr/bin/env bash
# Score one assessment against the ArdmerePoRAnchor record, print JSON, decide nothing.
#
#   por-check.sh docs/reports/okx/506872725-assessment.json
#   POR_ANCHOR_CONTRACT=0x... POR_RPC_URL=https://... por-check.sh <assessment.json>
#   por-check.sh --anchor-fixture fixtures/anchor/sepolia-2026-11-07.json <assessment.json>
#
# The chain read comes from the live Sepolia proxy unless --anchor-fixture or
# --no-anchor is given. A recorded fixture is a replay: the output says live=false.
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
exec python3 "$HERE/risk_score.py" "$@"
