package verifier

import (
	"encoding/json"
	"os"
	"testing"
)

func TestOKXMerkleRootFromPublicInputs506872725(t *testing.T) {
	path := os.Getenv("OKX_SUM_PROOF_PATH")
	if path == "" {
		path = "../../artifacts/okx/506872725/raw/319a01e352db65d17c3a6fd32b89ef35501c50e9a6b8714608c65112c36eeaee.zip"
	}
	raw, err := readSumProofData(path)
	if err != nil {
		t.Skipf("sum proof not available: %v", err)
	}
	var doc struct {
		Proof okxProofEnvelope `json:"proof"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	agg, err := okxZKAggregatesFromProof(doc.Proof)
	if err != nil {
		t.Fatal(err)
	}
	want := "dd32f9c29d1840dc444b0c0780782cd04f6840870b934a0a8f607fc61b3e8d6b"
	if agg.MerkleHex != want {
		t.Fatalf("merkle hex: got %s want %s", agg.MerkleHex, want)
	}
	if agg.Equity != 33938289499439 || agg.Debt != 2058816990077 {
		t.Fatalf("unexpected equity/debt: %+v", agg)
	}
	if agg.Liability != 31879472509362 {
		t.Fatalf("liability: got %d", agg.Liability)
	}
}

func TestMerkleHexEqual(t *testing.T) {
	if !merkleHexEqual("0xAbC", "abc") {
		t.Fatal("expected equal")
	}
}

func TestOKXProofRoundString(t *testing.T) {
	got := okxProofRoundString(map[string]any{"round_num": float64(506872725)})
	if got != "506872725" {
		t.Fatalf("got %q", got)
	}
}
