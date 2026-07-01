# PoR Snapshot History (Public Evaluation Set)

> Latest snapshot per exchange: [Exchange comparison](./exchange-comparison.md)  
> Methodology: [PoR Stage Framework](../por-transparency-framework.md)

This index lists every assessed snapshot in the ardmere public evaluation set. The [homepage](../../index.html) shows the six most recent rows; the table below is the complete history.

**Assessed:** 2026-07-01 · Schema: `ardmere/exchange-assessment@1`

## All public reports

| Exchange | Snapshot | Date | Stage | Report |
| --- | --- | --- | --- | --- |
| OKX | `508399035` | 2026-06-18 | **Stage 1** | [report](./okx/508399035-transparency-report.md) |
| Binance | `PR01JUN26` | 2026-01-06 | Stage 0 | [report](./binance/PR01JUN26-transparency-report.md) |
| Bitget | `202605` | 2026-05-27 | Stage 0 | [report](./bitget/202605-transparency-report.md) |
| OKX | `506872725` | 2026-05-06 | **Stage 1** | [report](./okx/506872725-transparency-report.md) |
| Binance | `PR01MAY26` | 2026-01-05 | Stage 0 | [report](./binance/PR01MAY26-transparency-report.md) |
| OKX | `507918525` | 2026-04-19 | **Stage 1** | [report](./okx/507918525-transparency-report.md) |
| Binance | `PR01APR26` | 2026-01-04 | Stage 0 | [report](./binance/PR01APR26-transparency-report.md) |
| Gate.io | `20260316` | 2026-03-16 | Stage 0 | [report](./gateio/20260316-transparency-report.md) |
| Bybit | `2025061709` | 2025-06-17 | Stage 0 | [report](./bybit/2025061709-transparency-report.md) |
| HTX | `20230910` | 2023-09-10 | Stage 0 | [report](./htx/20230910-transparency-report.md) |

Dates are snapshot times from each assessment JSON (`snapshotTime`, UTC date).

## OKX (monthly; full public artifact stack)

| Snapshot | Snapshot time | Stage | Effective PoR | Report |
| --- | --- | --- | --- | --- |
| `508399035` | 2026-06-18 | **Stage 1** | Yes | [report](./okx/508399035-transparency-report.md) |
| `506872725` | 2026-05-06 | **Stage 1** | Yes | [report](./okx/506872725-transparency-report.md) |
| `507918525` | 2026-04-19 | **Stage 1** | Yes | [report](./okx/507918525-transparency-report.md) |

## Binance (monthly; BAPI summary only for current snapshot)

| Snapshot | List time (UTC) | Stage | Effective PoR | Notes | Report |
| --- | --- | --- | --- | --- | --- |
| `PR01JUN26` | 01/06/26 | Stage 0 | No | Full BAPI summary + wallet zip | [report](./binance/PR01JUN26-transparency-report.md) |
| `PR01MAY26` | 01/05/26 | Stage 0 | No | Wallet zip + list metadata; BAPI summary not public for historical period | [report](./binance/PR01MAY26-transparency-report.md) |
| `PR01APR26` | 01/04/26 | Stage 0 | No | Wallet zip + list metadata; BAPI summary not public for historical period | [report](./binance/PR01APR26-transparency-report.md) |

## Gate.io, Bybit, Bitget, HTX (single snapshot in set)

Public APIs, WAF, login walls, or fixture-only pipelines limit these venues to **one** assessed snapshot each until additional periods are captured manually.

| Exchange | Snapshot | Stage | Report |
| --- | --- | --- | --- |
| Gate.io | `20260316` | Stage 0 | [report](./gateio/20260316-transparency-report.md) |
| Bybit | `2025061709` | Stage 0 | [report](./bybit/2025061709-transparency-report.md) |
| Bitget | `202605` | Stage 0 | [report](./bitget/202605-transparency-report.md) |
| HTX | `20230910` | Stage 0 | [report](./htx/20230910-transparency-report.md) |

## Regenerating this set

```bash
./scripts/por-batch-recent.sh 3 okx,binance
./scripts/build-public-html.sh
```

Use `go run ./cmd/por batch -help` for flags (`-skip-rpc`, `-full-rpc-latest-only`, etc.).
