package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <server>",
	Short: "Remove an MCP server from your client config.",
	Long: `Remove an MCP server from your client config.

Examples:

  apm uninstall github
  apm uninstall github --client cursor`,
	Args: cobra.ExactArgs(1),
	RunE: runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringP("config-path", "", "", "Path to MCP config file.")
	uninstallCmd.Flags().StringP("client", "", "", "Target client. Required when multiple clients are configured.")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	name := args[0]
	configPathFlag, _ := cmd.Flags().GetString("config-path")
	clientFlag, _ := cmd.Flags().GetString("client")

	client, err := resolveClientActive(clientFlag)
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
	if _, exists := servers[name]; !exists {
		die("'%s' not found in %s", boldStyle.Render(name), rootKey)
	}

	delete(servers, name)
	cfg[rootKey] = servers

	if err := config.WriteConfig(path, cfg); err != nil {
		die("writing config: %v", err)
	}

	clientLabel := clients.DisplayName(client) + " config"
	fmt.Printf("%s %s removed from %s.\n", successStyle.Render("✓"), boldStyle.Render(name), clientLabel)
	fmt.Printf("\nRestart %s for the change to take effect.\n", clients.DisplayName(client))
	return nil
}
