package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags "-X github.com/arcmesh-labs/arcmesh-pm/cmd.version=<tag>"
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "apm",
	Short:   "arcmesh-pm (apm) — the package manager for MCP servers.",
	Version: version,
}

var _sections = []struct {
	title string
	cmds  []string
}{
	{"Registry", []string{"search", "list"}},
	{"Install", []string{"install", "uninstall", "set-env"}},
	{"Config", []string{"config edit", "config path"}},
	{"Clients", []string{"clients", "status", "doctor"}},
}

func init() {
	rootCmd.SetUsageFunc(sectionedUsage)
}

func sectionedUsage(cmd *cobra.Command) error {
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "Usage:\n  %s [command]\n\n", cmd.Name())
	fmt.Fprintf(w, "  %s\n\n", cmd.Short)

	for _, section := range _sections {
		fmt.Fprintf(w, "%s:\n", section.title)
		for _, path := range section.cmds {
			c := findNestedCmd(cmd, path)
			if c == nil || c.Hidden {
				continue
			}
			fmt.Fprintf(w, "  %-16s %s\n", path, c.Short)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "Flags:\n%s\n", cmd.Flags().FlagUsages())
	fmt.Fprintf(w, "Use \"%s [command] --help\" for more information about a command.\n", cmd.Name())
	return nil
}

func findNestedCmd(root *cobra.Command, path string) *cobra.Command {
	parts := strings.SplitN(path, " ", 2)
	for _, c := range root.Commands() {
		if c.Name() != parts[0] {
			continue
		}
		if len(parts) == 1 {
			return c
		}
		for _, sub := range c.Commands() {
			if sub.Name() == parts[1] {
				return sub
			}
		}
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
