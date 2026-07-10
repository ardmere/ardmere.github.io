package rpc

import (
	"errors"
	"testing"
)

func TestIsNearGarbageCollected(t *testing.T) {
	err := errors.New("near ft: The data for block #200773913 is garbage collected on this node")
	if !isNearGarbageCollected(err) {
		t.Fatal("expected garbage collected detection")
	}
	if isNearGarbageCollected(nil) {
		t.Fatal("nil should be false")
	}
}

func TestParseAptosFABalanceBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    int64
		wantErr bool
	}{
		{name: "bare zero", body: "0", want: 0},
		{name: "bare number", body: "123456789", want: 123456789},
		{name: "object string", body: `{"balance":"987654321"}`, want: 987654321},
		{name: "object zero", body: `{"balance":"0"}`, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAptosFABalanceBody([]byte(tt.body))
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("got %d want %d", got, tt.want)
			}
		})
	}
}
