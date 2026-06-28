#!/usr/bin/env bash
# Import browser-captured Gate PoR API JSON into artifacts/gateio/<auditId>/raw/
#
# Usage:
#   ./scripts/gateio/gate-import-browser.sh ./info.json [./coinList.json] [./list.json]
#
# Capture steps:
#   1. Open https://www.gate.com/proof-of-reserves
#   2. DevTools → Network → filter "proof-of-reserves"
#   3. Save response bodies for getProofOfReservesInfo (+ optional CoinList / List)
set -euo pipefail
cd "$(dirname "$0")/../.."

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <info.json> [coinList.json] [list.json]" >&2
  exit 1
fi

args=(-info-file "$1")
[[ $# -ge 2 && -n "${2:-}" ]] && args+=(-coins-file "$2")
[[ $# -ge 3 && -n "${3:-}" ]] && args+=(-list-file "$3")

exec go run ./cmd/por fetch gateio "${args[@]}"
