package registry

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const officialBase = "https://registry.modelcontextprotocol.io/v0.1"
const officialMetaKey = "io.modelcontextprotocol.registry/official"
const officialMaxPages = 5

type OfficialServerEntry struct {
	Name        string
	FullName    string
	Description string
	Publisher   string
	Verified    bool
	Source      string // always "official"
}

// SearchOfficial searches the official MCP registry for servers matching query.
// Returns an empty slice (not an error) on network failure, matching Python behaviour.
func SearchOfficial(query string) ([]OfficialServerEntry, error) {
	var results []OfficialServerEntry
	cursor := ""
	q := strings.ToLower(query)

	for range officialMaxPages {
		servers, next, err := fetchOfficialPage(cursor)
		if err != nil {
			return results, nil // silent on connection errors, per Python
		}

		for _, item := range servers {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			meta, _ := itemMap["_meta"].(map[string]interface{})
			official, _ := meta[officialMetaKey].(map[string]interface{})
			if status, _ := official["status"].(string); status != "active" {
				continue
			}

			server, _ := itemMap["server"].(map[string]interface{})
			name, _ := server["name"].(string)
			desc, _ := server["description"].(string)

			if !strings.Contains(strings.ToLower(name), q) &&
				!strings.Contains(strings.ToLower(desc), q) {
				continue
			}

			shortName, publisher := splitName(name)
			results = append(results, OfficialServerEntry{
				Name:        shortName,
				FullName:    name,
				Description: desc,
				Publisher:   publisher,
				Verified:    false,
				Source:      "official",
			})
		}

		if next == "" {
			break
		}
		cursor = next
	}

	return results, nil
}

// FetchOfficialManifest fetches and converts a manifest from the official MCP registry.
// Returns (nil, nil) on 404 or connection error.
func FetchOfficialManifest(serverName string) (*Manifest, error) {
	encoded := url.PathEscape(serverName)
	reqURL := fmt.Sprintf("%s/servers/%s/versions/latest", officialBase, encoded)

	resp, err := httpClient.Get(reqURL)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil
	}

	server, _ := data["server"].(map[string]interface{})
	packages, _ := server["packages"].([]interface{})
	if len(packages) == 0 {
		return nil, nil
	}

	return ConvertToArcmeshManifest(server), nil
}

// ConvertToArcmeshManifest converts the official MCP registry server JSON
// to an ArcMesh Manifest. Returns nil if the server has no packages.
func ConvertToArcmeshManifest(serverJSON map[string]interface{}) *Manifest {
	packages, _ := serverJSON["packages"].([]interface{})
	if len(packages) == 0 {
		return nil
	}

	pkg, _ := packages[0].(map[string]interface{})

	name, _ := serverJSON["name"].(string)
	shortName, publisherName := splitName(name)

	registryType, _ := pkg["registryType"].(string)
	runtimeHint, _ := pkg["runtimeHint"].(string)
	identifier, _ := pkg["identifier"].(string)

	var installType string
	switch {
	case registryType == "npm":
		installType = "npx"
	case registryType == "pypi" && runtimeHint == "uvx":
		installType = "uvx"
	case registryType == "pypi":
		installType = "pip"
	default:
		installType = "npx"
	}

	var command string
	var args []string
	switch installType {
	case "npx":
		command, args = "npx", []string{"-y", identifier}
	case "uvx":
		command, args = "uvx", []string{identifier}
	default: // pip
		command, args = "python", []string{"-m", shortName}
	}

	env := buildEnv(pkg)
	cfg := ManifestConfig{Command: command, Args: args}
	if len(env) > 0 {
		cfg.Env = env
	}

	version, _ := serverJSON["version"].(string)
	if version == "" {
		version = "latest"
	}
	desc, _ := serverJSON["description"].(string)

	var sourceURL string
	if repo, ok := serverJSON["repository"].(map[string]interface{}); ok {
		sourceURL, _ = repo["url"].(string)
	}

	return &Manifest{
		Name:        shortName,
		Version:     version,
		Description: desc,
		Publisher:   ManifestPublisher{Name: publisherName, Verified: false},
		SourceURL:   sourceURL,
		Install:     ManifestInstall{Type: installType, Package: identifier},
		Config:      cfg,
	}
}

// fetchOfficialPage fetches one page from the official registry.
// Extracted to avoid defer inside a loop.
func fetchOfficialPage(cursor string) (servers []interface{}, nextCursor string, err error) {
	params := url.Values{"limit": {"100"}}
	if cursor != "" {
		params.Set("cursor", cursor)
	}

	resp, err := httpClient.Get(officialBase + "/servers?" + params.Encode())
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("official registry returned %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, "", err
	}

	servers, _ = data["servers"].([]interface{})
	meta, _ := data["metadata"].(map[string]interface{})
	nextCursor, _ = meta["nextCursor"].(string)
	return servers, nextCursor, nil
}

func splitName(name string) (short, publisher string) {
	slash := strings.LastIndex(name, "/")
	if slash == -1 {
		return name, ""
	}
	return name[slash+1:], name[:slash]
}

func buildEnv(pkg map[string]interface{}) map[string]EnvSpec {
	args, _ := pkg["packageArguments"].([]interface{})
	env := make(map[string]EnvSpec, len(args))
	for _, a := range args {
		argMap, ok := a.(map[string]interface{})
		if !ok {
			continue
		}
		if argMap["type"] != "environment" {
			continue
		}
		key, _ := argMap["value"].(string)
		desc, _ := argMap["description"].(string)
		required, _ := argMap["isRequired"].(bool)
		secret, _ := argMap["isSecret"].(bool)
		env[key] = EnvSpec{Description: desc, Required: required, Secret: secret}
	}
	return env
}
