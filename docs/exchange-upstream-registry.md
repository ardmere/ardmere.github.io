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

## Data guides

| Exchange | Guide |
| --- | --- |
| Binance | [binance-por-data-guide.md](./binance-por-data-guide.md) |
| OKX | [okx-por-data-guide.md](./okx-por-data-guide.md) |
| Gate.io | [gate-por-data-guide.md](./gate-por-data-guide.md) |
| Bybit | [bybit-por-data-guide.md](./bybit-por-data-guide.md) |
| Bitget | [bitget-por-data-guide.md](./bitget-por-data-guide.md) |
| HTX | [htx-por-data-guide.md](./htx-por-data-guide.md) |

## Maintenance

1. Add or update entries in `config/exchanges/registry.yaml`.
2. Run `por exchanges` to verify load.
3. Sync this index table if roles or URLs changed materially.
4. Add `references` URLs to new `*-assessment.json` snapshots.
5. Link from the exchange `*-por-data-guide.md` header (do not duplicate full URL lists in guides).
