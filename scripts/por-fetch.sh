#!/usr/bin/env bash
# Deprecated: use ./scripts/por.sh fetch
set -euo pipefail
exec "$(dirname "$0")/por.sh" fetch "$@"
