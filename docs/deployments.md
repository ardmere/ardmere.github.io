# Deployments — ArdmerePoRAnchor

On-chain anchor contract deployments and first anchor transactions.

**Product target (ADR-002):** Base mainnet (`8453`) is the intended **production** primary anchor chain.  
**Operable today:** Ethereum Sepolia only. Base mainnet is **not deployed** in this refresh (requires deployer key + Base RPC — see checklist below).

| Network | Chain ID | Role | Status |
| --- | --- | --- | --- |
| Ethereum Sepolia | `11155111` | Active testnet / CLI default | **Live** (proxy schema **v7**; homepage six backfilled) |  # pragma: allowlist secret
| Base | `8453` | Production primary (ADR-002) | **Pending** — not deployed |
| Base Sepolia | `84532` | Optional L2 testnet | **Pending** — not deployed |

See [`decisions.md` ADR-002](./decisions.md#adr-002-primary-anchor-chain-base).

---

## Ethereum Sepolia (testnet, chain id `11155111`)  # pragma: allowlist secret

Primary **testnet** deployment (operable CLI default until Base mainnet ships). Contract is **UUPS upgradeable** — users interact with the **proxy** address; upgrades are authorized by the deployer (owner).

Deployment uses `PRIVATE_KEY` (upgrade/owner key used to broadcast) and `ANCHOR_SIGNER` (permissioned signer address). Supply them via `~/.zshenv` for local operators, or via the same env vars injected into a cloud Cursor agent environment (testnet only).

| Item | Value |
|---|---|
| **Proxy (ANCHOR_CONTRACT)** | [`0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7`](https://sepolia.etherscan.io/address/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7) |  # pragma: allowlist secret
| **Implementation (v4, stage + dual coverage)** | [`0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78`](https://sepolia.etherscan.io/address/0x71860B086ce8c550cBa2Af00950BDbe8C7a21e78) |  # pragma: allowlist secret
| **Implementation (v3, superseded)** | [`0x62050405283f222EFEdF3BF0d6Cb16541cd1327a`](https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a) |
| **Implementation (v2, superseded)** | [`0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392`](https://sepolia.etherscan.io/address/0x0c13bdD72e0a8439b0BEC29eB9e7bCb745BE9392) |
| **Anchor signer** | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |
| **Upgrade owner** | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |
| **Schema version** | `7` (backfilled rows store `schemaVersion` 7; untouched rows keep their write-time value) |
| **Storage version** | `4` |
| **Implementation (v7 + reserve backfill)** | [`0xe8DB40AA87138d79437ED0DceA42D3A7A2AB95EB`](https://sepolia.etherscan.io/address/0xe8DB40AA87138d79437ED0DceA42D3A7A2AB95EB) |
| **Impl deploy tx (v7 backfill)** | [`0x288973f92b5287a4f571fb82c5590ac79862005eba5716cb076696e2dd795c0d`](https://sepolia.etherscan.io/tx/0x288973f92b5287a4f571fb82c5590ac79862005eba5716cb076696e2dd795c0d) (block `12119163`) |
| **Upgrade tx (backfill impl)** | [`0x0b56907e4f925ab9ccd80cae23659e2ce124e90b9376b41ff2175015b3a62b70`](https://sepolia.etherscan.io/tx/0x0b56907e4f925ab9ccd80cae23659e2ce124e90b9376b41ff2175015b3a62b70) (block `12119253`) |
| **Implementation (v7, no backfill, superseded)** | [`0x08B0d5d421578f76862dB4B8A2B8969502638929`](https://sepolia.etherscan.io/address/0x08B0d5d421578f76862dB4B8A2B8969502638929) |
| **Impl deploy tx (v7)** | [`0x7ecd31a796094c0eeabc2f9bff5841f77c56f813ed2b99b83928e93d5da246d9`](https://sepolia.etherscan.io/tx/0x7ecd31a796094c0eeabc2f9bff5841f77c56f813ed2b99b83928e93d5da246d9) (block `12115204`) |
| **Upgrade tx (v6→v7)** | [`0x64e2a1d69cd22a9fb04800c6afde98d5b949f7a43a40aa389226631078f09f81`](https://sepolia.etherscan.io/tx/0x64e2a1d69cd22a9fb04800c6afde98d5b949f7a43a40aa389226631078f09f81) (block `12115207`) |
| **Implementation (v6, composite snapshotId, superseded)** | [`0xB0A204285e6Cc7Eaf9fA21E327668Fc808032DEe`](https://sepolia.etherscan.io/address/0xB0A204285e6Cc7Eaf9fA21E327668Fc808032DEe) |
| **Impl deploy tx (v6)** | [`0xe6b5687974f34214a94c60651cc47085ced433ba8d0447c7af28a77b96b2b898`](https://sepolia.etherscan.io/tx/0xe6b5687974f34214a94c60651cc47085ced433ba8d0447c7af28a77b96b2b898) |
| **Upgrade tx (v5→v6)** | [`0x43ac3aa0d85bb826aedb70cac2cdd4860e0f42639f12013d3f0c607a26d25dc3`](https://sepolia.etherscan.io/tx/0x43ac3aa0d85bb826aedb70cac2cdd4860e0f42639f12013d3f0c607a26d25dc3) |
| **Implementation (v5, custody risk, superseded)** | [`0x8A0E063Ad547adBA8fC5B7D0d0F9Cf55dF671Baf`](https://sepolia.etherscan.io/address/0x8A0E063Ad547adBA8fC5B7D0d0F9Cf55dF671Baf) |
| **Impl deploy tx (v5)** | [`0x88afa3f6482d5207fd652e7ce1c8df466021372f1989444b8885fa4d78631a9e`](https://sepolia.etherscan.io/tx/0x88afa3f6482d5207fd652e7ce1c8df466021372f1989444b8885fa4d78631a9e) |
| **Upgrade tx (v4→v5)** | [`0x6d41d2bc5af19a583e028a24ff324f32807008c27bc56cf817bc0b51f9ede682`](https://sepolia.etherscan.io/tx/0x6d41d2bc5af19a583e028a24ff324f32807008c27bc56cf817bc0b51f9ede682) |
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

## Schema v5 upgrade (custody risk on-chain)

UUPS upgrade of the Sepolia proxy (address unchanged). Adds `custodyRiskScore` (`uint8`, **0–100** or **`255`** = unavailable/inconclusive) to `SnapshotRecord` and `anchorSnapshot` calldata. Aligns with NEAR sibling semantics.

```bash
export PATH="$HOME/.foundry/bin:$PATH"
cd contracts
# Required: PRIVATE_KEY (upgrade owner), ANCHOR_CONTRACT (proxy), RPC_URL
forge script script/Upgrade.s.sol:Upgrade --rpc-url "$RPC_URL" --broadcast -vvvv
cast call "$ANCHOR_CONTRACT" "SCHEMA_VERSION()(uint8)" --rpc-url "$RPC_URL"   # expect 5
cast call "$ANCHOR_CONTRACT" "STORAGE_VERSION()(uint8)" --rpc-url "$RPC_URL"  # expect 3
```

**Anchor-time score:** `por anchor` / `internal/pipeline` invoke `skills/por-custody-risk/scripts/risk_score.py --no-anchor` on `docs/reports/<exchange>/<auditId>-assessment.json` so the score uses the **same bundle roots** being anchored (no live anchor read; avoids `ANCHOR_MISMATCH` chicken-and-egg). If scoring is inconclusive, calldata uses `255`.

**Homepage re-anchor (human-operated):** `./scripts/reanchor-homepage-v5.sh` lists the current `HOMEPAGE_REPORT_SET`; `scripts/anchor-assessment-v5.py` broadcasts from assessment JSON + `risk_score.py --no-anchor` custody score. Immutable `snapshotId` rows cannot be updated on-chain.

### Homepage v5 anchors (2026-11-09)

| Exchange | auditId | custodyRiskScore (at anchor) | `anchorSnapshot` tx |
| --- | --- | --- | --- |
| binance | `PR04NOV26` | 91 | [`0xd582cd5893d60d7f3f46dc623fac4dbc280e45c3ff1a7f02222dd15c7e4c4ea8`](https://sepolia.etherscan.io/tx/0xd582cd5893d60d7f3f46dc623fac4dbc280e45c3ff1a7f02222dd15c7e4c4ea8) |
| bitunix | `20261103` | 100 | [`0xfd2cdca7e98895de5a0bf5fa4259149862728571ea6c895fb10589b06e5841bb`](https://sepolia.etherscan.io/tx/0xfd2cdca7e98895de5a0bf5fa4259149862728571ea6c895fb10589b06e5841bb) |
| phemex | `20261101` | 100 | [`0x262cc925d88fff16f3ab3f134e0d18e88ab97bc42cb93a9b332ec336264d7cc7`](https://sepolia.etherscan.io/tx/0x262cc925d88fff16f3ab3f134e0d18e88ab97bc42cb93a9b332ec336264d7cc7) |
| okx | `508756245` | 10 (periodSeq **4** — assessment `0` was taken) | [`0x1924fa2c126f6c0bf5d10781eeceea7a754caa94ad07a4b1e16b5f83c95124cc`](https://sepolia.etherscan.io/tx/0x1924fa2c126f6c0bf5d10781eeceea7a754caa94ad07a4b1e16b5f83c95124cc) |

Legacy v5 rows keyed by `keccak256(bytes(auditId))` only; schema v6 composite keys disambiguate colliding audit ids (see homepage v6 table below).

**Operator cast signature (v5):**

```bash
cast send "$ANCHOR_CONTRACT" \
  'anchorSnapshot(string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,bool,uint16,uint16,uint8)' \
  "<exchange>" <snapshotId_b32> <periodSeq> <snapshotTimeUnix> <btcBlockHeight> \
  <exchangeMerkleRoot> <artifactBundleRoot> <verificationBundleRoot> \
  <verdictSummary> <coverageBps> <porStage> <effectivePoR> <walletCoverageBps> <reserveCoverageBps> <custodyRiskScore> \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

Signer on Sepolia: `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F`.

## Schema v6 upgrade (composite snapshotId)

UUPS upgrade of the Sepolia proxy (address unchanged). **`snapshotId` for new anchors** is `keccak256(abi.encodePacked(normalizeExchange(exchange), 0x00, auditId))`; `anchorSnapshot` calldata gains `string auditId` and reverts `InvalidSnapshotId` when the bytes32 id does not match `expectedSnapshotId`. Event layout unchanged; `SCHEMA_VERSION` becomes **6**. Legacy v3–v5 keys remain readable via `legacySnapshotId` / client fallback.

```bash
export PATH="$HOME/.foundry/bin:$PATH"
cd contracts
forge script script/Upgrade.s.sol:Upgrade --rpc-url "$RPC_URL" --broadcast -vvvv
cast call "$ANCHOR_CONTRACT" "SCHEMA_VERSION()(uint8)" --rpc-url "$RPC_URL"   # expect 6
```

**Operator cast signature (v6):**

```bash
cast send "$ANCHOR_CONTRACT" \
  'anchorSnapshot(string,string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,bool,uint16,uint16,uint8)' \
  "<exchange>" "<auditId>" <snapshotId_b32> <periodSeq> <snapshotTimeUnix> <btcBlockHeight> \
  <exchangeMerkleRoot> <artifactBundleRoot> <verificationBundleRoot> \
  <verdictSummary> <coverageBps> <porStage> <effectivePoR> <walletCoverageBps> <reserveCoverageBps> <custodyRiskScore> \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

Derive `snapshotId_b32` with `cast call "$ANCHOR_CONTRACT" "expectedSnapshotId(string,string)(bytes32)" "<exchange>" "<auditId>" --rpc-url "$RPC_URL"` or `./por anchor -exchange <id> -audit <auditId>` (prints schema v6 calldata). **Do not broadcast** without signer keys.

### Homepage v6 composite anchors (2026-11-09)

After the v5→v6 UUPS upgrade, homepage rows were re-anchored under composite `snapshotId` keys (`expectedSnapshotId(exchange, auditId)`). `periodSeq` values below are the on-chain index used at anchor time (`getSnapshotCount` when the prior period was already consumed).

| Exchange | auditId | custodyRiskScore | periodSeq | `anchorSnapshot` tx |
| --- | --- | ---: | ---: | --- |
| ourbit | `20261101` | 100 | 0 | [`0x2afe397dc58654957d013fefc4736cf2fea713086c847a47d2d4d8930861f9ef`](https://sepolia.etherscan.io/tx/0x2afe397dc58654957d013fefc4736cf2fea713086c847a47d2d4d8930861f9ef) |
| htx | `20261101` | 99 | 1 | [`0x80960fdd7339ef5dc6ebe291b60080723281b7c0aa75d6d598e08151a1f7fcd0`](https://sepolia.etherscan.io/tx/0x80960fdd7339ef5dc6ebe291b60080723281b7c0aa75d6d598e08151a1f7fcd0) |
| binance | `PR04NOV26` | 91 | 6 | [`0x390bbf0a0eaae3d1cc94604496e9ff18ceb22fbe18565b71079977f1f471bdf8`](https://sepolia.etherscan.io/tx/0x390bbf0a0eaae3d1cc94604496e9ff18ceb22fbe18565b71079977f1f471bdf8) |
| bitunix | `20261103` | 100 | 1 | [`0x63068fcb8e6a9951ed2995a008d014bfbf0b7d41b9e924cac3b575dca29472f8`](https://sepolia.etherscan.io/tx/0x63068fcb8e6a9951ed2995a008d014bfbf0b7d41b9e924cac3b575dca29472f8) |
| phemex | `20261101` | 100 | 1 | [`0xadb39d62b65ea7d229b7ee561971ec6774a593fe9a4facd85b6d702ff9bf1aa4`](https://sepolia.etherscan.io/tx/0xadb39d62b65ea7d229b7ee561971ec6774a593fe9a4facd85b6d702ff9bf1aa4) |
| okx | `508756245` | 10 | 5 | [`0xe274ec7808eb4eda9918f4267fdf14679e342c8534d5aa8d18f92d73ed92ceab`](https://sepolia.etherscan.io/tx/0xe274ec7808eb4eda9918f4267fdf14679e342c8534d5aa8d18f92d73ed92ceab) |

**Colliding audit ids (composite keys, v6+):**

| Exchange | auditId | `expectedSnapshotId` (composite) |
| --- | --- | --- |
| phemex | `20261101` | `0x12abd55db2e6d9b44edf24364ab8735226d612d54b228074210cb4933ec39f80` |
| htx | `20261101` | `0xbdbd096d9043bf0b55246ef612e4db87ac27eb5718e8debd15fdd77db30e6f20` |
| ourbit | `20261101` | `0xcf880c45c95eb72602f327c669c070e2cf047bb236b57614ed774c1fca9e6412` |

Legacy `keccak256(bytes("20261101"))` = `0x9565ad3fb66b361d5dcf03c74124363797d48a3b2688f8bccfa728c066d061a5` — historical v3–v5 rows only; do **not** use for new anchors.

### Full manifest v6 anchors

Every file under `docs/reports/*/*-assessment.json` should have a composite v6 row with a live custody score that matches `skills/por-custody-risk/scripts/risk_score.py` (see `custodyRiskScoreMatchesComputed`).

```bash
# Dry-run: list missing/mismatched rows and print cast send calldata
./scripts/batch-anchor-v6-composite.py

# Broadcast (maintainer key; Sepolia only)
export ANCHOR_CONTRACT=<proxy from table above>
export RPC_URL=<Sepolia RPC>
export PRIVATE_KEY=<anchor signer key>
./scripts/batch-anchor-v6-composite.py --apply
```

Single-file helper: `scripts/anchor-assessment-v6.py [--dry-run] [--period-seq N] docs/reports/<ex>/<audit>-assessment.json`.

---

## Schema v7 upgrade (self-reported totals and major-coin ratios)

UUPS upgrade of the Sepolia proxy (address unchanged), broadcast 2026-11-12. `SCHEMA_VERSION()` is **7** and `STORAGE_VERSION()` is **4**. Implementation runtime matches the current `ArdmerePoRAnchor` build aside from the deploy-time `address(this)` immutables.

| Item | Value |
| --- | --- |
| Implementation | [`0x08B0d5d421578f76862dB4B8A2B8969502638929`](https://sepolia.etherscan.io/address/0x08B0d5d421578f76862dB4B8A2B8969502638929) |
| Impl deploy tx | [`0x7ecd31a796094c0eeabc2f9bff5841f77c56f813ed2b99b83928e93d5da246d9`](https://sepolia.etherscan.io/tx/0x7ecd31a796094c0eeabc2f9bff5841f77c56f813ed2b99b83928e93d5da246d9) (block `12115204`, gas used `16863696`) |
| Upgrade tx (v6→v7) | [`0x64e2a1d69cd22a9fb04800c6afde98d5b949f7a43a40aa389226631078f09f81`](https://sepolia.etherscan.io/tx/0x64e2a1d69cd22a9fb04800c6afde98d5b949f7a43a40aa389226631078f09f81) (block `12115207`) |

Two earlier creates ran out of gas and have no code (`0x53D6eE03E77c4030c91e88744D7B54bdF9931A70`, `0x5525d28ec38ff3ad879e9b5b7fae022ab303a80f`, `0x258a3e4a363b5be842eca60c6add7d95a34459cc`). They are not the implementation. The successful create used gas limit `17000000`.

`snapshotId` stays `expectedSnapshotId(exchange, auditId)`. A new `periodSeq` does not create a second record (`SnapshotAlreadyExists`, selector `0xb296876b`). `backfillReserveFigures` is owner-only: the row must already exist, and a second call reverts `ReserveBackfillAlreadyApplied` once `schemaVersion` is 7. It writes only the two USD totals, the four ratio bps, and `schemaVersion`. Score, time, roots, `periodSeq`, and `snapshotId` stay as they were.

### Homepage reserve backfill (2026-11-12)

Figures are whole USD from the assessment JSON (`round`; null encodes as `0`) and ratio bps (`round(reserveRate × 10000)`; missing coin `65535`; saturate at `65534`). Read-back: `schemaVersion` 7, custody score unchanged.

| Exchange | auditId | Custody | Liability USD | Reserve USD | BTC | ETH | USDT | USDC | `backfillReserveFigures` tx |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| binance | `PR04NOV26` | 91 | 132692864216 | 133789440153 | 10003 | 10002 | 10155 | 10312 | [`0xbba18f7ee15dc08e67eadcb8f04c58a13161b0aaadc2196fa45ffe84e301ccbc`](https://sepolia.etherscan.io/tx/0xbba18f7ee15dc08e67eadcb8f04c58a13161b0aaadc2196fa45ffe84e301ccbc) |
| bitunix | `20261103` | 100 | 0 | 0 | 16100 | 10500 | 10600 | 65535 | [`0x47ed420f459ed8653921978167284bdcc69cf2319366afb4ba52ec17061befae`](https://sepolia.etherscan.io/tx/0x47ed420f459ed8653921978167284bdcc69cf2319366afb4ba52ec17061befae) |
| phemex | `20261101` | 100 | 128774624 | 140796757 | 10329 | 11354 | 65535 | 65535 | [`0x806d7453a60237728b19a4d7a5ff752ff6d62852b643ad755a178964f039cbb5`](https://sepolia.etherscan.io/tx/0x806d7453a60237728b19a4d7a5ff752ff6d62852b643ad755a178964f039cbb5) |
| htx | `20261101` | 99 | 5182100387 | 5305874609 | 10219 | 10145 | 3502 | 65534 | [`0xa12e73b7ee366b0ef963d032f2f548a0ab19c4b1e3265d28a741fb95ec3bd128`](https://sepolia.etherscan.io/tx/0xa12e73b7ee366b0ef963d032f2f548a0ab19c4b1e3265d28a741fb95ec3bd128) |
| okx | `508756245` | 10 | 30967089361 | 32618136617 | 10839 | 10051 | 10643 | 10002 | [`0x2c8f05ac44d33f0793c38c64976d9dd47259ccdb742840aef068692c9ea37da0`](https://sepolia.etherscan.io/tx/0x2c8f05ac44d33f0793c38c64976d9dd47259ccdb742840aef068692c9ea37da0) |
| ourbit | `20261101` | 100 | 85351838 | 93871851 | 11073 | 10810 | 10913 | 65535 | [`0x80dae8a901aefaa2ca7d7ccacdd2773987ea801718012adfa630614628c77013`](https://sepolia.etherscan.io/tx/0x80dae8a901aefaa2ca7d7ccacdd2773987ea801718012adfa630614628c77013) |

### Remaining assessment reserve backfill (2026-11-12)

Batch sweep of all `docs/reports/*/*-assessment.json` (57 files) via `scripts/batch-backfill-reserve-v7.py`. Only **composite** `expectedSnapshotId(exchange, auditId)` rows are eligible; missing composite anchors are listed and skipped (no new `anchorSnapshot` in this pass). Full per-row status + tx appendix: [`docs/v7-reserve-backfill-log.json`](./v7-reserve-backfill-log.json).

| Outcome | Count |
| --- | ---: |
| Already schema ≥ 7 (homepage six) | 6 |
| Backfilled this pass | 3 |
| Missing composite anchor (skipped) | 48 |
| Null USD figures but written (`0`/`65535`) | 3 |
| Failed | 0 |

The three remaining composite rows (all binance, schema 6 → 7; assessment totals null → on-chain `0` USD and `65535` bps):

| Exchange | auditId | Liability USD | Reserve USD | BTC | ETH | USDT | USDC | `backfillReserveFigures` tx |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| binance | `PR01APR26` | 0 | 0 | 65535 | 65535 | 65535 | 65535 | [`0x0da58a544b8ff50ddc21f714502c2b89048c0e5fa0400ff6d4a116c5c1779be0`](https://sepolia.etherscan.io/tx/0x0da58a544b8ff50ddc21f714502c2b89048c0e5fa0400ff6d4a116c5c1779be0) |
| binance | `PR01AUG26` | 0 | 0 | 65535 | 65535 | 65535 | 65535 | [`0xf84d0c57c5771f47402cb06f08a5d9675310f750c3a4bb456aff4ebf5f30416a`](https://sepolia.etherscan.io/tx/0xf84d0c57c5771f47402cb06f08a5d9675310f750c3a4bb456aff4ebf5f30416a) |
| binance | `PR01JUL26` | 0 | 0 | 65535 | 65535 | 65535 | 65535 | [`0x404eb801784fc37e7b4442e1aa768e39298bcc2d29f922b9ba29deb77025046d`](https://sepolia.etherscan.io/tx/0x404eb801784fc37e7b4442e1aa768e39298bcc2d29f922b9ba29deb77025046d) |

```bash
# dry-run then apply (owner key; skips schemaVersion >= 7)
python3 scripts/batch-backfill-reserve-v7.py
python3 scripts/batch-backfill-reserve-v7.py --apply

cast send "$ANCHOR_CONTRACT" \
  'backfillReserveFigures(string,string,uint128,uint128,uint16,uint16,uint16,uint16)' \
  "<exchange>" "<auditId>" <totalLiabilityUsd> <totalReserveUsd> \
  <btcReserveRatioBps> <ethReserveRatioBps> <usdtReserveRatioBps> <usdcReserveRatioBps> \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

New `SnapshotRecord` slots, appended after `custodyRiskScore`:

| Field | Type | Sentinel |
| --- | --- | --- |
| `totalLiabilityUsd` | `uint128` | `0` = unavailable (no published amounts, or no priced coins). Whole USD, nearest dollar. |
| `totalReserveUsd` | `uint128` | `0` = unavailable. Self-custody + third-party, same CoinGecko-at-`snapshotTime` path as `reserveCoverageBps`. |
| `btcReserveRatioBps` | `uint16` | `10000` = 100%. `65535` = no BTC row or no ratio. |
| `ethReserveRatioBps` | `uint16` | same |
| `usdtReserveRatioBps` | `uint16` | same |
| `usdcReserveRatioBps` | `uint16` | same |

Ratios above 655.34% saturate at `65534` so they cannot collide with `65535`. Pre-v7 storage rows still decode the new slots as **0**. Readers must treat both USD totals and the four ratios as unavailable unless `schemaVersion >= 7` (a 0 bps reading on an old row is not a 0% reserve ratio). These figures are self-reported snapshot amounts, not a solvency opinion.

```bash
cast call "$ANCHOR_CONTRACT" "SCHEMA_VERSION()(uint8)" --rpc-url "$RPC_URL"   # 7
cast call "$ANCHOR_CONTRACT" "STORAGE_VERSION()(uint8)" --rpc-url "$RPC_URL"  # 4
```

**Operator cast signature (v7):**

```bash
cast send "$ANCHOR_CONTRACT" \
  'anchorSnapshot(string,string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,bool,uint16,uint16,uint8,uint128,uint128,uint16,uint16,uint16,uint16)' \
  "<exchange>" "<auditId>" <snapshotId_b32> <periodSeq> <snapshotTimeUnix> <btcBlockHeight> \
  <exchangeMerkleRoot> <artifactBundleRoot> <verificationBundleRoot> \
  <verdictSummary> <coverageBps> <porStage> <effectivePoR> <walletCoverageBps> <reserveCoverageBps> <custodyRiskScore> \
  <totalLiabilityUsd> <totalReserveUsd> <btcReserveRatioBps> <ethReserveRatioBps> <usdtReserveRatioBps> <usdcReserveRatioBps> \
  --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY"
```

`./por anchor` prints this v7 payload and does not broadcast. `scripts/anchor-assessment-v6.py` still targets the v6 selector and must not be pointed at v7 calldata. Homepage totals are the `backfillReserveFigures` rows above, not a second `anchorSnapshot`.

### Missing-composite v7 anchor batch (2026-11-12)

Anchored the 48 assessment rows that still lacked a composite `expectedSnapshotId` via `scripts/batch-anchor-v7.py` (full v7 `anchorSnapshot` in one shot: roots, custody score from `risk_score.py --no-anchor`, USD totals, four ratio bps). The nine rows already at schema ≥ 7 were skipped. Full per-row tx appendix: [`docs/v7-anchor-missing-log.json`](./v7-anchor-missing-log.json).

| Outcome | Count |
| --- | ---: |
| Anchored (new composite schema 7) | 48 |
| Skipped already exists | 9 |
| Failed | 0 |

Sample txs (read-back confirmed `schemaVersion=7` and encoded figures):

| Exchange | auditId | Custody | Liability USD | Reserve USD | BTC | ETH | USDT | USDC | `anchorSnapshot` tx |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| binance | `PR01JUN26` | 91 | 0 | 0 | 65535 | 65535 | 65535 | 65535 | [`0xb9dffe6a9c0be9315cdb702d1df061536f85866723434347c66af5f97384e43c`](https://sepolia.etherscan.io/tx/0xb9dffe6a9c0be9315cdb702d1df061536f85866723434347c66af5f97384e43c) |
| poloniex | `20261001` | 99 | 2312405057 | 2323214460 | 10010 | 10059 | 10033 | 65535 | [`0xb90a4e28ed8ab909f8b714aae7cde3a4f9992f04bcc819b9e753a5e5c6122fc1`](https://sepolia.etherscan.io/tx/0xb90a4e28ed8ab909f8b714aae7cde3a4f9992f04bcc819b9e753a5e5c6122fc1) |
| toobit | `20261001` | 100 | 0 | 0 | 10400 | 10700 | 10500 | 10800 | [`0xe3941cf5ffd67d42f3c29188afef72af2badd3c6e302b96bf6d04db41fb6645c`](https://sepolia.etherscan.io/tx/0xe3941cf5ffd67d42f3c29188afef72af2badd3c6e302b96bf6d04db41fb6645c) |
| okx | `506872725` | 10 | 22998236121 | 24652491059 | 10607 | 10245 | 10989 | 65535 | [`0xbd1fbdd67aa20e46c62c7820f34090fd8b7cd88658fd86de3ea3eae54a49784f`](https://sepolia.etherscan.io/tx/0xbd1fbdd67aa20e46c62c7820f34090fd8b7cd88658fd86de3ea3eae54a49784f) |
| gateio | `20261018` | 99 | 4245923696 | 5023237868 | 12160 | 12378 | 11168 | 10561 | [`0x33bddc0de453c40593b24dd4b53c0ae71f96ec953fa6032231ba72b33db5276c`](https://sepolia.etherscan.io/tx/0x33bddc0de453c40593b24dd4b53c0ae71f96ec953fa6032231ba72b33db5276c) |
| woox | `20261015` | 100 | 0 | 0 | 65535 | 65535 | 65535 | 65535 | [`0x7590057a16110861df19619eeb44b9f1aa76a83d43a079275069825327fa932b`](https://sepolia.etherscan.io/tx/0x7590057a16110861df19619eeb44b9f1aa76a83d43a079275069825327fa932b) |

```bash
python3 scripts/batch-anchor-v7.py          # dry-run
python3 scripts/batch-anchor-v7.py --apply  # signer key; skips existing composites
```

After that batch, every assessment JSON then on `main` had a composite Sepolia row with `schemaVersion` 7.

### Binance `PR12NOV26` (2026-11-14)

Post-batch assessment from the merged scan report ([PR #57](https://github.com/ardmere/ardmere.github.io/pull/57)). Composite `expectedSnapshotId(binance, PR12NOV26)` was already on Sepolia schema **7** when the cloud env checked (no duplicate `anchorSnapshot`). Read-back + `SnapshotAnchored` log appendix: [`docs/v7-anchor-binance-pr12nov26-log.json`](./v7-anchor-binance-pr12nov26-log.json).

| Exchange | auditId | periodSeq | Custody | Liability USD | Reserve USD | BTC | ETH | USDT | USDC | Block | `anchorSnapshot` tx |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| binance | `PR12NOV26` | 15 | 91 | 128582258733 | 130015809492 | 10005 | 10002 | 10434 | 10007 | `12129857` | [`0xaf3cdd6dfce7e044f25af2ff5e8c0ade4c88de621701084b6877bea090730bae`](https://sepolia.etherscan.io/tx/0xaf3cdd6dfce7e044f25af2ff5e8c0ade4c88de621701084b6877bea090730bae) |

- **snapshotId:** `0xe0cbbc36668618dced5cc703e0b7d1c16eefe2d4e1d9432ce04364cd357fcb77`
- **Proxy / impl / schema:** unchanged (`SCHEMA_VERSION()=7`, live impl `0xe8DB40AA87138d79437ED0DceA42D3A7A2AB95EB`)
- **Signer:** `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F`

All **58** `docs/reports/*/*-assessment.json` files now have a composite Sepolia row with `schemaVersion` 7.

## NEAR testnet (composite parity) — pending on-chain

Sibling deployment in [`ardmere/ardmere-near-demo`](https://github.com/ardmere/ardmere-near-demo) (schema v5 custody field + composite `snapshot_id` semantics). **On-chain NEAR backfill is pending** until the testnet contract WASM is upgraded and the owner key is available.

| Item | Status |
| --- | --- |
| WASM schema v5/v6 fields | **Pending** — redeploy `por_digest_anchor.wasm` |
| `backfill-custody --apply` | **Pending** |
| `reconcile --apply` (htx / ourbit composite ids) | **Pending** |

**Operator checklist** (human; requires `NEAR_TESTNET_OWNER_KEY` / keychain — not on cloud agents):

1. `cd ardmere-near-demo && cargo build --manifest-path contracts/por_digest_anchor/Cargo.toml --target wasm32-unknown-unknown --release`
2. Deploy WASM to `NEAR_CONTRACT_ID` on testnet (`near contract deploy … use-file …/por_digest_anchor.wasm`)
3. `cast call`-equivalent: `near contract call-function as-read-only "$NEAR_CONTRACT_ID" schema_version json-args '{}' network-config testnet now` → expect v5+
4. `python3 scripts/near_digest.py --network testnet --contract "$NEAR_CONTRACT_ID" backfill-custody` (dry-run)
5. `python3 scripts/near_digest.py --network testnet --contract "$NEAR_CONTRACT_ID" backfill-custody --apply`
6. `python3 scripts/near_digest.py --network testnet --contract "$NEAR_CONTRACT_ID" reconcile --apply` for composite collisions

See `ardmere-near-demo` README and `.github/workflows/anchor-testnet.yml` for exact `near-cli-rs` flags.
