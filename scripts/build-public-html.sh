#!/usr/bin/env bash
# Generate styled HTML pages from public-site markdown sources.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TEMPLATE="$ROOT/scripts/site/doc-template.html"

if ! command -v pandoc >/dev/null 2>&1; then
  echo "pandoc is required to build public HTML pages" >&2
  exit 1
fi

fix_links() {
  local file="$1"
  # Rewrite relative *.md links to *.html for GitHub Pages.
  # Skip absolute URLs (http:, https:, mailto:, etc.).
  python3 - "$file" <<'PY'
import re, sys
path = sys.argv[1]
text = open(path, encoding="utf-8").read()

def repl(m: re.Match[str]) -> str:
    href = m.group(1)
    if re.match(r"^[a-zA-Z][a-zA-Z0-9+.-]*:", href):
        return m.group(0)
    if href.endswith(".md"):
        return f'href="{href[:-3]}.html"'
    if ".md#" in href:
        return f'href="{href.replace(".md#", ".html#", 1)}"'
    return m.group(0)

text2 = re.sub(r'href="([^"]+)"', repl, text)
if text2 != text:
    open(path, "w", encoding="utf-8").write(text2)
PY
}

build_page() {
  local md_rel="$1"
  local nav_key="${2:-}"
  local html_rel="${md_rel%.md}.html"
  local md="$ROOT/$md_rel"
  local html="$ROOT/$html_rel"
  local title
  title="$(grep -m1 '^# ' "$md" | sed 's/^# //')"

  if [[ -n "$nav_key" ]]; then
    pandoc "$md" \
      --standalone \
      --template="$TEMPLATE" \
      --from markdown \
      --to html5 \
      --metadata "title=$title" \
      --metadata "nav-$nav_key=true" \
      -o "$html"
  else
    pandoc "$md" \
      --standalone \
      --template="$TEMPLATE" \
      --from markdown \
      --to html5 \
      --metadata "title=$title" \
      -o "$html"
  fi

  fix_links "$html"
  echo "  $html_rel"
}

echo "Building public HTML pages..."
go run ./cmd/por site
build_page "docs/por-transparency-framework.md" "methodology"
build_page "docs/effective-por-standard.md" "standard"
build_page "docs/contract-query-guide.md"
build_page "docs/ardmere-service-audience.md"
build_page "docs/exchange-reserve-transparency-whitepaper.md"
build_page "docs/reports/exchange-comparison.md" "reports"
build_page "docs/reports/snapshot-history.md" "reports"
build_page "docs/reports/artifact-archive-index.md" "reports"
for f in docs/reports/okx/*-transparency-report.md docs/reports/binance/*-transparency-report.md \
  docs/reports/bitget/*-transparency-report.md docs/reports/bybit/*-transparency-report.md \
  docs/reports/gateio/*-transparency-report.md docs/reports/htx/*-transparency-report.md \
  docs/reports/mexc/*-transparency-report.md docs/reports/kucoin/*-transparency-report.md \
  docs/reports/bingx/*-transparency-report.md; do
  [[ -f "$ROOT/$f" ]] && build_page "$f"
done
build_page "docs/insights/index.md" "insights"
build_page "docs/insights/why-audited-is-not-effective-por.md" "insights"
echo "Done."
