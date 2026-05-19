package install

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
)

// ErrAborted is returned when the user declines an overwrite prompt.
// Callers should exit with code 0, not 1.
var ErrAborted = errors.New("aborted")

var writeWarnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
var writeErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)

// WriteToConfig merges configEntry into the client config file at configPath.
// Prompts the user before overwriting an existing server entry.
func WriteToConfig(configEntry map[string]interface{}, serverName, alias, configPath, client string) error {
	if configPath == "" {
		return fmt.Errorf("configPath is required")
	}

	rootKey := config.RootKey(client)
	configKey := serverName
	if alias != "" {
		configKey = alias
	}

	// Read or initialise the config.
	existing, err := config.ReadConfig(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("reading config: %w", err)
		}
		existing = map[string]interface{}{}
	}

	// Validate that the root key, if present, is an object.
	if raw, ok := existing[rootKey]; ok {
		if _, isMap := raw.(map[string]interface{}); !isMap {
			return fmt.Errorf("config key '%s' is not an object. Fix your config file before installing", rootKey)
		}
	}

	// Check for an existing entry and prompt for overwrite.
	if servers, ok := existing[rootKey].(map[string]interface{}); ok {
		if _, exists := servers[configKey]; exists {
			fmt.Printf("\n%s  '%s' is already configured.\n", writeWarnStyle.Render("⚠"), configKey)
			fmt.Println("  1. Yes — overwrite")
			fmt.Println("  2. No  — abort")
			fmt.Print("Choose [1/2]: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				if strings.TrimSpace(scanner.Text()) != "1" {
					fmt.Println("Aborted.")
					return ErrAborted
				}
			}
		}
	}

	// Ensure parent directory exists.
	if err := os.MkdirAll(dirOf(configPath), 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Merge entry.
	servers, _ := existing[rootKey].(map[string]interface{})
	if servers == nil {
		servers = map[string]interface{}{}
	}
	servers[configKey] = configEntry
	existing[rootKey] = servers

	if err := config.WriteConfig(configPath, existing); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied writing to %s. Check file permissions", configPath)
		}
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

func dirOf(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}
