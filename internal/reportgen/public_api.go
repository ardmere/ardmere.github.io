package reportgen

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

// PageOnlyPoRSourceMessage explains why page/press-only PoR data is fragile.
// Used in risk flags, limitations, and report notes when publicPoRApi is null.
const PageOnlyPoRSourceMessage = "No documented public PoR API. Snapshot data depends on web pages and/or press materials; page scraping can break when resource paths migrate or the front end is a SPA that ships an HTML shell without usable totals."

type publicPoRApiInfo struct {
	// URL is the documented permissionless PoR HTTP API endpoint, or empty for null.
	URL   string
	Notes string
}

func (i publicPoRApiInfo) HasAPI() bool {
	return i.URL != ""
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

// lookupPublicPoRApi returns a documented PoR API URL from the exchange registry,
// or empty URL (assessment null) when none is documented. Does not invent endpoints.
func lookupPublicPoRApi(exchangeID string) publicPoRApiInfo {
	info := publicPoRApiInfo{}
	reg, err := loadExchangeRegistry()
	if err != nil {
		return info
	}
	ex := reg.Find(exchangeID)
	if ex == nil {
		return info
	}
	info.Notes = strings.TrimSpace(ex.PublicPoRApiNote)
	info.URL = normalizePublicPoRApiURL(ex.PublicPoRApi)
	if !info.HasAPI() && info.Notes == "" {
		info.Notes = PageOnlyPoRSourceMessage
	}
	return info
}

func normalizePublicPoRApiURL(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	lower := strings.ToLower(s)
	switch lower {
	case "null", "none", "no", "false", "unknown", "n/a", "-":
		return ""
	}
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return ""
	}
	return s
}

// publicPoRApiAssessmentValue is the assessment JSON value: URL string or null.
func publicPoRApiAssessmentValue(info publicPoRApiInfo) any {
	if !info.HasAPI() {
		return nil
	}
	return info.URL
}

func publicPoRApiDisplay(info publicPoRApiInfo) string {
	if !info.HasAPI() {
		return "null"
	}
	return info.URL
}

func noPublicPoRApiRiskFlag(info publicPoRApiInfo) (riskFlag, bool) {
	if info.HasAPI() {
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
