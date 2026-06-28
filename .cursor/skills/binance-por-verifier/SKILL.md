---
name: binance-por-verifier
description: Work on ardmere's Binance Proof-of-Reserves verifier. Use when the user mentions Binance PoR, zkPoR, PR01JUN26, wallet zip, HotCold, Deposit.csv, BAPI, merkleRootHash, global zk proof, onchain-balance-hot, Stake Hub, Eth2 deposits, ERC20/BEP20 balanceOf, artifact bundles, verification bundles, or por anchor.
---

# Binance PoR Verifier

Shared verdict / onchain triage → [por-onchain-triage](../por-onchain-triage/SKILL.md).  
Adapter checklist / tiers → [exchange-por-integration](../exchange-por-integration/SKILL.md).

## Start Here

- `docs/binance-por-data-guide.md` — BAPI, wallet zip, public vs unavailable data.
- `docs/verifier-architecture.md` — verifier matrix, staking-aware design.
- `artifacts/binance/PR01JUN26-20260601-period43/AUDIT-REPORT.md` — PR01JUN26 root-cause analysis.

## Core Mental Model

Three data planes (do not conflate):

- **BAPI summary:** `merkleRootHash`, per-coin liabilities/reserves.
- **Wallet zip:** HotCold + Deposit CSV — reserve-side evidence, not zk data.
- **zk proof data:** `proof.csv`, `cex_assets_info.json`, vk — **not public** via BAPI.

`merkleRootHash` = user account tree root, not wallet-tree root.

Keep `global-zk-proof@0` as **UNVERIFIABLE** unless proof artifacts are obtained.

## HotCold Rules (Binance-specific)

CSV `balance` = accounted balance (may include staking/deposits).

`onchain-balance-hot@2.1`:

- **ETH|ETH:** `eth_getBalance + Eth2 DepositContract` via `ETHERSCAN_API_KEY`.
- **BNB|BSC:** liquid + Stake Hub `getPooledBNB` + `lockedBNBs`.

Contracts: Eth2 `0x000…0fa`, Stake Hub `0x…2002`, FDUSD `0xc5f0f7b66764F6ec8C8Dff7BA683102295E16409`.

BNB archive gaps → **WARN**, not FAIL. See §BSC archive below.

## BSC Archive (PR01JUN26-tested)

| Endpoint | Archive? |
|----------|----------|
| `bsc-rpc.publicnode.com` | No (~100–200 blocks) |
| `bsc-dataseed*.bnbchain.org` | Full node; `missing trie node` on old blocks |
| `bsc-mainnet.public.blastapi.io` | Yes; Stake Hub scan OK |
| `bsc.rpc.sentio.xyz` | Yes |
| `bsc.drpc.org` | Yes but frequent 429/408 |

```bash
go run ./cmd/por probe rpc -network BSC -height 101590091 -chainlist
go run ./cmd/por probe stakehub -address 0x86523c87c8ec98c7539e2c58cd813ee9d1a08d96 -height 101590091
```

## Binance-Specific Triage

Use [por-onchain-triage](../por-onchain-triage/SKILL.md) first; then apply:

| Pattern | Binance note |
|---------|--------------|
| BNB Stake Hub **>** CSV | Stable over H-1000..H; rewards/exchange-rate vs conservative CSV — **WARN** |
| ETH/USDT fail at **H** only | Often matches **H-1** — height boundary |
| BSC native small gap | Try **H-100** window |
| `POL|ETH` CSV >> ERC20 | Polygon-native / internal ledger — see AUDIT-REPORT §7.8 |
| FDUSD `0 bytes` | Was mistyped contract — verify mapping |

```bash
go test ./internal/verifier/ -run TestPolEthProbe0xa64b -v
zsh -lc 'go run ./cmd/por verify -snapshot PR01JUN26'   # ETHERSCAN in profile
```

Tron `USDT|TRX` rules → [por-onchain-triage](../por-onchain-triage/SKILL.md) §Tron.

## PR01JUN26 Baseline

Regression snapshot `PR01JUN26-20260601-period43`:

| Verifier | Supported-row result |
|----------|---------------------|
| internal-consistency@1.0 | PASS |
| onchain-balance-hot@2.1 | 48 pass / 2 warn / 1 fail / 0 rpc-error |
| onchain-balance-token@1.7 | 558 pass / 5 warn / 5 fail / 1 unv / 142 rpc-error **of 711** |

Interpretation: no clear insufficiency signal; BNB surpluses = staking accounting; ETH/USDT fails often H-1 boundary; rpc-errors = L2 archive / wrong mappings.

`verificationBundleRoot`: `0x2c0baa8006e34e0a7682e284273aac01b8e48ad97de097db62356b3a5ea95528`.

Sonic SFC: `go run ./cmd/por probe sonic-sfc`; residual −6M S at `0x64de13c46f…`.

## Verdict Discipline

→ [por-onchain-triage](../por-onchain-triage/SKILL.md). Never frame RPC/staking gaps as Binance fraud.

## Operational Notes

- `ETHERSCAN_API_KEY` required for ETH deposit indexing.
- `SOLANA_INDEX_API_KEY` (or `SOLANA_HISTORY_API_KEY`) enables slot-indexed SPL balance via Solana Index; without it, `BOME|SOL` etc. compare live RPC vs CSV slot → WARN.
- PR01JUN26 = main regression for ETH deposit + BNB Stake Hub.
- Regenerate verification bundles after verifier changes.

## Current Next Steps

1. Close `0x64de13c46f…` −6M S and `0xf977…` −157 WBTC ETH row.
2. Expand EVM mappings (POL|MATIC, L2 long-tail).
3. Native-chain coverage gaps → ledger.json where applicable.
