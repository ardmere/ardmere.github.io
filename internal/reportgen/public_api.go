package reportgen

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

// PageOnlyPoRSourceMessage explains why page/press-only PoR data is fragile.
// Used in risk flags, limitations, and report notes when publicPoRApi.provided is "no".
const PageOnlyPoRSourceMessage = "No documented public PoR API. Snapshot data depends on web pages and/or press materials; page scraping can break when resource paths migrate or the front end is a SPA that ships an HTML shell without usable totals."

type publicPoRApiInfo struct {
	Provided string
	Notes    string
}

var (
	registryOnce sync.Once
	registry     exchangeregistry.Registry
	registryErr  error
)

func loadExchangeRegistry() (exchangeregistry.Registry, error) {
	registryOnce.Do(func() {
		var lastErr error
		for _, path := range candidateRegistryPaths() {
			reg, err := exchangeregistry.Load(path)
			if err == nil {
				registry = reg
				registryErr = nil
				return
			}
			lastErr = err
		}
		registryErr = lastErr
	})
	return registry, registryErr
}

func candidateRegistryPaths() []string {
	rel := exchangeregistry.DefaultPath()
	seen := map[string]bool{}
	var out []string
	add := func(p string) {
		if p == "" || seen[p] {
			return
		}
		seen[p] = true
		out = append(out, p)
	}
	add(rel)
	if wd, err := os.Getwd(); err == nil {
		dir := wd
		for i := 0; i < 8; i++ {
			add(filepath.Join(dir, rel))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return out
}

// lookupPublicPoRApi returns yes/no/unknown from the exchange registry.
// Missing or invalid registry values default to unknown (do not invent "yes").
func lookupPublicPoRApi(exchangeID string) publicPoRApiInfo {
	info := publicPoRApiInfo{Provided: "unknown"}
	reg, err := loadExchangeRegistry()
	if err != nil {
		return info
	}
	ex := reg.Find(exchangeID)
	if ex == nil {
		return info
	}
	switch strings.ToLower(strings.TrimSpace(ex.PublicPoRApi)) {
	case "yes", "true":
		// "true" only if a YAML bool was stringified; prefer quoted yes/no/unknown in registry.
		info.Provided = "yes"
	case "no", "false":
		info.Provided = "no"
	case "unknown":
		info.Provided = "unknown"
	default:
		info.Provided = "unknown"
	}
	info.Notes = strings.TrimSpace(ex.PublicPoRApiNote)
	if info.Provided == "no" && info.Notes == "" {
		info.Notes = PageOnlyPoRSourceMessage
	}
	return info
}

func publicPoRApiAssessmentMap(info publicPoRApiInfo) map[string]any {
	out := map[string]any{"provided": info.Provided}
	if info.Notes != "" {
		out["notes"] = info.Notes
	} else if info.Provided == "no" {
		out["notes"] = PageOnlyPoRSourceMessage
	}
	return out
}

func noPublicPoRApiRiskFlag(info publicPoRApiInfo) (riskFlag, bool) {
	if info.Provided != "no" {
		return riskFlag{}, false
	}
	msg := info.Notes
	if msg == "" {
		msg = PageOnlyPoRSourceMessage
	}
	return riskFlag{
		ID:          "NO_PUBLIC_POR_API",
		Severity:    "medium",
		StageEffect: "none",
		Message:     msg,
	}, true
}
