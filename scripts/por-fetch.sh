#!/usr/bin/env bash
# Deprecated: use ./scripts/por.sh fetch
set -euo pipefail
cd "$(dirname "$0")/.."
exec go run ./cmd/por fetch "$@"
