# Exchange PoR Transparency Comparison

> Methodology: [PoR Stage Framework](../por-transparency-framework.md)  
> Policy standard: [Effective PoR Standard](../effective-por-standard.md)  
> **All 27 public reports:** [Snapshot history](./snapshot-history.md)

## Summary (latest snapshot per exchange)

| Exchange | Snapshot | PoR Stage | Gen / Evidence | Effective PoR | Confidence | Last Verified Snapshot | Upstream verifier | Report |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Binance | `PR01OCT26` | `Stage 0` | `Gen 2 / E2` | `false` | `medium` | `2026-10-01T00:00:00Z` | [binance/zkmerkle-proof-of-solvency](https://github.com/binance/zkmerkle-proof-of-solvency) | [report](./binance/PR01OCT26-transparency-report.md) |
| OKX | `502260945` | `Stage 1` | `Gen 2 / E2` | `true` | `medium` | `2026-10-02T16:00:00Z` | [okx/proof-of-reserves-v2](https://github.com/okx/proof-of-reserves-v2) (+ [v1 sigs](https://github.com/okx/proof-of-reserves)) | [report](./okx/502260945-transparency-report.md) |
| Gate.io | `20260618` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-06-18T00:00:00Z` | [gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves) | [report](./gateio/20260618-transparency-report.md) |
| Bybit | `2026052709` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-05-27T09:00:00Z` | [bybit-exchange/merkle-proof](https://github.com/bybit-exchange/merkle-proof) | [report](./bybit/2026052709-transparency-report.md) |
| Bitget | `202606` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-06-18T10:00:00Z` | [BitgetLimited/proof-of-reserves](https://github.com/BitgetLimited/proof-of-reserves) | [report](./bitget/202606-transparency-report.md) |
| HTX | `20230910` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2023-09-10T18:10:33Z` | [huobiapi/Tool-Go-MerkleVerify](https://github.com/huobiapi/Tool-Go-MerkleVerify) | [report](./htx/20230910-transparency-report.md) |
| MEXC | `202610` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-15T00:00:00Z` | [mexcmerkle/mexc-merkle-proof-of-reserves](https://github.com/mexcmerkle/mexc-merkle-proof-of-reserves) | [report](./mexc/202610-transparency-report.md) |
| KuCoin | `20260930` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-09-30T15:59:59Z` | — | [report](./kucoin/20260930-transparency-report.md) |

## By exchange

- **Binance (`Stage 0`)** — Latest `PR01OCT26` (Oct 2026, period 47): Binance snapshot PR01OCT26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts. [Full report](./binance/PR01OCT26-transparency-report.md) · [history](./snapshot-history.md#binance-monthly-bapi-summary-only-for-current-snapshot)
- **OKX (`Stage 1`)** — Latest `502260945` (Oct 2026): OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit 502260945. ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1. [Full report](./okx/502260945-transparency-report.md) · [history](./snapshot-history.md#okx-monthly-full-public-artifact-stack)
- **Gate.io (`Stage 0`)** — Latest `20260618` (Jun 2026): gateio snapshot 20260618 remains Stage 0 in the ardmere public evaluation set. [Full report](./gateio/20260618-transparency-report.md)
- **Bybit (`Stage 0`)** — Latest `2026052709` (May 2026): bybit snapshot 2026052709 remains Stage 0 in the ardmere public evaluation set. [Full report](./bybit/2026052709-transparency-report.md)
- **Bitget (`Stage 0`)** — Latest `202606` (Jun 2026): bitget snapshot 202606 remains Stage 0 in the ardmere public evaluation set. [Full report](./bitget/202606-transparency-report.md)
- **HTX (`Stage 0`)** — Latest `20230910` (Sep 2023): htx snapshot 20230910 remains Stage 0 in the ardmere public evaluation set. [Full report](./htx/20230910-transparency-report.md)
- **MEXC (`Stage 0`)** — Latest `202610` (Oct 2026): MEXC snapshot 202610 is a weak Stage 0 (press-release / import summary): self-reported ratios only; no public period merkle root in the archived package; user Merkle remains login-gated. [Full report](./mexc/202610-transparency-report.md)
- **KuCoin (`Stage 0`)** — Latest `20260930` (Sep 2026): kucoin snapshot 20260930 remains Stage 0 in the ardmere public evaluation set. [Full report](./kucoin/20260930-transparency-report.md)

Gap detail, risk flags, and verifier output are in each report. Missing data is marked `UNVERIFIABLE`, not `PASS` ([methodology §1.3](../por-transparency-framework.md#gap-reporting-discipline)).

Canonical upstream registry (PoR pages + GitHub repos): [`config/exchanges/registry.yaml`](../../config/exchanges/registry.yaml) · CLI: `por exchanges`.

## Evidence archive

Assessment JSON, bundle roots, and hashes for this report set: [Artifact archive index](./artifact-archive-index.md).
