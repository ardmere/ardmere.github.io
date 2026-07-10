# Repository structure

> Repository map for contributors. Not part of the public site navigation. Public doc map: [README.md](./README.md).

ardmere is a GitHub Pages site plus a Go PoR verification toolkit.

## Top level

```
.
├── cmd/por/                # single CLI binary
├── internal/               # libraries (porcli, porrun, pipeline, verifier, …)
├── config/exchanges/
│   ├── registry.yaml       # official PoR pages + GitHub upstream repos
│   └── <exchange>/         # onchain.json, ledger.json, …
├── config/rating/          # machine-readable PoR Stage / E Level / risk flag rules
├── docs/
├── fixtures/
├── schemas/                # JSON schemas for reports, assessments, manifests
├── scripts/por.sh
└── artifacts/              # gitignored local snapshots
```

## Go layout (`internal/`)

```
internal/
├── porcli/         # CLI dispatch: anchor | verify | fetch | probe | exchanges
├── exchangeregistry/  # config/exchanges/registry.yaml loader
├── porrun/         # anchor + verify
├── porfetch/       # fetch gateio | okx
├── porprobe/       # rpc | tron | stakehub | sonic-sfc
├── pipeline/       # verify run, bundle build, anchor calldata
├── exchange/       # Adapter interface + Capabilities
├── exchanges/      # binance, gateio, okx
├── verifyrun/      # verifier registry
├── verifier/
├── artifacts/      # layout + FetchAndStore
└── …
```

## CLI

One binary, four command groups:

```bash
go run ./cmd/por anchor  -exchange okx
go run ./cmd/por verify  -snapshot PR01JUN26
go run ./cmd/por fetch   gateio
go run ./cmd/por probe   rpc -network BSC -chainlist
go run ./cmd/por exchanges
go run ./cmd/por exchanges okx
```

Details: [por-cli.md](./por-cli.md). Exchange upstream index: [exchange-upstream-registry.md](./exchange-upstream-registry.md).

## Rating and reporting

```
config/rating/
├── stage_requirements.yaml
├── evidence_level.yaml
└── risk_flags.yaml

schemas/
├── artifact-metadata.v1.schema.json
├── exchange-assessment.v1.schema.json
└── deposit-sample-manifest.v1.schema.json

docs/templates/
├── exchange-transparency-report.md
└── exchange-assessment.v1.example.json

docs/reports/
├── exchange-comparison.md
└── artifact-archive-index.md

docs/reports/binance/
├── PR01JUN26-assessment.json
└── PR01JUN26-transparency-report.md

docs/reports/okx/
├── 506872725-assessment.json
└── 506872725-transparency-report.md

docs/reports/bybit/
├── 2025061709-assessment.json
└── 2025061709-transparency-report.md
```

Methodology: [por-transparency-framework.md](./por-transparency-framework.md). Product audience: [ardmere-service-audience.md](./ardmere-service-audience.md).

## Scripts

```
scripts/por.sh              primary wrapper (loads ~/.zshenv + .env)
scripts/_por-env.sh         shared env bootstrap for por wrappers
scripts/por-run.sh          alias → por
scripts/por-batch-recent.sh batch verify + reports (via por.sh)
scripts/por-fetch.sh        alias → por fetch
scripts/gateio/               → por fetch gateio
scripts/okx/                  → por fetch okx
```

## Adding an exchange

Register adapter in `internal/exchangereg/reg.go`, add `por fetch <id>` subcommand in `internal/porfetch`, document in `docs/<id>-por-data-guide.md`.
