package registry

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	fallbackDimStyle  = lipgloss.NewStyle().Faint(true)
	fallbackWarnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	fallbackBoldStyle = lipgloss.NewStyle().Bold(true)
)

// SearchResult is a unified search hit from either registry.
type SearchResult struct {
	Name        string
	Description string
	Publisher   string
	Verified    bool
	Source      string // "arcmesh" or "official"
}

// FetchManifestWithFallback fetches a manifest, trying ArcMesh first then the
// official MCP Registry on 404. Prints the same informational messages as the
// Python version. Returns an error (not an exit) on fatal conditions.
func FetchManifestWithFallback(name string) (*Manifest, string, error) {
	manifest, err := FetchManifest(name)
	if err != nil {
		return nil, "", err
	}
	if manifest != nil {
		return manifest, "arcmesh", nil
	}

	// 404 from ArcMesh — try official registry
	fmt.Println(fallbackDimStyle.Render("Not found in ArcMesh registry. Trying official MCP Registry..."))

	official, err := FetchOfficialManifest(name)
	if err != nil {
		return nil, "", err
	}
	if official != nil {
		fmt.Printf("%s  %s not found in ArcMesh registry.\n"+
			"   Installing from official MCP Registry — manifest converted automatically.\n"+
			"   Review config before restarting your client.\n",
			fallbackWarnStyle.Render("⚠"),
			"'"+fallbackBoldStyle.Render(name)+"'")
		return official, "official", nil
	}

	return nil, "", fmt.Errorf("server '%s' not found in ArcMesh or official MCP Registry", name)
}

// SearchWithFallback searches ArcMesh first; falls back to the official MCP
// Registry when there are no ArcMesh matches. Prints a dim status line when
// falling back (matching Python output). The caller is responsible for rendering
// the returned results.
func SearchWithFallback(query string) ([]SearchResult, error) {
	index, err := FetchIndex()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var arcmeshHits []SearchResult
	for _, s := range index {
		if strings.Contains(strings.ToLower(s.Name), q) ||
			strings.Contains(strings.ToLower(s.Description), q) {
			arcmeshHits = append(arcmeshHits, SearchResult{
				Name:        s.Name,
				Description: s.Description,
				Publisher:   s.Publisher,
				Verified:    s.Verified,
				Source:      "arcmesh",
			})
		}
	}

	if len(arcmeshHits) > 0 {
		return arcmeshHits, nil
	}

	fmt.Println(fallbackDimStyle.Render("No results in ArcMesh registry. Searching official MCP Registry..."))

	official, err := SearchOfficial(query)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(official))
	for i, s := range official {
		results[i] = SearchResult{
			Name:        s.Name,
			Description: s.Description,
			Publisher:   s.Publisher,
			Verified:    false,
			Source:      "official",
		}
	}
	return results, nil
}
