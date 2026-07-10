#!/usr/bin/env bash
# Fetch, verify, and publish reports for the last N PoR snapshots per exchange.
# Requires RPC keys in ~/.zshenv for full on-chain coverage (see docs/ledger-rpc-runbook.md).
set -euo pipefail
cd "$(dirname "$0")/.."

COUNT="${1:-3}"
EXCHANGES="${2:-okx,binance,gateio,bybit,bitget,htx}"

exec ./scripts/por.sh batch \
	-count "$COUNT" \
	-exchanges "$EXCHANGES" \
	-artifacts ./artifacts \
	-reports ./docs/reports \
	-full-rpc-latest-only=true
