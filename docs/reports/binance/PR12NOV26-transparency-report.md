# BINANCE PR12NOV26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR12NOV26-assessment.json](./PR12NOV26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR12NOV26` |
| Snapshot time | `2026-11-12T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR12NOV26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

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
| Total liabilities (USD) | `128582258733.22` | self-reported |
| Total reserves (USD) | `130015809491.98` | self-reported |
| Unpriced coins | `1INCH`, `ACH`, `AGLD`, `AIGENSYN`, `ALICE`, `ALLO`, `AMP`, `ANIME`, `ANKR`, `APE`, `API3`, `ARKM`, `ARPA`, `ASTER`, `AUCTION`, `AUDIO`, `AVA`, `AXL`, `AXS`, `BARD`, `BAT`, `BEAMX`, `BICO`, `BIGTIME`, `BIO`, `BLUR`, `BREV`, `CAKE`, `CELR`, `CFG`, `CHR`, `CHZ`, `COMP`, `COTI`, `COW`, `CTSI`, `CVC`, `CVX`, `CYBER`, `DEXE`, `DIA`, `DODO`, `DOLO`, `DUSK`, `EDEN`, `EIGEN`, `ENA`, `ENJ`, `ENS`, `ENSO`, `ERA`, `ESP`, `ETHFI`, `EUL`, `EURI`, `F`, `FET`, `FF`, `FORM`, `FTT`, `G`, `GALA`, `GLM`, `GRT`, `GTC`, `HEI`, `HEMI`, `HOT`, `HYPE`, `HYPER`, `ID`, `ILV`, `IMX`, `IQ`, `JASMY`, `KERNEL`, `KITE`, `KNC`, `LA`, `LQTY`, `LSK`, `MANA`, `MASK`, `MEME`, `MORPHO`, `NEIRO`, `NEWT`, `NEXO`, `NIL`, `NMR`, `NOM`, `OGN`, `ONDO`, `OPEN`, `PAXG`, `PENDLE`, `PHA`, `PLUME`, `PORTAL`, `POWR`, `PROM`, `PROVE`, `PUMP`, `PUNDIX`, `QKC`, `QNT`, `RAD`, `RARE`, `RE`, `RED`, `RESOLV`, `REZ`, `RLC`, `RLUSD`, `ROBO`, `RSR`, `S`, `SAHARA`, `SAND`, `SENT`, `SHELL`, `SIGN`, `SKL`, `SKY`, `SLP`, `SNX`, `SOLV`, `SOPH`, `SPELL`, `SPK`, `SSV`, `STRK`, `SUPER`, `SUSHI`, `SXT`, `SYN`, `SYRUP`, `T`, `TLM`, `TRB`, `TREE`, `TRUMP`, `TURBO`, `TURTLE`, `U`, `UMA`, `USD1`, `USDS`, `VIRTUAL`, `WBETH`, `WIF`, `WLD`, `WLFI`, `WOO`, `XAUT`, `XUSD`, `YFI`, `YGG`, `ZAMA`, `ZBT`, `ZEC`, `ZKC`, `ZKP`, `ZRO`, `ZRX` | — |

| Coin | Reserve ratio | Basis | Label |
| --- | --- | --- | --- |
| `BTC` | `1.00054022` | `amounts` | self-reported |
| `ETH` | `1.00020038` | `amounts` | self-reported |
| `USDT` | `1.04342673` | `amounts` | self-reported |
| `USDC` | `1.0007431` | `amounts` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it, or when neither amounts nor a reported rate were published (Unavailable). On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-11-12T00:00:00Z` |
| Previous snapshot | `2026-11-04T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `10 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for binance.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `294130128f956d3e30aab2b0e7eb18dc0bb675e8feb4b06efb3d0a37dd9ec582` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `86254575aa1b630b51048fa1b95bf65fb1098e0d866020e91b2396cebf3a6292` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20261111.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x977240cf1d233e02a4492d97c54dddbf8b3dedb3b0edc8274f615b33fa440714` |
| Verification bundle root | `0x2a147e401f3c28346bdb9d3421382316d330893aa3c5253ea6c1ea00b2cca8e7` |
| Artifact bundle SHA-256 | `e366f0c2938bb22ce3bbf8c1396c2c04f095113b4729662d23c4fce2fd14d4be` |
| Verification bundle SHA-256 | `71402ca6ddb227de79930c62d372d8778ec5b467992ca11d9632e67db7f71223` |

Local bundle paths: [PR12NOV26.artifact-bundle.json](../../../artifacts/binance/PR12NOV26/bundles/PR12NOV26.artifact-bundle.json), [PR12NOV26.verification-bundle.v2.json](../../../artifacts/binance/PR12NOV26/bundles/PR12NOV26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0372` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.3406` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.0280` |  |
| `onchain-balance-deposit` | `1.2` | `WARN` | `0.9775` | more than 25% of sampled deposit rows failed RPC |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR12NOV26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `onchain-balance-hot` (`FAIL`)

Finding counts: FAIL 8, UNVERIFIABLE 294, PASS 174

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x86523c87c8ec98c7539e2c58cd813ee9d1a08d96` | `BNB@BSC#133060467` | 173140.7487 | 173139.098538279702150564 | shortfall / 不足 (CSV exceeds on-chain) by 1.650161720297849436 (provider=https://bsc.rpc.sentio.xyz,cache,https://bsc-mainnet.nodereal.io/v1/[REDACTED]) |
| `0xb2a6a72843db0f508204a56448413f3867ea691a` | `BNB@BSC#133060467` | 24024.01911 | 0.737809937535090657 | shortfall / 不足 (CSV exceeds on-chain) by 24023.281300062464909343 (provider=https://bsc.rpc.sentio.xyz,cache,https://bsc-mainnet.public.blastapi.io) |
| `0xbf83d18a46325acb7d8f40a462d23a92f467ed7a` | `BNB@BSC#133060467` | 6341713.185 | 6341570.522025983181439649 | shortfall / 不足 (CSV exceeds on-chain) by 142.662974016818560351 (provider=https://bsc.rpc.sentio.xyz,cache,https://bsc-mainnet.public.blastapi.io) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `ETH@ETH#26395133` | 1598.971278 | 1598.966407095141905532 | shortfall / 不足 (CSV exceeds on-chain) by 0.004870904858094468 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `ETH@ETH#26395133` | 1139.223472 | 1139.203735935993209669 | shortfall / 不足 (CSV exceeds on-chain) by 0.019736064006790331 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xb772fd70ea293335e5095597083e7da79e22ada7` | `ETH@ETH#26395133` | 876.2637994 | 873.353154390770632927 | shortfall / 不足 (CSV exceeds on-chain) by 2.910645009229367073 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `ETH@ETH#26395133` | 995.8296926 | 995.595417919064782187 | shortfall / 不足 (CSV exceeds on-chain) by 0.234274680935217813 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `ETH@ETH#26395133` | 455.0292372 | 455.020462382068132893 | shortfall / 不足 (CSV exceeds on-chain) by 0.008774817931867107 (provider=https://mainnet.infura.io/v3/[REDACTED]) |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `ACH\|BSC` | 2 | no native verifier for this (coin,network) pair yet |
| `ACH\|ETH` | 15 | no native verifier for this (coin,network) pair yet |
| `AGLD\|ETH` | 11 | no native verifier for this (coin,network) pair yet |
| `AIGENSYN\|ETH` | 11 | no native verifier for this (coin,network) pair yet |
| `AIGENSYN\|GENSYN` | 2 | no native verifier for this (coin,network) pair yet |
| `ALICE\|BSC` | 8 | no native verifier for this (coin,network) pair yet |
| `ALICE\|ETH` | 17 | no native verifier for this (coin,network) pair yet |
| `ALLO\|ALLORA` | 6 | no native verifier for this (coin,network) pair yet |
| `ALLO\|BASE` | 2 | no native verifier for this (coin,network) pair yet |
| `ALLO\|BSC` | 3 | no native verifier for this (coin,network) pair yet |
| `ALLO\|ETH` | 12 | no native verifier for this (coin,network) pair yet |
| `AMP\|ETH` | 15 | no native verifier for this (coin,network) pair yet |
| `ANIME\|ARBITRUM` | 5 | no native verifier for this (coin,network) pair yet |
| `ANIME\|ETH` | 13 | no native verifier for this (coin,network) pair yet |
| `ANKR\|BSC` | 13 | no native verifier for this (coin,network) pair yet |
| `ANKR\|ETH` | 14 | no native verifier for this (coin,network) pair yet |
| `APE\|ARBITRUM` | 2 | no native verifier for this (coin,network) pair yet |
| `APE\|BSC` | 2 | no native verifier for this (coin,network) pair yet |
| `APE\|ETH` | 18 | no native verifier for this (coin,network) pair yet |
| `API3\|ETH` | 16 | no native verifier for this (coin,network) pair yet |
| `ARKM\|ETH` | 14 | no native verifier for this (coin,network) pair yet |
| `ARPA\|BSC` | 3 | no native verifier for this (coin,network) pair yet |
| `ARPA\|ETH` | 14 | no native verifier for this (coin,network) pair yet |
| `AUCTION\|BSC` | 6 | no native verifier for this (coin,network) pair yet |
| `AUCTION\|ETH` | 18 | no native verifier for this (coin,network) pair yet |

_269 additional unsupported pairs omitted; see verification bundle._

#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 101, WARN 20, UNVERIFIABLE 6, PASS 1556

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `AAVE@ETH#26395133` | 1319.364018 | 1319.30374333 | shortfall / 不足 (CSV exceeds on-chain) by 0.06027467 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x28e8d4eab37cfbf49edf7ff39c27377b7ae3a2a0` | `ASTER@BSC#133060467` | 30200045 | 35 | shortfall / 不足 (CSV exceeds on-chain) by 30200010 (provider=https://bsc.rpc.sentio.xyz) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BCH@ETH#26395133` | 1496 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1496 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `BCH@ETH#26395133` | 3.012 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.012 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0x564286362092d8e7936f0549571a803b203aaced` | `BCH@ETH#26395133` | 125 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 125 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#26395133` | 1093.5871 | 401.265 | shortfall / 不足 (CSV exceeds on-chain) by 692.3221 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#26395133` | 1047.9 | 106.362 | shortfall / 不足 (CSV exceeds on-chain) by 941.538 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4aefa39caeadd662ae31ab0ce7c8c2c9c0a013e8` | `BUSD@AVAXC#100064891` | 133923.604 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 133923.604 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9020475aa99c045f6cd9eb7ec048519ea022d4b3` | `BUSD@AVAXC#100064891` | 7660.28733 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 7660.28733 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0x9f8c163cba728e99993abe7495f06c0a3c8ac8b9` | `BUSD@AVAXC#100064891` | 3.00005477 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.00005477 (provider=https://api.avax.network/ext/bc/C/rpc) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BUSD@MATIC#97157428` | 156931.468 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 156931.468 (provider=https://polygon.gateway.tenderly.co) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `BUSD@MATIC#97157428` | 0.8 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 0.8 (provider=https://polygon.gateway.tenderly.co) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `BUSD@OPTIMISM#159421411` | 3.00061923 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.00061923 (provider=https://mainnet.optimism.io) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BUSD@OPTIMISM#159421411` | 4121.634 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 4121.634 (provider=https://mainnet.optimism.io) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `BUSD@MATIC#97157428` | 3.00088312 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.00088312 (provider=https://polygon.gateway.tenderly.co) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `CAKE@ETH#26395133` | 23954.30402 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 23954.30402 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `CAKE@ETH#26395133` | 24593.94636 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 24593.94636 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `CAKE@ETH#26395133` | 37726.69426 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 37726.69426 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `CAKE@ETH#26395133` | 26670.01947 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 26670.01947 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `CAKE@ETH#26395133` | 23474.00522 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 23474.00522 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `CAKE@ETH#26395133` | 314432.2739 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 314432.2739 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `CAKE@ETH#26395133` | 33366.0616 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 33366.0616 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `ENA@ETH#26395133` | 1314432.87 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1314432.87 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `ENA@ETH#26395133` | 2950272.549 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2950272.549 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `ENA@ETH#26395133` | 1821762.817 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1821762.817 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `ENA@ETH#26395133` | 27343446.22 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 27343446.22 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `ENA@ETH#26395133` | 7509397.046 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 7509397.046 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x5a1e5b56aa33c056ed7bb93b88b1774a456f6694` | `ENA@ETH#26395133` | 37081438.59 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 37081438.59 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `ENA@ETH#26395133` | 610706426.5 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 610706426.5 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `ENA@ETH#26395133` | 43300489.38 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 43300489.38 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `ENA@ETH#26395133` | 54484458.94 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 54484458.94 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `ENA@ETH#26395133` | 1355520.166 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1355520.166 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `ENA@ETH#26395133` | 2.05552835 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2.05552835 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `ENA@ETH#26395133` | 2.15510575 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2.15510575 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `ENA@ETH#26395133` | 1804739.474 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1804739.474 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `ENA@ETH#26395133` | 5792467.017 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 5792467.017 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `ENA@ETH#26395133` | 2076022.479 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2076022.479 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ENA@ETH#26395133` | 787759017 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 787759017 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0x5ce4825e2c5d7025b110d7bbeb5b1f03976e1443` | `GRT@ARBITRUM#523417793` | 82610815.47 | 100 | shortfall / 不足 (CSV exceeds on-chain) by 82610715.47 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `GRT@ETH#26395133` | 8136385.616 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 8136385.616 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x605d0b92cde1e110d5ae9123969097ea2489d1aa` | `GRT@ARBITRUM#523417793` | 6000113 | 113 | shortfall / 不足 (CSV exceeds on-chain) by 6000000 (provider=https://arbitrum.gateway.tenderly.co) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `GRT@ETH#26395133` | 2.50053131 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 2.50053131 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `GRT@ETH#26395133` | 1181769.439 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1181769.439 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4fdfe365436b5273a42f135c6a6244a20404271e` | `GRT@ETH#26395133` | 0.092126 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 0.092126 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `GRT@ETH#26395133` | 124497423.4 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 124497423.4 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `GRT@ETH#26395133` | 1315581.858 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1315581.858 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `GRT@ETH#26395133` | 4810843.501 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 4810843.501 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4aec0e98fc1fb55b9cc2faaa7a81acca42cb4e96` | `GRT@ETH#26395133` | 187827.7132 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 187827.7132 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `GRT@ETH#26395133` | 4645209.564 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 4645209.564 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `GRT@ETH#26395133` | 11091360.57 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 11091360.57 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@ETH#26395133` | 13473559.78 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 13473559.78 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `GRT@ETH#26395133` | 6682732.766 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 6682732.766 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `GRT@ETH#26395133` | 0.50038197 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 0.50038197 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `GRT@ETH#26395133` | 100000000 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 100000000 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x85b722aad2468e4794485df3bbbe1148aba070be` | `GRT@MATIC#97157428` | 4.18298127 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 4.18298127 (provider=https://polygon.gateway.tenderly.co) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `GRT@MATIC#97157428` | 73.1902 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 73.1902 (provider=https://polygon.gateway.tenderly.co) |
| `0xe69f81b825d7dc31ee9becef4dbeab5cf30e3abb` | `GRT@ETH#26395133` | 50000000 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 50000000 (provider=https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `GRT@ETH#26395133` | 6468171.921 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 6468171.921 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xe7804c37c13166ff0b37f5ae0bb07a3aebb6e245` | `GRT@MATIC#97157428` | 3.16925866 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.16925866 (provider=https://polygon.gateway.tenderly.co) |
| `0xc34db06fd02be471dc4c1d8b0316a63ee37903bd` | `OP@OPTIMISM#159421411` | 9176644.51 | 9176634.666759003382204843 | shortfall / 不足 (CSV exceeds on-chain) by 9.843240996617795157 (provider=https://mainnet.optimism.io) |
| `0xa64b436964e7415c0e70b9989a53e1fb9a90e726` | `POL@ETH#26395133` | 122000099.1 | 99.066047574468237635 | POL+MATIC ERC20 on ETH still << CSV; residual may be Polygon-native POL or internal ledger; shortfall / 不足 (CSV exceeds on-chain) by 122000000.033952425531762365 (provider=https://mainnet.infura.io/v3/[REDACTED],https://eth-mainnet.g.alchemy.com/v2/[REDACTED]) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `SHIB@ETH#26395133` | 35223032528 | 32223147670.66 | shortfall / 不足 (CSV exceeds on-chain) by 2999884857.34 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x26bff98f6f37c8d560684eaf2afd1149bf76633e` | `SSV@ETH#26395133` | 650019 | 14 | shortfall / 不足 (CSV exceeds on-chain) by 650005 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x18e226459ccf0eec276514a4fd3b226d8961e4d1` | `TUSD@BSC#133060467` | 15000 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 15000 (provider=https://bsc.rpc.sentio.xyz) |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `TUSD@BSC#133060467` | 256388.6563 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 256388.6563 (provider=https://bsc.rpc.sentio.xyz) |
| `0x52854a35cf896a57b3a861ba7fbd2171f2445ef2` | `TUSD@BSC#133060467` | 3731.581989 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3731.581989 (provider=https://bsc.rpc.sentio.xyz) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@BSC#133060467` | 396917.499 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 396917.499 (provider=https://bsc.rpc.sentio.xyz) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `TUSD@ETH#26395133` | 25584 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25584 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x21a31ee1afc51d94c2efccaa2092ad1028285549` | `TUSD@ETH#26395133` | 7262.2 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 7262.2 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x28c6c06298d514db089934071355e5743bf21d60` | `TUSD@ETH#26395133` | 3.00029591 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3.00029591 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `TUSD@ETH#26395133` | 25564.38 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25564.38 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `TUSD@ETH#26395133` | 1077.401833 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1077.401833 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@BSC#133060467` | 9.9151 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 9.9151 (provider=https://bsc.rpc.sentio.xyz) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@BSC#133060467` | 1127096.623 | 0.01899189 | shortfall / 不足 (CSV exceeds on-chain) by 1127096.60400811 (provider=https://bsc.rpc.sentio.xyz) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `TUSD@ETH#26395133` | 529912.638 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 529912.638 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x98adef6f2ac8572ec48965509d69a8dd5e8bba9d` | `TUSD@ETH#26395133` | 31703.33056 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 31703.33056 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `TUSD@ETH#26395133` | 25584 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25584 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `TUSD@ETH#26395133` | 25574.3 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25574.3 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xdfd5293d8e347dfe59e90efd55b2956a1343963d` | `TUSD@ETH#26395133` | 8693.109 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 8693.109 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `TUSD@ETH#26395133` | 25475.27 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 25475.27 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `TUSD@ETH#26395133` | 5757027.079 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 5757027.079 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `TUSD@ETH#26395133` | 172019.41 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 172019.41 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x5a1e5b56aa33c056ed7bb93b88b1774a456f6694` | `TUSD@ETH#26395133` | 46838.56089 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 46838.56089 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x59fbc084aeed2c8e6e9b6ebe8a11d443efcb40d2` | `USD1@AB#16577265` | 16383889.29 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 16383889.29 (provider=https://api.mainnet.abs.xyz) |
| `0xce43851485f437005ef3c1821e250c4d3c53b946` | `USDC@ETH#26395133` | 80760147.87 | 80752392.679983 | shortfall / 不足 (CSV exceeds on-chain) by 7755.190017 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x282e0bd2d6314ac0ec139ab3c40bb368df699491` | `USDC@LINEA#32579493` | 0.460905 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 0.460905 (provider=https://rpc.linea.build) |
| `0x037741ef196087b26abfb173b64d05df23fadc72` | `USDT@ETH#26395133` | 4859522.123 | 4845654.469506 | shortfall / 不足 (CSV exceeds on-chain) by 13867.653494 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x6781b0098e3238caeede00647c9b69dd98a86a83` | `USDT@ETH#26395133` | 4635697.419 | 4622023.119043 | shortfall / 不足 (CSV exceeds on-chain) by 13674.299957 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `USDT@ETH#26395133` | 82153126.29 | 82150800.32817 | shortfall / 不足 (CSV exceeds on-chain) by 2325.96183 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xb772fd70ea293335e5095597083e7da79e22ada7` | `USDT@ETH#26395133` | 5618786.606 | 5618037.205518 | shortfall / 不足 (CSV exceeds on-chain) by 749.400482 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xc6fe7b7fbc7ef114284896a6631251cf54a8c304` | `USDT@ETH#26395133` | 6231106.369 | 6092506.368724 | shortfall / 不足 (CSV exceeds on-chain) by 138600.000276 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xde1c31ef18a36e44125c641063544ed5d44bfe78` | `USDT@ETH#26395133` | 5282784.157 | 5281912.733556 | shortfall / 不足 (CSV exceeds on-chain) by 871.423444 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x3d5dc8d8ef62348e8e2e1abe4aa6d1153b48f232` | `USDT@ETH#26395133` | 3299444.179 | 3299319.179005 | shortfall / 不足 (CSV exceeds on-chain) by 124.999995 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0xf3e9fdd1ae54e5ad88ee540855f83cf3918f4385` | `USDT@ETH#26395133` | 3144276.146 | 3142777.446279 | shortfall / 不足 (CSV exceeds on-chain) by 1498.699721 (provider=https://mainnet.infura.io/v3/[REDACTED]) |
| `0x85b722aad2468e4794485df3bbbe1148aba070be` | `USDT@MATIC#97157428` | 3890307.433 | 3889809.393345 | shortfall / 不足 (CSV exceeds on-chain) by 498.039655 (provider=https://polygon.gateway.tenderly.co) |
| `0x54df11cbb1e7c58755208a5e10f2a2e1b949c1dd` | `USDT@OPTIMISM#159421411` | 12 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 12 (provider=https://mainnet.optimism.io) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `USDT@OPTIMISM#159421411` | 713.507801 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 713.507801 (provider=https://mainnet.optimism.io) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `USDT@OPTIMISM#159421411` | 3737941.574 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 3737941.574 (provider=https://mainnet.optimism.io) |
| `0xc34db06fd02be471dc4c1d8b0316a63ee37903bd` | `USDT@OPTIMISM#159421411` | 1421898.506 | 0 | shortfall / 不足 (CSV exceeds on-chain) by 1421898.506 (provider=https://mainnet.optimism.io) |
| `0xc8daf809c7d7c27dd62d006196d8901ba57e5eae` | `USDT@PLASMA#37498220` | 6500001 | 1.00001 | shortfall / 不足 (CSV exceeds on-chain) by 6499999.99999 (provider=https://rpc.plasma.to) |
| `0x690c3eebf6ee237f00fc96f9766282f096c191cc` | `XRP@ETH#26395133` | 49189.01247 | 49179.502465165160096938 | shortfall / 不足 (CSV exceeds on-chain) by 9.510004834839903062 (provider=https://mainnet.infura.io/v3/[REDACTED]) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x02e48b9e60ea882d773ffb2551f8f07733492e3950ba9fa09cfd0c36c127f26b` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x0213c67ed78bc280887234fe5ed5e77272465317978ae86c25a71531d9332a2d` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x04164013f90b05d67f026779bf96e9c401c96f3485b645a786166e6935fba116` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x0785ba32c2d07462534d523b48f908ff4a6360cb702749e83566239341ed6aae` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `TJTKCzQpsG19CSMYF9TPtxz6wSHPSyeKfa` | `TUSD@TRX#87922590` | 546047.2153 | 545954.755327560768729417 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 92.459972439231270583 (provider=https://api.trongrid.io) |
| `TCQ8X6x9GiK4khydWVKfLJvm8H6WK3RE8B` | `U@TRX#87922590` | 1250000 | 1249999 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 1 (provider=https://tron-rpc.publicnode.com) |
| `0x9dfb9014e88087fba78cc9309c64031d02be9a33` | `USDC@ARBITRUM` |  |  | rpc error (provider=): https://arbitrum.gateway.tenderly.co -> HTTP 429 |
| `0x0213c67ed78bc280887234fe5ed5e77272465317978ae86c25a71531d9332a2d` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x04164013f90b05d67f026779bf96e9c401c96f3485b645a786166e6935fba116` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x0785ba32c2d07462534d523b48f908ff4a6360cb702749e83566239341ed6aae` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x02e48b9e60ea882d773ffb2551f8f07733492e3950ba9fa09cfd0c36c127f26b` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x480e8eb029fb996bfe4569bd7bfdd398f8c38dc0f6838fb23a2ceecf2d3754c` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `TADfh7Q4KA9LD6WfjjjcPMNra33gNRrCkb` | `USDT@TRX#87922590` | 4744693.082 | 3466688.101669 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 1278004.980331 (provider=https://api.trongrid.io) |
| `TCQ8X6x9GiK4khydWVKfLJvm8H6WK3RE8B` | `USDT@TRX#87922590` | 2274996.772 | 1882655.174978 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 392341.597022 (provider=https://tron-rpc.publicnode.com) |
| `TE4mtETxcuHHbprQmQkyVSaCC48Q5yk21i` | `USDT@TRX#87922590` | 3384400.465 | 2807411.707022 | live TRC20 balance vs snapshot CSV; Tron has no historical contract state (block_num ineffective); shortfall / 不足 (CSV exceeds on-chain) by 576988.757978 (provider=https://tron-rpc.publicnode.com) |


_5 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `DOT\|ETH` | 1 | no native verifier for this (coin,network) pair yet |
| `USDC\|ARC` | 6 | no native verifier for this (coin,network) pair yet |
| `USDT\|EOS` | 3 | no native verifier for this (coin,network) pair yet |
| `U\|ROBINHOOD` | 2 | no native verifier for this (coin,network) pair yet |
Other unverifiable rows:

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `1wzfee3q4xyudur6xgmzwwjbs3xasm91o` | `USDC@ETH` |  |  | invalid EVM wallet address for USDC@ETH |
| `0x564286362092d8e7936f0549571a803b203aaced` | `BTC@ETH#26395133` | 0.24702247 | 0 | no WBTC at snapshot height; BTC likely native/off-chain custody — EVM WBTC balanceOf cannot verify this row |

#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 2, WARN 450, UNVERIFIABLE 6, PASS 135

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `bc1q9n7j6a5wwru6tad9xrlncnagresvyjt96pl3en` | `BTC\|BTC#975493` | 1300.019572 | 100.01957174 | shortfall / 不足 (CSV exceeds on-chain) by 1200.00000026 (provider=https://blockstream.info/api) |
| `LNhQm4699LqWSUQ4qNLNsr9815tJK2y2GW` | `LTC\|LTC#3211492` | 103686.6161 | 103679.98494306 | shortfall / 不足 (CSV exceeds on-chain) by 6.63115694 (provider=https://litecoinspace.org/api) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x265ca6100138fd7274fa66f043b7b259c44cdc64f75ffd634a4fb523d9d47d8c/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x1d14ee0c332546658b13965a39faf5ec24ad195b722435d9fe23dc55487e67e3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x292f853b48a28864755c971299ce8a73a3e32c19a0f7b8dbbf782482396e8ef3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x33f91e694d40ca0a14cb84e1f27a4d03de5cf292b07ed75ed3286e4f243dab34/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x61e28472a24bc69824a036bfafb20fc0077e61a9265a335d2db2f52e3b2108f2` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x61e28472a24bc69824a036bfafb20fc0077e61a9265a335d2db2f52e3b2108f2/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x6d7ca7650519c77afd767da9b8f8abd4234de0bb63a0047386d939bf9344e5c2` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x6d7ca7650519c77afd767da9b8f8abd4234de0bb63a0047386d939bf9344e5c2/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x716666e019eb2cd1eea5ae29760e064f14984d8d6db2ff9ee56d0bd994e8c9b3/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x9577885192fe855857dfd699ba15ef89e7eac15f770f1d89565d7188339a7fbb` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9577885192fe855857dfd699ba15ef89e7eac15f770f1d89565d7188339a7fbb/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x9bb0eda169c10d09c042733125bd43a5ddc16fde27af994b884fd63445ac945d` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9bb0eda169c10d09c042733125bd43a5ddc16fde27af994b884fd63445ac945d/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x9a77147bad715dcfa5899c591e37310a0e587b05d4ad96d59db310c0d1cf36e5` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9a77147bad715dcfa5899c591e37310a0e587b05d4ad96d59db310c0d1cf36e5/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0xaa367ea8e145e3af459d3b8f6d50b96c76d9cf9c7dabedf9fff192f36eda0ee0` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xaa367ea8e145e3af459d3b8f6d50b96c76d9cf9c7dabedf9fff192f36eda0ee0/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x9fbc354d59041b8b1b8368e3e7397ac943a3c7c6da3ffde3aa4f4d221a1d205d/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `APT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/resource/0x1::coin::CoinStore%3C0x1::aptos_coin::AptosCoin%3E?ledger_version=1177771868): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1177771868) has been pruned","hint":"The requested data has bee |


_435 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `DOT\|ETH` | 1 | no native verifier for this (coin,network) pair yet |
| `TRX\|BSC` | 6 | no native verifier for this (coin,network) pair yet |
| `TRX\|ETH` | 10 | no native verifier for this (coin,network) pair yet |
| `TRX\|TRX` | 46 | no native verifier for this (coin,network) pair yet |
| `USDC\|ARC` | 6 | no native verifier for this (coin,network) pair yet |
| `USDT\|EOS` | 3 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-deposit` (`WARN`)

Summary: more than 25% of sampled deposit rows failed RPC

Finding counts: WARN 500, PASS 2

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0.0.1873771` | `HBAR\|HBAR` |  |  | deposit rpc error (provider=https://mainnet-public.mirrornode.hedera.com/api/v1/blocks/102227827): Get "https://mainnet-public.mirrornode.hedera.com/api/v1/blocks/102227827": context deadline exceeded |
| `0x064c8e2091103256d27e32a9a1a53b418d9870d5d82cff6a2a4a2f3db422a180` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x064c8e2091103256d27e32a9a1a53b418d9870d5d82cff6a2a4a2f3db422a180/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x064c8e2091103256d27e32a9a1a53b418d9870d5d82cff6a2a4a2f3db422a180/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0x744bf58c5562a9ebc71ef07743e8e3385eec8e5743fde745d937ba22d37e9eb1` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x744bf58c5562a9ebc71ef07743e8e3385eec8e5743fde745d937ba22d37e9eb1/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x744bf58c5562a9ebc71ef07743e8e3385eec8e5743fde745d937ba22d37e9eb1/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0xdb7b78ba6ac802ce5588b7de5b24e6f7827ac452c2e81b778a5ea1336ad555c8` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xdb7b78ba6ac802ce5588b7de5b24e6f7827ac452c2e81b778a5ea1336ad555c8/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xdb7b78ba6ac802ce5588b7de5b24e6f7827ac452c2e81b778a5ea1336ad555c8/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0x37a24b59cdee30ceb4394d530a873739c0888d3554e3e4d3db0423896afc9b02` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x37a24b59cdee30ceb4394d530a873739c0888d3554e3e4d3db0423896afc9b02/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x37a24b59cdee30ceb4394d530a873739c0888d3554e3e4d3db0423896afc9b02/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0x6dffa1ec1ba337320484a46fe9886d5327e67fa6c922b390d69da7f19f964ec7` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x6dffa1ec1ba337320484a46fe9886d5327e67fa6c922b390d69da7f19f964ec7/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x6dffa1ec1ba337320484a46fe9886d5327e67fa6c922b390d69da7f19f964ec7/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0x127cab605fd38774b1234a8ec9ef412ca292f183cffaae0fbd2d273ec5c77aa7` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x127cab605fd38774b1234a8ec9ef412ca292f183cffaae0fbd2d273ec5c77aa7/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x127cab605fd38774b1234a8ec9ef412ca292f183cffaae0fbd2d273ec5c77aa7/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0x65499f2565f95cc9c2007252e3957220f7c1beeaa25899ffda78665b780ca864` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x65499f2565f95cc9c2007252e3957220f7c1beeaa25899ffda78665b780ca864/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x65499f2565f95cc9c2007252e3957220f7c1beeaa25899ffda78665b780ca864/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0x850f082c906b6d09b4f3497375f803c30340898cc9dc8683ffd2434e1500bbd3` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x850f082c906b6d09b4f3497375f803c30340898cc9dc8683ffd2434e1500bbd3/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x850f082c906b6d09b4f3497375f803c30340898cc9dc8683ffd2434e1500bbd3/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0xf21a04ef463acfcdb9b460301b4e5b689d3a86b5a14d6231b0df93d60c01f4cd` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xf21a04ef463acfcdb9b460301b4e5b689d3a86b5a14d6231b0df93d60c01f4cd/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xf21a04ef463acfcdb9b460301b4e5b689d3a86b5a14d6231b0df93d60c01f4cd/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0xf8f89ae516a6384708cc6d3129c7c2a71494f7bcdb9d99838afa725678f1b368` | `USDC\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xf8f89ae516a6384708cc6d3129c7c2a71494f7bcdb9d99838afa725678f1b368/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868): Get "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xf8f89ae516a6384708cc6d3129c7c2a71494f7bcdb9d99838afa725678f1b368/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1177771868": context deadline exceeded |
| `0x13f12f5d8b269197be17f0ad4f0a5b27463c7839` | `USDC\|SONIC` |  |  | deposit rpc error (provider=): Post "https://rpc.soniclabs.com": context deadline exceeded |
| `0xd224e961f33ce0b5e59bf1e2235d0a11bcd1ab28` | `ARB\|ETH` |  |  | deposit rpc error (provider=): Post "https://eth.drpc.org": context deadline exceeded |
| `0xd224e961f33ce0b5e59bf1e2235d0a11bcd1ab28` | `CRV\|ETH` |  |  | deposit rpc error (provider=): no usable provider |
| `0xa9c61fe59b5702b1d382fd1d5e495887ff34c21d` | `CRV\|ETH` |  |  | deposit rpc error (provider=): no usable provider |


_485 additional `WARN` rows omitted; see verification bundle._

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR12NOV26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-11-12T00:00:00Z |
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
