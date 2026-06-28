---
name: exchange-por-integration
description: Integrate a new exchange Proof-of-Reserves adapter into ardmere (fetch, normalize, verifiers, artifacts, CLI). Use when adding or extending gateio, okx, binance, or any exchange PoR; when discussing data planes, VerifierProfile, UNVERIFIABLE stubs, appState/CDN fetch, browser-api import, solvency-claim, internal-consistency, address-ownership, or exchange-por data guides.
---

# Exchange PoR Integration

**Skill map** (read the narrowest skill that matches the task):

| Task | Skill |
|------|--------|
| Add adapter, stubs, tiers, checklist | **This file** |
| Onchain FAIL/WARN/RPC triage (any exchange) | [por-onchain-triage](../por-onchain-triage/SKILL.md) |
| Binance deep (PR01JUN26, Stake Hub, BSC) | [binance-por-verifier](../binance-por-verifier/SKILL.md) |
| OKX deep (zk v2, wallet CSV, 506872725) | [okx-por-verifier](../okx-por-verifier/SKILL.md) |
| Field URLs / artifact layout | `docs/{binance,okx,gate}-por-data-guide.md` |

## Start Here

1. Map **data planes** (summary / wallet / liability-zk / user-inclusion) and mark each **public | login | absent**.
2. Pick **integration tier** (weak / strong no-login) — see §Tiers.
3. Follow **Integration checklist** — do not skip `VerifierProfile` stubs with reasons.
4. Run `go test ./...` after verifier changes.

## Data Planes (never conflate)

| Plane | Typical contents | Common mistake |
|-------|------------------|----------------|
| Summary | ratios, per-coin liabilities/reserves, merkle root | Treating root as wallet-tree hash |
| Wallet / reserve | address CSV, signatures, block heights | Treating CSV balance as liquid EOA only |
| Liability / global ZK | proof.csv, sum_proof_data.json, vk | Expecting summary alone to prove solvency |
| User inclusion | user_config, merkle path | Confusing with exchange reserve evidence |

## Source Discovery Order

1. Official Open API docs — if no PoR endpoints, **stop** (Gate API v4 has none).
2. Page SSR / `<script id="appState">` embedded JSON (OKX).
3. Page-internal Web API (Gate `/api/web/v1/proof-of-reserves/*`; Akamai may block datacenter IPs).
4. Browser DevTools capture → fixtures + `ImportSource: browser-api`.
5. Static CDN URLs (OKX `static.okx.com/cdn/okx/por/...`).

## Integration Checklist

```
[ ] internal/<exchange>api/          fetch + parse
[ ] internal/exchanges/<id>/         adapter, normalize, store
[ ] artifacts/<id>/<auditId>/raw/  content-addressed artifacts + fetch.json
[ ] exchangereg.Register
[ ] VerifierProfile { Shared, Stubs } with honest UNVERIFIABLE reasons
[ ] verifyrun dispatch (pass Exchange to onchain verifiers)
[ ] config/exchanges/<id>/onchain.json + ledger.json (if onchain active)
[ ] por-fetch subcommand; por anchor / verify
[ ] docs/<id>-por-data-guide.md + fixtures/ (small samples only)
[ ] go test ./...
```

**Layout:** `artifacts/<exchange>/<auditId>/{raw,bundles,fetch.json}`.

**FetchOpts:** `SummaryPath`, `WalletZipPath`, `LiabilityZipPath`, `SkipWalletZip`, `SkipLiabilityZip`, `ImportSource`.

## Verifier Profile Rules

### Always active (if summary exists)

- `artifact-integrity@1` — SHA256 of archived artifacts.
- `solvency-claim@1` — **self-reported only**; finding note: *"self-reported exchange summary only; does not prove reserve authenticity"*.

### Activate only when public data exists

| Verifier | Needs |
|----------|--------|
| `internal-consistency@1` | Wallet artifact matching summary semantics |
| `address-ownership@*` | Signed address CSV (OKX) or equivalent |
| `onchain-balance-*` | Address list + heights + per-exchange onchain/ledger config |
| `global-zk-proof@*` | Public proof bundle or official validator binary |

### Stubs — mandatory

Every absent dimension → `UNVERIFIABLE` + **specific reason** in `internal/verifier/stubs.go`. Never omit a listed verifier; never mark missing data as `FAIL`.

### Verdict discipline

See [por-onchain-triage](../por-onchain-triage/SKILL.md) for onchain-specific FAIL vs WARN rules.

- `PASS` / `FAIL` — evidence complete and compared.
- `WARN` / `PARTIAL` — RPC, custody, staking attribution, or zk binary incomplete.
- `UNVERIFIABLE` — required public artifact or verifier missing.

On-chain anchor ≠ strong verification (Gate weak anchor is valid but limited).

## Internal Consistency — per exchange

| Exchange | Compare | Aggregate source |
|----------|---------|------------------|
| Binance | `exchangeBalance`, `thirdPartyCustody` | HotCold + Deposit CSV by custodian |
| OKX | `exchangeReserveBalances` | CSV **top section** `coin,amount` (+ ETH staking file) |
| Gate | — | UNVERIFIABLE (no public wallet CSV) |

OKX: do not sum address rows by coin; `custodyReserveBalances` → WARN only.

## Official Verifier Integration

| Exchange | Repo | Pattern |
|----------|------|---------|
| Gate | `gateio/proof-of-reserves` | Future `global-zk-proof@gateio-1` on login tar.gz |
| OKX | v1 + v2 | v1 signatures; v2 `OKX_ZK_STARK_VALIDATOR` — [okx-por-verifier](../okx-por-verifier/SKILL.md) |
| Binance | `zkmerkle-proof-of-solvency` | [binance-por-verifier](../binance-por-verifier/SKILL.md) |

## Tiers

### Weak no-login (Gate)

```bash
go run ./cmd/por fetch gateio
go run ./cmd/por verify -exchange gateio -snapshot <auditId> \
  -artifacts ./artifacts/gateio/<auditId>
```

Profile: integrity + solvency only; wallet/onchain/global-zk stubbed.

### Strong no-login (OKX)

→ [okx-por-verifier](../okx-por-verifier/SKILL.md)

### Full onchain (Binance)

→ [binance-por-verifier](../binance-por-verifier/SKILL.md)

## Gate vs OKX (quick reference)

| Dimension | Gate | OKX |
|-----------|------|-----|
| No-login strength | Weak | Strong |
| Summary source | Web API | Page appState |
| Wallet addresses | Not public | Public + signatures |
| Global ZK | Login tar.gz | Public `sum_proof_data.json` |
| On-chain audit | UNVERIFIABLE | Runnable with RPC |

## Anti-Patterns

1. Gate Open API v4 for PoR (use page/CDN).
2. OKX signatures on `Network` column instead of `coin`.
3. OKX internal-consistency via address-row sum.
4. Treating audit PDFs as period PoR proof.
5. `UNVERIFIABLE` without reason, or missing data as `FAIL`.
6. Low **coverage** interpreted as mass onchain failure (check **summary** finding pass rate).
7. Import cycles: proof-bundle resolution in `verifyrun`, not `verifier` → `bundle`.

## Code Entry Points

| Area | Path |
|------|------|
| Registry | `internal/exchangereg/reg.go` |
| Profile dispatch | `internal/verifyrun/runner.go` |
| Stub reasons | `internal/verifier/stubs.go` |
| Onchain config | `config/exchanges/<id>/` |
| Gate / OKX / Binance | `internal/exchanges/{gateio,okx,binance}/` |

## Additional Resources

- `docs/verifier-architecture.md`, `docs/STRUCTURE.md`, `docs/exchange-tiers.md`
- Onchain triage: [por-onchain-triage](../por-onchain-triage/SKILL.md)
