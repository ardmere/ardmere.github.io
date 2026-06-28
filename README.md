# ardmere.github.io

Decentralized. Verified. Visible.

Public PoR transparency registry and Go verification toolkit.

## Public site

- [Methodology](docs/por-transparency-framework.md) — PoR Stage transparency framework
- [Exchange reports](docs/reports/exchange-comparison.md) — public comparison table
- [Artifact archive](docs/reports/artifact-archive-index.md) — report artifact index
- [OKX report](docs/reports/okx/506872725-transparency-report.md) — Stage 1 reference
- [Binance report](docs/reports/binance/PR01JUN26-transparency-report.md) — Stage 0 blocked
- [Bybit report](docs/reports/bybit/2025061709-transparency-report.md) — Stage 0 Merkle inclusion
- [Audience](docs/ardmere-service-audience.md) — product positioning
- [Deposit spot-check](verify/deposit/) — user spot-check page ([spec](docs/deposit-spot-check.md))

Full doc map: [docs/README.md](docs/README.md).

## Developer documentation (repository only)

Not linked from the public site navigation. See [docs/README.md](docs/README.md#repository-only-documentation).

```bash
go run ./cmd/por anchor -exchange okx
go run ./cmd/por fetch htx -zk-bundle ./public-data.zip
go run ./cmd/por verify -exchange htx -snapshot 20230910
```

Or: `./scripts/por.sh anchor -exchange okx`
