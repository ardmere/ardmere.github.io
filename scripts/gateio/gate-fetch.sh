#!/usr/bin/env bash
# Fetch Gate.com public PoR summary and save raw artifacts locally.
set -euo pipefail
cd "$(dirname "$0")/../.."
exec go run ./cmd/por fetch gateio "$@"
