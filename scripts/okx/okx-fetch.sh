#!/usr/bin/env bash
# Fetch OKX public PoR artifacts and save under artifacts/okx/<auditId>/.
set -euo pipefail
cd "$(dirname "$0")/../.."
exec go run ./cmd/por fetch okx "$@"
