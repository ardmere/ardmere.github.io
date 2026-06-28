---
name: okx-por-verifier
description: Work on ardmere's OKX Proof-of-Reserves verifier. Use when the user mentions OKX PoR, proof-of-reserves-v2, zkSTARKValidator, sum_proof_data.json, liability zip, wallet zip, appState, address ownership, global-zk-proof, OKX_ZK_STARK_VALIDATOR, onchain FAIL, omnibus, USDT-TRC20, snapshot 506872725, or por fetch okx.
---

# OKX PoR Verifier

Shared verdict / onchain triage → [por-onchain-triage](../por-onchain-triage/SKILL.md).  
Adapter checklist / tiers → [exchange-por-integration](../exchange-por-integration/SKILL.md).

## Start Here

- `docs/okx-por-data-guide.md` — CDN URLs, CLI, verifier profile.
- `docs/exchange-tiers.md` — OKX is **Tier 2** (wallet + sig + liability zk).

## Core Mental Model

| Plane | Source | Contents |
|-------|--------|----------|
| Summary | `appState.auditRootInfo` | ratios, `exchangeReserveBalances`, `custodyReserveBalances`, merkle hash |
| Wallet zip | CDN `por_csv_*_V2.zip` | Signed CSV + `okx_por_eth_staking_*.csv` |
| Liability zk | CDN `por_<auditId>_proof_data.zip` | **`sum_proof_data.json` only** |

- Merkle hash = **user liability tree root**, not wallet-tree root.
- Public liability zip has **no** user inclusion proofs.
- `custodyReserveBalances` → **WARN** only (summary-only).

## Two OKX Repos

| Repo | Role |
|------|------|
| `okx/proof-of-reserves` (v1) | Go `common.Verify*` for **address-ownership@okx-1** |
| `okx/proof-of-reserves-v2` | Plonky2; **`OKX_ZK_STARK_VALIDATOR`** = `zk-por-cli verify-global` |

Local clone (PoR monorepo): `../../MVP/proof-of-reserves-v2` from this repo root (`ardmere.github.io/`).

## Workflow

```bash
go run ./cmd/por fetch okx
export OKX_ZK_STARK_VALIDATOR=/path/to/proof-of-reserves-v2/target/release/zk-por-cli
go run ./cmd/por verify -exchange okx -snapshot <auditId> \
  -artifacts ./artifacts/okx/<auditId>
# zk-only iteration:
go run ./cmd/por verify ... -skip-rpc
```

## global-zk-proof@okx-1

`internal/verifier/global_zk_okx.go`:

1. Structure check on `sum_proof_data.json`.
2. **Summary binding:** `merkleHash` ↔ zk `public_inputs[2:6]`; `round_num` ↔ snapshotId.
3. Crypto (optional): `verify-global --proof-path` + stdin `\n` (~3–4 min circuit rebuild).

Without validator: **PARTIAL** (0.75) if binding OK; with validator: **PASS**.

| Wrong | Right |
|-------|-------|
| No-args binary in cwd | `verify-global --proof-path` |
| v1 Python validator on v2 proof | Build `proof-of-reserves-v2` |

## Wallet CSV

1. Top: `coin,amount` → `exchangeReserveBalances`.
2. Address rows: `coin,Network,Height,address,amount,message,signature…` — message **`I am an OKX address`**.

- **internal-consistency:** top-section totals (+ ETH staking file); never sum address rows by coin.
- **address-ownership:** verify on column **`coin`**, not `Network`. ~320k rows ≈ 1.5 min.

## Onchain Verifiers

Configs: `config/exchanges/okx/onchain.json`, `ledger.json`.  
Token/ledger loaders use **`loadTokenSupportedFor("okx")`** / **`loadLedgerSupportedFor("okx")`** via `Exchange` field in verifiers.

- **Subsample:** OKX wallets >800 supported rows → top-40 per pair, max 800 (`onchain_sample.go`).
- **Network aliases:** OKX CSV uses `TRON`/`POLYGON`/`APTOS` — config maps via `networkAliases` in `onchain.json` / `ledger.json` (~96% row coverage).
- **Never** `"BTC|BTC": "BTC"` in `native` — panics duplicate key.
- **Low coverage (~0.15%)** = 320k denominator; **~96% rows now mapped** via aliases — use **summary** supported-row stats.

Active profile (`internal/exchanges/okx/adapter.go`):

| Verifier | Notes |
|----------|--------|
| onchain-balance-hot@2.1 | ETH\|ETH, BNB\|BSC native |
| onchain-balance-token@2.0 | OKX coin labels (`USDT-TRC20`, …) |
| onchain-balance-ledger@1.2 | BTC, DOGE, LTC, SOL, APT, … |

## OKX-Specific Onchain Triage

Use [por-onchain-triage](../por-onchain-triage/SKILL.md) first; then:

| Pattern | OKX note |
|---------|----------|
| USDC/USDT **61M claim, ~1–2 on-chain** | `stablecoinEthOmnibusMismatch` → **WARN** |
| ETH hot `staked: 0`, claim >> liquid | `likelyEthDepositGap` + `ethHotInternalCustodyLikely` → **WARN** |
| USDT-TRC20 mixed FAIL + WARN surplus | All `TRX`/`TRON` mismatches → **WARN** (`tronSnapshotMismatchNote`) |
| **80 TRX WARN** | Subsample **40×2** (`TRX@TRX` + `USDT-TRC20@TRX`); live vs snapshot + omnibus labels — **not FAIL** |
| `TRX\|TRX` rpc `zero length contract` | **Fixed:** native TRX uses `TronNativeBalance` before TRC20 path |
| BTC claim **2000.00149786**, UTXO 0 | Template addresses — **UNVERIFIABLE** |
| BTC partial (e.g. 75% UTXO) | Omnibus custody |
| PEOPLE@ETH `0 bytes` | No contract code at snapshot block |
| SOL / APT | Live vs snapshot — **WARN** |

```bash
go run ./cmd/por probe rpc -network ETH -height 25037041 -chainlist
go run ./cmd/por probe tron -holder <T...> -height 82471660
zsh -lc 'go run ./cmd/por verify -exchange okx -snapshot 506872725 -artifacts ./artifacts/okx/506872725'
```

## Snapshot 506872725 Baseline

| Verifier | Result |
|----------|--------|
| artifact-integrity / solvency / internal-consistency | PASS |
| address-ownership | PASS (320,917/320,917) |
| global-zk-proof | PASS with validator; merkle bind OK |
| onchain-hot | **378 pass / 6 warn / 0 fail** of **384** ETH rows |
| onchain-token | **407 pass / 96 warn / 44 fail** of **800** sampled (cap); TRX/TRON **80 warn / 0 fail**; XLAYER/FEVM live WARN |
| onchain-ledger | **166 pass / 177 warn / 76 fail** of **546** sampled; FIL + TONCOIN-NEW active |

No insufficiency signal in verifiable surface; FAILs dominated by omnibus accounting + protocol/RPC limits.

Zk public inputs: equity 33,938,289,499,439; debt 2,058,816,990,077; liability 31,879,472,509,362.

Artifacts: `artifacts/okx/506872725/`.

## Verdict Discipline

→ [por-onchain-triage](../por-onchain-triage/SKILL.md). Do not equate low **coverage** or omnibus FAIL with insolvency.

## Code Entry Points

| Area | Path |
|------|------|
| Adapter | `internal/exchanges/okx/adapter.go` |
| Fetch / CDN | `internal/exchanges/okx/okxapi/` |
| Wallet parse | `internal/walletzip/okx.go` |
| Address sigs | `internal/verifier/address_ownership_okx.go` |
| Global zk | `internal/verifier/global_zk_okx.go`, `okx_zk_bind.go` |
| Onchain config | `config/exchanges/okx/` |

## Current Next Steps

1. ~~Fix **TRX native** routing~~ — done.
2. ~~Eth2 deposit heuristic + ETH hot omnibus WARN~~ — done (`likelyEthDepositGap`, `ethHotInternalCustodyLikely`).
3. ~~Omnibus mega-stablecoin → WARN~~ — done (`stablecoinEthOmnibusMismatch`).
4. Re-run verify to refresh bundle (needs `ETHERSCAN_API_KEY` for deposit PASS rows).
5. Expand onchain/ledger config for long-tail pairs.
