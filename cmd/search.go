package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/registry"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for MCP servers across ArcMesh and official MCP Registry.",
	Long: `Search for MCP servers across ArcMesh and official MCP Registry.

Examples:

  apm search github
  apm search database`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(_ *cobra.Command, args []string) error {
	query := args[0]
	results, err := registry.SearchWithFallback(query)
	if err != nil {
		die("%v", err)
	}
	if len(results) == 0 {
		fmt.Printf("No servers found matching '%s'.\n", boldStyle.Render(query))
		return nil
	}

	t := newTable().Headers("NAME", "DESCRIPTION", "PUBLISHER", "VERIFIED", "SOURCE")
	for _, r := range results {
		t.Row(
			cyanStyle.Render(r.Name),
			truncate(r.Description, 72),
			dimStyle.Render(r.Publisher),
			verifiedBadge(r.Verified),
			sourceBadge(r.Source),
		)
	}
	printTable(fmt.Sprintf("Results for '%s'", query), t)
	return nil
}
