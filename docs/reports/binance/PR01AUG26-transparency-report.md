# BINANCE PR01AUG26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR01AUG26-assessment.json](./PR01AUG26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01AUG26` |
| Snapshot time | `2026-08-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR01AUG26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

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
| Previous snapshot | `2026-07-01T00:00:00Z` |
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
| `bapiSnapshot` | `543a24d426a809a3ebf9395b5aa75b6f0e6ac10f5f76158a93204f48714adbe3` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `9297db70fbfe071bc069f076ea72ec523636058bb5a9526ca1ea5309be5772c8` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260801.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x563dbc4463368e71556b9a04c6a4c6b8e5c9115c3c9ca3ca9637daa269a2adf1` |
| Verification bundle root | `0xf5b8e5c7ae6fd49bb056f12938f01caf54c25eb4fe9ad7f75155a331d4fb04b7` |
| Artifact bundle SHA-256 | `6391a3c15d1fba588dcd2616da8db8a4c5374a1f982a0beb0fea138d9353bb5d` |
| Verification bundle SHA-256 | `66d36665b633f5ede25d8ff80fb585c32305f7caffa773072e896e41228dde8f` |

Local bundle paths: [PR01AUG26.artifact-bundle.json](../../../artifacts/binance/PR01AUG26/bundles/PR01AUG26.artifact-bundle.json), [PR01AUG26.verification-bundle.v2.json](../../../artifacts/binance/PR01AUG26/bundles/PR01AUG26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `FAIL` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | summary has no coin rows |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0400` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.1653` |  |
| `onchain-balance-ledger` | `1.4` | `WARN` | `0.0030` | more than 25% of supported ledger rows failed to query |
| `onchain-balance-deposit` | `1.2` | `WARN` | `0.1363` | more than 25% of sampled deposit rows failed RPC |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR01AUG26.verification-bundle.v2.json` in the local artifact bundle.

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
| `csv-extras` | `longTailCoins` |  | 53 coins not in summary snapshot: [1INCH AAVE APT ARB ASTER BCH BNB BOME BTC BUSD CAKE CHR CHZ CRV DOGE DOT ENA ENJ ETH FDUSD FORM GRT HBAR HFT LINK LTC MASK NEAR OP PAXG PENDLE PEPE POL RLUSD S SHIB SOL SSV SUI TRUMP TUSD U UNI USD1 USDC USDE USDT WIF WLD WLFI XLM XRP ZEC] | informational — binance summary lists top coins only |

#### `btc-anchor` (`FAIL`)

Finding counts: FAIL 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `btc:962079` | `snapshotTime` | 2026-08-01T00:00:00Z | 2026-08-12T03:16:44Z | BTC block 962079 timestamp via cache; drift 267h16m44s (max 30m0s) |

#### `solvency-claim` (`FAIL`)

Summary: summary has no coin rows

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 1, WARN 1, UNVERIFIABLE 2, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xdccf3b77da55107280bd850ea519df3705d1a75a` | `BNB@BSC#113293987` | 4007.165184 | 4007.122197389732389964 | accounted on-chain balance != csv claim by 0.042986610267610036 (provider=https://bsc.rpc.sentio.xyz,cache,https://bsc-mainnet.nodereal.io/v1/64a9df0874fb4a93b9d0a3849de012d3) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x001ceb373c83ae75b9f5cf78fc2aba3e185d09e2` | `BNB@OPBNB#169701876` | 7757.679688 | 3.00042829591268743 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=7754.67925970408731257 (provider=https://opbnb-mainnet-rpc.bnbchain.org) |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `XLM\|BSC` | 6 | no native verifier for this (coin,network) pair yet |
| `XLM\|XLM` | 11 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 21, WARN 685, UNVERIFIABLE 5, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `BCH@ETH#25656292` | 1499.012 | 0 | on-chain balance != csv claim by 1499.012 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#25656292` | 1156.4841 | 400 | on-chain balance != csv claim by 756.4841 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `BUSD@MATIC#91224919` | 0.8 | 0 | on-chain balance != csv claim by 0.8 (provider=https://polygon.gateway.tenderly.co) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `BUSD@MATIC#91224919` | 156821.4267 | 0 | on-chain balance != csv claim by 156821.4267 (provider=https://polygon.gateway.tenderly.co) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `BUSD@OPTIMISM#154971811` | 4097.158029 | 0 | on-chain balance != csv claim by 4097.158029 (provider=https://mainnet.optimism.io) |
| `0x9f8c163cba728e99993abe7495f06c0a3c8ac8b9` | `BUSD@AVAXC#91716609` | 128031.2016 | 0 | on-chain balance != csv claim by 128031.2016 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `CAKE@ETH#25656292` | 419266.0948 | 0 | on-chain balance != csv claim by 419266.0948 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CAKE@ETH#25656292` | 24593.94636 | 0 | on-chain balance != csv claim by 24593.94636 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `CAKE@ETH#25656292` | 70588.57984 | 0 | on-chain balance != csv claim by 70588.57984 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CAKE@ETH#25656292` | 348849 | 0 | on-chain balance != csv claim by 348849 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `CAKE@ETH#25656292` | 61223.1052 | 0 | on-chain balance != csv claim by 61223.1052 (provider=https://mainnet.infura.io/v3/0ec3a9a31c814957abffb5f51da4ebfd) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `ENA@ETH#25656292` | 105666030.9 | 0 | on-chain balance != csv claim by 105666030.9 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `ENA@ETH#25656292` | 12697512.7 | 0 | on-chain balance != csv claim by 12697512.7 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `ENA@ETH#25656292` | 11964755.19 | 0 | on-chain balance != csv claim by 11964755.19 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ENA@ETH#25656292` | 27343446.22 | 0 | on-chain balance != csv claim by 27343446.22 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ENA@ETH#25656292` | 650011267 | 0 | on-chain balance != csv claim by 650011267 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `ENA@ETH#25656292` | 2852693.4 | 0 | on-chain balance != csv claim by 2852693.4 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `ENA@ETH#25656292` | 530115216 | 0 | on-chain balance != csv claim by 530115216 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `ENA@ETH#25656292` | 7182591.196 | 0 | on-chain balance != csv claim by 7182591.196 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ENA@ETH#25656292` | 25914254.22 | 0 | on-chain balance != csv claim by 25914254.22 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ENA@ETH#25656292` | 425315.119 | 0 | on-chain balance != csv claim by 425315.119 (provider=https://eth-mainnet.g.alchemy.com/v2/xALpTiRVSIOkY1V1FPmG52PYkKzHfhyk) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `TDqSquXBgUCLYvYC4XZgrprLK589dkhSCf` | `BUSD@TRX#84957211` | 111074.1095 | 3.000644606715447029 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); on-chain balance != csv claim by 111071.108855393284552971 (provider=https://api.trongrid.io) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `CHZ@CHZ2` |  |  | rpc error (provider=): https://chiliz-mainnet.public.blastapi.io -> HTTP 403 |
| `0x87917d879ba83ce3ada6e02d49a10c1ec1988062` | `CHZ@CHZ2` |  |  | rpc error (provider=): no usable provider |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `CHZ@CHZ2` |  |  | rpc error (provider=): no usable provider |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CHZ@CHZ2` |  |  | rpc error (provider=): no usable provider |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CHZ@CHZ2` |  |  | rpc error (provider=): no usable provider |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `CHZ@CHZ2` |  |  | rpc error (provider=): no usable provider |
| `0xa84fd90d8640fa63d194601e0b2d1c9094297083` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): Post "https://mainnet.era.zksync.io": net/http: TLS handshake timeout |
| `0x4aec0e98fc1fb55b9cc2faaa7a81acca42cb4e96` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): Post "https://mainnet.era.zksync.io": net/http: TLS handshake timeout |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x7aed074ca56f5050d5a2e512ecc5bf7103937d76` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0xd3a22590f8243f8e83ac230d1842c9af0404c4a1` | `ETH@MANTA` |  |  | rpc error (provider=): Post "https://pacific-rpc.manta.network/http": context deadline exceeded (Client.Timeout exceeded while awaiting headers) |


_670 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@BASE` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `Bridging in progress` | `ETH@ARBITRUM` |  |  | placeholder wallet address (bridging in progress) — not a queryable on-chain address |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#25656292` | 1000 | 0 | no WBTC at snapshot height; BTC likely native/off-chain custody — EVM WBTC balanceOf cannot verify this row |

#### `onchain-balance-ledger` (`WARN`)

Summary: more than 25% of supported ledger rows failed to query

Finding counts: WARN 341, PASS 1

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xc6ea9c27cd6737031631f6d4a4258f7ff87f3642da964a5fad49594847cce386/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xbdb53eb583ba02ab0606bdfc71b59a191400f75fb62f9df124494ab877cdfe2a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=939716070": dial tcp: lookup fullnode.mainnet.aptoslabs.com: no such host |
| `1M4mQ3oqSvf42uqxEDJVdTyDRR9ubQBmT8` | `BCH\|BCH` |  |  | ledger rpc error (provider=): Get "https://api.blockchair.com/bitcoin-cash/dashboards/address/1M4mQ3oqSvf42uqxEDJVdTyDRR9ubQBmT8": dial tcp: lookup api.blockchair.com: no such host |
| `1KuPikhUYtHz3fmSQ2UvotpUuN672NuEcm` | `BCH\|BCH` |  |  | ledger rpc error (provider=): Get "https://api.blockchair.com/bitcoin-cash/dashboards/address/1KuPikhUYtHz3fmSQ2UvotpUuN672NuEcm": dial tcp: lookup api.blockchair.com: no such host |


_326 additional `WARN` rows omitted; see verification bundle._

#### `onchain-balance-deposit` (`WARN`)

Summary: more than 25% of sampled deposit rows failed RPC

Finding counts: WARN 500, PASS 2

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `69EH4KKKUtYtK2LGsC4FUGznhYpyprePyX2zT7Dj9rmG` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/69EH4KKKUtYtK2LGsC4FUGznhYpyprePyX2zT7Dj9rmG/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `Ff1wYqLwVQou6T2VbSNNUbF7YimdJ34qFvAHKJwLc9nL` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/Ff1wYqLwVQou6T2VbSNNUbF7YimdJ34qFvAHKJwLc9nL/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `7V6QzuvhEqygKVgGH6iTC1oRcwfr2ncnQ47MSfuzT1Hs` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/7V6QzuvhEqygKVgGH6iTC1oRcwfr2ncnQ47MSfuzT1Hs/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `56xu8tvobateVojWAWFRSuTeWEphQCiHCHwdm5CFg1j1` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/56xu8tvobateVojWAWFRSuTeWEphQCiHCHwdm5CFg1j1/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `3LoHa5URjRooX7sHVxLguq27Nb2FWSjuqyZjrvqx8Add` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/3LoHa5URjRooX7sHVxLguq27Nb2FWSjuqyZjrvqx8Add/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `616Sm336RzgDM8LcWVbJq5rSDoY3FTiedt5AiesJmJZN` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/616Sm336RzgDM8LcWVbJq5rSDoY3FTiedt5AiesJmJZN/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `CUC8M5iX7Dimijqf2bUZ4NJPwgL4EosDRSGgj3oVEA9W` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/CUC8M5iX7Dimijqf2bUZ4NJPwgL4EosDRSGgj3oVEA9W/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `43CMZzahPhsi467CLTtCMo7eyGp2eZEDNRgQytDNrCBV` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/43CMZzahPhsi467CLTtCMo7eyGp2eZEDNRgQytDNrCBV/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `8wWAh8WR5N2fY3YiQDPE7mwC4C5EmCp3Zx3BGstwv2nm` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/8wWAh8WR5N2fY3YiQDPE7mwC4C5EmCp3Zx3BGstwv2nm/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `GMjfiJrsthiw3SwducRgeUCWyU3pvcLEbuS6PJiAjQLY` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/GMjfiJrsthiw3SwducRgeUCWyU3pvcLEbuS6PJiAjQLY/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `6Hhe1yBYSKM7ED4qZzfDwiXnW8CdebZvfGcNux8ZZAyS` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/6Hhe1yBYSKM7ED4qZzfDwiXnW8CdebZvfGcNux8ZZAyS/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `BSzGGZb21pv6LGxWj6rgYH6vuHr6h6JCnNiNCmZvjfqM` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/BSzGGZb21pv6LGxWj6rgYH6vuHr6h6JCnNiNCmZvjfqM/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `5yQNRRM1H3PdHzLhVwsNjN49qxMaWnzDpJk6cRhGbzkb` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/5yQNRRM1H3PdHzLhVwsNjN49qxMaWnzDpJk6cRhGbzkb/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `4KnmzRxTvZ6RdKfd8Lmskcdtnd3r1AoFxTj6kptuC9L3` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/4KnmzRxTvZ6RdKfd8Lmskcdtnd3r1AoFxTj6kptuC9L3/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |
| `2xZawYM3cCdDJA2mypxeKCx6WneTzZ58ENzKUUNUfzW7` | `BOME\|SOL` |  |  | deposit rpc error (provider=): solana history: Get "https://api.helius.xyz/v1/wallet/2xZawYM3cCdDJA2mypxeKCx6WneTzZ58ENzKUUNUfzW7/balance-at?mint=ukHH6c7mMyiWCf1b9pnWe25TSpkDDt3H5pQZgZ74J82&slot=436454956": dial tcp: lookup api.helius.xyz: no such host |


_485 additional `WARN` rows omitted; see verification bundle._

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR01AUG26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-08-01T00:00:00Z |
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
