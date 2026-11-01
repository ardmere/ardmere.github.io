# BINANCE PR31OCT26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR31OCT26-assessment.json](./PR31OCT26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR31OCT26` |
| Snapshot time | `2026-10-31T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR31OCT26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

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
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `~monthly` |
| History available | `7 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `a89b66d7efdeddbc2ee58eecec1f489a1c68d208e56cf1eef9f3680d605e498b` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `a8c216bce529919a137dce8d52936471c59b14b5ec12377b9a5fb7d3758b0251` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20261031.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xb7b67ea4153601ff05322208000f6003a56c5e72c188f18b03d543ac39e3100c` |
| Verification bundle root | `0x8f008fd012b050ab7e03e9ab4610e7aa3a27b5dff234adb6cac05bfbf57e0466` |
| Artifact bundle SHA-256 | `b398554b320660c53fad2465f7f5a1292dffacf8ca1176f7f1ca351baed89a48` |
| Verification bundle SHA-256 | `07f3f41fa61c99c4ebbc2aa363bd9a7aade6dac9f229e7d62e67613e0fa7ee3e` |

Local bundle paths: [PR31OCT26.artifact-bundle.json](../../../artifacts/binance/PR31OCT26/bundles/PR31OCT26.artifact-bundle.json), [PR31OCT26.verification-bundle.v2.json](../../../artifacts/binance/PR31OCT26/bundles/PR31OCT26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `onchain-balance-hot` | `2.1` | `WARN` | `0.0940` | some rows have incomplete staking/deposit attribution below CSV — not reserve shortfall |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.6124` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.2366` |  |
| `onchain-balance-deposit` | `1.2` | `PASS` | `0.4254` |  |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR31OCT26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `onchain-balance-hot` (`WARN`)

Summary: some rows have incomplete staking/deposit attribution below CSV — not reserve shortfall

Finding counts: WARN 1, UNVERIFIABLE 1, PASS 1

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xb2a6a72843db0f508204a56448413f3867ea691a` | `BNB@BSC#130757300` | 35011.51204 | 0.737837337885090657 | chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve; observed delta=35010.774202662114909343 (provider=https://bsc.rpc.sentio.xyz,cache,https://bsc-mainnet.public.blastapi.io) |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `USDC\|ARC` | 1 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-token` (`FAIL`)

Finding counts: FAIL 5, WARN 11, UNVERIFIABLE 1, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `BTC@ETH#26309019` | 1047.9 | 106.362 | on-chain balance != csv claim by 941.538 (provider=https://eth-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}) |
| `0x5a52e96bacdabb82fd05763e25335261b270efcb` | `BTC@ETH#26309019` | 1106.4841 | 400 | on-chain balance != csv claim by 706.4841 (provider=https://eth-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}) |
| `0x59fbc084aeed2c8e6e9b6ebe8a11d443efcb40d2` | `USD1@AB#16231665` | 16383889.29 | 0 | on-chain balance != csv claim by 16383889.29 (provider=https://api.mainnet.abs.xyz) |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `USDT@OPTIMISM#158903011` | 21176713 | 0 | on-chain balance != csv claim by 21176713 (provider=https://mainnet.optimism.io) |
| `0xacd03d601e5bb1b275bb94076ff46ed9d753435a` | `USDT@OPTIMISM#158903011` | 3.587801 | 0 | on-chain balance != csv claim by 3.587801 (provider=https://mainnet.optimism.io) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x00c1ce9c68f1c4a1b05e8d043b9b05784f41685d68e5009ffe379aa7f7c95ccf` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x04164013f90b05d67f026779bf96e9c401c96f3485b645a786166e6935fba116` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0x0213c67ed78bc280887234fe5ed5e77272465317978ae86c25a71531d9332a2d` | `ETH@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `ETH@ZKSYNCERA` |  |  | rpc error (provider=): Post "https://mainnet.era.zksync.io": context deadline exceeded (Client.Timeout exceeded while awaiting headers) |
| `0x0213c67ed78bc280887234fe5ed5e77272465317978ae86c25a71531d9332a2d` | `USDC@STARKNET` |  |  | rpc error (provider=): HTTP 410: {"error":"This endpoint has been discontinued.","message":"For an up to date list of available endpoints, visit https:// |
| `0xf977814e90da44bfa03b6295a0616a897441acec` | `USDC@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0xa84fd90d8640fa63d194601e0b2d1c9094297083` | `USDC@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x4aec0e98fc1fb55b9cc2faaa7a81acca42cb4e96` | `USDC@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x4ed6cf63bd9c009d247ee51224fc1c7041f517f1` | `USDC@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x7aed074ca56f5050d5a2e512ecc5bf7103937d76` | `USDC@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |
| `0x43684d03d81d3a4c70da68febdd61029d426f042` | `USDC@ZKSYNCERA` |  |  | rpc error (provider=): no usable provider |

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `USDC\|ARC` | 1 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-ledger` (`FAIL`)

Finding counts: FAIL 6, WARN 21, UNVERIFIABLE 1, PASS 1

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `bc1q9n7j6a5wwru6tad9xrlncnagresvyjt96pl3en` | `BTC\|BTC#973895` | 1300.019572 | 100.01957174 | on-chain balance != csv claim by 1200.00000026 (provider=https://blockstream.info/api) |
| `c5f9zfpkKMD9N8uLqJcFeJAAz7v12vDMnup9Y6EeQkk` | `SOL\|SOL#463478911` | 4211389.625 | 1.064970212 | on-chain balance != csv claim by 4211388.560029788 (provider=https://api.helius.xyz/v1/wallet) |
| `9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM` | `USDC\|SOL#463478911` | 1108.086785 | 1107.086785 | on-chain balance != csv claim by 1 (provider=https://api.helius.xyz/v1/wallet) |
| `5SDrsMNTYdhmApjfqYHDvjoW92f2S42vcc7zNDVcQ9Ej` | `USDC\|SOL#463478911` | 3.53509 | 2.000936 | on-chain balance != csv claim by 1.534154 (provider=https://api.helius.xyz/v1/wallet) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDC\|SOL#463478911` | 1984800.492 | 3.000307 | on-chain balance != csv claim by 1984797.491693 (provider=https://api.helius.xyz/v1/wallet) |
| `5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9` | `USDT\|SOL#463478911` | 19114.68138 | 3.000038 | on-chain balance != csv claim by 19111.681342 (provider=https://api.helius.xyz/v1/wallet) |

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0xbe13ea18241ed83ada91e8f5ca37b5507b397cb66ae1e52151cf32a96547600b` | `USDC\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xbe13ea18241ed83ada91e8f5ca37b5507b397cb66ae1e52151cf32a96547600b/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `USDC\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f` | `USDC\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x80174e0fe8cb2d32b038c6c888dd95c3e1560736f0d4a6e8bed6ae43b5c91f6f/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `USDC\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00` | `USDC\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x5bd7de5c56d5691f32ea86c973c73fec7b1445e59736c97158020018c080bb00/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xaa367ea8e145e3af459d3b8f6d50b96c76d9cf9c7dabedf9fff192f36eda0ee0` | `USDC\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xaa367ea8e145e3af459d3b8f6d50b96c76d9cf9c7dabedf9fff192f36eda0ee0/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a` | `USDC\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a/balance/0xbae207659db88bea0cbead6da0ed00aac12edcdda169e591cd41c94180b46f3b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x7ab9a6a7109dcb9cb357a109f32dfcc78a7aa2d6029084eb924d95133fc71cec` | `USDC\|SUI` |  |  | ledger rpc error (provider=https://fullnode.mainnet.sui.io:443): sui: Method not found. JSON-RPC on public fullnodes has been deprecated. Please migrate to gRPC or GraphQL endpoints. See https://docs.sui.io/develop/accessing-data/json-rpc-migration for more information. |
| `0x935029ca5219502a47ac9b69f556ccf6e2198b5e7815cf50f68846f723739cbd` | `USDC\|SUI` |  |  | ledger rpc error (provider=https://fullnode.mainnet.sui.io:443): sui: Method not found. JSON-RPC on public fullnodes has been deprecated. Please migrate to gRPC or GraphQL endpoints. See https://docs.sui.io/develop/accessing-data/json-rpc-migration for more information. |
| `0xac5bceec1b789ff840d7d4e6ce4ce61c90d190a7f8c4f4ddf0bff6ee2413c33c` | `USDC\|SUI` |  |  | ledger rpc error (provider=https://fullnode.mainnet.sui.io:443): sui: Method not found. JSON-RPC on public fullnodes has been deprecated. Please migrate to gRPC or GraphQL endpoints. See https://docs.sui.io/develop/accessing-data/json-rpc-migration for more information. |
| `0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70` | `USDT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xae1a6f3d3daccaf77b55044cea133379934bba04a11b9d0bbd643eae5e6e9c70/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e` | `USDT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xed8c46bec9dbc2b23c60568f822b95b87ea395f7e3fdb5e3adc0a30c55c0a60e/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a` | `USDT\|APT` |  |  | ledger rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xd91c64b777e51395c6ea9dec562ed79a4afa0cd6dad5a87b187c37198a1f855a/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h` | `BTC\|BTC` |  |  | ledger rpc error (provider=https://blockstream.info/api): esplora tx pagination exceeded 400 pages for bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h |
| `t1M1j1rPwUFGz8GDFrgi8rs9R3c7XtCH4Gv` | `ZEC\|ZEC` |  |  | ledger rpc error (provider=https://cipherscan.app/api/address/t1M1j1rPwUFGz8GDFrgi8rs9R3c7XtCH4Gv): cipherscan HTTP 429 |


_6 additional `WARN` rows omitted; see verification bundle._

**UNVERIFIABLE (row-level)**

Unsupported (coin, network) pairs:

| Pair | Rows | Note |
| --- | --- | --- |
| `USDC\|ARC` | 1 | no native verifier for this (coin,network) pair yet |
#### `onchain-balance-deposit` (`PASS`)

Finding counts: WARN 11, PASS 2

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x4fbff26f917e6d11dbbee980a8f109c2f5b225df1819a84c93d86c4dd36ba8bd` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x4fbff26f917e6d11dbbee980a8f109c2f5b225df1819a84c93d86c4dd36ba8bd/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xfa15c476424db44f1ce6094740b960dfafb8336d430d2a58ec4ed6a46648bd3e` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xfa15c476424db44f1ce6094740b960dfafb8336d430d2a58ec4ed6a46648bd3e/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x0feb3f48af880494b4c808225992e865644eb75d1af19928888845a651c16502` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x0feb3f48af880494b4c808225992e865644eb75d1af19928888845a651c16502/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x92c30e02e202453dd5babfbc1ba587be09a08937d06a2d0d6a2a30d148d3573d` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x92c30e02e202453dd5babfbc1ba587be09a08937d06a2d0d6a2a30d148d3573d/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x6cd96e07a157d9cb36c2a20aec257be285e87cc23b3d806f47c27d4be40f9d90` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x6cd96e07a157d9cb36c2a20aec257be285e87cc23b3d806f47c27d4be40f9d90/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xf428ce24216a3db0aef70a0ce36627a1affcf128fe7e48100924a3f3287ee838` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xf428ce24216a3db0aef70a0ce36627a1affcf128fe7e48100924a3f3287ee838/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xc7f299a8705fd0cd04986a58b8beb0d2418d318f106767e539cdd17970731ea3` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xc7f299a8705fd0cd04986a58b8beb0d2418d318f106767e539cdd17970731ea3/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x13d8aafa64fe3c8748fb4e5e9458579eae5ca0386bdcf2c3a7e84263c1689107` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x13d8aafa64fe3c8748fb4e5e9458579eae5ca0386bdcf2c3a7e84263c1689107/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0xb819e0e24b4fb958cd18314740dce3278653a7f9f9d4d8ff23c744eefade2270` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0xb819e0e24b4fb958cd18314740dce3278653a7f9f9d4d8ff23c744eefade2270/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x88631990bbefd456a9e80fd82b04386debf66a11c15f4b6be1386cfe312fecbe` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x88631990bbefd456a9e80fd82b04386debf66a11c15f4b6be1386cfe312fecbe/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |
| `0x06863a01d5aac6ee9b153ffea8393022d4057b8d4fdf582d073ac5ecb2e0f93a` | `USDT\|APT` |  |  | deposit rpc error (provider=https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x06863a01d5aac6ee9b153ffea8393022d4057b8d4fdf582d073ac5ecb2e0f93a/balance/0x357b0b74bc833e95a115ad22604854d6b0fca151cecd94111770e5d6ffc9dc2b?ledger_version=1150394966): HTTP 410: {"error_code":"version_pruned","message":"Ledger version(1150394966) has been pruned","hint":"The requested data has bee |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR31OCT26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-10-31T00:00:00Z |
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
