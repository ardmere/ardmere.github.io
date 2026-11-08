# Deployments — ArdmerePoRAnchor

On-chain anchor contract deployments and first anchor transactions.

**Product target (ADR-002):** Base mainnet (`8453`) is the intended **production** primary anchor chain.  
**Operable today:** Ethereum Sepolia only. Base mainnet is **not deployed** in this refresh (requires deployer key + Base RPC — see checklist below).

| Network | Chain ID | Role | Status |
| --- | --- | --- | --- |
| Ethereum Sepolia | `11155111` | Active testnet / CLI default | **Live** (proxy + impl v4) |  # pragma: allowlist secret
| Base | `8453` | Production primary (ADR-002) | **Pending** — not deployed |
| Base Sepolia | `84532` | Optional L2 testnet | **Pending** — not deployed |

See [`decisions.md` ADR-002](./decisions.md#adr-002-primary-anchor-chain-base).

---

## Ethereum Sepolia (testnet, chain id `11155111`)  # pragma: allowlist secret

Primary **testnet** deployment (operable CLI default until Base mainnet ships). Contract is **UUPS upgradeable** — users interact with the **proxy** address; upgrades are authorized by the deployer (owner).

Deployment using `PRIVATE_KEY` / `ANCHOR_SIGNER` from `~/.zshenv` (testnet only).

| Item | Value |
|---|---|
| **Proxy (ANCHOR_CONTRACT)** | [`0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7`](https://sepolia.etherscan.io/address/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7) |  # pragma: allowlist secret
| **Implementation (v4, stage + dual coverage)** | [`0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78`](https://sepolia.etherscan.io/address/0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78) |  # pragma: allowlist secret
| **Implementation (v3, superseded)** | [`0x62050405283f222EFEdF3BF0d6Cb16541cd1327a`](https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a) |
| **Implementation (v2, superseded)** | [`0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392`](https://sepolia.etherscan.io/address/0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392) |
| **Anchor signer** | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |
| **Upgrade owner** | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |
| **Schema version** | `4` (stage + dual coverage; v3 rows remain readable) |
| **Storage version** | `2` (v1/v3 records still readable; new fields zero/false) |
| **Pattern** | UUPS (`ERC1967Proxy` + `ArdmerePoRAnchor` implementation) |
| **Impl deploy tx (v4)** | [`0xb746e715c68a5fd3a5e603b4f9f48e4f0158a45fc4f5225420c2b94a6c12a0f9`](https://sepolia.etherscan.io/tx/0xb746e715c68a5fd3a5e603b4f9f48e4f0158a45fc4f5225420c2b94a6c12a0f9) |
| **Upgrade tx (v3→v4)** | [`0xf4bdf0a43969ff31f344389d40bf7bd8032fd7bfc4f018843a9cb388032cf3b6`](https://sepolia.etherscan.io/tx/0xf4bdf0a43969ff31f344389d40bf7bd8032fd7bfc4f018843a9cb388032cf3b6) |
| **Impl deploy tx (v3)** | [`0x6aaa54524151d18be495f26fcf60b4f4ed155719b891b2ebeaf36bab9375b96f`](https://sepolia.etherscan.io/tx/0x6aaa54524151d18be495f26fcf60b4f4ed155719b891b2ebeaf36bab9375b96f) |
| **Upgrade tx (v2→v3)** | [`0x905a68386dc49096208e4f75a060aa2dddb55f1b36c82735f8de73f471e75699`](https://sepolia.etherscan.io/tx/0x905a68386dc49096208e4f75a060aa2dddb55f1b36c82735f8de73f471e75699) |
| **Impl deploy tx (v2)** | [`0xbf4684b13b133370e71ae946b2972e233d6156aa5c208c2d11b7b3778da3ca39`](https://sepolia.etherscan.io/tx/0xbf4684b13b133370e71ae946b2972e233d6156aa5c208c2d11b7b3778da3ca39) |
| **Proxy deploy tx** | [`0xbd3b1a874831f37768d08a28bb015bbd39261b3bf32ecad99df652a1cb772e7d`](https://sepolia.etherscan.io/tx/0xbd3b1a874831f37768d08a28bb015bbd39261b3bf32ecad99df652a1cb772e7d) |
| **Deploy block** | `11182041` |
| **Source verified (Etherscan)** | [Proxy](https://sepolia.etherscan.io/address/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7#code) ✓ · [Implementation v4](https://sepolia.etherscan.io/address/0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78#code) · [Implementation v3](https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a#code) ✓ · [Implementation v2 (superseded)](https://sepolia.etherscan.io/address/0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392#code) ✓ |  # pragma: allowlist secret
| **Source verified (Sourcify)** | [Proxy](https://repo.sourcify.dev/contracts/full_match/11155111/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7/) · [Implementation v4](https://repo.sourcify.dev/contracts/full_match/11155111/0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78/) · [Implementation v3](https://repo.sourcify.dev/contracts/full_match/11155111/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a/) · [Implementation v2 (superseded)](https://repo.sourcify.dev/contracts/full_match/11155111/0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392/) |  # pragma: allowlist secret

### Previous non-upgradeable deployment (superseded)

| Item | Value |
|---|---|
| **Contract** | [`0x6106d883F222DCfA2B797c2231984541BD878873`](https://sepolia.etherscan.io/address/0x6106d883F222DCfA2B797c2231984541BD878873) |
| **Deploy tx** | [`0x54fb600a67a85bf8d8639130e96eb01e5b4d0801e508dcd72f98c7d69b34330c`](https://sepolia.etherscan.io/tx/0x54fb600a67a85bf8d8639130e96eb01e5b4d0801e508dcd72f98c7d69b34330c) |
| **Deploy block** | `11179549` |
| **Source verified** | [Etherscan](https://sepolia.etherscan.io/address/0x6106d883F222DCfA2B797c2231984541BD878873#code) ✓ |

### Homepage assessment roots (`origin/main` @ `c652dd2`, 2026-10-31)

Anchored the six snapshots shown on the public homepage into Sepolia `ArdmerePoRAnchor` (schema v3). Roots taken from `docs/reports/<exchange>/<auditId>-assessment.json` artifact / verification bundle fields. Base mainnet was **not** attempted (`BASE_RPC_URL` missing; contract undeployed).

| Exchange | auditId | periodSeq | Block | `anchorSnapshot` tx |
| --- | --- | --- | --- | --- |
| bingx | `20260915` | 0 | `12034438` | [`0x9e5391322f199afb77e9b867a793ac5c85c70f9b91bd8b0ecbd94f544a8afabc`](https://sepolia.etherscan.io/tx/0x9e5391322f199afb77e9b867a793ac5c85c70f9b91bd8b0ecbd94f544a8afabc) |
| kucoin | `20260930` | 0 | `12034439` | [`0x694595c8bb857d2afb64d74c7f6ed0ec9f96606c30551c86832e14469192804c`](https://sepolia.etherscan.io/tx/0x694595c8bb857d2afb64d74c7f6ed0ec9f96606c30551c86832e14469192804c) |
| binance | `PR01OCT26` | 47 | `12034444` | [`0xb99fd5b688d39bfda5b0c23923b4ebd2d54319d6f76cace25655c2c7509ff233`](https://sepolia.etherscan.io/tx/0xb99fd5b688d39bfda5b0c23923b4ebd2d54319d6f76cace25655c2c7509ff233) |
| okx | `502260945` | 3 | `12034445` | [`0x27a239c3f19730718cfe54dbd4d65a93223729a0b7bb5962fccf2d35ac3e9722`](https://sepolia.etherscan.io/tx/0x27a239c3f19730718cfe54dbd4d65a93223729a0b7bb5962fccf2d35ac3e9722) |
| bingx | `20261015` | 1 | `12034446` | [`0xe75c010b8108f90e28232f734553426e3efe7d565bcb971d8b9a7d247267c3b2`](https://sepolia.etherscan.io/tx/0xe75c010b8108f90e28232f734553426e3efe7d565bcb971d8b9a7d247267c3b2) |
| mexc | `202610` | 0 | `12034449` | [`0xffc454a681132218060283fa9ef8feff43d8a8cf18eefe08078561000fe24930`](https://sepolia.etherscan.io/tx/0xffc454a681132218060283fa9ef8feff43d8a8cf18eefe08078561000fe24930) |

Post-anchor `getSnapshotCount` (Sepolia): binance `4`, okx `4`, bingx `2`, kucoin `1`, mexc `1`.

### Homepage assessment roots (`origin/main` @ `f11109d`, 2026-11-02)

Anchored the six snapshots shown on the public homepage into Sepolia `ArdmerePoRAnchor` (schema v3). Roots taken from `docs/reports/<exchange>/<auditId>-assessment.json` artifact / verification bundle fields. Includes Binance `PR31OCT26` after the refreshed exchange Merkle root (`169dd7e2…740ad0`). Base mainnet was **not** attempted (`BASE_RPC_URL` missing; contract undeployed).

| Exchange | auditId | periodSeq | Block | `anchorSnapshot` tx |
| --- | --- | --- | --- | --- |
| kucoin | `20261031` | 1 | `12048458` | [`0xe2bca30654f409efa23732a72d9cd3e8ab6783bfd3b2b226c2abdd2cce413615`](https://sepolia.etherscan.io/tx/0xe2bca30654f409efa23732a72d9cd3e8ab6783bfd3b2b226c2abdd2cce413615) |
| binance | `PR31OCT26` | 48 | `12048459` | [`0x83801ca5505cf0d282981aa53f093f87f0db18a6c4ab4ab3399ae93786be66a1`](https://sepolia.etherscan.io/tx/0x83801ca5505cf0d282981aa53f093f87f0db18a6c4ab4ab3399ae93786be66a1) |
| bybit | `2026102706` | 1 | `12048460` | [`0xd541525cce2a622df1a135f889f69d8acafcc14ff600d197a532fa3018e46ebe`](https://sepolia.etherscan.io/tx/0xd541525cce2a622df1a135f889f69d8acafcc14ff600d197a532fa3018e46ebe) |
| bitget | `20261021` | 1 | `12048461` | [`0xd1e93ea100929b79b9ba845e5ad6736e696cfd48a869e2e38ce7ed3c123e6f54`](https://sepolia.etherscan.io/tx/0xd1e93ea100929b79b9ba845e5ad6736e696cfd48a869e2e38ce7ed3c123e6f54) |
| gateio | `20261018` | 1 | `12048462` | [`0xdb3ef9a800dba37b283e23ba50700a01046b219de0166c461505065af0e188ca`](https://sepolia.etherscan.io/tx/0xdb3ef9a800dba37b283e23ba50700a01046b219de0166c461505065af0e188ca) |
| bingx | `20261015` | 1 | `12034446` | [`0xe75c010b8108f90e28232f734553426e3efe7d565bcb971d8b9a7d247267c3b2`](https://sepolia.etherscan.io/tx/0xe75c010b8108f90e28232f734553426e3efe7d565bcb971d8b9a7d247267c3b2) (already anchored 2026-10-31; roots still match) |

Post-anchor `getSnapshotCount` (Sepolia): binance `5`, okx `4`, bingx `2`, kucoin `2`, mexc `1`, bybit `2`, bitget `2`, gateio `2`.

---

## Base mainnet (chain id `8453`) — pending

Intended production deployment of the same UUPS `ArdmerePoRAnchor` pattern. **No proxy or implementation address yet.**

### Blockers for this environment

- No `BASE_RPC_URL` / BaseScan API key configured here (Sepolia signer is available for testnet only).
- Base proxy/implementation not deployed; deploy must be signed by the project maintainer.

### Deploy checklist (maintainer)

1. Fund deployer on Base; export `PRIVATE_KEY`, `ANCHOR_SIGNER`, `BASE_RPC_URL` (and BaseScan key if verifying).
2. Deploy implementation + `ERC1967Proxy` via Foundry (`contracts/`), schema v3 storage.
3. Verify proxy + impl on BaseScan and Sourcify.
4. Record addresses in this file; set `ANCHOR_CONTRACT` / chain id in `.env.example` comments for production.
5. Optional: Base Sepolia dry-run before mainnet.
6. First `anchorSnapshot` tx for a public evaluation-set root; link tx here.

Until then, treat Sepolia anchors as **testnet evidence only**.


## Schema v4 upgrade (stage + dual coverage)

UUPS upgrade of the existing Sepolia proxy `0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7` (address unchanged).  # pragma: allowlist secret

```bash
export PATH="$HOME/.foundry/bin:$PATH"
cd contracts
# Required: PRIVATE_KEY (upgrade owner), ANCHOR_CONTRACT (proxy), SEPOLIA_RPC_URL / RPC_URL
forge script script/Upgrade.s.sol:Upgrade --rpc-url "$RPC_URL" --broadcast -vvvv
cast call "$ANCHOR_CONTRACT" "SCHEMA_VERSION()(uint8)" --rpc-url "$RPC_URL"   # expect 4
cast call "$ANCHOR_CONTRACT" "STORAGE_VERSION()(uint8)" --rpc-url "$RPC_URL"  # expect 2
```

New `SnapshotRecord` fields (appended after `anchoredAt`): `porStage`, `effectivePoR`, `walletCoverageBps`, `reserveCoverageBps`. Deprecated `coverageBps` is retained. Pre-upgrade v3 storage rows remain readable with the new fields as zero/false.


### Schema v4 status (2026-11-08)

- Sepolia proxy upgraded to schema v4 / storage v2.
- **Implementation (v4):** [`0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78`](https://sepolia.etherscan.io/address/0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78)  # pragma: allowlist secret
- **Impl deploy tx:** [`0xb746e715c68a5fd3a5e603b4f9f48e4f0158a45fc4f5225420c2b94a6c12a0f9`](https://sepolia.etherscan.io/tx/0xb746e715c68a5fd3a5e603b4f9f48e4f0158a45fc4f5225420c2b94a6c12a0f9)
- **Upgrade tx (v3→v4):** [`0xf4bdf0a43969ff31f344389d40bf7bd8032fd7bfc4f018843a9cb388032cf3b6`](https://sepolia.etherscan.io/tx/0xf4bdf0a43969ff31f344389d40bf7bd8032fd7bfc4f018843a9cb388032cf3b6)
- Verified live: `SCHEMA_VERSION()=4`, `STORAGE_VERSION()=2`; existing v3 snapshots (e.g. `PR01JUN26`, `506872725`) still readable via `getSnapshot` with `porStage=0`, `effectivePoR=false`, `walletCoverageBps=0`, `reserveCoverageBps=0`.
