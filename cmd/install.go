package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arcmesh-labs/arcmesh-pm/internal/clients"
	"github.com/arcmesh-labs/arcmesh-pm/internal/config"
	installpkg "github.com/arcmesh-labs/arcmesh-pm/internal/install"
	"github.com/arcmesh-labs/arcmesh-pm/internal/registry"
)

var installCmd = &cobra.Command{
	Use:   "install <server>",
	Short: "Install an MCP server and register it with your AI client config.",
	Long: `Install an MCP server and register it with your AI client config.

Fetches the server manifest from the registry, runs the install, and writes
the config entry. Prompts for any required environment variables (e.g. API tokens).

Examples:

  apm install github
  apm install github --client cursor
  apm install github --name gh`,
	Args: cobra.ExactArgs(1),
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("name", "", "", "Custom name for this server entry in the client config.")
	installCmd.Flags().StringP("config-path", "", "", "Path to MCP config file.")
	installCmd.Flags().StringP("client", "", "", "Target client. Required when multiple clients are configured.")
}

func runInstall(cmd *cobra.Command, args []string) error {
	serverArg := args[0]
	alias, _ := cmd.Flags().GetString("name")
	configPathFlag, _ := cmd.Flags().GetString("config-path")
	clientFlag, _ := cmd.Flags().GetString("client")

	client, err := resolveClientActive(clientFlag)
	if err != nil {
		die("%v", err)
	}

	fmt.Printf("Fetching manifest for %s...\n", cyanStyle.Render(serverArg))
	manifest, source, err := registry.FetchManifestWithFallback(serverArg)
	if err != nil {
		die("%v", err)
	}
	_ = source

	if err := validateManifest(manifest, serverArg); err != nil {
		die("%v", err)
	}

	serverName := manifest.Name
	fmt.Printf("%s %s...\n", boldStyle.Render("Installing"), cyanStyle.Render(serverName))

	if manifest.Install.Type == "pip" {
		fmt.Println("Running pip install...")
	}
	if err := installpkg.RunInstall(manifest.Install); err != nil {
		die("%v", err)
	}

	// Prompt for env vars and build the config entry.
	configEntry, skippedRequired, err := buildConfigEntry(manifest, client)
	if err != nil {
		die("%v", err)
	}

	// Resolve config path and write.
	path, err := config.ResolveConfigPath(configPathFlag, client)
	if err != nil {
		die("%v", err)
	}

	clientLabel := clients.DisplayName(client) + " config"
	fmt.Printf("Writing %s...\n", clientLabel)

	if err := installpkg.WriteToConfig(configEntry, serverName, alias, path, client); err != nil {
		if err == installpkg.ErrAborted {
			return nil
		}
		die("%v", err)
	}

	configKey := serverName
	if alias != "" {
		configKey = alias
	}
	fmt.Printf("\n%s %s installed successfully.\n", successStyle.Render("✓"), boldStyle.Render(configKey))
	fmt.Printf("  Config written to: %s\n", dimStyle.Render(path))

	for _, key := range skippedRequired {
		fmt.Printf("\n%s %s not set. Run: apm set-env %s %s\n",
			warnStyle.Render("⚠"), key, serverArg, key)
	}

	fmt.Printf("\n%s\n", boldStyle.Render("Next steps:"))
	fmt.Printf("  1. Restart %s\n", clients.DisplayName(client))
	fmt.Printf("  2. Look for %s in the MCP tools panel\n", boldStyle.Render(serverName))
	return nil
}

func validateManifest(m *registry.Manifest, name string) error {
	var missing []string
	if m.Name == "" {
		missing = append(missing, "name")
	}
	if m.Install.Type == "" {
		missing = append(missing, "install.type")
	}
	if m.Config.Command == "" {
		missing = append(missing, "config.command")
	}
	if len(missing) > 0 {
		return fmt.Errorf("manifest for '%s' is invalid — missing fields: %s",
			boldStyle.Render(name), strings.Join(missing, ", "))
	}
	return nil
}

func buildConfigEntry(manifest *registry.Manifest, client string) (map[string]interface{}, []string, error) {
	args := make([]string, len(manifest.Config.Args))
	copy(args, manifest.Config.Args)

	var resolvedEnv map[string]string
	var skipped []string

	if len(manifest.Config.Env) > 0 {
		var err error
		resolvedEnv, skipped, err = installpkg.PromptEnv(manifest.Config.Env, client)
		if err != nil {
			return nil, nil, err
		}

		// Interpolate ${KEY} placeholders in args.
		for i, arg := range args {
			if strings.HasPrefix(arg, "${") && strings.HasSuffix(arg, "}") {
				key := arg[2 : len(arg)-1]
				if val, ok := resolvedEnv[key]; ok {
					args[i] = val
				}
			}
		}

		// Remove keys that were interpolated into args — they don't go in the config env.
		for _, arg := range manifest.Config.Args {
			if strings.HasPrefix(arg, "${") && strings.HasSuffix(arg, "}") {
				delete(resolvedEnv, arg[2:len(arg)-1])
			}
		}
	}

	entry := map[string]interface{}{
		"command": manifest.Config.Command,
		"args":    args,
	}
	if len(resolvedEnv) > 0 {
		envMap := make(map[string]interface{}, len(resolvedEnv))
		for k, v := range resolvedEnv {
			envMap[k] = v
		}
		entry["env"] = envMap
	}
	return entry, skipped, nil
}
