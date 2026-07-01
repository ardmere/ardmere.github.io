#!/usr/bin/env bash
# Download HTX sample public-data.zip and archive via por fetch.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
ZIP="${1:-$ROOT/fixtures/htx/public-data.zip}"
mkdir -p "$(dirname "$ZIP")"
if [[ ! -f "$ZIP" ]]; then
  curl -fsSL -o "$ZIP" \
    "https://github.com/huobiapi/Tool-Go-MerkleVerify/releases/download/2.0.0/public-data.zip"
fi
exec go run "$ROOT/cmd/por" fetch htx -zk-bundle "$ZIP" "$@"
