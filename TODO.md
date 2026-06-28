# TODO

## Now

- [x] Complete core PoR verification service code.
- [x] Complete PoR Stage transparency framework (`docs/por-transparency-framework.md`).
- [x] Document product audience and service positioning (`docs/ardmere-service-audience.md`).
- [x] Convert the Stage framework into machine-readable rules:
  - `stage_requirements.yaml`
  - `evidence_level.yaml`
  - `risk_flags.yaml`
  - `exchange-assessment.v1.schema.json`
- [x] Define a standard `exchange transparency report` template.
- [x] Define artifact metadata schema: snapshot, source URL, hash, fetch time, verifier output, missing artifacts, risk flags.

## Phase 1: Public Methodology + First Reports

- [x] Publish methodology page from `docs/por-transparency-framework.md`.
- [x] Create first exchange report: OKX as the Stage 1 reference sample.
- [x] Create first exchange report: Binance as Gen 2 / E2 but Stage 1 blocked sample.
- [x] Create first exchange report: Bybit as Stage 0 / Merkle inclusion sample.
- [x] Build the first public comparison table:
  - PoR Stage
  - Gen
  - E Level
  - missing artifacts
  - risk flags
  - last verified snapshot
- [x] Add artifact archive index for first reports.
- [x] Link first reports from the GitHub Pages landing/docs navigation.

## Phase 2: Institution / Regulator Package

- [ ] Produce a regulator-facing “effective PoR minimum standard” brief.
- [ ] Produce a risk matrix for exchanges:
  - Stage 0 not accepted as effective PoR
  - Stage 1 minimum effective PoR
  - Stage 2 long-term best practice
- [ ] Produce an institution-facing PoR risk memo template.
- [ ] Document traditional audit vs PoR complementarity as a standalone explainer.
- [ ] Define standard wording for `UNVERIFIABLE`, `Stage blocked`, `high-friction inclusion proof`, and `sample-only`.

## Phase 3: API and Tooling

- [ ] Expose verification API:
  - input artifacts or artifact URLs
  - return verifier results
  - return missing artifacts
  - return risk flags
  - return suggested Stage / E Level
- [ ] Implement rule engine that maps artifacts + verifier output to preliminary Stage.
- [ ] Implement historical artifact archive:
  - raw files
  - hashes
  - source URLs
  - fetch timestamps
  - verifier output
- [ ] Build frontend verdict cards for public dashboard.
- [ ] Add machine-readable report export.

## Phase 4: Exchange Readiness / Remediation

- [ ] Create Stage upgrade checklist for exchanges.
- [ ] Create artifacts schema and disclosure checklist for Stage 1.
- [ ] Create Stage 2 readiness checklist:
  - low-friction inclusion proof
  - canonical anchoring
  - DA
  - permissionless verification
  - business-consistent constraints
- [ ] Create user verification UX checklist for inclusion proof:
  - Web/WASM/GUI one-click verification
  - proof export
  - local reproducibility
  - clear error messages
  - verification rate / participation disclosure
- [ ] Define independence policy for paid technical assessments vs public ratings.

## Verifier Backlog

- [ ] Complete the BNB Stake Hub full validator scan for `0x86523...08d96` — blocked: BSC archive RPC unavailable (`drpc 408/429`, `publicnode` tip-only, Binance full nodes `missing trie node`). Use `go run ./cmd/por probe rpc -network BSC -height 101590091 -chainlist`.
- [ ] Retry BSC archive scan when `bsc.drpc.org` or another archive provider is healthy.
- [ ] Resolve remaining Binance verifier gaps:
  - POL|ETH `0xa64b...`
  - 2 BTC FAIL
  - Sonic −6M
  - Hot micro-FAIL
  - native chains: BTC / DOGE / XRP / SOL
- [ ] Add beacon effective balance support for ETH validators.
- [ ] Add `wallet_ownership_proof` verifier if Binance or other exchanges publish wallet signatures.
- [ ] Activate `global-zk-proof@1` if `proof.csv`, `cex_assets_info.json`, and verifying keys become available.
- [ ] Add wrapped / cross-chain asset reconciliation rules.

## Guardrails

- Do not treat unavailable RPC/indexer data as Binance fraud; emit `WARN` or `UNVERIFIABLE`.
- Do not commit `.env`, API keys, or wallet private keys. Keep `PRIVATE_KEY` in `~/.zshrc`, not in the repo.
- Public ratings must bind to artifacts, hash, URL, verifier output, and timestamp.
- Paid technical assessment must be separated from independent public rating.
- Missing data is not PASS; mark it as `UNVERIFIABLE`.
- Run `go test ./...` after verifier changes.
