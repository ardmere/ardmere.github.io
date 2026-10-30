#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
# shellcheck source=scripts/_por-env.sh
source "$(dirname "$0")/_por-env.sh"
exec go run ./cmd/por "$@"
