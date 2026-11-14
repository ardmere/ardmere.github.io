# BINANCE PR04NOV26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR04NOV26-assessment.json](./PR04NOV26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR04NOV26` |
| Snapshot time | `2026-11-04T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR04NOV26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

The available artifacts support Gen 2 / E2 classification and Stage 0 PoR disclosure. Users still need to trust Binance for wallet control, trusted setup honesty, and public availability of the full global proof stack.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_ownership_proof` | `NO_WALLET_OWNERSHIP_PROOF` | Stage 0 | No public batch-verifiable wallet_ownership_proof. |
| `global_proof.csv, verifying_key` | `UNVERIFIABLE` | Stage 0 | Public global proof/vk not available. |
| `trusted_setup_transcript` | `OPAQUE_TRUSTED_SETUP` | Stage 0 | Trusted setup transcript is not public. |

## Self-reported liabilities, reserves, and major-coin ratios

Label: **self-reported**. These rows repeat the exchange snapshot. USD totals sum per-coin customer liability and exchange reserve (self-custody + third-party custody) using CoinGecko prices at `snapshotTime`, the same price path as `reserveCoverageBps`. A coin with no price is left out of the USD total and listed as unpriced. `walletCoverageBps` and `reserveCoverageBps` are the separate **checked** coverage ratios. This section is not solvency advice.

| Figure | Value | Label |
| --- | --- | --- |
| Total liabilities (USD) | `132692864215.53` | self-reported |
| Total reserves (USD) | `133789440152.98` | self-reported |
| Unpriced coins | `ASTER`, `CAKE`, `CHR`, `CHZ`, `ENA`, `ENJ`, `FORM`, `PENDLE`, `PUMP`, `RLUSD`, `S`, `SSV`, `TRUMP`, `U`, `USD1`, `WIF`, `WLD`, `WLFI`, `ZEC` | — |

| Coin | Reserve ratio | Basis | Label |
| --- | --- | --- | --- |
| `BTC` | `1.00030363` | `amounts` | self-reported |
| `ETH` | `1.00015231` | `amounts` | self-reported |
| `USDT` | `1.01547623` | `amounts` | self-reported |
| `USDC` | `1.03120374` | `amounts` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it, or when neither amounts nor a reported rate were published (Unavailable). On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-11-12T00:00:00Z` |
| Previous snapshot | `2026-10-31T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `10 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Newer snapshot PR12NOV26 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `15b7c5eaf40db76b179b58da298c6263dc866b8d0f9a87985c70fd8d2c5ede5d` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `d2e048c5427b2c77b70a7e4a2a8bdc3899204c8f2b9426190ede2098a12c90b2` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20261103.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xbb073aab9606d14312320245b619c5186fb54a282a2cac2adff9d5408f38213e` |
| Verification bundle root | `0xf476bcf34eb92798da9b7c5292d87ce9b0670922562db127774a58ea6559b780` |
| Artifact bundle SHA-256 | `8f96cdfcdf5ce676a36ebd2b6a648e0ef7b27b9d400e2d6a3274e252454d4efe` |
| Verification bundle SHA-256 | `7e135fddb48d0e973cc52aa50560d20dfca05552790247ff2adf50f898d29f67` |

Local bundle paths: [PR04NOV26.artifact-bundle.json](../../../artifacts/binance/PR04NOV26/bundles/PR04NOV26.artifact-bundle.json), [PR04NOV26.verification-bundle.v2.json](../../../artifacts/binance/PR04NOV26/bundles/PR04NOV26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0558` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.6085` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.2017` |  |
| `onchain-balance-deposit` | `1.2` | `PASS` | `0.9904` |  |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR04NOV26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 8, UNVERIFIABLE 8, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x40207fdaa6015a765e5be0a7571823f11173eed5` | `BNB@BSC#131525096` | 1087.263581 | 1085.957904199650239514 | shortfall / 不足 (CSV exceeds on-chain) by 1.305676800349760486 (provider=cache) |
| `0x52854a35cf896a57b3a861ba7fbd2171f2445ef2` | `BNB@BSC#131525096` | 50474.13716 | 50474.06956648352931349 | shortfall / 不足 (CSV exceeds on-chain) by 0.06759351647068651 (provider=cache) |
| `0x6416c7690528117a633849be2cd06b9b5c55f171` | `BNB@BSC#131525096` | 1494.871777 | 1493.662140958755429666 | shortfall / 不足 (CSV exceeds on-chain) by 1.209636041244570334 (provider=cache) |
| `0x721093b510870ff0135be8ae2db2c80ff51f0908` | `BNB@BSC#131525096` | 4123.418755 | 3638.201983853059810745 | shortfall / 不足 (CSV exceeds on-chain) by 485.216771146940189255 (provider=cache) |
| `0xa7c0d36c4698981fab42a7d8c783674c6fe2592d` | `BNB@BSC#131525096` | 203862.3942 | 131546.179153457216230754 | shortfall / 不足 (CSV exceeds on-chain) by 72316.215046542783769246 (provider=cache) |
| `0xb2a6a72843db0f508204a56448413f3867ea691a` | `BNB@BSC#131525096` | 35011.51204 | 0.737837337885090657 | shortfall / 不足 (CSV exceeds on-chain) by 35010.774202662114909343 (provider=cache) |
| `0xbf83d18a46325acb7d8f40a462d23a92f467ed7a` | `BNB@BSC#131525096` | 6340754.812 | 6340654.273451315373361418 | shortfall / 不足 (CSV exceeds on-chain) by 100.538548684626638582 (provider=cache) |
| `0xb772fd70ea293335e5095597083e7da79e22ada7` | `ETH@ETH#26337722` | 1667.982379 | 1667.975792360534102986 | shortfall / 不足 (CSV exceeds on-chain) by 0.006586639465897014 (provider=cache) |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `PUMP\|SOL` | 11 | no native verifier for this (coin,network) pair yet |
| `TRX\|BSC` | 4 | no native verifier for this (coin,network) pair yet |
| `TRX\|ETH` | 4 | no native verifier for this (coin,network) pair yet |
| `TRX\|TRX` | 34 | no native verifier for this (coin,network) pair yet |
| `USDC\|ARC` | 2 | no native verifier for this (coin,network) pair yet |
| `U\|ROBINHOOD` | 1 | no native verifier for this (coin,network) pair yet |
| `XLM\|BSC` | 9 | no native verifier for this (coin,network) pair yet |
| `XLM\|XLM` | 11 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 71, WARN 18, UNVERIFIABLE 6, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x28e8d4eab37cfbf49edf7ff39c27377b7ae3a2a0` | `ASTER@BSC#131525096` | 30200045 | 35 | shortfall / 不足 (CSV exceeds on-chain) by 30200010 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `BCH@ETH#26337722` | 3.012 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.012 (provider=cache) |
| `0x52854a35cf896a57b3a861ba7fbd2171f2445ef2` | `BTC@BSC#131525096` | 782.28 | 782.277981970469216788 | shortfall / 不足 (CSV exceeds on-chain) by 0.002018029530783212 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#26337722` | 1105.5871 | 400 | shortfall / 不足 (CSV exceeds on-chain) by 705.5871 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#26337722` | 1047.9 | 106.362 | shortfall / 不足 (CSV exceeds on-chain) by 941.538 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CAKE@ETH#26337722` | 24593.94636 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 24593.94636 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CAKE@ETH#26337722` | 314432.2739 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 314432.2739 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `ENA@ETH#26337722` | 2.15510575 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2.15510575 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `ENA@ETH#26337722` | 2815479.927 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2815479.927 (provider=cache) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `ENA@ETH#26337722` | 9882010.636 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 9882010.636 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ENA@ETH#26337722` | 27343446.22 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 27343446.22 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ENA@ETH#26337722` | 5679532.089 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 5679532.089 (provider=cache) |
| `0x5a1e5b56aa33c056ed7bb93b88b1774a456f6694` | `ENA@ETH#26337722` | 35915486.42 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 35915486.42 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `ENA@ETH#26337722` | 510115216 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 510115216 (provider=cache) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `ENA@ETH#26337722` | 6001743.287 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 6001743.287 (provider=cache) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `ENA@ETH#26337722` | 96712658.75 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 96712658.75 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ENA@ETH#26337722` | 59634358.05 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 59634358.05 (provider=cache) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `ENA@ETH#26337722` | 9024549.068 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 9024549.068 (provider=cache) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `ENA@ETH#26337722` | 9523966.98 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 9523966.98 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `ENA@ETH#26337722` | 2.05552835 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2.05552835 (provider=cache) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `ENA@ETH#26337722` | 9176572.988 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 9176572.988 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ENA@ETH#26337722` | 787759017 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 787759017 (provider=cache) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `ENA@ETH#26337722` | 9196887.416 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 9196887.416 (provider=cache) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `LINK@ETH#26337722` | 256539.902 | 255157.402 | shortfall / 不足 (CSV exceeds on-chain) by 1382.5 (provider=cache) |
| `0xa64b436964e7415c0e70b9989a53e1fb9a90e726` | `POL@ETH#26337722` | 134000099.1 | 99.066047574468237635 | POL+MATIC ERC20 on ETH still << CSV; residual may be Polygon-native POL or internal ledger; shortfall / 不足 (CSV exceeds on-chain) by 134000000.033952425531762365 (provider=cache) |
| `0x26bff98f6f37c8d560684eaf2afd1149bf76633e` | `SSV@ETH#26337722` | 650019 | 14 | shortfall / 不足 (CSV exceeds on-chain) by 650005 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `TUSD@BSC#131525096` | 15000 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 15000 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `TUSD@BSC#131525096` | 256388.6563 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 256388.6563 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@BSC#131525096` | 396917.499 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 396917.499 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@BSC#131525096` | 9.9151 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 9.9151 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@BSC#131525096` | 1127096.623 | 0.01899189 | shortfall / 不足 (CSV exceeds on-chain) by 1127096.60400811 (provider=cache) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `TUSD@ETH#26337722` | 25584 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25584 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `TUSD@ETH#26337722` | 7262.2 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 7262.2 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `TUSD@ETH#26337722` | 3.00029591 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.00029591 (provider=cache) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `TUSD@ETH#26337722` | 25584 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25584 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `TUSD@ETH#26337722` | 1077.401833 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1077.401833 (provider=cache) |
| `0x5a1e5b56aa33c056ed7bb93b88b1774a456f6694` | `TUSD@ETH#26337722` | 46838.56089 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 46838.56089 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@ETH#26337722` | 529912.638 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 529912.638 (provider=cache) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `TUSD@ETH#26337722` | 172019.41 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 172019.41 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@ETH#26337722` | 31703.33056 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 31703.33056 (provider=cache) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `TUSD@ETH#26337722` | 25584 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25584 (provider=cache) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `TUSD@ETH#26337722` | 25584 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25584 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `TUSD@ETH#26337722` | 8693.109 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 8693.109 (provider=cache) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `TUSD@ETH#26337722` | 25584 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25584 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@ETH#26337722` | 5757027.079 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 5757027.079 (provider=cache) |
| `0x59fbc084aeed2c8e6e9b6ebe8a11d443efcb40d2` | `USD1@AB#16346865` | 16383889.29 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 16383889.29 (provider=cache) |
| `0x1d40c7717b8c43f877644ba86ec353bdc35c03dc` | `USDC@BSC#131525096` | 27707992.54 | 27707909.539366158 | shortfall / 不足 (CSV exceeds on-chain) by 83.000633842 (provider=cache) |
| `0x52854a35cf896a57b3a861ba7fbd2171f2445ef2` | `USDC@BSC#131525096` | 124365415.9 | 124365385.868053636577324716 | shortfall / 不足 (CSV exceeds on-chain) by 30.031946363422675284 (provider=cache) |
| `0xf17719d0a31afce4d88d57c83c4b2265f6a6b4cd` | `USDT@ARBITRUM#521018431` | 68442994.01 | 68442919.303954 | shortfall / 不足 (CSV exceeds on-chain) by 74.706046 (provider=cache) |
| `0x1d40c7717b8c43f877644ba86ec353bdc35c03dc` | `USDT@BSC#131525096` | 3027588.587 | 3024495.569715062213806256 | shortfall / 不足 (CSV exceeds on-chain) by 3093.017284937786193744 (provider=cache) |
| `0x40207fdaa6015a765e5be0a7571823f11173eed5` | `USDT@BSC#131525096` | 2309946.362 | 2309699.655537428 | shortfall / 不足 (CSV exceeds on-chain) by 246.706462572 (provider=cache) |
| `0x4540c69b4d42d190c99afce5a0aa4912beb6acb6` | `USDT@BSC#131525096` | 2712322.99 | 2710385.759820706612919423 | shortfall / 不足 (CSV exceeds on-chain) by 1937.230179293387080577 (provider=cache) |
| `0x52854a35cf896a57b3a861ba7fbd2171f2445ef2` | `USDT@BSC#131525096` | 146086864.3 | 146085284.450098007634429863 | shortfall / 不足 (CSV exceeds on-chain) by 1579.849901992365570137 (provider=cache) |
| `0x554a8dd77c0aa024e701e53d84103f83fa5bd54a` | `USDT@BSC#131525096` | 5312453.567 | 5310865.899557676710628042 | shortfall / 不足 (CSV exceeds on-chain) by 1587.667442323289371958 (provider=cache) |
| `0x6416c7690528117a633849be2cd06b9b5c55f171` | `USDT@BSC#131525096` | 3916489.091 | 3914646.10754111 | shortfall / 不足 (CSV exceeds on-chain) by 1842.98345889 (provider=cache) |
| `0x721093b510870ff0135be8ae2db2c80ff51f0908` | `USDT@BSC#131525096` | 4985160.184 | 4985014.907615778 | shortfall / 不足 (CSV exceeds on-chain) by 145.276384222 (provider=cache) |
| `0x7508f00e03a5996cf4b680a09f5ad05cc5c83225` | `USDT@BSC#131525096` | 2609120.795 | 2608861.774761943956604745 | shortfall / 不足 (CSV exceeds on-chain) by 259.020238056043395255 (provider=cache) |
| `0x9b594898e9441aabc7f97f02d5bd76b165732ed5` | `USDT@BSC#131525096` | 5467506.947 | 5467455.2867755578 | shortfall / 不足 (CSV exceeds on-chain) by 51.6602244422 (provider=cache) |
| `0xa010ec4e43c6448f2fe13340991bededf8763bf4` | `USDT@BSC#131525096` | 5608956.634 | 5608923.674122973 | shortfall / 不足 (CSV exceeds on-chain) by 32.959877027 (provider=cache) |
| `0xc55644e979f53358ffc7d39452e11f54ccd9b63d` | `USDT@BSC#131525096` | 2672493.029 | 2672418.7892597 | shortfall / 不足 (CSV exceeds on-chain) by 74.2397403 (provider=cache) |
| `0xc924b8347885e731912852933cb2f4cb42fa8170` | `USDT@BSC#131525096` | 2870632.338 | 2870460.239092952427233882 | shortfall / 不足 (CSV exceeds on-chain) by 172.098907047572766118 (provider=cache) |
| `0xd1d8e218550dd6c9cf3469305910bbc2b276cf4e` | `USDT@BSC#131525096` | 5634667.302 | 5633005.232224217185556053 | shortfall / 不足 (CSV exceeds on-chain) by 1662.069775782814443947 (provider=cache) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `USDT@ETH#26337722` | 11806884.16 | 11806874.76498 | shortfall / 不足 (CSV exceeds on-chain) by 9.39502 (provider=cache) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `USDT@ETH#26337722` | 9839544.562 | 9839390.509367 | shortfall / 不足 (CSV exceeds on-chain) by 154.052633 (provider=cache) |
| `0x6781b0098e3238caeede00647c9b69dd98a86a83` | `USDT@ETH#26337722` | 8190663.098 | 8190648.697943 | shortfall / 不足 (CSV exceeds on-chain) by 14.400057 (provider=cache) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `USDT@ETH#26337722` | 216309467.8 | 216309418.4044 | shortfall / 不足 (CSV exceeds on-chain) by 49.3956 (provider=cache) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `USDT@ETH#26337722` | 11065059.74 | 11064555.837358 | shortfall / 不足 (CSV exceeds on-chain) by 503.902642 (provider=cache) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `USDT@ETH#26337722` | 11541590.69 | 11433492.966765 | shortfall / 不足 (CSV exceeds on-chain) by 108097.723235 (provider=cache) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `USDT@OPTIMISM#159075811` | 203.587801 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 203.587801 (provider=cache) |
| `0xc34db06fd02be471dc4c1d8b0316a63ee37903bd` | `USDT@OPTIMISM#159075811` | 4652811.084 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 4652811.084 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `USDT@OPTIMISM#159075811` | 3000000 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3000000 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x00c1ce9c68f1c4a1b05e8d043b9b05784f41685d68e5009ffe379aa7f7c95ccf` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x0213c67ed78bc280887234fe5ed5e77272465317978ae86c25a71531d9332a2d` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x04164013f90b05d67f026779bf96e9c401c96f3485b645a786166e6935fba116` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `TJTKCzQpsG19CSMYF9TPtxz6wSHPSyeKfa` | `TUSD@TRX#87692258` | 1006263.211 | 606541.980457840768729417 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 399721.230542159231270583 (provider=https://api.trongrid.io) |
| `TJCo98saj6WND61g1uuKwJ9GMWMT9WkJFo` | `TUSD@TRX#87692258` | 2436174.701 | 1236174.701 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 1200000 (provider=https://api.trongrid.io) |
| `0x0213c67ed78bc280887234fe5ed5e77272465317978ae86c25a71531d9332a2d` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x04164013f90b05d67f026779bf96e9c401c96f3485b645a786166e6935fba116` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `TWd4WrZ9wn84f5x1hZhL4DHvk738ns5jwb` | `U@TRX#87692258` | 36094679 | 10000000 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 26094679 (provider=https://tron-rpc.publicnode.com) |
| `TADfh7Q4KA9LD6WfjjjcPMNra33gNRrCkb` | `USDT@TRX#87692258` | 8522644.706 | 4725815.836459 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 3796828.869541 (provider=https://tron-rpc.publicnode.com) |
| `TCQ8X6x9GiK4khydWVKfLJvm8H6WK3RE8B` | `USDT@TRX#87692258` | 8817651.578 | 3791429.726805 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 5026221.851195 (provider=https://api.trongrid.io) |
| `TDB1RpAHPQdGqfb1qi8J2jLY2oGe2jSy8A` | `USDT@TRX#87692258` | 8695658.309 | 1258910.398794 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 7436747.910206 (provider=https://tron-rpc.publicnode.com) |
| `TE4mtETxcuHHbprQmQkyVSaCC48Q5yk21i` | `USDT@TRX#87692258` | 7410996.609 | 4907657.042015 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 2503339.566985 (provider=https://tron-rpc.publicnode.com) |
| `TGnBDu95RaJEaw7dJpqzHt8Sy7faayBGvj` | `USDT@TRX#87692258` | 7947369.455 | 4696064.414258 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 3251305.040742 (provider=https://tron-rpc.publicnode.com) |
| `TGpTXS67vXyPERWMC2iQSG27TosfD13HLZ` | `USDT@TRX#87692258` | 7975215.984 | 3865675.527045 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 4109540.456955 (provider=https://tron-rpc.publicnode.com) |
| `TJTKCzQpsG19CSMYF9TPtxz6wSHPSyeKfa` | `USDT@TRX#87692258` | 266590758.2 | 23694005.988296 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 242896752.211704 (provider=https://tron-rpc.publicnode.com) |


_3 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `USDC\|ARC` | 2 | no native verifier for this (coin,network) pair yet |
| `U\|ROBINHOOD` | 1 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 9, WARN 93, UNVERIFIABLE 5, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `bc1q9n7j6a5wwru6tad9xrlncnagresvyjt96pl3en` | `BTC\|BTC#974457` | 1900.019572 | 100.01957174 | shortfall / 不足 (CSV exceeds on-chain) by 1800.00000026 (provider=https://blockstream.info/api) |
| `DBTrXhmzMbo1bkYAv6ZYodtwvTahhApf7y` | `DOGE\|DOGE#6442914` | 44296831.24 | 44296189.26763205 | shortfall / 不足 (CSV exceeds on-chain) by 641.97236795 (provider=https://api.blockcypher.com/v1/doge/main/addrs/DBTrXhmzMbo1bkYAv6ZYodtwvTahhApf7y?before=6442915) |
| `c5f9zfpkKMD9N8uLqJcFeJAAz7v12vDMnup9Y6EeQkk` | `SOL\|SOL#465092326` | 4211389.625 | 1.064970212 | shortfall / 不足 (CSV exceeds on-chain) by 4211388.560029788 (provider=https://api.helius.xyz/v1/wallet) |
| `5SDrsMNTYdhmApjfqYHDvjoW92f2S42vcc7zNDVcQ9Ej` | `USDC\|SOL#465092326` | 3.53509 | 2.000936 | shortfall / 不足 (CSV exceeds on-chain) by 1.534154 (provider=https://api.helius.xyz/v1/wallet) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDC\|SOL#465092326` | 1984800.492 | 3.000307 | shortfall / 不足 (CSV exceeds on-chain) by 1984797.491693 (provider=https://api.helius.xyz/v1/wallet) |
| `8DfHuYFv7VQ6wgcGtVDN2vtFn5ezF52hHNtPWYcqm8hU` | `USDC\|SOL#465092326` | 91759235.04 | 91757847.552578 | shortfall / 不足 (CSV exceeds on-chain) by 1387.487422 (provider=https://api.helius.xyz/v1/wallet) |
| `9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM` | `USDC\|SOL#465092326` | 1108.086885 | 1107.086885 | shortfall / 不足 (CSV exceeds on-chain) by 1 (provider=https://api.helius.xyz/v1/wallet) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDT\|SOL#465092326` | 21714.38138 | 2602.700038 | shortfall / 不足 (CSV exceeds on-chain) by 19111.681342 (provider=https://api.helius.xyz/v1/wallet) |
| `r4ncjTqY4ZmXdg7gXpkK9pfiythmyxeecq` | `XRP\|XRP#108098626` | 68859.55735 | 68844.149357 | shortfall / 不足 (CSV exceeds on-chain) by 15.407993 (provider=https://xrplcluster.com/) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0xe395e58840235d58dc8471c3351d5b248dbc408adc2c35c1f373f8042aa16619` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xe395e58840235d58dc8471c3351d5b248dbc408adc2c35c1f373f8042aa16619/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1159661950): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1159661950) has been pruned","hint":"The requested data has bee |
| `13oMyD7rJqSajN1b3hSDcScA5AYg83ANgE` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/13oMyD7rJqSajN1b3hSDcScA5AYg83ANgE): blockchair HTTP 430 |


_78 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `TRX\|BSC` | 4 | no native verifier for this (coin,network) pair yet |
| `TRX\|ETH` | 4 | no native verifier for this (coin,network) pair yet |
| `TRX\|TRX` | 34 | no native verifier for this (coin,network) pair yet |
| `USDC\|ARC` | 2 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `LZEjckteAtWrugbsy9zU8VHEZ4iUiXo9Nm` | `LTC\|LTC#3206938` | 3.00053979 | 686.34200343 | address has 1704537 txs; exceeds public API limit 100000 — set BLOCKCHAIR_API_KEY or use archive provider |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR04NOV26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-11-04T00:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary and wallet_address_list present. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `pass` | Wallet list present; hot-wallet on-chain verifier ran. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `unverifiable` | address-ownership verifier UNVERIFIABLE. |
| Is global_proof public and independently reproducible? (S1-3) | `unverifiable` | global-zk-proof verifier UNVERIFIABLE. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `unverifiable` | No public trusted setup transcript observed for zk-SNARK PoS. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P0` | Publish batch-verifiable wallet_ownership_proof for wallet_address_list. | `NO_WALLET_OWNERSHIP_PROOF` |
| `P0` | Make global proof.csv, verifying key, and parameter sources publicly downloadable. | `UNVERIFIABLE` |
| `P0` | Publish trusted setup transcript or migrate to a transparent-setup proof system. | `OPAQUE_TRUSTED_SETUP` |
| `P2` | Review onchain-balance verifier FAIL/WARN outliers in the verification bundle and reconcile CSV vs chain deltas. | `IMPLEMENTATION_SEMANTICS_RISK` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.
