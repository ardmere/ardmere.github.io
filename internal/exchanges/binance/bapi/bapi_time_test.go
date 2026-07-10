package bapi

import (
	"testing"
	"time"
)

func TestParseSnapshotTime_DDMMYY(t *testing.T) {
	tests := []struct {
		in   string
		want time.Time
	}{
		{"01/04/26 00:00:00 UTC", time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)},
		{"01/05/26 00:00:00 UTC", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)},
		{"01/06/26 00:00:00 UTC", time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tc := range tests {
		got, err := ParseSnapshotTime(tc.in)
		if err != nil {
			t.Fatalf("ParseSnapshotTime(%q): %v", tc.in, err)
		}
		if !got.Equal(tc.want) {
			t.Fatalf("ParseSnapshotTime(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
