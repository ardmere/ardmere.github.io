# OKX 502260945 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [502260945-assessment.json](./502260945-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `502260945` |
| Snapshot time | `2026-10-02T16:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

OKX audit 502260945 does not meet all Stage 1 public artifact thresholds in this pipeline.

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
| Observed cadence | `~55 days between snapshots` |
| History available | `4 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for okx.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `4e51d11a27988427ad6affdd165b7e8c392c4666297ea0b812306922393c91f2` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `5230bd14fe64f688c5fd480006f714e746a12b6065ac9d683ca02b89852e9e07` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026100300_V4.zip |
| `globalProofBundle` | `6a6da5f8afec8d93b922cca48d04c70b10f5ca2a28944487edf66d8f2b23283d` | https://static.okx.com/cdn/okx/por/merkel/por_502260945_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x3067ccba7bff9f7a4af2cf98615f3ac719939d30e45f46d9bd4b11210c2edd4c` |
| Verification bundle root | `0x088457d82fc40beb795d50dc75dc792fa2f9e45380898884bfba1530ff9375a5` |
| Artifact bundle SHA-256 | `fd8707038f134a0b8d34e4f32b1e16907f8cb37662ac60dd30da55b939400677` |
| Verification bundle SHA-256 | `00c24726c55bdc103e867c394ae91fc31b4199f3106d1866f56df33ef3cfb0df` |

Local bundle paths: [502260945.artifact-bundle.json](../../../artifacts/okx/502260945/bundles/502260945.artifact-bundle.json), [502260945.verification-bundle.v2.json](../../../artifacts/okx/502260945/bundles/502260945.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `FAIL` | `0.9994` | 371/588673 address signatures failed |
| `onchain-balance-hot` | `2.1` | `PARTIAL` | `0.0001` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.0010` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.0005` |  |
| `global-zk-proof` | `okx-1` | `PASS` | `1.0000` | zkSTARKValidator verify-global succeeded; summary merkle root bound to zk proof |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `502260945.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `btc-anchor` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `internal-consistency` (`PASS`)

Finding counts: WARN 12, PASS 22

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `BTC` | `custodyReserveBalances` | 6212.6345421722722054 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `ETH` | `custodyReserveBalances` | 76435.0842110844594599 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDT` | `custodyReserveBalances` | 140455480.2233823973766545 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDC` | `custodyReserveBalances` | 170948528.8094416760670415 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `XRP` | `custodyReserveBalances` | 9987386.9676800585949367 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `DOGE` | `custodyReserveBalances` | 71459950.2071206516013964 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `SOL` | `custodyReserveBalances` | 83913.4943315679429555 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `OKB` | `custodyReserveBalances` | 12900.49316739333 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `LINK` | `custodyReserveBalances` | 134254.8996063100999999 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `LTC` | `custodyReserveBalances` | 49433.5384791958234293 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `UNI` | `custodyReserveBalances` | 49994.1428860879014175 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDG` | `custodyReserveBalances` | 835031.78834801 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |

#### `address-ownership` (`FAIL`)

Summary: 371/588673 address signatures failed

Finding counts: FAIL 5

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x80299387783b6b8c0645a17ea72fdc3cd32e0c9367445fdbe7344742af521fc39afc1fbc3917efc3a4bd72c48814e9b0` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x8063539f96661a8b8b86222b7170ee312c0ecc523b89e3f938698f69d0f82908d3db19259e047afc5c780b12eda785bc` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x806e4c086bc3aff5b9940868fac93e11b139aa04cf749a1adf13e8c6058e861df0426488db4c04d249d26d7d55e9a650` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x8082858d8486e36a4102dffb38cd65ad34aa93b7e81ecca681004508872f4df261ecf7b9ac3e741288f9a0ccbc3a9d0b` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x808c89165c54339155a9447c617c363a6c04fcb31716030e32d4db483ca1698f566c3b87d1f38d46e33e4742f0ef7b89` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |

#### `onchain-balance-hot` (`PARTIAL`)

Finding counts: UNVERIFIABLE 20, PASS 2

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `ASSET-HUB\|DOT` | 60 | no native verifier for this (coin,network) pair yet |
| `A\|VAULTA` | 1 | no native verifier for this (coin,network) pair yet |
| `DOTK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ELF\|AELF` | 272 | no native verifier for this (coin,network) pair yet |
| `ETCK-KIP20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ETC\|ETC` | 41 | no native verifier for this (coin,network) pair yet |
| `ETHK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `FILK-OKC20\|OKC` | 8 | no native verifier for this (coin,network) pair yet |
| `LINKK-OKC20\|OKC` | 3 | no native verifier for this (coin,network) pair yet |
| `OKB-OKC20\|OKC` | 5 | no native verifier for this (coin,network) pair yet |
| `TRXK-KIP20\|OKC` | 14 | no native verifier for this (coin,network) pair yet |
| `USDC-OKC20\|OKC` | 254 | no native verifier for this (coin,network) pair yet |
| `USDT-OKC20\|OKC` | 107 | no native verifier for this (coin,network) pair yet |
| `XRPK-KIP20\|OKC` | 4 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x85e306287d845752a1716bfd5e90f266d08377dd1d5e4189d6f0bc05a9860dac6a0f73d2f6ec7b1458681ba0570afac3` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xb134297b6c0f624afab8d357c9c2f8324dfa844466c35a12d33ebe4959dfe06d88f9544593cf17996bbc0cca4d0a43d9` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x80670ceb430da829ff33566afdecede98d4bcdb1bdbc33acb16425e1cdd29a8088f29bcd7c5ba104171eaf3b1b369291` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x821b54086bedce6be4542660391d66f72779165fd854c59ecb904f3ae4171b17b52e3a8784c326db3de888e188cffb49` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0x8a7fa64cd9d0fcbc52f417d7b05271b151219659da4ecb306385dfb06d049af43783c4662e22e48175cedfc8e8438bf8` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |
| `0xb2d016a313305400b7125cc0c5bd9a6acb001cd640648084d63d49154451a700feca775c17c09468638e517a4ceb42f8` | `ETH@ETH` |  |  | invalid EVM wallet address for ETH@ETH |

#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 160, WARN 237, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `PEOPLE@ETH#26105691` | 1142999298.56695082 | 0 | on-chain balance != csv claim by 1142999298.56695082 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `PEOPLE@ETH#26105691` | 127975207.2551587972 | 0 | on-chain balance != csv claim by 127975207.2551587972 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `PEOPLE@ETH#26105691` | 106905378.7075217122 | 0 | on-chain balance != csv claim by 106905378.7075217122 (provider=https://eth.drpc.org) |
| `0x85dcd76d4fbd3aa0c85c27b9441222c19a14134b` | `PEOPLE@ETH#26105691` | 7409856.36392 | 0 | on-chain balance != csv claim by 7409856.36392 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xafee421482faea92292ed3ffe29371742542ad72` | `USDT-ARBITRUM@ARBITRUM#511033106` | 9397013.078172 | 0 | on-chain balance != csv claim by 9397013.078172 (provider=https://arb1.arbitrum.io/rpc) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `PEOPLE@ETH#26105691` | 5512347.1922450403 | 0 | on-chain balance != csv claim by 5512347.1922450403 (provider=https://eth.drpc.org) |
| `0x92ea7496eba5f001d620005f88f3e8e686e3d4ea` | `PEOPLE@ETH#26105691` | 4364549.34 | 0 | on-chain balance != csv claim by 4364549.34 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd8a3ed83a59d3e10d82a4b845d64b24ba0572baf` | `USDT-ARBITRUM@ARBITRUM#511033106` | 4210578.839134 | 0 | on-chain balance != csv claim by 4210578.839134 (provider=https://arb1.arbitrum.io/rpc) |
| `0xe7c5e276b2434743f6c79365d967cbe4291477bc` | `USDT-ARBITRUM@ARBITRUM#511033106` | 2885234.815313 | 0 | on-chain balance != csv claim by 2885234.815313 (provider=https://arb1.arbitrum.io/rpc) |
| `0x4a1ca3e90e9e3a00eb27e9f409b4db310aa7405f` | `USDT-ARBITRUM@ARBITRUM#511033106` | 1312475.044963 | 0 | on-chain balance != csv claim by 1312475.044963 (provider=https://arb1.arbitrum.io/rpc) |
| `0x3f03440b19d15ff2068daf68e57eb656bfd0a67e` | `USDC-AVAXC@AVAXC#96612935` | 1266015.632136 | 0 | on-chain balance != csv claim by 1266015.632136 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xa0c9b4ce1785953abdceae569789aa158a01da4c` | `USDT-ARBITRUM@ARBITRUM#511033106` | 1154991.054464 | 0 | on-chain balance != csv claim by 1154991.054464 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x343d752bb710c5575e417edb3f9fa06241a4749a` | `POLY-USDC@MATIC#94834439` | 1027589.45663 | 46.515055 | on-chain balance != csv claim by 1027542.941575 (provider=https://polygon.gateway.tenderly.co) |
| `0xde3c8d726e906b379a020cdb04bbe4b8f39a726f` | `USDC-AVAXC@AVAXC#96612935` | 1017398.531787 | 0 | on-chain balance != csv claim by 1017398.531787 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xfc9161b99a3b0691477e12737163d007600ac262` | `PEOPLE@ETH#26105691` | 879971.59990868 | 0 | on-chain balance != csv claim by 879971.59990868 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `PEOPLE@ETH#26105691` | 852992.7652159399 | 0 | on-chain balance != csv claim by 852992.7652159399 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9d1964437bb6d849c66ab2a2a033340dc2d999dc` | `USDC-AVAXC@AVAXC#96612935` | 963145.192222 | 0 | on-chain balance != csv claim by 963145.192222 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `PEOPLE@ETH#26105691` | 486437.6318074999 | 0 | on-chain balance != csv claim by 486437.6318074999 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `PEOPLE@ETH#26105691` | 445366.72893191 | 0 | on-chain balance != csv claim by 445366.72893191 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf418263a6468c80d3bcb71033c56aacf81b241ba` | `PEOPLE@ETH#26105691` | 257337.48 | 0 | on-chain balance != csv claim by 257337.48 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4cafbe2b5bd89f0c0b7d6fef75b47e7c2ae82f22` | `USDC-AVAXC@AVAXC#96612935` | 293773.655154 | 0 | on-chain balance != csv claim by 293773.655154 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xb1da6e546bd592218d6e849a8177d9d4fbf12986` | `POLY-USDC@MATIC#94834439` | 198609.106952 | 0 | on-chain balance != csv claim by 198609.106952 (provider=https://polygon.gateway.tenderly.co) |
| `0xc74e4c556a16390165c99b40aabc39a87c358305` | `POLY-USDC@MATIC#94834439` | 204532.591279 | 0 | on-chain balance != csv claim by 204532.591279 (provider=https://polygon.gateway.tenderly.co) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `PEOPLE@ETH#26105691` | 123276.91 | 0 | on-chain balance != csv claim by 123276.91 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xcf67640fd0716eb9c0d18c32214a36b8d8a2abaf` | `USDT-ARBITRUM@ARBITRUM#511033106` | 158374.459791 | 0 | on-chain balance != csv claim by 158374.459791 (provider=https://arb1.arbitrum.io/rpc) |
| `0x229b1c70124141ec96945aec2323a22b68bb093c` | `USDT-ARBITRUM@ARBITRUM#511033106` | 139389.702233 | 0 | on-chain balance != csv claim by 139389.702233 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x5ca1063894e5bbd688830d83fefa5ce3d339744d` | `USDT-ARBITRUM@ARBITRUM#511033106` | 118192.045139 | 0 | on-chain balance != csv claim by 118192.045139 (provider=https://arb1.arbitrum.io/rpc) |
| `0xbe9ac2d59743d69c4054bab09fa9821929d0f31b` | `USDT-ARBITRUM@ARBITRUM#511033106` | 112846.065793 | 0 | on-chain balance != csv claim by 112846.065793 (provider=https://arb1.arbitrum.io/rpc) |
| `0x9c308b1665097af7ab2e38a41577bfecdc2d1c5c` | `PEOPLE@ETH#26105691` | 87146.24000091 | 0 | on-chain balance != csv claim by 87146.24000091 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `OKB@ETH#26105691` | 83487.0654274097 | 0 | on-chain balance != csv claim by 83487.0654274097 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x2fb2841c747621a7187b6f2a5390ae77321e5257` | `USDC-AVAXC@AVAXC#96612935` | 102225.223961 | 0 | on-chain balance != csv claim by 102225.223961 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `PEOPLE@ETH#26105691` | 64579.22 | 0 | on-chain balance != csv claim by 64579.22 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf1e7b58ccd200fcb9400ddbc30867c43b7dc9ee1` | `USDT-ARBITRUM@ARBITRUM#511033106` | 76637.300513 | 0 | on-chain balance != csv claim by 76637.300513 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x3b50287cf8f073ff443b84de57e2375d446b962b` | `POLY-USDC@MATIC#94834439` | 56542.582567 | 0 | on-chain balance != csv claim by 56542.582567 (provider=https://polygon.gateway.tenderly.co) |
| `0x04ba1b380b1418007c1407bbe12ac21a596ae1e0` | `USDC-AVAXC@AVAXC#96612935` | 55589.556619 | 0 | on-chain balance != csv claim by 55589.556619 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x582988097ef9f2cf50cfac3f5f3b1ed54bc52a88` | `USDT-ARBITRUM@ARBITRUM#511033106` | 49669.675025 | 0 | on-chain balance != csv claim by 49669.675025 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x4a152242605ab21df5ac6e819fdb128a4efc33fc` | `POLY-USDC@MATIC#94834439` | 41593.500032 | 0 | on-chain balance != csv claim by 41593.500032 (provider=https://polygon.gateway.tenderly.co) |
| `0x7d328f5d5f4640f81b2f93d2903ec0a6c030eed6` | `PEOPLE@ETH#26105691` | 40971.83003613 | 0 | on-chain balance != csv claim by 40971.83003613 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6bea9e00a6bb63d22904d975bdf8a98f26f9f49d` | `USDC-AVAXC@AVAXC#96612935` | 41593.500081 | 0 | on-chain balance != csv claim by 41593.500081 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x99860ec99aae9c0aa5ab763ee73da1217a9da665` | `POLY-USDC@MATIC#94834439` | 38563.410765 | 0 | on-chain balance != csv claim by 38563.410765 (provider=https://polygon.gateway.tenderly.co) |
| `0x85f56cd748379755494f1365a3efc19569973e99` | `USDC-AVAXC@AVAXC#96612935` | 41451.552004 | 0 | on-chain balance != csv claim by 41451.552004 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `PEOPLE@ETH#26105691` | 34517.45259804 | 0 | on-chain balance != csv claim by 34517.45259804 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xe7f6a87b1950b687a31d068bd01c6be81bc7e1af` | `POLY-USDC@MATIC#94834439` | 35117.783539 | 0 | on-chain balance != csv claim by 35117.783539 (provider=https://polygon.gateway.tenderly.co) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `PEOPLE@ETH#26105691` | 31111.57 | 0 | on-chain balance != csv claim by 31111.57 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `OKB@ETH#26105691` | 30823.2491351555 | 0 | on-chain balance != csv claim by 30823.2491351555 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x1cdf5ba7297f276bf0320c88ca948d4311c913ce` | `USDT-ARBITRUM@ARBITRUM#511033106` | 24433.251864 | 0 | on-chain balance != csv claim by 24433.251864 (provider=https://arb1.arbitrum.io/rpc) |
| `0x2ef3ddc24e4020173da2c0f3a148eff4a9f6d9dd` | `USDC-AVAXC@AVAXC#96612935` | 25880.39594 | 0 | on-chain balance != csv claim by 25880.39594 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x012bcc0b88e1ebe69ec29713f0a8f94d012a9de7` | `POLY-USDC@MATIC#94834439` | 17325.655352 | 0 | on-chain balance != csv claim by 17325.655352 (provider=https://polygon.gateway.tenderly.co) |
| `0x41ddf01de355547d821873ebd4fb2ffa2dc63d4c` | `POLY-USDC@MATIC#94834439` | 10378.213796 | 0.00009 | on-chain balance != csv claim by 10378.213706 (provider=https://polygon.gateway.tenderly.co) |
| `0x356a5c73aa2378a9ec63cfa889f2954095d738bc` | `USDT-ARBITRUM@ARBITRUM#511033106` | 11005.9539 | 0 | on-chain balance != csv claim by 11005.9539 (provider=https://arb1.arbitrum.io/rpc) |
| `0xf4cc92fcc8f41cb52df0c40b89b9c6188b1ae4f5` | `USDT-ARBITRUM@ARBITRUM#511033106` | 10000 | 0 | on-chain balance != csv claim by 10000 (provider=https://arb1.arbitrum.io/rpc) |
| `0x7cfb1543e6dec0dadd2435f23b97e08d682b6241` | `USDT-ARBITRUM@ARBITRUM#511033106` | 10011.343634 | 0 | on-chain balance != csv claim by 10011.343634 (provider=https://arb1.arbitrum.io/rpc) |
| `0x6cc5f688a315f3dc28a7781717a9a798a59fda7b` | `OKB@ETH#26105691` | 11302.306303002590174878 | 0 | on-chain balance != csv claim by 11302.306303002590174878 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x00000000004e3d5628234f18b977041e5242651f` | `POLY-USDC@MATIC#94834439` | 8106.253163 | 0 | on-chain balance != csv claim by 8106.253163 (provider=https://polygon.gateway.tenderly.co) |
| `0x0d96b0525186c56aab001ff856313c85ced6bbe3` | `PEOPLE@ETH#26105691` | 4772.33 | 0 | on-chain balance != csv claim by 4772.33 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0f51a310a4dd79d373eb8be1c0ddd54570235443` | `POLY-USDC@MATIC#94834439` | 5210.792271 | 0 | on-chain balance != csv claim by 5210.792271 (provider=https://polygon.gateway.tenderly.co) |
| `0x3f41ea66a734db98720512df6bd12dbc87c81d30` | `PEOPLE@ETH#26105691` | 4645.3 | 0 | on-chain balance != csv claim by 4645.3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0145506e5169a5efcbfc048636553d532108d668` | `PEOPLE@ETH#26105691` | 3857.1065 | 0 | on-chain balance != csv claim by 3857.1065 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x8dbeda7085ddb63fd91bdfb8ba2c69a9a06c7589` | `PEOPLE@ETH#26105691` | 3716.7398 | 0 | on-chain balance != csv claim by 3716.7398 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x74f42a5bc01eaeaf724da49254da739c37b5a808` | `USDC-AVAXC@AVAXC#96612935` | 4341.998865 | 0 | on-chain balance != csv claim by 4341.998865 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd56a81842033ac2ab3e958d357501de052bcf7bc` | `PEOPLE@ETH#26105691` | 3637.76547135 | 0 | on-chain balance != csv claim by 3637.76547135 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6e1a0e4b46e66ed7aca551ee498729a4956b5a5d` | `PEOPLE@ETH#26105691` | 3547.123523170114121657 | 0 | on-chain balance != csv claim by 3547.123523170114121657 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x942fd3a57d045b83f1e8f1d0381b496950cc16d3` | `USDC-AVAXC@AVAXC#96612935` | 3933.710315 | 0 | on-chain balance != csv claim by 3933.710315 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x4e3dd9ec3047246a5fd849334438696df4851b97` | `PEOPLE@ETH#26105691` | 3530.58365806 | 0 | on-chain balance != csv claim by 3530.58365806 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x03c1f15c3f2461d6b53cd7ab799b95b6d2e669d0` | `PEOPLE@ETH#26105691` | 3500 | 0 | on-chain balance != csv claim by 3500 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x1caf4162b5c0d254b4d07809dba8abf34b895b86` | `PEOPLE@ETH#26105691` | 3199.16248961 | 0 | on-chain balance != csv claim by 3199.16248961 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x03426aacbcf9030a380b3ba4a1f16b968649ef3d` | `USDT-ARBITRUM@ARBITRUM#511033106` | 3533.698617 | 0 | on-chain balance != csv claim by 3533.698617 (provider=https://arb1.arbitrum.io/rpc) |
| `0xcc8e49f620f47dee27dd08caa8dd3554302388ea` | `PEOPLE@ETH#26105691` | 3036.885 | 0 | on-chain balance != csv claim by 3036.885 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6fb31ae81a72af0ce3acd41d5eefe57d2758f72c` | `PEOPLE@ETH#26105691` | 2863.17572044 | 0 | on-chain balance != csv claim by 2863.17572044 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x85af37f0f445e09302bb1f4d49766555d83b49d3` | `PEOPLE@ETH#26105691` | 2814.71191813 | 0 | on-chain balance != csv claim by 2814.71191813 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x66b178421b8198c422026b3a0ebf3176e23120f9` | `USDT-ARBITRUM@ARBITRUM#511033106` | 2999.87 | 0 | on-chain balance != csv claim by 2999.87 (provider=https://arb1.arbitrum.io/rpc) |
| `0x8505bc2f4e44844db42cd9935905fbdd6a96d897` | `PEOPLE@ETH#26105691` | 2690.7088 | 0 | on-chain balance != csv claim by 2690.7088 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x45d6e0a10d9a58bba0df2f8f8b66c7b447451bc7` | `PEOPLE@ETH#26105691` | 2578.16279 | 0 | on-chain balance != csv claim by 2578.16279 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd082b5c3b274d4433b2353b927d91a1082e5d1d1` | `PEOPLE@ETH#26105691` | 2459.19770969 | 0 | on-chain balance != csv claim by 2459.19770969 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x5cf70a6b562a6bc8628e85a3c41e44fc916f26f0` | `PEOPLE@ETH#26105691` | 2220.9778 | 0 | on-chain balance != csv claim by 2220.9778 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6671c4d6bbdff00ef9aa4ad9488d400a36272cf6` | `PEOPLE@ETH#26105691` | 3236.30079953 | 0 | on-chain balance != csv claim by 3236.30079953 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xed1808cc19a7c8ea66319e560a652b427038b939` | `PEOPLE@ETH#26105691` | 2206.5973858 | 0 | on-chain balance != csv claim by 2206.5973858 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xc83a70f90970cfed8774e2aac3025cacb872e9a3` | `PEOPLE@ETH#26105691` | 2200 | 0 | on-chain balance != csv claim by 2200 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6290a8e42010d4a68c8428920d5113300bdeaa40` | `PEOPLE@ETH#26105691` | 2162.286628 | 0 | on-chain balance != csv claim by 2162.286628 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xb10cd07b0fbec13a9bc1e89e57e51e00fd26586c` | `PEOPLE@ETH#26105691` | 2069.48525718 | 0 | on-chain balance != csv claim by 2069.48525718 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x02a899128336dcc1cd087386292153e01de3848a` | `PEOPLE@ETH#26105691` | 1900.10162555 | 0 | on-chain balance != csv claim by 1900.10162555 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x6682018cc6f3e6138ea8e63d6b87b2a48f445c1f` | `PEOPLE@ETH#26105691` | 1893.63922039 | 0 | on-chain balance != csv claim by 1893.63922039 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `OKB@ETH#26105691` | 1887.741668 | 0 | on-chain balance != csv claim by 1887.741668 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x126d0c3f792b99ec831990c5ef23badda9c95f5e` | `USDT-ARBITRUM@ARBITRUM#511033106` | 2099.252946 | 0 | on-chain balance != csv claim by 2099.252946 (provider=https://arb1.arbitrum.io/rpc) |
| `0x34bba0d0a1d0cab16ef6be267465efe5c4df201a` | `USDT-ARBITRUM@ARBITRUM#511033106` | 1439.092055 | 0 | on-chain balance != csv claim by 1439.092055 (provider=https://arb1.arbitrum.io/rpc) |
| `0x04e9d79a17f4894c76c65c52516f790531ef6bee` | `USDT-ARBITRUM@ARBITRUM#511033106` | 1814.633057 | 0 | on-chain balance != csv claim by 1814.633057 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xc8802feab2fafb48b7d1ade77e197002c210f391` | `USDC-AVAXC@AVAXC#96612935` | 1683.735722 | 0 | on-chain balance != csv claim by 1683.735722 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf7f767ec932a506efce0b644a2344f628f804bfd` | `USDT-ARBITRUM@ARBITRUM#511033106` | 1300.908513 | 0 | on-chain balance != csv claim by 1300.908513 (provider=https://arb1.arbitrum.io/rpc) |
| `0x10e7d149e73dae219bb517de0fcb6a9601ba0f02` | `POLY-USDC@MATIC#94834439` | 1262.279791 | 50 | on-chain balance != csv claim by 1212.279791 (provider=https://polygon.gateway.tenderly.co) |
| `0x23d90deb6252846090f83e4e39555199bcc3b729` | `POLY-USDC@MATIC#94834439` | 941.302247 | 0 | on-chain balance != csv claim by 941.302247 (provider=https://polygon.gateway.tenderly.co) |
| `0xb0c3bf12221b8f889278bea9fd6712a284c82fc6` | `USDT-ARBITRUM@ARBITRUM#511033106` | 950 | 0 | on-chain balance != csv claim by 950 (provider=https://arb1.arbitrum.io/rpc) |
| `0x1f1ca79462ec79445c6bd10b9225e1ad2ab9f3d1` | `USDC-AVAXC@AVAXC#96612935` | 949.419859 | 0 | on-chain balance != csv claim by 949.419859 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xb8b70cf4edd49aa6c7212286917bcefb42da01c3` | `USDT-ARBITRUM@ARBITRUM#511033106` | 856.859305 | 0 | on-chain balance != csv claim by 856.859305 (provider=https://arb1.arbitrum.io/rpc) |
| `0x4d15e28516444a40f661b7a7d06c680338539e45` | `USDT-ARBITRUM@ARBITRUM#511033106` | 658 | 0 | on-chain balance != csv claim by 658 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x000000000000d5775ff7721cefb8097af62e52dd` | `POLY-USDC@MATIC#94834439` | 459.213822 | 0 | on-chain balance != csv claim by 459.213822 (provider=https://polygon.gateway.tenderly.co) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `OKB@ETH#26105691` | 443.3135848397 | 0 | on-chain balance != csv claim by 443.3135848397 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xabf8a176e8077047cc009263e89081f1a7b6bff6` | `POLY-USDC@MATIC#94834439` | 434.265954 | 0 | on-chain balance != csv claim by 434.265954 (provider=https://polygon.gateway.tenderly.co) |
| `0x8c1dd7385efdd54b96dd417ec192f8dc17a13f92` | `USDC-AVAXC@AVAXC#96612935` | 401.083218 | 0 | on-chain balance != csv claim by 401.083218 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x6f7ba0eb23c8fd9c2601ce6c2248a037c732fa40` | `POLY-USDC@MATIC#94834439` | 326.990587 | 0 | on-chain balance != csv claim by 326.990587 (provider=https://polygon.gateway.tenderly.co) |
| `0x65a1f26453d752fc423afdf4bf78f1fc21d81dab` | `USDT-ARBITRUM@ARBITRUM#511033106` | 330.3273 | 0 | on-chain balance != csv claim by 330.3273 (provider=https://arb1.arbitrum.io/rpc) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `OKB@ETH#26105691` | 233.510185 | 0 | on-chain balance != csv claim by 233.510185 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xdc80e4412b4c6817318b3450b98f87f745206e0b` | `USDC-AVAXC@AVAXC#96612935` | 326.988151 | 0 | on-chain balance != csv claim by 326.988151 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x08ba649329583c7601999e541202bab5a2255928` | `USDT-ARBITRUM@ARBITRUM#511033106` | 255.3 | 0 | on-chain balance != csv claim by 255.3 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `OKB@ETH#26105691` | 185.3889805357 | 0 | on-chain balance != csv claim by 185.3889805357 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf389d76e4917090933599c418b36f4e52f12ff45` | `USDT-ARBITRUM@ARBITRUM#511033106` | 166.956253 | 0 | on-chain balance != csv claim by 166.956253 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x9ff2433700c9d72b76e2eed5919bb36a7e3e78be` | `USDT-ARBITRUM@ARBITRUM#511033106` | 160.08 | 0 | on-chain balance != csv claim by 160.08 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x766881dc7381d96cc6ef495c53ab702a98d1cdc5` | `USDT-ARBITRUM@ARBITRUM#511033106` | 139.06 | 0 | on-chain balance != csv claim by 139.06 (provider=https://arb1.arbitrum.io/rpc) |
| `0x62383739d68dd0f844103db8dfb05a7eded5bbe6` | `USDC-AVAXC@AVAXC#96612935` | 123.402721 | 0 | on-chain balance != csv claim by 123.402721 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe4dfe1b5392723d40b348a196e851bed1f7c2206` | `USDT-ARBITRUM@ARBITRUM#511033106` | 115.819142 | 0 | on-chain balance != csv claim by 115.819142 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xc496e5e40f18035d6f8e805f1e5595657d7dca9f` | `USDT-ARBITRUM@ARBITRUM#511033106` | 101.786855 | 0 | on-chain balance != csv claim by 101.786855 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xf28143f7b282e59bd5f979012982e7cb9d9b95b0` | `USDT-ARBITRUM@ARBITRUM#511033106` | 101.250523 | 0 | on-chain balance != csv claim by 101.250523 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xd953e4717c4aa75b08052645379cec64d539ea5a` | `USDT-ARBITRUM@ARBITRUM#511033106` | 96.207252 | 0 | on-chain balance != csv claim by 96.207252 (provider=https://arb1.arbitrum.io/rpc) |
| `0x0938c63109801ee4243a487ab84dffa2bba4589e` | `USDT-ARBITRUM@ARBITRUM#511033106` | 78.397627 | 0 | on-chain balance != csv claim by 78.397627 (provider=https://arb1.arbitrum.io/rpc) |
| `0xa9ac43f5b5e38155a288d1a01d2cbc4478e14573` | `OKB@ETH#26105691` | 35.38 | 0 | on-chain balance != csv claim by 35.38 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `OKB@ETH#26105691` | 34.5644520006 | 0 | on-chain balance != csv claim by 34.5644520006 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xaa488c5dd62d17465128c5034179e88a4f2f4170` | `USDC-AVAXC@AVAXC#96612935` | 28.416698 | 0 | on-chain balance != csv claim by 28.416698 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x6f4b48aa129d873507770cdccb25bd60daaf09aa` | `USDC-AVAXC@AVAXC#96612935` | 28.255737 | 0 | on-chain balance != csv claim by 28.255737 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `OKB@ETH#26105691` | 18.95 | 0 | on-chain balance != csv claim by 18.95 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xb0473ca01e91981623b65daf84ab7bc94a7e4000` | `USDC-AVAXC@AVAXC#96612935` | 24.9968 | 0 | on-chain balance != csv claim by 24.9968 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `OKB@ETH#26105691` | 14.607403 | 0 | on-chain balance != csv claim by 14.607403 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `OKB@ETH#26105691` | 8.2962431535 | 0 | on-chain balance != csv claim by 8.2962431535 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x7970901be38fc9235a5f650aa3708441b378f466` | `USDC-AVAXC@AVAXC#96612935` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x2ad0f758ce754653c79b681750bff16fe005f230` | `USDC-AVAXC@AVAXC#96612935` | 4.999999 | 0 | on-chain balance != csv claim by 4.999999 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xa5ec83fa772e0c2610eeaf8b87e6875a0c321af3` | `USDC-AVAXC@AVAXC#96612935` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xfe051bc735bc6c69d915a3a035d8e8747fdd8587` | `USDC-AVAXC@AVAXC#96612935` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe035cb2a4719c0439c81413543d4eee39c0aaeab` | `USDC-AVAXC@AVAXC#96612935` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9157b6c4f48112490962175b4058d485939653ca` | `USDC-AVAXC@AVAXC#96612935` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x136bc54e0950eae85591954ae2bdea2199681234` | `USDC-AVAXC@AVAXC#96612935` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf8fbf738259603bb43172cfe4e66953e333dba50` | `USDC-AVAXC@AVAXC#96612935` | 4.99979 | 0 | on-chain balance != csv claim by 4.99979 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x920c6944f9e1c5bb7a80d180a8df75f65b0b94e5` | `USDC-AVAXC@AVAXC#96612935` | 4.9997 | 0 | on-chain balance != csv claim by 4.9997 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xceca016c59a63cccdf7a703f52b9dee581e1fa40` | `USDC-AVAXC@AVAXC#96612935` | 4.99968 | 0 | on-chain balance != csv claim by 4.99968 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xc6eca660c5da896784071ebc252fc27184da54cf` | `USDC-AVAXC@AVAXC#96612935` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x1810a491db6f47ea99e31903db4b65f74eaf0b36` | `USDC-AVAXC@AVAXC#96612935` | 4.9996 | 0 | on-chain balance != csv claim by 4.9996 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x381686ce20da524ad9520a157cdf7c871827ce0a` | `USDC-AVAXC@AVAXC#96612935` | 4.999362 | 0 | on-chain balance != csv claim by 4.999362 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xe270b325c128ac89107ef4a078b303bf81dc7538` | `USDC-AVAXC@AVAXC#96612935` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd3c8185a0aaf673a5a94da03a2e48311085e6058` | `USDC-AVAXC@AVAXC#96612935` | 4.99936 | 0 | on-chain balance != csv claim by 4.99936 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xd26a7101d2212f4f23c9b8a6af35f7ec432e5c4d` | `USDC-AVAXC@AVAXC#96612935` | 4.999199 | 0 | on-chain balance != csv claim by 4.999199 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x53231830fc5b42204b83fd749d39f151f74168a3` | `USDC-AVAXC@AVAXC#96612935` | 4.99909 | 0 | on-chain balance != csv claim by 4.99909 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xfee1e3ec3577bc256b3c3531de48ea7a8eeff4db` | `USDC-AVAXC@AVAXC#96612935` | 4.99919 | 0 | on-chain balance != csv claim by 4.99919 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9bed615d3e388d91b7302d0710949b1cf5fb6544` | `USDC-AVAXC@AVAXC#96612935` | 4.999081 | 0 | on-chain balance != csv claim by 4.999081 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x62a5b0e0b15f850016783aa8bfc9fa19cdb69208` | `USDC-AVAXC@AVAXC#96612935` | 4.99906 | 0 | on-chain balance != csv claim by 4.99906 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xa765013bb5fbb4f16e12ede43288f915021fc5d7` | `OKB@ETH#26105691` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x5bfed738dd750fee869b2ae816827af9df123055` | `OKB@ETH#26105691` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x030ae17036335c1c8c1b166f419b3f0721a5f184` | `OKB@ETH#26105691` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf7bb83b23869fde84375aaffa036af33e9ba8f96` | `OKB@ETH#26105691` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x352644b25656b1f2503e448f472e77c078e2e354` | `OKB@ETH#26105691` | 4 | 0 | on-chain balance != csv claim by 4 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9f25df81d05476f5e56309dd95d7c90f47f55452` | `OKB@ETH#26105691` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x36e0c1f02da4aa7c640fd0f3a56e3ab08262b47d` | `OKB@ETH#26105691` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x2801845ac56549789ec24377981333b73b651b2d` | `OKB@ETH#26105691` | 3 | 0 | on-chain balance != csv claim by 3 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xd53d43e09def98bc4a8bdeb24be58a01dd806d53` | `OKB@ETH#26105691` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x68dfbb33c82bd43e6820d5d7bd50e27a72cfa6bf` | `OKB@ETH#26105691` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x26956cafa64673b8c4d552c75b1817b6e1539fda` | `OKB@ETH#26105691` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x57c6c5159a9eb7107d1caa25ad1fc8228b395cf8` | `OKB@ETH#26105691` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9601bd761dd4c9aeee10b5fee5347b36c6733ea1` | `OKB@ETH#26105691` | 2 | 0 | on-chain balance != csv claim by 2 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9ce5a8098270f07a0a56398febda76a8ee2858c4` | `OKB@ETH#26105691` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9d65bf1bce7a10814228f4d81d8cf99a6b0f44aa` | `OKB@ETH#26105691` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9f1cdda0fac4e16f4ee2a844ec6c47add8d55054` | `OKB@ETH#26105691` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x19abcd6b8516d2b4d9aa3bbf700b38a30dc154fc` | `OKB@ETH#26105691` | 1 | 0 | on-chain balance != csv claim by 1 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xbb96e1253431f8d0523cd8337ca85ce47f01e5e2` | `OKB@ETH#26105691` | 1.09 | 0 | on-chain balance != csv claim by 1.09 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x48120be2162bf6d7b1c8ea04a150a530d294dcd0` | `OKB@ETH#26105691` | 1.79601 | 0 | on-chain balance != csv claim by 1.79601 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `TCw6YaWm3y6DvxY7M8hrCDnrJGeGMumzGJ` | `TRX@TRX#86761390` | 470817842.286196 | 1464767.281247 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 469353075.004949 (provider=https://tron-rpc.publicnode.com) |
| `TLaGjwhvA8XQYSxFAcAXy7Dvuue9eGYitv` | `USDT-TRC20@TRX#86761390` | 179597523.651807 | 79382956.341429 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 100214567.310378 (provider=https://tron-rpc.publicnode.com) |
| `THwyeoCMSXQq2aGDpXpSqo89VnTtFktGGp` | `TRX@TRX#86761390` | 155224146.696409 | 126497653.924368 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 28726492.772041 (provider=https://tron-rpc.publicnode.com) |
| `TYE5mNXSQHuPAdCa1etk76BKybD4FKJto6` | `USDT-TRC20@TRX#86761390` | 113889658.306949 | 36523892.8007 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 77365765.506249 (provider=https://tron-rpc.publicnode.com) |
| `TWAY47W8ci9XtfMgvs4r2A7Z2NFGUDJqpp` | `TRX@TRX#86761390` | 74624185 | 988165.014013 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 73636019.985987 (provider=https://tron-rpc.publicnode.com) |
| `0xa09fe38187abecc2768b98b80ea30afbc2fbd774` | `USDC@ETH#26105691` | 89829906.642508 | 14.386516 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=89829892.255992 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `TEHPn1jYEePZz7v3Fp3VT9N5xDSKS8xTym` | `USDT-TRC20@TRX#86761390` | 19679235.785248 | 0.000032 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 19679235.785216 (provider=https://tron-rpc.publicnode.com) |
| `TBNHZP2iAVyDhpTSZ8ee2Xjci1JExZTPE5` | `USDT-TRC20@TRX#86761390` | 18383973.521466 | 11218651.68177 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 7165321.839696 (provider=https://tron-rpc.publicnode.com) |
| `0x7fea5b4568751533039179116e372e26b6b41b13` | `USDC@ETH#26105691` | 27095777.876123 | 4.707755 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=27095773.168368 (provider=https://eth.drpc.org) |
| `TBNHZP2iAVyDhpTSZ8ee2Xjci1JExZTPE5` | `TRX@TRX#86761390` | 18140451.012594 | 40451.336353 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 18099999.676241 (provider=https://tron-rpc.publicnode.com) |
| `TJVeQfDKsV2ZqrDzZKabgoE6De1vJhnP7h` | `TRX@TRX#86761390` | 7265355.404672 | 174798.404115 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 7090557.000557 (provider=https://tron-rpc.publicnode.com) |
| `TZ7wp7mcM2Qq4kmDNW3pS7ZZ8Y8ztXUVji` | `TRX@TRX#86761390` | 7215605.000296 | 14250.998942 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 7201354.001354 (provider=https://tron-rpc.publicnode.com) |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `OKB-X1@XLAYER` |  |  | rpc error (provider=): network XLAYER not configured |
| `TVvyKT63p66nqhQiNHX25izSUAvjcdER9X` | `TRX@TRX#86761390` | 7186174.398626 | 5992498.399953 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 1193675.998673 (provider=https://tron-rpc.publicnode.com) |
| `TFAXd82DvBpykbt9JyTpevDZWfkmqy2CnZ` | `TRX@TRX#86761390` | 6357561.000371 | 6317982.074002 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 39578.926369 (provider=https://tron-rpc.publicnode.com) |


_222 additional `WARN` rows omitted; see verification bundle._

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 17, WARN 241, UNVERIFIABLE 2, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `C68a6RCGLiPskbPYtAcsCjhG8tfTWYcoB4JjCrXFdqyo` | `USDG\|SOL#452657748` | 94150609.741376 | 94150490.164117 | on-chain balance != csv claim by 119.577259 (provider=https://api.helius.xyz/v1/wallet) |
| `EU3X4ihDECNrLEwrgkfgppwbGKgTdjQ1RFiEtFdVXExS` | `USDC-SPL\|SOL#452657748` | 20963935.314324 | 20956803.932985 | on-chain balance != csv claim by 7131.381339 (provider=https://api.helius.xyz/v1/wallet) |
| `FxmrpAFvRUYQwY67AstmZ4oDQWLZcTDzRst4mg3xi3BQ` | `USDC-SPL\|SOL#452657748` | 11262805.854823 | 11262785.977367 | on-chain balance != csv claim by 19.877456 (provider=https://api.helius.xyz/v1/wallet) |
| `3WjxkXbu9JBosUWUnJpjzec7Wa4rSkWo9khUbbLDCssk` | `SOL\|SOL#452657748` | 3977475.677924826 | 3.971873082 | on-chain balance != csv claim by 3977471.706051744 (provider=https://api.helius.xyz/v1/wallet) |
| `CKPbSRoW17z8fkFBy5NPcitiLZX4uqmfSv5MpFewkNoQ` | `USDC-SPL\|SOL#452657748` | 597198.7273 | 597137.11779 | on-chain balance != csv claim by 61.60951 (provider=https://api.helius.xyz/v1/wallet) |
| `A2nNTaAiDSYqCykEzXeK2hohAaxMpbpZQN` | `DOGE\|DOGE#6398667` | 510750.001 | 500921.001 | on-chain balance != csv claim by 9829 (provider=https://api.blockcypher.com/v1/doge/main/addrs/A2nNTaAiDSYqCykEzXeK2hohAaxMpbpZQN?before=6398668) |
| `CyKe1nTzB56CLtwgLiJ4WWTDvAjndwTfg4qs8j89k1NU` | `USDC-SPL\|SOL#452657748` | 313931.202213 | 313930.120493 | on-chain balance != csv claim by 1.08172 (provider=https://api.helius.xyz/v1/wallet) |
| `GD3tr9DR8ZnSy2pQ22RwzSR4vKY48SxHLdgCymXQ6W5i` | `SOL\|SOL#452657748` | 280610.551501448 | 0.165076123 | on-chain balance != csv claim by 280610.386425325 (provider=https://api.helius.xyz/v1/wallet) |
| `8wn2fZPKmUSMqNKLNsBfb3KkrVhChYYrL5Kzyfc2Dqox` | `USDC-SPL\|SOL#452657748` | 211507.1578 | 211506.1841 | on-chain balance != csv claim by 0.9737 (provider=https://api.helius.xyz/v1/wallet) |
| `92QMPBK852QdkGxyyuQZ8BmKczeQihEiGLHbFbx4tVdC` | `USDC-SPL\|SOL#452657748` | 78565.486116 | 78564.408759 | on-chain balance != csv claim by 1.077357 (provider=https://api.helius.xyz/v1/wallet) |
| `GHWmH2ZjB6jhv68dzRDvQeC3v23oyxh1W2F4Av6Gj2fy` | `USDC-SPL\|SOL#452657748` | 73998.989144 | 0.3427 | on-chain balance != csv claim by 73998.646444 (provider=https://api.helius.xyz/v1/wallet) |
| `4YWWA6aFe3E9vh7Bs3XqgiHqvD5pBZiuAZqk3K1X5drt` | `USDC-SPL\|SOL#452657748` | 68922.336373 | 68922.296373 | on-chain balance != csv claim by 0.04 (provider=https://api.helius.xyz/v1/wallet) |
| `B6gPEoYfgKDeMsT2Mxkie4SBsW5GMLQyCDQk42WTFDR` | `USDC-SPL\|SOL#452657748` | 47809.283616 | 47808.402335 | on-chain balance != csv claim by 0.881281 (provider=https://api.helius.xyz/v1/wallet) |
| `AisgV457n3GHkXhmqFh91uz3a9LcQi96if2r934WVXQE` | `USDC-SPL\|SOL#452657748` | 39054.546272 | 0.2422 | on-chain balance != csv claim by 39054.304072 (provider=https://api.helius.xyz/v1/wallet) |
| `Bou5nPZCcqQCNhSaKyJMLhqgoZVKBx6yoDkGEW5WY6W5` | `USDC-SPL\|SOL#452657748` | 19806.667559 | 19805.667358 | on-chain balance != csv claim by 1.000201 (provider=https://api.helius.xyz/v1/wallet) |
| `3n8t1n1NF6rhDeAPqQJF7RDD6JLGJPQftF1gMEKd4Sfs` | `SOL\|SOL#452657748` | 10283.081282864 | 0.164817137 | on-chain balance != csv claim by 10282.916465727 (provider=https://api.helius.xyz/v1/wallet) |
| `ltc1q007kfe7s6hsy8mf03jpzanrvrj8h6g2hk4h2wu34ggwcm4cyqersprfmnj` | `LTC\|LTC#3188165` | 7068.2631248 | 7035.85436179 | on-chain balance != csv claim by 32.40876301 (provider=https://api.blockcypher.com/v1/ltc/main/addrs/ltc1q007kfe7s6hsy8mf03jpzanrvrj8h6g2hk4h2wu34ggwcm4cyqersprfmnj?before=3188166) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `f1mfxobliwdamau2x25t4mcegzhfln5rwng4aemqi` | `FIL\|FIL#6421676` | 71742082.19019275710801678 | 71660790.894682732279470548 | live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes; on-chain balance != csv claim by 81291.295510024828546232 (provider=https://api.node.glif.io/rpc/v0) |
| `UQDY4-KtVxawZU_Vva7KTOhlhx8Ho0jI0ahyebYT5YuJkYSf` | `TONCOIN-NEW\|TON#96485222` | 16337432.383698934 | 14534894.655269916 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 1802537.728429018 (provider=https://tonapi.io/v2/accounts/UQDY4-KtVxawZU_Vva7KTOhlhx8Ho0jI0ahyebYT5YuJkYSf) |
| `0x51c6abe562e755582d268340b2cf0e2d8895a155dc9b7a7fb5465000d62d770b` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x51c6abe562e755582d268340b2cf0e2d8895a155dc9b7a7fb5465000d62d770b/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `0xd0c23b2613e9c14f4888ff6c4359f197424ca88cf2970bac9e862e5ee65d5fa8` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd0c23b2613e9c14f4888ff6c4359f197424ca88cf2970bac9e862e5ee65d5fa8/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572` | `APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `0xffd7f0d7b24ba690657dd2aa7ccb754c4be938e4fb970c72eef966d664bd1a45` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xffd7f0d7b24ba690657dd2aa7ccb754c4be938e4fb970c72eef966d664bd1a45/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `0xabc1f0788277b325b8357068e7c75a4e6614ed2a42cb5e950cfac4ec3599bb04` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xabc1f0788277b325b8357068e7c75a4e6614ed2a42cb5e950cfac4ec3599bb04/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `0x581bdd91ed2b08496cbab0d6b02674f8da3a2f128b086a3209d23ebafc5df553` | `USDC-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x581bdd91ed2b08496cbab0d6b02674f8da3a2f128b086a3209d23ebafc5df553/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `EQDbf_PU1DLAcneMzWjjUbK_Vc7oahP0iOpuHbxIwDfaI5Rm` | `TONCOIN-NEW\|TON#96485222` | 8910525 | 12.191071225 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 8910512.808928775 (provider=https://tonapi.io/v2/accounts/EQDbf_PU1DLAcneMzWjjUbK_Vc7oahP0iOpuHbxIwDfaI5Rm) |
| `0x612a398b1320b877cbc84fee033e88aa9e06cfa630b9d56d2b469ea7750bf157` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x612a398b1320b877cbc84fee033e88aa9e06cfa630b9d56d2b469ea7750bf157/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a` | `USDT-APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `UQBmoechltZNb69I-ksuqGH2ewSE0tFKWcr3AdjChs5E5Uii` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQBmoechltZNb69I-ksuqGH2ewSE0tFKWcr3AdjChs5E5Uii/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `UQAWBO4_6eixoz9q8Shzck9G2zOiu2S3p9T8jlgWYTpKBqB2` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQAWBO4_6eixoz9q8Shzck9G2zOiu2S3p9T8jlgWYTpKBqB2/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `0x2af5827f7017a7f6127184d73b2340f12e5dc248a48f48e2fe3fff08f6007beb` | `APTOS\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x2af5827f7017a7f6127184d73b2340f12e5dc248a48f48e2fe3fff08f6007beb/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1085724290): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1085724290) has been pruned","hint":"The requested data has bee |
| `f1rbhn7q4zynh44k4dbwxb3gp244wtttc6k4dkqfa` | `FIL\|FIL#6421676` | 2534759.598012690562093368 | 1941734.08620788299330102 | live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes; on-chain balance != csv claim by 593025.511804807568792348 (provider=https://api.node.glif.io/rpc/v0) |


_226 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `ADNbM5fBujCRBW1vqezNeAWmnsLp19ki3n` | `DOGE\|DOGE#6398667` | 472842600.92392216 | 473111276.63356829 | address has 227180 txs; exceeds public API limit 200000 — set BLOCKCHAIR_API_KEY or use archive provider |
| `ltc1qzvcgmntglcuv4smv3lzj6k8szcvsrmvk0phrr9wfq8w493r096ssm2fgsw` | `LTC\|LTC#3188165` | 57798.35009379 | 57231.93053992 | address has 318888 txs; exceeds public API limit 100000 — set BLOCKCHAIR_API_KEY or use archive provider |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 502260945 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-10-02T16:00:00Z |
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
