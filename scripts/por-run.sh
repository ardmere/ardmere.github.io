#!/usr/bin/env bash
# Deprecated: use ./scripts/por.sh
set -euo pipefail
cd "$(dirname "$0")/.."
exec go run ./cmd/por "$@"
