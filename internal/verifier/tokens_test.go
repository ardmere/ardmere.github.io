package verifier

import (
	"strings"
	"testing"

	"github.com/ardmere/ardmere/internal/rpc"
)

func TestTokenSupportedEntriesValid(t *testing.T) {
	knownNets := map[rpc.Network]bool{
		rpc.NetEthereum:  true,
		rpc.NetBSC:       true,
		rpc.NetArbitrum:  true,
		rpc.NetOptimism:  true,
		rpc.NetBase:      true,
		rpc.NetPolygon:   true,
		rpc.NetAvalanche: true,
		rpc.NetOpBNB:     true,
		rpc.NetSonic:     true,
		rpc.NetWorld:     true,
		rpc.NetTron:      true,
		rpc.NetCelo:      true,
		rpc.NetPlasma:    true,
		rpc.NetKaia:      true,
		rpc.NetKavaEVM:   true,
		rpc.NetZkSync:    true,
		rpc.NetRon:       true,
		rpc.NetSeiEVM:    true,
		rpc.NetManta:     true,
		rpc.NetXLayer:    true,
		rpc.NetFEVM:      true,
		rpc.Network("LINEA"):    true,
		rpc.Network("SCROLL"):   true,
		rpc.Network("CHZ2"):     true,
		rpc.Network("MTL"):      true,
		rpc.Network("STARKNET"): true,
		rpc.NetAB:               true,
	}
	for key, spec := range loadTokenSupported() {
		if !strings.Contains(key, "|") {
			t.Fatalf("bad key %q", key)
		}
		if !knownNets[spec.Net] {
			t.Fatalf("%s: unknown network %q", key, spec.Net)
		}
		if spec.Native {
			if spec.Decimals != 18 {
				t.Fatalf("%s: native spec wants 18 decimals, got %d", key, spec.Decimals)
			}
			continue
		}
		if spec.Net == rpc.NetTron {
			if !strings.HasPrefix(spec.Contract, "T") || len(spec.Contract) < 30 {
				t.Fatalf("%s: bad tron contract %q", key, spec.Contract)
			}
			continue
		}
		if spec.Net == rpc.Network("STARKNET") {
			if !strings.HasPrefix(spec.Contract, "0x") || len(spec.Contract) < 10 {
				t.Fatalf("%s: bad starknet contract %q", key, spec.Contract)
			}
			continue
		}
		if !strings.HasPrefix(spec.Contract, "0x") || len(spec.Contract) != 42 {
			t.Fatalf("%s: bad contract %q", key, spec.Contract)
		}
	}
	if len(loadTokenSupported()) < 120 {
		t.Fatalf("expected expanded token map, got %d entries", len(loadTokenSupported()))
	}
}
