package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List all installed MCP servers across configured clients.",
	Long: `List all installed MCP servers across configured clients.

Shows client, server name, install type, and environment variable status.
Use --client to filter to a single client.`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringP("config-path", "", "", "Path to MCP config file.")
	statusCmd.Flags().StringP("client", "", "", "Filter by client. Shows all clients if omitted.")
}

func runStatus(cmd *cobra.Command, _ []string) error {
	configPathFlag, _ := cmd.Flags().GetString("config-path")
	clientFlag, _ := cmd.Flags().GetString("client")

	if clientFlag != "" {
		return runStatusSingle(configPathFlag, clientFlag)
	}
	return runStatusAll()
}

func runStatusSingle(configPathFlag, client string) error {
	path, err := config.ResolveConfigPath(configPathFlag, client)
	if err != nil {
		die("%v", err)
	}

	cfg, err := config.ReadConfig(path)
	if err != nil {
		if os.IsNotExist(err) {
			die("Config not found at %s", path)
		}
		die("reading config: %v", err)
	}

	rootKey := config.RootKey(client)
	servers, _ := cfg[rootKey].(map[string]interface{})
	if len(servers) == 0 {
		fmt.Println("No servers installed.")
		return nil
	}

	t := newTable().Headers("NAME", "TYPE", "ENV")
	for name, raw := range servers {
		serverCfg, _ := raw.(map[string]interface{})
		t.Row(cyanStyle.Render(name), serverType(serverCfg), envCell(serverCfg))
	}
	printTable(fmt.Sprintf("Installed MCP servers (%d) — %s", len(servers), clients.DisplayName(client)), t)
	return nil
}

func runStatusAll() error {
	type row struct {
		client string
		name   string
		cfg    map[string]interface{}
	}
	var rows []row

	for _, client := range clients.ClientOrder {
		path, found := clients.FindConfig(client)
		if !found {
			continue
		}
		cfg, err := config.ReadConfig(path)
		if err != nil {
			continue
		}
		rootKey := config.RootKey(client)
		servers, _ := cfg[rootKey].(map[string]interface{})
		for name, raw := range servers {
			serverCfg, _ := raw.(map[string]interface{})
			rows = append(rows, row{client, name, serverCfg})
		}
	}

	if len(rows) == 0 {
		fmt.Println("No servers installed.")
		return nil
	}

	t := newTable().Headers("CLIENT", "NAME", "TYPE", "ENV")
	for _, r := range rows {
		t.Row(r.client, cyanStyle.Render(r.name), serverType(r.cfg), envCell(r.cfg))
	}
	printTable(fmt.Sprintf("Installed MCP servers (%d)", len(rows)), t)
	return nil
}
