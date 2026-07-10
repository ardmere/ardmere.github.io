package verifier

import (
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
)

func TestLedgerSupportedEntriesValid(t *testing.T) {
	for key, spec := range loadLedgerSupported() {
		if !strings.Contains(key, "|") {
			t.Fatalf("bad key %q", key)
		}
		if spec.Decimals <= 0 {
			t.Fatalf("%s: bad decimals %d", key, spec.Decimals)
		}
		switch spec.Kind {
		case onchainconfig.LedgerEsplora:
			if spec.EsploraBase == "" {
				t.Fatalf("%s: missing esploraBase", key)
			}
		case onchainconfig.LedgerAlchemy:
			if spec.Alchemy == "" {
				t.Fatalf("%s: missing alchemy chain", key)
			}
		case onchainconfig.LedgerBlockchair:
			if spec.Blockchair == "" {
				t.Fatalf("%s: missing blockchair chain", key)
			}
		case onchainconfig.LedgerBlockcypher:
			if spec.Blockcypher == "" {
				t.Fatalf("%s: missing blockcypher chain", key)
			}
		case onchainconfig.LedgerSolanaSPL, onchainconfig.LedgerAptosFA, onchainconfig.LedgerAptosCoin,
			onchainconfig.LedgerNearFT, onchainconfig.LedgerSuiCoin, onchainconfig.LedgerHbarHTS,
			onchainconfig.LedgerTonJetton, onchainconfig.LedgerTezosFA:
			if spec.Mint == "" {
				t.Fatalf("%s: missing mint/metadata", key)
			}
		case onchainconfig.LedgerChromia:
			if spec.Mint == "" || spec.RPCURL == "" {
				t.Fatalf("%s: missing mint (blockchain RID) or rpcUrl", key)
			}
		case onchainconfig.LedgerStellarAsset:
			if spec.Mint == "" || spec.Issuer == "" {
				t.Fatalf("%s: missing mint or issuer", key)
			}
		case onchainconfig.LedgerXRPLEmitted:
			if spec.Mint == "" || spec.Issuer == "" {
				t.Fatalf("%s: missing mint or issuer", key)
			}
		case onchainconfig.LedgerAlgoASA:
			if spec.AssetID <= 0 {
				t.Fatalf("%s: missing assetId", key)
			}
		case onchainconfig.LedgerSubstrateAsset:
			if spec.AssetID <= 0 || spec.RPCURL == "" {
				t.Fatalf("%s: missing assetId or rpcUrl", key)
			}
		case onchainconfig.LedgerSubstrate:
			if spec.RPCURL == "" {
				t.Fatalf("%s: missing rpcUrl", key)
			}
		case onchainconfig.LedgerSolanaNative, onchainconfig.LedgerXRPL, onchainconfig.LedgerAptos,
			onchainconfig.LedgerNear, onchainconfig.LedgerSui, onchainconfig.LedgerHbar:
		default:
			t.Fatalf("%s: unknown kind %q", key, spec.Kind)
		}
	}
	if len(loadLedgerSupported()) < 30 {
		t.Fatalf("expected expanded ledger map, got %d", len(loadLedgerSupported()))
	}
}
