#!/usr/bin/env bash
# Fetch, verify, and publish reports for the last N PoR snapshots per exchange.
set -euo pipefail
cd "$(dirname "$0")/.."

COUNT="${1:-3}"
EXCHANGES="${2:-okx,binance,gateio,bybit,bitget,htx}"

exec go run ./cmd/por batch \
  -count "$COUNT" \
  -exchanges "$EXCHANGES" \
  -artifacts ./artifacts \
  -reports ./docs/reports \
  -full-rpc-latest-only=true
