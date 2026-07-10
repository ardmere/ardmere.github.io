# BINANCE PR01JUL26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR01JUL26-assessment.json](./PR01JUL26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01JUL26` |
| Snapshot time | `2026-07-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR01JUL26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

The available artifacts support Gen 2 / E2 classification and Stage 0 PoR disclosure. Users still need to trust Binance for wallet control, trusted setup honesty, and public availability of the full global proof stack.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_ownership_proof` | `NO_WALLET_OWNERSHIP_PROOF` | Stage 0 | No public batch-verifiable wallet_ownership_proof. |
| `global_proof.csv, verifying_key` | `UNVERIFIABLE` | Stage 0 | Public global proof/vk not available. |
| `trusted_setup_transcript` | `OPAQUE_TRUSTED_SETUP` | Stage 0 | Trusted setup transcript is not public. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-07-01T00:00:00Z` |
| Previous snapshot | `2026-06-01T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `4 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for binance.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `2bf5b5ead956cf67bdfcd9af36bb1bc00b321bb6fa248333a694c92b7249afb7` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `01b8f59d890414e42bdd66091fbd0c0f8a56735c6edd90d7c43ace67e49a3fa0` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260701.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x1fc76bdbd43232fe2d50e0c9ac473db6c54080990e874daa400b4d64e32434f4` |
| Verification bundle root | `0xf1d0b2d1b08f4d931b50b6cbede27408f498c7eca45ed3a3745178201f8b860c` |
| Artifact bundle SHA-256 | `12b87a1e9102b33068f815193302529707a7dd2e294162e81826d29ada53a801` |
| Verification bundle SHA-256 | `b097d6528ffdb91f79c01ba4e8c1ce3df31d3289f9e4be2a271d41402db8d408` |

Local bundle paths: [PR01JUL26.artifact-bundle.json](../../../artifacts/binance/PR01JUL26/bundles/PR01JUL26.artifact-bundle.json), [PR01JUL26.verification-bundle.v2.json](../../../artifacts/binance/PR01JUL26/bundles/PR01JUL26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0416` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.6849` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.2003` |  |
| `onchain-balance-deposit` | `1.2` | `PASS` | `0.1245` |  |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR01JUL26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 2, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x73f5ebe90f27b46ea12e5795d16c4b408b19cc6f` | `BNB@BSC#107345136` | 3161.073862 | 3161.072958497720012181 | accounted on-chain balance != csv claim by 0.000903502279987819 (provider=cache) |
| `0x56eddb7aa87536c09ccc2793473599fd21a8b17f` | `ETH@ETH#25433938` | 16026.52437 | 16021.024377714519000495 | accounted on-chain balance != csv claim by 5.499992285480999505 (provider=cache) |

#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 57, WARN 9, UNVERIFIABLE 6, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `AAVE@ETH#25433938` | 21513.70239 | 21513.658281429757891541 | on-chain balance != csv claim by 0.044108570242108459 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `BCH@ETH#25433938` | 1499.012 | 0 | on-chain balance != csv claim by 1499.012 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#25433938` | 1156.4841 | 400 | on-chain balance != csv claim by 756.4841 (provider=cache) |
| `0x9f8c163cba728e99993abe7495f06c0a3c8ac8b9` | `BUSD@AVAXC#89166730` | 129697.2016 | 0 | on-chain balance != csv claim by 129697.2016 (provider=cache) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `BUSD@MATIC#89439336` | 156651.5241 | 0 | on-chain balance != csv claim by 156651.5241 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `BUSD@MATIC#89439336` | 0.8 | 0 | on-chain balance != csv claim by 0.8 (provider=cache) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `BUSD@OPTIMISM#153632611` | 3679.63027 | 0 | on-chain balance != csv claim by 3679.63027 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `CAKE@ETH#25433938` | 17439.91448 | 0 | on-chain balance != csv claim by 17439.91448 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CAKE@ETH#25433938` | 348849 | 0 | on-chain balance != csv claim by 348849 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CAKE@ETH#25433938` | 24593.94636 | 0 | on-chain balance != csv claim by 24593.94636 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `CAKE@ETH#25433938` | 56949.96523 | 0 | on-chain balance != csv claim by 56949.96523 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `CAKE@ETH#25433938` | 53521.95853 | 0 | on-chain balance != csv claim by 53521.95853 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `CHZ@ETH#25433938` | 36253376.09 | 36253226.913775091822482267 | on-chain balance != csv claim by 149.176224908177517733 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `ENA@ETH#25433938` | 17871141.52 | 0 | on-chain balance != csv claim by 17871141.52 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `ENA@ETH#25433938` | 18268232.75 | 0 | on-chain balance != csv claim by 18268232.75 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `ENA@ETH#25433938` | 13706637.76 | 0 | on-chain balance != csv claim by 13706637.76 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `ENA@ETH#25433938` | 2075727.888 | 0 | on-chain balance != csv claim by 2075727.888 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `ENA@ETH#25433938` | 630115216 | 0 | on-chain balance != csv claim by 630115216 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ENA@ETH#25433938` | 27343446.22 | 0 | on-chain balance != csv claim by 27343446.22 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ENA@ETH#25433938` | 1246038.24 | 0 | on-chain balance != csv claim by 1246038.24 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ENA@ETH#25433938` | 580173076 | 0 | on-chain balance != csv claim by 580173076 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ENA@ETH#25433938` | 35734592.1 | 0 | on-chain balance != csv claim by 35734592.1 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `GRT@ETH#25433938` | 1315581.858 | 0 | on-chain balance != csv claim by 1315581.858 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `GRT@ARBITRUM#479089705` | 494352067.1 | 432481760.000000000001450888 | on-chain balance != csv claim by 61870307.099999999998549112 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `GRT@ETH#25433938` | 15236583.59 | 0 | on-chain balance != csv claim by 15236583.59 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `GRT@ETH#25433938` | 3.44990978 | 0 | on-chain balance != csv claim by 3.44990978 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `GRT@ETH#25433938` | 897665.8017 | 0 | on-chain balance != csv claim by 897665.8017 (provider=cache) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `GRT@ETH#25433938` | 3040301.896 | 0 | on-chain balance != csv claim by 3040301.896 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `GRT@ETH#25433938` | 23.26231794 | 0 | on-chain balance != csv claim by 23.26231794 (provider=cache) |
| `0x4aec0e98fc1fb55b9cc2faaa7a81acca42cb4e96` | `GRT@ETH#25433938` | 189374.7477 | 0 | on-chain balance != csv claim by 189374.7477 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@ETH#25433938` | 20329156.3 | 0 | on-chain balance != csv claim by 20329156.3 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@MATIC#89439336` | 73.1902 | 0 | on-chain balance != csv claim by 73.1902 (provider=cache) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `GRT@MATIC#89439336` | 20551.90111 | 0 | on-chain balance != csv claim by 20551.90111 (provider=cache) |
| `0xa64b436964e7415c0e70b9989a53e1fb9a90e726` | `POL@ETH#25433938` | 160000099.1 | 49.065528965081548221 | POL+MATIC ERC20 on ETH still << CSV; residual may be Polygon-native POL or internal ledger; on-chain balance != csv claim by 160000050.034471034918451779 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `TUSD@BSC#107345136` | 15000 | 0 | on-chain balance != csv claim by 15000 (provider=cache) |
| `0xe2fc31f816a9b94326492132018c3aecc4a93ae1` | `TUSD@BSC#107345136` | 53613.76233 | 0.018991890000006253 | on-chain balance != csv claim by 53613.743338109999993747 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@BSC#107345136` | 896917.499 | 0 | on-chain balance != csv claim by 896917.499 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `TUSD@BSC#107345136` | 256388.6563 | 0 | on-chain balance != csv claim by 256388.6563 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@BSC#107345136` | 60337.586 | 0 | on-chain balance != csv claim by 60337.586 (provider=cache) |
| `0x8894e0a0c962cb723c1976a4421c95949be2d4e3` | `TUSD@BSC#107345136` | 495359.3826 | 71251.229140440670720258 | on-chain balance != csv claim by 424108.153459559329279742 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `TUSD@ETH#25433938` | 1077.401833 | 0 | on-chain balance != csv claim by 1077.401833 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@BSC#107345136` | 9.9151 | 0 | on-chain balance != csv claim by 9.9151 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `TUSD@ETH#25433938` | 352175.8936 | 0 | on-chain balance != csv claim by 352175.8936 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `TUSD@ETH#25433938` | 30913.82647 | 0 | on-chain balance != csv claim by 30913.82647 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `TUSD@ETH#25433938` | 64133.45955 | 0 | on-chain balance != csv claim by 64133.45955 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@ETH#25433938` | 43355.39948 | 0 | on-chain balance != csv claim by 43355.39948 (provider=cache) |
| `0x59fbc084aeed2c8e6e9b6ebe8a11d443efcb40d2` | `USD1@AB#12718065` | 15422456.79 | 0 | on-chain balance != csv claim by 15422456.79 (provider=cache) |
| `0xee7ae85f2fe2239e27d9c1e23fffe168d63b4055` | `USDC@WLD#31764180` | 1424627.066 | 1424601.066315 | on-chain balance != csv claim by 25.999685 (provider=cache) |
| `0xb38e8c17e38363af6ebdcb3dae12e0243582891d` | `USDT@ARBITRUM#479089705` | 23455350.36 | 23455250.357965 | on-chain balance != csv claim by 100.002035 (provider=cache) |
| `0x9696f59e4d72e237be84ffd425dcad154bf96976` | `USDT@ETH#25433938` | 55213921.49 | 55213841.089374 | on-chain balance != csv claim by 80.400626 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `USDT@ETH#25433938` | 79580872.66 | 79580827.462177 | on-chain balance != csv claim by 45.197823 (provider=cache) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `USDT@OPTIMISM#153632611` | 8390320.854 | 0 | on-chain balance != csv claim by 8390320.854 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `USDT@OPTIMISM#153632611` | 1774.870765 | 0 | on-chain balance != csv claim by 1774.870765 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `USDT@OPTIMISM#153632611` | 47035311.56 | 0 | on-chain balance != csv claim by 47035311.56 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `USDT@OPTIMISM#153632611` | 187097.6339 | 0 | on-chain balance != csv claim by 187097.6339 (provider=cache) |
| `0xf055616bfc551a314f24ab6c09e8e1582b49cb44` | `USDC@SEIEVM#217115721` | 11565979 | 1640322.694007 | on-chain balance != csv claim by 9925656.305993 (provider=https://evm-rpc.sei-apis.com) |
| `0x3b9ce4b73fb57181194d83ec44544c0ccc77319a` | `WLD@WLD#31764180` | 22259517.02 | 22259495.684554739598300955 | on-chain balance != csv claim by 21.335445260401699045 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@ETH#25433938` | 3843003.696 | 0 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=3843003.696 (provider=cache) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `TUSD@TRX#84064751` | 44862.19616 | 30874.917585402395005971 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 13987.278574597604994029 (provider=https://tron-rpc.publicnode.com) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `U@TRX#84064751` | 34217094.01 | 34133181.242840738726129221 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 83912.767159261273870779 (provider=https://tron-rpc.publicnode.com) |
| `TJDENsfBJs4RFETt1X1W8wMDc8M5XnJhCe` | `USDT@TRX#84064751` | 46588884.34 | 44933999.754206 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 1654884.585794 (provider=https://api.trongrid.io) |
| `TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G` | `USDT@TRX#84064751` | 53002388.36 | 31559370.664527 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 21443017.695473 (provider=https://api.trongrid.io) |
| `TJqwA7SoZnERE4zW5uDEiPkbz4B66h9TFj` | `USDT@TRX#84064751` | 45383370.47 | 27599493.33298 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 17783877.13702 (provider=https://tron-rpc.publicnode.com) |
| `TCLgK89AnXbC9rewvhNb9UgXCc2qJJpBXh` | `USDT@TRX#84064751` | 41198990.13 | 37569479.885798 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 3629510.244202 (provider=https://tron-rpc.publicnode.com) |
| `TWd4WrZ9wn84f5x1hZhL4DHvk738ns5jwb` | `USDT@TRX#84064751` | 100000000 | 0 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 100000000 (provider=https://tron-rpc.publicnode.com) |
| `TQq26fUorctUZvrAgKg8Wz6QyYHvYd6xWK` | `USDT@TRX#84064751` | 15476669.78 | 13376896.889797 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 2099772.890203 (provider=https://tron-rpc.publicnode.com) |

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#25433938` | 1000 | 0 | no WBTC at snapshot height; BTC likely native/off-chain custody — EVM WBTC balanceOf cannot verify this row |

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 13, WARN 71, UNVERIFIABLE 9, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `CAKE\|APT#867767788` | 3340818.333 | 800000 | on-chain balance != csv claim by 2540818.333 (provider=cache) |
| `c5f9zfpkKMD9N8uLqJcFeJAAz7v12vDMnup9Y6EeQkk` | `SOL\|SOL#429980964` | 4764335.668 | 1.085865132 | on-chain balance != csv claim by 4764334.582134868 (provider=https://api.helius.xyz) |
| `0xaa367ea8e145e3af459d3b8f6d50b96c76d9cf9c7dabedf9fff192f36eda0ee0` | `USDC\|APT#867767788` | 24629266.53 | 0 | on-chain balance != csv claim by 24629266.53 (provider=cache) |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `USDC\|APT#867767788` | 1791897.02 | 0 | on-chain balance != csv claim by 1791897.02 (provider=cache) |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `USDC\|APT#867767788` | 182487.1009 | 0 | on-chain balance != csv claim by 182487.1009 (provider=cache) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDC\|SOL#429980964` | 499679778.8 | 497753753.759298 | on-chain balance != csv claim by 1926025.040702 (provider=https://api.helius.xyz) |
| `9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM` | `USDC\|SOL#429980964` | 1091.417022 | 1090.417022 | on-chain balance != csv claim by 1 (provider=https://api.helius.xyz) |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `USDT\|APT#867767788` | 149870817.3 | 0 | on-chain balance != csv claim by 149870817.3 (provider=cache) |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `USDT\|APT#867767788` | 400000000 | 0 | on-chain balance != csv claim by 400000000 (provider=cache) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDT\|SOL#429980964` | 133918581.8 | 133900798.217789 | on-chain balance != csv claim by 17783.582211 (provider=https://api.helius.xyz) |
| `G9RCBaYb8aBRxoe8QBC2ucGrVqjuZFysRhY8d56cnNT1` | `USDT\|SOL#429980964` | 311489.4214 | 0.016752 | on-chain balance != csv claim by 311489.404648 (provider=https://api.helius.xyz) |
| `rhWj9gaovwu2hZxYW7p388P8GRbuXFLQkK` | `XRP\|XRP#105288795` | 11801998.42 | 11801843.918172 | on-chain balance != csv claim by 154.501828 (provider=cache) |
| `rDAE53VfMvftPB4ogpWGWvzkQxfht6JPxr` | `XRP\|XRP#105288795` | 28888703.75 | 28888607.727477 | on-chain balance != csv claim by 96.022523 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `APT\|APT#867767788` | 62626935.96 | 22665573.834647 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 39961362.125353 (provider=cache) |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `APT\|APT#867767788` | 4212325.373 | 2559057.63737864 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 1653267.73562136 (provider=cache) |
| `0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3` | `APT\|APT#867767788` | 118380.3993 | 0.003994 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 118380.395306 (provider=cache) |
| `0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3` | `APT\|APT#867767788` | 114195.4981 | 6.371591 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 114189.126509 (provider=cache) |
| `1JyNxYBkMSDRRaVzPZRzgJxpHkRoK28d7D` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1JyNxYBkMSDRRaVzPZRzgJxpHkRoK28d7D): blockchair HTTP 430 |
| `1NbpHNjzNeNSUfHoJFSDsTNBZ8PxSXwtjv` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1NbpHNjzNeNSUfHoJFSDsTNBZ8PxSXwtjv): blockchair HTTP 430 |
| `19dQkvaH2NGgkGomzZu3qrnqRGCicXwedM` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/19dQkvaH2NGgkGomzZu3qrnqRGCicXwedM): blockchair HTTP 430 |
| `1P86nZCNWUiynP52AK2eTuTGZXYUTwX6qQ` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1P86nZCNWUiynP52AK2eTuTGZXYUTwX6qQ): blockchair HTTP 430 |
| `1G949ToNzeWSHqBCmvYRukaK6SqQtLdc8C` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1G949ToNzeWSHqBCmvYRukaK6SqQtLdc8C): blockchair HTTP 430 |
| `1B6cHqiZhcztHWriiyVrBbS7orHkSbjwbU` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1B6cHqiZhcztHWriiyVrBbS7orHkSbjwbU): blockchair HTTP 430 |
| `1KuPikhUYtHz3fmSQ2UvotpUuN672NuEcm` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1KuPikhUYtHz3fmSQ2UvotpUuN672NuEcm): blockchair HTTP 430 |
| `1Np9RvgpNve8hQ9H5oeGQ71T9uerDx32rr` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1Np9RvgpNve8hQ9H5oeGQ71T9uerDx32rr): blockchair HTTP 430 |
| `1KMcFqQxgJr7X9ADVGy28FSji4BuADkqUr` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1KMcFqQxgJr7X9ADVGy28FSji4BuADkqUr): blockchair HTTP 430 |
| `1Gi1qU3esZcmJtPphC19ZpNpVNL8ZFWFsM` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1Gi1qU3esZcmJtPphC19ZpNpVNL8ZFWFsM): blockchair HTTP 430 |
| `1Pzaqw98PeRfyHypfqyEgg5yycJRsENrE7` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/1Pzaqw98PeRfyHypfqyEgg5yycJRsENrE7): blockchair HTTP 430 |


_56 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3` | `APT\|APT#867767788` | 4000000.299 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 4000000.299 (provider=cache) |
| `0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d` | `APT\|APT#867767788` | 8333874.056 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 8333874.056 (provider=cache) |
| `0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a` | `APT\|APT#867767788` | 14596.32301 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 14596.32301 (provider=cache) |
| `0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c` | `APT\|APT#867767788` | 402396.998 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 402396.998 (provider=cache) |
| `0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34` | `APT\|APT#867767788` | 1437737.908 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 1437737.908 (provider=cache) |
| `0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386` | `APT\|APT#867767788` | 2060923.69 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 2060923.69 (provider=cache) |
| `bc1p7x4aaws8t8cmmccu39kun6cajglx6ntuta7lhyxjt2cwrw4k89zqn0wley` | `BTC\|BTC#956140` | 2000.069978 | 0.0699781 | on-chain UTXO sum << CSV; likely omnibus/internal BTC custody not visible at this address |
| `DJfU2p6woQ9GiBdiXsWZWJnJ9uDdZfSSNC` | `DOGE\|DOGE#6270926` | 103254932.4 | 99984413.52166502 | address has 1531292 txs; exceeds public API limit 200000 — set BLOCKCHAIR_API_KEY or use archive provider |
| `LZEjckteAtWrugbsy9zU8VHEZ4iUiXo9Nm` | `LTC\|LTC#3134270` | 156129.41 | 156491.61377138 | address has 1618327 txs; exceeds public API limit 100000 — set BLOCKCHAIR_API_KEY or use archive provider |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR01JUL26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-07-01T00:00:00Z |
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
