package config

import (
	"encoding/json"
	"os"
)

func ReadConfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func WriteConfig(path string, cfg map[string]interface{}) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RootKey returns the JSON key used for MCP server entries.
// VS Code uses "servers"; all other clients use "mcpServers".
func RootKey(client string) string {
	if client == "vscode" {
		return "servers"
	}
	return "mcpServers"
}
