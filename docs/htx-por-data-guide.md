# HTX Proof-of-Reserves — data guide

> HTX (Huobi) zk PoR v2 uses Groth16 (`zkpor500`, 174 assets, batch 500).  
> Companion audit: see `PoR/htx-audit-report/HTX-zkPoR-设计审计报告.md` in the research repo.

## Public sources

| Source | URL | Content |
|--------|-----|---------|
| PoR page | https://www.htx.com/zh-cn/finance/merkle/ | Reserve ratios (browser); no stable public API |
| Sample zk bundle | https://github.com/huobiapi/Tool-Go-MerkleVerify/releases/download/2.0.0/public-data.zip | `config.json` + `proof0.csv` + `zkpor500.vk.save` (2023-09-10) |
| zk verifier binary | https://github.com/huobiapi/Tool-Go-MerkleVerify/releases/download/2.0.0/zkverifier-macos-x64.zip | `zkverifiermac` (expects `./public-data/` cwd layout) |
| Monthly reports | HTX PoR Reports (login / dashboard) | Production `proof0.csv` per audit period |

## Local layout

```text
artifacts/htx/<auditId>/
  raw/<sha256>.json          # summarySnapshot (may be zk-derived)
  raw/<sha256>.zip           # globalProofBundle (public-data.zip)
  fetch.json
  bundles/
```

## CLI

```bash
# Archive sample zk bundle (derives summary from proof0.csv)
go run ./cmd/por fetch htx -zk-bundle ./public-data.zip

# Full pipeline
go run ./cmd/por anchor -exchange htx -zk-bundle ./public-data.zip

# Cryptographic verify (optional; ~minutes on sample data)
export HTX_ZK_VERIFIER=/path/to/zkverifiermac   # run from dir containing public-data/
go run ./cmd/por verify -exchange htx -snapshot 20230910 -artifacts ./artifacts/htx/20230910
```

Flags: `-summary-path` (browser-captured reserve ratios), `-zk-bundle` (public-data.zip).

## Verifier profile (`htx`)

| Verifier | Status | Notes |
|----------|--------|-------|
| `artifact-integrity@1` | active | SHA256 over archived artifacts |
| `solvency-claim@1` | UNVERIFIABLE if zk-only summary | Needs browser reserve-ratio JSON |
| `global-zk-proof@htx-1` | PARTIAL / PASS | Structure + merkle bind; `HTX_ZK_VERIFIER` for Groth16 |
| `internal-consistency@0` | UNVERIFIABLE | No public wallet CSV |
| `onchain-balance-*@0` | UNVERIFIABLE | No public address list |

## ZK proof boundary (mandatory)

HTX Groth16 proves batch SMT insertion + platform `cumDebt ≤ cumEquity` + 174-asset commitment chain.  
It does **not** prove per-user solvency, consistent asset aggregation (HTXPOR-01), equity/debt binding (HTXPOR-02), or on-chain reserves.

## References

- [Tool-Go-MerkleVerify](https://github.com/huobiapi/Tool-Go-MerkleVerify)
- [HTX zk PoR upgrade announcement](https://www.htx.com/en-in/feed/news/724709/)
