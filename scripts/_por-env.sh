# shellcheck shell=bash
# Load operator RPC/API keys for por verify, anchor, and batch.
# See docs/ledger-rpc-runbook.md and docs/por-cli.md.

if [[ -f "${HOME}/.zshenv" ]]; then
	set -a
	# shellcheck disable=SC1091
	source "${HOME}/.zshenv"
	set +a
fi

if [[ -f .env ]]; then
	set -a
	# shellcheck disable=SC1091
	source .env
	set +a
fi
