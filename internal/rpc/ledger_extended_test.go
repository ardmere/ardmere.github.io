package rpc

import "testing"

func TestParseDecimalToSmallest(t *testing.T) {
	v, err := parseDecimalToSmallest("13046845.244527", 6)
	if err != nil {
		t.Fatal(err)
	}
	if v != 13046845244527 {
		t.Fatalf("got %d", v)
	}
	v, err = parseDecimalToSmallest("27280247.0925453", 7)
	if err != nil {
		t.Fatal(err)
	}
	if v != 272802470925453 {
		t.Fatalf("got %d", v)
	}
}
