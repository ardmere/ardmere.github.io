# Exchange upstream registry (index)

> **Canonical source**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml)  
> **CLI**: `por exchanges` · `por exchanges <id>`

Human-readable index of official PoR landing pages and GitHub repositories for exchanges integrated in ardmere. Edit the YAML file first; keep this page in sync when adding exchanges or upstream repos.

## Role vocabulary

| Role | Meaning |
| --- | --- |
| `verifier_toolchain` | Open-source verifier / circuit / CLI |
| `global_zk_verifier_spec` | Global ZK proof specification + implementation |
| `address_ownership_verifier` | Wallet control signature verification |
| `user_merkle_verifier` | User inclusion Merkle path verification (Gen 1) |
| `sample_artifact_release` | Sample period bundles on GitHub Releases only |

## Integrated exchanges

| ID | VC | PoR page | Upstream GitHub | ardmere use |
| --- | --- | --- | --- | --- |
| `binance` | VC3 | [proof-of-reserves](https://www.binance.com/en/proof-of-reserves) | [binance/zkmerkle-proof-of-solvency](https://github.com/binance/zkmerkle-proof-of-solvency) | `global-zk-proof@0` (stub) |
| `okx` | VC2 | [proof-of-reserves](https://www.okx.com/proof-of-reserves) | [okx/proof-of-reserves](https://github.com/okx/proof-of-reserves) · [okx/proof-of-reserves-v2](https://github.com/okx/proof-of-reserves-v2) | `address-ownership@okx-1` · `global-zk-proof@okx-1` |
| `gateio` | VC1 | [proof-of-reserves](https://www.gate.com/proof-of-reserves) | [gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves) | `global-zk-proof@gateio-1` (planned) |
| `bybit` | VC1 | [proof-of-reserves](https://www.bybit.com/en/proof-of-reserves) | [bybit-exchange/merkle-proof](https://github.com/bybit-exchange/merkle-proof) | `user-merkle-proof@bybit-1` |
| `bitget` | VC1 | [proof-of-reserves](https://www.bitget.com/proof-of-reserves) | [BitgetLimited/proof-of-reserves](https://github.com/BitgetLimited/proof-of-reserves) | `user-merkle-proof@bitget-1` |
| `htx` | VC1 | [merkle](https://www.htx.com/zh-cn/finance/merkle/) | [huobiapi/Tool-Go-MerkleVerify](https://github.com/huobiapi/Tool-Go-MerkleVerify) | `global-zk-proof@htx-1` |
| `mexc` | VC1 | [verify user assets](https://www.mexc.com/support/article/how-to-verify-user-assets-on-mexc-17827791528836) | [mexcmerkle/mexc-merkle-proof-of-reserves](https://github.com/mexcmerkle/mexc-merkle-proof-of-reserves) | `user-merkle-proof@mexc-1` |
| `kucoin` | VC1 | [proof-of-reserves](https://www.kucoin.com/proof-of-reserves) | — (no official OSS verifier found) | `user-merkle-proof@0` (stub) |
| `bingx` | VC1 | [balance-reserve](https://bingx.com/en/balance-reserve) | [BingX-API/BingX-proof-of-reserves](https://github.com/BingX-API/BingX-proof-of-reserves) | `user-merkle-proof@bingx-1` |
| `poloniex` | VC1 | [proof-of-reserves](https://www.poloniex.com/topic/security/proof-of-reserves) | [poloniex/tools-nodejs-address-verify](https://github.com/poloniex/tools-nodejs-address-verify) | `address-ownership@0` (CSV channel blocked) |
| `ourbit` | VC1 | [reserve-proof](https://www.ourbit.com/assets/reserve-proof) | [OurbitOfficial/merkle_validate](https://github.com/OurbitOfficial/merkle_validate) | `user-merkle-proof@ourbit-1` |
| `coinex` | VC1 | [reserve-proof](https://www.coinex.com/en/reserve-proof) | [coinexcom/proof-of-reserves](https://github.com/coinexcom/proof-of-reserves) | `user-merkle-proof@coinex-1` |
| `bitunix` | VC1 | [proof-of-reserve](https://www.bitunix.com/proof-of-reserve) | [BitunixOfficial/proof-of-reserves](https://github.com/BitunixOfficial/proof-of-reserves) | `user-merkle-proof@bitunix-1` |
| `woox` | VC1 | [proof-of-reserves](https://woox.io/proof-of-reserves) | [woo-x/merkle-tree](https://github.com/woo-x/merkle-tree) | `user-merkle-proof@0` (stub; tree+login hash) |
| `toobit` | VC1 | [proof-of-reserves](https://www.toobit.com/proof-of-reserves) | [toobit-docs/por-verifier](https://github.com/toobit-docs/por-verifier) | `user-merkle-proof@toobit-1` |
| `phemex` | VC1 | [proof-of-reserves](https://phemex.com/proof-of-reserves) | — (no official OSS verifier found) | `user-merkle-proof@0` (stub) |
| `btcc` | VC1 | [proof-of-reserves](https://www.btcc.com/en-US/proof-of-reserves) | — (no official OSS verifier found) | `user-merkle-proof@0` (stub) |

## Data guides

| Exchange | Guide |
| --- | --- |
| Binance | [binance-por-data-guide.md](./binance-por-data-guide.md) |
| OKX | [okx-por-data-guide.md](./okx-por-data-guide.md) |
| Gate.io | [gate-por-data-guide.md](./gate-por-data-guide.md) |
| Bybit | [bybit-por-data-guide.md](./bybit-por-data-guide.md) |
| Bitget | [bitget-por-data-guide.md](./bitget-por-data-guide.md) |
| HTX | [htx-por-data-guide.md](./htx-por-data-guide.md) |
| MEXC | [mexc-por-data-guide.md](./mexc-por-data-guide.md) |
| KuCoin | [kucoin-por-data-guide.md](./kucoin-por-data-guide.md) |
| BingX | [bingx-por-data-guide.md](./bingx-por-data-guide.md) |
| Poloniex | [poloniex-por-data-guide.md](./poloniex-por-data-guide.md) |
| Ourbit | [ourbit-por-data-guide.md](./ourbit-por-data-guide.md) |
| CoinEx | [coinex-por-data-guide.md](./coinex-por-data-guide.md) |
| Bitunix | [bitunix-por-data-guide.md](./bitunix-por-data-guide.md) |
| WOO X | [woox-por-data-guide.md](./woox-por-data-guide.md) |
| Toobit | [toobit-por-data-guide.md](./toobit-por-data-guide.md) |
| Phemex | [phemex-por-data-guide.md](./phemex-por-data-guide.md) |
| BTCC | [btcc-por-data-guide.md](./btcc-por-data-guide.md) |

## Maintenance

1. Add or update entries in `config/exchanges/registry.yaml`.
2. Run `por exchanges` to verify load.
3. Sync this index table if roles or URLs changed materially.
4. Add `references` URLs to new `*-assessment.json` snapshots.
5. Link from the exchange `*-por-data-guide.md` header (do not duplicate full URL lists in guides).
