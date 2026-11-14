# ardmere

**Having a Safe CEX.**

Decentralized. Verified. Visible.

[Public registry](https://ardmere.github.io/) · [Reports](docs/reports/exchange-comparison.html) · [Standard](docs/effective-por-standard.html) · [Methodology](docs/por-transparency-framework.html)

---

ardmere is an independent **Proof-of-Reserves transparency registry** and Go verification toolkit. It helps external observers determine whether an exchange PoR disclosure is genuinely verifiable, whether it meets the minimum **effective PoR** standard, and which artifacts or proof constraints are missing.

Public ratings bind to artifacts, hashes, verifier outputs, timestamps, and explicit `UNVERIFIABLE` gaps — not exchange marketing claims.

## Mission

Make centralized exchange reserve transparency **evidence-bound and publicly auditable**.

We publish exchange PoR Stage ratings, transparency reports, artifact indexes, and reproducible verifier results so regulators, institutional allocators, and independent reviewers can assess custody risk without trusting slogans like “100% backed”, “audited”, or “ZK verified” on their own.

## Minimum standard

1. **Stage 0 is not effective PoR.** Incomplete or unverifiable disclosures do not meet the minimum evidentiary bar.
2. **Stage 1 is the minimum effective threshold.** Effective PoR requires public evidence for reserves, liabilities, and key proofs.
3. **PoR and traditional audit answer different questions.** They complement each other; neither substitutes for the other.

Gap reporting (`UNVERIFIABLE`, not PASS) is defined in the [PoR Stage Framework — §1.3](docs/por-transparency-framework.html#gap-reporting-discipline).

Full framework: [PoR Stage Framework](docs/por-transparency-framework.html).  
Public policy standard: [Effective PoR Standard](docs/effective-por-standard.html).

## Who it is for

| Audience | Use |
| --- | --- |
| Regulators & policy researchers | Define and apply an actionable minimum PoR standard |
| Institutional allocators | Compare exchange transparency and integrate PoR signals into risk workflows |
| Independent reviewers | Re-run verifiers, inspect artifacts, and audit evidence gaps |

Retail users benefit indirectly; exchanges are assessed subjects, not the primary audience. See [service audience](docs/ardmere-service-audience.html).

## Public site

| Resource | Description |
| --- | --- |
| [Homepage](index.html) | Public PoR transparency registry |
| [Methodology](docs/por-transparency-framework.html) | PoR Stage framework |
| [Effective PoR standard](docs/effective-por-standard.html) | Minimum policy standard for public ratings |
| [Insights](docs/insights/index.html) | Editorial explainers, research, and event analysis (not ratings) |
| [Exchange reports](docs/reports/exchange-comparison.html) | Comparison table and per-exchange reports; [evidence archive](docs/reports/artifact-archive-index.html) for bundles and hashes |
| [OKX report](docs/reports/okx/508399035-transparency-report.html) | Stage 1 reference |
| [Binance report](docs/reports/binance/PR01JUN26-transparency-report.html) | Stage 0 blocked |
| [Bybit report](docs/reports/bybit/2025061709-transparency-report.html) | Stage 0 Merkle inclusion |
| [Deposit spot-check](verify/deposit/) | Independent deposit sample verification ([spec](docs/deposit-spot-check.md)) |
| [On-chain query](docs/contract-query-guide.html) | How to query `ArdmerePoRAnchor` for homepage audit / anchor records |

Full doc map: [docs/README.md](docs/README.md).

Public HTML pages are generated from markdown with `./scripts/build-public-html.sh` (requires **pandoc 3.1.3** — see [docs/por-cli.md](docs/por-cli.md#regenerating-reports-and-public-html); other versions rewrite table markup and bury content changes in noise). The homepage, snapshot history, and exchange comparison pages are regenerated from `docs/reports/*/*-assessment.json` (`go run ./cmd/por site`, also run automatically after `por report` / `por batch`). Edit per-snapshot editorial notes in `config/site/snapshot-notes.yaml` when needed; rebuild HTML after markdown changes and commit both.

## Verification toolkit

This repository also ships the Go CLI used to fetch PoR artifacts, run exchange adapters, produce verification bundles, and anchor report digests on-chain.

```bash
go run ./cmd/por anchor -exchange okx
go run ./cmd/por fetch htx -zk-bundle ./public-data.zip
go run ./cmd/por verify -exchange htx -snapshot 20230910
```

Or: `./scripts/por.sh anchor -exchange okx` (loads `~/.zshenv` and `.env` for RPC keys)

Copy [`.env.example`](.env.example) to `.env` for local RPC and contract settings. API keys for verify (`INFURA_KEY`, `BLOCKCHAIR_API_KEY`, …) belong in `~/.zshenv` or the process environment — see [ledger-rpc-runbook](docs/ledger-rpc-runbook.md). Testnet `PRIVATE_KEY` / `ANCHOR_SIGNER` may come from `~/.zshenv` or cloud Cursor agent env injection; keep private keys outside the repo.

## How to trust these results

Independent reviewers should treat ardmere as **evidence machinery**, not a rating oracle. This section states what you can verify yourself, what on-chain anchors attest, how agent skills fit in, and what we do *not* claim.

### What we claim — and what we do not

| We publish | We do **not** claim |
| --- | --- |
| Open-source verifiers, assessment JSON, bundle Merkle roots, explicit `UNVERIFIABLE` gaps | That any exchange is solvent, that your balance was included, or any trading / allocation advice |
| Immutable on-chain digests (and, from schema v5, `custodyRiskScore` at anchor time) for anchored snapshots | That Sepolia anchors are production finality ([Base mainnet](docs/decisions.md) is the stated production target) |
| Reproducible custody and inclusion **readouts** from skills | A `decision` field or allow / reduce / block output from skills |

> **Operational trust today:** Anchors are written by a **permissioned `signer`** on **Ethereum Sepolia** (testnet). Anyone can read records with `eth_call`, but you should not equate “on-chain” with “decentralized witness” until production deployment and policy are published. See [`deployments.md`](docs/deployments.md).

### Verification hosting and on-chain push (today vs roadmap)

**Today.** Fetch → verify → report runs on ordinary machines (your laptop, CI, or operator infrastructure). A permissioned signer broadcasts `anchorSnapshot` with bundle roots and summary fields. Readers verify by **re-running the same open-source verifier** and comparing roots — not by trusting the operator’s narrative.

**Exploring (no committed date).** We are considering **TEE-backed hosted verification**: run the fetch → verify pipeline inside a confidential execution environment (enclave / confidential VM — concrete stack TBD), produce assessment and bundle roots there, and use **remote attestation** to tie “this digest set was produced by this measured build” to **authorized snapshot anchoring**. TEE would bind *who ran* the replay under a published policy; it would **not** replace public verifier replay, and it would **not** turn `UNVERIFIABLE` into `PASS`.

Full verifier replay remains the ground truth. Details: [verifier architecture §6.4](docs/verifier-architecture.md#64-verification-loop-closure).

### Trusting assessment outputs (`*-assessment.json` and public reports)

Each snapshot assessment is produced by `por fetch` → `por verify` → `por report`. It carries `verifierResults`, `riskFlags`, `porStage`, `effectivePoR`, and `artifactBundleRoot` / `verificationBundleRoot`.

**How to trust it:**

1. Check out the **same git revision** (or release tag) as the published report.
2. Re-run `go run ./cmd/por verify -exchange <slug> -snapshot <auditId>` (or the matching `./scripts/por.sh` invocation).
3. Recompute bundle Merkle roots and compare with the assessment fields.
4. Accept gaps only when the verifier returns **`UNVERIFIABLE` with a reason** — never when marketing copy fills the gap ([framework §1.3](docs/por-transparency-framework.html#gap-reporting-discipline)).

Site HTML is generated from these JSON files (`go run ./cmd/por site`); editing HTML without regenerating does not change the evidentiary record.

### Trusting on-chain anchors (`ArdmerePoRAnchor`)

The contract stores **digests and summary metadata** per snapshot, not raw artifacts or proofs.

| On-chain (v5+) | Off-chain (this repo) |
| --- | --- |
| `artifactBundleRoot`, `verificationBundleRoot`, verdict summary, stage, coverage bps, `custodyRiskScore` | Full artifacts, verifier output, transparency prose |

**`snapshotId` (schema v6):** `keccak256(lowercaseExchange ‖ 0x00 ‖ auditId)` so identical audit labels on different exchanges do not collide; legacy rows used `keccak256(auditId)` only — clients try composite first, then legacy ([ADR-006](docs/decisions.md)). Query helpers: [contract query guide](docs/contract-query-guide.html), [anchor query API](docs/anchor-query-api.md).

**Independent loop (abbreviated):**

1. Resolve `snapshotId` for `(exchange, auditId)` (composite or legacy).
2. `getSnapshot` on the proxy (Sepolia today).
3. Load `docs/reports/<exchange>/<auditId>-assessment.json` and compare roots field-by-field.
4. If roots match, the file you hold is the bundle that was anchored; if not, treat as `ANCHOR_MISMATCH` and re-verify — do not edit the assessment to match the chain.

An optional **NEAR testnet** mirror in [`ardmere-near-demo`](https://github.com/ardmere/ardmere-near-demo) follows the same composite id semantics for digest replication; production anchoring policy is documented there separately.

### Trusting agent skill outputs

Skills are **read-only interpreters** over assessments and chain state. Run `./scripts/check-skills.sh` locally (`--live` hits Sepolia RPC).

| Skill | Trust model | Mis-trust guard |
| --- | --- | --- |
| [`por-custody-risk`](skills/por-custody-risk/SKILL.md) | Applies the public [decision table](skills/por-custody-risk/references/decision-table.md); can cross-check on-chain `custodyRiskScore` when schema ≥ 5 | `conclusive: false` / `riskScore: null` when RPC fails — **not** “zero risk” |
| [`por-inclusion-proof`](skills/por-inclusion-proof/SKILL.md) | User-held proof verified locally against the published Merkle root | Inclusion success does not prove exchange-level PoR |

Neither skill emits a `decision` field or trade instructions.

### Reviewer checklist

1. Pin the repository revision used for the published assessment.
2. Re-run `por verify` for that exchange and snapshot id.
3. Match bundle roots to the assessment and, if present, to `getSnapshot`.
4. Run `skills/por-custody-risk/scripts/por-check.sh <assessment.json>` (or `risk_score.py`) and note `custodyRiskScoreMatchesComputed` when anchored.
5. Optionally `cast call` the proxy using [anchor-query-api.md](docs/anchor-query-api.md) examples.
6. Run inclusion checks only when you hold a user proof file; bind roots to the same snapshot.

## Agent skills

Two skills read what this repository already publishes and decide nothing on their own. Both run on a bare `python3` with no web3 library and no Go build, and both read the `ArdmerePoRAnchor` proxy directly over JSON-RPC (`POR_ANCHOR_CONTRACT` / `POR_RPC_URL` override the defaults).

| Skill | Question it answers | Read |
| --- | --- | --- |
| [`por-custody-risk`](skills/por-custody-risk/SKILL.md) | How much flagged custody risk does this assessment show, and are these the bytes that were anchored? | [`SKILL.md`](skills/por-custody-risk/SKILL.md) · [`decision-table.md`](skills/por-custody-risk/references/decision-table.md) |
| [`por-inclusion-proof`](skills/por-inclusion-proof/SKILL.md) | Is this user's own proof leaf in the tree that produced the published root? | [`SKILL.md`](skills/por-inclusion-proof/SKILL.md) · [`schemes.md`](skills/por-inclusion-proof/references/schemes.md) |

```bash
# Score an assessment against the on-chain record
skills/por-custody-risk/scripts/por-check.sh docs/reports/kucoin/20261031-assessment.json

# Verify a user's own proof locally, bound to the anchored root
skills/por-inclusion-proof/scripts/inclusion_check.py --list-schemes
skills/por-inclusion-proof/scripts/inclusion_check.py --exchange bybit \
  --onchain 2025061709 ~/Downloads/myProof.json
```

Neither skill has a `decision` field. `por-custody-risk` returns a 0–100 score with its factors, and returns `conclusive: false` with `riskScore: null` when the chain cannot be read. `por-inclusion-proof` verifies the proof on the user's own machine and prints verdicts only: it never echoes the account id, nonce or balances, and writes no file. Run the checks with `./scripts/check-skills.sh` (add `--live` to also read Sepolia).

## Developer documentation

Contributor and operator docs are in-repo but not linked from the public site navigation. See [docs/README.md — Repository-only documentation](docs/README.md#repository-only-documentation).

Product context: [PRODUCT.md](PRODUCT.md) · Architecture: [docs/verifier-architecture.md](docs/verifier-architecture.md) · On-chain query: [docs/contract-query-guide.html](docs/contract-query-guide.html) ([ABI reference](docs/anchor-query-api.md))
