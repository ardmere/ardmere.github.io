# OKX 508399035 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [508399035-assessment.json](./508399035-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `508399035` |
| Snapshot time | `2026-06-18T16:00:00Z` |
| PoR Stage | `Stage 1 — Verifiable Disclosure` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `true` |

OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit 508399035. ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1.

Stage 1 is supported by public wallet_address_list, wallet_ownership_proof, global_proof, and parameter/proof artifacts. Stage 2 is not reached because canonical official anchoring, stable DA, low-friction user inclusion proof, stronger publication frequency, and full business-consistent constraints are not established in this report.

## 2. Stage Decision

### Stage 1, Effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `canonical on-chain/DA anchor` | `NO_CANONICAL_ANCHOR` | Stage 1 | No official exchange canonical on-chain/DA anchor established. |
| `weekly full PoR / daily anchor` | `HIGH_FREQUENCY_GAP` | Stage 1 | Monthly cadence does not satisfy Stage 2 frequency expectations. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-06-18T16:00:00Z` |
| Previous snapshot | `2026-05-06T16:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `3 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for okx.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `eadf8fe305b6893c36c3de295b5ceb59a8f4e9e3e12db2b49675bc6701a8a9c3` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `1fd7764c9beafcf34141eaffd0352d1653040bb0cd30a1d8d1c2cd9ae408d139` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026061900_V3.zip |
| `globalProofBundle` | `e1b2a413d2a0cd7464a3ad9e9ec2e2705e673240689455ce889e55723b33a075` | https://static.okx.com/cdn/okx/por/merkel/por_508399035_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xf9d4380bc4a8ff2146cd305090b353b22776054bd583094c0e3aabbca86f7102` |
| Verification bundle root | `0x1bf05b5c8e8882acdf58b2ee3c9414cf338d99566eb522e6283e1a80d969c120` |
| Artifact bundle SHA-256 | `d15001dced11ff7fc3c8aeda80a21a340eaaa8f2f83d5b4609cdce569e0c88b1` |
| Verification bundle SHA-256 | `0d4a6bf9fbe9eeb3289c77b4bc5eba15ac91db13c00a683d29260b44f06ccfb9` |

Local bundle paths: [508399035.artifact-bundle.json](../../../artifacts/okx/508399035/bundles/508399035.artifact-bundle.json), [508399035.verification-bundle.v2.json](../../../artifacts/okx/508399035/bundles/508399035.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `PASS` | `1.0000` | verified 196411/196411 address signatures |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0022` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.0030` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.0023` |  |
| `global-zk-proof` | `okx-1` | `PASS` | `1.0000` | zkSTARKValidator verify-global succeeded; summary merkle root bound to zk proof |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `508399035.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `btc-anchor` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `internal-consistency` (`PASS`)

Finding counts: WARN 11, PASS 22

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `BTC` | `custodyReserveBalances` | 5531.9478778728929917 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `ETH` | `custodyReserveBalances` | 39558.1772869728048806 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDT` | `custodyReserveBalances` | 191946635.8045224507315334 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDC` | `custodyReserveBalances` | 177423039.1015769740354214 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `XRP` | `custodyReserveBalances` | 963903.2137709051666404 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `DOGE` | `custodyReserveBalances` | 75415352.5328960016013964 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `SOL` | `custodyReserveBalances` | 128690.0037699438799555 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `OKB` | `custodyReserveBalances` | 12671.21093626265 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `LINK` | `custodyReserveBalances` | 296771.4034853100999999 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `LTC` | `custodyReserveBalances` | 54009.628879997753657 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `UNI` | `custodyReserveBalances` | 60916.2943337618763291 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 1, WARN 4, UNVERIFIABLE 16, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x05e851e9196fa1fc55f92ddd07ee83a978db35d0` | `ETH@ETH#25345500` | 832.12497732 | 17.4891558088222 | accounted on-chain balance != csv claim by 814.6358215111778 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x9b4495fcc3c99b3cb73a3722320ff1858c81c2f0` | `ETH@ETH#25345500` | 839.17505406 | 0.04997816 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=839.1250759 (provider=cache) |
| `0x4ccda666c69cd35be70aabd8931d824b2cfa3d2d` | `ETH@ETH#25345500` | 1236.35102649 | 0.499127259078130314 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=1235.851899230921869686 (provider=cache) |
| `0xd5f8c864f21824daf3a4b22d9c33693642ed7634` | `ETH@ETH#25345500` | 207342.52232797 | 0.799072178747936154 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=207341.723255791252063846 (provider=cache) |
| `0xe839a3e9efb32c6a56ab7128e51056585275506c` | `ETH@ETH#25345500` | 284131.31916117 | 5597.839126934988361811 | likely ETH2 deposit balance; deposit indexer unavailable: etherscan txlist: status=0 message=NOTOK; observed delta=278533.480034235011638189 (provider=cache) |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `ASSET-HUB\|DOT` | 50 | no native verifier for this (coin,network) pair yet |
| `A\|VAULTA` | 2 | no native verifier for this (coin,network) pair yet |
| `DOTK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ELF\|AELF` | 263 | no native verifier for this (coin,network) pair yet |
| `ETCK-KIP20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `ETC\|ETC` | 42 | no native verifier for this (coin,network) pair yet |
| `ETHK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `FILK-OKC20\|OKC` | 8 | no native verifier for this (coin,network) pair yet |
| `LINKK-OKC20\|OKC` | 3 | no native verifier for this (coin,network) pair yet |
| `LTCK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `OKB-OKC20\|OKC` | 3 | no native verifier for this (coin,network) pair yet |
| `TRXK-KIP20\|OKC` | 14 | no native verifier for this (coin,network) pair yet |
| `UNIK-OKC20\|OKC` | 1 | no native verifier for this (coin,network) pair yet |
| `USDC-OKC20\|OKC` | 237 | no native verifier for this (coin,network) pair yet |
| `USDT-OKC20\|OKC` | 100 | no native verifier for this (coin,network) pair yet |
| `XRPK-KIP20\|OKC` | 5 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 157, WARN 213, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `PEOPLE@ETH#25345500` | 1119325648.11 | 0 | on-chain balance != csv claim by 1119325648.11 (provider=cache) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `PEOPLE@ETH#25345500` | 175889537.9756854757 | 0 | on-chain balance != csv claim by 175889537.9756854757 (provider=cache) |
| `0xb0a27099582833c0cb8c7a0565759ff145113d64` | `PEOPLE@ETH#25345500` | 85271343.345626 | 0 | on-chain balance != csv claim by 85271343.345626 (provider=cache) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `PEOPLE@ETH#25345500` | 21744011.0089032823 | 0 | on-chain balance != csv claim by 21744011.0089032823 (provider=cache) |
| `0x343d752bb710c5575e417edb3f9fa06241a4749a` | `POLY-USDC@MATIC#88728936` | 9966364.122794 | 46.412112 | on-chain balance != csv claim by 9966317.710682 (provider=cache) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `PEOPLE@ETH#25345500` | 8788890.8489837 | 0 | on-chain balance != csv claim by 8788890.8489837 (provider=cache) |
| `0xafee421482faea92292ed3ffe29371742542ad72` | `USDT-ARBITRUM@ARBITRUM#474829257` | 8767014.800987 | 0 | on-chain balance != csv claim by 8767014.800987 (provider=cache) |
| `0x85dcd76d4fbd3aa0c85c27b9441222c19a14134b` | `PEOPLE@ETH#25345500` | 5941615.515201 | 0 | on-chain balance != csv claim by 5941615.515201 (provider=cache) |
| `0xb5216cb558cb018583bed009ee25ca73eb27bb1d` | `USDT-OPTIMISM@OPTIMISM#153099811` | 5015693.220212 | 0 | on-chain balance != csv claim by 5015693.220212 (provider=cache) |
| `0xd8a3ed83a59d3e10d82a4b845d64b24ba0572baf` | `USDT-ARBITRUM@ARBITRUM#474829257` | 3911578.164781 | 0 | on-chain balance != csv claim by 3911578.164781 (provider=cache) |
| `0x3f03440b19d15ff2068daf68e57eb656bfd0a67e` | `USDC-AVAXC@AVAXC#88306240` | 3759840.63621 | 0 | on-chain balance != csv claim by 3759840.63621 (provider=cache) |
| `0xde3c8d726e906b379a020cdb04bbe4b8f39a726f` | `USDC-AVAXC@AVAXC#88306240` | 2674234.876199 | 0 | on-chain balance != csv claim by 2674234.876199 (provider=cache) |
| `0xe7c5e276b2434743f6c79365d967cbe4291477bc` | `USDT-ARBITRUM@ARBITRUM#474829257` | 2636733.295912 | 0 | on-chain balance != csv claim by 2636733.295912 (provider=cache) |
| `0x4a1ca3e90e9e3a00eb27e9f409b4db310aa7405f` | `USDT-ARBITRUM@ARBITRUM#474829257` | 1966160.455769 | 0 | on-chain balance != csv claim by 1966160.455769 (provider=cache) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `PEOPLE@ETH#25345500` | 932963.7752159399 | 0 | on-chain balance != csv claim by 932963.7752159399 (provider=cache) |
| `0x9098575016440086c694c85c5010a0ed85249675` | `USDT-OPTIMISM@OPTIMISM#153099811` | 879733.358725 | 0 | on-chain balance != csv claim by 879733.358725 (provider=cache) |
| `0x92ea7496eba5f001d620005f88f3e8e686e3d4ea` | `PEOPLE@ETH#25345500` | 379523.52 | 0 | on-chain balance != csv claim by 379523.52 (provider=cache) |
| `0x4cafbe2b5bd89f0c0b7d6fef75b47e7c2ae82f22` | `USDC-AVAXC@AVAXC#88306240` | 276669.708697 | 0 | on-chain balance != csv claim by 276669.708697 (provider=cache) |
| `0xa0c9b4ce1785953abdceae569789aa158a01da4c` | `USDT-ARBITRUM@ARBITRUM#474829257` | 271720.894945 | 0 | on-chain balance != csv claim by 271720.894945 (provider=cache) |
| `0xc74e4c556a16390165c99b40aabc39a87c358305` | `POLY-USDC@MATIC#88728936` | 235758.271106 | 0 | on-chain balance != csv claim by 235758.271106 (provider=cache) |
| `0xf418263a6468c80d3bcb71033c56aacf81b241ba` | `PEOPLE@ETH#25345500` | 200190.64 | 0 | on-chain balance != csv claim by 200190.64 (provider=cache) |
| `0x3b50287cf8f073ff443b84de57e2375d446b962b` | `POLY-USDC@MATIC#88728936` | 158913.879751 | 0 | on-chain balance != csv claim by 158913.879751 (provider=cache) |
| `0xfbabce3f61f043cad4516e8ec526717133c12057` | `USDT-OPTIMISM@OPTIMISM#153099811` | 135768.992487 | 0 | on-chain balance != csv claim by 135768.992487 (provider=cache) |
| `0x76705493d53beefadffc253acb8aef72b7e7f6c6` | `USDT-OPTIMISM@OPTIMISM#153099811` | 119173.24072 | 0 | on-chain balance != csv claim by 119173.24072 (provider=cache) |
| `0xb57198a6e5af11544936a520a54268bc6ababb0c` | `USDT-OPTIMISM@OPTIMISM#153099811` | 112390.602791 | 0 | on-chain balance != csv claim by 112390.602791 (provider=cache) |
| `0x04e9d79a17f4894c76c65c52516f790531ef6bee` | `USDT-ARBITRUM@ARBITRUM#474829257` | 106057.793878 | 0 | on-chain balance != csv claim by 106057.793878 (provider=cache) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `PEOPLE@ETH#25345500` | 102150.99 | 0 | on-chain balance != csv claim by 102150.99 (provider=cache) |
| `0x9d1964437bb6d849c66ab2a2a033340dc2d999dc` | `USDC-AVAXC@AVAXC#88306240` | 95786.116837 | 0 | on-chain balance != csv claim by 95786.116837 (provider=cache) |
| `0x2fb2841c747621a7187b6f2a5390ae77321e5257` | `USDC-AVAXC@AVAXC#88306240` | 87525.643751 | 0 | on-chain balance != csv claim by 87525.643751 (provider=cache) |
| `0x91d40e4818f4d4c57b4578d9eca6afc92ac8debe` | `OKB@ETH#25345500` | 81236.6762816413 | 0 | on-chain balance != csv claim by 81236.6762816413 (provider=cache) |
| `0x229b1c70124141ec96945aec2323a22b68bb093c` | `USDT-ARBITRUM@ARBITRUM#474829257` | 78299.07596 | 0 | on-chain balance != csv claim by 78299.07596 (provider=cache) |
| `0x99860ec99aae9c0aa5ab763ee73da1217a9da665` | `POLY-USDC@MATIC#88728936` | 73169.442597 | 0 | on-chain balance != csv claim by 73169.442597 (provider=cache) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `PEOPLE@ETH#25345500` | 67183.6 | 0 | on-chain balance != csv claim by 67183.6 (provider=cache) |
| `0x74f42a5bc01eaeaf724da49254da739c37b5a808` | `USDC-AVAXC@AVAXC#88306240` | 65134.755561 | 0 | on-chain balance != csv claim by 65134.755561 (provider=cache) |
| `0x6bea9e00a6bb63d22904d975bdf8a98f26f9f49d` | `USDC-AVAXC@AVAXC#88306240` | 64344.197481 | 0 | on-chain balance != csv claim by 64344.197481 (provider=cache) |
| `0x4a152242605ab21df5ac6e819fdb128a4efc33fc` | `POLY-USDC@MATIC#88728936` | 64344.197432 | 0 | on-chain balance != csv claim by 64344.197432 (provider=cache) |
| `0x5ca1063894e5bbd688830d83fefa5ce3d339744d` | `USDT-ARBITRUM@ARBITRUM#474829257` | 61121.706922 | 0 | on-chain balance != csv claim by 61121.706922 (provider=cache) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `PEOPLE@ETH#25345500` | 52369.72407388 | 0 | on-chain balance != csv claim by 52369.72407388 (provider=cache) |
| `0x582988097ef9f2cf50cfac3f5f3b1ed54bc52a88` | `USDT-ARBITRUM@ARBITRUM#474829257` | 49669.675025 | 0 | on-chain balance != csv claim by 49669.675025 (provider=cache) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `PEOPLE@ETH#25345500` | 44142.30371176 | 0 | on-chain balance != csv claim by 44142.30371176 (provider=cache) |
| `0x2ef3ddc24e4020173da2c0f3a148eff4a9f6d9dd` | `USDC-AVAXC@AVAXC#88306240` | 43777.044961 | 0 | on-chain balance != csv claim by 43777.044961 (provider=cache) |
| `0x9c308b1665097af7ab2e38a41577bfecdc2d1c5c` | `PEOPLE@ETH#25345500` | 42354.84532091 | 0 | on-chain balance != csv claim by 42354.84532091 (provider=cache) |
| `0x4a8c00e80e3cfae4c903f9462b4ee2826404fe05` | `USDT-OPTIMISM@OPTIMISM#153099811` | 41151.531102 | 0 | on-chain balance != csv claim by 41151.531102 (provider=cache) |
| `0xf1e7b58ccd200fcb9400ddbc30867c43b7dc9ee1` | `USDT-ARBITRUM@ARBITRUM#474829257` | 41151.331175 | 0 | on-chain balance != csv claim by 41151.331175 (provider=cache) |
| `0x7d328f5d5f4640f81b2f93d2903ec0a6c030eed6` | `PEOPLE@ETH#25345500` | 40963.67003613 | 0 | on-chain balance != csv claim by 40963.67003613 (provider=cache) |
| `0x04ba1b380b1418007c1407bbe12ac21a596ae1e0` | `USDC-AVAXC@AVAXC#88306240` | 39608.565961 | 0 | on-chain balance != csv claim by 39608.565961 (provider=cache) |
| `0x03426aacbcf9030a380b3ba4a1f16b968649ef3d` | `USDT-ARBITRUM@ARBITRUM#474829257` | 35911.962433 | 0 | on-chain balance != csv claim by 35911.962433 (provider=cache) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `PEOPLE@ETH#25345500` | 32622.72 | 0 | on-chain balance != csv claim by 32622.72 (provider=cache) |
| `0x4a4aaa0155237881fbd5c34bfae16e985a7b068d` | `OKB@ETH#25345500` | 30822.1988656955 | 0 | on-chain balance != csv claim by 30822.1988656955 (provider=cache) |
| `0xbe9ac2d59743d69c4054bab09fa9821929d0f31b` | `USDT-ARBITRUM@ARBITRUM#474829257` | 29456.250044 | 0 | on-chain balance != csv claim by 29456.250044 (provider=cache) |
| `0x85f56cd748379755494f1365a3efc19569973e99` | `USDC-AVAXC@AVAXC#88306240` | 29016.392484 | 0 | on-chain balance != csv claim by 29016.392484 (provider=cache) |
| `0x0ca0e2aa6f2203ef8fa6e1f7f1903fb3a8ec1d35` | `USDT-OPTIMISM@OPTIMISM#153099811` | 23407.566664 | 0 | on-chain balance != csv claim by 23407.566664 (provider=cache) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `PEOPLE@ETH#25345500` | 23007.43 | 0 | on-chain balance != csv claim by 23007.43 (provider=cache) |
| `0xc30082de324fbfb7b293bcda397e2da2942d861d` | `USDT-OPTIMISM@OPTIMISM#153099811` | 22994.919244 | 0 | on-chain balance != csv claim by 22994.919244 (provider=cache) |
| `0xb1da6e546bd592218d6e849a8177d9d4fbf12986` | `POLY-USDC@MATIC#88728936` | 20656.310731 | 0 | on-chain balance != csv claim by 20656.310731 (provider=cache) |
| `0x1cdf5ba7297f276bf0320c88ca948d4311c913ce` | `USDT-ARBITRUM@ARBITRUM#474829257` | 19083.586352 | 0 | on-chain balance != csv claim by 19083.586352 (provider=cache) |
| `0x41ddf01de355547d821873ebd4fb2ffa2dc63d4c` | `POLY-USDC@MATIC#88728936` | 18295.126814 | 0.00009 | on-chain balance != csv claim by 18295.126724 (provider=cache) |
| `0xe7f6a87b1950b687a31d068bd01c6be81bc7e1af` | `POLY-USDC@MATIC#88728936` | 16805.684276 | 0 | on-chain balance != csv claim by 16805.684276 (provider=cache) |
| `0x28734d8ccfa30a3270513b40fd147071b79d7638` | `USDT-OPTIMISM@OPTIMISM#153099811` | 14848.093439 | 0 | on-chain balance != csv claim by 14848.093439 (provider=cache) |
| `0x42cf18596ee08e877d532df1b7cf763059a7ea57` | `USDT-ARBITRUM@ARBITRUM#474829257` | 14758.390218 | 0 | on-chain balance != csv claim by 14758.390218 (provider=cache) |
| `0x6cc5f688a315f3dc28a7781717a9a798a59fda7b` | `OKB@ETH#25345500` | 11302.306303002590174878 | 0 | on-chain balance != csv claim by 11302.306303002590174878 (provider=cache) |
| `0x5d002b87b6e1bdf390027c1210fdc1f6dd2ed85c` | `USDT-OPTIMISM@OPTIMISM#153099811` | 10339.783718 | 0 | on-chain balance != csv claim by 10339.783718 (provider=cache) |
| `0x00000000004e3d5628234f18b977041e5242651f` | `POLY-USDC@MATIC#88728936` | 8106.253163 | 0 | on-chain balance != csv claim by 8106.253163 (provider=cache) |
| `0x7cfb1543e6dec0dadd2435f23b97e08d682b6241` | `USDT-ARBITRUM@ARBITRUM#474829257` | 7490.499101 | 0 | on-chain balance != csv claim by 7490.499101 (provider=cache) |
| `0x48120be2162bf6d7b1c8ea04a150a530d294dcd0` | `PEOPLE@ETH#25345500` | 7271.51 | 0 | on-chain balance != csv claim by 7271.51 (provider=cache) |
| `0x26738cb185d888bb8c77ef0821dae410d0a748de` | `USDT-OPTIMISM@OPTIMISM#153099811` | 6580.761654 | 0 | on-chain balance != csv claim by 6580.761654 (provider=cache) |
| `0x4610dd69a52fc33fb50e284eec80cdfb0049caf7` | `PEOPLE@ETH#25345500` | 6266.34406 | 0 | on-chain balance != csv claim by 6266.34406 (provider=cache) |
| `0xcbabcb8d416591cc30c875f0a76df4dd1a53ba65` | `PEOPLE@ETH#25345500` | 6097.332 | 0 | on-chain balance != csv claim by 6097.332 (provider=cache) |
| `0x0d96b0525186c56aab001ff856313c85ced6bbe3` | `PEOPLE@ETH#25345500` | 5740.62 | 0 | on-chain balance != csv claim by 5740.62 (provider=cache) |
| `0xa58fe8afb81d4d4f94b9d3523845fe5b2077f779` | `POLY-USDC@MATIC#88728936` | 5823.6435 | 0.026754 | on-chain balance != csv claim by 5823.616746 (provider=cache) |
| `0xcf67640fd0716eb9c0d18c32214a36b8d8a2abaf` | `USDT-ARBITRUM@ARBITRUM#474829257` | 5700.363529 | 0 | on-chain balance != csv claim by 5700.363529 (provider=cache) |
| `0x787d2b2c401d2635fa4251fe6cf267177f5e45d5` | `PEOPLE@ETH#25345500` | 4912.04882206 | 0 | on-chain balance != csv claim by 4912.04882206 (provider=cache) |
| `0x942fd3a57d045b83f1e8f1d0381b496950cc16d3` | `USDC-AVAXC@AVAXC#88306240` | 4884.342349 | 0 | on-chain balance != csv claim by 4884.342349 (provider=cache) |
| `0xdf18e02f023c2f87245e30509ae273854d825573` | `PEOPLE@ETH#25345500` | 4239.74390235 | 0 | on-chain balance != csv claim by 4239.74390235 (provider=cache) |
| `0x0145506e5169a5efcbfc048636553d532108d668` | `PEOPLE@ETH#25345500` | 3857.1065 | 0 | on-chain balance != csv claim by 3857.1065 (provider=cache) |
| `0xe57245fb0e5fc8d3cb02fb05e1cb0a8939970d6d` | `USDT-OPTIMISM@OPTIMISM#153099811` | 3901.736277 | 0 | on-chain balance != csv claim by 3901.736277 (provider=cache) |
| `0x7d8b449450b78d6ef2af6e69d0b818cc822f2804` | `USDT-ARBITRUM@ARBITRUM#474829257` | 3801.08989 | 0 | on-chain balance != csv claim by 3801.08989 (provider=cache) |
| `0x8dbeda7085ddb63fd91bdfb8ba2c69a9a06c7589` | `PEOPLE@ETH#25345500` | 3716.7398 | 0 | on-chain balance != csv claim by 3716.7398 (provider=cache) |
| `0x4e3dd9ec3047246a5fd849334438696df4851b97` | `PEOPLE@ETH#25345500` | 3530.58365806 | 0 | on-chain balance != csv claim by 3530.58365806 (provider=cache) |
| `0x03c1f15c3f2461d6b53cd7ab799b95b6d2e669d0` | `PEOPLE@ETH#25345500` | 3500 | 0 | on-chain balance != csv claim by 3500 (provider=cache) |
| `0x5a741eab50d90a955392fcb06d77e81fccefdb3d` | `USDT-OPTIMISM@OPTIMISM#153099811` | 3265.078569 | 0 | on-chain balance != csv claim by 3265.078569 (provider=cache) |
| `0x6671c4d6bbdff00ef9aa4ad9488d400a36272cf6` | `PEOPLE@ETH#25345500` | 3236.30079953 | 0 | on-chain balance != csv claim by 3236.30079953 (provider=cache) |
| `0x1caf4162b5c0d254b4d07809dba8abf34b895b86` | `PEOPLE@ETH#25345500` | 3199.16248961 | 0 | on-chain balance != csv claim by 3199.16248961 (provider=cache) |
| `0xcc8e49f620f47dee27dd08caa8dd3554302388ea` | `PEOPLE@ETH#25345500` | 3036.885 | 0 | on-chain balance != csv claim by 3036.885 (provider=cache) |
| `0x6fb31ae81a72af0ce3acd41d5eefe57d2758f72c` | `PEOPLE@ETH#25345500` | 2863.17572044 | 0 | on-chain balance != csv claim by 2863.17572044 (provider=cache) |
| `0x8505bc2f4e44844db42cd9935905fbdd6a96d897` | `PEOPLE@ETH#25345500` | 2690.7088 | 0 | on-chain balance != csv claim by 2690.7088 (provider=cache) |
| `0x45d6e0a10d9a58bba0df2f8f8b66c7b447451bc7` | `PEOPLE@ETH#25345500` | 2578.16279 | 0 | on-chain balance != csv claim by 2578.16279 (provider=cache) |
| `0xd082b5c3b274d4433b2353b927d91a1082e5d1d1` | `PEOPLE@ETH#25345500` | 2459.19770969 | 0 | on-chain balance != csv claim by 2459.19770969 (provider=cache) |
| `0x6290a8e42010d4a68c8428920d5113300bdeaa40` | `PEOPLE@ETH#25345500` | 2162.286628 | 0 | on-chain balance != csv claim by 2162.286628 (provider=cache) |
| `0xb10cd07b0fbec13a9bc1e89e57e51e00fd26586c` | `PEOPLE@ETH#25345500` | 2069.48525718 | 0 | on-chain balance != csv claim by 2069.48525718 (provider=cache) |
| `0x1f1ca79462ec79445c6bd10b9225e1ad2ab9f3d1` | `USDC-AVAXC@AVAXC#88306240` | 2028.7302 | 0 | on-chain balance != csv claim by 2028.7302 (provider=cache) |
| `0x23d90deb6252846090f83e4e39555199bcc3b729` | `POLY-USDC@MATIC#88728936` | 1969.612683 | 0 | on-chain balance != csv claim by 1969.612683 (provider=cache) |
| `0x02a899128336dcc1cd087386292153e01de3848a` | `PEOPLE@ETH#25345500` | 1900.10162555 | 0 | on-chain balance != csv claim by 1900.10162555 (provider=cache) |
| `0x6682018cc6f3e6138ea8e63d6b87b2a48f445c1f` | `PEOPLE@ETH#25345500` | 1893.63922039 | 0 | on-chain balance != csv claim by 1893.63922039 (provider=cache) |
| `0x0af6b0b2f1e8b1a259fffc70f2c19e177e4638c4` | `OKB@ETH#25345500` | 1887.741668 | 0 | on-chain balance != csv claim by 1887.741668 (provider=cache) |
| `0xc8802feab2fafb48b7d1ade77e197002c210f391` | `USDC-AVAXC@AVAXC#88306240` | 1683.735722 | 0 | on-chain balance != csv claim by 1683.735722 (provider=cache) |
| `0x10e7d149e73dae219bb517de0fcb6a9601ba0f02` | `POLY-USDC@MATIC#88728936` | 1663.986314 | 50 | on-chain balance != csv claim by 1613.986314 (provider=cache) |
| `0xc310dd8b3b57748bf20d45243527768a09da7661` | `PEOPLE@ETH#25345500` | 1641.877204 | 0 | on-chain balance != csv claim by 1641.877204 (provider=cache) |
| `0x32de6ee2e3edb1405eb712277a100f3924075700` | `PEOPLE@ETH#25345500` | 1619.4748 | 0 | on-chain balance != csv claim by 1619.4748 (provider=cache) |
| `0x34bba0d0a1d0cab16ef6be267465efe5c4df201a` | `USDT-ARBITRUM@ARBITRUM#474829257` | 1439.092055 | 0 | on-chain balance != csv claim by 1439.092055 (provider=cache) |
| `0xe73856b980d2765790e6653228bfc534629e25b3` | `USDT-ARBITRUM@ARBITRUM#474829257` | 999.9 | 0 | on-chain balance != csv claim by 999.9 (provider=cache) |
| `0xbf0538e028264414f705ec5ad17e0c923ccdbc26` | `USDT-ARBITRUM@ARBITRUM#474829257` | 888.138339 | 0 | on-chain balance != csv claim by 888.138339 (provider=cache) |
| `0x7fd1cf2c0347ccb62dc2ca6db0031abc35a18e0e` | `USDT-ARBITRUM@ARBITRUM#474829257` | 799.9 | 0 | on-chain balance != csv claim by 799.9 (provider=cache) |
| `0x42cf18596ee08e877d532df1b7cf763059a7ea57` | `POLY-USDC@MATIC#88728936` | 700.688187 | 0 | on-chain balance != csv claim by 700.688187 (provider=cache) |
| `0xabf8a176e8077047cc009263e89081f1a7b6bff6` | `POLY-USDC@MATIC#88728936` | 661.6497 | 0 | on-chain balance != csv claim by 661.6497 (provider=cache) |
| `0x8c1dd7385efdd54b96dd417ec192f8dc17a13f92` | `USDC-AVAXC@AVAXC#88306240` | 624.449831 | 0 | on-chain balance != csv claim by 624.449831 (provider=cache) |
| `0x0b38179722faecaeaabf19477374f75b58645894` | `USDT-ARBITRUM@ARBITRUM#474829257` | 630.090064 | 0 | on-chain balance != csv claim by 630.090064 (provider=cache) |
| `0x4b3da5f48b5fc6552734f55f9f8fadbc902235f6` | `USDT-ARBITRUM@ARBITRUM#474829257` | 499 | 0 | on-chain balance != csv claim by 499 (provider=cache) |
| `0x000000000000d5775ff7721cefb8097af62e52dd` | `POLY-USDC@MATIC#88728936` | 459.213822 | 0 | on-chain balance != csv claim by 459.213822 (provider=cache) |
| `0x3273efbac62e9a4d91f454e04bc2969c3c42c23b` | `USDT-ARBITRUM@ARBITRUM#474829257` | 405.629 | 0 | on-chain balance != csv claim by 405.629 (provider=cache) |
| `0xf7f767ec932a506efce0b644a2344f628f804bfd` | `USDT-ARBITRUM@ARBITRUM#474829257` | 348.920896 | 0 | on-chain balance != csv claim by 348.920896 (provider=cache) |
| `0x8b688535f823d2a4c6b57997850f8d4fa3ff0135` | `USDT-ARBITRUM@ARBITRUM#474829257` | 329.4254 | 0 | on-chain balance != csv claim by 329.4254 (provider=cache) |
| `0x177dd993ad65ea85b5f85281e68893a46122f230` | `USDT-ARBITRUM@ARBITRUM#474829257` | 301.924862 | 0 | on-chain balance != csv claim by 301.924862 (provider=cache) |
| `0x0e081924b84fb5e378d7b2e3b0d3e7e3a5269f68` | `USDT-ARBITRUM@ARBITRUM#474829257` | 301.59 | 0 | on-chain balance != csv claim by 301.59 (provider=cache) |
| `0xbf90664017218ac7d9db5ae0c84873628e9ae1ae` | `OKB@ETH#25345500` | 253.2335271573 | 0 | on-chain balance != csv claim by 253.2335271573 (provider=cache) |
| `0x55823001c8cd87a6e86716ca768bce51f2d89c9c` | `OKB@ETH#25345500` | 233.510185 | 0 | on-chain balance != csv claim by 233.510185 (provider=cache) |
| `0xdd765ebdd3437b9a7674b6667d13ade7416e8237` | `POLY-USDC@MATIC#88728936` | 216.905631 | 0 | on-chain balance != csv claim by 216.905631 (provider=cache) |
| `0xc9b00b361785f3af5fd1027554ff1509923691b6` | `USDT-OPTIMISM@OPTIMISM#153099811` | 214.203491 | 0 | on-chain balance != csv claim by 214.203491 (provider=cache) |
| `0x8ed12b6e826589bfb32a3d0149d3fbd12d548722` | `USDT-ARBITRUM@ARBITRUM#474829257` | 212.401406 | 0 | on-chain balance != csv claim by 212.401406 (provider=cache) |
| `0xff8a035ea6c80673f741c2265985ed976a40d390` | `OKB@ETH#25345500` | 185.12488712 | 0 | on-chain balance != csv claim by 185.12488712 (provider=cache) |
| `0x73238f23192e190753f327d78394779e54249038` | `USDT-ARBITRUM@ARBITRUM#474829257` | 173.922828 | 0 | on-chain balance != csv claim by 173.922828 (provider=cache) |
| `0xaa5504106eaaab3df740671d85319a5fe81511f4` | `USDT-ARBITRUM@ARBITRUM#474829257` | 147.5 | 0 | on-chain balance != csv claim by 147.5 (provider=cache) |
| `0x62383739d68dd0f844103db8dfb05a7eded5bbe6` | `USDC-AVAXC@AVAXC#88306240` | 123.402721 | 0 | on-chain balance != csv claim by 123.402721 (provider=cache) |
| `0x105e68f301dd2ccc169243951da88c74fceb81e0` | `USDT-ARBITRUM@ARBITRUM#474829257` | 111.401379 | 0 | on-chain balance != csv claim by 111.401379 (provider=cache) |
| `0xe995fe7dc72f7a817ca38c70004fa5388e7f4c2e` | `USDT-OPTIMISM@OPTIMISM#153099811` | 114.71218 | 0 | on-chain balance != csv claim by 114.71218 (provider=cache) |
| `0xebeb422937ea3c652926dde7c978cea631e9a0e5` | `USDT-OPTIMISM@OPTIMISM#153099811` | 105.25121 | 0 | on-chain balance != csv claim by 105.25121 (provider=cache) |
| `0xf1d885f9579e211db28fbc1df31859fe16137eee` | `USDT-ARBITRUM@ARBITRUM#474829257` | 111.214236 | 0 | on-chain balance != csv claim by 111.214236 (provider=cache) |
| `0xf28143f7b282e59bd5f979012982e7cb9d9b95b0` | `USDT-ARBITRUM@ARBITRUM#474829257` | 101.250523 | 0 | on-chain balance != csv claim by 101.250523 (provider=cache) |
| `0x6f4b48aa129d873507770cdccb25bd60daaf09aa` | `USDC-AVAXC@AVAXC#88306240` | 96.434108 | 0 | on-chain balance != csv claim by 96.434108 (provider=cache) |
| `0xd953e4717c4aa75b08052645379cec64d539ea5a` | `USDT-ARBITRUM@ARBITRUM#474829257` | 96.207252 | 0 | on-chain balance != csv claim by 96.207252 (provider=cache) |
| `0xb8b70cf4edd49aa6c7212286917bcefb42da01c3` | `USDT-ARBITRUM@ARBITRUM#474829257` | 94.655592 | 0 | on-chain balance != csv claim by 94.655592 (provider=cache) |
| `0x5ebfc3de7cef8a26775b15913130831ee9159109` | `USDT-OPTIMISM@OPTIMISM#153099811` | 90.563481 | 0 | on-chain balance != csv claim by 90.563481 (provider=cache) |
| `0xc3520b4a68f7fd5e0832e90187a8fad3d7d3c880` | `USDT-ARBITRUM@ARBITRUM#474829257` | 90.48124 | 0 | on-chain balance != csv claim by 90.48124 (provider=cache) |
| `0xc7f758b2126a33239fb7e80f91c722caf7126c15` | `USDT-OPTIMISM@OPTIMISM#153099811` | 90.195751 | 0 | on-chain balance != csv claim by 90.195751 (provider=cache) |
| `0xc496e5e40f18035d6f8e805f1e5595657d7dca9f` | `USDT-ARBITRUM@ARBITRUM#474829257` | 88.95282 | 0 | on-chain balance != csv claim by 88.95282 (provider=cache) |
| `0x0938c63109801ee4243a487ab84dffa2bba4589e` | `USDT-ARBITRUM@ARBITRUM#474829257` | 78.397627 | 0 | on-chain balance != csv claim by 78.397627 (provider=cache) |
| `0xdc80e4412b4c6817318b3450b98f87f745206e0b` | `USDC-AVAXC@AVAXC#88306240` | 54.039151 | 0 | on-chain balance != csv claim by 54.039151 (provider=cache) |
| `0xaa488c5dd62d17465128c5034179e88a4f2f4170` | `USDC-AVAXC@AVAXC#88306240` | 47.995717 | 0 | on-chain balance != csv claim by 47.995717 (provider=cache) |
| `0xa9ac43f5b5e38155a288d1a01d2cbc4478e14573` | `OKB@ETH#25345500` | 35.38 | 0 | on-chain balance != csv claim by 35.38 (provider=cache) |
| `0x3c9df24f248ff949161a5a43bbf6dd69378d7799` | `OKB@ETH#25345500` | 34.5644520006 | 0 | on-chain balance != csv claim by 34.5644520006 (provider=cache) |
| `0xb0473ca01e91981623b65daf84ab7bc94a7e4000` | `USDC-AVAXC@AVAXC#88306240` | 24.9968 | 0 | on-chain balance != csv claim by 24.9968 (provider=cache) |
| `0x993fade64545c0d82d121fab5232b2072da98daf` | `OKB@ETH#25345500` | 18.95 | 0 | on-chain balance != csv claim by 18.95 (provider=cache) |
| `0x0799ddbf6f14db566ca4df4ff0575c4cc1e7749c` | `USDC-AVAXC@AVAXC#88306240` | 14.009382 | 8.843332 | on-chain balance != csv claim by 5.16605 (provider=cache) |
| `0x6cbb2add2a126e65cef029afca45576910a26b5f` | `OKB@ETH#25345500` | 14.607403 | 0 | on-chain balance != csv claim by 14.607403 (provider=cache) |
| `0xe2bc72b7f2c15476cdda714a44d08432d9302487` | `USDC-AVAXC@AVAXC#88306240` | 12 | 0 | on-chain balance != csv claim by 12 (provider=cache) |
| `0xaf6e08926da465fb3c59c2ad5dad47f6dca252ab` | `OKB@ETH#25345500` | 8.2962431535 | 0 | on-chain balance != csv claim by 8.2962431535 (provider=cache) |
| `0x079959dc64e2f6126414b66ac994600a2dd73099` | `USDC-AVAXC@AVAXC#88306240` | 5.604602 | 0 | on-chain balance != csv claim by 5.604602 (provider=cache) |
| `0x2ad0f758ce754653c79b681750bff16fe005f230` | `USDC-AVAXC@AVAXC#88306240` | 4.999999 | 0 | on-chain balance != csv claim by 4.999999 (provider=cache) |
| `0x9157b6c4f48112490962175b4058d485939653ca` | `USDC-AVAXC@AVAXC#88306240` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=cache) |
| `0x136bc54e0950eae85591954ae2bdea2199681234` | `USDC-AVAXC@AVAXC#88306240` | 4.99976 | 0 | on-chain balance != csv claim by 4.99976 (provider=cache) |
| `0x920c6944f9e1c5bb7a80d180a8df75f65b0b94e5` | `USDC-AVAXC@AVAXC#88306240` | 4.9997 | 0 | on-chain balance != csv claim by 4.9997 (provider=cache) |
| `0xceca016c59a63cccdf7a703f52b9dee581e1fa40` | `USDC-AVAXC@AVAXC#88306240` | 4.99968 | 0 | on-chain balance != csv claim by 4.99968 (provider=cache) |
| `0x3716ae413fcdc5c177735ce6cb38e68ed4c08b1d` | `USDC-AVAXC@AVAXC#88306240` | 4.99961 | 0 | on-chain balance != csv claim by 4.99961 (provider=cache) |
| `0x1810a491db6f47ea99e31903db4b65f74eaf0b36` | `USDC-AVAXC@AVAXC#88306240` | 4.9996 | 0 | on-chain balance != csv claim by 4.9996 (provider=cache) |
| `0xc6eca660c5da896784071ebc252fc27184da54cf` | `USDC-AVAXC@AVAXC#88306240` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=cache) |
| `0xe270b325c128ac89107ef4a078b303bf81dc7538` | `USDC-AVAXC@AVAXC#88306240` | 4.99956 | 0 | on-chain balance != csv claim by 4.99956 (provider=cache) |
| `0x381686ce20da524ad9520a157cdf7c871827ce0a` | `USDC-AVAXC@AVAXC#88306240` | 4.999362 | 0 | on-chain balance != csv claim by 4.999362 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `TLaGjwhvA8XQYSxFAcAXy7Dvuue9eGYitv` | `USDT-TRC20@TRX#83709660` | 243232693.234473 | 179134665.740288 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 64098027.494185 (provider=https://api.trongrid.io) |
| `TCw6YaWm3y6DvxY7M8hrCDnrJGeGMumzGJ` | `TRX@TRX#83709660` | 465294638.105331 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 465294638.105331 (provider=https://tron-rpc.publicnode.com) |
| `TYE5mNXSQHuPAdCa1etk76BKybD4FKJto6` | `USDT-TRC20@TRX#83709660` | 95974697.146829 | 65227120.737843 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 30747576.408986 (provider=https://tron-rpc.publicnode.com) |
| `0x7fea5b4568751533039179116e372e26b6b41b13` | `USDC@ETH#25345500` | 41621800.609858 | 4.517902 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=41621796.091956 (provider=cache) |
| `THdXEcBCZ7nagnmfRrfSHoiLRTbSRueoVW` | `TRX@TRX#83709660` | 71436182.072994 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 71436182.072994 (provider=https://tron-rpc.publicnode.com) |
| `TWAY47W8ci9XtfMgvs4r2A7Z2NFGUDJqpp` | `TRX@TRX#83709660` | 68289440 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 68289440 (provider=https://tron-rpc.publicnode.com) |
| `THwyeoCMSXQq2aGDpXpSqo89VnTtFktGGp` | `TRX@TRX#83709660` | 51572829.264957 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 51572829.264957 (provider=https://tron-rpc.publicnode.com) |
| `TD9GNrgmrg46yistfTMJpRQZcDmDRx1Bmq` | `TRX@TRX#83709660` | 33417832.529417 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 33417832.529417 (provider=https://tron-rpc.publicnode.com) |
| `TJVeQfDKsV2ZqrDzZKabgoE6De1vJhnP7h` | `TRX@TRX#83709660` | 18167856.407142 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 18167856.407142 (provider=https://tron-rpc.publicnode.com) |
| `TBNHZP2iAVyDhpTSZ8ee2Xjci1JExZTPE5` | `USDT-TRC20@TRX#83709660` | 8463250.31129 | 7978778.541636 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 484471.769654 (provider=https://tron-rpc.publicnode.com) |
| `0x611f7bf868a6212f871e89f7e44684045ddfb09d` | `OKB-X1@XLAYER` |  |  | rpc error (provider=): network XLAYER not configured |
| `TBwBJwj81yXc4DNKS19GJcpUUzfSWRbBzS` | `USDT-TRC20@TRX#83709660` | 13180213.989589 | 11881743.887357 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 1298470.102232 (provider=https://tron-rpc.publicnode.com) |
| `TNkUL9Za9uBnSJ5MQMDaf2aSSFpuuoKw3k` | `USDT-TRC20@TRX#83709660` | 8382310.886537 | 3438149.909432 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 4944160.977105 (provider=https://api.trongrid.io) |
| `TXGRttpC4PZkyid4WYFoWDJdekT1MmnGVC` | `TRX@TRX#83709660` | 7945378.727503 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 7945378.727503 (provider=https://tron-rpc.publicnode.com) |
| `TVvyKT63p66nqhQiNHX25izSUAvjcdER9X` | `TRX@TRX#83709660` | 6643124.263812 | 0 | live TRX balance vs snapshot CSV; Tron historical native balance limited on public nodes; on-chain balance != csv claim by 6643124.263812 (provider=https://tron-rpc.publicnode.com) |


_198 additional `WARN` rows omitted; see verification bundle._

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 111, WARN 123, UNVERIFIABLE 2, PASS 2

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `C68a6RCGLiPskbPYtAcsCjhG8tfTWYcoB4JjCrXFdqyo` | `USDG\|SOL#427319948` | 88840703.689716 | 88840584.112457 | on-chain balance != csv claim by 119.577259 (provider=https://api.helius.xyz) |
| `0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a` | `USDT-APTOS\|APT#840200933` | 7769877.194035 | 0 | on-chain balance != csv claim by 7769877.194035 (provider=cache) |
| `0x4c1ef44079cb31851349fba50d385f708a10ec7ac612859fdbf28888d1f7b572` | `APTOS\|APT#840200933` | 8478221.9460576 | 0 | on-chain balance != csv claim by 8478221.9460576 (provider=cache) |
| `0x51c6abe562e755582d268340b2cf0e2d8895a155dc9b7a7fb5465000d62d770b` | `USDC-APTOS\|APT#840200933` | 6250423.430363 | 0 | on-chain balance != csv claim by 6250423.430363 (provider=cache) |
| `0xd7eb876597276ad1c4159eb7a03dd4b694db9c142f28a9472a42510bc6a13a6e` | `USDC-APTOS\|APT#840200933` | 3769985.088027 | 0 | on-chain balance != csv claim by 3769985.088027 (provider=cache) |
| `0x2af5827f7017a7f6127184d73b2340f12e5dc248a48f48e2fe3fff08f6007beb` | `APTOS\|APT#840200933` | 4066856.13184567 | 0 | on-chain balance != csv claim by 4066856.13184567 (provider=cache) |
| `0x612a398b1320b877cbc84fee033e88aa9e06cfa630b9d56d2b469ea7750bf157` | `USDT-APTOS\|APT#840200933` | 3075172.237336 | 0 | on-chain balance != csv claim by 3075172.237336 (provider=cache) |
| `3WjxkXbu9JBosUWUnJpjzec7Wa4rSkWo9khUbbLDCssk` | `SOL\|SOL#427319948` | 2823088.987601145 | 4.047288082 | on-chain balance != csv claim by 2823084.940313063 (provider=https://api.helius.xyz) |
| `0x51c6abe562e755582d268340b2cf0e2d8895a155dc9b7a7fb5465000d62d770b` | `USDT-APTOS\|APT#840200933` | 2282189.077211 | 0 | on-chain balance != csv claim by 2282189.077211 (provider=cache) |
| `0xd7eb876597276ad1c4159eb7a03dd4b694db9c142f28a9472a42510bc6a13a6e` | `USDT-APTOS\|APT#840200933` | 2052655.909869 | 0 | on-chain balance != csv claim by 2052655.909869 (provider=cache) |
| `0x781fa312385e11c39c0c901f7d76d2892e32d272d09fad90f728656a542c2995` | `APTOS\|APT#840200933` | 1518257.6 | 46.92014438 | on-chain balance != csv claim by 1518210.67985562 (provider=cache) |
| `0x6f64b49c6c9899abaab8b748aa37c5776105a4e243cc4ee0a62d7aea4fb13865` | `USDT-APTOS\|APT#840200933` | 1103134.803228 | 0 | on-chain balance != csv claim by 1103134.803228 (provider=cache) |
| `0xf3a9bd8f8f9c7e07bef7d0a841e2165b6ed7a6767934975684f982d1c4b78581` | `USDC-APTOS\|APT#840200933` | 1101268.002415 | 0 | on-chain balance != csv claim by 1101268.002415 (provider=cache) |
| `0xc3ffdaa47eb60296d098f44b8e569e40543e733e01e78bf1abe96cbbcf442491` | `APTOS\|APT#840200933` | 1082874.7374205 | 0 | on-chain balance != csv claim by 1082874.7374205 (provider=cache) |
| `0xcf884dde6b75437e22b0860f6346ca3e385c35b28b891f16fa9da8d64cdd37` | `USDC-APTOS\|APT#840200933` | 1000478.712428 | 0 | on-chain balance != csv claim by 1000478.712428 (provider=cache) |
| `0x8f347361a9461e9312a4d2b5b5b928c65c3a740965705361317e3ca0015c64d8` | `APTOS\|APT#840200933` | 922537.7213531 | 0 | on-chain balance != csv claim by 922537.7213531 (provider=cache) |
| `0x76e20f5a1053f14843e9667917eb105d4213597a4b3033d1daa9da8e81bf55da` | `USDT-APTOS\|APT#840200933` | 775000.78229 | 0 | on-chain balance != csv claim by 775000.78229 (provider=cache) |
| `0xd2b472ad7168b53759e6cdccc2ef83964209806d28844d42f188e48f6b87fffd` | `USDC-APTOS\|APT#840200933` | 732361.705319 | 0 | on-chain balance != csv claim by 732361.705319 (provider=cache) |
| `A2nNTaAiDSYqCykEzXeK2hohAaxMpbpZQN` | `DOGE\|DOGE#6254074` | 510750.001 | 500921.001 | on-chain balance != csv claim by 9829 (provider=cache) |
| `0x26809f4fab301efb45d91ca4bd468cef3244438e2361dfc139f489a118ed95b7` | `APTOS\|APT#840200933` | 396111.4977252 | 0 | on-chain balance != csv claim by 396111.4977252 (provider=cache) |
| `0x612a398b1320b877cbc84fee033e88aa9e06cfa630b9d56d2b469ea7750bf157` | `USDC-APTOS\|APT#840200933` | 321140.509823 | 0 | on-chain balance != csv claim by 321140.509823 (provider=cache) |
| `0xdf13371a54322a65f3e049c07fcc0124ab1e4e22385a95d829b751f220ee872d` | `APTOS\|APT#840200933` | 321484.93783229 | 0 | on-chain balance != csv claim by 321484.93783229 (provider=cache) |
| `0xdd9d51fd5def96dfa15ebc8ea9a98d818e2ae86e6c4f5a40419c8bb043c8af32` | `APTOS\|APT#840200933` | 226678.74202629 | 0 | on-chain balance != csv claim by 226678.74202629 (provider=cache) |
| `0xb936d83662102e666b0863a7cab427ae1f710a3a9d8f315e454edc21a39f28cd` | `APTOS\|APT#840200933` | 149262.7879144 | 0 | on-chain balance != csv claim by 149262.7879144 (provider=cache) |
| `0xc9bc867420c8bc13add4008261774acd063435e2e77fd7b35975b57bdfc7716a` | `APTOS\|APT#840200933` | 84523.89 | 0 | on-chain balance != csv claim by 84523.89 (provider=cache) |
| `0x719f198d51db6d507bdd1bd4e43a98085b4b58428d01d24f3781d1863fecdd15` | `APTOS\|APT#840200933` | 82959.7177549 | 0 | on-chain balance != csv claim by 82959.7177549 (provider=cache) |
| `0x3db375ba86a32d0245170ef665db1b96072a435565fe492ffd92396fd7e97531` | `APTOS\|APT#840200933` | 74055.09987765 | 0 | on-chain balance != csv claim by 74055.09987765 (provider=cache) |
| `0xa507d2395d68f7fb6c9345bc39403981b50b83c41607893ae71af94e4877d637` | `USDC-APTOS\|APT#840200933` | 67308.964414 | 0 | on-chain balance != csv claim by 67308.964414 (provider=cache) |
| `0xe63c17bec2269d3cf0fedf88b16f4d05995fadaed3a97361c178af7e6c53895e` | `USDC-APTOS\|APT#840200933` | 64344.197471 | 0 | on-chain balance != csv claim by 64344.197471 (provider=cache) |
| `0x76e20f5a1053f14843e9667917eb105d4213597a4b3033d1daa9da8e81bf55da` | `USDC-APTOS\|APT#840200933` | 56205.143828 | 0 | on-chain balance != csv claim by 56205.143828 (provider=cache) |
| `0x9f4ac0affc4e79a0023b403e8455fb584c3118b8f6fb17732522ec4a3143294c` | `APTOS\|APT#840200933` | 42251.6389748 | 0 | on-chain balance != csv claim by 42251.6389748 (provider=cache) |
| `0xe63c17bec2269d3cf0fedf88b16f4d05995fadaed3a97361c178af7e6c53895e` | `USDT-APTOS\|APT#840200933` | 41151.531177 | 0 | on-chain balance != csv claim by 41151.531177 (provider=cache) |
| `0xaa9363f60cee8b21abf35b49ec0eb8cd86dfa9d479aa3734483ae93ef1c66669` | `APTOS\|APT#840200933` | 40684.7394852 | 0 | on-chain balance != csv claim by 40684.7394852 (provider=cache) |
| `0xfc4b8005085837503d583649b0daa2f68d2c69255c0f3656819232889e1711a5` | `APTOS\|APT#840200933` | 39908.5482851 | 0 | on-chain balance != csv claim by 39908.5482851 (provider=cache) |
| `0xcf884dde6b75437e22b0860f6346ca3e385c35b28b891f16fa9da8d64cdd37` | `USDT-APTOS\|APT#840200933` | 38437.899287 | 0 | on-chain balance != csv claim by 38437.899287 (provider=cache) |
| `GD3tr9DR8ZnSy2pQ22RwzSR4vKY48SxHLdgCymXQ6W5i` | `SOL\|SOL#427319948` | 35375.333790745 | 0.240511123 | on-chain balance != csv claim by 35375.093279622 (provider=https://api.helius.xyz) |
| `0xa507d2395d68f7fb6c9345bc39403981b50b83c41607893ae71af94e4877d637` | `USDT-APTOS\|APT#840200933` | 30927.363119 | 0 | on-chain balance != csv claim by 30927.363119 (provider=cache) |
| `0x3362d9ac7de31e010c333d1274aebb3f99bc58de864eb5c0c71e7ca2d9ebad70` | `APTOS\|APT#840200933` | 27816.52821165 | 0 | on-chain balance != csv claim by 27816.52821165 (provider=cache) |
| `0xe24c5f9ceddf4917ecf707cca259fa336b1ce4e3db6790584bec1747de0ba590` | `USDT-APTOS\|APT#840200933` | 26091.136751 | 0 | on-chain balance != csv claim by 26091.136751 (provider=cache) |
| `0x5bb0aea1ca2ae630beb1d238b2dac2505ec7528783a1460e9041b2341ad42188` | `USDT-APTOS\|APT#840200933` | 25919.640655 | 0 | on-chain balance != csv claim by 25919.640655 (provider=cache) |
| `0xfb9709f3bdcb816cf5ed2bbdc4e52387100faa499c80a26b56c8fc9b4693ea90` | `APTOS\|APT#840200933` | 24696.76750938 | 0 | on-chain balance != csv claim by 24696.76750938 (provider=cache) |
| `0xf3a9bd8f8f9c7e07bef7d0a841e2165b6ed7a6767934975684f982d1c4b78581` | `USDT-APTOS\|APT#840200933` | 23885.697585 | 0 | on-chain balance != csv claim by 23885.697585 (provider=cache) |
| `0x2cd5c058cf0594c70949c86b4bc548f871babd86f6e5250caab210cd9602d892` | `USDT-APTOS\|APT#840200933` | 23328.037091 | 0 | on-chain balance != csv claim by 23328.037091 (provider=cache) |
| `0x4e30a7c354208ab73e4702b30dc2bd617860c01a25601e728c78b6d064ac3f43` | `APTOS\|APT#840200933` | 21100.0476812 | 0 | on-chain balance != csv claim by 21100.0476812 (provider=cache) |
| `0x40e18b05ca9707b56d65720bec375c916d23df84fbf3a8c9db9d55c6ed0aad5d` | `APTOS\|APT#840200933` | 18455.45437125 | 0 | on-chain balance != csv claim by 18455.45437125 (provider=cache) |
| `0xd2b472ad7168b53759e6cdccc2ef83964209806d28844d42f188e48f6b87fffd` | `USDT-APTOS\|APT#840200933` | 15452.118755 | 0 | on-chain balance != csv claim by 15452.118755 (provider=cache) |
| `0xf719790bf38d1a981e6bf849602394e2531f1553bfcc26569a1ba4f9628f09ec` | `APTOS\|APT#840200933` | 15360.94 | 0 | on-chain balance != csv claim by 15360.94 (provider=cache) |
| `0xa5ab0d8982da365c24feef396a6884c41d01716158171de6a226f7b2d4704ae4` | `APTOS\|APT#840200933` | 13462.63114351 | 0 | on-chain balance != csv claim by 13462.63114351 (provider=cache) |
| `0xcd1020ed5014029f11a64737dcfb5f49c619a0ac19a35c1d91cc23e71f950422` | `APTOS\|APT#840200933` | 12411.39188521 | 0 | on-chain balance != csv claim by 12411.39188521 (provider=cache) |
| `3n8t1n1NF6rhDeAPqQJF7RDD6JLGJPQftF1gMEKd4Sfs` | `SOL\|SOL#427319948` | 10699.109934254 | 0.140801137 | on-chain balance != csv claim by 10698.969133117 (provider=https://api.helius.xyz) |
| `0xb3003185561644a9b70281af99381651a15315c75808a7ed4abd353a769e32c2` | `APTOS\|APT#840200933` | 8458.44929702 | 0 | on-chain balance != csv claim by 8458.44929702 (provider=cache) |
| `0xa07c09757eee7051143b0790c99f14ac0fc0d2fb40e4bec5aed863881e41e1c3` | `USDT-APTOS\|APT#840200933` | 6861 | 0 | on-chain balance != csv claim by 6861 (provider=cache) |
| `0x54b2980f24eb570c45b88cd02cc9e968912632f1f7840b23d3b4483591db2786` | `USDT-APTOS\|APT#840200933` | 6214.206338 | 0 | on-chain balance != csv claim by 6214.206338 (provider=cache) |
| `0x5bb0aea1ca2ae630beb1d238b2dac2505ec7528783a1460e9041b2341ad42188` | `USDC-APTOS\|APT#840200933` | 5262.5304 | 0 | on-chain balance != csv claim by 5262.5304 (provider=cache) |
| `0xa072b25a76fac70d5ecd5460ae855750bbb8dba454a7320e8f398c71a4cb69b9` | `APTOS\|APT#840200933` | 3187.5702114 | 0 | on-chain balance != csv claim by 3187.5702114 (provider=cache) |
| `0x54b2980f24eb570c45b88cd02cc9e968912632f1f7840b23d3b4483591db2786` | `USDC-APTOS\|APT#840200933` | 3024.688709 | 0 | on-chain balance != csv claim by 3024.688709 (provider=cache) |
| `0xe24c5f9ceddf4917ecf707cca259fa336b1ce4e3db6790584bec1747de0ba590` | `USDC-APTOS\|APT#840200933` | 1972.823232 | 0 | on-chain balance != csv claim by 1972.823232 (provider=cache) |
| `0x612a398b1320b877cbc84fee033e88aa9e06cfa630b9d56d2b469ea7750bf157` | `APTOS\|APT#840200933` | 1121.89293327 | 0 | on-chain balance != csv claim by 1121.89293327 (provider=cache) |
| `0xb1af32dc49767be11d794c7ac6a67c4a3838e5a3483c64549dfb4734dcba1ad3` | `APTOS\|APT#840200933` | 1010.9994907 | 0 | on-chain balance != csv claim by 1010.9994907 (provider=cache) |
| `0xf12ab2006b1ba6768e67eb4162fad0e5a1f33d33ed449340ed90f94b7485276e` | `APTOS\|APT#840200933` | 994.963871 | 0 | on-chain balance != csv claim by 994.963871 (provider=cache) |
| `0x34651505953861253b2796483433b500e3ceff98e38e259770a417b4bec30bb4` | `APTOS\|APT#840200933` | 957.94722598 | 0 | on-chain balance != csv claim by 957.94722598 (provider=cache) |
| `0x2b9378a4473fa663fdbe6a9673898d4ea2e9f6fe08afcb962ca1cbd81999bf4f` | `USDT-APTOS\|APT#840200933` | 720.227841 | 0 | on-chain balance != csv claim by 720.227841 (provider=cache) |
| `0x25ee5cdcd37bdb510e4d70c240ed666cae441b86085546d56edeea1c4fa589cf` | `APTOS\|APT#840200933` | 608.3976766 | 0 | on-chain balance != csv claim by 608.3976766 (provider=cache) |
| `0xfd84b64cfa48aee4de081341080eb421154f6414705a60da5a015958717bdb87` | `APTOS\|APT#840200933` | 502.15 | 0 | on-chain balance != csv claim by 502.15 (provider=cache) |
| `0x709b8e7b649b6de24552f8154c5c5edd4be950145d48f74829511e3cbef41b5a` | `APTOS\|APT#840200933` | 502.1494786 | 0 | on-chain balance != csv claim by 502.1494786 (provider=cache) |
| `0xd6880103d691d630ab40443ea1f8b4038bed1ff6a6bea10d8c41df6ce6d53dc3` | `USDT-APTOS\|APT#840200933` | 499.9 | 0 | on-chain balance != csv claim by 499.9 (provider=cache) |
| `0xef98dd35f4320b79038c861c52b0b3b6d3dde4cd74902452a260f554d0472ce4` | `USDT-APTOS\|APT#840200933` | 462.6289 | 0 | on-chain balance != csv claim by 462.6289 (provider=cache) |
| `0xdc8eb927f3db62950a6eb21c6a96531448aa385d2e34b66628ce48bae559e106` | `APTOS\|APT#840200933` | 412.3533615 | 0 | on-chain balance != csv claim by 412.3533615 (provider=cache) |
| `0x7dbc6d82ffd5fd7d564ef411a5b343191194475f3798f03cb37373753100fd6d` | `APTOS\|APT#840200933` | 329.43921465 | 0 | on-chain balance != csv claim by 329.43921465 (provider=cache) |
| `0xb7809da0d0716b12abbba6a41f1cdf3264b77a7719c2d1d396ca7516862e5e6f` | `USDT-APTOS\|APT#840200933` | 324.4112 | 0 | on-chain balance != csv claim by 324.4112 (provider=cache) |
| `0x688fd776608906a4396100762875ccd34696d12d75a1c72f6df97b28d3968d56` | `USDT-APTOS\|APT#840200933` | 307.74 | 0 | on-chain balance != csv claim by 307.74 (provider=cache) |
| `0x5c77f9f1244ccb1fba334aaaf4425add2e00420017527bcece3500119276335b` | `USDT-APTOS\|APT#840200933` | 299.9 | 0 | on-chain balance != csv claim by 299.9 (provider=cache) |
| `0xad19eea305cd97a7677e4864b0915ae3f205894602b33dd01dad532b4c8d4559` | `USDT-APTOS\|APT#840200933` | 299.9 | 0 | on-chain balance != csv claim by 299.9 (provider=cache) |
| `0xe63c17bec2269d3cf0fedf88b16f4d05995fadaed3a97361c178af7e6c53895e` | `APTOS\|APT#840200933` | 294.49640111 | 0 | on-chain balance != csv claim by 294.49640111 (provider=cache) |
| `0x59d8c3aee4dd4f8aef9d0f9cc5a149c3f0cf47632b301a6e04ad434328e28a94` | `USDT-APTOS\|APT#840200933` | 281.001 | 0 | on-chain balance != csv claim by 281.001 (provider=cache) |
| `0xf3a9bd8f8f9c7e07bef7d0a841e2165b6ed7a6767934975684f982d1c4b78581` | `APTOS\|APT#840200933` | 264.30330816 | 0 | on-chain balance != csv claim by 264.30330816 (provider=cache) |
| `0x5bb0aea1ca2ae630beb1d238b2dac2505ec7528783a1460e9041b2341ad42188` | `APTOS\|APT#840200933` | 259.24846477 | 0 | on-chain balance != csv claim by 259.24846477 (provider=cache) |
| `0xd7eb876597276ad1c4159eb7a03dd4b694db9c142f28a9472a42510bc6a13a6e` | `APTOS\|APT#840200933` | 246.63533888 | 0 | on-chain balance != csv claim by 246.63533888 (provider=cache) |
| `0xe24c5f9ceddf4917ecf707cca259fa336b1ce4e3db6790584bec1747de0ba590` | `APTOS\|APT#840200933` | 224.72117016 | 0 | on-chain balance != csv claim by 224.72117016 (provider=cache) |
| `0xdce51d4a200fda0bbf0343b5c0b3c9047afe0bf173b4207d490e9efceb3071d8` | `USDT-APTOS\|APT#840200933` | 211.48944 | 0 | on-chain balance != csv claim by 211.48944 (provider=cache) |
| `0xa507d2395d68f7fb6c9345bc39403981b50b83c41607893ae71af94e4877d637` | `APTOS\|APT#840200933` | 207.58645637 | 0 | on-chain balance != csv claim by 207.58645637 (provider=cache) |
| `0xcb27e898947904bff31dbbb815b05e5fb2ac13fa68f6975607a16712552aef76` | `APTOS\|APT#840200933` | 205.0986908 | 0 | on-chain balance != csv claim by 205.0986908 (provider=cache) |
| `0xb6d6e3d5a584540abb4b233145cf333213e220ff903808e5222386ba4a4d63ae` | `USDC-APTOS\|APT#840200933` | 87.794426 | 0 | on-chain balance != csv claim by 87.794426 (provider=cache) |
| `0xa0cbb53bc9742600daf01010b45e0723ee94ffd7491a1d0eed7938d560f1310f` | `USDC-APTOS\|APT#840200933` | 54.04153 | 0 | on-chain balance != csv claim by 54.04153 (provider=cache) |
| `0x2cd5c058cf0594c70949c86b4bc548f871babd86f6e5250caab210cd9602d892` | `USDC-APTOS\|APT#840200933` | 47.820919 | 0 | on-chain balance != csv claim by 47.820919 (provider=cache) |
| `0x30edf91802149dc71665e4258dae6c19abda611a8fec1bc59057eb4d4c28d5f8` | `USDC-APTOS\|APT#840200933` | 9.981616 | 0 | on-chain balance != csv claim by 9.981616 (provider=cache) |
| `0x8264eeafa0767bcb41b7711d476f6c87eb206e5ed20f2e880601c5d0b3d715a8` | `USDC-APTOS\|APT#840200933` | 9.980062 | 0 | on-chain balance != csv claim by 9.980062 (provider=cache) |
| `0x4b1326729f64c4eda98c6d597edfa48f5a10a69f840dc47006cb32ded66856b7` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x720f53e724eb334ec5ed4cfa3a208571bf000d578fef9cf719b697679556c8e5` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x5945aa11e6717ced3c0b959cce8a142b6ac6d10f1a83f3b5bd4306658800758d` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x47e16f9db551a90ec8bf2fd017d7d21797cba6c3ac8bd77711d59044a06ccb55` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x6d67f3901606ec31ee3021dfba49992c3e71c8e11e5da1b7baefe875a930474c` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x81d7cfc8a55333fa3bd2c749782b78395af1316a38ef1806feff240e9655e7e1` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x756ae256187f2fb867d52e505b179e038d8b5aecbab9e7c11d0c2c7b9d09f52b` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x1c21e926fcc9059907eb71b2826d2f2525f70dd89b61945e1a181a3e423cd392` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x60204b5571f8c20a925c6b835d567971c924f0de5da710504f76a6082fdcebf2` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x84cdacbf0cf376f16d762a45c2105862a766a89b7c623a5c1673a09d62baab45` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x216f4ad3dacef611d00456506fc8397184e4176efa405b0c2e394179cd1b63a8` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x4d2097a9c499c64373ffcdcbfa5b70c4e77c2ebae696684d3cad22e5281008bd` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x7bbc45a4c5cabcec15a9f37b318fd7d9f4912c7eddd433670e545522bdf15d09` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x599ab0e8174023f1d2bca49b2aaf86488fb3bad00af56f7fd3059e2a799dd7c6` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x135ba3467d14a3613ae17e5df30e7710f0d04b027fad532414d4a235344d6671` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0xcc5ad1895e56158e04638c3885c1c4f66d41bb00cbea56b70ba464d7da0b6d2c` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x87cf5735202ff86aff450721c9c3abe2b89e4a207e5380ae82fa95313d834d2e` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x4ef1b714b8d54de76737bc9452a07e73e1e2c0d2a86346c63b723ac8a96155db` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0xaffe1d91b202c3b28adccf15a6044d5a9c2cff2501661dd19e69f20676d244f0` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x81a215ba108615faaa347ce1d4f41f303bd45b2e692710915f20996d91c4cf86` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x2551114bcda955ef5fce627e5046ee0e04d9b9ecff3508a178cfc3124afdd2c3` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0x90bb51eaad22c38ebea8ee281f367972f63b2cc5c214ad118e31331bbe00bde1` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `0xba99b0ed33094dadf256789475ac4d3635e3a109b4c40c75cb1272367b32fa3c` | `USDC-APTOS\|APT#840200933` | 9.98 | 0 | on-chain balance != csv claim by 9.98 (provider=cache) |
| `ltc1q007kfe7s6hsy8mf03jpzanrvrj8h6g2hk4h2wu34ggwcm4cyqersprfmnj` | `LTC\|LTC#3127194` | 891.70290657 | 890.59758131 | on-chain balance != csv claim by 1.10532526 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `UQDY4-KtVxawZU_Vva7KTOhlhx8Ho0jI0ahyebYT5YuJkYSf` | `TONCOIN-NEW\|TON#74139837` | 13375926.688197012 | 12230464.075189614 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 1145462.613007398 (provider=cache) |
| `EQDbf_PU1DLAcneMzWjjUbK_Vc7oahP0iOpuHbxIwDfaI5Rm` | `TONCOIN-NEW\|TON#74139837` | 5117343.224889212 | 12.301580252 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 5117330.92330896 (provider=cache) |
| `0xdf7c04c9bebf4b35bd8c66a92469f0b66cf77ce8586b6262709a897659f4e772` | `USDC-SUI\|SUI#288486982` | 4408270.628906 | 0 | Sui public RPC may return latest checkpoint balance rather than snapshot checkpoint; on-chain balance != csv claim by 4408270.628906 (provider=cache) |
| `0x2bf102613d75bf4da4e954948d50636656301a916d13f8995180f84b1794eecf` | `USDC-SUI\|SUI#288486982` | 3003055.288068 | 0 | Sui public RPC may return latest checkpoint balance rather than snapshot checkpoint; on-chain balance != csv claim by 3003055.288068 (provider=cache) |
| `f1iesbg2vfubsunwyjjcvaeh3wuxoi4neyatprhmq` | `FIL\|FIL#6116396` | 2413113.463034402378767142 | 2003800.975933774438171252 | live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes; on-chain balance != csv claim by 409312.48710062794059589 (provider=cache) |
| `UQDYa6Lp99YwuDrDmCsngxs5dKzL_P7epA4iP0T7dsoY1po4` | `TONCOIN-NEW\|TON#74139837` | 2239827.886274203 | 2215579.697504624 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 24248.188769579 (provider=cache) |
| `UQBmoechltZNb69I-ksuqGH2ewSE0tFKWcr3AdjChs5E5Uii` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQBmoechltZNb69I-ksuqGH2ewSE0tFKWcr3AdjChs5E5Uii/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `UQAHsamD0jg4uG0FtCCkKBP2JPLLT02RaU2Ous8PX0sCrbCt` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQAHsamD0jg4uG0FtCCkKBP2JPLLT02RaU2Ous8PX0sCrbCt/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `UQAWBO4_6eixoz9q8Shzck9G2zOiu2S3p9T8jlgWYTpKBqB2` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQAWBO4_6eixoz9q8Shzck9G2zOiu2S3p9T8jlgWYTpKBqB2/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `UQBmoechltZNb69I-ksuqGH2ewSE0tFKWcr3AdjChs5E5Uii` | `TONCOIN-NEW\|TON#74139837` | 877259.558910857 | 785234.467266321 | live TON balance via tonapi vs snapshot CSV; no public historical seqno API; on-chain balance != csv claim by 92025.091644536 (provider=cache) |
| `f1ohs3ah3pj6s7rna2furnpvcoeztzef5eded6u2y` | `FIL\|FIL#6116396` | 318306.583979347047787358 | 292712.460997920031276636 | live Filecoin StateGetActor balance vs snapshot CSV; Glif lookback limited on public nodes; on-chain balance != csv claim by 25594.122981427016510722 (provider=cache) |
| `0x33fbf89217a07be0b7d5525a3bc17119a8d7bd70aa3b16eb77d836a4e6c8a026` | `USDC-SUI\|SUI#288486982` | 258228.977221 | 0 | Sui public RPC may return latest checkpoint balance rather than snapshot checkpoint; on-chain balance != csv claim by 258228.977221 (provider=cache) |
| `UQA2ogqg4sy7KfGuQVXWtbOm3vnE7NVGCeBCl6-XdJkqy3pq` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQA2ogqg4sy7KfGuQVXWtbOm3vnE7NVGCeBCl6-XdJkqy3pq/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `UQANL6-PmVenYBTCpKXDwNYaQPJ-fJFL__hWR0OHAtYRCoci` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQANL6-PmVenYBTCpKXDwNYaQPJ-fJFL__hWR0OHAtYRCoci/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |
| `UQDbSRhbADuce1uOWHQj5oQpgWCVSFA5ddQ0NRnUg26tIUON` | `USDT-TON\|TON` |  |  | ledger rpc error (provider=https://tonapi.io/v2/accounts/UQDbSRhbADuce1uOWHQj5oQpgWCVSFA5ddQ0NRnUg26tIUON/jettons/EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr): HTTP 400: {"error":"can't decode address EQCxE6mUTQJKFnGh1kS7u5DqsD2je6B84774Y_8f5g3L-Gr"} |


_108 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `ADNbM5fBujCRBW1vqezNeAWmnsLp19ki3n` | `DOGE\|DOGE#6254074` | 438568208.78528263 | 438836884.49492876 | address has 212361 txs; exceeds public API limit 200000 — set BLOCKCHAIR_API_KEY or use archive provider |
| `ltc1qzvcgmntglcuv4smv3lzj6k8szcvsrmvk0phrr9wfq8w493r096ssm2fgsw` | `LTC\|LTC#3127194` | 68727.08945253 | 68732.14920214 | address has 295453 txs; exceeds public API limit 100000 — set BLOCKCHAIR_API_KEY or use archive provider |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 508399035 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-06-18T16:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary and wallet_address_list present. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `pass` | Wallet list present; hot-wallet on-chain verifier ran. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `pass` | address-ownership verifier PASS. |
| Is global_proof public and independently reproducible? (S1-3) | `pass` | global-zk-proof verifier PASS. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `pass` | Plonky2 / zk-STARK uses transparent setup; no trusted setup ceremony required. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `fail` | No official exchange canonical on-chain/DA anchor established. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `fail` | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P1` | Anchor root/proof/commitment on-chain or publish immutable DA/archive records. | `NO_CANONICAL_ANCHOR` |
| `P1` | Add daily anchoring, weekly full PoR, or event-triggered updates to reduce timing risk. | `HIGH_FREQUENCY_GAP` |
| `P2` | Review onchain-balance verifier FAIL/WARN outliers in the verification bundle and reconcile CSV vs chain deltas. | `IMPLEMENTATION_SEMANTICS_RISK` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.
