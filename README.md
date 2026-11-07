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

Copy [`.env.example`](.env.example) to `.env` for local RPC and contract settings. API keys for verify (`INFURA_KEY`, `BLOCKCHAIR_API_KEY`, …) belong in `~/.zshenv` — see [ledger-rpc-runbook](docs/ledger-rpc-runbook.md). Private keys stay outside the repo.

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
