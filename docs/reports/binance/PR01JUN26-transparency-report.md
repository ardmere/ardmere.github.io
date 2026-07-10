# BINANCE PR01JUN26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR01JUN26-assessment.json](./PR01JUN26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01JUN26` |
| Snapshot time | `2026-06-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR01JUN26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

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
| Latest snapshot (evaluation set) | `2026-06-01T00:00:00Z` |
| Previous snapshot | `2026-05-01T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `3 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for binance.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `7898c1147f470b08404618c7f485664a4cb91740be3415d469d52f23c64f4202` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `722466eec2376e28cfeb825eeee1a08a50b7459e640d63f78a506386de65c8b2` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260601.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xf80414370da052e721aabd9845f917e0aa944e1e0e50de4ab77210a46d8b02ca` |
| Verification bundle root | `0xa919262721abc4ab660633c288e0c65161e89931781dc858f7538e92e7682807` |
| Artifact bundle SHA-256 | `2752c789500d3a642aa8b8a8cf0169ad6c99701a64e5a0d93eebcad2a18ffe52` |
| Verification bundle SHA-256 | `c4263772299ce76cea419bde6c13973033508e0eef41e9ae984cfdbaf809663c` |

Local bundle paths: [PR01JUN26.artifact-bundle.json](../../../artifacts/binance/PR01JUN26/bundles/PR01JUN26.artifact-bundle.json), [PR01JUN26.verification-bundle.v2.json](../../../artifacts/binance/PR01JUN26/bundles/PR01JUN26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0428` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.6815` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.2147` |  |
| `onchain-balance-deposit` | `1.2` | `PASS` | `0.1958` |  |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR01JUN26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 1, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x73f5ebe90f27b46ea12e5795d16c4b408b19cc6f` | `BNB@BSC#101590091` | 4634.557026 | 4634.555633295605962181 | accounted on-chain balance != csv claim by 0.001392704394037819 (provider=cache) |

#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 52, WARN 12, UNVERIFIABLE 7, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `GRT@ARBITRUM#468748167` | 879568408.4 | 819176091 | on-chain balance != csv claim by 60392317.4 (provider=cache) |
| `0x9f8c163cba728e99993abe7495f06c0a3c8ac8b9` | `BUSD@AVAXC#86865826` | 129697.2016 | 0 | on-chain balance != csv claim by 129697.2016 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `TUSD@BSC#101590091` | 15000 | 0 | on-chain balance != csv claim by 15000 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@BSC#101590091` | 60337.586 | 0 | on-chain balance != csv claim by 60337.586 (provider=cache) |
| `0xe2fc31f816a9b94326492132018c3aecc4a93ae1` | `TUSD@BSC#101590091` | 38824.78879 | 0.018991890000006253 | on-chain balance != csv claim by 38824.769798109999993747 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `TUSD@BSC#101590091` | 256388.6563 | 0 | on-chain balance != csv claim by 256388.6563 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@BSC#101590091` | 896917.499 | 0 | on-chain balance != csv claim by 896917.499 (provider=cache) |
| `0x8894e0a0c962cb723c1976a4421c95949be2d4e3` | `TUSD@BSC#101590091` | 471103.6579 | 69272.225457565597150464 | on-chain balance != csv claim by 401831.432442434402849536 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@BSC#101590091` | 9.9151 | 0 | on-chain balance != csv claim by 9.9151 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `BCH@ETH#25218797` | 1499.012 | 0 | on-chain balance != csv claim by 1499.012 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#25218797` | 1457.247 | 1300 | on-chain balance != csv claim by 157.247 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CAKE@ETH#25218797` | 348849 | 0 | on-chain balance != csv claim by 348849 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `CAKE@ETH#25218797` | 318601.2265 | 0 | on-chain balance != csv claim by 318601.2265 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `CAKE@ETH#25218797` | 57792.54209 | 0 | on-chain balance != csv claim by 57792.54209 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `CAKE@ETH#25218797` | 72191.8459 | 0 | on-chain balance != csv claim by 72191.8459 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CAKE@ETH#25218797` | 24593.94636 | 0 | on-chain balance != csv claim by 24593.94636 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `ENA@ETH#25218797` | 10655061.06 | 0 | on-chain balance != csv claim by 10655061.06 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `ENA@ETH#25218797` | 62366413.91 | 0 | on-chain balance != csv claim by 62366413.91 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `ENA@ETH#25218797` | 9856023.17 | 0 | on-chain balance != csv claim by 9856023.17 (provider=cache) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ENA@ETH#25218797` | 27343446.22 | 0 | on-chain balance != csv claim by 27343446.22 (provider=cache) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `ENA@ETH#25218797` | 310000000 | 0 | on-chain balance != csv claim by 310000000 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `ENA@ETH#25218797` | 2123113.539 | 0 | on-chain balance != csv claim by 2123113.539 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ENA@ETH#25218797` | 1166631.31 | 0 | on-chain balance != csv claim by 1166631.31 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ENA@ETH#25218797` | 762126001 | 0 | on-chain balance != csv claim by 762126001 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ENA@ETH#25218797` | 23619795.32 | 0 | on-chain balance != csv claim by 23619795.32 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `GRT@ETH#25218797` | 28835942.1 | 0 | on-chain balance != csv claim by 28835942.1 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `GRT@ETH#25218797` | 10247710.51 | 0 | on-chain balance != csv claim by 10247710.51 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `GRT@ETH#25218797` | 1315581.858 | 0 | on-chain balance != csv claim by 1315581.858 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `GRT@ETH#25218797` | 6171751.875 | 0 | on-chain balance != csv claim by 6171751.875 (provider=cache) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `GRT@ETH#25218797` | 2849341.533 | 0 | on-chain balance != csv claim by 2849341.533 (provider=cache) |
| `0x4aec0e98fc1fb55b9cc2faaa7a81acca42cb4e96` | `GRT@ETH#25218797` | 189374.7477 | 0 | on-chain balance != csv claim by 189374.7477 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `GRT@ETH#25218797` | 897665.8017 | 0 | on-chain balance != csv claim by 897665.8017 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@ETH#25218797` | 17422667.87 | 0 | on-chain balance != csv claim by 17422667.87 (provider=cache) |
| `0xa64b436964e7415c0e70b9989a53e1fb9a90e726` | `POL@ETH#25218797` | 160000099.1 | 49.065528965081548221 | POL+MATIC ERC20 on ETH still << CSV; residual may be Polygon-native POL or internal ledger; on-chain balance != csv claim by 160000050.034471034918451779 (provider=cache) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `TUSD@ETH#25218797` | 419177.1617 | 0 | on-chain balance != csv claim by 419177.1617 (provider=cache) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `TUSD@ETH#25218797` | 54703.3 | 0 | on-chain balance != csv claim by 54703.3 (provider=cache) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `TUSD@ETH#25218797` | 37540.74557 | 0 | on-chain balance != csv claim by 37540.74557 (provider=cache) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `TUSD@ETH#25218797` | 1077.401833 | 0 | on-chain balance != csv claim by 1077.401833 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@ETH#25218797` | 55418.60361 | 0 | on-chain balance != csv claim by 55418.60361 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@MATIC#87732730` | 73.1902 | 0 | on-chain balance != csv claim by 73.1902 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `BUSD@MATIC#87732730` | 0.8 | 0 | on-chain balance != csv claim by 0.8 (provider=cache) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `BUSD@MATIC#87732730` | 156622.6274 | 0 | on-chain balance != csv claim by 156622.6274 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `GRT@MATIC#87732730` | 6034171 | 0 | on-chain balance != csv claim by 6034171 (provider=cache) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `GRT@MATIC#87732730` | 9466.485118 | 0 | on-chain balance != csv claim by 9466.485118 (provider=cache) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `BUSD@OPTIMISM#152336611` | 3671.672463 | 0 | on-chain balance != csv claim by 3671.672463 (provider=cache) |
| `0xee7ae85f2fe2239e27d9c1e23fffe168d63b4055` | `USDC@OPTIMISM#152336611` | 9107611.187 | 9107397.23734 | on-chain balance != csv claim by 213.94966 (provider=cache) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `USDT@OPTIMISM#152336611` | 27784190.89 | 0 | on-chain balance != csv claim by 27784190.89 (provider=cache) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `USDT@OPTIMISM#152336611` | 20000000 | 0 | on-chain balance != csv claim by 20000000 (provider=cache) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `USDT@OPTIMISM#152336611` | 1774.870765 | 0 | on-chain balance != csv claim by 1774.870765 (provider=cache) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `USDT@OPTIMISM#152336611` | 187097.6339 | 0 | on-chain balance != csv claim by 187097.6339 (provider=cache) |
| `0xb32e9a84ae0b55b8ab715e4ac793a61b277bafa3` | `USDC@RON#56401606` | 1728823.962 | 1728809.322359 | on-chain balance != csv claim by 14.639641 (provider=cache) |
| `0x64de13c46f627d9c86212050d48756fb65c06d8a` | `S@SONIC#72001540` | 82000048.8 | 76000048.80334127621235255 | on-chain balance != csv claim by 5999999.99665872378764745 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@ETH#25218797` | 3030595.297 | 0 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=3030595.297 (provider=cache) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `TUSD@TRX#83201055` | 53587.60406 | 30874.917585402395005971 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 22712.686474597604994029 (provider=https://api.trongrid.io) |
| `TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y` | `USDT@TRX#83201055` | 49025134.17 | 30375704.237039 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 18649429.932961 (provider=https://tron-rpc.publicnode.com) |
| `TJqwA7SoZnERE4zW5uDEiPkbz4B66h9TFj` | `USDT@TRX#83201055` | 37833490.54 | 30171905.601535 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 7661584.938465 (provider=https://tron-rpc.publicnode.com) |
| `TYASr5UV6HEcXatwdFQfmLVUqQQQMUxHLS` | `USDT@TRX#83201055` | 53860294.75 | 47913121.260032 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 5947173.489968 (provider=https://tron-rpc.publicnode.com) |
| `TJDENsfBJs4RFETt1X1W8wMDc8M5XnJhCe` | `USDT@TRX#83201055` | 49185201.95 | 35455161.590048 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 13730040.359952 (provider=https://tron-rpc.publicnode.com) |
| `TCLgK89AnXbC9rewvhNb9UgXCc2qJJpBXh` | `USDT@TRX#83201055` | 57078438.35 | 53133547.857503 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 3944890.492497 (provider=https://tron-rpc.publicnode.com) |
| `TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G` | `USDT@TRX#83201055` | 47573015.43 | 27288904.534162 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 20284110.895838 (provider=https://tron-rpc.publicnode.com) |
| `TQrY8tryqsYVCYS3MFbtffiPp2ccyn4STm` | `USDT@TRX#83201055` | 31361706.05 | 26322097.582938 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 5039608.467062 (provider=https://api.trongrid.io) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `USDT@TRX#83201055` | 407520197.9 | 263424522.160267 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 144095675.739733 (provider=https://tron-rpc.publicnode.com) |
| `TAzsQ9Gx8eqFNFSKbeXrbi45CuVPHzA8wr` | `USDT@TRX#83201055` | 45331338.13 | 31845301.085686 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 13486037.044314 (provider=https://tron-rpc.publicnode.com) |
| `TWd4WrZ9wn84f5x1hZhL4DHvk738ns5jwb` | `USDT@TRX#83201055` | 1857436958 | 0 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 1857436958 (provider=https://tron-rpc.publicnode.com) |

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#25218797` | 699.2371 | 0 | no WBTC at snapshot height; BTC likely native/off-chain custody — EVM WBTC balanceOf cannot verify this row |

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 14, WARN 49, UNVERIFIABLE 9, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `CAKE\|APT#800530436` | 3332037.46 | 800000 | on-chain balance != csv claim by 2532037.46 (provider=cache) |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `USDC\|APT#800530436` | 31294357.34 | 0 | on-chain balance != csv claim by 31294357.34 (provider=cache) |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `USDC\|APT#800530436` | 1791897.02 | 0 | on-chain balance != csv claim by 1791897.02 (provider=cache) |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `USDT\|APT#800530436` | 173846627.7 | 0 | on-chain balance != csv claim by 173846627.7 (provider=cache) |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `USDT\|APT#800530436` | 286111302.8 | 0 | on-chain balance != csv claim by 286111302.8 (provider=cache) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `HFT\|SOL#423478906` | 451556.0691 | 450541.945588 | on-chain balance != csv claim by 1014.123512 (provider=https://api.helius.xyz) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `SOL\|SOL#423478906` | 2192861.712 | 2192798.807060176 | on-chain balance != csv claim by 62.904939824 (provider=https://api.helius.xyz) |
| `c5f9zfpkKMD9N8uLqJcFeJAAz7v12vDMnup9Y6EeQkk` | `SOL\|SOL#423478906` | 4719193.116 | 1.085865132 | on-chain balance != csv claim by 4719192.030134868 (provider=https://api.helius.xyz) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDC\|SOL#423478906` | 626600914.8 | 624685713.392221 | on-chain balance != csv claim by 1915201.407779 (provider=https://api.helius.xyz) |
| `9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM` | `USDC\|SOL#423478906` | 1089.417022 | 1088.417022 | on-chain balance != csv claim by 1 (provider=https://api.helius.xyz) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDT\|SOL#423478906` | 433237360.5 | 433219712.950744 | on-chain balance != csv claim by 17647.549256 (provider=https://api.helius.xyz) |
| `G9RCBaYb8aBRxoe8QBC2ucGrVqjuZFysRhY8d56cnNT1` | `USDT\|SOL#423478906` | 311489.4214 | 0.016752 | on-chain balance != csv claim by 311489.404648 (provider=https://api.helius.xyz) |
| `9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM` | `USDT\|SOL#423478906` | 20000022.59 | 22.585 | on-chain balance != csv claim by 20000000.005 (provider=https://api.helius.xyz) |
| `rhWj9gaovwu2hZxYW7p388P8GRbuXFLQkK` | `XRP\|XRP#104615295` | 4336748.088 | 4336731.398347 | on-chain balance != csv claim by 16.689653 (provider=cache) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `APT\|APT#800530436` | 4400633.125 | 2518326.31792997 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 1882306.80707003 (provider=cache) |
| `0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3` | `APT\|APT#800530436` | 79641.49137 | 0.003994 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 79641.487376 (provider=cache) |
| `0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3` | `APT\|APT#800530436` | 77790.29454 | 6.371591 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 77783.922949 (provider=cache) |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `APT\|APT#800530436` | 51625683.41 | 22665573.834647 | partial CoinStore+FA APT vs CSV row allocation; other custody forms not queried; on-chain balance != csv claim by 28960109.575353 (provider=cache) |
| `19dQkvaH2NGgkGomzZu3qrnqRGCicXwedM` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/19dQkvaH2NGgkGomzZu3qrnqRGCicXwedM?to=953422&details=basic): alchemy HTTP 403 |
| `1Gi1qU3esZcmJtPphC19ZpNpVNL8ZFWFsM` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1Gi1qU3esZcmJtPphC19ZpNpVNL8ZFWFsM?to=953422&details=basic): alchemy HTTP 403 |
| `1KuPikhUYtHz3fmSQ2UvotpUuN672NuEcm` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1KuPikhUYtHz3fmSQ2UvotpUuN672NuEcm?to=953422&details=basic): alchemy HTTP 403 |
| `1DHLUNgGMib8sB5G7U8VttsrXig7i8T6Gf` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1DHLUNgGMib8sB5G7U8VttsrXig7i8T6Gf?to=953422&details=basic): alchemy HTTP 403 |
| `19UKPC3A7L842ZzP4Rw4Vg4vQ1rrsvmuJ1` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/19UKPC3A7L842ZzP4Rw4Vg4vQ1rrsvmuJ1?to=953422&details=basic): alchemy HTTP 403 |
| `1KMcFqQxgJr7X9ADVGy28FSji4BuADkqUr` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1KMcFqQxgJr7X9ADVGy28FSji4BuADkqUr?to=953422&details=basic): alchemy HTTP 403 |
| `1Np9RvgpNve8hQ9H5oeGQ71T9uerDx32rr` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1Np9RvgpNve8hQ9H5oeGQ71T9uerDx32rr?to=953422&details=basic): alchemy HTTP 403 |
| `1P86nZCNWUiynP52AK2eTuTGZXYUTwX6qQ` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1P86nZCNWUiynP52AK2eTuTGZXYUTwX6qQ?to=953422&details=basic): alchemy HTTP 403 |
| `1NbpHNjzNeNSUfHoJFSDsTNBZ8PxSXwtjv` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1NbpHNjzNeNSUfHoJFSDsTNBZ8PxSXwtjv?to=953422&details=basic): alchemy HTTP 403 |
| `1JyNxYBkMSDRRaVzPZRzgJxpHkRoK28d7D` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1JyNxYBkMSDRRaVzPZRzgJxpHkRoK28d7D?to=953422&details=basic): alchemy HTTP 403 |
| `1B6cHqiZhcztHWriiyVrBbS7orHkSbjwbU` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://bitcoincash-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}/api/v2/address/1B6cHqiZhcztHWriiyVrBbS7orHkSbjwbU?to=953422&details=basic): alchemy HTTP 403 |


_34 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386` | `APT\|APT#800530436` | 1087020.95 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 1087020.95 (provider=cache) |
| `0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d` | `APT\|APT#800530436` | 5320002.889 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 5320002.889 (provider=cache) |
| `0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c` | `APT\|APT#800530436` | 402396.998 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 402396.998 (provider=cache) |
| `0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3` | `APT\|APT#800530436` | 4000000.299 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 4000000.299 (provider=cache) |
| `0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34` | `APT\|APT#800530436` | 1437737.908 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 1437737.908 (provider=cache) |
| `0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a` | `APT\|APT#800530436` | 14596.32301 | 0 | no AptosCoin/FungibleStore APT at ledger_version; account empty or CSV row is internal allocation; on-chain balance != csv claim by 14596.32301 (provider=cache) |
| `bc1p7x4aaws8t8cmmccu39kun6cajglx6ntuta7lhyxjt2cwrw4k89zqn0wley` | `BTC\|BTC#951913` | 2000.069978 | 0.0699781 | on-chain UTXO sum << CSV; likely omnibus/internal BTC custody not visible at this address |
| `DJfU2p6woQ9GiBdiXsWZWJnJ9uDdZfSSNC` | `DOGE\|DOGE#6230017` | 457117030.3 | 453845695.7881638 | address has 1502429 txs; exceeds public API limit 200000 — set BLOCKCHAIR_API_KEY or use archive provider |
| `LZEjckteAtWrugbsy9zU8VHEZ4iUiXo9Nm` | `LTC\|LTC#3117262` | 299054.098 | 299649.83998825 | address has 1595315 txs; exceeds public API limit 100000 — set BLOCKCHAIR_API_KEY or use archive provider |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR01JUN26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-06-01T00:00:00Z |
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
