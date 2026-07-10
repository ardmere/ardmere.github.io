# OKX Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: okx`).

OKX publishes PoR artifacts without login. ardmere fetches them via page-embedded JSON and static CDN downloads.

## Public sources

| Source | URL | Content |
|--------|-----|---------|
| Detail page | `https://www.okx.com/proof-of-reserves/detail` | `appState` JSON → `auditRootInfo` (ratios, balances, merkle hash) |
| Download page | `https://www.okx.com/proof-of-reserves/download` | `auditList` with CDN URLs per audit |
| Wallet zip | `https://static.okx.com/cdn/okx/por/chain/por_csv_*_V2.zip` | Signed address CSV + ETH staking CSV |
| Liability zip | `https://static.okx.com/cdn/okx/por/merkel/por_<auditId>_proof_data.zip` | `sum_proof_data.json` (global zk-STARK) |

## Local layout

```
artifacts/okx/<auditId>/
  raw/          — content-addressed summary, wallet zip, liability zip
  bundles/      — artifact + verification bundles
  fetch.json    — fetch metadata
```

## CLI

```bash
# Archive only (~48 MB wallet + liability zips)
go run ./cmd/por fetch okx

# Full pipeline (fetch + verify + anchor calldata)
go run ./cmd/por anchor -exchange okx

# Re-verify cached snapshot (full: disclosure + zk + on-chain RPC)
go run ./cmd/por verify -exchange okx -snapshot 508399035 \
  -artifacts ./artifacts/okx/508399035

# Zk-only iteration (skip on-chain RPC; ~90s with validator built)
go run ./cmd/por verify -exchange okx -snapshot 508399035 \
  -artifacts ./artifacts/okx/508399035 -skip-rpc -write-bundle
```

Flags: `--skip-zip`, `--skip-liability`, `--skip-rpc`, `-summary-path`, `-wallet-zip`, `-liability-zip`.

## Local toolchain (zk validator)

OKX liability proofs use Plonky2 (`okx/proof-of-reserves-v2`). ardmere calls `zk-por-cli verify-global` when `OKX_ZK_STARK_VALIDATOR` is set.

### One-time setup

```bash
# Rust (if missing)
curl -fsSL https://sh.rustup.rs | sh -s -- -y

# Source (git clone or tarball if GitHub is slow)
cd ~/Project
curl -fsSL https://github.com/okx/proof-of-reserves-v2/archive/refs/heads/main.tar.gz \
  | tar -xz && mv proof-of-reserves-v2-main proof-of-reserves-v2

cd proof-of-reserves-v2
cargo build --features zk-por-core/verifier --release \
  --package zk-por-cli --bin zk-por-cli
mkdir -p bin && cp target/release/zk-por-cli bin/
```

Add to `~/.zshrc`:

```bash
export GOMODCACHE="${GOMODCACHE:-$HOME/go/pkg/mod}"
export OKX_ZK_STARK_VALIDATOR="$HOME/Project/proof-of-reserves-v2/bin/zk-por-cli"
```

### Verify modes

| Mode | Command | Duration (508399035) | global-zk-proof |
|------|---------|----------------------|-----------------|
| Full | `verify` (no `--skip-rpc`) | ~12 min | PASS with validator |
| Zk + disclosure | `verify -skip-rpc` | ~90 s | PASS with validator |
| No validator | omit env var | same | PARTIAL (0.75) |

Without `OKX_ZK_STARK_VALIDATOR`, structure + merkle binding still run; cryptographic `verify-global` is skipped.

### Cursor / sandbox gotchas

- If `go run` fails with missing files under `/var/folders/.../go-mod/`, set `GOMODCACHE=$HOME/go/pkg/mod`.
- If `cargo build` says `Finished` but `target/release/zk-por-cli` is missing, `CARGO_TARGET_DIR` may point at a sandbox cache — build with a normal shell or copy the binary to `~/Project/proof-of-reserves-v2/bin/`.

## Verifier profile (`okx`)

| Verifier | Status | Notes |
|----------|--------|-------|
| `artifact-integrity@1` | active | SHA256 over archived artifacts |
| `solvency-claim@1` | active | `capitalRatio >= 100%` per reserve currency (self-reported) |
| `internal-consistency@1` | active | `exchangeReserveBalances` vs signed address CSV aggregate (+ ETH staking) |
| `address-ownership@okx-1` | active | Uses [okx/proof-of-reserves](https://github.com/okx/proof-of-reserves) signature verification |
| `onchain-balance-hot@2.1` | active | Native balances at snapshot height (`config/exchanges/okx/onchain.json`); subsamples top balances when >800 rows |
| `onchain-balance-token@2.0` | active | ERC20/BEP20/TRC20 where configured (OKX coin labels e.g. `USDT-TRC20`) |
| `onchain-balance-ledger@1.2` | active | UTXO/XRPL/Solana/Aptos via `config/exchanges/okx/ledger.json` (BTC, DOGE, LTC, …) |
| `global-zk-proof@okx-1` | active | Structure + summary merkle bind + optional `OKX_ZK_STARK_VALIDATOR` crypto |

## Wallet CSV format

Multi-section CSV inside the wallet zip:

1. Top: `coin,amount` — exchange reserve totals
2. Blank line
3. Address rows: `coin,Network,Snapshot Height,address,amount,message,signature1,...`

Message is always `I am an OKX address`.

## Global zk verification

Public liability zip contains **`sum_proof_data.json` only** (no user inclusion proofs).

Build OKX v2 validator from [`okx/proof-of-reserves-v2`](https://github.com/okx/proof-of-reserves-v2):

```bash
cargo build --features zk-por-core/verifier --release \
  --package zk-por-cli --bin zk-por-cli
export OKX_ZK_STARK_VALIDATOR=/path/to/target/release/zk-por-cli
```

ardmere invokes `verify-global --proof-path` (circuit rebuild, secure mode). Without the env var, verdict is `PARTIAL` (structure only).

Detail: `.cursor/skills/okx-por-verifier/SKILL.md`.

## References

- [OKX PoR page](https://www.okx.com/proof-of-reserves)
- [OKX verification guide](https://www.okx.com/support/hc/en-us/articles/10781041719437)
- Upstream repos: [registry `okx`](../config/exchanges/registry.yaml) — [okx/proof-of-reserves](https://github.com/okx/proof-of-reserves) (address signatures), [okx/proof-of-reserves-v2](https://github.com/okx/proof-of-reserves-v2) (Plonky2 global zk)
