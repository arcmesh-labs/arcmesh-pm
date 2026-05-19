package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const RegistryBase = "https://raw.githubusercontent.com/arcmesh-labs/arcmesh-registry/master/servers"

var httpClient = &http.Client{Timeout: 10 * time.Second}

// ServerEntry is one item from the registry index.
type ServerEntry struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Publisher   string   `json:"publisher"`
	Verified    bool     `json:"verified"`
	Tags        []string `json:"tags"`
	Path        string   `json:"path"`
}

type ManifestPublisher struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Verified bool   `json:"verified"`
}

type ManifestInstall struct {
	Type    string `json:"type"`
	Package string `json:"package"`
	Version string `json:"version"`
}

type EnvSpec struct {
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Default     *string `json:"default"`
	Secret      bool    `json:"secret"`
}

type ManifestConfig struct {
	Command string             `json:"command"`
	Args    []string           `json:"args"`
	Env     map[string]EnvSpec `json:"env"`
}

type Manifest struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Publisher   ManifestPublisher `json:"publisher"`
	Install     ManifestInstall   `json:"install"`
	Config      ManifestConfig    `json:"config"`
	SourceURL   string            `json:"source_url"`
	Tags        []string          `json:"tags"`
	Clients     []string          `json:"clients"`
}

// FetchIndex fetches the full server list from the ArcMesh registry.
func FetchIndex() ([]ServerEntry, error) {
	resp, err := httpClient.Get(RegistryBase + "/index.json")
	if err != nil {
		return nil, fmt.Errorf("could not connect to registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("registry returned %d", resp.StatusCode)
	}

	var payload struct {
		Servers []ServerEntry `json:"servers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("invalid registry response: %w", err)
	}
	return payload.Servers, nil
}

// FetchManifest fetches a server manifest by name.
// Returns (nil, nil) when the server is not found (404).
// Returns (nil, err) on network or non-404 HTTP errors.
func FetchManifest(name string) (*Manifest, error) {
	url := fmt.Sprintf("%s/%s/manifest.json", RegistryBase, name)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not connect to registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("registry returned %d", resp.StatusCode)
	}

	var m Manifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}
	return &m, nil
}
