# Exchange PoR Transparency Comparison

> Methodology: [PoR Stage Framework](../por-transparency-framework.md)  
> Policy standard: [Effective PoR Standard](../effective-por-standard.md)  
> **All 50 public reports:** [Snapshot history](./snapshot-history.md)

## Summary (latest snapshot per exchange)

| Exchange | Snapshot | PoR Stage | Gen / Evidence | Effective PoR | Confidence | Last Verified Snapshot | Upstream verifier | Report |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Binance | `PR04NOV26` | `Stage 0` | `Gen 2 / E2` | `false` | `medium` | `2026-11-04T00:00:00Z` | [binance/zkmerkle-proof-of-solvency](https://github.com/binance/zkmerkle-proof-of-solvency) | [report](./binance/PR04NOV26-transparency-report.md) |
| OKX | `508756245` | `Stage 1` | `Gen 2 / E2` | `true` | `medium` | `2026-10-31T16:00:00Z` | [okx/proof-of-reserves-v2](https://github.com/okx/proof-of-reserves-v2) (+ [v1 sigs](https://github.com/okx/proof-of-reserves)) | [report](./okx/508756245-transparency-report.md) |
| Gate.io | `20261018` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-18T00:00:00Z` | [gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves) | [report](./gateio/20261018-transparency-report.md) |
| Bybit | `2026102706` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-27T06:00:00Z` | [bybit-exchange/merkle-proof](https://github.com/bybit-exchange/merkle-proof) | [report](./bybit/2026102706-transparency-report.md) |
| Bitget | `20261021` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-21T04:00:00Z` | [BitgetLimited/proof-of-reserves](https://github.com/BitgetLimited/proof-of-reserves) | [report](./bitget/20261021-transparency-report.md) |
| HTX | `20261001` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-01T00:00:00Z` | [huobiapi/Tool-Go-MerkleVerify](https://github.com/huobiapi/Tool-Go-MerkleVerify) | [report](./htx/20261001-transparency-report.md) |
| MEXC | `202610` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-15T00:00:00Z` | [mexcmerkle/mexc-merkle-proof-of-reserves](https://github.com/mexcmerkle/mexc-merkle-proof-of-reserves) | [report](./mexc/202610-transparency-report.md) |
| KuCoin | `20261031` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-31T15:59:59Z` | — | [report](./kucoin/20261031-transparency-report.md) |
| BingX | `20261015` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-15T00:00:00Z` | [BingX-API/BingX-proof-of-reserves](https://github.com/BingX-API/BingX-proof-of-reserves) | [report](./bingx/20261015-transparency-report.md) |
| Poloniex | `20261001` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-01T00:00:00Z` | [poloniex/tools-nodejs-address-verify](https://github.com/poloniex/tools-nodejs-address-verify) | [report](./poloniex/20261001-transparency-report.md) |
| Ourbit | `20261101` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-10-31T16:00:00Z` | [OurbitOfficial/merkle_validate](https://github.com/OurbitOfficial/merkle_validate) | [report](./ourbit/20261101-transparency-report.md) |
| CoinEx | `20260914` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-09-14T07:00:00Z` | [coinexcom/proof-of-reserves](https://github.com/coinexcom/proof-of-reserves) | [report](./coinex/20260914-transparency-report.md) |
| Bitunix | `20261103` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2026-11-03T05:42:40Z` | [BitunixOfficial/proof-of-reserves](https://github.com/BitunixOfficial/proof-of-reserves) | [report](./bitunix/20261103-transparency-report.md) |

## By exchange

- **Binance (`Stage 0`)** — Latest `PR04NOV26` (Nov 2026, period 49): Binance snapshot PR04NOV26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts. [Full report](./binance/PR04NOV26-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **OKX (`Stage 1`)** — Latest `508756245` (Oct 2026): OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit 508756245. ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1. [Full report](./okx/508756245-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **Gate.io (`Stage 0`)** — Latest `20261018` (Oct 2026): gateio snapshot 20261018 remains Stage 0 in the ardmere public evaluation set. [Full report](./gateio/20261018-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **Bybit (`Stage 0`)** — Latest `2026102706` (Oct 2026): bybit snapshot 2026102706 remains Stage 0 in the ardmere public evaluation set. [Full report](./bybit/2026102706-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **Bitget (`Stage 0`)** — Latest `20261021` (Oct 2026): bitget snapshot 20261021 remains Stage 0 in the ardmere public evaluation set. [Full report](./bitget/20261021-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **HTX (`Stage 0`)** — Latest `20261001` (Oct 2026): htx snapshot 20261001 remains Stage 0 in the ardmere public evaluation set. [Full report](./htx/20261001-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **MEXC (`Stage 0`)** — Latest `202610` (Oct 2026): MEXC snapshot 202610 is a weak Stage 0 (press-release / import summary): self-reported ratios only; no public period merkle root in the archived package; user Merkle remains login-gated. [Full report](./mexc/202610-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **KuCoin (`Stage 0`)** — Latest `20261031` (Oct 2026): kucoin snapshot 20261031 remains Stage 0 in the ardmere public evaluation set. [Full report](./kucoin/20261031-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **BingX (`Stage 0`)** — Latest `20261015` (Oct 2026): bingx snapshot 20261015 remains Stage 0 in the ardmere public evaluation set. [Full report](./bingx/20261015-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **Poloniex (`Stage 0`)** — Latest `20261001` (Oct 2026): poloniex snapshot 20261001 remains Stage 0 in the ardmere public evaluation set. [Full report](./poloniex/20261001-transparency-report.md) · [history](./snapshot-history.md#all-public-reports)
- **Ourbit (`Stage 0`)** — Latest `20261101` (Oct 2026): Ourbit snapshot 20261101 is Stage 0: public stock_info ratios + live multi-chain wallet table; no snapshot-bound walletZip/ownership; user Merkle login-gated; hash truncated to 128 bits. [Full report](./ourbit/20261101-transparency-report.md)
- **CoinEx (`Stage 0`)** — Latest `20260914` (Sep 2026): CoinEx snapshot 20260914 is Stage 0: public reserve ratios + period merkle root + XLSX address list; user Merkle login-gated; no ownership proofs or global ZK. [Full report](./coinex/20260914-transparency-report.md)
- **Bitunix (`Stage 0`)** — Latest `20261103` (Nov 2026): Bitunix snapshot 20261103 is Stage 0: public BTC/ETH/USDT reserve ratios + period merkle root; user merkle_tree_bitunix.json login-gated; no ownership proofs or global ZK. [Full report](./bitunix/20261103-transparency-report.md)

Gap detail, risk flags, and verifier output are in each report. Missing data is marked `UNVERIFIABLE`, not `PASS` ([methodology §1.3](../por-transparency-framework.md#gap-reporting-discipline)).

Canonical upstream registry (PoR pages + GitHub repos): [`config/exchanges/registry.yaml`](../../config/exchanges/registry.yaml) · CLI: `por exchanges`.

## Evidence archive

Assessment JSON, bundle roots, and hashes for this report set: [Artifact archive index](./artifact-archive-index.md).
