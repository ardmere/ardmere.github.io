# Gate PoR browser-capture fixtures (2026-03-16 audit)

Public dashboard fields from https://www.gate.com/proof-of-reserves (2026-06-13).

These files mirror the Gate Web API response shape. The coin list includes only
the six assets shown on the landing page; the live API returns ~83 coins.

When Akamai blocks datacenter API access, run:

```bash
./scripts/gateio/gate-save-local.sh
```

For a full coin list, capture `getProofOfReservesCoinList` in browser DevTools and:

```bash
./scripts/gateio/gate-import-browser.sh ./info.json ./coinList.json
```
