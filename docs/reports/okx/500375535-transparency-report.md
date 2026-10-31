# OKX 500375535 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [500375535-assessment.json](./500375535-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `500375535` |
| Snapshot time | `2026-08-10T16:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

OKX audit 500375535 does not meet all Stage 1 public artifact thresholds in this pipeline.

Missing or failing verifiers block Stage 1 effective PoR for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `canonical on-chain/DA anchor` | `NO_CANONICAL_ANCHOR` | Stage 1 | No official exchange canonical on-chain/DA anchor established. |
| `weekly full PoR / daily anchor` | `HIGH_FREQUENCY_GAP` | Stage 1 | Monthly cadence does not satisfy Stage 2 frequency expectations. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-02T16:00:00Z` |
| Previous snapshot | `2026-06-18T16:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `6 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Older snapshot 502260945 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `a52bb5a36482ee00ac08ae04d3c73b38234d7fe2907fac418ed91d9c7910e291` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `7db12d643664e0b84ff1e54ffdba339d47298a249afe494c7e770d63dc117462` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026081100_V6.zip |
| `globalProofBundle` | `c38e880b939d664e91f4e5f14b6c307295b4920156b65f3185ba1b5b6460c63e` | https://static.okx.com/cdn/okx/por/merkel/por_500375535_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x8c29f98a96bb8249a39271bcee97f8a6b8372c17464fd5df343fd553dc56b3c0` |
| Verification bundle root | `0xa702d9edd173c41208a96a92efc2e28fa97ba0423bb372af024300c1fcf01355` |
| Artifact bundle SHA-256 | `66a673b2404118aca399869e8a914a0e1b22240fea324c9fb121900d29443b68` |
| Verification bundle SHA-256 | `f875d24e1e623008b2790a685924ce767e3ea37e55aec311aee8098ddbd7faf6` |

Local bundle paths: [500375535.artifact-bundle.json](../../../artifacts/okx/500375535/bundles/500375535.artifact-bundle.json), [500375535.verification-bundle.v2.json](../../../artifacts/okx/500375535/bundles/500375535.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | okx summary has no reserve currencies |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `FAIL` | `0.9983` | 351/212313 address signatures failed |
| `onchain-balance-hot` | `2.1` | `PARTIAL` | `0.0002` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.0030` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.0013` |  |
| `global-zk-proof` | `okx-1` | `PASS` | `1.0000` | zkSTARKValidator verify-global succeeded; summary merkle root bound to zk proof |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `500375535.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `btc-anchor` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `solvency-claim` (`FAIL`)

Summary: okx summary has no reserve currencies

#### `address-ownership` (`FAIL`)

Summary: 351/212313 address signatures failed

Finding counts: FAIL 5

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x80299387783b6b8c0645a17ea72fdc3cd32e0c9367445fdbe7344742af521fc39afc1fbc3917efc3a4bd72c48814e9b0` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x8063539f96661a8b8b86222b7170ee312c0ecc523b89e3f938698f69d0f82908d3db19259e047afc5c780b12eda785bc` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x806e4c086bc3aff5b9940868fac93e11b139aa04cf749a1adf13e8c6058e861df0426488db4c04d249d26d7d55e9a650` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x80a935b3b8a82f6a032002dd04b9cac3f0844c9b73b3513067f97b1141e6c0b5b89ed36d0a4c74c59d98498220805c27` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x80b57d57a80e34e207622869dcb263220623483c61392b95c762f40e073f52cf207808635f1824d69338275927f55509` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |

#### `onchain-balance-hot` (`PARTIAL`)

Finding counts: UNVERIFIABLE 21, PASS 2

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `ASSET-HUB\|DOT` | 56 | no native verifier for this (coin,network) pair yet |
| `A\|VAULTA` | 2 | no native verifier for this (coin,network) pair yet |
| `DOTK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ELF\|AELF` | 273 | no native verifier for this (coin,network) pair yet |
| `ETCK-KIP20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ETC\|ETC` | 40 | no native verifier for this (coin,network) pair yet |
| `ETHK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `FILK-OKC20\|OKC` | 8 | no native verifier for this (coin,network) pair yet |
| `LINKK-OKC20\|OKC` | 3 | no native verifier for this (coin,network) pair yet |
| `OKB-OKC20\|OKC` | 4 | no native verifier for this (coin,network) pair yet |
| `TRXK-KIP20\|OKC` | 13 | no native verifier for this (coin,network) pair yet |
| `UNIK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `USDC-OKC20\|OKC` | 244 | no native verifier for this (coin,network) pair yet |
| `USDT-OKC20\|OKC` | 105 | no native verifier for this (coin,network) pair yet |
| `XRPK-KIP20\|OKC` | 5 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x85e306287d845752a1716bfd5e90f266d08377dd1d5e4189d6f0bc05a9860dac6a0f73d2f6ec7b1458681ba0570afac3` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xb134297b6c0f624afab8d357c9c2f8324dfa844466c35a12d33ebe4959dfe06d88f9544593cf17996bbc0cca4d0a43d9` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x80670ceb430da829ff33566afdecede98d4bcdb1bdbc33acb16425e1cdd29a8088f29bcd7c5ba104171eaf3b1b369291` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xa2b4a69dc29719fcf4c400c5ba4c60966f216fafea399cfc2024cf07ad60d9258291d5f83b8094225bbe86f646a67f1b` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x8a7fa64cd9d0fcbc52f417d7b05271b151219659da4ecb306385dfb06d049af43783c4662e22e48175cedfc8e8438bf8` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xb2d016a313305400b7125cc0c5bd9a6acb001cd640648084d63d49154451a700feca775c17c09468638e517a4ceb42f8` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |

#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 163, WARN 173, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `PEOPLE@ETH#25725648` | 1133754327.14 | 0 | on-chain balance != csv claim by 1133754327.14 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `PEOPLE@ETH#25725648` | 111957008.5713641423 | 0 | on-chain balance != csv claim by 111957008.5713641423 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `PEOPLE@ETH#25725648` | 112213388.2783534113 | 0 | on-chain balance != csv claim by 112213388.2783534113 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xb0a27099582833c0cb8c7a0565759ff145113d64` | `PEOPLE@ETH#25725648` | 13481229.915432 | 0 | on-chain balance != csv claim by 13481229.915432 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x85dcd76d4fbd3aa0c85c27b9441222c19a14134b` | `PEOPLE@ETH#25725648` | 8341552.467069 | 0 | on-chain balance != csv claim by 8341552.467069 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `PEOPLE@ETH#25725648` | 5929789.9465202031 | 0 | on-chain balance != csv claim by 5929789.9465202031 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0xafee421482faea92292ed3ffe29371742542ad72` | `USDT-ARBITRUM@ARBITRUM#493128998` | 6015680.143667 | 0 | on-chain balance != csv claim by 6015680.143667 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xb5216cb558cb018583bed009ee25ca73eb27bb1d` | `USDT-OPTIMISM@OPTIMISM#155389411` | 5047539.797394 | 0 | on-chain balance != csv claim by 5047539.797394 (provider=https://mainnet.optimism.io) |
| `0xd8a3ed83a59d3e10d82a4b845d64b24ba0572baf` | `USDT-ARBITRUM@ARBITRUM#493128998` | 3999114.809395 | 0 | on-chain balance != csv claim by 3999114.809395 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x92ea7496eba5f001d620005f88f3e8e686e3d4ea` | `PEOPLE@ETH#25725648` | 2454976.41 | 0 | on-chain balance != csv claim by 2454976.41 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x343d752bb710c5575e417edb3f9fa06241a4749a` | `POLY-USDC@MATIC#91781706` | 2246620.224877 | 46.514712 | on-chain balance != csv claim by 2246573.710165 (provider=https://polygon.gateway.tenderly.co) |
| `0xde3c8d726e906b379a020cdb04bbe4b8f39a726f` | `USDC-AVAXC@AVAXC#92491108` | 3880903.783612 | 0 | on-chain balance != csv claim by 3880903.783612 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x4a1ca3e90e9e3a00eb27e9f409b4db310aa7405f` | `USDT-ARBITRUM@ARBITRUM#493128998` | 1394907.78412 | 0 | on-chain balance != csv claim by 1394907.78412 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xe7c5e276b2434743f6c79365d967cbe4291477bc` | `USDT-ARBITRUM@ARBITRUM#493128998` | 1196452.048765 | 0 | on-chain balance != csv claim by 1196452.048765 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `PEOPLE@ETH#25725648` | 862400.3152159399 | 0 | on-chain balance != csv claim by 862400.3152159399 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x3f03440b19d15ff2068daf68e57eb656bfd0a67e` | `USDC-AVAXC@AVAXC#92491108` | 1412091.176149 | 0 | on-chain balance != csv claim by 1412091.176149 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `PEOPLE@ETH#25725648` | 518806.37407388 | 0 | on-chain balance != csv claim by 518806.37407388 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xa0c9b4ce1785953abdceae569789aa158a01da4c` | `USDT-ARBITRUM@ARBITRUM#493128998` | 584029.984272 | 0 | on-chain balance != csv claim by 584029.984272 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x9098575016440086c694c85c5010a0ed85249675` | `USDT-OPTIMISM@OPTIMISM#155389411` | 482027.402086 | 0 | on-chain balance != csv claim by 482027.402086 (provider=https://mainnet.optimism.io) |
| `0x9d1964437bb6d849c66ab2a2a033340dc2d999dc` | `USDC-AVAXC@AVAXC#92491108` | 717519.967569 | 0 | on-chain balance != csv claim by 717519.967569 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xc74e4c556a16390165c99b40aabc39a87c358305` | `POLY-USDC@MATIC#91781706` | 296744.92669 | 0 | on-chain balance != csv claim by 296744.92669 (provider=https://polygon.gateway.tenderly.co) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `PEOPLE@ETH#25725648` | 302640.3289310999 | 0 | on-chain balance != csv claim by 302640.3289310999 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4cafbe2b5bd89f0c0b7d6fef75b47e7c2ae82f22` | `USDC-AVAXC@AVAXC#92491108` | 404058.620078 | 0 | on-chain balance != csv claim by 404058.620078 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf418263a6468c80d3bcb71033c56aacf81b241ba` | `PEOPLE@ETH#25725648` | 213037.2 | 0 | on-chain balance != csv claim by 213037.2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a8c00e80e3cfae4c903f9462b4ee2826404fe05` | `USDT-OPTIMISM@OPTIMISM#155389411` | 190453.953202 | 0 | on-chain balance != csv claim by 190453.953202 (provider=https://mainnet.optimism.io) |
| `0x76705493d53beefadffc253acb8aef72b7e7f6c6` | `USDT-OPTIMISM@OPTIMISM#155389411` | 177718.422256 | 0 | on-chain balance != csv claim by 177718.422256 (provider=https://mainnet.optimism.io) |
| `0xf1e7b58ccd200fcb9400ddbc30867c43b7dc9ee1` | `USDT-ARBITRUM@ARBITRUM#493128998` | 190453.753212 | 0 | on-chain balance != csv claim by 190453.753212 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xb1da6e546bd592218d6e849a8177d9d4fbf12986` | `POLY-USDC@MATIC#91781706` | 146792.377334 | 0 | on-chain balance != csv claim by 146792.377334 (provider=https://polygon.gateway.tenderly.co) |
| `0x229b1c70124141ec96945aec2323a22b68bb093c` | `USDT-ARBITRUM@ARBITRUM#493128998` | 160453.675941 | 0 | on-chain balance != csv claim by 160453.675941 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `PEOPLE@ETH#25725648` | 137719.01 | 0 | on-chain balance != csv claim by 137719.01 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x04e9d79a17f4894c76c65c52516f790531ef6bee` | `USDT-ARBITRUM@ARBITRUM#493128998` | 121228.731027 | 0 | on-chain balance != csv claim by 121228.731027 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `OKB@ETH#25725648` | 78512.7414021786 | 0 | on-chain balance != csv claim by 78512.7414021786 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x5ca1063894e5bbd688830d83fefa5ce3d339744d` | `USDT-ARBITRUM@ARBITRUM#493128998` | 69336.14213 | 0 | on-chain balance != csv claim by 69336.14213 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `PEOPLE@ETH#25725648` | 65075.91 | 0 | on-chain balance != csv claim by 65075.91 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x99860ec99aae9c0aa5ab763ee73da1217a9da665` | `POLY-USDC@MATIC#91781706` | 64739.147268 | 0 | on-chain balance != csv claim by 64739.147268 (provider=https://polygon.gateway.tenderly.co) |
| `0xbe9ac2d59743d69c4054bab09fa9821929d0f31b` | `USDT-ARBITRUM@ARBITRUM#493128998` | 62995.895889 | 0 | on-chain balance != csv claim by 62995.895889 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x2fb2841c747621a7187b6f2a5390ae77321e5257` | `USDC-AVAXC@AVAXC#92491108` | 74234.314579 | 0 | on-chain balance != csv claim by 74234.314579 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x3b50287cf8f073ff443b84de57e2375d446b962b` | `POLY-USDC@MATIC#91781706` | 51160.610126 | 0 | on-chain balance != csv claim by 51160.610126 (provider=https://polygon.gateway.tenderly.co) |
| `0x582988097ef9f2cf50cfac3f5f3b1ed54bc52a88` | `USDT-ARBITRUM@ARBITRUM#493128998` | 49669.675025 | 0 | on-chain balance != csv claim by 49669.675025 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x9c308b1665097af7ab2e38a41577bfecdc2d1c5c` | `PEOPLE@ETH#25725648` | 46334.55532091 | 0 | on-chain balance != csv claim by 46334.55532091 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `PEOPLE@ETH#25725648` | 46289.52 | 0 | on-chain balance != csv claim by 46289.52 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a152242605ab21df5ac6e819fdb128a4efc33fc` | `POLY-USDC@MATIC#91781706` | 45772.035132 | 0 | on-chain balance != csv claim by 45772.035132 (provider=https://polygon.gateway.tenderly.co) |
| `0xb57198a6e5af11544936a520a54268bc6ababb0c` | `USDT-OPTIMISM@OPTIMISM#155389411` | 43982.29452 | 0 | on-chain balance != csv claim by 43982.29452 (provider=https://mainnet.optimism.io) |
| `0x04ba1b380b1418007c1407bbe12ac21a596ae1e0` | `USDC-AVAXC@AVAXC#92491108` | 49442.470531 | 0 | on-chain balance != csv claim by 49442.470531 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x7d328f5d5f4640f81b2f93d2903ec0a6c030eed6` | `PEOPLE@ETH#25725648` | 40963.67003613 | 0 | on-chain balance != csv claim by 40963.67003613 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `OKB@ETH#25725648` | 30822.1988656955 | 0 | on-chain balance != csv claim by 30822.1988656955 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6bea9e00a6bb63d22904d975bdf8a98f26f9f49d` | `USDC-AVAXC@AVAXC#92491108` | 45772.035181 | 0 | on-chain balance != csv claim by 45772.035181 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `PEOPLE@ETH#25725648` | 27925.52259804 | 0 | on-chain balance != csv claim by 27925.52259804 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x2ef3ddc24e4020173da2c0f3a148eff4a9f6d9dd` | `USDC-AVAXC@AVAXC#92491108` | 44915.110406 | 0 | on-chain balance != csv claim by 44915.110406 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x5d002b87b6e1bdf390027c1210fdc1f6dd2ed85c` | `USDT-OPTIMISM@OPTIMISM#155389411` | 21539.135155 | 0 | on-chain balance != csv claim by 21539.135155 (provider=https://mainnet.optimism.io) |
| `0xcf67640fd0716eb9c0d18c32214a36b8d8a2abaf` | `USDT-ARBITRUM@ARBITRUM#493128998` | 23749.40165 | 0 | on-chain balance != csv claim by 23749.40165 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x1cdf5ba7297f276bf0320c88ca948d4311c913ce` | `USDT-ARBITRUM@ARBITRUM#493128998` | 19295.428303 | 0 | on-chain balance != csv claim by 19295.428303 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x85f56cd748379755494f1365a3efc19569973e99` | `USDC-AVAXC@AVAXC#92491108` | 30345.483935 | 0 | on-chain balance != csv claim by 30345.483935 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xfbabce3f61f043cad4516e8ec526717133c12057` | `USDT-OPTIMISM@OPTIMISM#155389411` | 17186.807277 | 0 | on-chain balance != csv claim by 17186.807277 (provider=https://mainnet.optimism.io) |
| `0x74f42a5bc01eaeaf724da49254da739c37b5a808` | `USDC-AVAXC@AVAXC#92491108` | 27771.381957 | 0 | on-chain balance != csv claim by 27771.381957 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x0ca0e2aa6f2203ef8fa6e1f7f1903fb3a8ec1d35` | `USDT-OPTIMISM@OPTIMISM#155389411` | 16099.499621 | 0 | on-chain balance != csv claim by 16099.499621 (provider=https://mainnet.optimism.io) |
| `0x41ddf01de355547d821873ebd4fb2ffa2dc63d4c` | `POLY-USDC@MATIC#91781706` | 15144.0145 | 0.00009 | on-chain balance != csv claim by 15144.01441 (provider=https://polygon.gateway.tenderly.co) |
| `0xc30082de324fbfb7b293bcda397e2da2942d861d` | `USDT-OPTIMISM@OPTIMISM#155389411` | 12949.280826 | 0 | on-chain balance != csv claim by 12949.280826 (provider=https://mainnet.optimism.io) |
| `0x6cc5f688a315f3dc28a7781717a9a798a59fda7b` | `OKB@ETH#25725648` | 11302.306303002590174878 | 0 | on-chain balance != csv claim by 11302.306303002590174878 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x28734d8ccfa30a3270513b40fd147071b79d7638` | `USDT-OPTIMISM@OPTIMISM#155389411` | 8593.202578 | 0 | on-chain balance != csv claim by 8593.202578 (provider=https://mainnet.optimism.io) |
| `0x00000000004e3d5628234f18b977041e5242651f` | `POLY-USDC@MATIC#91781706` | 8106.253163 | 0 | on-chain balance != csv claim by 8106.253163 (provider=https://polygon.gateway.tenderly.co) |
| `0x26738cb185d888bb8c77ef0821dae410d0a748de` | `USDT-OPTIMISM@OPTIMISM#155389411` | 6450.351135 | 0 | on-chain balance != csv claim by 6450.351135 (provider=https://mainnet.optimism.io) |
| `0x48120be2162bf6d7b1c8ea04a150a530d294dcd0` | `PEOPLE@ETH#25725648` | 6356.7 | 0 | on-chain balance != csv claim by 6356.7 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x787d2b2c401d2635fa4251fe6cf267177f5e45d5` | `PEOPLE@ETH#25725648` | 4912.04882206 | 0 | on-chain balance != csv claim by 4912.04882206 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0d96b0525186c56aab001ff856313c85ced6bbe3` | `PEOPLE@ETH#25725648` | 4822.62 | 0 | on-chain balance != csv claim by 4822.62 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x7cfb1543e6dec0dadd2435f23b97e08d682b6241` | `USDT-ARBITRUM@ARBITRUM#493128998` | 6317.650153 | 0 | on-chain balance != csv claim by 6317.650153 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x0145506e5169a5efcbfc048636553d532108d668` | `PEOPLE@ETH#25725648` | 3857.1065 | 0 | on-chain balance != csv claim by 3857.1065 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xe57245fb0e5fc8d3cb02fb05e1cb0a8939970d6d` | `USDT-OPTIMISM@OPTIMISM#155389411` | 4109.351872 | 0 | on-chain balance != csv claim by 4109.351872 (provider=https://mainnet.optimism.io) |
| `0xdf18e02f023c2f87245e30509ae273854d825573` | `PEOPLE@ETH#25725648` | 4239.74390235 | 0 | on-chain balance != csv claim by 4239.74390235 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd56a81842033ac2ab3e958d357501de052bcf7bc` | `PEOPLE@ETH#25725648` | 3637.76547135 | 0 | on-chain balance != csv claim by 3637.76547135 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x8dbeda7085ddb63fd91bdfb8ba2c69a9a06c7589` | `PEOPLE@ETH#25725648` | 3716.7398 | 0 | on-chain balance != csv claim by 3716.7398 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4e3dd9ec3047246a5fd849334438696df4851b97` | `PEOPLE@ETH#25725648` | 3530.58365806 | 0 | on-chain balance != csv claim by 3530.58365806 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x03c1f15c3f2461d6b53cd7ab799b95b6d2e669d0` | `PEOPLE@ETH#25725648` | 3500 | 0 | on-chain balance != csv claim by 3500 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6671c4d6bbdff00ef9aa4ad9488d400a36272cf6` | `PEOPLE@ETH#25725648` | 3236.30079953 | 0 | on-chain balance != csv claim by 3236.30079953 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x1caf4162b5c0d254b4d07809dba8abf34b895b86` | `PEOPLE@ETH#25725648` | 3199.16248961 | 0 | on-chain balance != csv claim by 3199.16248961 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xcc8e49f620f47dee27dd08caa8dd3554302388ea` | `PEOPLE@ETH#25725648` | 3036.885 | 0 | on-chain balance != csv claim by 3036.885 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xe7f6a87b1950b687a31d068bd01c6be81bc7e1af` | `POLY-USDC@MATIC#91781706` | 3003.505924 | 0 | on-chain balance != csv claim by 3003.505924 (provider=https://polygon.gateway.tenderly.co) |
| `0x6fb31ae81a72af0ce3acd41d5eefe57d2758f72c` | `PEOPLE@ETH#25725648` | 2863.17572044 | 0 | on-chain balance != csv claim by 2863.17572044 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x03426aacbcf9030a380b3ba4a1f16b968649ef3d` | `USDT-ARBITRUM@ARBITRUM#493128998` | 3416.372239 | 0 | on-chain balance != csv claim by 3416.372239 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x85af37f0f445e09302bb1f4d49766555d83b49d3` | `PEOPLE@ETH#25725648` | 2814.71191813 | 0 | on-chain balance != csv claim by 2814.71191813 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x8505bc2f4e44844db42cd9935905fbdd6a96d897` | `PEOPLE@ETH#25725648` | 2690.7088 | 0 | on-chain balance != csv claim by 2690.7088 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x45d6e0a10d9a58bba0df2f8f8b66c7b447451bc7` | `PEOPLE@ETH#25725648` | 2578.16279 | 0 | on-chain balance != csv claim by 2578.16279 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd082b5c3b274d4433b2353b927d91a1082e5d1d1` | `PEOPLE@ETH#25725648` | 2459.19770969 | 0 | on-chain balance != csv claim by 2459.19770969 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xed1808cc19a7c8ea66319e560a652b427038b939` | `PEOPLE@ETH#25725648` | 2206.5973858 | 0 | on-chain balance != csv claim by 2206.5973858 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xc83a70f90970cfed8774e2aac3025cacb872e9a3` | `PEOPLE@ETH#25725648` | 2200 | 0 | on-chain balance != csv claim by 2200 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6290a8e42010d4a68c8428920d5113300bdeaa40` | `PEOPLE@ETH#25725648` | 2162.286628 | 0 | on-chain balance != csv claim by 2162.286628 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xb10cd07b0fbec13a9bc1e89e57e51e00fd26586c` | `PEOPLE@ETH#25725648` | 2069.48525718 | 0 | on-chain balance != csv claim by 2069.48525718 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x23d90deb6252846090f83e4e39555199bcc3b729` | `POLY-USDC@MATIC#91781706` | 1956.222623 | 0 | on-chain balance != csv claim by 1956.222623 (provider=https://polygon.gateway.tenderly.co) |
| `0x6682018cc6f3e6138ea8e63d6b87b2a48f445c1f` | `PEOPLE@ETH#25725648` | 1893.63922039 | 0 | on-chain balance != csv claim by 1893.63922039 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x02a899128336dcc1cd087386292153e01de3848a` | `PEOPLE@ETH#25725648` | 1900.10162555 | 0 | on-chain balance != csv claim by 1900.10162555 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `OKB@ETH#25725648` | 1887.741668 | 0 | on-chain balance != csv claim by 1887.741668 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x942fd3a57d045b83f1e8f1d0381b496950cc16d3` | `USDC-AVAXC@AVAXC#92491108` | 2935.486461 | 0 | on-chain balance != csv claim by 2935.486461 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x012bcc0b88e1ebe69ec29713f0a8f94d012a9de7` | `POLY-USDC@MATIC#91781706` | 1269.682734 | 0 | on-chain balance != csv claim by 1269.682734 (provider=https://polygon.gateway.tenderly.co) |
| `0x5a741eab50d90a955392fcb06d77e81fccefdb3d` | `USDT-OPTIMISM@OPTIMISM#155389411` | 1230.714201 | 0 | on-chain balance != csv claim by 1230.714201 (provider=https://mainnet.optimism.io) |
| `0x34bba0d0a1d0cab16ef6be267465efe5c4df201a` | `USDT-ARBITRUM@ARBITRUM#493128998` | 1439.092055 | 0 | on-chain balance != csv claim by 1439.092055 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x1f1ca79462ec79445c6bd10b9225e1ad2ab9f3d1` | `USDC-AVAXC@AVAXC#92491108` | 1964.34014 | 0 | on-chain balance != csv claim by 1964.34014 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x10e7d149e73dae219bb517de0fcb6a9601ba0f02` | `POLY-USDC@MATIC#91781706` | 949.469514 | 50 | on-chain balance != csv claim by 899.469514 (provider=https://polygon.gateway.tenderly.co) |
| `0x072ecd8fac201bb3951d48e2796c917939c2ecf1` | `USDT-ARBITRUM@ARBITRUM#493128998` | 1000 | 0 | on-chain balance != csv claim by 1000 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xabf8a176e8077047cc009263e89081f1a7b6bff6` | `POLY-USDC@MATIC#91781706` | 813.75688 | 0 | on-chain balance != csv claim by 813.75688 (provider=https://polygon.gateway.tenderly.co) |
| `0xc8802feab2fafb48b7d1ade77e197002c210f391` | `USDC-AVAXC@AVAXC#92491108` | 1683.735722 | 0 | on-chain balance != csv claim by 1683.735722 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf7f767ec932a506efce0b644a2344f628f804bfd` | `USDT-ARBITRUM@ARBITRUM#493128998` | 522.641425 | 0 | on-chain balance != csv claim by 522.641425 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x000000000000d5775ff7721cefb8097af62e52dd` | `POLY-USDC@MATIC#91781706` | 459.213822 | 0 | on-chain balance != csv claim by 459.213822 (provider=https://polygon.gateway.tenderly.co) |
| `0x8c1dd7385efdd54b96dd417ec192f8dc17a13f92` | `USDC-AVAXC@AVAXC#92491108` | 742.005374 | 0 | on-chain balance != csv claim by 742.005374 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `OKB@ETH#25725648` | 396.1273811663 | 0 | on-chain balance != csv claim by 396.1273811663 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x8b98797da896969b57109c84ab5dba69815b97f6` | `USDT-OPTIMISM@OPTIMISM#155389411` | 399.95 | 0 | on-chain balance != csv claim by 399.95 (provider=https://mainnet.optimism.io) |
| `0x7037f56eaa9b91ef5de0fd9e4a2016d5a914e070` | `USDT-ARBITRUM@ARBITRUM#493128998` | 377.425583 | 0 | on-chain balance != csv claim by 377.425583 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xc9b00b361785f3af5fd1027554ff1509923691b6` | `USDT-OPTIMISM@OPTIMISM#155389411` | 326.429142 | 0 | on-chain balance != csv claim by 326.429142 (provider=https://mainnet.optimism.io) |
| `0x769fd5192ffd6bb8645b1f8f7f60cae2854992e7` | `USDT-ARBITRUM@ARBITRUM#493128998` | 322.681304 | 0 | on-chain balance != csv claim by 322.681304 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `OKB@ETH#25725648` | 233.510185 | 0 | on-chain balance != csv claim by 233.510185 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xcfb0295ae71f3f149a1b9249b3955f84f18e6dcd` | `USDT-ARBITRUM@ARBITRUM#493128998` | 263.92 | 0 | on-chain balance != csv claim by 263.92 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xd82b9773eb816ef935c0f820a754900109730f6d` | `USDT-OPTIMISM@OPTIMISM#155389411` | 233.607407 | 0 | on-chain balance != csv claim by 233.607407 (provider=https://mainnet.optimism.io) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `OKB@ETH#25725648` | 185.12488712 | 0 | on-chain balance != csv claim by 185.12488712 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf2364aef78ce6914c37e412bc42d8dea3212538a` | `USDT-ARBITRUM@ARBITRUM#493128998` | 201 | 0 | on-chain balance != csv claim by 201 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x261c9dcc80c35c366e4fa942b7505a02a2f89b32` | `USDT-ARBITRUM@ARBITRUM#493128998` | 160 | 0 | on-chain balance != csv claim by 160 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xe995fe7dc72f7a817ca38c70004fa5388e7f4c2e` | `USDT-OPTIMISM@OPTIMISM#155389411` | 114.71218 | 0 | on-chain balance != csv claim by 114.71218 (provider=https://mainnet.optimism.io) |
| `0x6f7ba0eb23c8fd9c2601ce6c2248a037c732fa40` | `POLY-USDC@MATIC#91781706` | 109.220287 | 0 | on-chain balance != csv claim by 109.220287 (provider=https://polygon.gateway.tenderly.co) |
| `0xaa488c5dd62d17465128c5034179e88a4f2f4170` | `USDC-AVAXC@AVAXC#92491108` | 133.769454 | 0 | on-chain balance != csv claim by 133.769454 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x62383739d68dd0f844103db8dfb05a7eded5bbe6` | `USDC-AVAXC@AVAXC#92491108` | 123.402721 | 0 | on-chain balance != csv claim by 123.402721 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf28143f7b282e59bd5f979012982e7cb9d9b95b0` | `USDT-ARBITRUM@ARBITRUM#493128998` | 101.250523 | 0 | on-chain balance != csv claim by 101.250523 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x6e1702bde7bf3dee0044a06f7eda171023afc171` | `USDT-ARBITRUM@ARBITRUM#493128998` | 99.9 | 0 | on-chain balance != csv claim by 99.9 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xd953e4717c4aa75b08052645379cec64d539ea5a` | `USDT-ARBITRUM@ARBITRUM#493128998` | 96.207252 | 0 | on-chain balance != csv claim by 96.207252 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xc7f758b2126a33239fb7e80f91c722caf7126c15` | `USDT-OPTIMISM@OPTIMISM#155389411` | 89.733851 | 0 | on-chain balance != csv claim by 89.733851 (provider=https://mainnet.optimism.io) |
| `0xdc80e4412b4c6817318b3450b98f87f745206e0b` | `USDC-AVAXC@AVAXC#92491108` | 109.217851 | 0 | on-chain balance != csv claim by 109.217851 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x5aa807ac72fce975be826529d9d72cd70fb8e76b` | `USDT-ARBITRUM@ARBITRUM#493128998` | 89.025 | 0 | on-chain balance != csv claim by 89.025 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xc496e5e40f18035d6f8e805f1e5595657d7dca9f` | `USDT-ARBITRUM@ARBITRUM#493128998` | 88.896744 | 0 | on-chain balance != csv claim by 88.896744 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x0938c63109801ee4243a487ab84dffa2bba4589e` | `USDT-ARBITRUM@ARBITRUM#493128998` | 78.397627 | 0 | on-chain balance != csv claim by 78.397627 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x40c436c9b5890e3dce36e524b835e95a504bf305` | `USDT-ARBITRUM@ARBITRUM#493128998` | 75.057882 | 0 | on-chain balance != csv claim by 75.057882 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xa9ac43f5b5e38155a288d1a01d2cbc4478e14573` | `OKB@ETH#25725648` | 35.38 | 0 | on-chain balance != csv claim by 35.38 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `OKB@ETH#25725648` | 34.5644520006 | 0 | on-chain balance != csv claim by 34.5644520006 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `OKB@ETH#25725648` | 18.95 | 0 | on-chain balance != csv claim by 18.95 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `OKB@ETH#25725648` | 14.607403 | 0 | on-chain balance != csv claim by 14.607403 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `OKB@ETH#25725648` | 8.2962431535 | 0 | on-chain balance != csv claim by 8.2962431535 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xb0473ca01e91981623b65daf84ab7bc94a7e4000` | `USDC-AVAXC@AVAXC#92491108` | 24.9968 | 0 | on-chain balance != csv claim by 24.9968 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x2ad0f758ce754653c79b681750bff16fe005f230` | `USDC-AVAXC@AVAXC#92491108` | 4.999999 | 0 | on-chain balance != csv claim by 4.999999 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x136bc54e0950eae85591954ae2bdea2199681234` | `USDC-AVAXC@AVAXC#92491108` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x7970901be38fc9235a5f650aa3708441b378f466` | `USDC-AVAXC@AVAXC#92491108` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe035cb2a4719c0439c81413543d4eee39c0aaeab` | `USDC-AVAXC@AVAXC#92491108` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9157b6c4f48112490962175b4058d485939653ca` | `USDC-AVAXC@AVAXC#92491108` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x1810a491db6f47ea99e31903db4b65f74eaf0b36` | `USDC-AVAXC@AVAXC#92491108` | 4.9996 | 0 | on-chain balance != csv claim by 4.9996 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xceca016c59a63cccdf7a703f52b9dee581e1fa40` | `USDC-AVAXC@AVAXC#92491108` | 4.99968 | 0 | on-chain balance != csv claim by 4.99968 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x920c6944f9e1c5bb7a80d180a8df75f65b0b94e5` | `USDC-AVAXC@AVAXC#92491108` | 4.9997 | 0 | on-chain balance != csv claim by 4.9997 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xc6eca660c5da896784071ebc252fc27184da54cf` | `USDC-AVAXC@AVAXC#92491108` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x381686ce20da524ad9520a157cdf7c871827ce0a` | `USDC-AVAXC@AVAXC#92491108` | 4.999362 | 0 | on-chain balance != csv claim by 4.999362 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe270b325c128ac89107ef4a078b303bf81dc7538` | `USDC-AVAXC@AVAXC#92491108` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd3c8185a0aaf673a5a94da03a2e48311085e6058` | `USDC-AVAXC@AVAXC#92491108` | 4.99936 | 0 | on-chain balance != csv claim by 4.99936 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd26a7101d2212f4f23c9b8a6af35f7ec432e5c4d` | `USDC-AVAXC@AVAXC#92491108` | 4.999199 | 0 | on-chain balance != csv claim by 4.999199 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9bed615d3e388d91b7302d0710949b1cf5fb6544` | `USDC-AVAXC@AVAXC#92491108` | 4.999081 | 0 | on-chain balance != csv claim by 4.999081 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xfee1e3ec3577bc256b3c3531de48ea7a8eeff4db` | `USDC-AVAXC@AVAXC#92491108` | 4.99919 | 0 | on-chain balance != csv claim by 4.99919 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x62a5b0e0b15f850016783aa8bfc9fa19cdb69208` | `USDC-AVAXC@AVAXC#92491108` | 4.99906 | 0 | on-chain balance != csv claim by 4.99906 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf52ac68c0b89cc62846edba7733f645e17f8f98e` | `USDC-AVAXC@AVAXC#92491108` | 4.998932 | 0 | on-chain balance != csv claim by 4.998932 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x52aeace7f732f4f93ec47caac1bf0752fa2b09dd` | `USDC-AVAXC@AVAXC#92491108` | 4.998905 | 0 | on-chain balance != csv claim by 4.998905 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x454bfa69931e6f7e27ad630dd37531241cb73394` | `USDC-AVAXC@AVAXC#92491108` | 4.998883 | 0 | on-chain balance != csv claim by 4.998883 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x7d2968ddb30db723ff62896821d506234f0a18a4` | `USDC-AVAXC@AVAXC#92491108` | 4.998758 | 0 | on-chain balance != csv claim by 4.998758 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x352644b25656b1f2503e448f472e77c078e2e354` | `OKB@ETH#25725648` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x5bfed738dd750fee869b2ae816827af9df123055` | `OKB@ETH#25725648` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xa765013bb5fbb4f16e12ede43288f915021fc5d7` | `OKB@ETH#25725648` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf7bb83b23869fde84375aaffa036af33e9ba8f96` | `OKB@ETH#25725648` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x030ae17036335c1c8c1b166f419b3f0721a5f184` | `OKB@ETH#25725648` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x08085c84b8e5fa4fe61f568b9990ae100935e079` | `USDC-AVAXC@AVAXC#92491108` | 4.998588 | 0 | on-chain balance != csv claim by 4.998588 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x988033c6762141b3f8e20bedfb4bc74b508a9ae0` | `USDC-AVAXC@AVAXC#92491108` | 4.998522 | 0 | on-chain balance != csv claim by 4.998522 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9f25df81d05476f5e56309dd95d7c90f47f55452` | `OKB@ETH#25725648` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x2801845ac56549789ec24377981333b73b651b2d` | `OKB@ETH#25725648` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x36e0c1f02da4aa7c640fd0f3a56e3ab08262b47d` | `OKB@ETH#25725648` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `TCw6YaWm3y6DvxY7M8hrCDnrJGeGMumzGJ` | `TRX@TRX#85235532` | 448873177.240876 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 448873177.240876 (provider=https://api.trongrid.io) |
| `TLaGjwhvA8XQYSxFAcAXy7Dvuue9eGYitv` | `USDT-TRC20@TRX#85235532` | 172761635.851229 | 94371239.829759 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 78390396.02147 (provider=https://tron-rpc.publicnode.com) |
| `0xa09fe38187abecc2768b98b80ea30afbc2fbd774` | `USDC@ETH#25725648` | 92957342.3756 | 14.232575 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=92957328.143025 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `TYE5mNXSQHuPAdCa1etk76BKybD4FKJto6` | `USDT-TRC20@TRX#85235532` | 55965608.089424 | 31592459.6782 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 24373148.411224 (provider=https://api.trongrid.io) |
| `TWAY47W8ci9XtfMgvs4r2A7Z2NFGUDJqpp` | `TRX@TRX#85235532` | 55745387 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 55745387 (provider=https://api.trongrid.io) |
| `0x7fea5b4568751533039179116e372e26b6b41b13` | `USDC@ETH#25725648` | 34643768.911411 | 4.641389 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=34643764.270022 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `THdXEcBCZ7nagnmfRrfSHoiLRTbSRueoVW` | `TRX@TRX#85235532` | 36520573.613888 | 9560158.059565 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 26960415.554323 (provider=https://tron-rpc.publicnode.com) |
| `TNYowvaHncRBWAKDzAySA3uVHLkB9qs62p` | `USDT-TRC20@TRX#85235532` | 32475866.202703 | 0.000037 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 32475866.202666 (provider=https://tron-rpc.publicnode.com) |
| `TLaGjwhvA8XQYSxFAcAXy7Dvuue9eGYitv` | `TRX@TRX#85235532` | 31258245.302261 | 1970319.151802 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 29287926.150459 (provider=https://tron-rpc.publicnode.com) |
| `TJVeQfDKsV2ZqrDzZKabgoE6De1vJhnP7h` | `TRX@TRX#85235532` | 18166524.40581 | 174790.404107 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 17991734.001703 (provider=https://tron-rpc.publicnode.com) |
| `TBNHZP2iAVyDhpTSZ8ee2Xjci1JExZTPE5` | `TRX@TRX#85235532` | 18140451.003736 | 40451.336404 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 18099999.667332 (provider=https://tron-rpc.publicnode.com) |
| `TD9GNrgmrg46yistfTMJpRQZcDmDRx1Bmq` | `TRX@TRX#85235532` | 14151058.435524 | 1058.435543 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 14149999.999981 (provider=https://tron-rpc.publicnode.com) |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `OKB-X1@XLAYER` |  |  | rpc error (provider=): network XLAYER not configured |
| `TFAXd82DvBpykbt9JyTpevDZWfkmqy2CnZ` | `TRX@TRX#85235532` | 8106729.088911 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 8106729.088911 (provider=https://api.trongrid.io) |
| `TZ7wp7mcM2Qq4kmDNW3pS7ZZ8Y8ztXUVji` | `TRX@TRX#85235532` | 7219069.00376 | 14198.99889 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 7204870.00487 (provider=https://tron-rpc.publicnode.com) |


_158 additional `WARN` rows omitted; see verification bundle._

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 13, WARN 267, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `C68a6RCGLiPskbPYtAcsCjhG8tfTWYcoB4JjCrXFdqyo` | `USDG\|SOL#438431613` | 96736951.014949 | 96736831.43769 | on-chain balance != csv claim by 119.577259 (provider=https://api.helius.xyz/v1/wallet) |
| `3WjxkXbu9JBosUWUnJpjzec7Wa4rSkWo9khUbbLDCssk` | `SOL\|SOL#438431613` | 3310582.822857096 | 3.971948082 | on-chain balance != csv claim by 3310578.850909014 (provider=https://api.helius.xyz/v1/wallet) |
| `CKPbSRoW17z8fkFBy5NPcitiLZX4uqmfSv5MpFewkNoQ` | `USDC-SPL\|SOL#438431613` | 866263.490112 | 866202.624302 | on-chain balance != csv claim by 60.86581 (provider=https://api.helius.xyz/v1/wallet) |
| `Bou5nPZCcqQCNhSaKyJMLhqgoZVKBx6yoDkGEW5WY6W5` | `USDC-SPL\|SOL#438431613` | 322988.606546 | 322988.218746 | on-chain balance != csv claim by 0.3878 (provider=https://api.helius.xyz/v1/wallet) |
| `CyKe1nTzB56CLtwgLiJ4WWTDvAjndwTfg4qs8j89k1NU` | `USDC-SPL\|SOL#438431613` | 241998.153924 | 0.4307 | on-chain balance != csv claim by 241997.723224 (provider=https://api.helius.xyz/v1/wallet) |
| `GD3tr9DR8ZnSy2pQ22RwzSR4vKY48SxHLdgCymXQ6W5i` | `SOL\|SOL#438431613` | 210602.398535009 | 0.165141123 | on-chain balance != csv claim by 210602.233393886 (provider=https://api.helius.xyz/v1/wallet) |
| `GHWmH2ZjB6jhv68dzRDvQeC3v23oyxh1W2F4Av6Gj2fy` | `USDC-SPL\|SOL#438431613` | 208608.159359 | 208608.042059 | on-chain balance != csv claim by 0.1173 (provider=https://api.helius.xyz/v1/wallet) |
| `8wn2fZPKmUSMqNKLNsBfb3KkrVhChYYrL5Kzyfc2Dqox` | `USDC-SPL\|SOL#438431613` | 158261.865644 | 0.3959 | on-chain balance != csv claim by 158261.469744 (provider=https://api.helius.xyz/v1/wallet) |
| `92QMPBK852QdkGxyyuQZ8BmKczeQihEiGLHbFbx4tVdC` | `USDC-SPL\|SOL#438431613` | 136977.500328 | 0.4325 | on-chain balance != csv claim by 136977.067828 (provider=https://api.helius.xyz/v1/wallet) |
| `B6gPEoYfgKDeMsT2Mxkie4SBsW5GMLQyCDQk42WTFDR` | `USDC-SPL\|SOL#438431613` | 86551.727026 | 86551.381926 | on-chain balance != csv claim by 0.3451 (provider=https://api.helius.xyz/v1/wallet) |
| `AisgV457n3GHkXhmqFh91uz3a9LcQi96if2r934WVXQE` | `USDC-SPL\|SOL#438431613` | 74829.816836 | 0.0782 | on-chain balance != csv claim by 74829.738636 (provider=https://api.helius.xyz/v1/wallet) |
| `4YWWA6aFe3E9vh7Bs3XqgiHqvD5pBZiuAZqk3K1X5drt` | `USDC-SPL\|SOL#438431613` | 21977.557929 | 21977.517929 | on-chain balance != csv claim by 0.04 (provider=https://api.helius.xyz/v1/wallet) |
| `3n8t1n1NF6rhDeAPqQJF7RDD6JLGJPQftF1gMEKd4Sfs` | `SOL\|SOL#438431613` | 12529.097548696 | 0.165247137 | on-chain balance != csv claim by 12528.932301559 (provider=https://api.helius.xyz/v1/wallet) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `AC8azEyeWHk1L7XKvyeBjJtVwRwXgmP5hQ` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/AC8azEyeWHk1L7XKvyeBjJtVwRwXgmP5hQ?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."}
 |
| `ADNbM5fBujCRBW1vqezNeAWmnsLp19ki3n` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/ADNbM5fBujCRBW1vqezNeAWmnsLp19ki3n?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."}
 |
| `f1mfxobliwdamau2x25t4mcegzhfln5rwng4aemqi` | `FIL\|FIL#6269036` | 74868592.010029404977985748 | 71938987.044683732279470548 | live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes; on-chain balance != csv claim by 2929604.9653456726985152 (provider=https://api.node.glif.io/rpc/v0) |
| `9vPesMs2Md65mu2JBEccp8RKBqccqaLMRL` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/9vPesMs2Md65mu2JBEccp8RKBqccqaLMRL?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."}
 |
| `ABKETAr5U7DP4k1yE8x9z2eFAmM1LxbP62` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/ABKETAr5U7DP4k1yE8x9z2eFAmM1LxbP62?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."} |
| `0xd0c23b2613e9c14f4888ff6c4359f197424ca88cf2970bac9e862e5ee65d5fa8` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd0c23b2613e9c14f4888ff6c4359f197424ca88cf2970bac9e862e5ee65d5fa8/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=962785370): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(962785370) has been pruned","hint":"The requested data has been |
| `A7zXZBP6SzR1yNaFxD3jKAT8YPkY8Uz9Ch` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/A7zXZBP6SzR1yNaFxD3jKAT8YPkY8Uz9Ch?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."} |
| `AA431oUHgBPpuEbGt6G1Z96bgkZhyE4KTx` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/AA431oUHgBPpuEbGt6G1Z96bgkZhyE4KTx?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."} |
| `0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572` | `APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=962785370): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(962785370) has been pruned","hint":"The requested data has been |
| `9za4ts8qmKDqNqPZeMj9ceB7nH1eMo6k6P` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/9za4ts8qmKDqNqPZeMj9ceB7nH1eMo6k6P?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."} |
| `A4DCEMmVLD39r4iUzS6tZTbD8VbsXpfCVH` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/A4DCEMmVLD39r4iUzS6tZTbD8VbsXpfCVH?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."} |
| `0xffd7f0d7b24ba690657dd2aa7ccb754c4be938e4fb970c72eef966d664bd1a45` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xffd7f0d7b24ba690657dd2aa7ccb754c4be938e4fb970c72eef966d664bd1a45/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=962785370): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(962785370) has been pruned","hint":"The requested data has been |
| `0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=962785370): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(962785370) has been pruned","hint":"The requested data has been |
| `EQDbf_PU1DLAcneMzWjjUbK_Vc7oahP0iOpuHbxIwDfaI5Rm` | `TONCOIN-NEW\|TON#85240351` | 6643521.230872047 | 12.191071225 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 6643509.039800822 (provider=https://tonapi.io/v2/accounts/EQDbf_PU1DLAcneMzWjjUbK_Vc7oahP0iOpuHbxIwDfaI5Rm) |
| `A2vrh5EQPv6vcTedvMZK9AYwLSsx1KkEUV` | `DOGE\|DOGE` |  |  | ledger rpc error (provider=https://api.blockcypher.com/v1/doge/main/addrs/A2vrh5EQPv6vcTedvMZK9AYwLSsx1KkEUV?before=6326404): blockcypher HTTP 429: {"error": "Limits reached."} |


_252 additional `WARN` rows omitted; see verification bundle._

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 500375535 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-08-10T16:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary and wallet_address_list present. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `pass` | Wallet list present; hot-wallet on-chain verifier ran. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `fail` | address-ownership verifier FAIL. |
| Is global_proof public and independently reproducible? (S1-3) | `pass` | global-zk-proof verifier PASS. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `pass` | Plonky2 / zk-STARK uses transparent setup; no trusted setup ceremony required. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P1` | Anchor root/proof/commitment on-chain or publish immutable DA/archive records. | `NO_CANONICAL_ANCHOR` |
| `P1` | Add daily anchoring, weekly full PoR, or event-triggered updates to reduce timing risk. | `HIGH_FREQUENCY_GAP` |
| `P2` | Review onchain-balance verifier FAIL/WARN outliers in the verification bundle and reconcile CSV vs chain deltas. | `IMPLEMENTATION_SEMANTICS_RISK` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.
