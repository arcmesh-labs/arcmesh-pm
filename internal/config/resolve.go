package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
)

var SettingsPath = filepath.Join(mustHome(), ".arcmesh-pm", "config.json")

var errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)

func mustHome() string {
	h, _ := os.UserHomeDir()
	return h
}

// ResolveConfigPath returns the config file path for a client, with fallback to
// saved settings, auto-detection, and interactive prompt (faithful port of resolve_config_path).
func ResolveConfigPath(override string, client string) (string, error) {
	if override != "" {
		return override, nil
	}

	// Check ~/.arcmesh-pm/config.json for a saved path
	if data, err := os.ReadFile(SettingsPath); err == nil {
		var saved map[string]interface{}
		if json.Unmarshal(data, &saved) == nil {
			switch v := saved["config_path"].(type) {
			case map[string]interface{}:
				if p, ok := v[client].(string); ok {
					return p, nil
				}
			case string:
				if client == "claude-desktop" {
					return v, nil
				}
			}
		}
	}

	// Auto-detect
	if found, ok := clients.FindConfig(client); ok {
		return found, nil
	}

	// Interactive fallback
	fmt.Fprintf(os.Stderr, "%s %s config not found. Checked:\n",
		errStyle.Render("Error:"), clients.DisplayName(client))
	for _, p := range clients.CheckedPaths(client) {
		prefix := ""
		if strings.Contains(p, "*") {
			prefix = "(WSL) "
		}
		fmt.Fprintf(os.Stderr, "  - %s%s\n", prefix, p)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Enter path to %s: ", clients.ConfigFilename(client))
		if !scanner.Scan() {
			return "", fmt.Errorf("no input")
		}
		candidate := strings.TrimSpace(scanner.Text())
		if _, err := os.Stat(candidate); err == nil {
			if err := saveSettingsPath(client, candidate); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not save config path: %v\n", err)
			}
			return candidate, nil
		}
		fmt.Fprintf(os.Stderr, "%s File not found: %s\n", errStyle.Render("Error:"), candidate)
	}
}

func saveSettingsPath(client, path string) error {
	existing := map[string]interface{}{}
	if data, err := os.ReadFile(SettingsPath); err == nil {
		json.Unmarshal(data, &existing) //nolint:errcheck — parse errors mean we start fresh
	}

	var paths map[string]interface{}
	switch v := existing["config_path"].(type) {
	case map[string]interface{}:
		paths = v
	case string:
		paths = map[string]interface{}{"claude-desktop": v}
	default:
		paths = map[string]interface{}{}
	}
	paths[client] = path
	existing["config_path"] = paths

	if err := os.MkdirAll(filepath.Dir(SettingsPath), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(existing)
	if err != nil {
		return err
	}
	return os.WriteFile(SettingsPath, data, 0644)
}
