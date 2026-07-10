# Ledger & RPC Runbook

> Operator guide for `onchain-balance-*` verifiers: configuration, environment variables, common failure modes, and fix order.  
> Companion: [`verifier-architecture.md`](./verifier-architecture.md) §5–§7, [`por-cli.md`](./por-cli.md).

Last updated: 2026-07-07 (PR01JUN26 regression)

---

## 1. Three on-chain verifiers

HotCold CSV rows are routed by `(coin, network)` into **one** verifier:

| Verifier | Config | Chains / assets |
|---|---|---|
| `onchain-balance-hot` | `config/exchanges/<ex>/onchain.json` → `native` | EVM native (ETH, BNB, …) |
| `onchain-balance-token` | `onchain.json` → `tokens[]` | ERC20 / BEP20 / TRC20 / Starknet ETH, … |
| `onchain-balance-ledger` | `config/exchanges/<ex>/ledger.json` | UTXO, XRPL, Solana, Substrate, NEAR, Aptos, Sui, HBAR, TON, … |

**Rule:** A pair listed in `ledger.json` is **not** an unsupported pair for hot/token — do not duplicate it in `onchain.json` unless intentional.

Findings use `status`: `PASS`, `WARN`, `FAIL`, `UNVERIFIABLE`.  
`ledger rpc error` in `note` means the query failed before a balance comparison — fix RPC/config first.

---

## 2. Configuration files

| File | Purpose |
|---|---|
| `config/rpc-providers.json` | EVM JSON-RPC pool (`ETH`, `BSC`, `ARBITRUM`, …). Supports `${ENV_VAR}` expansion via `os.ExpandEnv`; entries with unset vars are skipped. |
| `config/exchanges/binance/onchain.json` | Native + token contract addresses per `(coin, network)`. |
| `config/exchanges/binance/ledger.json` | Non-EVM backends: `esplora`, `blockchair`, `substrate`, `near-ft`, `ton-jetton`, … |

Override paths for tests or local runs:

- `RPC_PROVIDERS_CONFIG`
- `LEDGER_CONFIG`

---

## 3. Environment variables

Load before `por verify` (e.g. `source ~/.zshenv`):

| Variable | Used by | Notes |
|---|---|---|
| `INFURA_KEY` | EVM (`rpc-providers.json` ETH) | Archive `eth_call` / `balanceOf` at snapshot block |
| `ALCHEMY_KEY` | EVM ETH + ledger UTXO (`alchemy` backend) | ETH: archive `eth_call`. BTC/BCH: `GET .../api/v2/address/{addr}?to={height}` — indexed balance at block height, no Esplora pagination. |
| `BLOCKCHAIR_API_KEY` | Ledger UTXO (`blockchair` backend / fallback) | **Historical** `?state=<height>` in one request. Optional if `ALCHEMY_KEY` is set for BTC/BCH. |
| `HELIUS_API_KEY` | Solana ledger / token | Historical slot when configured |
| `NEAR_RPC` | NEAR native / NEP-141 | Default: `https://archival-rpc.mainnet.fastnear.com` |
| `NEAR_RPC_FALLBACK` | NEAR | Comma-separated extra endpoints |
| `CHROMIA_NODE` | CHR (`chromia` backend) | Postchain query base URL (port `7740`) |
| `HBAR_MIRROR` | HBAR / HTS | Default: `mainnet-public.mirrornode.hedera.com` |
| `TONAPI_URL` | TON native / jetton | Default: `https://tonapi.io/v2` |
| `XRPL_RPC` | XRP ledger rows | Default: `https://xrplcluster.com/` |
| `SOLANA_RPC` | SOL / SPL | Overrides `rpc-providers.json` SOL entries |

**Security:** Never commit keys; verify bundle notes may echo provider URLs (including expanded Infura paths).

---

## 4. Recommended fix order

When `por verify` shows many on-chain WARNs, work in this order (empirically fastest signal):

| Priority | Layer | Typical symptom | Action |
|---|---|---|---|
| P0 | Wrong contract / mint in config | `short eth_call`, zero bytecode | Fix `onchain.json` / `ledger.json` addresses (BSC peg typos, USDT@ARB, native USDC, …) |
| P1 | Wrong address encoding | TRX base58, placeholder `Bridging in progress` | Fix `onchain.json`; invalid addresses → `UNVERIFIABLE` not RPC error |
| P2 | RPC provider pool | HTTP 403/429, all providers exhausted | Edit `rpc-providers.json`; demote bad blastapi/drpc; add Infura/Alchemy for ETH |
| P3 | Client semantics | Empty `0x` = zero balance; batch `len<32` | `internal/rpc/evm.go` — normalize historical `eth_call` results |
| P4 | Per-chain ledger | `ledger rpc error` | This document §5–§6 |

Re-run after each layer:

```bash
source ~/.zshenv   # if keys live there
cd /path/to/ardmere.github.io
go run ./cmd/por verify -exchange binance -snapshot PR01JUN26 \
  -artifacts ./artifacts/binance/PR01JUN26
go run ./cmd/por report -exchange binance -snapshot PR01JUN26
```

Results cache under RPC cache dir — delete stale cache if a provider was fixed but old errors persist.

---

## 5. EVM lessons (token / hot)

### 5.1 Empty `eth_call` is valid zero

Archive nodes often return `0x` for `balanceOf` when balance is zero. Treat as 32 zero bytes, not `len<32` failure — otherwise Infura/Alchemy are skipped and traffic falls through to broken public nodes.

### 5.2 Env-expanded providers

`https://mainnet.infura.io/v3/${INFURA_KEY}` in `rpc-providers.json` only works when `INFURA_KEY` is set in the process environment. Confirm keys with a one-off `eth_blockNumber` / historical `eth_getBalance` at the snapshot block.

### 5.3 `short eth_call` after provider fix

If contracts are correct and archive keys work, `short eth_call` WARNs on ETH should drop to ~0. Remaining FAILs are usually real balance mismatch or omnibus semantics (mega-stablecoin rows).

---

## 6. Ledger lessons by chain

### 6.1 BTC (Esplora → Alchemy)

- **Symptom:** `esplora tx pagination exceeded N pages` (e.g. `34xp4vRo…`, `bc1q…`).
- **Cause:** Summing UTXOs by paging `/address/…/txs/chain` does not scale for exchange hot wallets.
- **Fix:** Set `ALCHEMY_KEY`; ledger entry has `"alchemy": "bitcoin"`. On Esplora failure, verifier falls back to Alchemy `GET .../api/v2/address/{addr}?to={height}`. Optional: `BLOCKCHAIR_API_KEY` as further fallback.
- **Code:** Pagination limit raised (80 → 400); still insufficient for megawallets without indexed provider.

### 6.2 BCH (Alchemy)

- **Symptom:** `blockchair HTTP 430` (IP rate limit / blacklist).
- **Fix:** Primary backend is `alchemy` with `"alchemy": "bitcoin-cash"` and `ALCHEMY_KEY`. Blockchair remains optional fallback when key is set.

### 6.3 ZEC

- **Symptom:** Blockcypher `HTTP 404` — ZEC mainnet not supported on Blockcypher v1.
- **Fix:** `ledger.json` backend `blockchair` / `zcash`. On Blockchair 430, fallback **CipherScan** (`cipherscan.app`) live balance + `liveSnapshot: true`.

### 6.4 NEAR (native + NEP-141)

- **Symptom:** `garbage collected` / deprecated `archival-rpc.mainnet.near.org`.
- **Fix:** Endpoints: FastNear archival, Lava, then live `rpc.mainnet.near.org`. Detect GC with `strings.Contains(err, "garbage collected")` (not `GARBAGE_COLLECTED` — underscore mismatch broke live fallback).
- **USDC mint:** Native Circle USDC account id is `17208628f84f5d6ad33f0da3bbbeb27ffcb398eac501a31bd6ad2011e36133a1` (not truncated hex). Bridged USDC.e: `a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48.factory.bridge.near`.
- **Outcome:** Historical miss → live balance WARN (`liveSnapshot`), not RPC error.

### 6.5 ENJ (Substrate)

- **Symptom:** `required result to be 32 bytes, but got 0` on Matrix RPC.
- **Cause:** Binance `en…` addresses are **Relay chain** SS58 (prefix **2135**), not Matrix (prefix **1110**).
- **Fix:** `rpcUrl`: `https://enjin-relay-chain.matrixed.link` (not matrix chain). Pruned state → query latest block; `liveSnapshot: true`.

### 6.6 HBAR HTS

- **Symptom:** `Unknown query parameter: timestamp` on `/accounts/{id}/tokens?timestamp=…`.
- **Fix:** Use `/accounts/{id}?timestamp=lte:<consensus_ts>`; read `balance.tokens[]` from account response (mirror node API).

### 6.7 TON jetton (USDT)

- **Symptom:** `can't decode address` for jetton master in tonapi path.
- **Fix:** Correct mainnet USDT master: `EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs` (URL-encode in client). Public tonapi has no historical seqno — `liveSnapshot: true`.

### 6.8 CHR (Chromia)

- **Symptom:** `dial tcp: lookup bootstrap1.chromia.com` or postchain query failure.
- **Fix:** Set `CHROMIA_NODE` to a reachable postchain endpoint (`:7740/query/<blockchainRID>/<height>`). Code tries multiple default hosts; mainnet economy RID in `ledger.json` `mint` field.

### 6.9 STARKNET / KAVAEVM / others

- Starknet: official ETH contract + correct `balanceOf` selector + felt padding (`internal/rpc/starknet.go`).
- KAVAEVM: add archive-capable RPC in `rpc-providers.json` (e.g. `evm.kava-rpc.com`).

---

## 7. Triaging bundle WARNs

Quick Python over `artifacts/.../PR01JUN26.verification-bundle.v2.json`:

```python
import json
from collections import Counter

with open("artifacts/binance/PR01JUN26/bundles/PR01JUN26.verification-bundle.v2.json") as f:
    b = json.load(f)

for v in b["verifications"]:
    if not v["verifierId"].startswith("onchain-balance"):
        continue
    warns = [x for x in v["findings"] if x.get("status") == "WARN"]
    rpc = [x for x in warns if "ledger rpc error" in x.get("note", "")]
    print(v["verifierId"], "coverage", round(v["coverage"], 4),
          "WARN", len(warns), "ledger-rpc", len(rpc))
```

| Pattern in `note` | Class | Next step |
|---|---|---|
| `ledger rpc error` | Infrastructure | §6 chain table |
| `short eth_call` | EVM archive / contract | §5, P0–P3 |
| `live … vs snapshot` | Expected without archive API | `liveSnapshot` semantics; not a bug |
| `on-chain balance != csv claim` | Mismatch or semantics | Investigate row; may be omnibus / internal ledger label |

---

## 8. PR01JUN26 regression snapshot

Approximate progression on Binance PR01JUN26 (2026-07):

| Stage | `onchain-balance-token` coverage | `ledger rpc error` (ledger verifier) |
|---|---|---|
| Before P0–P4 | ~55% | — |
| After contract + RPC pool fixes | ~65% | 44 |
| After ETH Infura/Alchemy + empty-`0x` fix | ~68% | 44 |
| After ledger RPC runbook fixes | ~68% (token) | **29** |
| Ledger coverage | — | 20.7% → **21.4%** |

Remaining ledger RPC errors (without `ALCHEMY_KEY` / `BLOCKCHAIR_API_KEY`): mostly **BCH 430** (25), **BTC Esplora pagination** (3), **CHR node** (1).

With `ALCHEMY_KEY` in environment, expect BCH + BTC historical rows to clear in one verify pass (Blockchair optional).

---

## 9. Tests

```bash
go test ./internal/rpc/... -short -count=1
```

Live probes (optional): `LEDGER_LIVE_PROBE=1 go test ./internal/rpc/... -run Live`.

Key unit tests: `TestResolveProvidersExpandsEnv`, `TestNormalizeEthCallResultEmpty`, `TestIsNearGarbageCollected`, `TestLoadProviderConfigInfuraAlchemyFromEnv`, `TestAlchemyBalanceAtHeight`.

---

## 10. Checklist before closing a ledger RPC ticket

- [ ] Bundle `ledger rpc error` count down or reclassified to `live snapshot` WARN
- [ ] Config change in `ledger.json` / `onchain.json` with correct mainnet address (cross-check block explorer)
- [ ] Env vars documented in operator shell (`~/.zshenv`)
- [ ] `go test ./internal/rpc/... -short` passes
- [ ] `por verify` + `por report` for affected snapshot
- [ ] No secrets in git diff
