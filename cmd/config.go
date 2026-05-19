package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage MCP client config.",
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open your MCP client config file in $EDITOR.",
	Long: `Open your MCP client config file in $EDITOR.

If $EDITOR is not set, prints the file path for manual editing.`,
	RunE: runConfigEdit,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the path to your MCP client config file.",
	RunE:  runConfigPath,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configPathCmd)

	for _, sub := range []*cobra.Command{configEditCmd, configPathCmd} {
		sub.Flags().StringP("config-path", "", "", "Path to MCP config file.")
		sub.Flags().StringP("client", "", "", "Target client. Required when multiple clients are configured.")
	}
}

func runConfigEdit(cmd *cobra.Command, _ []string) error {
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

	editor := os.Getenv("EDITOR")
	if editor == "" {
		var tip string
		switch runtime.GOOS {
		case "windows":
			tip = "Set EDITOR=notepad in System Environment Variables\n    (System Properties → Advanced → Environment Variables)"
		case "darwin":
			tip = "echo 'export EDITOR=nano' >> ~/.zshrc && source ~/.zshrc"
		default:
			tip = "echo 'export EDITOR=nano' >> ~/.bashrc && source ~/.bashrc"
		}
		fmt.Printf("$EDITOR is not set. Edit the config manually:\n  %s\n\nTip: Set your preferred editor:\n  %s\n", path, tip)
		return nil
	}

	c := exec.Command(editor, path)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func runConfigPath(cmd *cobra.Command, _ []string) error {
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

	fmt.Println(path)
	return nil
}
