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

---

## Base mainnet (chain id `8453`) — pending

Intended production deployment of the same UUPS `ArdmerePoRAnchor` pattern. **No proxy or implementation address yet.**

### Blockers for this environment

- Operator `PRIVATE_KEY` / `ANCHOR_SIGNER` not present in Cloud Agent env (must stay in maintainer `~/.zshenv`).
- No `BASE_RPC_URL` / BaseScan API key configured here.
- Deploy must be signed by the project maintainer; this refresh only documents the gap.

### Deploy checklist (maintainer)

1. Fund deployer on Base; export `PRIVATE_KEY`, `ANCHOR_SIGNER`, `BASE_RPC_URL` (and BaseScan key if verifying).
2. Deploy implementation + `ERC1967Proxy` via Foundry (`contracts/`), schema v3 storage.
3. Verify proxy + impl on BaseScan and Sourcify.
4. Record addresses in this file; set `ANCHOR_CONTRACT` / chain id in `.env.example` comments for production.
5. Optional: Base Sepolia dry-run before mainnet.
6. First `anchorSnapshot` tx for a public evaluation-set root; link tx here.

Until then, treat Sepolia anchors as **testnet evidence only**.
