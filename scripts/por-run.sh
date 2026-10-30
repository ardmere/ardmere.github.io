#!/usr/bin/env bash
# Deprecated: use ./scripts/por.sh
set -euo pipefail
exec "$(dirname "$0")/por.sh" "$@"
