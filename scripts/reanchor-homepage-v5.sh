#!/usr/bin/env bash
# Print human-operated steps to upgrade ArdmerePoRAnchor to schema v6 and anchor
# the current homepage report set with custodyRiskScore.
#
# Does not broadcast: requires PRIVATE_KEY for the permissioned signer on Sepolia.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

ANCHOR_SIG='anchorSnapshot(string,string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,bool,uint16,uint16,uint8)'

echo "== 1) UUPS upgrade (proxy address unchanged)"
echo "export PATH=\"\$HOME/.foundry/bin:\$PATH\""
echo "export ANCHOR_CONTRACT=\${ANCHOR_CONTRACT:-<proxy from docs/deployments.md>}"
echo "export RPC_URL=\${RPC_URL:-<Sepolia RPC>}"
echo "export PRIVATE_KEY=\${PRIVATE_KEY:-<upgrade owner key>}"
echo "cd contracts && forge script script/Upgrade.s.sol:Upgrade --rpc-url \"\$RPC_URL\" --broadcast -vvvv"
echo "cast call \"\$ANCHOR_CONTRACT\" \"SCHEMA_VERSION()(uint8)\" --rpc-url \"\$RPC_URL\"   # expect 6"
echo "cast call \"\$ANCHOR_CONTRACT\" \"STORAGE_VERSION()(uint8)\" --rpc-url \"\$RPC_URL\"  # expect 3"
echo

echo "== 2) Homepage snapshots (current index.html HOMEPAGE_REPORT_SET)"
echo "# Immutable records cannot be updated; anchor only snapshotIds not yet on-chain."
echo "# Prefer: go build -o por ./cmd/por && ./por anchor -exchange <id> -audit <auditId>"
echo "# That prints schema v6 cast send (composite snapshotId + auditId string)."
echo

HOMEPAGE=(
  "binance:PR04NOV26"
  "bitunix:20261103"
  "phemex:20261101"
  "htx:20261101"
  "okx:508756245"
  "ourbit:20261101"
)

for row in "${HOMEPAGE[@]}"; do
  exchange="${row%%:*}"
  audit="${row##*:}"
  assessment="docs/reports/${exchange}/${audit}-assessment.json"
  if [[ ! -f "$assessment" ]]; then
    echo "# SKIP missing $assessment"
    continue
  fi
  score="$(
    python3 skills/por-custody-risk/scripts/risk_score.py --no-anchor "$assessment" \
      | python3 -c 'import json,sys; d=json.load(sys.stdin); s=d.get("riskScore"); print(255 if s is None else int(s))'
  )"
  echo "# --- ${exchange} ${audit} custodyRiskScore=${score} (file-only; roots from assessment) ---"
  echo "# Run ./por anchor -exchange ${exchange} -audit ${audit} for full calldata, or paste roots from:"
  echo "#   ${assessment}"
  echo "# cast send \"\$ANCHOR_CONTRACT\" '${ANCHOR_SIG}' \\"
  echo "#   ... --rpc-url \"\$RPC_URL\" --private-key <ANCHOR_SIGNER_PK>"
done

echo
echo "Signer (Sepolia): 0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F"
