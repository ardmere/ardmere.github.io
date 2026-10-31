# BINANCE PR01SEP26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR01SEP26-assessment.json](./PR01SEP26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01SEP26` |
| Snapshot time | `2026-09-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR01SEP26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

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
| Previous snapshot | `2026-08-01T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `7 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Older snapshot PR01OCT26 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `b53c6132fd91689454294ad5cad66168116400ca2f740f5b8bff75ab551ed564` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `4f4777b7cea57165d60563e3c1238ad6bc2e749bb39b684c06fca3971d804b35` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260901.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x86842b275a42f5ad31f45ba3fa522526843fe084c64d6cd1a2bc8b5079ba8117` |
| Verification bundle root | `0x72d7e998e36e8cd3f4d22328bcbffbcc9ccec4e3537996e858b010a27226bbc4` |
| Artifact bundle SHA-256 | `090be3e0d4592d4670dc72e5f4eb829d7dc3eaead75b9b1cdeac910321d0ed19` |
| Verification bundle SHA-256 | `5db23d846fb45ea54c9d62e26cbd0a8ce7ee7748882044fb08a2b30da624c671` |

Local bundle paths: [PR01SEP26.artifact-bundle.json](../../../artifacts/binance/PR01SEP26/bundles/PR01SEP26.artifact-bundle.json), [PR01SEP26.verification-bundle.v2.json](../../../artifacts/binance/PR01SEP26/bundles/PR01SEP26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | summary has no coin rows |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0396` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.6613` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.1268` |  |
| `onchain-balance-deposit` | `1.2` | `PASS` | `0.1619` |  |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR01SEP26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `internal-consistency` (`PASS`)

Finding counts: WARN 1, PASS 15

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `csv-extras` | `longTailCoins` |  | 54 coins not in summary snapshot: [1INCH AAVE APT ARB ASTER BCH BNB BOME BTC BUSD CAKE CHR CHZ CRV DOGE DOT ENA ENJ ETH FDUSD FORM GRT HBAR HFT LINK LTC MASK NEAR OP PAXG PENDLE PEPE POL PUMP RLUSD S SHIB SOL SSV SUI TRUMP TUSD U UNI USD1 USDC USDE USDT WIF WLD WLFI XLM XRP ZEC] | informational — binance summary lists top coins only |

#### `solvency-claim` (`FAIL`)

Summary: summary has no coin rows

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 1, WARN 1, UNVERIFIABLE 3, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xb2a6A72843db0F508204A56448413f3867Ea691a` | `BNB@BSC#119243645` | 23000.06899 | 6000.068986839 | accounted on-chain balance != csv claim by 17000.000003161 (provider=https://bsc.rpc.sentio.xyz,cache,https://bsc-mainnet.nodereal.io/v1/64a9df0874fb4a93b9d0a3849de012d3) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x001ceb373c83ae75b9f5cf78fc2aba3e185d09e2` | `BNB@OPBNB#180415476` | 7623.741552 | 3.00042829591268743 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=7620.74112370408731257 (provider=https://opbnb-mainnet-rpc.bnbchain.org) |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `PUMP\|SOL` | 10 | no native verifier for this (coin,network) pair yet |
| `XLM\|BSC` | 6 | no native verifier for this (coin,network) pair yet |
| `XLM\|XLM` | 11 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 58, WARN 19, UNVERIFIABLE 11, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x3931dab967c3e2dbb492fe12460a66d0fe4cc857` | `ARB@ARBITRUM#500455224` | 10819440.9 | 10803990.68029034864000375 | on-chain balance != csv claim by 15450.21970965135999625 (provider=https://arb1.arbitrum.io/rpc) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `BCH@ETH#25878704` | 1499.012 | 0 | on-chain balance != csv claim by 1499.012 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x28e8d4eab37cfbf49edf7ff39c27377b7ae3a2a0` | `ASTER@BSC#119243645` | 15000045 | 35 | on-chain balance != csv claim by 15000010 (provider=https://bsc.rpc.sentio.xyz) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#25878704` | 1156.4841 | 400 | on-chain balance != csv claim by 756.4841 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `BUSD@MATIC#93010514` | 156880.6127 | 0 | on-chain balance != csv claim by 156880.6127 (provider=https://polygon.gateway.tenderly.co) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `BUSD@OPTIMISM#156311011` | 4097.158029 | 0 | on-chain balance != csv claim by 4097.158029 (provider=https://mainnet.optimism.io) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `BUSD@MATIC#93010514` | 0.8 | 0 | on-chain balance != csv claim by 0.8 (provider=https://polygon.gateway.tenderly.co) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `CAKE@ETH#25878704` | 368217.8018 | 0 | on-chain balance != csv claim by 368217.8018 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CAKE@ETH#25878704` | 348849 | 0 | on-chain balance != csv claim by 348849 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CAKE@ETH#25878704` | 24593.94636 | 0 | on-chain balance != csv claim by 24593.94636 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `CAKE@ETH#25878704` | 75000 | 0 | on-chain balance != csv claim by 75000 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `CAKE@ETH#25878704` | 53461.19192 | 0 | on-chain balance != csv claim by 53461.19192 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x9f8c163cba728e99993abe7495f06c0a3c8ac8b9` | `BUSD@AVAXC#94159927` | 128031.2016 | 0 | on-chain balance != csv claim by 128031.2016 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `ENA@ETH#25878704` | 154092094.8 | 0 | on-chain balance != csv claim by 154092094.8 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ENA@ETH#25878704` | 27343446.22 | 0 | on-chain balance != csv claim by 27343446.22 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ENA@ETH#25878704` | 550000000 | 0 | on-chain balance != csv claim by 550000000 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `ENA@ETH#25878704` | 530115216 | 0 | on-chain balance != csv claim by 530115216 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `ENA@ETH#25878704` | 17591004.15 | 0 | on-chain balance != csv claim by 17591004.15 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `ENA@ETH#25878704` | 21607375.64 | 0 | on-chain balance != csv claim by 21607375.64 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `ENA@ETH#25878704` | 3023937.431 | 0 | on-chain balance != csv claim by 3023937.431 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ENA@ETH#25878704` | 37886440.55 | 0 | on-chain balance != csv claim by 37886440.55 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `ENA@ETH#25878704` | 6349900.76 | 0 | on-chain balance != csv claim by 6349900.76 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ENA@ETH#25878704` | 3893767.901 | 0 | on-chain balance != csv claim by 3893767.901 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `GRT@ETH#25878704` | 1315581.858 | 0 | on-chain balance != csv claim by 1315581.858 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `GRT@ETH#25878704` | 3860556.119 | 0 | on-chain balance != csv claim by 3860556.119 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `GRT@ETH#25878704` | 42.12983058 | 0 | on-chain balance != csv claim by 42.12983058 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `GRT@ETH#25878704` | 23.26231794 | 0 | on-chain balance != csv claim by 23.26231794 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `GRT@ETH#25878704` | 897665.8017 | 0 | on-chain balance != csv claim by 897665.8017 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@ETH#25878704` | 13473559.78 | 0 | on-chain balance != csv claim by 13473559.78 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `GRT@ETH#25878704` | 0.092126 | 0 | on-chain balance != csv claim by 0.092126 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4aec0e98fc1fb55b9cc2faaa7a81acca42cb4e96` | `GRT@ETH#25878704` | 189374.7477 | 0 | on-chain balance != csv claim by 189374.7477 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `GRT@MATIC#93010514` | 776293.5762 | 0 | on-chain balance != csv claim by 776293.5762 (provider=https://polygon.gateway.tenderly.co) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@MATIC#93010514` | 73.1902 | 0 | on-chain balance != csv claim by 73.1902 (provider=https://polygon.gateway.tenderly.co) |
| `0x605d0b92cde1e110d5ae9123969097ea2489d1aa` | `GRT@ARBITRUM#500455224` | 6000930.84 | 113 | on-chain balance != csv claim by 6000817.84 (provider=https://arb1.arbitrum.io/rpc) |
| `0x5ce4825e2c5d7025b110d7bbeb5b1f03976e1443` | `GRT@ARBITRUM#500455224` | 70000100 | 100 | on-chain balance != csv claim by 70000000 (provider=https://arbitrum.gateway.tenderly.co) |
| `0xa64b436964e7415c0e70b9989a53e1fb9a90e726` | `POL@ETH#25878704` | 152000099.1 | 49.065998286027503462 | POL+MATIC ERC20 on ETH still << CSV; residual may be Polygon-native POL or internal ledger; on-chain balance != csv claim by 152000050.034001713972496538 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xa4e471dbfe8c95d4c44f520b19cee436c01c3267` | `S@SONIC#78496667` | 39059148.26 | 39059067.838049811619039375 | on-chain balance != csv claim by 80.421950188380960625 (provider=https://rpc.soniclabs.com) |
| `0x26BFf98f6f37c8d560684eAF2AFd1149bF76633E` | `SSV@ETH#25878704` | 300019 | 14 | on-chain balance != csv claim by 300005 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `TUSD@ETH#25878704` | 55254.20406 | 0 | on-chain balance != csv claim by 55254.20406 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `TUSD@ETH#25878704` | 7385.587173 | 0 | on-chain balance != csv claim by 7385.587173 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `TUSD@ETH#25878704` | 5225.631036 | 0 | on-chain balance != csv claim by 5225.631036 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@ETH#25878704` | 31703.33056 | 0 | on-chain balance != csv claim by 31703.33056 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `TUSD@ETH#25878704` | 1077.401833 | 0 | on-chain balance != csv claim by 1077.401833 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `TUSD@BSC#119243645` | 15000 | 0 | on-chain balance != csv claim by 15000 (provider=https://bsc.rpc.sentio.xyz) |
| `0xe2fc31f816a9b94326492132018c3aecc4a93ae1` | `TUSD@BSC#119243645` | 57249.14952 | 0.018991890000006253 | on-chain balance != csv claim by 57249.130528109999993747 (provider=https://bsc.rpc.sentio.xyz) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `TUSD@BSC#119243645` | 256388.6563 | 0 | on-chain balance != csv claim by 256388.6563 (provider=https://bsc.rpc.sentio.xyz) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@BSC#119243645` | 60337.586 | 0 | on-chain balance != csv claim by 60337.586 (provider=https://bsc.rpc.sentio.xyz) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@BSC#119243645` | 896917.499 | 0 | on-chain balance != csv claim by 896917.499 (provider=https://bsc.rpc.sentio.xyz) |
| `0x8894e0a0c962cb723c1976a4421c95949be2d4e3` | `TUSD@BSC#119243645` | 312267.842 | 153639.229140440670720258 | on-chain balance != csv claim by 158628.612859559329279742 (provider=https://bsc.rpc.sentio.xyz) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@BSC#119243645` | 9.9151 | 0 | on-chain balance != csv claim by 9.9151 (provider=https://bsc.rpc.sentio.xyz) |
| `0x59fbc084aeed2c8e6e9b6ebe8a11d443efcb40d2` | `USD1@AB#14503665` | 15423241.92 | 0 | on-chain balance != csv claim by 15423241.92 (provider=https://api.mainnet.abs.xyz) |
| `0xe2fc31f816a9b94326492132018c3aecc4a93ae1` | `USDC@BSC#119243645` | 38058315.83 | 38058115.728915971180723677 | on-chain balance != csv claim by 200.101084028819276323 (provider=https://bsc.rpc.sentio.xyz) |
| `0xee7ae85f2fe2239e27d9c1e23fffe168d63b4055` | `USDC@SEIEVM#229479307` | 544986.8531 | 3.000109 | on-chain balance != csv claim by 544983.852991 (provider=https://evm-rpc.sei-apis.com) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `USDT@OPTIMISM#156311011` | 15573155.45 | 0 | on-chain balance != csv claim by 15573155.45 (provider=https://mainnet.optimism.io) |
| `0x25681ab599b4e2ceea31f8b498052c53fc2d74db` | `USDT@ARBITRUM#500455224` | 7198025.11 | 7198020.409594 | on-chain balance != csv claim by 4.700406 (provider=https://arb1.arbitrum.io/rpc) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `USDT@OPTIMISM#156311011` | 1783.335686 | 0 | on-chain balance != csv claim by 1783.335686 (provider=https://mainnet.optimism.io) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `USDT@OPTIMISM#156311011` | 187097.6339 | 0 | on-chain balance != csv claim by 187097.6339 (provider=https://mainnet.optimism.io) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `USDT@OPTIMISM#156311011` | 73341107.18 | 0 | on-chain balance != csv claim by 73341107.18 (provider=https://mainnet.optimism.io) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `BUSD@TRX#85849746` | 111074.1095 | 3.000644606715447029 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 111071.108855393284552971 (provider=https://api.trongrid.io) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@ETH#25878704` | 3072391.422 | 0 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=3072391.422 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@ETH#25878704` | 2728912.638 | 0 | on-chain ERC20 balance << CSV mega-stablecoin allocation; likely omnibus/internal ledger label — not single-address custody; delta=2728912.638 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `TUSD@TRX#85849746` | 1550578.963 | 3.000412606533647589 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 1550575.962587393466352411 (provider=https://api.trongrid.io) |
| `TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G` | `TUSD@TRX#85849746` | 991929.56 | 0.75 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 991928.81 (provider=https://api.trongrid.io) |
| `TNXoiAJ3dct8Fjg4M9fkLFh9S2v9TXc32G` | `U@TRX#85849746` | 46495.74655 | 3.536147 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 46492.210403 (provider=https://api.trongrid.io) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `U@TRX#85849746` | 241189.186 | 3.481774710853280259 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 241185.704225289146719741 (provider=https://api.trongrid.io) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `USDC@TRX#85849746` | 79913.95874 | 600.000933 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 79313.957807 (provider=https://tron-rpc.publicnode.com) |
| `TYASr5UV6HEcXatwdFQfmLVUqQQQMUxHLS` | `USDT@TRX#85849746` | 49959069.41 | 0.100282 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 49959069.309718 (provider=https://api.trongrid.io) |
| `TJDENsfBJs4RFETt1X1W8wMDc8M5XnJhCe` | `USDT@TRX#85849746` | 46788664.35 | 98.5 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 46788565.85 (provider=https://api.trongrid.io) |
| `TK4ykR48cQQoyFcZ5N4xZCbsBaHcg6n3gJ` | `USDT@TRX#85849746` | 38987865.48 | 0.100853 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 38987865.379147 (provider=https://api.trongrid.io) |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `USDT@TRX#85849746` | 708267249.6 | 2.000197 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 708267247.599803 (provider=https://tron-rpc.publicnode.com) |
| `TJ5usJLLwjwn7Pw3TPbdzreG7dvgKzfQ5y` | `USDT@TRX#85849746` | 40525326.24 | 0 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 40525326.24 (provider=https://tron-rpc.publicnode.com) |
| `TJqwA7SoZnERE4zW5uDEiPkbz4B66h9TFj` | `USDT@TRX#85849746` | 34674262.43 | 0.1 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 34674262.33 (provider=https://tron-rpc.publicnode.com) |
| `TQrY8tryqsYVCYS3MFbtffiPp2ccyn4STm` | `USDT@TRX#85849746` | 42678967.5 | 0.10034 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 42678967.39966 (provider=https://tron-rpc.publicnode.com) |


_4 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `Bridging in progress` | `ETH@OPTIMISM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `GRT@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `GRT@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `GRT@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `GRT@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `WLD@WLD` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#25878704` | 1000 | 0 | no WBTC at snapshot height; BTC likely native/off-chain custody — EVM WBTC balanceOf cannot verify this row |

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 12, WARN 176, UNVERIFIABLE 1, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `BOME\|SOL#443282097` | 2648114223 | 2646589890.650728 | on-chain balance != csv claim by 1524332.349272 (provider=https://api.helius.xyz/v1/wallet) |
| `0.0.1030878` | `HBAR\|HBAR#99493123` | 125878685 | 125546267.67528923 | on-chain balance != csv claim by 332417.32471077 (provider=https://mainnet-public.mirrornode.hedera.com/api/v1/accounts/0.0.1030878?timestamp=lte:1788220798.345897104) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `SOL\|SOL#443282097` | 1790580.231 | 1790571.630786998 | on-chain balance != csv claim by 8.600213002 (provider=https://api.helius.xyz/v1/wallet) |
| `c5f9zfpkKMD9N8uLqJcFeJAAz7v12vDMnup9Y6EeQkk` | `SOL\|SOL#443282097` | 4235909.878 | 1.074220732 | on-chain balance != csv claim by 4235908.803779268 (provider=https://api.helius.xyz/v1/wallet) |
| `GC5LF63GRVIT5ZXXCXLPI3RX2YXKJQFZVBSAO6AUELN3YIMSWPD6Z6FH` | `USDC\|XLM#64216184` | 99396596.73 | 3.0007133 | on-chain balance != csv claim by 99396593.7292867 (provider=https://horizon.stellar.org/accounts/GC5LF63GRVIT5ZXXCXLPI3RX2YXKJQFZVBSAO6AUELN3YIMSWPD6Z6FH?at=64216184) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDC\|SOL#443282097` | 365841869 | 363894235.972549 | on-chain balance != csv claim by 1947633.027451 (provider=https://api.helius.xyz/v1/wallet) |
| `GBAIA5U6E3FSRUW55AXACIVGX2QR5JYAS74OWLED3S22EGXVYEHPLGPA` | `USDC\|XLM#64216184` | 1675167.198 | 1235907.292 | on-chain balance != csv claim by 439259.906 (provider=https://horizon.stellar.org/accounts/GBAIA5U6E3FSRUW55AXACIVGX2QR5JYAS74OWLED3S22EGXVYEHPLGPA?at=64216184) |
| `5SDrsMNTYdhmApjfqYHDvjoW92f2S42vcc7zNDVcQ9Ej` | `USDC\|SOL#443282097` | 581765.9675 | 581764.436339 | on-chain balance != csv claim by 1.531161 (provider=https://api.helius.xyz/v1/wallet) |
| `9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM` | `USDC\|SOL#443282097` | 1094.24345 | 1093.24345 | on-chain balance != csv claim by 1 (provider=https://api.helius.xyz/v1/wallet) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDT\|SOL#443282097` | 476368442.4 | 476349731.652257 | on-chain balance != csv claim by 18710.747743 (provider=https://api.helius.xyz/v1/wallet) |
| `tz2WDATNYnp7FdsmuZDYSidioZqeoLNZqXvE` | `USDT\|XTZ#14748075` | 2937517.122 | 3.000651 | on-chain balance != csv claim by 2937514.121349 (provider=https://api.tzkt.io/v1/tokens/balances?account=tz2WDATNYnp7FdsmuZDYSidioZqeoLNZqXvE&token.contract=KT1XnTn74bUtxHfDtBmm2bGZAQfhPbvKWR8o&level=14748075) |
| `tz1Q3jvYU9knekDYJfyvj3GjUy6898MNjvb2` | `USDT\|XTZ#14748075` | 14579292.92 | 2290793.961 | on-chain balance != csv claim by 12288498.959 (provider=https://api.tzkt.io/v1/tokens/balances?account=tz1Q3jvYU9knekDYJfyvj3GjUy6898MNjvb2&token.contract=KT1XnTn74bUtxHfDtBmm2bGZAQfhPbvKWR8o&level=14748075) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1012869185): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1012869185) has been pruned","hint":"The requested data has bee |
| `181Hj4Ek3BKy6g6MJgC7y4dBuGCfb7vLgr` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/181Hj4Ek3BKy6g6MJgC7y4dBuGCfb7vLgr): blockchair HTTP 430 |
| `19dQkvaH2NGgkGomzZu3qrnqRGCicXwedM` | `BCH\|BCH` |  |  | ledger rpc error (provider=https://api.blockchair.com/bitcoin-cash/dashboards/address/19dQkvaH2NGgkGomzZu3qrnqRGCicXwedM): blockchair HTTP 430 |


_161 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `bc1q9n7j6a5wwru6tad9xrlncnagresvyjt96pl3en` | `BTC\|BTC#964957` | 1200.019572 | 0.01957174 | on-chain UTXO sum << CSV; likely omnibus/internal BTC custody not visible at this address |

#### `onchain-balance-deposit` (`PASS`)

Finding counts: WARN 56, PASS 2

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `BSzGGZb21pv6LGxWj6rgYH6vuHr6h6JCnNiNCmZvjfqM` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `239uubEyVqrH7s3fZfCXRnPfUQZzD6z9QEDJSwRPhHmG` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `2MvJsKMw33D6kJN24D3MgDoW48mho3mUNP2Qom2PDpPv` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `6ncSuUvMAfomF3NQLadHezFUrhDSgxi4B6sZB82CtnmN` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `EPf2vCzPXrBwLrYyeWGVkCQrMDMXosvK2ajpDWjYX7gb` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `724eNdouTiQF4mwUY8QG1q81jbkQ2UXudwYTvL5VA8LH` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `AMMUbKU1PPjmidNyknt3g9Q4RUFLUvSuDK94mH2k11rG` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `CzWy1LF4XZaZ8NjoNUrqcV2ACggwdKZnNMf9i6GT1QHm` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `Dw362osAz6pb9S2d1GoKZa13zxxeqyth3TEuWhsLEZEo` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `DQByzsbUrDSyeE5gNVpSXSdy51AG9esykvj6WHtvcjbw` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `CNe5HBgpxFCg3eantqUjVFmBPin6gtMk3zej6ijhHuYn` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `GiwwjZr1JrPfCgrv4o3WjtRN4SZ71LGLXeoYFtxG6eAL` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `EV1TCHkue8dc6uLoJK79SCW8XV3y29NJ6aQyznm4DF25` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `47E97pM3ceUjUeudF5a2nnRSiDU3XxDEY2yMBjEqqg9T` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |
| `CWBUBaNSamFnnzQkQFCPGxDuWDeQe23tnKq1sA5z8y7A` | `BOME\|SOL` |  |  | deposit rpc error (provider=https://api.helius.xyz/v1/wallet): solana history HTTP 429: Too Many Requests |


_41 additional `WARN` rows omitted; see verification bundle._

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR01SEP26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-09-01T00:00:00Z |
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
