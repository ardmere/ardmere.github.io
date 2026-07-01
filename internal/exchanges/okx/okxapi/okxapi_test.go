package okxapi

import (
	"testing"
)

func TestParseSummaryBytes(t *testing.T) {
	raw := []byte(`{
	  "fetchedAt":"2026-06-14T00:00:00Z",
	  "source":"fixture",
	  "audit":{"auditId":"506872725","merkleHash":"abc","createTime":1778083200000}
	}`)
	bundle, err := ParseSummaryBytes(raw)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Audit.AuditID != "506872725" {
		t.Fatalf("audit: %+v", bundle.Audit)
	}
}
