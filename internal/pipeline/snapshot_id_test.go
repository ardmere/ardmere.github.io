package pipeline

import (
	"encoding/hex"
	"testing"
)

func TestCompositeSnapshotID_collisionRegression(t *testing.T) {
	audit := "20261101"
	want := map[string]string{
		"phemex": "12abd55db2e6d9b44edf24364ab8735226d612d54b228074210cb4933ec39f80",
		"htx":    "bdbd096d9043bf0b55246ef612e4db87ac27eb5718e8debd15fdd77db30e6f20",
		"ourbit": "cf880c45c95eb72602f327c669c070e2cf047bb236b57614ed774c1fca9e6412",
	}
	ids := make(map[[32]byte]string)
	for ex, hexWant := range want {
		got := CompositeSnapshotID(ex, audit)
		if hex.EncodeToString(got[:]) != hexWant {
			t.Fatalf("%s: got %x want %s", ex, got, hexWant)
		}
		if prev, ok := ids[got]; ok {
			t.Fatalf("collision between %s and %s", prev, ex)
		}
		ids[got] = ex
	}
	legacy := LegacySnapshotID(audit)
	if _, ok := ids[legacy]; ok {
		t.Fatal("legacy id must differ from composite ids")
	}
}

func TestCompositeSnapshotID_normalizesExchangeCase(t *testing.T) {
	lo := CompositeSnapshotID("binance", "PR01JUN26")
	hi := CompositeSnapshotID("Binance", "PR01JUN26")
	if lo != hi {
		t.Fatal("exchange slug should be normalized to lowercase")
	}
}
