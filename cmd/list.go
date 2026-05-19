package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/registry"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available MCP servers in the ArcMesh registry.",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(_ *cobra.Command, _ []string) error {
	servers, err := registry.FetchIndex()
	if err != nil {
		die("%v", err)
	}
	if len(servers) == 0 {
		fmt.Println("Registry is empty.")
		return nil
	}

	t := newTable().Headers("NAME", "DESCRIPTION", "PUBLISHER", "VERIFIED")
	for _, s := range servers {
		t.Row(
			cyanStyle.Render(s.Name),
			truncate(s.Description, 72),
			dimStyle.Render(s.Publisher),
			verifiedBadge(s.Verified),
		)
	}
	printTable("Available MCP Servers", t)
	return nil
}
