package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
	"github.com/arcmesh-labs/arcmesh-pm/internal/install"
)

var setEnvCmd = &cobra.Command{
	Use:   "set-env <server> <key>",
	Short: "Update an environment variable for an installed MCP server.",
	Long: `Update an environment variable for an installed MCP server.

Use this to set or rotate API tokens without reinstalling the server.

Examples:

  apm set-env github GITHUB_PERSONAL_ACCESS_TOKEN
  apm set-env github GITHUB_PERSONAL_ACCESS_TOKEN --client cursor`,
	Args: cobra.ExactArgs(2),
	RunE: runSetEnv,
}

func init() {
	rootCmd.AddCommand(setEnvCmd)
	setEnvCmd.Flags().StringP("config-path", "", "", "Path to MCP config file.")
	setEnvCmd.Flags().StringP("client", "", "", "Target client. Required when multiple clients are configured.")
}

func runSetEnv(cmd *cobra.Command, args []string) error {
	server, key := args[0], args[1]
	configPathFlag, _ := cmd.Flags().GetString("config-path")
	clientFlag, _ := cmd.Flags().GetString("client")

	client, err := resolveClientConfigured(clientFlag)
	if err != nil {
		die("%v", err)
	}

	path, err := config.ResolveConfigPath(configPathFlag, client)
	if err != nil {
		die("%v", err)
	}

	cfg, err := config.ReadConfig(path)
	if err != nil {
		die("reading config: %v", err)
	}

	rootKey := config.RootKey(client)
	servers, _ := cfg[rootKey].(map[string]interface{})
	serverEntry, ok := servers[server].(map[string]interface{})
	if !ok {
		die("'%s' not found in %s", boldStyle.Render(server), rootKey)
	}

	envRaw, hasEnv := serverEntry["env"]
	if !hasEnv {
		die("'%s' has no env configuration", boldStyle.Render(server))
	}
	env, _ := envRaw.(map[string]interface{})
	if _, hasKey := env[key]; !hasKey {
		die("'%s' not found in env for '%s'", boldStyle.Render(key), boldStyle.Render(server))
	}

	current, _ := env[key].(string)
	if strings.TrimSpace(current) != "" {
		masked := current[:min(4, len(current))] + strings.Repeat("*", max(10, len(current)-4))
		fmt.Printf("%s %s is already set (current value: %s).\n\n", warnStyle.Render("⚠"), key, masked)
		fmt.Println("  Overwrite?")
		fmt.Println("    1. Yes")
		fmt.Println("    2. No")
		fmt.Print("Choose [1/2]: ")
		var answer string
		fmt.Scanln(&answer)
		if strings.TrimSpace(answer) != "1" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	value, err := install.ReadSecret(key)
	if err != nil {
		die("reading input: %v", err)
	}
	if strings.TrimSpace(value) == "" {
		die("Value cannot be empty.")
	}

	env[key] = value
	if err := config.WriteConfig(path, cfg); err != nil {
		die("writing config: %v", err)
	}

	fmt.Printf("%s %s updated for %s.\n",
		successStyle.Render("✓"), boldStyle.Render(key), boldStyle.Render(server))
	fmt.Printf("\nRestart %s for the change to take effect.\n", clients.DisplayName(client))
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
