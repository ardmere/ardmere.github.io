#!/usr/bin/env bash
# Portable temp path ending in .json (BSD mktemp requires XXXXXX at end of template).
por_mktemp_json() {
  local prefix="${1:-por-temp}"
  local _tmp
  _tmp="$(mktemp "${TMPDIR:-/tmp}/${prefix}-XXXXXX")"
  local path="${_tmp}.json"
  mv "$_tmp" "$path"
  printf '%s' "$path"
}
