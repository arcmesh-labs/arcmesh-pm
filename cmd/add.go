package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
	installpkg "github.com/arcmesh-labs/arcmesh-pm/internal/install"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Scaffold a local MCP server for the current directory and register it with your AI client config.",
	Long: `Scaffold a local MCP server for the current directory and register it with your AI client config.

Creates a .mcp/ directory in cwd, writes a fastmcp server.py with read_file,
list_directory, and search_content tools, writes a .mcp/config.json, and
registers the server in the AI client config.

Examples:

  apm add
  apm add --name my-server
  apm add --client cursor`,
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("name", "", "", "Server name in the client config. Default: current directory name.")
	addCmd.Flags().StringP("client", "", "", "Target client. Required when multiple clients are configured.")
	addCmd.Flags().StringP("config-path", "", "", "Path to MCP config file.")
}

func runAdd(cmd *cobra.Command, _ []string) error {
	nameFlag, _ := cmd.Flags().GetString("name")
	clientFlag, _ := cmd.Flags().GetString("client")
	configPathFlag, _ := cmd.Flags().GetString("config-path")

	// ── cwd + name ────────────────────────────────────────────────────────────

	cwd, err := os.Getwd()
	if err != nil {
		die("could not determine working directory: %v", err)
	}

	name := nameFlag
	if name == "" {
		name = filepath.Base(cwd)
	}

	// ── .mcp/ directory ───────────────────────────────────────────────────────

	mcpDir := filepath.Join(cwd, ".mcp")
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		die("could not create .mcp directory: %v", err)
	}

	// ── .mcp/server.py ────────────────────────────────────────────────────────

	serverPyPath := filepath.Join(mcpDir, "server.py")
	if _, err := os.Stat(serverPyPath); err == nil {
		fmt.Printf("\n%s  '.mcp/server.py' already exists.\n", warnStyle.Render("⚠"))
		fmt.Println("  1. Yes — overwrite")
		fmt.Println("  2. No  — abort")
		fmt.Print("Choose [1/2]: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) != "1" {
				fmt.Println("Aborted.")
				return nil
			}
		}
	}

	if err := os.WriteFile(serverPyPath, []byte(buildServerPy(cwd, name)), 0644); err != nil {
		die("could not write .mcp/server.py: %v", err)
	}

	// ── .mcp/config.json ──────────────────────────────────────────────────────

	mcpConfig := map[string]interface{}{
		"version": "1.0",
		"servers": map[string]interface{}{
			name: map[string]interface{}{
				"command": "python",
				"args":    []string{serverPyPath},
			},
		},
	}
	configJSON, err := json.MarshalIndent(mcpConfig, "", "  ")
	if err != nil {
		die("could not marshal .mcp/config.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(mcpDir, "config.json"), configJSON, 0644); err != nil {
		die("could not write .mcp/config.json: %v", err)
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

	// ── register in AI client config ──────────────────────────────────────────

	python, err := installpkg.FindPython()
	if err != nil {
		die("%v", err)
	}

	var entry map[string]interface{}
	if clients.IsWSL() {
		entry = map[string]interface{}{
			"command": "wsl.exe",
			"args":    []string{"-d", clients.WSLDistro(), "-e", "bash", "-lc", fmt.Sprintf("%s \"%s\"", python, serverPyPath)},
		}
	} else {
		entry = map[string]interface{}{
			"command": python,
			"args":    []string{serverPyPath},
		}
	}

	if err := installpkg.WriteToConfig(entry, name, "", configPath, client); err != nil {
		if err == installpkg.ErrAborted {
			return nil
		}
		die("%v", err)
	}

	fmt.Printf("\n%s %s added successfully.\n", successStyle.Render("✓"), boldStyle.Render(name))
	fmt.Printf("  Config written to: %s\n", dimStyle.Render(configPath))
	fmt.Printf("  Server script:     %s\n", dimStyle.Render(serverPyPath))
	fmt.Printf("\n%s\n", boldStyle.Render("Next steps:"))
	fmt.Printf("  1. Restart %s\n", clients.DisplayName(client))
	fmt.Printf("  2. Look for %s in the MCP tools panel\n", boldStyle.Render(name))
	return nil
}

func buildServerPy(baseDirPath, repoName string) string {
	return fmt.Sprintf(`from pathlib import Path

from mcp.server.fastmcp import FastMCP

BASE_DIR = Path(%q)

mcp = FastMCP(%q)


def _check_path(path: str) -> tuple[Path, str | None]:
    if path in ("", "/"):
        return BASE_DIR.resolve(), None
    target = (BASE_DIR / path).resolve()
    if not str(target).startswith(str(BASE_DIR.resolve())):
        return target, "Error: path is outside BASE_DIR"
    return target, None


@mcp.tool()
def read_file(path: str) -> str:
    """Read and return the contents of a file."""
    target, err = _check_path(path)
    if err:
        return err
    if not target.is_file():
        return f"Error: not a file: {path}"
    return target.read_text()


@mcp.tool()
def list_directory(path: str) -> str:
    """List files and directories in a directory."""
    target, err = _check_path(path)
    if err:
        return err
    if not target.is_dir():
        return f"Error: not a directory: {path}"
    return "\n".join(entry.name for entry in sorted(target.iterdir()))


@mcp.tool()
def search_content(path: str, query: str) -> list[str]:
    """Recursively search for files containing query, returning matching paths and lines."""
    target, err = _check_path(path)
    if err:
        return [err]
    if not target.is_dir():
        return [f"Error: not a directory: {path}"]
    results = []
    for file in target.rglob("*"):
        if file.is_file():
            try:
                for lineno, line in enumerate(file.read_text().splitlines(), 1):
                    if query in line:
                        results.append(f"{file}:{lineno}: {line}")
            except (OSError, UnicodeDecodeError):
                pass
    return results


if __name__ == "__main__":
    mcp.run()
`, baseDirPath, repoName)
}
