# ardmere Verifier Coverage Tiers

Each registered exchange declares a **Verifier Coverage Tier (VC)** via `Adapter.Capabilities()`. VC tiers set expectations for what ardmere can independently verify today; they do not change verdict semantics (`PASS` / `FAIL` / `UNVERIFIABLE`).

For exchange-level transparency ratings, use [`por-transparency-framework.md`](./por-transparency-framework.md). Do not conflate Transparency Tier (T0-T5) with ardmere implementation coverage (VC1-VC3).

| VC | Name | Exchanges | What we can verify independently |
|----|------|-----------|----------------------------------|
| **VC1** | Summary / user proof | Gate.io, **HTX**, **Bybit**, **Bitget** | Artifact integrity, self-reported solvency ratios; HTX adds `global-zk-proof@htx-1` when zk bundle present; Bybit/Bitget add user Merkle proof verifiers when login-gated proof JSON is imported |
| **VC2** | Wallet + sig / zk | OKX | VC1 + internal consistency, address ownership, global-zk structure; partial onchain with RPC |
| **VC3** | Full onchain | Binance | VC2-like surface where available + deep multi-chain hot/token/ledger/deposit onchain audit |

## Capability flags

```go
type Capabilities struct {
    Tier             int  // VC1 | VC2 | VC3, encoded as 1 | 2 | 3
    WalletZip        bool
    AddressOwnership bool
    GlobalZK         bool
    OnchainHot       bool
    OnchainToken     bool
    OnchainLedger    bool
}
```

Inspect at runtime:

```bash
go run ./cmd/por anchor -exchange okx -skip-zip   # logs VC tier in [1/5] line
```

## Adding an exchange

Pick the VC tier that matches **ardmere's implemented verifier coverage**, then implement only the verifiers your profile activates. Separately assign the exchange's Transparency Tier in audit reports using [`por-transparency-framework.md`](./por-transparency-framework.md).
