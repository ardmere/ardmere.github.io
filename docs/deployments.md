# Deployments — ArdmerePoRAnchor

On-chain anchor contract deployments and first anchor transactions.

**Product target (ADR-002):** Base mainnet (`8453`) is the intended **production** primary anchor chain.  
**Operable today:** Ethereum Sepolia only. Base mainnet is **not deployed** in this refresh (requires deployer key + Base RPC — see checklist below).

| Network | Chain ID | Role | Status |
| --- | --- | --- | --- |
| Ethereum Sepolia | `11155111` | Active testnet / CLI default | **Live** (proxy + impl v3) |
| Base | `8453` | Production primary (ADR-002) | **Pending** — not deployed |
| Base Sepolia | `84532` | Optional L2 testnet | **Pending** — not deployed |

See [`decisions.md` ADR-002](./decisions.md#adr-002-primary-anchor-chain-base).

---

## Ethereum Sepolia (testnet, chain id `11155111`)

Primary **testnet** deployment (operable CLI default until Base mainnet ships). Contract is **UUPS upgradeable** — users interact with the **proxy** address; upgrades are authorized by the deployer (owner).

Deployment using `PRIVATE_KEY` / `ANCHOR_SIGNER` from `~/.zshenv` (testnet only).

| Item | Value |
|---|---|
| **Proxy (ANCHOR_CONTRACT)** | [`0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7`](https://sepolia.etherscan.io/address/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7) |
| **Implementation (v3, query storage)** | [`0x62050405283f222EFEdF3BF0d6Cb16541cd1327a`](https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a) |
| **Implementation (v2, superseded)** | [`0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392`](https://sepolia.etherscan.io/address/0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392) |
| **Anchor signer** | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |
| **Upgrade owner** | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |
| **Schema version** | `3` (on-chain storage + merged single-tx anchor) |
| **Storage version** | `1` |
| **Pattern** | UUPS (`ERC1967Proxy` + `ArdmerePoRAnchor` implementation) |
| **Impl deploy tx (v3)** | [`0x6aaa54524151d18be495f26fcf60b4f4ed155719b891b2ebeaf36bab9375b96f`](https://sepolia.etherscan.io/tx/0x6aaa54524151d18be495f26fcf60b4f4ed155719b891b2ebeaf36bab9375b96f) |
| **Upgrade tx (v2→v3)** | [`0x905a68386dc49096208e4f75a060aa2dddb55f1b36c82735f8de73f471e75699`](https://sepolia.etherscan.io/tx/0x905a68386dc49096208e4f75a060aa2dddb55f1b36c82735f8de73f471e75699) |
| **Impl deploy tx (v2)** | [`0xbf4684b13b133370e71ae946b2972e233d6156aa5c208c2d11b7b3778da3ca39`](https://sepolia.etherscan.io/tx/0xbf4684b13b133370e71ae946b2972e233d6156aa5c208c2d11b7b3778da3ca39) |
| **Proxy deploy tx** | [`0xbd3b1a874831f37768d08a28bb015bbd39261b3bf32ecad99df652a1cb772e7d`](https://sepolia.etherscan.io/tx/0xbd3b1a874831f37768d08a28bb015bbd39261b3bf32ecad99df652a1cb772e7d) |
| **Deploy block** | `11182041` |
| **Source verified (Etherscan)** | [Proxy](https://sepolia.etherscan.io/address/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7#code) ✓ · [Implementation v3](https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a#code) ✓ · [Implementation v2 (superseded)](https://sepolia.etherscan.io/address/0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392#code) ✓ |
| **Source verified (Sourcify)** | [Proxy](https://repo.sourcify.dev/contracts/full_match/11155111/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7/) · [Implementation v3](https://repo.sourcify.dev/contracts/full_match/11155111/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a/) · [Implementation v2 (superseded)](https://repo.sourcify.dev/contracts/full_match/11155111/0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392/) |

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
