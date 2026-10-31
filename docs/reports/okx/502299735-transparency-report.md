# OKX 502299735 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [502299735-assessment.json](./502299735-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `502299735` |
| Snapshot time | `2026-09-07T16:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

OKX audit 502299735 does not meet all Stage 1 public artifact thresholds in this pipeline.

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
| Previous snapshot | `2026-08-10T16:00:00Z` |
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
| `summarySnapshot` | `1e10e866ace60e02a4d551516034b3eb6542e51009f0414c4c7cb17285c4f052` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `fa84690d06e7b5c883f3e6f24c635781a74321a67a76f1d0a7aab45528364641` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026090800_V1.zip |
| `globalProofBundle` | `9b0508cab98f21bdbf391bc8e5b4aae8fd9acff6ef09eb182bacaa171c8f42d1` | https://static.okx.com/cdn/okx/por/merkel/por_502299735_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x808ffacc9f2e261ce92e8cd824d274ebafce53028d31ed357ef85fc9a6ae45c3` |
| Verification bundle root | `0xa8eb5f38c988bd7fefb7deebebfaa8979a6b0b1d002d81795ae8e584a38a7c6c` |
| Artifact bundle SHA-256 | `d329da9901b57664bb2e78bceb5721d3e37f0c0d01e154be138b40ff1844f5d9` |
| Verification bundle SHA-256 | `f2356ada437d81505521dc2589f5b38f69e4cc2cbbac4b7958447e49f0161b47` |

Local bundle paths: [502299735.artifact-bundle.json](../../../artifacts/okx/502299735/bundles/502299735.artifact-bundle.json), [502299735.verification-bundle.v2.json](../../../artifacts/okx/502299735/bundles/502299735.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | okx summary has no reserve currencies |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `FAIL` | `0.9993` | 364/496948 address signatures failed |
| `onchain-balance-hot` | `2.1` | `PARTIAL` | `0.0001` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.0013` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.0006` |  |
| `global-zk-proof` | `okx-1` | `PASS` | `1.0000` | zkSTARKValidator verify-global succeeded; summary merkle root bound to zk proof |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `502299735.verification-bundle.v2.json` in the local artifact bundle.

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

Summary: 364/496948 address signatures failed

Finding counts: FAIL 5

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x80299387783b6b8c0645a17ea72fdc3cd32e0c9367445fdbe7344742af521fc39afc1fbc3917efc3a4bd72c48814e9b0` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x8063539f96661a8b8b86222b7170ee312c0ecc523b89e3f938698f69d0f82908d3db19259e047afc5c780b12eda785bc` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x806e4c086bc3aff5b9940868fac93e11b139aa04cf749a1adf13e8c6058e861df0426488db4c04d249d26d7d55e9a650` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x8082858d8486e36a4102dffb38cd65ad34aa93b7e81ecca681004508872f4df261ecf7b9ac3e741288f9a0ccbc3a9d0b` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x80a935b3b8a82f6a032002dd04b9cac3f0844c9b73b3513067f97b1141e6c0b5b89ed36d0a4c74c59d98498220805c27` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |

#### `onchain-balance-hot` (`PARTIAL`)

Finding counts: UNVERIFIABLE 22, PASS 2

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `ASSET-HUB\|DOT` | 59 | no native verifier for this (coin,network) pair yet |
| `A\|VAULTA` | 2 | no native verifier for this (coin,network) pair yet |
| `DOTK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ELF\|AELF` | 273 | no native verifier for this (coin,network) pair yet |
| `ETCK-KIP20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ETC\|ETC` | 40 | no native verifier for this (coin,network) pair yet |
| `ETHK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `FILK-OKC20\|OKC` | 8 | no native verifier for this (coin,network) pair yet |
| `LINKK-OKC20\|OKC` | 3 | no native verifier for this (coin,network) pair yet |
| `OKB-OKC20\|OKC` | 4 | no native verifier for this (coin,network) pair yet |
| `TRXK-KIP20\|OKC` | 14 | no native verifier for this (coin,network) pair yet |
| `USDC-OKC20\|OKC` | 248 | no native verifier for this (coin,network) pair yet |
| `USDT-OKC20\|OKC` | 105 | no native verifier for this (coin,network) pair yet |
| `XRPK-KIP20\|OKC` | 5 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x85e306287d845752a1716bfd5e90f266d08377dd1d5e4189d6f0bc05a9860dac6a0f73d2f6ec7b1458681ba0570afac3` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xb134297b6c0f624afab8d357c9c2f8324dfa844466c35a12d33ebe4959dfe06d88f9544593cf17996bbc0cca4d0a43d9` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x80670ceb430da829ff33566afdecede98d4bcdb1bdbc33acb16425e1cdd29a8088f29bcd7c5ba104171eaf3b1b369291` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x821b54086bedce6be4542660391d66f72779165fd854c59ecb904f3ae4171b17b52e3a8784c326db3de888e188cffb49` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x8a7fa64cd9d0fcbc52f417d7b05271b151219659da4ecb306385dfb06d049af43783c4662e22e48175cedfc8e8438bf8` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xad3359fe5b728c2fc57647d04ab0eacddbdb4da1f4820c72b5a11c203e26fec333e0ea771e245262ba06d2d3d75a0e31` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xb2d016a313305400b7125cc0c5bd9a6acb001cd640648084d63d49154451a700feca775c17c09468638e517a4ceb42f8` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xa2b4a69dc29719fcf4c400c5ba4c60966f216fafea399cfc2024cf07ad60d9258291d5f83b8094225bbe86f646a67f1b` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |

#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 186, WARN 157, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `PEOPLE@ETH#25926528` | 1129696969.05960099 | 0 | on-chain balance != csv claim by 1129696969.05960099 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `PEOPLE@ETH#25926528` | 135187613.9389226963 | 0 | on-chain balance != csv claim by 135187613.9389226963 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `PEOPLE@ETH#25926528` | 75572205.5509888722 | 0 | on-chain balance != csv claim by 75572205.5509888722 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x85dcd76d4fbd3aa0c85c27b9441222c19a14134b` | `PEOPLE@ETH#25926528` | 6861400.987827 | 0 | on-chain balance != csv claim by 6861400.987827 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `PEOPLE@ETH#25926528` | 6262293.9383380403 | 0 | on-chain balance != csv claim by 6262293.9383380403 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xafee421482faea92292ed3ffe29371742542ad72` | `USDT-ARBITRUM@ARBITRUM#502728533` | 7871433.934394 | 0 | on-chain balance != csv claim by 7871433.934394 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xb5216cb558cb018583bed009ee25ca73eb27bb1d` | `USDT-OPTIMISM@OPTIMISM#156599011` | 4835707.283356 | 0 | on-chain balance != csv claim by 4835707.283356 (provider=https://mainnet.optimism.io) |
| `0xd8a3ed83a59d3e10d82a4b845d64b24ba0572baf` | `USDT-ARBITRUM@ARBITRUM#502728533` | 4043997.159049 | 0 | on-chain balance != csv claim by 4043997.159049 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x92ea7496eba5f001d620005f88f3e8e686e3d4ea` | `PEOPLE@ETH#25926528` | 2859913.01 | 0 | on-chain balance != csv claim by 2859913.01 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x343d752bb710c5575e417edb3f9fa06241a4749a` | `POLY-USDC@MATIC#93394510` | 2619922.738323 | 46.514712 | on-chain balance != csv claim by 2619876.223611 (provider=https://polygon.gateway.tenderly.co) |
| `0xe7c5e276b2434743f6c79365d967cbe4291477bc` | `USDT-ARBITRUM@ARBITRUM#502728533` | 1867984.662931 | 0 | on-chain balance != csv claim by 1867984.662931 (provider=https://arb1.arbitrum.io/rpc) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `PEOPLE@ETH#25926528` | 1226569.4852159399 | 0 | on-chain balance != csv claim by 1226569.4852159399 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a1ca3e90e9e3a00eb27e9f409b4db310aa7405f` | `USDT-ARBITRUM@ARBITRUM#502728533` | 1012122.215906 | 0 | on-chain balance != csv claim by 1012122.215906 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xa0c9b4ce1785953abdceae569789aa158a01da4c` | `USDT-ARBITRUM@ARBITRUM#502728533` | 756108.097645 | 0 | on-chain balance != csv claim by 756108.097645 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xde3c8d726e906b379a020cdb04bbe4b8f39a726f` | `USDC-AVAXC@AVAXC#94703902` | 1099692.422633 | 0 | on-chain balance != csv claim by 1099692.422633 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x3f03440b19d15ff2068daf68e57eb656bfd0a67e` | `USDC-AVAXC@AVAXC#94703902` | 1044342.751921 | 0 | on-chain balance != csv claim by 1044342.751921 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9098575016440086c694c85c5010a0ed85249675` | `USDT-OPTIMISM@OPTIMISM#156599011` | 551968.8371 | 0 | on-chain balance != csv claim by 551968.8371 (provider=https://mainnet.optimism.io) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `PEOPLE@ETH#25926528` | 447707.19893191 | 0 | on-chain balance != csv claim by 447707.19893191 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9d1964437bb6d849c66ab2a2a033340dc2d999dc` | `USDC-AVAXC@AVAXC#94703902` | 889242.491949 | 0 | on-chain balance != csv claim by 889242.491949 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x3b50287cf8f073ff443b84de57e2375d446b962b` | `POLY-USDC@MATIC#93394510` | 339577.571006 | 0 | on-chain balance != csv claim by 339577.571006 (provider=https://polygon.gateway.tenderly.co) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `PEOPLE@ETH#25926528` | 318247.6953893999 | 0 | on-chain balance != csv claim by 318247.6953893999 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf418263a6468c80d3bcb71033c56aacf81b241ba` | `PEOPLE@ETH#25926528` | 256304.43 | 0 | on-chain balance != csv claim by 256304.43 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `PEOPLE@ETH#25926528` | 215204.98 | 0 | on-chain balance != csv claim by 215204.98 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xc74e4c556a16390165c99b40aabc39a87c358305` | `POLY-USDC@MATIC#93394510` | 203736.733187 | 0 | on-chain balance != csv claim by 203736.733187 (provider=https://polygon.gateway.tenderly.co) |
| `0x76705493d53beefadffc253acb8aef72b7e7f6c6` | `USDT-OPTIMISM@OPTIMISM#156599011` | 202319.05514 | 0 | on-chain balance != csv claim by 202319.05514 (provider=https://mainnet.optimism.io) |
| `0xb1da6e546bd592218d6e849a8177d9d4fbf12986` | `POLY-USDC@MATIC#93394510` | 187855.42838 | 0 | on-chain balance != csv claim by 187855.42838 (provider=https://polygon.gateway.tenderly.co) |
| `0xb57198a6e5af11544936a520a54268bc6ababb0c` | `USDT-OPTIMISM@OPTIMISM#156599011` | 152998.2332 | 0 | on-chain balance != csv claim by 152998.2332 (provider=https://mainnet.optimism.io) |
| `0x9c308b1665097af7ab2e38a41577bfecdc2d1c5c` | `PEOPLE@ETH#25926528` | 126466.45532091 | 0 | on-chain balance != csv claim by 126466.45532091 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a8c00e80e3cfae4c903f9462b4ee2826404fe05` | `USDT-OPTIMISM@OPTIMISM#156599011` | 117794.00938 | 0 | on-chain balance != csv claim by 117794.00938 (provider=https://mainnet.optimism.io) |
| `0x4cafbe2b5bd89f0c0b7d6fef75b47e7c2ae82f22` | `USDC-AVAXC@AVAXC#94703902` | 285612.69009 | 0 | on-chain balance != csv claim by 285612.69009 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf1e7b58ccd200fcb9400ddbc30867c43b7dc9ee1` | `USDT-ARBITRUM@ARBITRUM#502728533` | 117793.809316 | 0 | on-chain balance != csv claim by 117793.809316 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x5ca1063894e5bbd688830d83fefa5ce3d339744d` | `USDT-ARBITRUM@ARBITRUM#502728533` | 107051.201179 | 0 | on-chain balance != csv claim by 107051.201179 (provider=https://arb1.arbitrum.io/rpc) |
| `0x229b1c70124141ec96945aec2323a22b68bb093c` | `USDT-ARBITRUM@ARBITRUM#502728533` | 103703.860873 | 0 | on-chain balance != csv claim by 103703.860873 (provider=https://arb1.arbitrum.io/rpc) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `OKB@ETH#25926528` | 82886.7363764665 | 0 | on-chain balance != csv claim by 82886.7363764665 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xfbabce3f61f043cad4516e8ec526717133c12057` | `USDT-OPTIMISM@OPTIMISM#156599011` | 73105.112454 | 0 | on-chain balance != csv claim by 73105.112454 (provider=https://mainnet.optimism.io) |
| `0xbe9ac2d59743d69c4054bab09fa9821929d0f31b` | `USDT-ARBITRUM@ARBITRUM#502728533` | 73864.370378 | 0 | on-chain balance != csv claim by 73864.370378 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `PEOPLE@ETH#25926528` | 65097.57 | 0 | on-chain balance != csv claim by 65097.57 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `PEOPLE@ETH#25926528` | 57350.76259804 | 0 | on-chain balance != csv claim by 57350.76259804 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x582988097ef9f2cf50cfac3f5f3b1ed54bc52a88` | `USDT-ARBITRUM@ARBITRUM#502728533` | 49669.675025 | 0 | on-chain balance != csv claim by 49669.675025 (provider=https://arb1.arbitrum.io/rpc) |
| `0x7d328f5d5f4640f81b2f93d2903ec0a6c030eed6` | `PEOPLE@ETH#25926528` | 43996.91003613 | 0 | on-chain balance != csv claim by 43996.91003613 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x2fb2841c747621a7187b6f2a5390ae77321e5257` | `USDC-AVAXC@AVAXC#94703902` | 58878.778092 | 0 | on-chain balance != csv claim by 58878.778092 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x04ba1b380b1418007c1407bbe12ac21a596ae1e0` | `USDC-AVAXC@AVAXC#94703902` | 53271.738263 | 0 | on-chain balance != csv claim by 53271.738263 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x4a152242605ab21df5ac6e819fdb128a4efc33fc` | `POLY-USDC@MATIC#93394510` | 42144.864932 | 0 | on-chain balance != csv claim by 42144.864932 (provider=https://polygon.gateway.tenderly.co) |
| `0x99860ec99aae9c0aa5ab763ee73da1217a9da665` | `POLY-USDC@MATIC#93394510` | 37730.849791 | 0 | on-chain balance != csv claim by 37730.849791 (provider=https://polygon.gateway.tenderly.co) |
| `0x2ef3ddc24e4020173da2c0f3a148eff4a9f6d9dd` | `USDC-AVAXC@AVAXC#94703902` | 47908.361516 | 0 | on-chain balance != csv claim by 47908.361516 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x6bea9e00a6bb63d22904d975bdf8a98f26f9f49d` | `USDC-AVAXC@AVAXC#94703902` | 42144.864981 | 0 | on-chain balance != csv claim by 42144.864981 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xc30082de324fbfb7b293bcda397e2da2942d861d` | `USDT-OPTIMISM@OPTIMISM#156599011` | 35686.46323 | 0 | on-chain balance != csv claim by 35686.46323 (provider=https://mainnet.optimism.io) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `OKB@ETH#25926528` | 30822.1988656955 | 0 | on-chain balance != csv claim by 30822.1988656955 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `PEOPLE@ETH#25926528` | 29298.35 | 0 | on-chain balance != csv claim by 29298.35 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xe7f6a87b1950b687a31d068bd01c6be81bc7e1af` | `POLY-USDC@MATIC#93394510` | 30327.669655 | 0 | on-chain balance != csv claim by 30327.669655 (provider=https://polygon.gateway.tenderly.co) |
| `0x5d002b87b6e1bdf390027c1210fdc1f6dd2ed85c` | `USDT-OPTIMISM@OPTIMISM#156599011` | 24738.391455 | 0 | on-chain balance != csv claim by 24738.391455 (provider=https://mainnet.optimism.io) |
| `0x28734d8ccfa30a3270513b40fd147071b79d7638` | `USDT-OPTIMISM@OPTIMISM#156599011` | 20747.108888 | 0 | on-chain balance != csv claim by 20747.108888 (provider=https://mainnet.optimism.io) |
| `0x3f41ea66a734db98720512df6bd12dbc87c81d30` | `PEOPLE@ETH#25926528` | 15483.1 | 0 | on-chain balance != csv claim by 15483.1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x85f56cd748379755494f1365a3efc19569973e99` | `USDC-AVAXC@AVAXC#94703902` | 36983.700033 | 0 | on-chain balance != csv claim by 36983.700033 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd82b9773eb816ef935c0f820a754900109730f6d` | `USDT-OPTIMISM@OPTIMISM#156599011` | 15795.609613 | 0 | on-chain balance != csv claim by 15795.609613 (provider=https://mainnet.optimism.io) |
| `0x04e9d79a17f4894c76c65c52516f790531ef6bee` | `USDT-ARBITRUM@ARBITRUM#502728533` | 15179.150407 | 0 | on-chain balance != csv claim by 15179.150407 (provider=https://arb1.arbitrum.io/rpc) |
| `0x1cdf5ba7297f276bf0320c88ca948d4311c913ce` | `USDT-ARBITRUM@ARBITRUM#502728533` | 15473.610922 | 0 | on-chain balance != csv claim by 15473.610922 (provider=https://arb1.arbitrum.io/rpc) |
| `0x3c8f7f247c5a494aaf89d1296ef62497e73d15d8` | `USDT-ARBITRUM@ARBITRUM#502728533` | 13241.66 | 0 | on-chain balance != csv claim by 13241.66 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x0ca0e2aa6f2203ef8fa6e1f7f1903fb3a8ec1d35` | `USDT-OPTIMISM@OPTIMISM#156599011` | 12050.763738 | 0 | on-chain balance != csv claim by 12050.763738 (provider=https://mainnet.optimism.io) |
| `0x41ddf01de355547d821873ebd4fb2ffa2dc63d4c` | `POLY-USDC@MATIC#93394510` | 12008.70316 | 0.00009 | on-chain balance != csv claim by 12008.70307 (provider=https://polygon.gateway.tenderly.co) |
| `0x6cc5f688a315f3dc28a7781717a9a798a59fda7b` | `OKB@ETH#25926528` | 11302.306303002590174878 | 0 | on-chain balance != csv claim by 11302.306303002590174878 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x74f42a5bc01eaeaf724da49254da739c37b5a808` | `USDC-AVAXC@AVAXC#94703902` | 22319.051483 | 0 | on-chain balance != csv claim by 22319.051483 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x012bcc0b88e1ebe69ec29713f0a8f94d012a9de7` | `POLY-USDC@MATIC#93394510` | 8607.699744 | 0 | on-chain balance != csv claim by 8607.699744 (provider=https://polygon.gateway.tenderly.co) |
| `0x00000000004e3d5628234f18b977041e5242651f` | `POLY-USDC@MATIC#93394510` | 8106.253163 | 0 | on-chain balance != csv claim by 8106.253163 (provider=https://polygon.gateway.tenderly.co) |
| `0xcf67640fd0716eb9c0d18c32214a36b8d8a2abaf` | `USDT-ARBITRUM@ARBITRUM#502728533` | 6501.233553 | 0 | on-chain balance != csv claim by 6501.233553 (provider=https://arb1.arbitrum.io/rpc) |
| `0xc4eec7444e3b8564b89c6879423e8d66badcebfd` | `USDT-ARBITRUM@ARBITRUM#502728533` | 5899.49 | 0 | on-chain balance != csv claim by 5899.49 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x90cad9621cb2717208fca2f78d97fa32a53da6b1` | `PEOPLE@ETH#25926528` | 5000 | 0 | on-chain balance != csv claim by 5000 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x26738cb185d888bb8c77ef0821dae410d0a748de` | `USDT-OPTIMISM@OPTIMISM#156599011` | 5169.533166 | 0 | on-chain balance != csv claim by 5169.533166 (provider=https://mainnet.optimism.io) |
| `0x7cfb1543e6dec0dadd2435f23b97e08d682b6241` | `USDT-ARBITRUM@ARBITRUM#502728533` | 5447.282449 | 0 | on-chain balance != csv claim by 5447.282449 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x0d96b0525186c56aab001ff856313c85ced6bbe3` | `PEOPLE@ETH#25926528` | 4822.62 | 0 | on-chain balance != csv claim by 4822.62 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0145506e5169a5efcbfc048636553d532108d668` | `PEOPLE@ETH#25926528` | 3857.1065 | 0 | on-chain balance != csv claim by 3857.1065 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x8dbeda7085ddb63fd91bdfb8ba2c69a9a06c7589` | `PEOPLE@ETH#25926528` | 3716.7398 | 0 | on-chain balance != csv claim by 3716.7398 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd56a81842033ac2ab3e958d357501de052bcf7bc` | `PEOPLE@ETH#25926528` | 3637.76547135 | 0 | on-chain balance != csv claim by 3637.76547135 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4e3dd9ec3047246a5fd849334438696df4851b97` | `PEOPLE@ETH#25926528` | 3530.58365806 | 0 | on-chain balance != csv claim by 3530.58365806 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x03c1f15c3f2461d6b53cd7ab799b95b6d2e669d0` | `PEOPLE@ETH#25926528` | 3500 | 0 | on-chain balance != csv claim by 3500 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6671c4d6bbdff00ef9aa4ad9488d400a36272cf6` | `PEOPLE@ETH#25926528` | 3236.30079953 | 0 | on-chain balance != csv claim by 3236.30079953 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x1caf4162b5c0d254b4d07809dba8abf34b895b86` | `PEOPLE@ETH#25926528` | 3199.16248961 | 0 | on-chain balance != csv claim by 3199.16248961 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xcc8e49f620f47dee27dd08caa8dd3554302388ea` | `PEOPLE@ETH#25926528` | 3036.885 | 0 | on-chain balance != csv claim by 3036.885 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xe57245fb0e5fc8d3cb02fb05e1cb0a8939970d6d` | `USDT-OPTIMISM@OPTIMISM#156599011` | 2959.276143 | 0 | on-chain balance != csv claim by 2959.276143 (provider=https://mainnet.optimism.io) |
| `0x6fb31ae81a72af0ce3acd41d5eefe57d2758f72c` | `PEOPLE@ETH#25926528` | 2863.17572044 | 0 | on-chain balance != csv claim by 2863.17572044 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4b4b37f62437ba4e669c418f261176c2982bc149` | `USDT-ARBITRUM@ARBITRUM#502728533` | 3150 | 0 | on-chain balance != csv claim by 3150 (provider=https://arb1.arbitrum.io/rpc) |
| `0x85af37f0f445e09302bb1f4d49766555d83b49d3` | `PEOPLE@ETH#25926528` | 2814.71191813 | 0 | on-chain balance != csv claim by 2814.71191813 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x8505bc2f4e44844db42cd9935905fbdd6a96d897` | `PEOPLE@ETH#25926528` | 2690.7088 | 0 | on-chain balance != csv claim by 2690.7088 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x45d6e0a10d9a58bba0df2f8f8b66c7b447451bc7` | `PEOPLE@ETH#25926528` | 2578.16279 | 0 | on-chain balance != csv claim by 2578.16279 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd082b5c3b274d4433b2353b927d91a1082e5d1d1` | `PEOPLE@ETH#25926528` | 2459.19770969 | 0 | on-chain balance != csv claim by 2459.19770969 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xed1808cc19a7c8ea66319e560a652b427038b939` | `PEOPLE@ETH#25926528` | 2206.5973858 | 0 | on-chain balance != csv claim by 2206.5973858 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xc83a70f90970cfed8774e2aac3025cacb872e9a3` | `PEOPLE@ETH#25926528` | 2200 | 0 | on-chain balance != csv claim by 2200 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6290a8e42010d4a68c8428920d5113300bdeaa40` | `PEOPLE@ETH#25926528` | 2162.286628 | 0 | on-chain balance != csv claim by 2162.286628 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xbb96e1253431f8d0523cd8337ca85ce47f01e5e2` | `PEOPLE@ETH#25926528` | 2107.4 | 0 | on-chain balance != csv claim by 2107.4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xb10cd07b0fbec13a9bc1e89e57e51e00fd26586c` | `PEOPLE@ETH#25926528` | 2069.48525718 | 0 | on-chain balance != csv claim by 2069.48525718 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x02a899128336dcc1cd087386292153e01de3848a` | `PEOPLE@ETH#25926528` | 1900.10162555 | 0 | on-chain balance != csv claim by 1900.10162555 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6682018cc6f3e6138ea8e63d6b87b2a48f445c1f` | `PEOPLE@ETH#25926528` | 1893.63922039 | 0 | on-chain balance != csv claim by 1893.63922039 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `OKB@ETH#25926528` | 1887.741668 | 0 | on-chain balance != csv claim by 1887.741668 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x23d90deb6252846090f83e4e39555199bcc3b729` | `POLY-USDC@MATIC#93394510` | 1863.62079 | 0 | on-chain balance != csv claim by 1863.62079 (provider=https://polygon.gateway.tenderly.co) |
| `0x922a2e84faf1b264182d92c2c5eff013547d02e8` | `PEOPLE@ETH#25926528` | 1817.25 | 0 | on-chain balance != csv claim by 1817.25 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0591199632c24da82ec82d853efcc960347d255c` | `USDT-ARBITRUM@ARBITRUM#502728533` | 1498.626359 | 0 | on-chain balance != csv claim by 1498.626359 (provider=https://arb1.arbitrum.io/rpc) |
| `0x942fd3a57d045b83f1e8f1d0381b496950cc16d3` | `USDC-AVAXC@AVAXC#94703902` | 2583.587644 | 0 | on-chain balance != csv claim by 2583.587644 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x34bba0d0a1d0cab16ef6be267465efe5c4df201a` | `USDT-ARBITRUM@ARBITRUM#502728533` | 1439.092055 | 0 | on-chain balance != csv claim by 1439.092055 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x03426aacbcf9030a380b3ba4a1f16b968649ef3d` | `USDT-ARBITRUM@ARBITRUM#502728533` | 1058.168507 | 0 | on-chain balance != csv claim by 1058.168507 (provider=https://arb1.arbitrum.io/rpc) |
| `0xa8dc30904dc5132f1b86f0a2124414e2f75dfaa0` | `USDT-ARBITRUM@ARBITRUM#502728533` | 848.02478 | 0 | on-chain balance != csv claim by 848.02478 (provider=https://arb1.arbitrum.io/rpc) |
| `0x10e7d149e73dae219bb517de0fcb6a9601ba0f02` | `POLY-USDC@MATIC#93394510` | 839.787582 | 50 | on-chain balance != csv claim by 789.787582 (provider=https://polygon.gateway.tenderly.co) |
| `0x1f1ca79462ec79445c6bd10b9225e1ad2ab9f3d1` | `USDC-AVAXC@AVAXC#94703902` | 1871.738361 | 0 | on-chain balance != csv claim by 1871.738361 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf7f767ec932a506efce0b644a2344f628f804bfd` | `USDT-ARBITRUM@ARBITRUM#502728533` | 740.086967 | 0 | on-chain balance != csv claim by 740.086967 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xc8802feab2fafb48b7d1ade77e197002c210f391` | `USDC-AVAXC@AVAXC#94703902` | 1683.735722 | 0 | on-chain balance != csv claim by 1683.735722 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xabf8a176e8077047cc009263e89081f1a7b6bff6` | `POLY-USDC@MATIC#93394510` | 479.584999 | 0 | on-chain balance != csv claim by 479.584999 (provider=https://polygon.gateway.tenderly.co) |
| `0xc9b00b361785f3af5fd1027554ff1509923691b6` | `USDT-OPTIMISM@OPTIMISM#156599011` | 473.038988 | 0 | on-chain balance != csv claim by 473.038988 (provider=https://mainnet.optimism.io) |
| `0x000000000000d5775ff7721cefb8097af62e52dd` | `POLY-USDC@MATIC#93394510` | 459.213822 | 0 | on-chain balance != csv claim by 459.213822 (provider=https://polygon.gateway.tenderly.co) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `OKB@ETH#25926528` | 441.4124507897 | 0 | on-chain balance != csv claim by 441.4124507897 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `OKB@ETH#25926528` | 233.510185 | 0 | on-chain balance != csv claim by 233.510185 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `OKB@ETH#25926528` | 185.3889805357 | 0 | on-chain balance != csv claim by 185.3889805357 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6f7ba0eb23c8fd9c2601ce6c2248a037c732fa40` | `POLY-USDC@MATIC#93394510` | 151.526887 | 0 | on-chain balance != csv claim by 151.526887 (provider=https://polygon.gateway.tenderly.co) |
| `0xb8b70cf4edd49aa6c7212286917bcefb42da01c3` | `USDT-ARBITRUM@ARBITRUM#502728533` | 148.917 | 0 | on-chain balance != csv claim by 148.917 (provider=https://arb1.arbitrum.io/rpc) |
| `0x8c1dd7385efdd54b96dd417ec192f8dc17a13f92` | `USDC-AVAXC@AVAXC#94703902` | 439.504105 | 0 | on-chain balance != csv claim by 439.504105 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe995fe7dc72f7a817ca38c70004fa5388e7f4c2e` | `USDT-OPTIMISM@OPTIMISM#156599011` | 114.71218 | 0 | on-chain balance != csv claim by 114.71218 (provider=https://mainnet.optimism.io) |
| `0x81595bfaf0acf2ca9ba357676a8743ce1ec1c67e` | `USDT-OPTIMISM@OPTIMISM#156599011` | 99.96 | 0 | on-chain balance != csv claim by 99.96 (provider=https://mainnet.optimism.io) |
| `0xf28143f7b282e59bd5f979012982e7cb9d9b95b0` | `USDT-ARBITRUM@ARBITRUM#502728533` | 101.250523 | 0 | on-chain balance != csv claim by 101.250523 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xd953e4717c4aa75b08052645379cec64d539ea5a` | `USDT-ARBITRUM@ARBITRUM#502728533` | 96.207252 | 0 | on-chain balance != csv claim by 96.207252 (provider=https://arb1.arbitrum.io/rpc) |
| `0xdc80e4412b4c6817318b3450b98f87f745206e0b` | `USDC-AVAXC@AVAXC#94703902` | 151.524451 | 0 | on-chain balance != csv claim by 151.524451 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x5a741eab50d90a955392fcb06d77e81fccefdb3d` | `USDT-OPTIMISM@OPTIMISM#156599011` | 91.009573 | 0 | on-chain balance != csv claim by 91.009573 (provider=https://mainnet.optimism.io) |
| `0x011935ed307ef3865dc616bf5f1923f545e0fc97` | `USDT-ARBITRUM@ARBITRUM#502728533` | 93.531021 | 0 | on-chain balance != csv claim by 93.531021 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xfdbc535ef5bc2b99e4e35b0b053422cc0d1dd389` | `USDT-OPTIMISM@OPTIMISM#156599011` | 90.546289 | 0 | on-chain balance != csv claim by 90.546289 (provider=https://mainnet.optimism.io) |
| `0xc7f758b2126a33239fb7e80f91c722caf7126c15` | `USDT-OPTIMISM@OPTIMISM#156599011` | 89.991773 | 0 | on-chain balance != csv claim by 89.991773 (provider=https://mainnet.optimism.io) |
| `0x62383739d68dd0f844103db8dfb05a7eded5bbe6` | `USDC-AVAXC@AVAXC#94703902` | 123.402721 | 0 | on-chain balance != csv claim by 123.402721 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x4a39790a5afa234a5b241ec052cbb1573d709c6f` | `USDT-ARBITRUM@ARBITRUM#502728533` | 89.9666 | 0 | on-chain balance != csv claim by 89.9666 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xc496e5e40f18035d6f8e805f1e5595657d7dca9f` | `USDT-ARBITRUM@ARBITRUM#502728533` | 89.481355 | 0 | on-chain balance != csv claim by 89.481355 (provider=https://arb1.arbitrum.io/rpc) |
| `0x0938c63109801ee4243a487ab84dffa2bba4589e` | `USDT-ARBITRUM@ARBITRUM#502728533` | 78.397627 | 0 | on-chain balance != csv claim by 78.397627 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x712d8412cf5bed5e4b08f1ae86e6506eb5a4a686` | `USDT-ARBITRUM@ARBITRUM#502728533` | 75.2147 | 0 | on-chain balance != csv claim by 75.2147 (provider=https://arb1.arbitrum.io/rpc) |
| `0x5ebfc3de7cef8a26775b15913130831ee9159109` | `USDT-OPTIMISM@OPTIMISM#156599011` | 66.194211 | 0 | on-chain balance != csv claim by 66.194211 (provider=https://mainnet.optimism.io) |
| `0xc3520b4a68f7fd5e0832e90187a8fad3d7d3c880` | `USDT-ARBITRUM@ARBITRUM#502728533` | 66.112069 | 0 | on-chain balance != csv claim by 66.112069 (provider=https://arb1.arbitrum.io/rpc) |
| `0xa9ac43f5b5e38155a288d1a01d2cbc4478e14573` | `OKB@ETH#25926528` | 35.38 | 0 | on-chain balance != csv claim by 35.38 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `OKB@ETH#25926528` | 34.5644520006 | 0 | on-chain balance != csv claim by 34.5644520006 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `OKB@ETH#25926528` | 18.95 | 0 | on-chain balance != csv claim by 18.95 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `OKB@ETH#25926528` | 14.607403 | 0 | on-chain balance != csv claim by 14.607403 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `OKB@ETH#25926528` | 8.2962431535 | 0 | on-chain balance != csv claim by 8.2962431535 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6f4b48aa129d873507770cdccb25bd60daaf09aa` | `USDC-AVAXC@AVAXC#94703902` | 28.255737 | 0 | on-chain balance != csv claim by 28.255737 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xb0473ca01e91981623b65daf84ab7bc94a7e4000` | `USDC-AVAXC@AVAXC#94703902` | 24.9968 | 0 | on-chain balance != csv claim by 24.9968 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x7970901be38fc9235a5f650aa3708441b378f466` | `USDC-AVAXC@AVAXC#94703902` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf8fbf738259603bb43172cfe4e66953e333dba50` | `USDC-AVAXC@AVAXC#94703902` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x2ad0f758ce754653c79b681750bff16fe005f230` | `USDC-AVAXC@AVAXC#94703902` | 4.999999 | 0 | on-chain balance != csv claim by 4.999999 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe035cb2a4719c0439c81413543d4eee39c0aaeab` | `USDC-AVAXC@AVAXC#94703902` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x920c6944f9e1c5bb7a80d180a8df75f65b0b94e5` | `USDC-AVAXC@AVAXC#94703902` | 4.9997 | 0 | on-chain balance != csv claim by 4.9997 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x136bc54e0950eae85591954ae2bdea2199681234` | `USDC-AVAXC@AVAXC#94703902` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9157b6c4f48112490962175b4058d485939653ca` | `USDC-AVAXC@AVAXC#94703902` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xfe051bc735bc6c69d915a3a035d8e8747fdd8587` | `USDC-AVAXC@AVAXC#94703902` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xc6eca660c5da896784071ebc252fc27184da54cf` | `USDC-AVAXC@AVAXC#94703902` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xceca016c59a63cccdf7a703f52b9dee581e1fa40` | `USDC-AVAXC@AVAXC#94703902` | 4.99968 | 0 | on-chain balance != csv claim by 4.99968 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x1810a491db6f47ea99e31903db4b65f74eaf0b36` | `USDC-AVAXC@AVAXC#94703902` | 4.9996 | 0 | on-chain balance != csv claim by 4.9996 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe270b325c128ac89107ef4a078b303bf81dc7538` | `USDC-AVAXC@AVAXC#94703902` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x381686ce20da524ad9520a157cdf7c871827ce0a` | `USDC-AVAXC@AVAXC#94703902` | 4.999362 | 0 | on-chain balance != csv claim by 4.999362 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd26a7101d2212f4f23c9b8a6af35f7ec432e5c4d` | `USDC-AVAXC@AVAXC#94703902` | 4.999199 | 0 | on-chain balance != csv claim by 4.999199 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xfee1e3ec3577bc256b3c3531de48ea7a8eeff4db` | `USDC-AVAXC@AVAXC#94703902` | 4.99919 | 0 | on-chain balance != csv claim by 4.99919 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd3c8185a0aaf673a5a94da03a2e48311085e6058` | `USDC-AVAXC@AVAXC#94703902` | 4.99936 | 0 | on-chain balance != csv claim by 4.99936 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x53231830fc5b42204b83fd749d39f151f74168a3` | `USDC-AVAXC@AVAXC#94703902` | 4.99909 | 0 | on-chain balance != csv claim by 4.99909 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x62a5b0e0b15f850016783aa8bfc9fa19cdb69208` | `USDC-AVAXC@AVAXC#94703902` | 4.99906 | 0 | on-chain balance != csv claim by 4.99906 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf52ac68c0b89cc62846edba7733f645e17f8f98e` | `USDC-AVAXC@AVAXC#94703902` | 4.998932 | 0 | on-chain balance != csv claim by 4.998932 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9bed615d3e388d91b7302d0710949b1cf5fb6544` | `USDC-AVAXC@AVAXC#94703902` | 4.999081 | 0 | on-chain balance != csv claim by 4.999081 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x5bfed738dd750fee869b2ae816827af9df123055` | `OKB@ETH#25926528` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x352644b25656b1f2503e448f472e77c078e2e354` | `OKB@ETH#25926528` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf7bb83b23869fde84375aaffa036af33e9ba8f96` | `OKB@ETH#25926528` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x030ae17036335c1c8c1b166f419b3f0721a5f184` | `OKB@ETH#25926528` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xa765013bb5fbb4f16e12ede43288f915021fc5d7` | `OKB@ETH#25926528` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x454bfa69931e6f7e27ad630dd37531241cb73394` | `USDC-AVAXC@AVAXC#94703902` | 4.998883 | 0 | on-chain balance != csv claim by 4.998883 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x52aeace7f732f4f93ec47caac1bf0752fa2b09dd` | `USDC-AVAXC@AVAXC#94703902` | 4.998905 | 0 | on-chain balance != csv claim by 4.998905 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9f25df81d05476f5e56309dd95d7c90f47f55452` | `OKB@ETH#25926528` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x36e0c1f02da4aa7c640fd0f3a56e3ab08262b47d` | `OKB@ETH#25926528` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x2801845ac56549789ec24377981333b73b651b2d` | `OKB@ETH#25926528` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9601bd761dd4c9aeee10b5fee5347b36c6733ea1` | `OKB@ETH#25926528` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x68dfbb33c82bd43e6820d5d7bd50e27a72cfa6bf` | `OKB@ETH#25926528` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x26956cafa64673b8c4d552c75b1817b6e1539fda` | `OKB@ETH#25926528` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd53d43e09def98bc4a8bdeb24be58a01dd806d53` | `OKB@ETH#25926528` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x48120be2162bf6d7b1c8ea04a150a530d294dcd0` | `OKB@ETH#25926528` | 1.79601 | 0 | on-chain balance != csv claim by 1.79601 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xbb96e1253431f8d0523cd8337ca85ce47f01e5e2` | `OKB@ETH#25926528` | 1.09 | 0 | on-chain balance != csv claim by 1.09 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x57c6c5159a9eb7107d1caa25ad1fc8228b395cf8` | `OKB@ETH#25926528` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x3329218b3cbac9127e0a990f2cf3fe474e3a8877` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x80335e658ee73756bb8220bdb77f99e98322035f` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x7926c83479e2708edfd83a39d62084693e192298` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4d619eb6db0b9b6184d9881c4828c0effff84eda` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x8ad740adce08f2f3579884eac0a12e581304ac3b` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9ce5a8098270f07a0a56398febda76a8ee2858c4` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9d65bf1bce7a10814228f4d81d8cf99a6b0f44aa` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x7d4a2c448fe5b07c62228140417b3565bbe1b86b` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xc0d01a27fd89e6bbc68f43e94ade3dfca2c27674` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xc6724cf116573eab588fb0d73dcdc16929493cea` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x203094b137651e05c334e119f6143f65595a8420` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9f1cdda0fac4e16f4ee2a844ec6c47add8d55054` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xbfd30f47fd0af30a5147f42154e9f29ad2d94bda` | `OKB@ETH#25926528` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `TCw6YaWm3y6DvxY7M8hrCDnrJGeGMumzGJ` | `TRX@TRX#86041672` | 501754060.133808 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 501754060.133808 (provider=https://api.trongrid.io) |
| `0xa09fe38187abecc2768b98b80ea30afbc2fbd774` | `USDC@ETH#25926528` | 112068568.816273 | 14.316356 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=112068554.499917 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `TLaGjwhvA8XQYSxFAcAXy7Dvuue9eGYitv` | `USDT-TRC20@TRX#86041672` | 170185039.05659 | 88303518.727888 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 81881520.328702 (provider=https://tron-rpc.publicnode.com) |
| `TYE5mNXSQHuPAdCa1etk76BKybD4FKJto6` | `USDT-TRC20@TRX#86041672` | 67627478.839697 | 31313404.0782 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 36314074.761497 (provider=https://api.trongrid.io) |
| `TWAY47W8ci9XtfMgvs4r2A7Z2NFGUDJqpp` | `TRX@TRX#86041672` | 73351863 | 1001357.825295 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 72350505.174705 (provider=https://api.trongrid.io) |
| `0x7fea5b4568751533039179116e372e26b6b41b13` | `USDC@ETH#25926528` | 29338672.815116 | 4.657162 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=29338668.157954 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `THdXEcBCZ7nagnmfRrfSHoiLRTbSRueoVW` | `TRX@TRX#86041672` | 36531497.455287 | 9560158.059565 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 26971339.395722 (provider=https://tron-rpc.publicnode.com) |
| `TBwBJwj81yXc4DNKS19GJcpUUzfSWRbBzS` | `USDT-TRC20@TRX#86041672` | 18798223.830006 | 16215327.201417 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 2582896.628589 (provider=https://tron-rpc.publicnode.com) |
| `TBNHZP2iAVyDhpTSZ8ee2Xjci1JExZTPE5` | `TRX@TRX#86041672` | 18140451.006077 | 40451.336404 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 18099999.669673 (provider=https://tron-rpc.publicnode.com) |
| `TJVeQfDKsV2ZqrDzZKabgoE6De1vJhnP7h` | `TRX@TRX#86041672` | 18165851.405137 | 174790.404107 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 17991061.00103 (provider=https://tron-rpc.publicnode.com) |
| `TBNHZP2iAVyDhpTSZ8ee2Xjci1JExZTPE5` | `USDT-TRC20@TRX#86041672` | 16088358.559474 | 11318402.817686 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 4769955.741788 (provider=https://tron-rpc.publicnode.com) |
| `TEKQdq94xCKNUQwbw6hPMaFCDYoHcdzGzB` | `TRX@TRX#86041672` | 9988133.738557 | 2961.32199 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 9985172.416567 (provider=https://tron-rpc.publicnode.com) |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `OKB-X1@XLAYER` |  |  | rpc error (provider=): network XLAYER not configured |
| `TZ7wp7mcM2Qq4kmDNW3pS7ZZ8Y8ztXUVji` | `TRX@TRX#86041672` | 7216820.001511 | 14205.998897 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 7202614.002614 (provider=https://tron-rpc.publicnode.com) |
| `TEHPn1jYEePZz7v3Fp3VT9N5xDSKS8xTym` | `USDT-TRC20@TRX#86041672` | 6695389.292677 | 0.000032 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 6695389.292645 (provider=https://tron-rpc.publicnode.com) |


_142 additional `WARN` rows omitted; see verification bundle._

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 15, WARN 229, UNVERIFIABLE 1, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `C68a6RCGLiPskbPYtAcsCjhG8tfTWYcoB4JjCrXFdqyo` | `USDG\|SOL#445105348` | 91803081.369402 | 91802961.792143 | on-chain balance != csv claim by 119.577259 (provider=https://api.helius.xyz/v1/wallet) |
| `EU3X4ihDECNrLEwrgkfgppwbGKgTdjQ1RFiEtFdVXExS` | `USDC-SPL\|SOL#445105348` | 6028455.576808 | 6021407.932608 | on-chain balance != csv claim by 7047.6442 (provider=https://api.helius.xyz/v1/wallet) |
| `3WjxkXbu9JBosUWUnJpjzec7Wa4rSkWo9khUbbLDCssk` | `SOL\|SOL#445105348` | 3938217.356102703 | 3.971923082 | on-chain balance != csv claim by 3938213.384179621 (provider=https://api.helius.xyz/v1/wallet) |
| `CKPbSRoW17z8fkFBy5NPcitiLZX4uqmfSv5MpFewkNoQ` | `USDC-SPL\|SOL#445105348` | 687045.5512 | 686984.25029 | on-chain balance != csv claim by 61.30091 (provider=https://api.helius.xyz/v1/wallet) |
| `A2nNTaAiDSYqCykEzXeK2hohAaxMpbpZQN` | `DOGE\|DOGE#6364631` | 510750.001 | 500921.001 | on-chain balance != csv claim by 9829 (provider=https://api.blockcypher.com/v1/doge/main/addrs/A2nNTaAiDSYqCykEzXeK2hohAaxMpbpZQN?before=6364632) |
| `CyKe1nTzB56CLtwgLiJ4WWTDvAjndwTfg4qs8j89k1NU` | `USDC-SPL\|SOL#445105348` | 286918.882107 | 286918.067707 | on-chain balance != csv claim by 0.8144 (provider=https://api.helius.xyz/v1/wallet) |
| `GD3tr9DR8ZnSy2pQ22RwzSR4vKY48SxHLdgCymXQ6W5i` | `SOL\|SOL#445105348` | 262932.32334521 | 0.165106123 | on-chain balance != csv claim by 262932.158239087 (provider=https://api.helius.xyz/v1/wallet) |
| `92QMPBK852QdkGxyyuQZ8BmKczeQihEiGLHbFbx4tVdC` | `USDC-SPL\|SOL#445105348` | 250278.626935 | 250277.814035 | on-chain balance != csv claim by 0.8129 (provider=https://api.helius.xyz/v1/wallet) |
| `Bou5nPZCcqQCNhSaKyJMLhqgoZVKBx6yoDkGEW5WY6W5` | `USDC-SPL\|SOL#445105348` | 259026.204198 | 259025.466798 | on-chain balance != csv claim by 0.7374 (provider=https://api.helius.xyz/v1/wallet) |
| `GHWmH2ZjB6jhv68dzRDvQeC3v23oyxh1W2F4Av6Gj2fy` | `USDC-SPL\|SOL#445105348` | 144481.444886 | 0.2104 | on-chain balance != csv claim by 144481.234486 (provider=https://api.helius.xyz/v1/wallet) |
| `8wn2fZPKmUSMqNKLNsBfb3KkrVhChYYrL5Kzyfc2Dqox` | `USDC-SPL\|SOL#445105348` | 140717.3157 | 140716.6049 | on-chain balance != csv claim by 0.7108 (provider=https://api.helius.xyz/v1/wallet) |
| `B6gPEoYfgKDeMsT2Mxkie4SBsW5GMLQyCDQk42WTFDR` | `USDC-SPL\|SOL#445105348` | 52794.607416 | 52793.992916 | on-chain balance != csv claim by 0.6145 (provider=https://api.helius.xyz/v1/wallet) |
| `AisgV457n3GHkXhmqFh91uz3a9LcQi96if2r934WVXQE` | `USDC-SPL\|SOL#445105348` | 46731.36744 | 0.163 | on-chain balance != csv claim by 46731.20444 (provider=https://api.helius.xyz/v1/wallet) |
| `4YWWA6aFe3E9vh7Bs3XqgiHqvD5pBZiuAZqk3K1X5drt` | `USDC-SPL\|SOL#445105348` | 31938.751648 | 31938.711648 | on-chain balance != csv claim by 0.04 (provider=https://api.helius.xyz/v1/wallet) |
| `3n8t1n1NF6rhDeAPqQJF7RDD6JLGJPQftF1gMEKd4Sfs` | `SOL\|SOL#445105348` | 11093.07244576 | 0.165037137 | on-chain balance != csv claim by 11092.907408623 (provider=https://api.helius.xyz/v1/wallet) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `f1mfxobliwdamau2x25t4mcegzhfln5rwng4aemqi` | `FIL\|FIL#6349676` | 72472753.665746784192885828 | 71938987.044683732279470548 | live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes; on-chain balance != csv claim by 533766.62106305191341528 (provider=https://api.node.glif.io/rpc/v0) |
| `0x51c6abe562e755582d268340b2cf0e2d8895a155dc9b7a7fb5465000d62d770b` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x51c6abe562e755582d268340b2cf0e2d8895a155dc9b7a7fb5465000d62d770b/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `0xd0c23b2613e9c14f4888ff6c4359f197424ca88cf2970bac9e862e5ee65d5fa8` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd0c23b2613e9c14f4888ff6c4359f197424ca88cf2970bac9e862e5ee65d5fa8/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `0xffd7f0d7b24ba690657dd2aa7ccb754c4be938e4fb970c72eef966d664bd1a45` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xffd7f0d7b24ba690657dd2aa7ccb754c4be938e4fb970c72eef966d664bd1a45/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572` | `APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `EQDbf_PU1DLAcneMzWjjUbK_Vc7oahP0iOpuHbxIwDfaI5Rm` | `TONCOIN-NEW\|TON#91210105` | 8113508.797494951 | 12.191071225 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 8113496.606423726 (provider=https://tonapi.io/v2/accounts/EQDbf_PU1DLAcneMzWjjUbK_Vc7oahP0iOpuHbxIwDfaI5Rm) |
| `0x612a398b1320b877cbc84fee033e88aa9e06cfa630b9d56d2b469ea7750bf157` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x612a398b1320b877cbc84fee033e88aa9e06cfa630b9d56d2b469ea7750bf157/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `UQAWBO4_6eixoz9q8Shzck9G2zOiu2S3p9T8jlgWYTpKBqB2` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQAWBO4_6eixoz9q8Shzck9G2zOiu2S3p9T8jlgWYTpKBqB2/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `0x2af5827f7017a7f6127184d73b2340f12e5dc248a48f48e2fe3fff08f6007beb` | `APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x2af5827f7017a7f6127184d73b2340f12e5dc248a48f48e2fe3fff08f6007beb/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `f1iesbg2vfubsunwyjjcvaeh3wuxoi4neyatprhmq` | `FIL\|FIL#6349676` | 3433456.331782142548838724 | 2787236.512524495142029556 | live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes; on-chain balance != csv claim by 646219.819257647406809168 (provider=https://api.node.glif.io/rpc/v0) |
| `0xdf7c04c9bebf4b35bd8c66a92469f0b66cf77ce8586b6262709a897659f4e772` | `USDC-SUI\|SUI` |  |  | ledger rpc error (provider=https://fullnode.mainnet.sui.io:443): sui: Method not found. JSON-RPC on public fullnodes has been deprecated. Please migrate to gRPC or GraphQL endpoints. See https://docs.sui.io/develop/accessing-data/json-rpc-migration for more information. |
| `0x581bdd91ed2b08496cbab0d6b02674f8da3a2f128b086a3209d23ebafc5df553` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x581bdd91ed2b08496cbab0d6b02674f8da3a2f128b086a3209d23ebafc5df553/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |
| `UQDYa6Lp99YwuDrDmCsngxs5dKzL_P7epA4iP0T7dsoY1po4` | `TONCOIN-NEW\|TON#91210105` | 2677209.521003911 | 1133816.241107775 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 1543393.279896136 (provider=https://tonapi.io/v2/accounts/UQDYa6Lp99YwuDrDmCsngxs5dKzL_P7epA4iP0T7dsoY1po4) |
| `0x781fa312385e11c39c0c901f7d76d2892e32d272d09fad90f728656a542c2995` | `APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x781fa312385e11c39c0c901f7d76d2892e32d272d09fad90f728656a542c2995/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1028481988): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1028481988) has been pruned","hint":"The requested data has bee |


_214 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `ADNbM5fBujCRBW1vqezNeAWmnsLp19ki3n` | `DOGE\|DOGE#6364631` | 441101126.77242735 | 441369802.48207348 | address has 223762 txs; exceeds public API limit 200000 — set BLOCKCHAIR_API_KEY or use archive provider |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 502299735 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-09-07T16:00:00Z |
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
