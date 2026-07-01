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

Full doc map: [docs/README.md](docs/README.md).

Public HTML pages are generated from markdown with `./scripts/build-public-html.sh` (requires `pandoc`). Edit the `.md` sources, rebuild, and commit both.

## Verification toolkit

This repository also ships the Go CLI used to fetch PoR artifacts, run exchange adapters, produce verification bundles, and anchor report digests on-chain.

```bash
go run ./cmd/por anchor -exchange okx
go run ./cmd/por fetch htx -zk-bundle ./public-data.zip
go run ./cmd/por verify -exchange htx -snapshot 20230910
```

Or: `./scripts/por.sh anchor -exchange okx`

Copy [`.env.example`](.env.example) to `.env` for local RPC and contract settings. Private keys stay outside the repo.

## Developer documentation

Contributor and operator docs are in-repo but not linked from the public site navigation. See [docs/README.md — Repository-only documentation](docs/README.md#repository-only-documentation).

Product context: [PRODUCT.md](PRODUCT.md) · Architecture: [docs/verifier-architecture.md](docs/verifier-architecture.md) · On-chain query: [docs/anchor-query-api.md](docs/anchor-query-api.md)
