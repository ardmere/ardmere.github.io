---
name: por-onchain-triage
description: Triage onchain-balance-hot, token, and ledger FAIL/WARN/RPC errors across Binance, OKX, and other exchange PoR verifiers. Use when debugging chain vs CSV mismatches, surplus rows, archive RPC, Tron historical state, omnibus addresses, eth_getCode zero bytes, or low coverage percentages.
---

# PoR Onchain Failure Triage

Shared playbook for **all exchanges**. Exchange-specific baselines and configs:

- [binance-por-verifier](../binance-por-verifier/SKILL.md) — PR01JUN26, Stake Hub, BSC archive
- [okx-por-verifier](../okx-por-verifier/SKILL.md) — 506872725, omnibus CSV, 320k subsample
- [exchange-por-integration](../exchange-por-integration/SKILL.md) — verdict rules, adapter checklist

## Core Rule

**CSV `balance` is an accounted allocation, not necessarily liquid EOA/token balance at that address.**

Before interpreting FAIL as reserve shortfall, classify the row (see §Decision tree).

## Read Results Correctly

- **`coverage`** = `(pass + fail) / total wallet rows` — mostly reflects **unconfigured coin|network pairs**, not pass rate.
- **`summary` finding** = supported-row stats — **use this for pass rate** (e.g. `378 pass / 5 fail — of 384 supported`).
- OKX subsamples up to **800** top-balance rows per verifier when exchange is `okx` (`internal/verifier/onchain_sample.go`).

## Decision Tree

Classify each failing row **before** changing verifier code or verdict semantics:

| Symptom | Likely cause | Prefer |
|---------|--------------|--------|
| `short balanceOf response: 0 bytes` | Wrong token contract or **no code at snapshot block** | Check `eth_getCode`; **WARN/UNVERIFIABLE** if mapping wrong |
| Chain **>** CSV (surplus) | Staking rewards, Tron latest state, internal accounting conservative | **WARN** — not insolvency |
| Chain **<<** CSV, liquid-only, large gap | Omnibus / internal ledger / cold elsewhere | **WARN** with note; not fraud without other evidence |
| Chain **<<** CSV, ETH native, `staked: 0` | Missing Eth2 deposit indexer | **WARN** until `ETHERSCAN_API_KEY` + deposit scan |
| Small mismatch at exact `H` | Snapshot height boundary | Sample **H-1, H, H+1** (EVM) or wider BSC window |
| Tron TRC20 / native mismatch | **No historical state** on public nodes; any `TRX`/`TRON` row → **WARN** (`tronSnapshotMismatchNote`) |
| BTC UTXO **0** with round claim (e.g. 2000 BTC) | Template/placeholder address | **UNVERIFIABLE** (`btcNativeCustodyUnverifiable`) |
| BTC partial UTXO (e.g. 75% of claim) | Omnibus; not all UTXOs at single address | **WARN** or document gap |
| Solana / Aptos mismatch | Public RPC returns **latest** not snapshot slot/version | **WARN** (`ledgerSnapshotNote`) |
| DOGE/LTC tx count > API limit | BlockCypher/Esplora pagination cap | **UNVERIFIABLE** with limit note |
| RPC 429 / timeout / archive missing | Infrastructure | **WARN**, not FAIL |

**Notes must state what was observed:** `chain observed > CSV`, `H-1 matches CSV`, `archive unavailable`, `omnibus allocation`, `staking not indexed`.

Never use broad insolvency/fraud language from RPC gaps alone.

## Verdict Discipline (shared)

| Verdict | When |
|---------|------|
| **PASS** | Complete observation within tolerance |
| **FAIL** | Complete observation, material mismatch, no known attribution gap |
| **WARN** | Surplus, incomplete staking/indexer, archive/RPC limits, Tron live state |
| **UNVERIFIABLE** | No verifier, no contract at block, API tx limit, omnibus UTXO |
| **PARTIAL** | Verifier surface intentionally incomplete (e.g. zk binary missing) |

Full integration context → [exchange-por-integration](../exchange-por-integration/SKILL.md) §Verifier Profile.

## Probe Commands

```bash
# EVM archive at snapshot height
go run ./cmd/por probe rpc -network ETH -height <H> -chainlist
go run ./cmd/por probe rpc -network BSC -height <H> -chainlist

# Tron TRC20 (historical block often ineffective)
go run ./cmd/por probe tron -holder <base58> -height <H>

# Binance BNB staking
go run ./cmd/por probe stakehub -address <addr> -height <H>

# ETH deposit indexer needs login shell if key in profile
zsh -lc 'go run ./cmd/por verify -exchange okx -snapshot <id> -artifacts ./artifacts/okx/<id>'
```

Manual checks:

```bash
# Token exists at block?
curl -s -X POST https://eth.drpc.org -H 'content-type: application/json' \
  -d '{"jsonrpc":"2.0","method":"eth_getCode","params":["<token>","0x<hex_height>"],"id":1}'

# Native balance at block
curl -s -X POST https://eth.drpc.org -H 'content-type: application/json' \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["<addr>","0x<hex_height>"],"id":1}'
```

## Tron (all exchanges)

- Providers: `https://api.trongrid.io`, `https://tron-rpc.publicnode.com`
- Do **not** use `https://trx.publicnode.com` (404)
- java-tron: **no true archive** for TRC20 — treat mismatches as protocol limitation first
- Binance: `USDT|TRX`; OKX: `USDT-TRC20|TRX`, `TRX|TRX` (native TRX needs dedicated path, not empty-contract TRC20)

## Config Layout

Per exchange:

- `config/exchanges/<id>/onchain.json` — EVM native + ERC20/BEP20/TRC20 (`loadTokenSupportedFor(exchange)`)
- `config/exchanges/<id>/ledger.json` — UTXO, XRPL, Solana, Aptos (`loadLedgerSupportedFor(exchange)`)

OKX wallet rows use exchange-specific coin labels (`USDT-TRC20`, `ETH-ARBITRUM`); map keys must match **CSV coin|network** after `normalizeOKXNetwork`.

## When to Escalate vs Downgrade

**Downgrade to WARN** when: surplus, Tron, Solana live, missing Eth2 indexer, height boundary after H±1 check, omnibus pattern on mega stablecoin rows.

**Keep FAIL** when: complete ERC20/native observation at correct block, mapping verified, no attribution gap, material shortfall — rare on OKX address-level audit.

**Fix code** when: wrong routing (e.g. native TRX hitting TRC20), wrong contract address, exchange config not passed to token loader.
