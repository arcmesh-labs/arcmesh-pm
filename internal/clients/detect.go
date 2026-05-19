package clients

import (
	"encoding/json"
	"os"
)

// DetectConfiguredClients returns clients that have an existing config file.
func DetectConfiguredClients() []string {
	var result []string
	for _, client := range ClientOrder {
		if _, found := FindConfig(client); found {
			result = append(result, client)
		}
	}
	return result
}

// DetectActiveClients returns clients that have at least one server entry.
func DetectActiveClients() []string {
	var active []string
	for _, client := range DetectConfiguredClients() {
		path, ok := FindConfig(client)
		if !ok {
			continue
		}
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var cfg map[string]interface{}
		if json.Unmarshal(data, &cfg) != nil {
			continue
		}
		rootKey := "mcpServers"
		if client == "vscode" {
			rootKey = "servers"
		}
		if servers, ok := cfg[rootKey].(map[string]interface{}); ok && len(servers) > 0 {
			active = append(active, client)
		}
	}
	return active
}
