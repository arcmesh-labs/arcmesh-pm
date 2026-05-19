package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check MCP client config health.",
	Long: `Check MCP client config health.

Validates config file exists, is valid JSON, has the correct root key,
and that each server entry has a command field. Warns on missing or empty
environment variables.`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().StringP("config-path", "", "", "Path to MCP config file.")
	doctorCmd.Flags().StringP("client", "", "", "Target client. Checks all clients if omitted.")
}

func runDoctor(cmd *cobra.Command, _ []string) error {
	configPathFlag, _ := cmd.Flags().GetString("config-path")
	clientFlag, _ := cmd.Flags().GetString("client")

	if clientFlag != "" {
		path, err := config.ResolveConfigPath(configPathFlag, clientFlag)
		if err != nil {
			die("%v", err)
		}
		if !checkClient(clientFlag, path) {
			os.Exit(1)
		}
		return nil
	}

	ok := true
	for _, client := range clients.ClientOrder {
		fmt.Printf("\n%s\n", dimStyle.Render("── "+client+" ──"))
		path, found := clients.FindConfig(client)
		if !found {
			fmt.Printf("%s\n", dimStyle.Render(fmt.Sprintf("No config found for %s — skipping", client)))
			continue
		}
		if !checkClient(client, path) {
			ok = false
		}
	}
	if !ok {
		os.Exit(1)
	}
	return nil
}

func checkClient(client, configPath string) bool {
	ok := true
	clientLabel := clients.DisplayName(client) + " config"
	rootKey := config.RootKey(client)

	if _, err := os.Stat(configPath); err != nil {
		fmt.Printf("%s Config not found at %s\n", errStyle.Render("Error:"), configPath)
		return false
	}
	fmt.Printf("%s %s found: %s\n", successStyle.Render("✓"), clientLabel, dimStyle.Render(configPath))

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("%s Could not read config: %v\n", errStyle.Render("✗"), err)
		return false
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("%s Config is not valid JSON: %v\n", errStyle.Render("✗"), err)
		return false
	}
	fmt.Printf("%s Config is valid JSON\n", successStyle.Render("✓"))

	if _, hasKey := cfg[rootKey]; !hasKey {
		fmt.Printf("%s Config is missing '%s' key\n", errStyle.Render("✗"), rootKey)
		ok = false
	} else {
		fmt.Printf("%s '%s' key exists\n", successStyle.Render("✓"), rootKey)
		if servers, ok2 := cfg[rootKey].(map[string]interface{}); ok2 {
			for name, raw := range servers {
				serverCfg, _ := raw.(map[string]interface{})
				if _, hasCmd := serverCfg["command"]; !hasCmd {
					fmt.Printf("%s Server '%s' is missing 'command' field\n", errStyle.Render("✗"), boldStyle.Render(name))
					ok = false
				} else {
					fmt.Printf("%s Server '%s' has 'command' field\n", successStyle.Render("✓"), boldStyle.Render(name))
				}
			}
		}
	}

	if servers, ok2 := cfg[rootKey].(map[string]interface{}); ok2 {
		for name, raw := range servers {
			serverCfg, _ := raw.(map[string]interface{})
			env, _ := serverCfg["env"].(map[string]interface{})
			for key, val := range env {
				valStr, _ := val.(string)
				if strings.TrimSpace(valStr) == "" {
					fmt.Printf("%s %s: %s is not set. Run: apm set-env %s %s\n",
						warnStyle.Render("⚠"), name, key, name, key)
				}
			}
		}
	}

	return ok
}
