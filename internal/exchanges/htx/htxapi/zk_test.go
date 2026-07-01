package htxapi

import (
	"os"
	"testing"
)

func TestParseZkBundleZipSample(t *testing.T) {
	path := os.Getenv("HTX_PUBLIC_DATA_ZIP")
	if path == "" {
		path = "/tmp/htx-pub/public-data.zip"
	}
	if _, err := os.Stat(path); err != nil {
		t.Skipf("public-data.zip not available: %v", err)
	}
	meta, err := ParseZkBundleZip(path)
	if err != nil {
		t.Fatal(err)
	}
	wantRoot := "2c51a40e991be73032f034c61b5d6a620c6da1fc9a7e0fdf104ca396e45d9ee3"
	if meta.MerkleRoot != wantRoot {
		t.Fatalf("merkle root: got %s want %s", meta.MerkleRoot, wantRoot)
	}
	if meta.AuditID != "20230910" {
		t.Fatalf("audit id: got %s", meta.AuditID)
	}
	if meta.BatchCount < 1000 {
		t.Fatalf("expected many batches, got %d", meta.BatchCount)
	}
	if len(meta.Config.CexAssetsInfo) < 100 {
		t.Fatalf("expected 174 assets, got %d", len(meta.Config.CexAssetsInfo))
	}
}
