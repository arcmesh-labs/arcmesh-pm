package cmd

import (
	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
)

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "List all supported MCP clients and whether they are configured on this system.",
	Long: `List all supported MCP clients and whether they are configured on this system.

Shows config file path and status (found / not found) for each client:
claude-desktop, vscode, cursor, windsurf.`,
	RunE: runClients,
}

func init() {
	rootCmd.AddCommand(clientsCmd)
}

func runClients(_ *cobra.Command, _ []string) error {
	t := newTable().Headers("CLIENT", "STATUS", "CONFIG PATH")

	for _, client := range clients.ClientOrder {
		path, found := clients.FindConfig(client)
		var status, configPath string
		if found {
			status = successStyle.Render("✓ configured")
			configPath = path
		} else {
			status = dimStyle.Render("✗ not configured")
			configPath = dimStyle.Render("—")
		}
		t.Row(cyanStyle.Render(client), status, configPath)
	}

	printTable("MCP clients", t)
	return nil
}
