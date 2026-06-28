# `por` CLI

Unified Proof-of-Reserves command-line tool. Implementation lives in `internal/porcli`, `internal/porrun`, `internal/porfetch`, `internal/porprobe`.

## Commands

| Command | Purpose |
|---------|---------|
| `por anchor` | fetch → verify → bundles → anchor calldata |
| `por verify` | re-verify cached artifacts |
| `por fetch gateio\|okx` | archive-only download |
| `por probe rpc\|tron\|…` | Binance regression / RPC diagnostics |

## Build & run

```bash
go build -o por ./cmd/por
./por anchor -exchange binance
./por anchor -exchange binance --skip-zip
./por verify -snapshot PR01JUN26
./por fetch gateio
./por probe rpc -network BSC -chainlist
```

Shell wrapper: `./scripts/por.sh anchor -exchange okx`

Legacy `./scripts/por-run.sh` and `./scripts/por-fetch.sh` forward to the same binary.

## Pipeline (anchor)

See [verifier-architecture.md](./verifier-architecture.md) for verifier matrix and verdict semantics.

## Trust boundary

- Does not require a private key
- Does not broadcast transactions (prints `cast send` for the operator)
- Uses public RPC providers only
