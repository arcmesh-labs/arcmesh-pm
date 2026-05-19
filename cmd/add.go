package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
	installpkg "github.com/arcmesh-labs/arcmesh-pm/internal/install"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Register the MCP server in the current directory with your AI client config.",
	Long: `Register the MCP server in the current directory with your AI client config.

Expects server.py to exist in the current directory.
In WSL the entry is wrapped with wsl.exe so Claude Desktop (a Windows process)
can launch the server.

Examples:

  apm add
  apm add --name my-server
  apm add --venv ~/.venv --client cursor`,
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("name", "", "", "Server name in the client config. Default: current directory name.")
	addCmd.Flags().StringP("venv", "", "", "Path to Python venv. Default: ~/tools/mcp/venv.")
	addCmd.Flags().StringP("client", "", "", "Target client. Required when multiple clients are configured.")
	addCmd.Flags().StringP("config-path", "", "", "Path to MCP config file.")
}

func runAdd(cmd *cobra.Command, _ []string) error {
	nameFlag, _ := cmd.Flags().GetString("name")
	venvFlag, _ := cmd.Flags().GetString("venv")
	clientFlag, _ := cmd.Flags().GetString("client")
	configPathFlag, _ := cmd.Flags().GetString("config-path")

	// ── cwd + server.py ───────────────────────────────────────────────────────

	cwd, err := os.Getwd()
	if err != nil {
		die("could not determine working directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cwd, "server.py")); err != nil {
		die("server.py not found in %s", cwd)
	}

	// ── name ──────────────────────────────────────────────────────────────────

	name := nameFlag
	if name == "" {
		name = filepath.Base(cwd)
	}

	// ── client + config path ──────────────────────────────────────────────────

	client, err := resolveClientActive(clientFlag)
	if err != nil {
		die("%v", err)
	}

	configPath, err := config.ResolveConfigPath(configPathFlag, client)
	if err != nil {
		die("%v", err)
	}

	// ── venv ──────────────────────────────────────────────────────────────────

	venv := venvFlag
	if venv == "" {
		home, _ := os.UserHomeDir()
		candidate := filepath.Join(home, "tools", "mcp", "venv")
		if _, err := os.Stat(candidate); err == nil {
			venv = candidate
		} else {
			die("could not find a Python venv. Pass --venv <path>")
		}
	}

	// ── config entry ──────────────────────────────────────────────────────────

	scriptPath := filepath.Join(cwd, "server.py")
	var entry map[string]interface{}

	if clients.IsWSL() {
		shellStr := fmt.Sprintf("source %s/bin/activate && python %s", venv, scriptPath)
		entry = map[string]interface{}{
			"command": "wsl.exe",
			"args": []string{
				"-d", clients.WSLDistro(),
				"-e", "bash", "-lc",
				shellStr,
			},
		}
	} else {
		entry = map[string]interface{}{
			"command": "python",
			"args":    []string{scriptPath},
		}
	}

	// ── write ─────────────────────────────────────────────────────────────────

	if err := installpkg.WriteToConfig(entry, name, "", configPath, client); err != nil {
		if err == installpkg.ErrAborted {
			return nil
		}
		die("%v", err)
	}

	fmt.Printf("\n%s %s added successfully.\n", successStyle.Render("✓"), boldStyle.Render(name))
	fmt.Printf("  Config written to: %s\n", dimStyle.Render(configPath))
	fmt.Printf("\n%s\n", boldStyle.Render("Next steps:"))
	fmt.Printf("  1. Restart %s\n", clients.DisplayName(client))
	fmt.Printf("  2. Look for %s in the MCP tools panel\n", boldStyle.Render(name))
	return nil
}
