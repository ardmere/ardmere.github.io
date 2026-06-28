# PoR Artifact Archive Index

> Methodology: [`por-transparency-framework.md`](../por-transparency-framework.md)
> Comparison: [`exchange-comparison.md`](./exchange-comparison.md)

This index lists the report artifacts used by the first ardmere public transparency report set. Each report binds its conclusion to an assessment JSON, local artifact bundle, verification bundle, bundle roots, hashes, and timestamped verifier output.

## Report Set

| Exchange | Snapshot | Assessment JSON | Transparency Report | Supporting Audit Report |
| --- | --- | --- | --- | --- |
| OKX | `506872725` | [`506872725-assessment.json`](./okx/506872725-assessment.json) | [`506872725-transparency-report.md`](./okx/506872725-transparency-report.md) | [`AUDIT-REPORT.md`](../../artifacts/okx/506872725/AUDIT-REPORT.md) |
| Binance | `PR01JUN26` | [`PR01JUN26-assessment.json`](./binance/PR01JUN26-assessment.json) | [`PR01JUN26-transparency-report.md`](./binance/PR01JUN26-transparency-report.md) | [`AUDIT-REPORT.md`](../../artifacts/binance/PR01JUN26-20260601-period43/AUDIT-REPORT.md) |
| Bybit | `2025061709` | [`2025061709-assessment.json`](./bybit/2025061709-assessment.json) | [`2025061709-transparency-report.md`](./bybit/2025061709-transparency-report.md) | [`AUDIT-REPORT.md`](../../artifacts/bybit/2025061709/AUDIT-REPORT.md) |

## Bundle References

| Exchange | Artifact Bundle | Artifact Bundle Root | Artifact Bundle SHA-256 | Verification Bundle | Verification Bundle Root | Verification Bundle SHA-256 |
| --- | --- | --- | --- | --- | --- | --- |
| OKX | [`506872725.artifact-bundle.json`](../../artifacts/okx/506872725/bundles/506872725.artifact-bundle.json) | `0x8c6ce11cf467a4c74a971c8008d487be586ec9988ab8dde910c5cf110c312921` | `0ba21485473840cbc1f25a6f4906f7cb9d15eeb645c608f526d05c3a3f971023` | [`506872725.verification-bundle.v2.json`](../../artifacts/okx/506872725/bundles/506872725.verification-bundle.v2.json) | `0x71a08d13032475f58b2c306cfe2807ca2483912bc16be298b0bf8626a07b15c4` | `a0d627a67e0e2497dd788a8f01da5bf3fdbb82edd51eb8bc990269b4c6773976` |
| Binance | [`PR01JUN26.artifact-bundle.json`](../../artifacts/binance/PR01JUN26-20260601-period43/bundles/PR01JUN26.artifact-bundle.json) | `0xf452e47dd22ed63dc4a905fe79da6c6f7a6975cc0d775f50d01879e97616671f` | `ffc546892f40613f04434624c3ae14e3d5bb9f23e38bc971731c071b656ade50` | [`PR01JUN26.verification-bundle.v2.json`](../../artifacts/binance/PR01JUN26-20260601-period43/bundles/PR01JUN26.verification-bundle.v2.json) | `0x84e7f008460dcfa4ae4969c9a6e2a8b2585d5dcbe6a33cfead74b175782f6f42` | `b3b5a3dc2f08f198110131c787c52a666b3af5b5c3fbf4eddff4dff6e7c1e709` |
| Bybit | [`2025061709.artifact-bundle.json`](../../artifacts/bybit/2025061709/bundles/2025061709.artifact-bundle.json) | `0x7e2710dde862676f03cdde12e4590d451efca54b597049c5a732884c9fd9e69d` | `c00cc424d1218befb358b0fecb6046499dc6175b79b013079f99b347df989550` | [`2025061709.verification-bundle.json`](../../artifacts/bybit/2025061709/bundles/2025061709.verification-bundle.json) | `0xd9079261cfa47c44357e30119fdf95a9d077fda29a2a2026cc4a1bd067872c21` | `7d7f949786a4c1d48be58b10858a32e51a7930b2bf71102de2c2d9b1291c83d6` |

## Evidence Boundaries

- The archive index is a pointer map, not a substitute for verifier results. Read the exchange transparency report and assessment JSON for the Stage decision.
- Third-party ardmere bundle roots do not replace an official exchange-controlled canonical on-chain or DA anchor.
- Missing artifacts remain `UNVERIFIABLE`; they are not inferred from public summaries, screenshots, or marketing pages.
- Public ratings must bind to artifact hash, source URL or local path, verifier output, and timestamp.
