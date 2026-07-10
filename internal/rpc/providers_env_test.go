package rpc

import (
	"path/filepath"
	"testing"
)

func TestLoadProviderConfigInfuraAlchemyFromEnv(t *testing.T) {
	t.Setenv("INFURA_KEY", "test-infura")
	t.Setenv("ALCHEMY_KEY", "test-alchemy")
	t.Setenv("RPC_PROVIDERS_CONFIG", filepath.Join("..", "..", "config", "rpc-providers.json"))
	cfg, err := LoadProviderConfig()
	if err != nil {
		t.Fatal(err)
	}
	eth := cfg[NetEthereum]
	var infura, alchemy bool
	for _, p := range eth {
		switch p.URL {
		case "https://mainnet.infura.io/v3/test-infura":
			infura = true
		case "https://eth-mainnet.g.alchemy.com/v2/test-alchemy":
			alchemy = true
		}
	}
	if !infura {
		t.Fatalf("infura provider missing from ETH list: %+v", eth)
	}
	if !alchemy {
		t.Fatalf("alchemy provider missing from ETH list: %+v", eth)
	}
}
