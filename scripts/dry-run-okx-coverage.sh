#!/usr/bin/env bash
# Dry-run schema-v4 wallet/reserve coverage ratios for OKX audit 508756245.
# Uses the public assessment (ownership PASS, on-chain skip-rpc UNVERIFIABLE) plus
# a live OKX summary fetch for claimed asset amounts × CoinGecko at snapshotTime.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
AUDIT_ID="${1:-508756245}"
go run ./cmd/okxcoverage -audit "$AUDIT_ID"
