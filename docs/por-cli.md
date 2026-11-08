# `por` CLI

Unified Proof-of-Reserves command-line tool. Implementation lives in `internal/porcli`, `internal/porrun`, `internal/porfetch`, `internal/porprobe`.

## Commands

| Command | Purpose |
|---------|---------|
| `por anchor` | fetch → verify → bundles → anchor calldata |
| `por verify` | re-verify cached artifacts |
| `por fetch gateio\|okx\|mexc\|…` | archive-only download |
| `por exchanges [id]` | Official PoR pages + GitHub repos (`config/exchanges/registry.yaml`) |
| `por report -exchange <id> -snapshot <id>` | assessment JSON + transparency report from cached artifacts |
| `por refresh-frequency` | recompute Section 3 of every existing report from assessment JSON |
| `por probe rpc\|tron\|…` | Binance regression / RPC diagnostics |

## Build & run

```bash
go build -o por ./cmd/por
./por anchor -exchange binance
./por anchor -exchange binance --skip-zip
./por verify -snapshot PR01JUN26
./por fetch gateio
./por probe rpc -network BSC -chainlist
./por exchanges
./por exchanges okx
```

Shell wrapper: `./scripts/por.sh anchor -exchange okx`

Legacy `./scripts/por-run.sh` and `./scripts/por-fetch.sh` forward to the same binary.

## Pipeline (anchor)

See [verifier-architecture.md](./verifier-architecture.md) for verifier matrix and verdict semantics.

## Regenerating reports and public HTML

`por report` writes `<snapshot>-assessment.json` and `<snapshot>-transparency-report.md` under `docs/reports/<exchange>/`. The HTML page is derived from the markdown:

```bash
./scripts/build-public-html.sh
```

**pandoc is pinned to 3.1.3.** Later versions changed table markup (`<tr class="odd">` → `<tr>`) and paragraph line breaking, so building the committed pages with anything else produces a large diff that has nothing to do with the content change. The script warns when the detected version differs; install 3.1.3 or point `PANDOC` at it:

```bash
PANDOC=/opt/pandoc-3.1.3/bin/pandoc ./scripts/build-public-html.sh
```

If only Section 3 (Frequency and Freshness) is stale — for example after fixing the assessment set — recompute it without re-running fetch/verify, then rebuild HTML:

```bash
go run ./cmd/por refresh-frequency --reports ./docs/reports
./scripts/build-public-html.sh
```

## Trust boundary

- Does not require a private key
- Does not broadcast transactions (prints `cast send` for the operator)
- Uses public RPC providers only

## Environment variables (verify / anchor)

Load API keys before `por verify` when using paid or quota-backed providers. Full chain-by-chain notes: [ledger-rpc-runbook.md](./ledger-rpc-runbook.md).

| Variable | Purpose |
|----------|---------|
| `INFURA_KEY` / `ALCHEMY_KEY` | ETH archive RPC in `config/rpc-providers.json` |
| `BLOCKCHAIR_API_KEY` | BCH/BTC/ZEC historical UTXO via Blockchair `state=` |
| `HELIUS_API_KEY` | Solana historical slot |
| `NEAR_RPC` / `NEAR_RPC_FALLBACK` | NEAR archival + failover |
| `CHROMIA_NODE` | Chromia postchain query base URL |

Example:

```bash
source ~/.zshenv
go run ./cmd/por verify -exchange binance -snapshot PR01JUN26 \
  -artifacts ./artifacts/binance/PR01JUN26
```

## Query anchored data on-chain

After anchoring (schema v3+), read snapshot records directly from the proxy contract — no log scan required. Full query guide: [anchor-query-api.md](./anchor-query-api.md) (view functions, `cast` / `ethers.js` examples, pre-v3 event queries, bundle verification loop).
