package gateapi

import (
	"encoding/json"
	"os"
	"time"
)

// BuildSummaryBundle merges browser-captured API JSON files into one archive blob.
func BuildSummaryBundle(infoPath, coinsPath, listPath string) ([]byte, error) {
	info, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, err
	}
	bundle := SummaryBundle{
		FetchedAt: time.Now().UTC(),
		Info:      json.RawMessage(info),
	}
	if coinsPath != "" {
		raw, err := os.ReadFile(coinsPath)
		if err != nil {
			return nil, err
		}
		bundle.CoinList = json.RawMessage(raw)
	}
	if listPath != "" {
		raw, err := os.ReadFile(listPath)
		if err != nil {
			return nil, err
		}
		bundle.List = json.RawMessage(raw)
	}
	return json.Marshal(bundle)
}
