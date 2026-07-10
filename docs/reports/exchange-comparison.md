# Exchange PoR Transparency Comparison

> Methodology: [PoR Stage Framework](../por-transparency-framework.md)  
> Policy standard: [Effective PoR Standard](../effective-por-standard.md)  
> **All 10 public reports:** [Snapshot history](./snapshot-history.md)

## Summary (latest snapshot per exchange)

| Exchange | Snapshot | PoR Stage | Gen / Evidence | Effective PoR | Confidence | Last Verified Snapshot | Upstream verifier | Report |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Binance | `PR01JUL26` | `Stage 0` | `Gen 2 / E2` | `false` | `medium` | `2026-07-01T00:00:00Z` | [binance/zkmerkle-proof-of-solvency](https://github.com/binance/zkmerkle-proof-of-solvency) | [report](./binance/PR01JUL26-transparency-report.md) |
| OKX | `508399035` | `Stage 1` | `Gen 2 / E2` | `true` | `medium` | `2026-06-18T16:00:00Z` | [okx/proof-of-reserves-v2](https://github.com/okx/proof-of-reserves-v2) (+ [v1 sigs](https://github.com/okx/proof-of-reserves)) | [report](./okx/508399035-transparency-report.md) |
| Gate.io | `20260618` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-06-18T00:00:00Z` | [gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves) | [report](./gateio/20260618-transparency-report.md) |
| Bitget | `202606` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-06-18T10:00:00Z` | [BitgetLimited/proof-of-reserves](https://github.com/BitgetLimited/proof-of-reserves) | [report](./bitget/202606-transparency-report.md) |
| Bybit | `2026052709` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-05-27T09:00:00Z` | [bybit-exchange/merkle-proof](https://github.com/bybit-exchange/merkle-proof) | [report](./bybit/2026052709-transparency-report.md) |
| HTX | `20230910` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2023-09-10T18:10:33Z` | [huobiapi/Tool-Go-MerkleVerify](https://github.com/huobiapi/Tool-Go-MerkleVerify) | [report](./htx/20230910-transparency-report.md) |

## By exchange

- **Binance (`Stage 0`)** — Latest `PR01JUL26` (Jul 2026, period 44): public summary and wallet list; missing ownership proof, public global proof, trusted-setup transparency. [Full report](./binance/PR01JUL26-transparency-report.md) · [history](./snapshot-history.md#binance-monthly-bapi-summary-only-for-current-snapshot)
- **OKX (`Stage 1`)** — Latest `508399035` (Jun 2026): public reserves, ownership proof, global zk proof. [Full report](./okx/508399035-transparency-report.md) · [history](./snapshot-history.md#okx-monthly-full-public-artifact-stack)
- **Gate.io (`Stage 0`)** — Latest `20260618`: public reserve-ratio summary (page-scrape import); login-gated zk and no public wallet artifacts. [Full report](./gateio/20260618-transparency-report.md)
- **Bitget (`Stage 0`)** — Latest `202606`: announcement-scrape summary; merkle root and user proof still need browser capture. [Full report](./bitget/202606-transparency-report.md)
- **Bybit (`Stage 0`)** — Latest `2026052709` (period 36, May 27 2026): press-release scrape; merkle root needs browser API capture. [Full report](./bybit/2026052709-transparency-report.md)
- **HTX (`Stage 0`)** — Still on archived zk sample `20230910` (2023-09); June 2026 bundle requires manual login download. [Full report](./htx/20230910-transparency-report.md)

Gap detail, risk flags, and verifier output are in each report. Missing data is marked `UNVERIFIABLE`, not `PASS` ([methodology §1.3](../por-transparency-framework.md#gap-reporting-discipline)).

Canonical upstream registry (PoR pages + GitHub repos): [`config/exchanges/registry.yaml`](../../config/exchanges/registry.yaml) · CLI: `por exchanges`.

## Evidence archive

Assessment JSON, bundle roots, and hashes for this report set: [Artifact archive index](./artifact-archive-index.md).
