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
  if [[ "$(uname)" == "Darwin" ]]; then
    sed -i '' \
      -e 's/href="\([^"]*\)\.md"/href="\1.html"/g' \
      -e 's/href="\([^"]*\)\.md#/href="\1.html#/g' \
      -e 's/(\([^)]*\)\.md)/(\1.html)/g' \
      "$file"
  else
    sed -i \
      -e 's/href="\([^"]*\)\.md"/href="\1.html"/g' \
      -e 's/href="\([^"]*\)\.md#/href="\1.html#/g' \
      -e 's/(\([^)]*\)\.md)/(\1.html)/g' \
      "$file"
  fi
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
build_page "docs/por-transparency-framework.md" "methodology"
build_page "docs/effective-por-standard.md" "standard"
build_page "docs/ardmere-service-audience.md"
build_page "docs/exchange-reserve-transparency-whitepaper.md"
build_page "docs/reports/exchange-comparison.md" "reports"
build_page "docs/reports/snapshot-history.md" "reports"
build_page "docs/reports/artifact-archive-index.md" "reports"
for f in docs/reports/okx/*-transparency-report.md docs/reports/binance/*-transparency-report.md \
  docs/reports/bitget/*-transparency-report.md docs/reports/bybit/*-transparency-report.md \
  docs/reports/gateio/*-transparency-report.md docs/reports/htx/*-transparency-report.md; do
  [[ -f "$ROOT/$f" ]] && build_page "$f"
done
build_page "docs/insights/index.md" "insights"
build_page "docs/insights/why-audited-is-not-effective-por.md" "insights"
echo "Done."
