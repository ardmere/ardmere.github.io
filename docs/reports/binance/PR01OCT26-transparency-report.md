# BINANCE PR01OCT26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR01OCT26-assessment.json](./PR01OCT26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01OCT26` |
| Snapshot time | `2026-10-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR01OCT26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

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
| Latest snapshot (evaluation set) | `2026-10-01T00:00:00Z` |
| Previous snapshot | `2026-09-01T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `7 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for binance.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `a44f71850124c5304fed5be4f1e85b640934acab7d9bb1e063e036c29620f128` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `656c334685e45635daabde9f2f446ba01b9d8b2f7a2ba03396327218f43ae685` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20261001.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x8c7964c29251e1b2571fe28fd433b48e36c34e506efb9072e504899215803663` |
| Verification bundle root | `0x0ba36daabab39503df1decec99fe1f67eed3a65b496924aa974aced689355a2c` |
| Artifact bundle SHA-256 | `606f39c8348aa72adc7d7387cfcf782f30e53ac40f1ffa75ae635a891f3da22f` |
| Verification bundle SHA-256 | `6137f87d46ca93da8aa2210e3c7080f64ac8ca446fa89b4e29ca0813db49d217` |

Local bundle paths: [PR01OCT26.artifact-bundle.json](../../../artifacts/binance/PR01OCT26/bundles/PR01OCT26.artifact-bundle.json), [PR01OCT26.verification-bundle.v2.json](../../../artifacts/binance/PR01OCT26/bundles/PR01OCT26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0380` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.6492` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.1693` |  |
| `onchain-balance-deposit` | `1.2` | `PASS` | `0.1762` |  |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR01OCT26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 4, WARN 2, UNVERIFIABLE 8, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x01c952174c24e1210d26961d456a77a39e1f0bb0` | `BNB@BSC#125000754` | 7556.919323 | 7556.915162868843981037 | accounted on-chain balance != csv claim by 0.004160131156018963 (provider=cache) |
| `0x29bdfbf7d27462a2d115748ace2bd71a2646946c` | `BNB@BSC#125000754` | 6196.599247 | 6196.523311463425233539 | accounted on-chain balance != csv claim by 0.075935536574766461 (provider=cache) |
| `0x3c783c21a0383057d128bae431894a5c19f9cf06` | `BNB@BSC#125000754` | 6794.552475 | 6794.0886404984867466 | accounted on-chain balance != csv claim by 0.4638345015132534 (provider=cache) |
| `0x4976a4a02f38326660d17bf34b431dc6e2eb2327` | `ETH@ETH#26093737` | 16360.57195 | 16360.555002904883598668 | accounted on-chain balance != csv claim by 0.016947095116401332 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xb2a6a72843db0f508204a56448413f3867ea691a` | `BNB@BSC#125000754` | 110113.6444 | 0.06894914595 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=110113.57545085405 (provider=cache) |
| `0x001ceb373c83ae75b9f5cf78fc2aba3e185d09e2` | `BNB@OPBNB#190783476` | 1311.974816 | 3.00042829591268743 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=1308.97438770408731257 (provider=cache) |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `PUMP\|SOL` | 9 | no native verifier for this (coin,network) pair yet |
| `TRX\|BSC` | 3 | no native verifier for this (coin,network) pair yet |
| `TRX\|ETH` | 3 | no native verifier for this (coin,network) pair yet |
| `TRX\|TRX` | 21 | no native verifier for this (coin,network) pair yet |
| `USDC\|ARC` | 1 | no native verifier for this (coin,network) pair yet |
| `U\|ROBINHOOD` | 1 | no native verifier for this (coin,network) pair yet |
| `XLM\|BSC` | 6 | no native verifier for this (coin,network) pair yet |
| `XLM\|XLM` | 11 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 58, WARN 18, UNVERIFIABLE 10, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x28e8d4eab37cfbf49edf7ff39c27377b7ae3a2a0` | `ASTER@BSC#125000754` | 21200045 | 35 | on-chain balance != csv claim by 21200010 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `BCH@ETH#26093737` | 1499.012 | 0 | on-chain balance != csv claim by 1499.012 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#26093737` | 1106.4841 | 400 | on-chain balance != csv claim by 706.4841 (provider=cache) |
| `0x9f8c163cba728e99993abe7495f06c0a3c8ac8b9` | `BUSD@AVAXC#96488899` | 128407.9711 | 0 | on-chain balance != csv claim by 128407.9711 (provider=cache) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `BUSD@OPTIMISM#157607011` | 4097.158029 | 0 | on-chain balance != csv claim by 4097.158029 (provider=cache) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `BUSD@MATIC#94738439` | 156934.4689 | 0 | on-chain balance != csv claim by 156934.4689 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `BUSD@MATIC#94738439` | 0.8 | 0 | on-chain balance != csv claim by 0.8 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `CAKE@ETH#26093737` | 57197.43806 | 0 | on-chain balance != csv claim by 57197.43806 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `CAKE@ETH#26093737` | 132648.4801 | 0 | on-chain balance != csv claim by 132648.4801 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CAKE@ETH#26093737` | 24593.94636 | 0 | on-chain balance != csv claim by 24593.94636 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `CAKE@ETH#26093737` | 67601.63877 | 0 | on-chain balance != csv claim by 67601.63877 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CAKE@ETH#26093737` | 200000 | 0 | on-chain balance != csv claim by 200000 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `ENA@ETH#26093737` | 2448187.762 | 0 | on-chain balance != csv claim by 2448187.762 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `ENA@ETH#26093737` | 25652679.97 | 0 | on-chain balance != csv claim by 25652679.97 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `ENA@ETH#26093737` | 111780222.7 | 0 | on-chain balance != csv claim by 111780222.7 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ENA@ETH#26093737` | 27343446.22 | 0 | on-chain balance != csv claim by 27343446.22 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ENA@ETH#26093737` | 5396031.043 | 0 | on-chain balance != csv claim by 5396031.043 (provider=cache) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `ENA@ETH#26093737` | 5736145.254 | 0 | on-chain balance != csv claim by 5736145.254 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `ENA@ETH#26093737` | 530115216 | 0 | on-chain balance != csv claim by 530115216 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ENA@ETH#26093737` | 56397870.49 | 0 | on-chain balance != csv claim by 56397870.49 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `ENA@ETH#26093737` | 27014615.13 | 0 | on-chain balance != csv claim by 27014615.13 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ENA@ETH#26093737` | 679017816 | 0 | on-chain balance != csv claim by 679017816 (provider=cache) |
| `0x605d0b92cde1e110d5ae9123969097ea2489d1aa` | `GRT@ARBITRUM#510496863` | 6000930.84 | 113 | on-chain balance != csv claim by 6000817.84 (provider=cache) |
| `0x5ce4825e2c5d7025b110d7bbeb5b1f03976e1443` | `GRT@ARBITRUM#510496863` | 140000100 | 100 | on-chain balance != csv claim by 140000000 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `GRT@ETH#26093737` | 1315581.858 | 0 | on-chain balance != csv claim by 1315581.858 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `GRT@ETH#26093737` | 16369441.24 | 0 | on-chain balance != csv claim by 16369441.24 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `GRT@ETH#26093737` | 61833129.46 | 0 | on-chain balance != csv claim by 61833129.46 (provider=cache) |
| `0x4aec0e98fc1fb55b9cc2faaa7a81acca42cb4e96` | `GRT@ETH#26093737` | 296513.3798 | 0 | on-chain balance != csv claim by 296513.3798 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `GRT@ETH#26093737` | 1181769.439 | 0 | on-chain balance != csv claim by 1181769.439 (provider=cache) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `GRT@ETH#26093737` | 0.092126 | 0 | on-chain balance != csv claim by 0.092126 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@ETH#26093737` | 13473559.78 | 0 | on-chain balance != csv claim by 13473559.78 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `GRT@ETH#26093737` | 12022036.31 | 0 | on-chain balance != csv claim by 12022036.31 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `GRT@ETH#26093737` | 80000000 | 0 | on-chain balance != csv claim by 80000000 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@MATIC#94738439` | 73.1902 | 0 | on-chain balance != csv claim by 73.1902 (provider=cache) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `GRT@MATIC#94738439` | 1142256.12 | 0 | on-chain balance != csv claim by 1142256.12 (provider=cache) |
| `0xa64b436964e7415c0e70b9989a53e1fb9a90e726` | `POL@ETH#26093737` | 152000149.1 | 99.066047574468237635 | POL+MATIC ERC20 on ETH still << CSV; residual may be Polygon-native POL or internal ledger; on-chain balance != csv claim by 152000050.033952425531762365 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `SHIB@ETH#26093737` | 104188955348 | 104188275472.480000006903899337 | on-chain balance != csv claim by 679875.519999993096100663 (provider=cache) |
| `0x26BFf98f6f37c8d560684eAF2AFd1149bF76633E` | `SSV@ETH#26093737` | 400019 | 14 | on-chain balance != csv claim by 400005 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `TUSD@BSC#125000754` | 256388.6563 | 0 | on-chain balance != csv claim by 256388.6563 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `TUSD@BSC#125000754` | 15000 | 0 | on-chain balance != csv claim by 15000 (provider=cache) |
| `0x8894e0a0c962cb723c1976a4421c95949be2d4e3` | `TUSD@BSC#125000754` | 966381.832 | 153639.229140440670720258 | on-chain balance != csv claim by 812742.602859559329279742 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@BSC#125000754` | 396917.499 | 0 | on-chain balance != csv claim by 396917.499 (provider=cache) |
| `0xe2fc31f816a9b94326492132018c3aecc4a93ae1` | `TUSD@BSC#125000754` | 41590.41498 | 0.018991890000006253 | on-chain balance != csv claim by 41590.395988109999993747 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@BSC#125000754` | 9.9151 | 0 | on-chain balance != csv claim by 9.9151 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@BSC#125000754` | 60337.586 | 0 | on-chain balance != csv claim by 60337.586 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `TUSD@ETH#26093737` | 5974.572 | 0 | on-chain balance != csv claim by 5974.572 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `TUSD@ETH#26093737` | 87315.99685 | 0 | on-chain balance != csv claim by 87315.99685 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `TUSD@ETH#26093737` | 1077.401833 | 0 | on-chain balance != csv claim by 1077.401833 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@ETH#26093737` | 828912.638 | 0 | on-chain balance != csv claim by 828912.638 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@ETH#26093737` | 31703.33056 | 0 | on-chain balance != csv claim by 31703.33056 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `TUSD@ETH#26093737` | 7284.413174 | 0 | on-chain balance != csv claim by 7284.413174 (provider=cache) |
| `0x59fbc084aeed2c8e6e9b6ebe8a11d443efcb40d2` | `USD1@AB#15367665` | 16383789.79 | 0 | on-chain balance != csv claim by 16383789.79 (provider=cache) |
| `0x25681ab599b4e2ceea31f8b498052c53fc2d74db` | `USDT@ARBITRUM#510496863` | 5883943.563 | 5883869.702849 | on-chain balance != csv claim by 73.860151 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `USDT@ETH#26093737` | 78823476.34 | 78822395.76939 | on-chain balance != csv claim by 1080.57061 (provider=cache) |
| `0x56eddb7aa87536c09ccc2793473599fd21a8b17f` | `USDT@ETH#26093737` | 72591040.4 | 72590977.363474 | on-chain balance != csv claim by 63.036526 (provider=cache) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `USDT@OPTIMISM#157607011` | 6624258.848 | 0 | on-chain balance != csv claim by 6624258.848 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `USDT@OPTIMISM#157607011` | 24000000 | 0 | on-chain balance != csv claim by 24000000 (provider=cache) |
| `0xee7ae85f2fe2239e27d9c1e23fffe168d63b4055` | `USDC@SEIEVM#235054691` | 1206526.245 | 3.000109 | on-chain balance != csv claim by 1206523.244891 (provider=https://evm-rpc.sei-apis.com) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@ETH#26093737` | 5264112.656 | 0 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=5264112.656 (provider=cache) |
| `TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G` | `TUSD@TRX#86713408` | 606254.76 | 0.75 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 606254.01 (provider=https://api.trongrid.io) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `BUSD@TRX#86713408` | 114364.4826 | 3.000644606715447029 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 114361.481955393284552971 (provider=https://api.trongrid.io) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `TUSD@TRX#86713408` | 3409166.959 | 3.000412606533647589 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 3409163.958587393466352411 (provider=https://tron-rpc.publicnode.com) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `U@TRX#86713408` | 174835.9011 | 3.481774710853280259 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 174832.419325289146719741 (provider=https://tron-rpc.publicnode.com) |
| `TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G` | `U@TRX#86713408` | 88470.72016 | 3.536147 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 88467.184013 (provider=https://tron-rpc.publicnode.com) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `USDC@TRX#86713408` | 85343.98693 | 600.000933 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 84743.985997 (provider=https://tron-rpc.publicnode.com) |
| `TAzsQ9Gx8eqFNFSKbeXrbi45CuVPHzA8wr` | `USDT@TRX#86713408` | 43855313.66 | 0.100814 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 43855313.559186 (provider=https://tron-rpc.publicnode.com) |
| `TCLgK89AnXbC9rewvhNb9UgXCc2qJJpBXh` | `USDT@TRX#86713408` | 40307130.56 | 0.100504 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 40307130.459496 (provider=https://tron-rpc.publicnode.com) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `USDT@TRX#86713408` | 613473361.4 | 2.000197 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 613473359.399803 (provider=https://tron-rpc.publicnode.com) |
| `TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y` | `USDT@TRX#86713408` | 49458892.27 | 0 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 49458892.27 (provider=https://tron-rpc.publicnode.com) |
| `TJDENsfBJs4RFETt1X1W8wMDc8M5XnJhCe` | `USDT@TRX#86713408` | 32434522.25 | 98.5 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 32434423.75 (provider=https://tron-rpc.publicnode.com) |
| `TK4ykR48cQQoyFcZ5N4xZCbsBaHcg6n3gJ` | `USDT@TRX#86713408` | 42667147.28 | 0.100853 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 42667147.179147 (provider=https://api.trongrid.io) |
| `TJqwA7SoZnERE4zW5uDEiPkbz4B66h9TFj` | `USDT@TRX#86713408` | 40727335.42 | 0.1 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 40727335.32 (provider=https://tron-rpc.publicnode.com) |
| `TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G` | `USDT@TRX#86713408` | 43338473.5 | 0.100546 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 43338473.399454 (provider=https://tron-rpc.publicnode.com) |


_3 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `USDC\|ARC` | 1 | no native verifier for this (coin,network) pair yet |
| `U\|ROBINHOOD` | 1 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `GRT@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `WLD@WLD` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#26093737` | 1000 | 0 | no WBTC at snapshot height; BTC likely native/off-chain custody — EVM WBTC balanceOf cannot verify this row |

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 11, WARN 114, UNVERIFIABLE 7, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `SOL\|SOL#452119503` | 1118679.401 | 1118670.581683223 | on-chain balance != csv claim by 8.819316777 (provider=https://api.helius.xyz) |
| `c5f9zfpkKMD9N8uLqJcFeJAAz7v12vDMnup9Y6EeQkk` | `SOL\|SOL#452119503` | 4252274.141 | 1.074220732 | on-chain balance != csv claim by 4252273.066779268 (provider=https://api.helius.xyz) |
| `5SDrsMNTYdhmApjfqYHDvjoW92f2S42vcc7zNDVcQ9Ej` | `USDC\|SOL#452119503` | 5147380.337 | 5147378.803354 | on-chain balance != csv claim by 1.533646 (provider=https://api.helius.xyz) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDC\|SOL#452119503` | 797194427.9 | 795231561.036757 | on-chain balance != csv claim by 1962866.863243 (provider=https://api.helius.xyz) |
| `9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM` | `USDC\|SOL#452119503` | 1096.076785 | 1095.076785 | on-chain balance != csv claim by 1 (provider=https://api.helius.xyz) |
| `GBAIA5U6E3FSRUW55AXACIVGX2QR5JYAS74OWLED3S22EGXVYEHPLGPA` | `USDC\|XLM#64705403` | 69700893.33 | 1235907.292 | on-chain balance != csv claim by 68464986.038 (provider=cache) |
| `GC5LF63GRVIT5ZXXCXLPI3RX2YXKJQFZVBSAO6AUELN3YIMSWPD6Z6FH` | `USDC\|XLM#64705403` | 44210927.97 | 3.0007133 | on-chain balance != csv claim by 44210924.9692867 (provider=cache) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDT\|SOL#452119503` | 807206387.1 | 807187424.95287 | on-chain balance != csv claim by 18962.14713 (provider=https://api.helius.xyz) |
| `tz1Q3jvYU9knekDYJfyvj3GjUy6898MNjvb2` | `USDT\|XTZ#15176643` | 2860007.38 | 2290793.961 | on-chain balance != csv claim by 569213.419 (provider=cache) |
| `tz2WDATNYnp7FdsmuZDYSidioZqeoLNZqXvE` | `USDT\|XTZ#15176643` | 39006.71833 | 3.000651 | on-chain balance != csv claim by 39003.717679 (provider=cache) |
| `rBtttd61FExHC68vsZ8dqmS3DfjFEceA1A` | `XRP\|XRP#107347788` | 11341503.47 | 11341313.669541 | on-chain balance != csv claim by 189.800459 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1081927336): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1081927336) has been pruned","hint":"The requested data has bee |
| `13oMyD7rJqSajN1b3hSDcScA5AYg83ANgE` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/13oMyD7rJqSajN1b3hSDcScA5AYg83ANgE): blockchair HTTP 430 |
| `14MS72FgWMqQmqGfqEVuh73efqZTzisQx2` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/14MS72FgWMqQmqGfqEVuh73efqZTzisQx2): blockchair HTTP 430 |


_99 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `TRX\|BSC` | 3 | no native verifier for this (coin,network) pair yet |
| `TRX\|ETH` | 3 | no native verifier for this (coin,network) pair yet |
| `TRX\|TRX` | 21 | no native verifier for this (coin,network) pair yet |
| `USDC\|ARC` | 1 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `bc1q9n7j6a5wwru6tad9xrlncnagresvyjt96pl3en` | `BTC\|BTC#969369` | 1200.019572 | 0.01957174 | on-chain UTXO sum << CSV; likely omnibus/internal BTC custody not visible at this address |
| `DJfU2p6woQ9GiBdiXsWZWJnJ9uDdZfSSNC` | `DOGE\|DOGE#6396399` | 418444275.2 | 415172394.45657312 | address has 1611900 txs; exceeds public API limit 200000 — set BLOCKCHAIR_API_KEY or use archive provider |
| `LZEjckteAtWrugbsy9zU8VHEZ4iUiXo9Nm` | `LTC\|LTC#3187197` | 306354.8969 | 307005.83097135 | address has 1683811 txs; exceeds public API limit 100000 — set BLOCKCHAIR_API_KEY or use archive provider |

#### `onchain-balance-deposit` (`PASS`)

Finding counts: WARN 34, PASS 2

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `2VRgRxyfFANV41oeogC4CMm1jRNvHo2ybh9BZxjxEicC` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `82BNLH9t1AcvYBLP52Nn6VXyANXarjvfSd2o8XPtYt4A` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `57v9pKHg2ZeitL9R3E5ogWMBgGT7QUYSYUcdDT4BFUrd` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `Ax51HLNLbiGbfPWc6e1HHEVW9DMhjxBLNoa9EeCvWrQM` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `AGktEKEQ9a7nr6SJ9NmzfYTRqqpoWDFRKKBnD8v8XUJU` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `8MuyW9P8e3JyEc2S96kwNyXHt8aQzznkPTkxGkek3wRp` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `CFuBZCBigBnwxsc8wwUiytn2eHfKxSSagkER4CBcarxi` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `HrniVxv8QmcdaXSaH76RC65CsH5vhTXkzqToSVNcry1u` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `ELi4EVfZNj9wFqZRdYJ73nWCxmBktLMFvSajoWQ2TapG` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `CeJCS4nsptSXAuhf1MEWbKHJVFuzicWcPmNMjsm3HMjz` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `8krQPwcxVsKYoLAq7k52CkcuPoPXWuXVqB9vJBZbHk8J` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `5CcnZj1sx5jx2xyMkAejBqY7KVypinCyYLHrv2kvNruy` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `9xqtuifqHkNc1jbcKGdbMP5hMW1yTL4NuiguzTYSfCGq` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `4qvbYSn2XEzhsQt9q5TcUSr4CNfmMSmU3CA8XLE9eQxs` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `9RE2f9mZrQR78ddCdFAcLiNzrcJ2kD5URmVJwXM3KPu2` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |


_19 additional `WARN` rows omitted; see verification bundle._

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR01OCT26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-10-01T00:00:00Z |
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
