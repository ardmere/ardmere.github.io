// Package exchangeregistry loads config/exchanges/registry.yaml (official PoR upstream sources).
package exchangeregistry

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

const defaultPath = "config/exchanges/registry.yaml"

// UpstreamRepo is an official exchange PoR GitHub repo or release source.
type UpstreamRepo struct {
	URL        string `yaml:"url"`
	Role       string `yaml:"role"`
	ArdmereUse string `yaml:"ardmereUse,omitempty"`
	Note       string `yaml:"note,omitempty"`
}

// Exchange is one integrated exchange entry.
type Exchange struct {
	ID               string         `yaml:"id"`
	Name             string         `yaml:"name"`
	PorPage          string         `yaml:"porPage"`
	DataGuide        string         `yaml:"dataGuide"`
	VCTier           int            `yaml:"vcTier,omitempty"`
	UpstreamRepos    []UpstreamRepo `yaml:"upstreamRepos"`
	IntegrationNotes string         `yaml:"integrationNotes,omitempty"`
}

// Registry is the full upstream registry document.
type Registry struct {
	Version   int        `yaml:"version"`
	Updated   string     `yaml:"updated"`
	Exchanges []Exchange `yaml:"exchanges"`
}

// DefaultPath returns the canonical registry file path (relative to repo root).
func DefaultPath() string { return defaultPath }

// Load reads and parses the registry YAML at path.
func Load(path string) (Registry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Registry{}, err
	}
	var reg Registry
	if err := yaml.Unmarshal(raw, &reg); err != nil {
		return Registry{}, fmt.Errorf("parse registry %s: %w", path, err)
	}
	if len(reg.Exchanges) == 0 {
		return Registry{}, fmt.Errorf("registry %s: no exchanges", path)
	}
	return reg, nil
}

// Find returns the exchange with the given id, or nil.
func (r Registry) Find(id string) *Exchange {
	want := strings.ToLower(strings.TrimSpace(id))
	for i := range r.Exchanges {
		if strings.ToLower(r.Exchanges[i].ID) == want {
			return &r.Exchanges[i]
		}
	}
	return nil
}

// PrintTable writes a human-readable summary to w.
func PrintTable(reg Registry, w *tabwriter.Writer) {
	fmt.Fprintf(w, "Exchange upstream registry\t(updated %s)\n\n", reg.Updated)
	fmt.Fprintf(w, "ID\tVC\tPoR page\tUpstream repos\n")
	for _, ex := range reg.Exchanges {
		repos := make([]string, 0, len(ex.UpstreamRepos))
		for _, repo := range ex.UpstreamRepos {
			repos = append(repos, shortenRepoURL(repo.URL))
		}
		if len(repos) == 0 {
			repos = append(repos, "—")
		}
		fmt.Fprintf(w, "%s\tVC%d\t%s\t%s\n", ex.ID, ex.VCTier, ex.PorPage, strings.Join(repos, ", "))
	}
	_ = w.Flush()
}

// PrintDetail writes one exchange entry.
func PrintDetail(ex Exchange) {
	fmt.Printf("%s (%s)\n", ex.Name, ex.ID)
	fmt.Printf("  PoR page:   %s\n", ex.PorPage)
	fmt.Printf("  Data guide: %s\n", ex.DataGuide)
	if ex.VCTier > 0 {
		fmt.Printf("  VC tier:    VC%d\n", ex.VCTier)
	}
	if ex.IntegrationNotes != "" {
		fmt.Printf("  Notes:      %s\n", ex.IntegrationNotes)
	}
	for _, repo := range ex.UpstreamRepos {
		fmt.Printf("  upstream %s:\n", shortenRepoURL(repo.URL))
		fmt.Printf("    role:       %s\n", repo.Role)
		if repo.ArdmereUse != "" {
			fmt.Printf("    ardmereUse: %s\n", repo.ArdmereUse)
		}
		if repo.Note != "" {
			fmt.Printf("    note:       %s\n", repo.Note)
		}
	}
}

func shortenRepoURL(u string) string {
	const prefix = "https://github.com/"
	if strings.HasPrefix(u, prefix) {
		return u[len(prefix):]
	}
	return u
}
