#!/usr/bin/env bash
# Save latest Gate PoR public summary into ./artifacts/gateio/<auditId>/raw/
#
# Tries live API first; on Akamai block, imports bundled browser-capture fixtures.
set -euo pipefail
cd "$(dirname "$0")/../.."

FIXTURE_DIR="fixtures/gateio/20260316"

if go run ./cmd/por fetch gateio "$@"; then
  exit 0
fi

echo "Live API unavailable; importing fixtures from ${FIXTURE_DIR}..." >&2
exec go run ./cmd/por fetch gateio \
  -info-file "${FIXTURE_DIR}/getProofOfReservesInfo.json" \
  -coins-file "${FIXTURE_DIR}/getProofOfReservesCoinList.json" \
  "$@"
